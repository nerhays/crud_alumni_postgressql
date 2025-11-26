package route

import (
	"crud_alumni/app/model"
	"crud_alumni/app/service"
	"crud_alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    api := app.Group("/api")

    // Public
    api.Post("/login", func(c *fiber.Ctx) error {
        var req model.LoginRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
        }
        resp, err := service.Login(req)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(resp)
    })

    // Protected
    protected := api.Group("", middleware.AuthRequired())

    alumni := protected.Group("/alumni")
    alumni.Get("/pag", service.GetAlumniPagination)
    alumni.Get("/", service.GetAllAlumni)  // Admin + User
    alumni.Get("/:id", service.GetAlumniByID) // Admin + User
    alumni.Post("/", middleware.AdminOnly(), service.CreateAlumni)
    alumni.Put("/:id", middleware.AdminOnly(), service.UpdateAlumni)
    alumni.Delete("/:id", middleware.AdminOnly(), service.DeleteAlumni)

    pekerjaan := protected.Group("/pekerjaan")
    pekerjaan.Get("/", service.GetAllPekerjaan)
    pekerjaan.Get("/trash", service.Trash)
    pekerjaan.Get("/pag", service.GetPekerjaanPagination) // Admin + User
    pekerjaan.Put("/:id/soft-delete", service.SoftDeletePekerjaan)
    pekerjaan.Put("/:id/restore", service.RestorePekerjaan)
    pekerjaan.Get("/:id", service.GetPekerjaanByID) // Admin + User
    pekerjaan.Get("/tahun/:tahun", middleware.AdminOnly(), service.GetPekerjaanByTahun)
    pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
    pekerjaan.Post("/", middleware.AdminOnly(), service.CreatePekerjaan)
    pekerjaan.Put("/:id", middleware.AdminOnly(), service.UpdatePekerjaan)
    pekerjaan.Delete("/:id", middleware.AdminOnly(), service.DeletePekerjaan)
    pekerjaan.Delete("/hard/:id", service.HardDeletePekerjaan)

}
