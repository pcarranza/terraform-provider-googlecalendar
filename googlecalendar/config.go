package googlecalendar

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

var oauthScopes = []string{
	calendar.CalendarScope,
}

// Config is the structure used to instantiate the Google Calendar provider.
type Config struct {
	calendar *calendar.Service
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *Config) loadAndValidate() error {
	config, err := configFromFile()
	if err != nil {
		log.Fatalf("Unable to load configuration from file: %v", err)
	}

	tok, err := tokenFromFile()
	if err != nil {
		log.Fatalf("Unable to load oauth2 token from file: %v", err)
	}

	client := config.Client(context.Background(), tok)

	// Use a custom user-agent string. This helps google with analytics and it's
	// just a nice thing to do.
	client.Transport = logging.NewTransport("Google", client.Transport)
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, terraform.VersionString())

	// Create the calendar service.
	calendarSvc, err := calendar.New(client)
	if err != nil {
		return fmt.Errorf("failed creating a new calendar service: %v", err)
	}
	calendarSvc.UserAgent = userAgent
	c.calendar = calendarSvc

	return nil
}

func configFromFile() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(getEnvWithDefault("CALENDAR_CREDENTIALS_FILE", "credentials.json"))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read client secret file")
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse client secret file to config")
	}

	return config, nil
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(getEnvWithDefault("CALENDAR_OAUTH2_TOKEN_FILE", "token.json"))
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.Unmarshal(b, tok)
	return tok, err
}

// Reads an environment variable with a default in case it's not defined
func getEnvWithDefault(key, defaltValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaltValue
	}
	return value
}
