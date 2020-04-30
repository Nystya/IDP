package main

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
)

//go:generate easytags $GOFILE

type Movie struct {
	MovieId    string   `json:"movieId"`
	Title      string   `json:"title"`
	Year       int      `json:"year"`
	Runtime    int      `json:"runtime"`
	TmdbId     string   `json:"tmdbId"`
	ImdbId     string   `json:"imdbId"`
	ImdbRating int      `json:"imdbRating"`
	ImdbVotes  int      `json:"imdbVotes"`
	Plot       string   `json:"plot"`
	Poster     string   `json:"poster"`
	Released   string   `json:"released"`
	Languages  []string `json:"languages"`
	Countries  []string `json:"countries"`
}

const (
	getJobByIDQuery = "MATCH(m:Movie) WHERE(m.movieId=$id OR m.movieId=$id2) RETURN properties(m)"
)

func main() {
	driver, err := neo4j.NewDriver("bolt://100.24.206.62:33829", neo4j.BasicAuth("neo4j", "nets-refunds-preference", ""))
	if err != nil {
		log.Println("Could not connect to Neo4j: ", err.Error())
		panic(err.Error())
	}

	session, err := driver.Session(neo4j.AccessModeRead)
	if err != nil {
		panic(err.Error())
	}

	defer session.Close()

	movieData, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		movieData, err := tx.Run(getJobByIDQuery, map[string]interface{}{
			"id":  "2713",
			"id2": "2714",
		})

		if err != nil {
			return nil, err
		}

		return movieData, nil
	})

	if err != nil {
		panic(err.Error())
	}

	for movieData.(neo4j.Result).Next() {
		movie := &Movie{}

		mapstructure.Decode(movieData.(neo4j.Result).Record().GetByIndex(0), movie)

		log.Println(movie.TmdbId)
	}

}
