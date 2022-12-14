package service

import (
	"context"
	"errors"

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
	if _, err := s.repository.FindPassword(in.GetKey()); err == nil {
		return nil, errors.New("Key already used")
	}

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

func (s *PasswordServer) FindPassword(ctx context.Context, in *pb.FindPasswordRequest) (*pb.PasswordResponse, error) {
	password, err := s.repository.FindPassword(in.GetKey())
	if err != nil {
		return nil, err
	}
	response := &pb.PasswordResponse{
		Id:       password.Id,
		Key:      password.Key,
		Password: password.Pwd,
	}

	return response, nil
}

func (s *PasswordServer) DeletePassword(ctx context.Context, in *pb.DeletePasswordRequest) (*pb.DeletePasswordResponse, error) {
	if err := s.repository.Delete(in.GetKey()); err != nil {
		return nil, err
	}
	return &pb.DeletePasswordResponse{
		OK: true,
	}, nil
}

func (s *PasswordServer) UpdatePassword(ctx context.Context, in *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	password, err := s.repository.FindPassword(in.GetKey())
	if err != nil {
		return nil, errors.New("Key not found")
	}

	password.Key = in.GetKey()
	password.Pwd = in.GetPassword()

	if err := s.repository.Update(password); err != nil {
		return nil, err
	}

	return &pb.UpdatePasswordResponse{OK: true}, nil
}
