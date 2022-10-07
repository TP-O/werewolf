package instance

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FBAuth *auth.Client

func initFirebaseInstance() {
	opt := option.WithCredentialsFile("firebase-service-account.json")

	if fb, err1 := firebase.NewApp(context.Background(), nil, opt); err1 != nil {
		log.Fatal(err1)
	} else if authClient, err2 := fb.Auth(context.Background()); err2 != nil {
		log.Fatal(err2)
	} else {
		FBAuth = authClient
	}
}
