package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/HalilFocic/chirpy/internal/database"
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

	c, err := cfg.DB.CreateChirp(cleanedMessage)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating chips")
		return
	}

	respondWithJSON(w, http.StatusCreated, database.Chirp{
		Id:   c.Id,
		Body: c.Body,
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
			Id:   dbChirp.Id,
			Body: dbChirp.Body,
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
