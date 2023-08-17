package datastore

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/alexanderlesser/sv-cli/testUtil"
	"github.com/alexanderlesser/sv-cli/types"
	"github.com/stretchr/testify/assert"
)

func TestInitDataStorage(t *testing.T) {

	defer os.Remove(constants.DATASTORE_NAME)

	// setup
	defer func() {
		restore := testUtil.SetupTestEnvironment()
		defer restore() // Restore the original environment after the tests are done
	}()

	err := InitDataStorage()
	assert.NoError(t, err)

	// Clean up after the test
	defer os.Remove(constants.DATASTORE_NAME)
}

func TestSaveAndLoad(t *testing.T) {
	defer os.Remove(constants.DATASTORE_NAME)

	// setup
	defer func() {
		restore := testUtil.SetupTestEnvironment()
		defer restore() // Restore the original environment after the tests are done
	}()

	record := types.Record{
		// Initialize the fields as needed
		ID:       1,
		Username: "testuser",
		Password: "testpass",
		Name:     "Test Record",
		Domain:   "example.com",
		Path:     "/path",
		Entries:  []types.Entry{},
	}

	err := Save(record)
	assert.NoError(t, err)

	loadedData, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, []types.Record{record}, loadedData)

	// Clean up after the test
	defer os.Remove(constants.DATASTORE_NAME)
}

func TestUpdateRecord(t *testing.T) {
	// Clean up any existing data file before testing
	defer os.Remove(constants.DATASTORE_NAME)

	defer func() {
		restore := testUtil.SetupTestEnvironment()
		defer restore() // Restore the original environment after the tests are done
	}()

	// Prepare initial data
	initialData := []types.Record{
		{
			ID:       1,
			Username: "user1",
		},
		{
			ID:       2,
			Username: "user2",
		},
	}

	// Save initial data
	dataBytes, _ := json.Marshal(initialData)
	_ = os.WriteFile(constants.DATASTORE_NAME, dataBytes, 0644)

	// Call UpdateRecord
	updatedRecord := types.Record{
		ID:       2,
		Username: "updateduser2",
	}
	err := UpdateRecord(updatedRecord)
	assert.NoError(t, err)

	// Load data and verify updated record
	loadedData, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(loadedData))
	assert.Equal(t, "updateduser2", loadedData[1].Username)

	// Clean up after the test
	defer os.Remove(constants.DATASTORE_NAME)
}

func TestDeleteRecord(t *testing.T) {
	// Clean up any existing data file before testing
	defer os.Remove(constants.DATASTORE_NAME)

	defer func() {
		restore := testUtil.SetupTestEnvironment()
		defer restore() // Restore the original environment after the tests are done
	}()

	initialData := []types.Record{
		{
			ID:       1,
			Username: "user1",
		},
		{
			ID:       2,
			Username: "user2",
		},
	}

	// Save initial data
	dataBytes, _ := json.Marshal(initialData)
	_ = os.WriteFile(constants.DATASTORE_NAME, dataBytes, 0644)

	recordIDToDelete := int32(1)
	err := DeleteRecord(recordIDToDelete)
	assert.NoError(t, err)

	// Load data and verify deleted record
	loadedData, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loadedData))
	assert.Equal(t, int32(2), loadedData[0].ID)

	defer os.Remove(constants.DATASTORE_NAME)
}

// Run tests
func TestMain(m *testing.M) {
	// Set up any necessary environment or mocks

	// Run the tests
	exitCode := m.Run()

	os.Exit(exitCode)
}
