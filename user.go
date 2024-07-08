package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/HalilFocic/chirpy/internal/database"
)

func (cfg *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong!")
		return
	}
	users := []database.User{}
	for _, u := range dbUsers {
		users = append(users, database.User{
			Id:    u.Id,
			Email: u.Email,
		})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})
	respondWithJSON(w, http.StatusOK, users)

}

type userBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handleAddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while decoding body")
		return
	}
	u, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Couldnt create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
	return
}
