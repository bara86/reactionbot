package main

import (
	"fmt"

	"reactionbot/environment"
	"reactionbot/handlers"
	"reactionbot/storageondb"
)

func main() {

	if err := environment.LoadEnvironmentVariables(); err != nil {
		panic(err)
	}

	storage, err := storageondb.SetUp()

	if err != nil {
		panic(err)
	}

	fmt.Println("Ready to react!!1!")
	if err := handlers.StartServer(storage); err != nil {
		fmt.Println(err)
	}
}
