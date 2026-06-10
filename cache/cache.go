package cache

import (
	"encoding/json"
	"os"

	"github.com/visagex/osrsdb-api/models"
)

func LoadCache(path string) ([]models.OsrsItem, error) {
	var err error
	var data []byte
	data, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var items []models.OsrsItem

	err = json.Unmarshal(data, &items)
	return items, err
}

func SaveCache(items []models.OsrsItem, path string) error {
	var err error
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	return err
}

func CreateCache() (*os.File, error) {
	file, err := os.Create("osrs-db-cache")
	if err != nil {
		return nil, err
	}

	return file, err
}
