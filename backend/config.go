package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	goaway "github.com/TwiN/go-away"
	goValidator "github.com/go-playground/validator/v10"
)

type CustomProfanityDictionary struct {
	Profanities    []string
	FalsePositives []string
	FalseNegatives []string
}

var emailRegex *regexp.Regexp

// uninit'ed variables
var baseUrl string
var gAppCredPath string
var chromeDevtoolsURL string
var readOnly = false

// TODO: add custom dictionary for bisaya and tagalog
var profanityDetector *goaway.ProfanityDetector

// init'ed variables
var serverPort = 4000
var targetEnv = "development"
var postalOfficeAddress = "localhost:3350"
var dataDirPath = filepath.Join(".", "_data")
var invitationCookieName = "vw-invitation_ref"

var sendPrice = float64(150.0)

var validator = goValidator.New()

func loadCustomProfanityDetector(customDictionary *CustomProfanityDictionary) *goaway.ProfanityDetector {
	return goaway.NewProfanityDetector().WithCustomDictionary(
		append(goaway.DefaultProfanities, customDictionary.Profanities...),
		append(goaway.DefaultFalsePositives, customDictionary.FalsePositives...),
		append(goaway.DefaultFalseNegatives, customDictionary.FalseNegatives...),
	)
}

func init() {
	log.Println("compiling email regex...")
	{
		var err error
		emailRegex, err = regexp.Compile(`\A[a-z]+_([0-9]+)@uic.edu.ph\z`)
		if err != nil {
			log.Panicln(err)
		}
	}

	// if err := godotenv.Load("./.server.env"); err != nil {
	// 	log.Panicln(err)
	// }

	if gotPort, exists := os.LookupEnv("PORT"); exists {
		var err error
		serverPort, err = strconv.Atoi(gotPort)
		if err != nil {
			log.Panicln(err)
		}
	}

	if gotBaseUrl, exists := os.LookupEnv("BASE_URL"); exists {
		baseUrl = gotBaseUrl
	} else {
		baseUrl = fmt.Sprintf("http://localhost:%d", serverPort)
	}

	if gotTargetEnv, exists := os.LookupEnv("ENV"); exists {
		switch gotTargetEnv {
		case "development", "production", "staging":
			targetEnv = gotTargetEnv
		default:
			log.Fatalf("invalid environment '%s'\n", gotTargetEnv)
		}
	}

	if gotProfanityListFilePath, exists := os.LookupEnv("PROFANITY_JSON_FILE_PATH"); exists {
		var data []byte
		var err error

		if data, err = os.ReadFile(gotProfanityListFilePath); err != nil {
			log.Panicln(err)
		}

		customDictionary := &CustomProfanityDictionary{}
		if err := json.Unmarshal(data, customDictionary); err != nil {
			log.Panicln(err)
		}

		log.Println("loading custom profanity detector...")
		profanityDetector = loadCustomProfanityDetector(customDictionary)
	} else {
		profanityDetector = goaway.NewProfanityDetector()
	}

	if gotPoAddress, exists := os.LookupEnv("POSTAL_OFFICE_ADDRESS"); exists {
		postalOfficeAddress = gotPoAddress
	}

	if gotChromeDevtoolsURL, exists := os.LookupEnv("CHROME_DEVTOOLS_URL"); exists {
		chromeDevtoolsURL = gotChromeDevtoolsURL
	}

	if gotInvitationCookieName, exists := os.LookupEnv("INVITATION_COOKIE_NAME"); exists {
		invitationCookieName = gotInvitationCookieName
	}

	if gotReadOnly, exists := os.LookupEnv("READ_ONLY"); exists {
		readOnly = gotReadOnly == "true"
	}
}
