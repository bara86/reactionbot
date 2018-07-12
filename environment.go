package main

import "os"

const clientID = "CLIENT_ID"

func getClientID() string {
	return getEnvVariable(clientID)
}

func getEnvVariable(name string) string {
	return os.Getenv(name)
}
