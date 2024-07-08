package main

import (
	"encoding/json"
	"github.com/HalilFocic/chirpy/internal/auth"
	"net/http"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserDTO struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while decoding body")
		return
	}
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user with that email")
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	respondWithJSON(w, http.StatusOK, UserDTO{
		Id:    user.Id,
		Email: user.Email,
	})
	return

}
