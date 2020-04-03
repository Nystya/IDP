package database

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
)
import "idp/Jobs/models"

const (
	createJobQuery = "CREATE (a: Job {$prop})"
)

type Database interface {
	CreateJob(job models.Job) error

}

type Neo4jDatabase struct {
	Driver neo4j.Driver
}

func NewNeo4jDatabase(uri, username, password string) *Neo4jDatabase {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Println("Could not connect to Neo4j: ", err.Error())
		return nil
	}

	return &Neo4jDatabase{
		Driver: driver,
	}
}

func (nj *Neo4jDatabase) CreateJob(job models.Job) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		_, err = tx.Run(createJobQuery, map[string]interface{}{
			"prop": map[string]interface{}{
				"ID": job.ID,
				"descr": job.Description,
				"wage": job.Wage,
				"places": job.Places,
				"title": job.Title,
				"exp": job.Experience,
				"postTime": job.PostTime,
			},
		})

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}


