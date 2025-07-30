package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/diogocarasco/go-pharmacy-service/internal/database"
	"github.com/diogocarasco/go-pharmacy-service/internal/models"
)

type ClaimLoader struct {
	DBRepo database.DBRepository
}

func NewClaimLoader(dbRepo database.DBRepository) *ClaimLoader {
	return &ClaimLoader{DBRepo: dbRepo}
}

// LoadAndSaveClaimsFromDir reads all JSON files from a directory and saves them to the database.
func (cl *ClaimLoader) LoadAndSaveClaimsFromDir(dirPath string) error {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("error obtaining absolute path of claims directory %s: %w", dirPath, err)
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		return fmt.Errorf("error reading claims directory %s: %w", absPath, err)
	}

	var allClaims []models.Claim

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".json") {
			log.Printf("INFO: Ignoring non-JSON file: %s/%s", absPath, file.Name())
			continue
		}

		filePath := filepath.Join(absPath, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("ERROR: Error reading claims file %s: %v", filePath, err)
			continue
		}

		var claimsFromFile []models.Claim
		if err := json.Unmarshal(data, &claimsFromFile); err != nil {
			log.Printf("ERROR: Error decoding JSON from claims file %s: %v", filePath, err)
			continue
		}

		log.Printf("INFO: Loaded %d claims from file: %s", len(claimsFromFile), file.Name())
		allClaims = append(allClaims, claimsFromFile...)
	}

	log.Printf("INFO: Finished loading claims from all JSON files. Total of %d records found.", len(allClaims))

	if len(allClaims) > 0 {
		log.Println("INFO: Starting to save all claims to the database...")
		if err := cl.DBRepo.SaveClaims(allClaims); err != nil {
			return fmt.Errorf("error saving claims to the database: %w", err)
		}
		log.Println("INFO: All claims saved successfully to the database.")
	} else {
		log.Println("INFO: No claims to save to the database.")
	}

	return nil
}
