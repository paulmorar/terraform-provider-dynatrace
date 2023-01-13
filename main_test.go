package main_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

var client = http.DefaultClient

func GET(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(res.Body)
}

func POST(url string, v any, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 204 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return data, errors.New("Validation failed")
	}
	return io.ReadAll(res.Body)
}

func token(tok string) map[string]string {
	return map[string]string{"Authorization": "Api-Token " + tok}
}

func TestFoo(t *testing.T) {
	t.Skip()
	demoLiveToken := os.Getenv("DEMO_LIVE_TOKEN")
	sizToken := os.Getenv("SIZ_TOKEN")
	hide(demoLiveToken)
	hide(sizToken)
	data, err := GET("https://guu84124.live.dynatrace.com/api/config/v1/dashboards", token(demoLiveToken))
	if err != nil {
		panic(err)
	}
	var m map[string]any
	err = json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	ids := []string{}
	stubs := m["dashboards"].([]any)
	for _, untypedStub := range stubs {
		stub := untypedStub.(map[string]any)
		ids = append(ids, stub["id"].(string))
	}
	// ids = []string{"9b386608-d386-a9c0-163f-aa67d28208ec"}
	for _, id := range ids {
		data, err := GET(fmt.Sprintf("https://guu84124.live.dynatrace.com/api/config/v1/dashboards/%s", id), token(demoLiveToken))
		if err != nil {
			panic(err)
		}
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			panic(err)
		}
		delete(m, "id")
		delete(m, "metadata")
		data, err = POST("https://siz65484.live.dynatrace.com/api/config/v1/dashboards/validator", m, token(sizToken))
		if err != nil {
			if err.Error() == "Validation failed" {
				fmt.Println(id, string(data))
			} else {
				panic(err)
			}
		}
	}

}

func hide(v any) {}
