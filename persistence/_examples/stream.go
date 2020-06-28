package persistence

import (
	// "io/ioutil"
	"context"
	"encoding/json"
	"io"

	"github.com/ortymid/projector/projector"
)

// type jsonData []jsonUser

type jsonUser struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Boards []jsonBoard `json:"boards"`
}

type jsonBoard struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StreamUserRepo struct {
	r io.ReadWriter
}

func NewStreamUserRepo(r io.ReadWriter) StreamUserRepo {
	return StreamUserRepo{r: r}
}

func (repo StreamUserRepo) All(ctx context.Context) ([]projector.User, error) {
	var data []jsonUser

	dec := json.NewDecoder(repo.r)
	if err := dec.Decode(data); err == io.EOF {
		// log.Println("EOF:", data)
		// return nil, errors.New("corrupted data")
	} else if err != nil {
		return nil, err
	}

	users := []projector.User{}

	for _, u := range data {
		user, err := projector.NewUser(u.ID, u.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo StreamUserRepo) Create(ctx context.Context, u projector.User) (err error) {
	data := jsonUser{
		ID:     u.ID,
		Name:   u.Name,
		Boards: []jsonBoard{},
	}

	enc := json.NewEncoder(repo.r)
	err = enc.Encode(&data)

	return err
}
