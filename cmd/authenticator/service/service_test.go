package service

import (
	"testing"

	"github.com/ViniciusMiana/sensor-metadata/cmd/authenticator/db"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLogin(t *testing.T) {
	service, err := NewAuthenticatorService("mongodb://localhost:27017", "users"+primitive.NewObjectID().Hex())
	require.NoError(t, err)
	result, err := service.Login(db.User{
		Username: "root",
		Password: "1234",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
}
