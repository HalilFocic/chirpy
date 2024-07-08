package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DbStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

var ErrNotExist = errors.New("Resource does not exist")

func NewDB(path string) (*DB, error) {
	db := DB{
		mux:  &sync.RWMutex{},
		path: path,
	}
	err := db.ensureDB()
	return &db, err
}


func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}
func (db *DB) createDB() error {
	dbStructure := DbStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	return db.writeDB(dbStructure)
}
func (db *DB) loadDB() (DbStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStructure := DbStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil

}
func (db *DB) writeDB(dbStuct DbStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStuct)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil

}



