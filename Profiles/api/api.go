package api

//go:generate protoc --plugin=/home/skanda/go/bin/protoc-gen-go --proto_path=../proto/ --go_out=plugins=grpc:../proto/ ../proto/api.proto

import (
	"context"
	"google.golang.org/grpc"
	pb "idp/Profiles/proto"
	"idp/Profiles/models"
	"idp/Profiles/usecases"
	"log"
	"net"
	"os"
	"os/signal"
)

type ServerConfig struct {
	ProfileManagerConf usecases.ProfileManagerConfig
	Transport string
	Port string
}

type ProfileServer struct {
	ProfileManager usecases.ProfileManager
}

type SendError struct {
	Err string
}

func (err SendError) Error() string {
	return err.Err
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

func convertEmployerProfileToProto(employer *models.Employer) *pb.EmployerProfile {
	pbEmployer := &pb.EmployerProfile{
		EUID:      	&pb.ID{ID: employer.ID},
		Phone:     	employer.Phone,
		LastName:   employer.LastName,
		FirstName:  employer.FirstName,
		Rating:     employer.Rating,
		JobsPosted: employer.JobsPosted,
		MoneySpent: employer.MoneySpent,
	}

	return pbEmployer
}

func NewProfileServer(conf *ServerConfig) *ProfileServer {
	return &ProfileServer{
		ProfileManager: usecases.NewProfileManagerImpl(conf.ProfileManagerConf),
	}
}

func (p ProfileServer) CreateEmployerProfile(ctx context.Context, request *pb.EditEmployerProfileRequest) (*pb.ErrorCode, error) {
	editProfileRequest := &models.EditEmployerRequest{
		ID:        request.GetEUID().GetID(),
		Phone:     request.GetPhone(),
		LastName:  request.GetLastName(),
		FirstName: request.GetFirstName(),
	}

	if err := p.ProfileManager.CreateEmployerProfile(editProfileRequest); err != nil {
		return &pb.ErrorCode{Err: 500, Msg:""}, err
	}

	return &pb.ErrorCode{Err: 200, Msg:""}, nil
}

func (p ProfileServer) CreateFreelancerProfile(ctx context.Context, request *pb.EditFreelancerProfileRequest) (*pb.ErrorCode, error) {
	editProfileRequest := &models.EditFreelancerRequest{
		ID:              request.GetFUID().GetID(),
		Phone:           request.GetPhone(),
		LastName:        request.GetLastName(),
		FirstName:       request.GetFirstName(),
		Description:     request.GetDescription(),
		Skills:          make([]*models.Skill, 0),
		SkillCategories: make([]*models.SkillCategory, 0),
	}

	if err := p.ProfileManager.CreateFreelancerProfile(editProfileRequest); err != nil {
		return &pb.ErrorCode{Err: 500, Msg:""}, err
	}

	return &pb.ErrorCode{Err: 200, Msg:""}, nil
}

func (p ProfileServer) EditEmployerProfile(ctx context.Context, request *pb.EditEmployerProfileRequest) (*pb.ErrorCode, error) {
	editProfileRequest := &models.EditEmployerRequest{
		ID:        request.GetEUID().GetID(),
		Phone:     request.GetPhone(),
		LastName:  request.GetLastName(),
		FirstName: request.GetFirstName(),
	}

	if err := p.ProfileManager.EditEmployerProfile(editProfileRequest); err != nil {
		return &pb.ErrorCode{Err: 500, Msg:""}, err
	}

	return &pb.ErrorCode{Err: 200, Msg:""}, nil
}

func (p ProfileServer) EditFreelancerProfile(ctx context.Context, request *pb.EditFreelancerProfileRequest) (*pb.ErrorCode, error) {
	editProfileRequest := &models.EditFreelancerRequest{
		ID:              request.GetFUID().GetID(),
		Phone:           request.GetPhone(),
		LastName:        request.GetLastName(),
		FirstName:       request.GetFirstName(),
		Description:     request.GetDescription(),
		Skills:          make([]*models.Skill, 0),
		SkillCategories: make([]*models.SkillCategory, 0),
	}

	for _, pbskc := range request.GetSkillCategories() {
		skc := &models.SkillCategory{
			ID:       pbskc.GetID().GetID(),
			Category: pbskc.GetCategory(),
		}
		editProfileRequest.SkillCategories = append(editProfileRequest.SkillCategories, skc)
	}

	for _, pbsk := range request.GetSkills() {
		sk := &models.Skill{
			ID:       pbsk.GetID().GetID(),
			SCID:     models.SkillCategory{
				ID:       pbsk.GetCategory().GetID().GetID(),
				Category: pbsk.GetCategory().GetCategory(),
			},
			Skill: 	  pbsk.GetSkill(),
		}
		editProfileRequest.Skills = append(editProfileRequest.Skills, sk)
	}
	
	if err := p.ProfileManager.EditFreelancerProfile(editProfileRequest); err != nil {
		return &pb.ErrorCode{Err: 500, Msg:""}, err
	}

	return &pb.ErrorCode{Err: 200, Msg:""}, nil
}

func (p ProfileServer) GetEmployerProfile(ctx context.Context, id *pb.ID) (*pb.EmployerProfile, error) {
    log.Println(id.GetID())
	employer, err := p.ProfileManager.GetEmployerProfile(id.GetID())
	if err != nil {
		return nil, err
	}

	return convertEmployerProfileToProto(employer), nil
}

func (p ProfileServer) GetFreelancerProfile(ctx context.Context, id *pb.ID) (*pb.FreelancerProfile, error) {
	freelancer, err := p.ProfileManager.GetFreelancerProfile(id.GetID())
	if err != nil {
		return nil, err
	}

	return convertFreelancerProfileToProto(freelancer), nil
}

func RunServer(ctx context.Context, conf *ServerConfig) error {
	profileServer := NewProfileServer(conf)

	listener, err := net.Listen(conf.Transport, ":"+conf.Port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProfilesServer(grpcServer, profileServer)

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

