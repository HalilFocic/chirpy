package main

import (
	"fmt"
	"net/http"

	"github.com/HalilFocic/chirpy/internal/auth"
)

type RefreshResposne struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}
	user, err := cfg.DB.GetUserByRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find refresh in db")
		return
	}
	newToken, err := auth.MakeJWT(user.Id, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error making JWT")
		return
	}
	respondWithJSON(w, http.StatusOK, RefreshResposne{
		Token: newToken,
	})
	return
}
func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}
	err = cfg.DB.DeleteUserToken(token)
	if err != nil {
		fmt.Println("bio error")
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	fmt.Println("Nije bio error")
	respondWithJSON(w, http.StatusNoContent, nil)
	return
}
