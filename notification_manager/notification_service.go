package notificationmanager

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var notificationRepo = NewNotificationManagerRepo()

type NotificationManagerService interface{}

type NotificationManager struct {
	FCMApp *firebase.App
}

func (nm *NotificationManager) InitApp() {
	fmt.Println("InitApp get called...")
	opt := option.WithCredentialsFile("./maharaj-fcm.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return
	}
	nm.FCMApp = app
}

func (nm *NotificationManager) NotifyAllUser(message *messaging.Message) {
	//check fcm app is availabel or not
	if nm.FCMApp == nil {
		nm.InitApp()
		log.Println("App has been init")
	}

	// fetch all tokens
	tokens, err := notificationRepo.GetTokens()

	if err != nil {
		return
	}

	fmt.Println("Tokens to send notification : ", len(tokens))

	// Get the FCM client.
	client, err := nm.FCMApp.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting FCM client: %v\n", err)
		return
	}

	// notify all users
	for _, tokenData := range tokens {
		message.Token = tokenData.Token

		response, err := client.Send(context.Background(), message)

		if err != nil {
			fmt.Println("Notification Error : ", err.Error())
			continue
		}
		log.Println("Notification Response : ", response)
	}

	fmt.Println("Notification has been send!")
}
