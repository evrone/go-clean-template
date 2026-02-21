package middleware

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func buildPanicMessage(ctx *fiber.Ctx, err any) string {
	var result strings.Builder

	result.WriteString(ctx.IP())
	result.WriteString(" - ")
	result.WriteString(ctx.Method())
	result.WriteString(" ")
	result.WriteString(ctx.OriginalURL())
	result.WriteString(" PANIC DETECTED: ")
	result.WriteString(fmt.Sprintf("%v\n%s\n", err, debug.Stack())) //nolint: staticcheck,gocritic // it's okay for panic

	return result.String()
}

func logPanic(l logger.Interface) func(c *fiber.Ctx, err any) {
	return func(ctx *fiber.Ctx, err any) {
		l.Error(buildPanicMessage(ctx, err))
	}
}

func Recovery(l logger.Interface) func(c *fiber.Ctx) error {
	return fiberRecover.New(fiberRecover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: logPanic(l),
	})
}
