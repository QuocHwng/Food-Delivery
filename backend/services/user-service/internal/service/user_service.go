package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"user-service/internal/model"
	"user-service/internal/repository"
)

type UserService struct {
	userRepo    *repository.UserRepository
	addressRepo *repository.AddressRepository
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		userRepo:    repository.NewUserRepository(db),
		addressRepo: repository.NewAddressRepository(db),
	}
}

// ─── Profile ──────────────────────────────────────────────────────────────────

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("GetProfile: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req model.UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("UpdateProfile find: %w", err)
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("UpdateProfile update: %w", err)
	}
	return user, nil
}

// ─── Addresses ────────────────────────────────────────────────────────────────

func (s *UserService) GetAddresses(ctx context.Context, userID uuid.UUID) ([]model.UserAddress, error) {
	addresses, err := s.addressRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAddresses: %w", err)
	}
	return addresses, nil
}

func (s *UserService) CreateAddress(ctx context.Context, userID uuid.UUID, req model.CreateAddressRequest) (*model.UserAddress, error) {
	addr := &model.UserAddress{
		ID:     uuid.New(),
		UserID: userID,
		Street: req.Street,
		Lat:    req.Lat,
		Lng:    req.Lng,
	}
	if req.Label != "" {
		addr.Label = &req.Label
	}
	if req.District != "" {
		addr.District = &req.District
	}
	if req.City != "" {
		addr.City = &req.City
	}

	// Địa chỉ đầu tiên → tự động là mặc định
	count, _ := s.addressRepo.CountByUserID(ctx, userID)
	if count == 0 {
		addr.IsDefault = true
	}

	if err := s.addressRepo.Create(ctx, addr); err != nil {
		return nil, fmt.Errorf("CreateAddress: %w", err)
	}
	return addr, nil
}

func (s *UserService) UpdateAddress(ctx context.Context, userID, addrID uuid.UUID, req model.UpdateAddressRequest) (*model.UserAddress, error) {
	addr, err := s.addressRepo.FindByID(ctx, addrID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("address not found")
		}
		return nil, fmt.Errorf("UpdateAddress find: %w", err)
	}

	if req.Label != "" {
		addr.Label = &req.Label
	}
	if req.Street != "" {
		addr.Street = req.Street
	}
	if req.District != "" {
		addr.District = &req.District
	}
	if req.City != "" {
		addr.City = &req.City
	}
	addr.Lat = req.Lat
	addr.Lng = req.Lng

	if err := s.addressRepo.Update(ctx, addr); err != nil {
		return nil, fmt.Errorf("UpdateAddress update: %w", err)
	}
	return addr, nil
}

func (s *UserService) DeleteAddress(ctx context.Context, userID, addrID uuid.UUID) error {
	return s.addressRepo.Delete(ctx, addrID, userID)
}

func (s *UserService) SetDefaultAddress(ctx context.Context, userID, addrID uuid.UUID) error {
	// Kiểm tra địa chỉ thuộc về user
	if _, err := s.addressRepo.FindByID(ctx, addrID, userID); err != nil {
		return fmt.Errorf("address not found")
	}
	return s.addressRepo.SetDefault(ctx, addrID, userID)
}
