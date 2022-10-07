package service

import (
	"context"
	"errors"
	"uwwolf/app/instance"
)

func Verify(authorization string) (string, error) {
	if token, err := instance.FBAuth.VerifyIDToken(context.Background(), authorization); err != nil {
		return "", errors.New("Invalid token!")
	} else {
		return token.UID, nil
	}
}
