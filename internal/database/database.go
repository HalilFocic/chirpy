package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DbStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	db := DB{
		mux:  &sync.RWMutex{},
		path: path,
	}
	return &db, nil
}

func (db *DB) GetChrips() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, nil
	}

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		fmt.Println("Error loading db")
		return Chirp{}, nil
	}
	id := len(dbStruct.Chirps) + 1
	c := Chirp{
		id,
		body,
	}
	dbStruct.Chirps[id] = c
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, nil
	}
	return c, nil

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
