package wordcounter

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type WordCounterServer struct {
	srv *fiber.App
}

type CountBody struct {
	Content string `json:"content"`
}

func NewWordCounterServer() *WordCounterServer {
	srv := fiber.New(fiber.Config{
		AppName: "WordCounter",
	})
	return &WordCounterServer{srv: srv}
}

func (s *WordCounterServer) Count(ctx *fiber.Ctx) error {
	body := new(CountBody)
	errMsg := ""
	tc := NewTextCounter()

	if err := ctx.BodyParser(body); err != nil {
		errMsg = fmt.Sprintf("%s", err)
		return ctx.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"msg": "parse failed", "error": errMsg})
	}

	err := tc.Count(body.Content)
	if err != nil {
		errMsg = fmt.Sprintf("%s", err)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "ok", "data": tc.S, "error": errMsg})
}

func (s *WordCounterServer) Run(port int) error {
	s.srv.Get("/v1/wordcounter/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("pong")
	})
	s.srv.Post("/v1/wordcounter/count", s.Count)

	return s.srv.Listen(fmt.Sprintf(":%d", port))
}
