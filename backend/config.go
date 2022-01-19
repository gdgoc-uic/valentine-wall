package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dghubble/oauth1"
	twitterOauth1 "github.com/dghubble/oauth1/twitter"
)

// uninit'ed variables
var databasePrefix = "valentine-wall"
var databasePath string
var baseUrl string
var twitterOauth1Config *oauth1.Config

// init'ed variables
var serverPort = 4000
var sessionName = "vw-session"
var frontendUrl = "http://localhost:3000"
var targetEnv = "development"

func init() {
	// if err := godotenv.Load("./.server.env"); err != nil {
	// 	log.Fatalln(err)
	// }

	twitterOauth1Config = &oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_CLIENT_ID"),
		ConsumerSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
		CallbackURL:    os.Getenv("TWITTER_CALLBACK_URL"),
		Endpoint:       twitterOauth1.AuthenticateEndpoint,
	}

	if gotPort, exists := os.LookupEnv("PORT"); exists {
		var err error
		serverPort, err = strconv.Atoi(gotPort)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if gotSessionName, exists := os.LookupEnv("SESSION_COOKIE_NAME"); exists {
		sessionName = gotSessionName
	}

	if _, exists := os.LookupEnv("FIREBASE_CONFIG"); !exists {
		log.Fatalln("path to firebase admin file is required")
	}

	if gotFrontendUrl, exists := os.LookupEnv("FRONTEND_URL"); exists {
		frontendUrl = gotFrontendUrl
	}

	if gotBaseUrl, exists := os.LookupEnv("BASE_URL"); exists {
		baseUrl = gotBaseUrl
	} else {
		baseUrl = fmt.Sprintf("http://localhost:%d", serverPort)
	}

	if gotDatabasePrefix, exists := os.LookupEnv("DATABASE_PREFIX"); exists {
		databasePath = gotDatabasePrefix
	}

	if gotTargetEnv, exists := os.LookupEnv("ENV"); exists {
		switch gotTargetEnv {
		case "development", "production", "staging":
			targetEnv = gotTargetEnv
		default:
			log.Fatalf("invalid environment '%s'\n", gotTargetEnv)
		}
	}

	databasePath = configureDatabasePath(targetEnv)
}
