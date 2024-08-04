package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	pb "github.com/AnatolySnegovskiy/metric/internal/services/grpc"
	"io"
	"net/http"
	"os"

	"github.com/AnatolySnegovskiy/metric/internal/services/dto"
	"github.com/mailru/easyjson"
)

func (a *Agent) sendMetricsPeriodically(ctx context.Context) error {
	metricDtoCollection := dto.MetricsCollection{}
	var metricsGrpcCollection []*pb.MetricRequest

	for storageType, storage := range a.storage.GetList() {
		if storage == nil {
			continue
		}

		list, _ := storage.GetList(ctx)

		for metricName, metric := range list {

			metricDto := dto.Metrics{
				ID:    metricName,
				MType: storageType,
			}
			metricsGrpc := &pb.MetricRequest{
				Id:   metricName,
				Type: storageType,
			}

			if storageType == "counter" {
				iv := int64(metric)
				newIv := iv
				metricDto.Delta = &newIv
				metricsGrpc.Delta = iv
			} else {
				newMetric := metric
				metricDto.Value = &newMetric
				metricsGrpc.Value = float32(metric)
			}
			metricsGrpcCollection = append(metricsGrpcCollection, metricsGrpc)
			metricDtoCollection = append(metricDtoCollection, metricDto)
		}
	}
	body, _ := easyjson.Marshal(metricDtoCollection)
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, _ = gw.Write(body)
	_ = gw.Close()
	body = buf.Bytes()

	if a.cryptoKey != "" {
		bodyEncrypted, err := encryptMessage(buf.Bytes(), a.cryptoKey)

		if err != nil {
			return err
		}

		body = bodyEncrypted
	}

	url := fmt.Sprintf("http://%s/updates/", a.sendAddr)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", "127.0.0.1")

	if a.shaKey != "" {
		hash := hmac.New(sha256.New, []byte(a.shaKey))
		hash.Write(buf.Bytes())
		req.Header.Set("HashSHA256", fmt.Sprintf("%x", hash.Sum(nil)))
	}

	_, err := a.grpcClient.UpdateMany(ctx, &pb.MetricRequestMany{Requests: metricsGrpcCollection})

	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)

	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(body))
		}
	}

	return err
}

func encryptMessage(message []byte, publicKeyPath string) ([]byte, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyData)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, _ := publicKey.(*rsa.PublicKey)
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, message)

	return encryptedMessage, err
}
