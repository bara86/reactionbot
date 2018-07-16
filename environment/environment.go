package environment

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	clientID           = "CLIENT_ID"
	appURL             = "APP_URL"
	clientSecret       = "CLIENT_SECRET"
	slackTokenEnv      = "SLACK_TOKEN"
	connectionPort     = "PORT"
	slackOauthBotToken = "SLACK_OAUTH_BOT_TOKEN"
	saveOnFile         = "SAVE_ON_FILE"
	saveFileName       = "SAVE_FILE_NAME"
)

func LoadEnvironmentVariables() error {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Missing .env file, try to read env variables anyway")
	}

	if missingVariables := checkEnvVariables(); len(missingVariables) != 0 {
		return fmt.Errorf("Missing env variables %v, can't continue", missingVariables)
	}
	return nil
}

func checkEnvVariables() []string {

	var missingVariables []string
	checkedEnvVariables := []string{
		clientID,
		appURL,
		clientSecret,
		slackTokenEnv,
		connectionPort,
		slackOauthBotToken,
		saveOnFile,
	}
	for _, envVariable := range checkedEnvVariables {
		if _, ok := os.LookupEnv(envVariable); !ok {
			missingVariables = append(missingVariables, envVariable)
		}
	}

	if v, _ := GetSaveOnFile(); v {
		if _, ok := os.LookupEnv(saveFileName); !ok {
			missingVariables = append(missingVariables, saveFileName)
		}
	}

	return missingVariables
}

func GetSaveOnFile() (bool, error) {
	return strconv.ParseBool(getEnvVariable(saveOnFile))
}

func GetSaveFileName() string {
	return getEnvVariable(saveFileName)
}

func GetOauthToken() string {
	return getEnvVariable(slackOauthBotToken)
}

func GetConnectionPort() string {
	return getEnvVariable(connectionPort)
}

func GetSlackToken() string {
	return getEnvVariable(slackTokenEnv)
}

func GetClientID() string {
	return getEnvVariable(clientID)
}

func GetClientSecret() string {
	return getEnvVariable(clientSecret)
}

func GetAppURL() string {
	return getEnvVariable(appURL)
}

func getEnvVariable(name string) string {
	return os.Getenv(name)
}
