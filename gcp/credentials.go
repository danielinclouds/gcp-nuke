package gcp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	Email string `json:"client_email"`
	JSON  []byte
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func FindCredentials(credentialsPath string) (Credentials, error) {

	if credentialsPath != "" {
		b, err := ioutil.ReadFile(credentialsPath)
		if err != nil {
			return Credentials{}, err
		}

		var creds Credentials
		creds.JSON = b
		err = json.Unmarshal(b, &creds)
		if err != nil {
			return Credentials{}, err
		}

		log.Debugf("Found credentials in file: %s", credentialsPath)
		return creds, nil

	} else if os.Getenv("GOOGLE_CREDENTIALS") != "" {
		b := []byte(os.Getenv("GOOGLE_CREDENTIALS"))

		var creds Credentials
		creds.JSON = b
		err := json.Unmarshal(b, &creds)
		if err != nil {
			return Credentials{}, err
		}

		log.Debug("Found credentials in env: GOOGLE_CREDENTIALS")
		return creds, nil

	} else if os.Getenv("GOOGLE_CLOUD_KEYFILE_JSON") != "" {
		b := []byte(os.Getenv("GOOGLE_CLOUD_KEYFILE_JSON"))

		var creds Credentials
		creds.JSON = b
		err := json.Unmarshal(b, &creds)
		if err != nil {
			return Credentials{}, err
		}

		log.Debug("Found credentials in env: GOOGLE_CLOUD_KEYFILE_JSON")
		return creds, nil
	}

	return Credentials{}, errors.New("Could not find credentials")
}
