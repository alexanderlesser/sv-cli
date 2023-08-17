package datastore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/alexanderlesser/sv-cli/types"
)

var dataStorage *FileStorage

type FileStorage struct {
	FilePath string
}

// Initialize the data file
func InitDataStorage() error {
	dataDir, err := getDataDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(dataDir, constants.DATASTORE_NAME)
	dataStorage = &FileStorage{FilePath: filePath}

	//Check if the data file exists
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// Create the data file and initialize it with an empty array
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Initialize the JSON array
		emptyArray := []interface{}{}
		encoder := json.NewEncoder(file)
		err = encoder.Encode(emptyArray)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get the data directory
func getDataDir() (string, error) {
	if constants.DEVELOPMENT {
		// Use the current directory for development mode
		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return currentDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dataDir := filepath.Join(homeDir, constants.CONFIG_DIR_NAME)
	return dataDir, nil
}

func ensureDataDirectory() error {
	// Get the directory path
	dir := filepath.Dir(dataStorage.FilePath)

	// Check if the directory exists, create it if not
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// Load loads data from the data file.
func Load() ([]types.Record, error) {
	var existingData []types.Record
	dataBytes, err := os.ReadFile(dataStorage.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return existingData, nil
		}
		return existingData, err
	}

	err = json.Unmarshal(dataBytes, &existingData)
	if err != nil {
		return existingData, err
	}

	return existingData, nil
}

// Save record to data file
func Save(record types.Record) error {
	var existingData []types.Record
	existingData, err := Load()
	if err != nil {
		if os.IsNotExist(err) {
			existingData = []types.Record{}
		} else {
			return err
		}
	}

	record.Entries = []types.Entry{}
	record.ID = int32(len(existingData) + 1)
	existingData = append(existingData, record)

	dataBytes, err := json.Marshal(existingData)
	if err != nil {
		return err
	}

	// Ensure the data directory exists
	if err := ensureDataDirectory(); err != nil {
		return err
	}

	return os.WriteFile(dataStorage.FilePath, dataBytes, 0644)
}

// Update record in data file
func UpdateRecord(updatedRecord types.Record) error {
	var data []types.Record
	data, err := Load()
	if err != nil {
		return err
	}

	var updated = false

	for i, record := range data {
		if record.ID == updatedRecord.ID {
			data[i] = updatedRecord

			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("record with ID %d not found", updatedRecord.ID)
	}

	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	if err := ensureDataDirectory(); err != nil {
		return err
	}

	err = os.WriteFile(dataStorage.FilePath, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Delete record by ID from data file
func DeleteRecord(recordID int32) error {
	var data []types.Record
	data, err := Load()
	if err != nil {
		return err
	}

	// Find and remove the record with the given ID
	indexToDelete := -1
	for i, record := range data {
		if record.ID == recordID {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return fmt.Errorf("record with ID %d not found", recordID)
	}

	data = append(data[:indexToDelete], data[indexToDelete+1:]...)

	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	if err := ensureDataDirectory(); err != nil {
		return err
	}

	err = os.WriteFile(dataStorage.FilePath, dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
