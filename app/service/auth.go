package service

import (
	"context"
	"errors"
	"strings"

	"uwwolf/app/instance"
)

func Verify(authorization string) (string, error) {
	var token string

	if result := strings.Split(authorization, "Bearer "); len(result) != 2 {
		return "", errors.New("Missing access token!")
	} else {
		token = result[1]
	}

	if token, err := instance.FBAuth.VerifyIDToken(context.Background(), token); err != nil {
		return "", errors.New("Invalid access token!")
	} else {
		return token.UID, nil
	}
}
