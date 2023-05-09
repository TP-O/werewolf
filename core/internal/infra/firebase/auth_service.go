package firebase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"uwwolf/internal/app/game/logic/types"
	"uwwolf/internal/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type AuthService interface {
	VerifyAuthorization(authorization string) (types.PlayerID, error)
}

type authService struct {
	auth *auth.Client
}

func NewAuthService(config config.Firebase) (AuthService, error) {
	opt := option.WithCredentialsJSON([]byte(fmt.Sprintf(`{
        "projectId": %v,
        "privateKey": %v,
        "clientEmail: %v
    }`,
		config.ProjectId,
		strings.ReplaceAll(config.PrivateKey, "\\n", "\n"),
		config.Email)),
	)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &authService{auth}, nil
}

func (as authService) VerifyAuthorization(authorization string) (types.PlayerID, error) {
	if len(authorization) == 0 {
		return "", errors.New("Empty authorization header!")
	}

	idToken := strings.TrimPrefix(authorization, "Bearer ")
	token, err := as.auth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", err
	}

	return types.PlayerID(token.UID), nil
}
