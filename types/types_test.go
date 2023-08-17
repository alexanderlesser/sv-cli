package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntrySerializationDeserialization(t *testing.T) {
	originalEntry := Entry{
		Name:         "EntryName",
		Time:         "12:00",
		Date:         "2023-08-17",
		ErrorWarning: true,
	}

	// Serialize the Entry
	serializedEntry, err := json.Marshal(originalEntry)
	assert.NoError(t, err)

	// Deserialize the Entry
	var deserializedEntry Entry
	err = json.Unmarshal(serializedEntry, &deserializedEntry)
	assert.NoError(t, err)

	// Check if the deserialized Entry matches the original
	assert.Equal(t, originalEntry, deserializedEntry)
}

func TestRecordSerializationDeserialization(t *testing.T) {
	originalRecord := Record{
		ID:       1,
		Username: "testuser",
		Password: "testpass",
		Name:     "Test Record",
		Domain:   "example.com",
		Path:     "/path",
		Entries: []Entry{
			{
				Name:         "EntryName",
				Time:         "12:00",
				Date:         "2023-08-17",
				ErrorWarning: true,
			},
		},
	}

	// Serialize the Record
	serializedRecord, err := json.Marshal(originalRecord)
	assert.NoError(t, err)

	// Deserialize the Record
	var deserializedRecord Record
	err = json.Unmarshal(serializedRecord, &deserializedRecord)
	assert.NoError(t, err)

	// Check if the deserialized Record matches the original
	assert.Equal(t, originalRecord, deserializedRecord)
}
