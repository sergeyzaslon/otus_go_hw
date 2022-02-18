package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	internalgrpc "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// test stores the HTTP testing client preconfigured
// var http = baloo.New(os.Getenv("APP_HOST"))

func TestMain(t *testing.T) {
	host := os.Getenv("APP_HOST")

	httpClient := http.Client{}

	res, err := httpClient.Get(host + "/events") // nolint: noctx
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf := strings.Builder{}
	io.Copy(&buf, res.Body)
	res.Body.Close()
	require.Equal(t, "[]", buf.String())

	// Создадим новое событие
	body := `{
		"id": "4927aa58-a175-429a-a125-c04765597152",
		"title": "Test Event 01",
		"description": "Test Event Description 01",
		"date": "2021-12-20 12:30:00",
		"duration": 60,
		"userId": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		"notifyBeforeSeconds": 60
	}`

	respCode, _, err := RESTPost(&httpClient, host+"/events", body)
	if err != nil {
		t.Errorf("Failed to create /events: %s", err)
		t.FailNow()
	}
	require.Equal(t, 201, respCode)

	// Созданное событие можно прочитать в общем списке
	resCode, resBody, err := RESTGet(&httpClient, host+"/events")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}

	bodyExpected := `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01","date":"2021-12-20 12:30:00","duration":60,"description":"Test Event Description 01","userId":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notifyBeforeSeconds":60}]` // nolint:lll
	require.Equal(t, bodyExpected, resBody)
	require.Equal(t, 200, resCode)

	// Попробуем создать ещё одно такое же события - должны получить ОШИБКУ
	respCode, respBody, err := RESTPost(&httpClient, host+"/events", body)
	if err != nil {
		t.Errorf("Failed to create /events: %s", err)
		t.FailNow()
	}
	bodyExpected = `{"success":false,"error":"validation error: event with such id already exists"}`
	require.Equal(t, bodyExpected, respBody)
	require.Equal(t, 400, respCode)

	// Обновим событие
	body = `{
		"title": "Test Event 01 UPD",
		"description": "Test Event Description 01 UPD",
		"date": "2021-12-20 12:30:30",
		"duration": 70,
		"userId": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		"notifyBeforeSeconds": 70
	}`
	respCode, _, err = RESTPut(&httpClient, host+"/events/4927aa58-a175-429a-a125-c04765597152", body)
	if err != nil {
		t.Errorf("Failed to update event: %s", err)
		t.FailNow()
	}
	require.Equal(t, 200, respCode)

	// Прочитаем обновления
	respCode, respBody, err = RESTGet(&httpClient, host+"/events")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
		t.FailNow()
	}
	bodyExpected = `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01 UPD","date":"2021-12-20 12:30:30","duration":70,"description":"Test Event Description 01 UPD","userId":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notifyBeforeSeconds":70}]` // nolint:lll
	require.Equal(t, bodyExpected, respBody)
	require.Equal(t, 200, respCode)

	// Удалим событие
	_, _, err = RESTDelete(&httpClient, host+"/events/4927aa58-a175-429a-a125-c04765597152")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
		t.FailNow()
	}

	// Удалённого события больше нет в списке
	// Созданное событие можно прочитать в общем списке
	respCode, respBody, err = RESTGet(&httpClient, host+"/events")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}

	bodyExpected = `[]`
	require.Equal(t, bodyExpected, respBody)
	require.Equal(t, 200, respCode)

	// Проверим, что было уведомление о нашем событии, так как мы его создали в прошлом
	time.Sleep(time.Second * 10)
	logFileName := "/var/logs/app.log"
	content, err := os.ReadFile(logFileName)
	if err != nil {
		t.Errorf("Failed to read sender logs")
		t.FailNow()
	}
	contentString := string(content)

	fmt.Println(contentString)

	require.Contains(t, contentString, "4927aa58-a175-429a-a125-c04765597152")
}

var host = os.Getenv("APP_HOST_GRPC")

func TestMainGrpc(t *testing.T) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("Failed to dial GRPC service: %s", err)
	}

	ctx := context.Background()

	client := internalgrpc.NewEventServiceClient(conn)

	event := internalgrpc.Event{
		Id:                  "4927aa58-a175-429a-a125-c04765597152",
		Title:               "Event",
		Description:         "Event Descr",
		Date:                "2021-12-20 12:30:30",
		Duration:            45,
		UserId:              "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		NotifyBeforeSeconds: 100,
	}
	_, err = client.Create(ctx, &event)
	assert.Nil(t, err)

	req := internalgrpc.EventListRequest{
		Date: "2021-12-01",
	}
	res, err := client.EventListMonth(ctx, &req)
	assert.Nil(t, err)
	assert.Len(t, res.GetEvents(), 1)
}

func RESTPost(httpClient *http.Client, url, body string) (int, string, error) {
	return RESTWithPayload(httpClient, url, body, "POST")
}

func RESTPut(httpClient *http.Client, url, body string) (int, string, error) {
	return RESTWithPayload(httpClient, url, body, "PUT")
}

func RESTDelete(httpClient *http.Client, url string) (int, string, error) {
	return RESTWithPayload(httpClient, url, "", "DELETE")
}

func RESTWithPayload(httpClient *http.Client, url, body, method string) (int, string, error) {
	bufWrite := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(method, url, bufWrite) // nolint: noctx
	if err != nil {
		return 0, "", err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	resBuilder := strings.Builder{}
	io.Copy(&resBuilder, res.Body)

	return res.StatusCode, resBuilder.String(), nil
}

func RESTGet(httpClient *http.Client, url string) (int, string, error) {
	res, err := httpClient.Get(url) // nolint: noctx
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	buf := strings.Builder{}
	io.Copy(&buf, res.Body)

	return res.StatusCode, buf.String(), nil
}
