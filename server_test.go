package wordcounter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	wcg "github.com/100gle/wordcounter"
	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
)

func TestNewWordCounterServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Testing NewWordCounterServer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wcg.NewWordCounterServer()
			if got.Srv == nil {
				t.Errorf("NewWordCounterServer() returned nil server")
			}
			if got.Srv.HideBanner != true {
				t.Errorf("NewWordCounterServer() HideBanner should be true")
			}
		})
	}
}

func TestWordCounterServer_Ping(t *testing.T) {
	app := echo.New()
	app.GET("/v1/wordcounter/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	server := httptest.NewServer(app)
	defer server.Close()

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	e.GET("/v1/wordcounter/ping").
		Expect().
		Status(http.StatusOK).
		Body().
		NotEmpty().
		IsEqual("pong")
}

func TestWordCounterServer_Count(t *testing.T) {
	app := echo.New()
	server := wcg.NewWordCounterServer()
	apiPath := "/v1/wordcounter/count"
	app.POST(apiPath, server.Count)

	testServer := httptest.NewServer(app)
	defer testServer.Close()

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  testServer.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
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
			statusCode: http.StatusOK,
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
			statusCode: http.StatusOK,
			want:       &wcg.Stats{},
		},
		{
			name: "Test Count without content field",
			content: &struct {
				Foo int `json:"foo"`
			}{Foo: 1},
			statusCode: http.StatusOK,
			want:       &wcg.Stats{},
		},
		{
			name:       "Test Count with invalid body",
			content:    nil,
			statusCode: http.StatusUnprocessableEntity,
			want:       &wcg.Stats{},
		},
	}

	for _, tt := range args {
		t.Run(tt.name, func(t *testing.T) {
			req := e.POST(apiPath)
			if tt.content == nil {
				req.Expect().
					Status(tt.statusCode).
					JSON().
					Object().
					ContainsKey("msg").
					ContainsKey("error")
			} else {
				req.WithJSON(tt.content).Expect().
					Status(tt.statusCode).
					JSON().
					Object().
					ContainsKey("msg").
					ContainsKey("data").
					ContainsKey("error")
			}
		})
	}
}
