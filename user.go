package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/HalilFocic/chirpy/internal/auth"
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
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while hashing password")
		return
	}
	u, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Couldnt create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
	return
}
func (cfg *apiConfig) handleUpateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		userBody
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		fmt.Printf("err: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Couldnt decode parameters")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while decoding body")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while hashing password")
		return
	}
	userId, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt convert user id")
		return
	}
	u, err := cfg.DB.UpdateUser(userId, params.Email, hashedPassword)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldnt create user")
		return
	}
	type updateResponse struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}
	respondWithJSON(w, http.StatusOK, updateResponse{
		Email: u.Email,
		Id:    u.Id,
	})
}
