package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"user-service/internal/model"
	"user-service/internal/repository"
)

// ─── Sentinel errors ──────────────────────────────────────────────────────────

var (
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidCreds = errors.New("invalid email or password")
	ErrInvalidToken = errors.New("invalid or expired refresh token")
	ErrUserNotFound = errors.New("user not found")
)

// ─── Service ──────────────────────────────────────────────────────────────────

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
	jwtCfg    JWTConfig
}

func NewAuthService(db *sqlx.DB, jwtCfg JWTConfig) *AuthService {
	return &AuthService{
		userRepo:  repository.NewUserRepository(db),
		tokenRepo: repository.NewTokenRepository(db),
		jwtCfg:    jwtCfg,
	}
}

// Register đăng ký user mới, trả về cặp token
func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	// Kiểm tra email đã tồn tại chưa
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Tạo user
	user := &model.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.issueTokens(ctx, user)
}

// Login xác thực email/password, trả về cặp token
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCreds
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	return s.issueTokens(ctx, user)
}

// Refresh xoay vòng refresh token (rotation), trả về cặp token mới
func (s *AuthService) Refresh(ctx context.Context, rawToken string) (*model.AuthResponse, error) {
	hash := HashRefreshToken(rawToken)

	tokenRecord, err := s.tokenRepo.FindByHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Thu hồi token cũ trước khi cấp mới (rotation để tránh reuse)
	if err := s.tokenRepo.Revoke(ctx, hash); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	return s.issueTokens(ctx, user)
}

// Logout thu hồi refresh token
func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	hash := HashRefreshToken(rawToken)
	return s.tokenRepo.Revoke(ctx, hash)
}

// issueTokens tạo token pair và lưu refresh hash vào DB
func (s *AuthService) issueTokens(ctx context.Context, user *model.User) (*model.AuthResponse, error) {
	pair, err := GenerateTokenPair(user.ID, user.Role, s.jwtCfg)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	refreshRecord := &model.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: pair.RefreshHash,
		ExpiresAt: pair.ExpiresAt,
	}
	if err := s.tokenRepo.Save(ctx, refreshRecord); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		User:         *user,
	}, nil
}
