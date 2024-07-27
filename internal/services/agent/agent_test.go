package agent

import (
	"bou.ke/monkey"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/mocks"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	options := Options{
		Storage:        storages.NewMemStorage(),
		PollInterval:   10,
		ReportInterval: 20,
		SendAddr:       "example.com:1234",
	}

	agent := New(options)
	assert.NotNil(t, agent.storage, "storage should not be nil")
	assert.NotNil(t, agent.pollInterval, "pollInterval should not be nil")
	assert.NotNil(t, agent.reportInterval, "reportInterval should not be nil")
	assert.Equal(t, "example.com:1234", agent.sendAddr, "send address should be example.com:1234")
}

func TestAgent(t *testing.T) {
	testCases := []struct {
		name          string
		statusCode    int
		doReturnError error
		expectedErr   bool
		mockStorage   func() *storages.MemStorage
	}{
		{"SuccessNil", http.StatusOK, nil, false, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			mockStorage.AddMetric("nil", nil)
			return mockStorage
		}},
		{"Success", http.StatusOK, nil, false, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			return mockStorage
		}},
		{"ErrorPoll", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			return mockStorage
		}},
		{"ErrorPollCounter", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			return mockStorage
		}},
		{"ErrorPollGauge", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			return mockStorage
		}},

		{"ErrorPollPollCount", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "PollCount", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("PollCount"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", mockEntity)
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			return mockStorage
		}},
		{"ErrorTotalMemory", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "TotalMemory", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("TotalMemory"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorFreeMemory", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "FreeMemory", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("FreeMemory"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()

			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorCPUutilization1", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("CPUutilization1"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().Process(gomock.Any(), "CPUutilization1", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorPollRandomValue", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "RandomValue", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("RandomValue"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorPollRuntimeEntityArray", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "Alloc", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("Alloc"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
	}
	ctrl := gomock.NewController(t)
	privateKey, publicKey := generateRSAKeys()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := mocks.NewMockHTTPClient(ctrl)
			resp := http.Response{StatusCode: tc.statusCode, Body: http.NoBody}
			httpClient.EXPECT().Do(gomock.Any()).Return(&resp, tc.doReturnError).AnyTimes()

			a := Agent{
				storage:        tc.mockStorage(),
				sendAddr:       "testAddr",
				client:         httpClient,
				pollInterval:   1,
				reportInterval: 3,
				maxRetries:     2,
				shaKey:         "testKey",
				cryptoKey:      publicKey,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
			defer cancel()

			err := a.Run(ctx)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	os.Remove(privateKey)
	os.Remove(publicKey)
	defer ctrl.Finish()
}

func TestAgentReportTickerEmpty(t *testing.T) {
	t.Run("EmptyStorage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mocks.NewMockHTTPClient(ctrl)
		resp := &http.Response{}
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).AnyTimes()
		a := &Agent{
			storage:        storages.NewMemStorage(),
			sendAddr:       "testAddr",
			client:         httpClient,
			pollInterval:   1,
			reportInterval: 1,
			shaKey:         "testKey",
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := a.Run(ctx)

		assert.Error(t, err)
	})
}

func TestAgentErrorCrypto(t *testing.T) {
	mockStorage := storages.NewMemStorage()
	mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
	mockStorage.AddMetric("counter", metrics.NewCounter(nil))
	mockStorage.AddMetric("nil", nil)
	t.Run("EmptyStorage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mocks.NewMockHTTPClient(ctrl)
		resp := &http.Response{}
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).AnyTimes()
		a := Agent{
			storage:        mockStorage,
			sendAddr:       "testAddr",
			client:         httpClient,
			pollInterval:   1,
			reportInterval: 3,
			maxRetries:     2,
			shaKey:         "testKey",
			cryptoKey:      "testKey",
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := a.sendMetricsPeriodically(ctx)

		assert.Error(t, err)
	})
}

func TestEncryptMessage(t *testing.T) {
	mockPublicKeyData := []byte("mocked_public_key_data")
	mockPublicKeyBlock := &pem.Block{Bytes: []byte("mocked_public_key_block")}
	mockRSAPublicKey, _ := x509.ParsePKIXPublicKey([]byte("mocked_rsa_public_key"))
	mockEncryptedMessage := []byte("mocked_encrypted_message")

	monkey.Patch(os.ReadFile, func(filename string) ([]byte, error) {
		return mockPublicKeyData, nil
	})
	monkey.Patch(pem.Decode, func(data []byte) (*pem.Block, []byte) {
		return mockPublicKeyBlock, nil
	})
	monkey.Patch(x509.ParsePKIXPublicKey, func(data []byte) (interface{}, error) {
		return mockRSAPublicKey, nil
	})
	monkey.Patch(rsa.EncryptPKCS1v15, func(rand io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error) {
		return mockEncryptedMessage, nil
	})
	defer monkey.UnpatchAll()

	// Test cases
	t.Run("Successful Encryption", func(t *testing.T) {
		encryptedMessage, err := encryptMessage([]byte("test_message"), "mocked_public_key_path")
		assert.NoError(t, err)
		assert.Equal(t, mockEncryptedMessage, encryptedMessage)
	})

	t.Run("Error Reading Public Key File", func(t *testing.T) {
		monkey.Patch(os.ReadFile, func(filename string) ([]byte, error) {
			return nil, errors.New("mocked read file error")
		})
		defer monkey.Unpatch(os.ReadFile)

		encryptedMessage, err := encryptMessage([]byte("test_message"), "invalid_public_key_path")
		assert.Error(t, err)
		assert.Nil(t, encryptedMessage)
	})

	monkey.Patch(x509.ParsePKIXPublicKey, func(data []byte) (interface{}, error) {
		return nil, errors.New("mocked parse error")
	})

	t.Run("Error Parsing Public Key", func(t *testing.T) {
		encryptedMessage, err := encryptMessage([]byte("test_message"), "mocked_public_key_path")
		assert.Error(t, err)
		assert.Nil(t, encryptedMessage)
	})
}

func generateRSAKeys() (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	privateKeyFile, err := os.Create("private_key.pem")
	if err != nil {
		panic(err)
	}

	defer privateKeyFile.Close()
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPEM := &pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}
	publicKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		panic(err)
	}
	defer publicKeyFile.Close()
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		panic(err)
	}

	privateKeyPath, _ := filepath.Abs(privateKeyFile.Name())
	publicKeyPath, _ := filepath.Abs(publicKeyFile.Name())
	return privateKeyPath, publicKeyPath
}
