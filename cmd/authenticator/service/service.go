package service

import (
	"context"
	"errors"
	"os"

	"github.com/ViniciusMiana/sensor-metadata/cmd/authenticator/db"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// TODO move to a common pkg folder
// TokenClaims contains a basic jwt claim structure
type TokenClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Valid overrides StandardClaims validation.
func (c TokenClaims) Valid() error {
	return nil
}

type AuthenticatorService interface {
	// Creates a new user
	Register(user db.User) error
	// Returns a JWT token for a user if exists
	Login(user db.User) (string, error)
}

type authenticatorService struct {
	userStore *db.UserStore
	key       []byte
}

func NewAuthenticatorService(uri, databaseName string) (*authenticatorService, error) {
	ss, err := db.NewUserStore(uri, databaseName)
	if err != nil {
		return nil, err
	}
	return &authenticatorService{
		userStore: ss,
		key:       []byte(os.Getenv("tls.key")),
	}, nil
}

func (as authenticatorService) Register(user db.User) error {
	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPasswordBytes)
	return as.userStore.AddUser(context.Background(), user)
}

// Returns a JWT token for a user if exists
func (as authenticatorService) Login(user db.User) (string, error) {
	dbUser, err := as.userStore.FindByUserName(context.Background(), user.Username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("Invalid Password")
	}
	// Generate the claim...
	token, err := as.GenerateToken(TokenClaims{
		Username: dbUser.Username,
		Role:     dbUser.Role,
	})
	return token, err
}

// GenerateRefreshToken generates a new jwt token
func (as authenticatorService) GenerateToken(claim TokenClaims) (string, error) {
	alg := jwt.SigningMethodRS256
	key, err := jwt.ParseRSAPrivateKeyFromPEM(as.key)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(alg, claim)
	return token.SignedString(key)
}
