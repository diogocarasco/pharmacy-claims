package loader

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/diogocarasco/go-pharmacy-service/internal/database"
	"github.com/diogocarasco/go-pharmacy-service/internal/models"
)

// LoadPharmaciesFromCSV loads pharmacies from a CSV file and saves them to the database.
func LoadPharmaciesFromCSV(filePath string, repo database.DBRepository) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV header: %w", err)
	}

	log.Println("Starting to load pharmacies from CSV...")
	recordsLoaded := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV line: %v. Skipping to next.", err)
			continue
		}

		if len(record) < 2 {
			log.Printf("Invalid CSV line (less than 2 columns): %v. Skipping.", record)
			continue
		}

		pharmacy := models.Pharmacy{
			Chain: record[0],
			NPI:   record[1],
		}

		err = repo.SavePharmacy(pharmacy)
		if err != nil {
			log.Printf("Error saving pharmacy %s to the database: %v", pharmacy.NPI, err)
			continue
		}
		recordsLoaded++
	}

	log.Printf("Finished loading pharmacies from CSV. %d records loaded.", recordsLoaded)
	return nil
}
