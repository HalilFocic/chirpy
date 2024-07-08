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
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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
		return Chirp{}, err
	}
	id := len(dbStruct.Chirps) + 1
	c := Chirp{
		id,
		body,
	}
	dbStruct.Chirps[id] = c
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
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

func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	val, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("Chirp with that id was not found")
	}
	return val, nil

}

func (db *DB) CreateUser(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStruct.Users) + 1
	c := User{
		id,
		email,
	}
	dbStruct.Users[id] = c
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return c, nil

}
func (db *DB) GetUsers() ([]User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return []User{}, err
	}

	users := make([]User, 0, len(dbStruct.Users))
	for _, u := range dbStruct.Users {
		users = append(users, u)
	}
	return users, nil
}
