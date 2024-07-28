package application

import (
	"context"
	"fmt"
	pb "github.com/synthao/meetme/gen/go/imgproc"
	"github.com/synthao/meetme/internal/user/domain"
)

type Service struct {
	userRepo      domain.Repository
	imgprocClient pb.ImageProcessingServiceClient
}

func NewService(userRepo domain.Repository, imgprocClient pb.ImageProcessingServiceClient) *Service {
	return &Service{userRepo: userRepo, imgprocClient: imgprocClient}
}

func (s *Service) Create(dto CreateUserDTO) (int, error) {
	user := &domain.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Gender:    domain.Gender(dto.Gender),
		BirthDate: dto.BirthDate,
	}

	id, err := s.userRepo.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetByID(id int) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *Service) Update(user *domain.User) error {
	err := s.userRepo.Update(user)
	if err != nil {
		return err
	}

	image, err := s.imgprocClient.ProcessImage(context.Background(), &pb.ProcessImageRequest{
		Path: "374h65f834765gf8234fh8453.jpeg",
		X:    0,
		Y:    0,
		W:    480,
		H:    640,
	})
	if err != nil {
		return err
	}

	fmt.Println(">>> grpc response", image.Small, image.Medium, image.Large) // TODO rm

	return nil
}

func (s *Service) GetList(limit, offset int) ([]domain.User, error) {
	return s.userRepo.GetList(limit, offset)
}

func (s *Service) Delete(id int) error {
	return s.userRepo.Delete(id)
}
