package service

import (
	"context"

	"github.com/danilomarques1/grpc-gopm/pb"
	"github.com/danilomarques1/grpc-gopm/server/model"
	"github.com/google/uuid"
)

type PasswordServer struct {
	pb.UnimplementedPasswordServer
	repository model.PasswordRepository
}

func NewPasswordServer(repository model.PasswordRepository) *PasswordServer {
	return &PasswordServer{repository: repository}
}

func (s *PasswordServer) SavePassword(ctx context.Context, in *pb.CreatePasswordRequest) (*pb.CreatePasswordResponse, error) {
	password := &model.Password{
		Id:  uuid.NewString(),
		Key: in.GetKey(),
		Pwd: in.GetPassword(),
	}
	if err := s.repository.Save(password); err != nil {
		return nil, err
	}

	return &pb.CreatePasswordResponse{OK: true}, nil
}

func (s *PasswordServer) FindAllKeys(ctx context.Context, in *pb.Empty) (*pb.Keys, error) {
	keys, err := s.repository.FindAllKeys()
	if err != nil {
		return nil, err
	}
	response := &pb.Keys{
		Keys: keys,
	}

	return response, nil
}
