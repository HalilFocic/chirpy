package database

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
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

func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStruct.Chirps) + 1
	c := Chirp{
		AuthorId: userID,
		Id:       id,
		Body:     body,
	}
	dbStruct.Chirps[id] = c
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return c, nil

}
func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	val, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}
	return val, nil

}
func (db *DB) DeleteChirp(chirpId int) (error ){
	dbStruct, err := db.loadDB()
	if err != nil {
		return  err
	}
    delete(dbStruct.Chirps,chirpId)
	return nil


}
