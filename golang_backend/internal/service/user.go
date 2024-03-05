package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"clean_arch/pkg"
	"clean_arch/pkg/phone"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func CodeGenerator() int {
	rand.Seed(time.Now().UnixNano())
	code := 1000 + rand.Intn(9999-1000)
	return code
}

type UserService struct {
	repo     repository.User
	log      *logrus.Logger
	tokenTTL string
	apiKey   string
}

func NewUserService(repo repository.User, log *logrus.Logger, tokenTTL string, apiKey string) *UserService {
	return &UserService{
		repo:     repo,
		log:      log,
		tokenTTL: tokenTTL,
		apiKey:   apiKey,
	}
}

type CustomClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int, log *logrus.Logger, tokenTTL int) string {
	claims := CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenTTL) * time.Hour)),
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Second)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNED_KEY")))
	if err != nil {
		log.Errorf("error occurred token service: %s", err.Error())
	}
	return tokenString
}

func (s *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.log.Error("error occurred parse token")
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("SIGNED_KEY")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		s.log.Error("error occurred parse token")
		return 0, errors.New("token claims are not of type *CustomClaims")
	}
	return claims.UserId, nil
}

func (s *UserService) RegisterUser(ctx context.Context, userDTO *dto.RegisterUser, session string, userId int) (string, error) {
	newUserId, err := s.repo.RegisterUser(ctx, userDTO, session, userId)
	if err != nil {
		return "", err
	} else {
		tokenTTL, err := strconv.Atoi(s.tokenTTL)
		if err != nil {
			return "", err
		}
		token := GenerateToken(newUserId, s.log, tokenTTL)
		return token, nil
	}
}

func (s *UserService) AuthenticateUser(ctx context.Context, userDTO *dto.UserAuth, session string, userId int) (string, error) {
	userId, err := s.repo.AuthenticateUser(ctx, userDTO, session, userId)
	if err != nil {
		return "", err
	} else {
		tokenTTL, err := strconv.Atoi(s.tokenTTL)
		if err != nil {
			return "", err
		}
		token := GenerateToken(userId, s.log, tokenTTL)
		return token, nil
	}
}

func (s *UserService) ValidateSession(sessionKey string) (string, int, error) {
	s.repo.SessionValidator()
	if sessionKey == "" {
		newSessionKey := pkg.GenerateSession()
		userId, err := s.repo.CreateSession(newSessionKey)
		return newSessionKey, userId, err
	} else {
		userId, sessionKey, err := s.repo.CheckSessionInDb(sessionKey)
		return sessionKey, userId, err
	}
}

func (s *UserService) GetUserData(ctx context.Context, userId int) (*dto.UserData, error) {
	return s.repo.GetUserData(ctx, userId)
}

// TODO dublicate code
func (s *UserService) RegCodeGenerator(ctx context.Context, code *dto.CodeGenerate) error {
	NewCode := CodeGenerator()
	err := s.repo.RegCodeGenerator(ctx, code, NewCode)
	if err != nil {
		return err
	}
	phoneCode := phone.NewPhone(s.apiKey)
	if err := phoneCode.SendCode(fmt.Sprint(NewCode), code.Phone); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *UserService) AuthCodeGenerator(ctx context.Context, code *dto.CodeGenerate) error {
	NewCode := CodeGenerator()
	err := s.repo.AuthCodeGenerator(ctx, code, NewCode)
	if err != nil {
		return err
	}
	phoneCode := phone.NewPhone(s.apiKey)
	if err = phoneCode.SendCode(fmt.Sprint(NewCode), code.Phone); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *UserService) UpdateEmailUser(ctx context.Context, emailInfo *dto.UpdateEmail, userId int) error {
	return s.repo.UpdateEmailUser(ctx, emailInfo, userId)
}

func (s *UserService) CanToSendEmail(ctx context.Context, emailInfo *dto.CanToSendEmail, userId int) error {
	return s.repo.CanToSendEmail(ctx, emailInfo, userId)
}

func (s *UserService) UpdateManagerName(ctx context.Context, userInfo *dto.UpdateManagerName, id int) error {
	return s.repo.UpdateManagerName(ctx, userInfo, id)
}

func (s *UserService) CreateSiteReview(ctx context.Context, siteReview *dto.SiteReview) error {
	return s.repo.CreateSiteReview(ctx, siteReview)
}
