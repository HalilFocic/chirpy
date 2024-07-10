package main

import (
	"encoding/json"
	"github.com/HalilFocic/chirpy/internal/auth"
	"github.com/HalilFocic/chirpy/internal/database"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handleAddChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return

	}
	cleanedMessage := BadWorkReplacement(params.Body)
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt decode parameters")
		return
	}
	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters")
		return
	}
	c, err := cfg.DB.CreateChirp(cleanedMessage, userID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating chips")
		return
	}

	respondWithJSON(w, http.StatusCreated, database.Chirp{
		Id:       c.Id,
		Body:     c.Body,
		AuthorId: c.AuthorId,
	})

}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChrips()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt retrieve chips!")
		return
	}
	chirps := []database.Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			Id:       dbChirp.Id,
			Body:     dbChirp.Body,
			AuthorId: dbChirp.AuthorId,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	respondWithJSON(w, http.StatusOK, chirps)

}
func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := strconv.Atoi(pathId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id was passed to request")
		return
	}

	chirp, err := cfg.DB.GetChirpById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp with that id doesnt exist")
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)

}
func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id was passed to request")
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnt decode parameters")
		return
	}
	userId, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Couldnt decode parameters")
		return
	}

	chirp, err := cfg.DB.GetChirpById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp with that id doesnt exist")
		return
	}
	if chirp.AuthorId != userId {
		respondWithError(w, http.StatusForbidden, "You dont have permisison to edit this!")
		return
	}
	err = cfg.DB.DeleteChirp(chirp.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while deleting")
		return

	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
