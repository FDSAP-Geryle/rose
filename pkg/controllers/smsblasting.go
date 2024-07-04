package controllers

import (
	"rose/pkg/models/response"
	models "rose/pkg/models/struct"
	"rose/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

// Get All SMS Type
func SMSBlastingContentType(c *fiber.Ctx) error {
	// Start of getting info of the current logged-in user
	// claims, ok := c.Locals("user").(*JWTClaims)
	// if !ok {
	// 	// Handle user claims not found
	// 	return c.Status(fiber.StatusInternalServerError).JSON(response.RBIGetSingleClientResponse{
	// 		RetCode: "500",
	// 		Message: "User claims not found",
	// 		Data:    nil,
	// 	})
	// }
	// fmt.Println(claims.Username)
	// End of getting info of the current logged-in user

	getAllSMSContent := []models.SMSContent{}

	if fetchErr := database.DBConn.Debug().
		Table("sms_blasting_content").
		Select("DISTINCT ON (sms_type) *"). // Select distinct records based on sms_type
		//Where("status = ?", "Active").                               // Add the condition for "status"
		//Order("sms_type, COALESCE(date_modified, date_added) DESC"). // Order by sms_type and date_modified if not empty, otherwise use date_added
		Order("sms_type ASC"). // Order by sms_type and date_modified if not empty, otherwise use date_added
		Find(&getAllSMSContent).
		Error; fetchErr != nil {
		return c.JSON(response.SMSContentResponse{
			RetCode: "400",
			Message: "Cannot Retrieve List of SMS Type",
			Data:    nil,
		})
	}

	if len(getAllSMSContent) == 0 {
		// Handle case where no clients are found
		return c.JSON(response.SMSContentResponse{
			RetCode: "404",
			Message: "No SMS Type Found",
			Data:    nil,
		})
	}

	return c.JSON(response.SMSContentResponse{
		RetCode: "200",
		Message: "SMS Type Found",
		Data:    getAllSMSContent,
	})
}
