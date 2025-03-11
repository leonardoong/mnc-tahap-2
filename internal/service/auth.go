package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leonardoong/e-wallet/config"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/repository"
	"github.com/leonardoong/e-wallet/internal/utils"
)

type IAuthService interface {
	Register(req *entity.RegisterUserRequest) (user *entity.User, err error)
	Login(req *entity.LoginRequest) (resp *entity.LoginResponse, err error)
	UpdateProfile(req *entity.UpdateProfileRequest) (resp *entity.UpdateProfileResponse, err error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type authService struct {
	config         *config.Config
	userRepository repository.IUserRepository
}

func NewAuthService(config *config.Config, userRepo repository.IUserRepository) IAuthService {
	return &authService{
		config:         config,
		userRepository: userRepo,
	}
}

func (s *authService) Register(req *entity.RegisterUserRequest) (user *entity.User, err error) {
	existingUser, _ := s.userRepository.FindByPhoneNumber(req.PhoneNumber)
	if existingUser != nil {
		return nil, errors.New("Phone number already registered")
	}

	hashedPin, err := utils.HashPin(req.Pin)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user = &entity.User{
		UserID:      uuid.New().String(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Pin:         hashedPin,
		Address:     req.Address,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.userRepository.Register(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req *entity.LoginRequest) (resp *entity.LoginResponse, err error) {
	user, err := s.userRepository.FindByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, errors.New("Phone Number and PIN doesn't match.")
	}

	if user == nil {
		return nil, errors.New("Phone Number and PIN doesn't match.")
	}

	if !utils.CheckPinHash(req.Pin, user.Pin) {
		return nil, errors.New("Phone Number and PIN doesn't match.")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone_number": req.PhoneNumber,
		"user_id":      user.UserID,
		"exp":          time.Now().Add(time.Minute * 60).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone_number": req.PhoneNumber,
		"user_id":      user.UserID,
		"exp":          time.Now().Add(time.Hour * time.Duration(s.config.JWTExpiryHours)).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return
	}

	return &entity.LoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, err
}

func (s *authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (s *authService) UpdateProfile(req *entity.UpdateProfileRequest) (resp *entity.UpdateProfileResponse, err error) {
	now := time.Now()

	user := entity.User{
		UserID:    req.UserID,
		Address:   req.Address,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UpdatedAt: now,
	}

	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	return &entity.UpdateProfileResponse{
		UserID:    req.UserID,
		Address:   req.Address,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UpdatedAt: now,
	}, err
}
