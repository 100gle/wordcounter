package wordcounter_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
			if got.Echo == nil {
				t.Errorf("NewWordCounterServer() returned nil server")
			}
			if got.Echo.HideBanner != true {
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

// TestWordCounterServer_Run tests the Run function with server startup and shutdown
func TestWordCounterServer_Run(t *testing.T) {
	server := wcg.NewWordCounterServer()

	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		err := server.Run(port)
		serverErr <- err
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Test that server is running and routes are registered
	baseURL := fmt.Sprintf("http://localhost:%d", port)

	// Test ping endpoint
	resp, err := http.Get(baseURL + "/v1/wordcounter/ping")
	if err != nil {
		t.Errorf("Failed to reach ping endpoint: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for ping, got %d", resp.StatusCode)
		}
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Echo.Shutdown(ctx); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}

	// Check if server stopped properly
	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server returned unexpected error: %v", err)
		}
	case <-time.After(6 * time.Second):
		t.Error("Server did not stop within timeout")
	}
}

// TestWordCounterServer_CountErrorHandling tests error handling in Count function
func TestWordCounterServer_CountErrorHandling(t *testing.T) {
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

	// Test with empty content - should trigger error in Count function
	e.POST(apiPath).
		WithJSON(&wcg.CountBody{Content: ""}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("msg").
		ContainsKey("data").
		ContainsKey("error").
		Value("error").NotEqual("")

	// Test with valid content but empty string (edge case)
	e.POST(apiPath).
		WithJSON(&wcg.CountBody{Content: " "}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("msg").
		ContainsKey("data").
		ContainsKey("error")
}

// TestWordCounterServer_CountEmptyBody tests Count function with empty request body
func TestWordCounterServer_CountEmptyBody(t *testing.T) {
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

	// Test with completely empty body (ContentLength = 0)
	e.POST(apiPath).
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().
		Object().
		ContainsKey("msg").
		ContainsKey("error").
		Value("msg").String().Equal("parse failed")
}

// TestWordCounterServer_CountInvalidJSON tests Count function with invalid JSON
func TestWordCounterServer_CountInvalidJSON(t *testing.T) {
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

	// Test with invalid JSON
	e.POST(apiPath).
		WithText("invalid json {").
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().
		Object().
		ContainsKey("msg").
		ContainsKey("error").
		Value("msg").Equal("parse failed")
}
