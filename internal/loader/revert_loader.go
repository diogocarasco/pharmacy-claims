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

type RevertLoader struct {
	DBRepo database.DBRepository
}

func NewRevertLoader(dbRepo database.DBRepository) *RevertLoader {
	return &RevertLoader{DBRepo: dbRepo}
}

// LoadAndSaveRevertsFromDir reads all JSON files from a directory and saves them to the database.
func (rl *RevertLoader) LoadAndSaveRevertsFromDir(dirPath string) error {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("error obtaining absolute path of reverts directory %s: %w", dirPath, err)
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		return fmt.Errorf("error reading reverts directory %s: %w", absPath, err)
	}

	var allReverts []models.Revert

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
			log.Printf("ERROR: Error reading reverts file %s: %v", filePath, err)
			continue
		}

		var revertsFromFile []models.Revert
		if err := json.Unmarshal(data, &revertsFromFile); err != nil {
			log.Printf("ERROR: Error decoding JSON from reverts file %s: %v", filePath, err)
			continue
		}

		log.Printf("INFO: Loaded %d reverts from file: %s", len(revertsFromFile), file.Name())
		allReverts = append(allReverts, revertsFromFile...)
	}

	log.Printf("INFO: Finished loading reverts from all JSON files. Total of %d records found.", len(allReverts))

	if len(allReverts) > 0 {
		log.Println("INFO: Starting to save all reverts to the database...")
		if err := rl.DBRepo.SaveReverts(allReverts); err != nil {
			return fmt.Errorf("error saving reverts to the database: %w", err)
		}
		log.Println("INFO: All reverts saved successfully to the database.")
	} else {
		log.Println("INFO: No reverts to save to the database.")
	}

	return nil
}
