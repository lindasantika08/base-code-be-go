package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"go-base-project/config"
	"go-base-project/internal/model"
	"go-base-project/internal/repository"
	"go-base-project/internal/utils"
)

type UserService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	GetByID(ctx context.Context, id int64) (*model.UserResponse, error)
	GetAll(ctx context.Context, query model.PaginationQuery) ([]*model.UserResponse, *model.Meta, error)
	Update(ctx context.Context, id int64, req model.UpdateUserRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository) UserService {
	cfg, _ := config.Load()
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *userService) Register(ctx context.Context, req model.RegisterRequest) (*model.UserResponse, error) {

	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal cek email: %w", err)
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// ── PROSES DATA ──────────────────────────────────────────

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password: %w", err)
	}

	user := &model.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "user",
		IsActive:  true,
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("gagal simpan user: %w", err)
	}

	response := created.ToResponse()
	return &response, nil
}

func (s *userService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal cek user: %w", err)
	}

	if user == nil {
		return nil, errors.New("email atau password salah")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	if !user.IsActive {
		return nil, errors.New("akun tidak aktif, hubungi admin")
	}

	accessToken, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpiryHours)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %w", err)
	}

	refreshToken, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshExpHours)
	if err != nil {
		return nil, fmt.Errorf("gagal generate refresh token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.cfg.JWT.ExpiryHours * 3600,
		User:         user.ToResponse(),
	}, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) GetAll(ctx context.Context, query model.PaginationQuery) ([]*model.UserResponse, *model.Meta, error) {
	query.DefaultPagination()

	users, total, err := s.userRepo.FindAll(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal ambil users: %w", err)
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		resp := user.ToResponse()
		responses[i] = &resp
	}

	meta := &model.Meta{
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(query.Limit))),
	}

	return responses, meta, nil
}

func (s *userService) Update(ctx context.Context, id int64, req model.UpdateUserRequest) (*model.UserResponse, error) {
	existing, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if req.Email != "" && req.Email != existing.Email {
		emailExists, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if emailExists != nil {
			return nil, errors.New("email sudah dipakai")
		}
	}

	updated, err := s.userRepo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("gagal update user: %w", err)
	}

	response := updated.ToResponse()
	return &response, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	existing, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("user tidak ditemukan")
	}

	return s.userRepo.Delete(ctx, id)
}

var _ UserService = (*userService)(nil)
var _ = time.Now 
