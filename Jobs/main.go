package main

import (
	"context"
	"idp/Jobs/api"
	"idp/Jobs/config"
	"log"
)

const (
	confFile="config.json"
)

func main() {
	ctx := context.Background()
	conf, err := config.NewMicroserviceConfig(confFile).GetConfig()
	if err != nil {
		log.Println("Could not get server config: " + err.Error())
		panic(1)
	}

	if err := api.RunServer(ctx, conf); err != nil {
		log.Println("Could not start server: " + err.Error())
		panic(2)
	}
}

//func main() {
//	conf, err := config.NewMicroserviceConfig(confFile).GetConfig()
//	if err != nil {
//		log.Println("Could not get server config: " + err.Error())
//		panic(1)
//	}
//
//	jm := usecases.NewJobManagerImpl(conf.JobManagerConf)
//	//job := &models.Job{
//	//	ID:           "1234",
//	//	EUID:         "1",
//	//	Title:        "Accountant",
//	//	ServiceCategories:      []*models.ServiceCategory{{ID:"1"}, {ID: "2"}},
//	//	SkillCategories:     []*models.SkillCategory{{ID: "3"}, {ID: "4"}, {ID: "7"}},
//	//	Experience:   "VHDL experience needed.",
//	//	Wage:         2500,
//	//	Places:       1,
//	//	Description:  "We are looking for an experienced VHDL developer",
//	//	Skills:       nil,
//	//}
//
//	//err = jm.AddJob(job)
//
//	jobs, err := jm.GetJobs(&models.Filter{
//		ID:      "",
//		Title:   "Developer",
//		WageMin: 3000,
//		ERating: 0,
//	})
//	if err != nil {
//		log.Println("Error creating job: ", err.Error())
//		return
//	}
//
//	for _, job := range jobs {
//		log.Println(job)
//	}
//
//}