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
	"io"
	"net/http"
	"os"

	"github.com/AnatolySnegovskiy/metric/internal/services/dto"
	"github.com/mailru/easyjson"
)

func (a *Agent) sendMetricsPeriodically(ctx context.Context) error {
	metricDtoCollection := dto.MetricsCollection{}

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

			if storageType == "counter" {
				iv := int64(metric)
				newIv := iv
				metricDto.Delta = &newIv
			} else {
				newMetric := metric
				metricDto.Value = &newMetric
			}

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

	if a.shaKey != "" {
		hash := hmac.New(sha256.New, []byte(a.shaKey))
		hash.Write(buf.Bytes())
		req.Header.Set("HashSHA256", fmt.Sprintf("%x", hash.Sum(nil)))
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
