package api

//go:generate protoc --plugin=/home/skanda/go/bin/protoc-gen-go --proto_path=../proto/ --go_out=plugins=grpc:../proto/ ../proto/api.proto

import (
	"context"
	"google.golang.org/grpc"
	"idp/Jobs/models"
	pb "idp/Jobs/proto"
	"idp/Jobs/usecases"
	"log"
	"net"
	"os"
	"os/signal"
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
	pbJob := &pb.Job{
		ID:                   &pb.ID{ID: job.ID},
		EUID:                 &pb.ID{ID: job.EUID},
		ServiceCategories:    make([]*pb.ServiceCategory, 0),
		SkillCategories:	  make([]*pb.SkillCategory, 0),
		Wage:                 job.Wage,
		Places:               job.Places,
		Title:                job.Title,
		Exp:                  job.Experience,
		Description:          job.Description,
		PostTime:             job.PostTime,
		NrOfCandidates:       int32(job.NrCandidates),
		EmployerRating:       job.ERating,
		MoneySpent:           float32(job.MoneySpent),
	}

	for _, sc := range job.ServiceCategories {
		pbJob.ServiceCategories = append(pbJob.ServiceCategories, &pb.ServiceCategory{
			ID:                   &pb.ID{ID: sc.ID},
			Service:              sc.Category,
		})
	}

	for _, skc := range job.SkillCategories {
		pbJob.SkillCategories = append(pbJob.SkillCategories, &pb.SkillCategory{
			ID:                   &pb.ID{ID: skc.ID},
			Category:             skc.Category,
		})
	}

	return pbJob
}

func convertFreelancerProfileToProto(freelancer *models.Freelancer) *pb.FreelancerProfile {
	pbFreelancer := &pb.FreelancerProfile{
		FUID:            &pb.ID{ID: freelancer.ID},
		Phone:           freelancer.Phone,
		LastName:        freelancer.LastName,
		FirstName:       freelancer.FirstName,
		Rating:          freelancer.Rating,
		Balance:         freelancer.Balance,
		Description:     freelancer.Description,
		Photo:           freelancer.Photo,
		SkillCategories: make([]*pb.SkillCategory, 0),
		Skills:          make([]*pb.Skill, 0),
	}

	for _, skc := range freelancer.SkillCategories {
		pbFreelancer.SkillCategories = append(pbFreelancer.SkillCategories, &pb.SkillCategory{
			ID:                   &pb.ID{ID: skc.ID},
			Category:             skc.Category,
		})
	}

	for _, sk := range freelancer.Skills {
		pbFreelancer.Skills = append(pbFreelancer.Skills, &pb.Skill{
			ID:                   &pb.ID{ID: sk.ID},
			Skill: 				  sk.Skill,
		})
	}

	return pbFreelancer
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
			Service:              sc.Category,
		}
		log.Println(pbsc)
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
		ID:       pbsc.GetID().GetID(),
		Category: pbsc.GetService(),
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
		ID:       pbskc.GetID().GetID(),
		Category: pbskc.GetCategory(),
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
		EUID:         job.GetEUID().GetID(),
		Title:        job.GetTitle(),
		ServiceCategories: nil,
		SkillCategories: nil,
		Experience:   job.GetExp(),
		Wage:         job.GetWage(),
		Places:       job.GetPlaces(),
		Description:  job.GetDescription(),
		Skills:       nil,
		PostTime:     time.Now().String(),
	}

	for _, sc := range job.GetServiceCategories() {
		jobData.ServiceCategories = append(jobData.ServiceCategories, &models.ServiceCategory{
			ID:       sc.ID.ID,
			Category: sc.Service,
		})
	}

	for _, skc := range job.GetSkillCategories() {
		jobData.SkillCategories = append(jobData.SkillCategories, &models.SkillCategory{
			ID: skc.ID.ID,
			Category: skc.Category,
		})
	}

	newJob, err := js.JobManager.AddJob(jobData)
	if err != nil {
		return nil, err
	}

	pbJob := convertJobToProto(newJob)

	return pbJob, nil
}

func (js *JobServer) GetJob(ctx context.Context, id *pb.ID) (*pb.Job, error) {
    log.Println("Here")
	jobData, err := js.JobManager.GetJob(id.GetID())
	if err != nil {
		return nil, err
	}

	return convertJobToProto(jobData), nil
}

func (js *JobServer) GetJobs(filter *pb.Filter, server pb.Jobs_GetJobsServer) error {
	errs := make([]string, 0)

	qfilter := &models.Filter{
		ID:      filter.GetID().GetID(),
		Title:   filter.GetTitle(),
		WageMin: filter.GetWageMin(),
		ERating: filter.GetEmployerRating(),
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

func (js *JobServer) ApplyForJob(ctx context.Context, ja *pb.JobApplication) (*pb.ErrorCode, error) {
	if err := js.JobManager.ApplyForJob(ja.GetJID().GetID(), ja.GetFUID().GetID()); err != nil {
		return &pb.ErrorCode{Err: 500, Msg: "Could not create link"}, err
	}

	return &pb.ErrorCode{Err: 200, Msg: ""}, nil
}

func (js *JobServer) GetApplicants(jid *pb.ID, server pb.Jobs_GetApplicantsServer) error {
	errs := make([]string, 0)

	freelancers, err := js.JobManager.GetApplicants(jid.GetID())
	if err != nil {
		return err
	}

	for _, freelancer := range freelancers {
		if err = server.Send(convertFreelancerProfileToProto(freelancer)); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) SelectForJob(ctx context.Context, jobSel *pb.JobSelection) (*pb.ErrorCode, error) {
	if err := js.JobManager.SelectFreelancerForJob(jobSel.GetJID().GetID(), jobSel.GetFUID().GetID()); err != nil {
		return &pb.ErrorCode{Err: 500, Msg: "Could not create link"}, err
	}

	return &pb.ErrorCode{Err: 200, Msg: ""}, nil
}

func (js *JobServer) GetAcceptedFreelancers(jid *pb.ID, server pb.Jobs_GetAcceptedFreelancersServer) error {
	errs := make([]string, 0)

	freelancers, err := js.JobManager.GetAcceptedFreelancers(jid.GetID())
	if err != nil {
		return err
	}

	for _, freelancer := range freelancers {
		if err = server.Send(convertFreelancerProfileToProto(freelancer)); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) GetAcceptedJobs(fid *pb.ID, server pb.Jobs_GetAcceptedJobsServer) error {
	errs := make([]string , 0)

    log.Println(fid.GetID())

	jobs, err := js.JobManager.GetAcceptedJobsForFreelancer(fid.GetID())
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err = server.Send(convertJobToProto(job)); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) GetFreelancerHistoryJobs(id *pb.ID, server pb.Jobs_GetFreelancerHistoryJobsServer) error {
	errs := make([]string, 0)

	jobs, err := js.JobManager.GetHistoryJobs(id.GetID(), models.ActorTypeFreelancer)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err = server.Send(convertJobToProto(job)); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) GetEmployerHistoryJobs(id *pb.ID, server pb.Jobs_GetEmployerHistoryJobsServer) error {
	errs := make([]string, 0)

	jobs, err := js.JobManager.GetHistoryJobs(id.GetID(), models.ActorTypeEmployer)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err = server.Send(convertJobToProto(job)); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return SendError{Err: strings.Join(errs, "\n")}
	}

	return nil
}

func (js *JobServer) FinishJob(ctx context.Context, jid *pb.ID) (*pb.ErrorCode, error) {
	if err := js.JobManager.FinishJob(jid.GetID()); err != nil {
	    return &pb.ErrorCode{Err: 500, Msg: "Could not create link"}, err
	}

	return &pb.ErrorCode{Err: 200, Msg: ""}, nil
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
