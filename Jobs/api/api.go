package api

import (
	"context"
	"google.golang.org/grpc"
	"idp/Jobs/models"
	"idp/Jobs/usecases"
	"log"
	"net"
	"os"
	"os/signal"
	pb "idp/Jobs/proto"
	"strings"
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

type SendError struct {
	Err string
}

func (err SendError) Error() string {
	return err.Err
}

func convertJobToProto(job *models.Job) *pb.Job {
	return &pb.Job{
		ID:                   &pb.ID{ID: job.ID},
		EUID:                 &pb.ID{ID: job.EUID},
		Service:              &pb.ServiceCategory{
			ID:                   	&pb.ID{ID: job.Service.ID},
			Service:      			job.Service.Service,
		},
		Category:             &pb.SkillCategory{
			ID:                   &pb.ID{ID: job.Category.ID},
			Category:             job.Category.Category,
		},
		Wage:                 job.Wage,
		Places:               job.Places,
		Title:                job.Title,
		Exp:                  job.Experience,
		Description:          job.Description,
		PostTime:             job.PostTime.String(),
		NrOfCandidates:       int32(job.NrCandidates),
		EmployerRating:       job.ERating,
		MoneySpent:           float32(job.MoneySpent),
	}
}

func (js *JobServer) GetAllServices(void *pb.Void, server pb.Jobs_GetAllServicesServer) error {
	errs := make([]string, 0)

	serviceCategories, err := js.JobManager.GetAllServiceCategories()
	if err != nil {
		return err
	}

	for _, sc := range serviceCategories {
		pbsc := &pb.ServiceCategory{
			ID:                   &pb.ID{ID: sc.ID},
			Service:              sc.Service,
		}
		if err := server.Send(pbsc); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) GetAllSkillCategories(pbsc *pb.ServiceCategory, server pb.Jobs_GetAllSkillCategoriesServer) error {
	errs := make([]string, 0)

	sc := &models.ServiceCategory{
		ID:      pbsc.ID.ID,
		Service: pbsc.Service,
	}

	skillCategories, err := js.JobManager.GetSkillCategoriesByServiceCategory(sc)
	if err != nil {
		return err
	}

	for _, skc := range skillCategories {
		pbskc := &pb.SkillCategory{
			ID:                   &pb.ID{ID: skc.ID},
			Category:             skc.Category,
		}
		if err := server.Send(pbskc); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) GetAllSkills(pbskc *pb.SkillCategory, server pb.Jobs_GetAllSkillsServer) error {
	errs := make([]string, 0)

	skc := &models.SkillCategory{
		ID:       pbskc.ID.ID,
		Category: pbskc.Category,
	}

	skills, err := js.JobManager.GetSkillByCategory(skc)
	if err != nil {
		return err
	}

	for _, sk := range skills {
		pbsc := &pb.Skill{
			ID:                   &pb.ID{ID: sk.ID},
			Category:             pbskc,
			Skill:                sk.Skill,
		}
		if err := server.Send(pbsc); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) PostJob(ctx context.Context, job *pb.CreateJob) (*pb.Job, error) {
	jobData := &models.Job{
		ID:           "",
		EUID:         job.GetEUID().GetID(),
		Title:        job.GetTitle(),
		Service:      models.ServiceCategory{
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

func (js *JobServer) GetJob(ctx context.Context, id *pb.ID) (*pb.Job, error) {
	jobData, err := js.JobManager.GetJob(id.ID)
	if err != nil {
		return nil, err
	}

	return convertJobToProto(jobData), nil
}

func (js *JobServer) GetJobs(filter *pb.Filter, server pb.Jobs_GetJobsServer) error {
	errs := make([]string, 0)

	qfilter := &models.Filter{
		ID:      filter.ID.ID,
		Title:   filter.Title,
		WageMin: filter.WageMin,
		ERating: filter.EmployerRating,
	}

	jobsData, err := js.JobManager.GetJobs(qfilter)
	if err != nil {
		return err
	}

	for _, jobData := range jobsData {
		if err := server.Send(convertJobToProto(jobData)); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) ApplyForJob(context.Context, *pb.JobApplication) (*pb.ErrorCode, error) {
	panic("implement me")
}

func (js *JobServer) GetApplicants(*pb.ID, pb.Jobs_GetApplicantsServer) error {
	panic("implement me")
}

func (js *JobServer) SelectForJob(context.Context, *pb.JobSelection) (*pb.ErrorCode, error) {
	panic("implement me")
}

func (js *JobServer) GetAcceptedFreelancers(*pb.ID, pb.Jobs_GetAcceptedFreelancersServer) error {
	panic("implement me")
}

func (js *JobServer) GetAcceptedJobs(*pb.ID, pb.Jobs_GetAcceptedJobsServer) error {
	panic("implement me")
}

func (js *JobServer) GetHistoryJobs(*pb.ID, pb.Jobs_GetHistoryJobsServer) error {
	panic("implement me")
}

func (js *JobServer) FinishJob(context.Context, *pb.ID) (*pb.ErrorCode, error) {
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
