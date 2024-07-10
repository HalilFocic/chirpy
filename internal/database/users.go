package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	RefreshToken   string `json:"refresh_token"`
}

var ErrAlreadyExists = errors.New("Already exists!")

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	randByte := make([]byte, 32)
	_, err = rand.Read(randByte)
	if err != nil {
		return User{}, err
	}
	encodedString := hex.EncodeToString(randByte)
	user := User{
		Id:             id,
		Email:          email,
		HashedPassword: hashedPassword,
		RefreshToken:   encodedString,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []User{}, err
	}
	users := make([]User, 0, len(dbStructure.Users))
	for _, v := range dbStructure.Users {
		users = append(users, v)
	}
	return users, nil

}
func (db *DB) UpdateUser(userID int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStructure.Users[userID]
	if !ok {
		return User{}, ErrNotExist
	}
	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[userID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUserByRefreshToken(rToken string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	users := dbStructure.Users
	for _, v := range users {
		if v.RefreshToken == rToken {
			return User{
				Email: v.Email,
				Id:    v.Id,
			}, nil
		}
	}
	return User{}, ErrNotExist
}
func (db *DB) DeleteUserToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	userId := 0
	for k, v := range dbStructure.Users {
		if v.RefreshToken == token {
			userId = k
		}
	}
	if userId == 0 {
		return ErrNotExist
	}
    user := dbStructure.Users[userId]
    newUser := User {
        Email: user.Email,
        Id: user.Id,
        HashedPassword: user.HashedPassword,
        RefreshToken: "",

    }
    dbStructure.Users[userId] = newUser
	db.writeDB(dbStructure)
    fmt.Println("vracam nil")
    return nil
}
