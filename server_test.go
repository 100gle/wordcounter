package wordcounter_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	wcg "github.com/100gle/wordcounter"
	"github.com/gofiber/fiber/v2"
)

func TestNewWordCounterServer(t *testing.T) {
	tests := []struct {
		name string
		want *wcg.WordCounterServer
	}{
		{
			name: "Testing NewWordCounterServer",
			want: &wcg.WordCounterServer{Srv: fiber.New(fiber.Config{
				AppName: "WordCounter",
			})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wcg.NewWordCounterServer()
			gotConfig := got.Srv.Config()
			wantConfig := tt.want.Srv.Config()

			if gotConfig.AppName != wantConfig.AppName {
				t.Errorf("NewWordCounterServer() = %v, want %v", gotConfig.AppName, wantConfig.AppName)
			}
		})
	}
}

func TestWordCounterServer_Count(t *testing.T) {
	app := wcg.NewWordCounterServer()

	go app.Run(8080)

	t.Run("Test ping", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:8080/v1/wordcounter/ping", nil)
		resp, _ := app.Srv.Test(req)
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("got %v, want %v", resp.StatusCode, fiber.StatusOK)
		}
		ret, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(ret) != "pong" {
			t.Errorf("got %v, want %v", string(ret), "pong")
		}
	})

	args := []struct {
		name       string
		content    interface{}
		statusCode int
		want       *wcg.Stats
	}{
		{
			name: "Testing Count",
			content: &wcg.CountBody{
				Content: "你好，世界！\n这是一个测试",
			},
			statusCode: fiber.StatusOK,
			want: &wcg.Stats{
				Lines:           2,
				NonChineseChars: 3,
				ChineseChars:    10,
				TotalChars:      13,
			},
		},
		{
			name:       "Testing Count with empty string",
			content:    &wcg.CountBody{Content: ""},
			statusCode: fiber.StatusOK,
			want:       &wcg.Stats{},
		},
		{
			name: "Test Count without content field",
			content: &struct {
				Foo int `json:"foo"`
			}{Foo: 1},
			statusCode: fiber.StatusOK,
			want:       &wcg.Stats{},
		},
		{
			name:       "Test Count with invalid body",
			content:    nil,
			statusCode: fiber.StatusUnprocessableEntity,
			want:       &wcg.Stats{},
		},
	}

	for _, tt := range args {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.content != nil {
				body, err := json.Marshal(tt.content)
				if err != nil {
					t.Fatal(err)
				}
				query := bytes.NewReader(body)
				req = httptest.NewRequest("POST", "http://localhost:8080/v1/wordcounter/count", query)
			} else {
				req = httptest.NewRequest("POST", "http://localhost:8080/v1/wordcounter/count", nil)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			resp, _ := app.Srv.Test(req)

			if resp.StatusCode != tt.statusCode {
				t.Errorf("Count() = %v, want %v", resp.StatusCode, tt.statusCode)
			}
		})
	}
}
