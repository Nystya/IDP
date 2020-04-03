package api

import (
	"context"
	"google.golang.org/grpc"
	"idp/Jobs/models"
	"idp/Jobs/proto"
	"idp/Jobs/usecases"
	"log"
	"net"
	"os"
	"os/signal"
	pb "idp/Jobs/proto"
	"time"
)

type ServerConfig struct {
	JobManagerConf usecases.JobManagerConfig
	Transport string
	Port string
}

type JobServer struct {
	JobManager usecases.JobManager
}

func NewJobServer(conf *ServerConfig) *JobServer {
	return &JobServer{
		JobManager: usecases.NewJobManagerImpl(conf.JobManagerConf),
	}
}

func (js *JobServer) GetAllServices(*proto.Void, proto.Jobs_GetAllServicesServer) error {
	panic("implement me")
}

func (js *JobServer) GetAllSkillCategories(*proto.Skill, proto.Jobs_GetAllSkillCategoriesServer) error {
	panic("implement me")
}

func (js *JobServer) GetAllSkills(*proto.SkillCategory, proto.Jobs_GetAllSkillsServer) error {
	panic("implement me")
}

func (js *JobServer) PostJob(ctx context.Context, job *proto.CreateJob) (*proto.Job, error) {
	jobData := models.Job{
		ID:           "",
		EUID:         job.GetEUID().GetID(),
		Title:        job.GetTitle(),
		Service:      models.Service{
			ID: job.GetService().GetID().GetID(),
			Service: job.GetService().GetService(),
		},
		Category:     models.SkillCategory{
			ID: job.GetCategory().GetID().GetID(),
			Category: job.GetCategory().GetCategory(),
		},
		Experience:   job.GetExp(),
		Wage:         job.GetWage(),
		Places:       job.GetPlaces(),
		Description:  job.GetDescription(),
		Skills:       nil,
		PostTime:     time.Now(),
		ERating:      0,
		NrCandidates: 0,
		MoneySpent:   0,
	}

	if err := js.JobManager.AddJob(jobData); err != nil {
		return nil, err
	}

	return nil, nil
}

func (js *JobServer) GetJobs(*proto.ID, proto.Jobs_GetJobsServer) error {
	panic("implement me")
}

func (js *JobServer) ApplyForJob(context.Context, *proto.JobApplication) (*proto.ErrorCode, error) {
	panic("implement me")
}

func (js *JobServer) GetApplicants(*proto.ID, proto.Jobs_GetApplicantsServer) error {
	panic("implement me")
}

func (js *JobServer) SelectForJob(context.Context, *proto.JobSelection) (*proto.ErrorCode, error) {
	panic("implement me")
}

func (js *JobServer) GetAcceptedFreelancers(*proto.ID, proto.Jobs_GetAcceptedFreelancersServer) error {
	panic("implement me")
}

func (js *JobServer) GetAcceptedJobs(*proto.ID, proto.Jobs_GetAcceptedJobsServer) error {
	panic("implement me")
}

func (js *JobServer) GetHistoryJobs(*proto.ID, proto.Jobs_GetHistoryJobsServer) error {
	panic("implement me")
}

func (js *JobServer) FinishJob(context.Context, *proto.ID) (*proto.ErrorCode, error) {
	panic("implement me")
}

func RunServer(ctx context.Context, conf *ServerConfig) error {
	jobServer := NewJobServer(conf)

	listener, err := net.Listen(conf.Transport, ":"+conf.Port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterJobsServer(grpcServer, jobServer)

	// Graceful stop in case of os signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Shuting down gRPC server...")
			grpcServer.GracefulStop()

			<- ctx.Done()
		}
	}()

	log.Println("Starting gRPC server...")
	return grpcServer.Serve(listener)
}
