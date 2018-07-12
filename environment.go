package main

import "os"

const (
	clientID = "CLIENT_ID"

	appURL = "APP_URL"

	clientSecret = "CLIENT_SECRET"
)

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
