package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	goaway "github.com/TwiN/go-away"
	"github.com/dghubble/oauth1"
	twitterOauth1 "github.com/dghubble/oauth1/twitter"
	goValidator "github.com/go-playground/validator/v10"
)

type CustomProfanityDictionary struct {
	Profanities    []string
	FalsePositives []string
	FalseNegatives []string
}

// uninit'ed variables
var databasePrefix = "valentine-wall"
var databasePath string
var baseUrl string
var twitterOauth1Config *oauth1.Config
var gAppCredPath string

// TODO: add custom dictionary for bisaya and tagalog
var profanityDetector *goaway.ProfanityDetector

// init'ed variables
var serverPort = 4000
var sessionName = "vw-session"
var frontendUrl = "http://localhost:3000"
var targetEnv = "development"
var postalOfficeAddress = "localhost:3350"

var validator = goValidator.New()
var giftList = []Gift{
	{1, "sunflower", "Sunflower"},
	{2, "rose", "Rose"},
	{3, "balloons", "Balloons"},
	{4, "teddy-bear", "Teddy Bear"},
	{5, "ring", "Ring"},
	{6, "money", "Money"},
	{7, "heart", "Heart"},
	{8, "chocolate", "Chocolate"},
	{9, "pizza", "Pizza"},
}
var collegeDepartments = []CollegeDepartment{
	{"CCS", "College of Computer Studies"},
}

func loadCustomProfanityDetector(customDictionary *CustomProfanityDictionary) *goaway.ProfanityDetector {
	return goaway.NewProfanityDetector().WithCustomDictionary(
		append(goaway.DefaultProfanities, customDictionary.Profanities...),
		append(goaway.DefaultFalsePositives, customDictionary.FalsePositives...),
		append(goaway.DefaultFalseNegatives, customDictionary.FalseNegatives...),
	)
}

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

	if gotGAppCredPath, exists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); exists {
		gAppCredPath = gotGAppCredPath
	}

	if len(gAppCredPath) == 0 {
		log.Fatalln("path to firebase/google service account file is required")
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

	if gotProfanityListFilePath, exists := os.LookupEnv("PROFANITY_JSON_FILE_PATH"); exists {
		var data []byte
		var err error

		if data, err = ioutil.ReadFile(gotProfanityListFilePath); err != nil {
			log.Fatalln(err)
		}

		customDictionary := &CustomProfanityDictionary{}
		if err := json.Unmarshal(data, customDictionary); err != nil {
			log.Fatalln(err)
		}

		profanityDetector = loadCustomProfanityDetector(customDictionary)
	} else {
		profanityDetector = goaway.NewProfanityDetector()
	}

	if gotPoAddress, exists := os.LookupEnv("POSTAL_OFFICE_ADDRESS"); exists {
		postalOfficeAddress = gotPoAddress
	}
}
