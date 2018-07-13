package main

import (
	"os"
)

const (
	clientID            = "CLIENT_ID"
	appURL              = "APP_URL"
	clientSecret        = "CLIENT_SECRET"
	slackTokenEnv       = "SLACK_TOKEN"
	connectionPort      = "PORT"
	slackOauthBotToken  = "SLACK_OAUTH_BOT_TOKEN"
	slackOauthUserToken = "SLACK_OAUTH_USER_TOKEN"
)

func checkEnvVariables() []string {

	var missingVariables []string
	checkedEnvVariables := []string{clientID, appURL, clientSecret, slackTokenEnv, connectionPort, slackOauthBotToken, slackOauthUserToken}
	for _, envVariable := range checkedEnvVariables {
		if _, ok := os.LookupEnv(envVariable); !ok {
			missingVariables = append(missingVariables, envVariable)
		}
	}
	return missingVariables
}

func getOauthToken(user bool) string {
	if user {
		return getEnvVariable(slackOauthUserToken)
	}
	return getEnvVariable(slackOauthBotToken)
}

func getConnectionPort() string {
	return getEnvVariable(connectionPort)
}

func getSlackToken() string {
	return getEnvVariable(slackTokenEnv)
}

func getClientID() string {
	return getEnvVariable(clientID)
}

func getClientSecret() string {
	return getEnvVariable(clientSecret)
}

func getAppURL() string {
	return getEnvVariable(appURL)
}

func getEnvVariable(name string) string {
	return os.Getenv(name)
}
