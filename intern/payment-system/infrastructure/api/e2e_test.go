package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"final-project/intern/payment-system/application"
	"final-project/intern/payment-system/config"
	"final-project/intern/payment-system/domain"
	"final-project/intern/payment-system/infrastructure/api"
	"final-project/intern/payment-system/infrastructure/render"
	"final-project/intern/payment-system/infrastructure/storage"
	"final-project/intern/payment-system/infrastructure/webhook"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type Tester struct {
	server *httptest.Server
	client *http.Client
	config *config.Config

	mockServer      *httptest.Server
	lastRequestData []byte
	requestRecieved bool
}

func (tester *Tester) CreateMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tester.lastRequestData, _ = io.ReadAll(r.Body)
		tester.requestRecieved = true
	}))
}

func (tester *Tester) New(config *config.Config) error {
	tester.mockServer = tester.CreateMockServer()

	redis := storage.RedisStorage{}
	err := redis.NewStorage(context.Background(), config)

	if err != nil {
		return err
	}

	generator := storage.RandomKeyGenerator{}
	render := render.HTMLRender{}
	render.New(config)

	sender := webhook.NetWebHookSender{
		Client: tester.mockServer.Client(),
	}

	usecase := application.PaymentsUseCase{
		Repository: &redis,
		Generator:  &generator,
		Render:     &render,
		Sender:     &sender,
		Config:     config,
	}

	server := api.CreateServer(context.Background(), &usecase, config)

	tester.server = httptest.NewServer(server.Handler)
	tester.config = config
	tester.client = tester.server.Client()

	return nil
}

func (tester *Tester) CreateLink(t *testing.T, webhook string) string {
	info := api.PaymentRequest{
		Amount:  "100",
		WebHook: webhook,
		Message: "test-message",
	}

	data, _ := json.Marshal(info)

	response, err := tester.client.Post(tester.server.URL+tester.config.RouterPrefix+"/link", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var link api.LinkResponse
	err = json.NewDecoder(response.Body).Decode(&link)
	defer response.Body.Close()
	assert.NoError(t, err)

	key_index := strings.LastIndex(link.Link, "/")
	assert.NotEqual(t, key_index, -1)
	key := link.Link[key_index+1:]
	assert.NotEqual(t, key, "")

	return key
}

func (tester *Tester) ConfirmLink(t *testing.T, key string) {
	response, err := tester.client.Post(tester.server.URL+tester.config.RouterPrefix+"/confirm/"+key, "message/http", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestFullCycle(t *testing.T) {
	test_env := os.Getenv("TEST_ENV")
	godotenv.Load(test_env)

	config := config.Config{}
	config.LoadFromEnv()

	tester := Tester{}
	err := tester.New(&config)
	assert.NoError(t, err)

	// test timeout
	tester.CreateLink(t, tester.mockServer.URL)

	time.Sleep(time.Duration(2*tester.config.PaymentTimeout) * time.Minute)

	assert.True(t, tester.requestRecieved)

	var webHookResult webhook.WebHookMessage
	err = json.Unmarshal(tester.lastRequestData, &webHookResult)
	assert.NoError(t, err)
	assert.Equal(t, string(domain.TIMEOUT), webHookResult.Status)
	tester.requestRecieved = false

	// test ok
	key := tester.CreateLink(t, tester.mockServer.URL)

	tester.ConfirmLink(t, key)
	time.Sleep(2 * time.Second)

	assert.True(t, tester.requestRecieved)

	err = json.Unmarshal(tester.lastRequestData, &webHookResult)
	assert.NoError(t, err)
	assert.Equal(t, string(domain.OK), webHookResult.Status)
}
