package uploadmerchant

import (
	"fmt"
	"log"
	models "rosei/pkg/models/merchant"
	"rosei/pkg/utils/go-utils/database"
	"time"

	"github.com/google/uuid"
)

// InsertReceiveRecord processes data by removing specific fields, setting merchant_name,
// and assigning random values to mid and mpan.
func InsertReceiveRecord(data []map[string]interface{}) {
	SnakeCase(data) // Assuming this function handles any necessary transformations.

	for _, record := range data {
		// Remove the specified fields
		delete(record, "no")

		// Set the merchant_name field to the value of store_name and remove store_name
		if storeName, exists := record["store_name"]; exists {
			record["merchant_name"] = storeName
			delete(record, "store_name")
		}

		if OwnerAddress, exists := record["owner_address"]; exists {
			record["merchant_address"] = OwnerAddress
			delete(record, "owner_address")
		}

		if AOName, exists := record["ao_name"]; exists {
			record["agent_name"] = AOName
			delete(record, "ao_name")
		}

		if AOCode, exists := record["ao_code"]; exists {
			record["agent_id"] = AOCode
			delete(record, "ao_code")
		}

		// Set uploaded_at to the current timestamp
		record["uploaded_at"] = time.Now().Format("2006-01-02 15:04:05")

		// Generate random UUIDs for mid and mpan
		record["mid"] = uuid.New().String()
		record["mpan"] = uuid.New().String()
	}

	// Insert the processed records into the database
	err := database.DBConn.Model(&models.ReceiveUploadActivated{}).Create(data).Error
	if err != nil {
		log.Println("Failed to save records:", err)
	} else {
		fmt.Println("Records inserted successfully.")
	}
}
