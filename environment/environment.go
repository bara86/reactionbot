package environment

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	clientID              = "CLIENT_ID"
	appURL                = "APP_URL"
	clientSecret          = "CLIENT_SECRET"
	slackTokenEnv         = "SLACK_TOKEN"
	connectionPort        = "PORT"
	slackOauthBotToken    = "SLACK_OAUTH_BOT_TOKEN"
	slackOauthAccessToken = "SLACK_OAUTH_ACCESS_TOKEN"
	saveFileName          = "SAVE_FILE_NAME"
	postgresDBURL         = "DATABASE_URL"
	botID                 = "BOT_ID"
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
		botID,
		slackOauthAccessToken,
	}
	for _, envVariable := range checkedEnvVariables {
		if _, ok := os.LookupEnv(envVariable); !ok {
			missingVariables = append(missingVariables, envVariable)
		}
	}

	if _, ok := os.LookupEnv(postgresDBURL); !ok {
		missingVariables = append(missingVariables, postgresDBURL)
	}

	return missingVariables
}

func GetOauthAccessToken() string {
	return getEnvVariable(slackOauthAccessToken)
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

func GetPostgresDBURL() string {
	return getEnvVariable(postgresDBURL)
}

func GetBotID() string {
	return getEnvVariable(botID)
}

func getEnvVariable(name string) string {
	return os.Getenv(name)
}
