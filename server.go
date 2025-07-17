package wordcounter

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WordCounterServer struct {
	Echo *echo.Echo
}

type CountBody struct {
	Content string `json:"content"`
}

func NewWordCounterServer() *WordCounterServer {
	echoServer := echo.New()
	echoServer.HideBanner = true
	return &WordCounterServer{Echo: echoServer}
}

func (s *WordCounterServer) Count(c echo.Context) error {
	body := new(CountBody)
	errMsg := ""
	counter := NewCounter()

	// Check if request has a body
	if c.Request().ContentLength == 0 {
		errMsg = "request body is empty"
		return c.JSON(http.StatusUnprocessableEntity, map[string]any{
			"msg":   "parse failed",
			"error": errMsg,
		})
	}

	if err := c.Bind(body); err != nil {
		errMsg = fmt.Sprintf("%s", err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]any{
			"msg":   "parse failed",
			"error": errMsg,
		})
	}

	err := counter.Count(body.Content)
	if err != nil {
		errMsg = fmt.Sprintf("%s", err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"msg":   "ok",
		"data":  counter.Stats,
		"error": errMsg,
	})
}

func (s *WordCounterServer) Run(port int) error {
	s.Echo.GET(PingEndpoint, func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	s.Echo.POST(CountEndpoint, s.Count)

	return s.Echo.Start(fmt.Sprintf(":%d", port))
}
