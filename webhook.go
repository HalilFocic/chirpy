package main

import (
	"encoding/json"
	"net/http"
)

type WebhookRequestBody struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := WebhookRequestBody{}

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while decoding body")
		return
	}
	if body.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	user, err := cfg.DB.SetUserToPremium(body.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user with that id")
		return
	}
	respondWithJSON(w, http.StatusNoContent, UserDTO{
		Email:       user.Email,
		Id:          user.Id,
		IsChirpyRed: user.IsChirpyRed,
	})

}
