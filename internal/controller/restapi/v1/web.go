package v1

import "github.com/gofiber/fiber/v2"

func (r *V1) showUI(ctx *fiber.Ctx) error {
	return ctx.SendFile("./internal/controller/restapi/v1/web/index.html")
}
