package wordcounter_test

import (
	"net/http"
	"testing"

	wcg "github.com/100gle/wordcounter"
	"github.com/gavv/httpexpect/v2"
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

func TestWordCounterServer_Ping(t *testing.T) {
	app := fiber.New()
	app.Get("/v1/wordcounter/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("pong")
	})

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(app.Handler()),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	e.GET("/v1/wordcounter/ping").
		Expect().
		Status(fiber.StatusOK).
		Body().
		NotEmpty().
		IsEqual("pong")
}

func TestWordCounterServer_Count(t *testing.T) {
	app := fiber.New()
	server := wcg.NewWordCounterServer()
	apiPath := "/v1/wordcounter/count"
	app.Post(apiPath, server.Count)

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(app.Handler()),
		},
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
