package database

import (
	"fmt"
	"os"

	"chirpyserver/types"
)

type Db struct {
	// mu   *sync.RWMutex
	path string
}

type DBStructure struct {
	Chirps map[int]types.ReturnVals `json:"chirps"`
}

// NewDb creates a new database connection
// and creates the database file if it doesn't exist
func NewDb(path string) (*Db, error) {
	db := &Db{
		path: path,
	}

	_, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("error reading file: ", err)
		return nil, err
	}

	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *Db) CreateChirp(body string) (types.Chirp, error) {
	return types.Chirp{}, nil
}

// GetChirps returns all chrips in the database
func (db *Db) GetChirps() ([]types.Chirp, error) {
	return []types.Chirp{}, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *Db) ensureDB() error {
	return nil
}

// loadDB reads the database file into memory
func (db *Db) loadDB() (DBStructure, error) {
	return DBStructure{}, nil
}

// writeDB writes the database file to disk
func (db *Db) writeDB(dbSt DBStructure) error {
	return nil
}
