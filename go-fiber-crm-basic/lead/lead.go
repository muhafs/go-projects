package lead

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-crm-basic/database"
	"gorm.io/gorm"
)

type Lead struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

func GetLeads(ctx *fiber.Ctx) error {
	db := database.ConnectDB

	var leads []Lead
	db.Find(&leads)

	return ctx.JSON(leads)
}

func GetLead(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	db := database.ConnectDB

	var lead Lead
	db.Find(&lead, id)
	if lead.Name == "" {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Lead not found",
		})
	}

	return ctx.JSON(lead)
}

func NewLead(ctx *fiber.Ctx) error {
	db := database.ConnectDB

	lead := new(Lead) // same as &Lead{}
	if err := ctx.BodyParser(lead); err != nil {
		return ctx.Status(503).JSON(fiber.Map{
			"message": err,
		})
	}

	db.Create(&lead)

	return ctx.JSON(lead)
}
func DeleteLead(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	db := database.ConnectDB

	var lead Lead
	db.Find(&lead, id)
	if lead.Name == "" {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Lead not found",
		})
	}

	db.Delete(&lead)

	return ctx.JSON(fiber.Map{
		"message": "lead has deleted",
	})
}
