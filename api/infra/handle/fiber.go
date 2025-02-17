package handle

/* type Request any
type Response any



func Handle[Req Request, Res Response](handler func(ctx context.Context, req *Req) Res) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req Req

		// body parser
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// param parser
		if err := c.ParamsParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// query parser
		if err := c.QueryParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// header parser
		if err := c.ReqHeaderParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// validation
		if err := req.Validate(); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		// timeout context
		ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
		defer cancel()

		res := handler(ctx, &req)

		return c.JSON(res)
	}
} */

// usage
//app.Get("/test", basehandler.Handle(userHandler.Login))
