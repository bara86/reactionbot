package main

import (
	"fmt"

	"reactionbot/commonstructure"
	"reactionbot/environment"
	"reactionbot/handlers"
	"reactionbot/storageondb"
)

func main() {

	if err := environment.LoadEnvironmentVariables(); err != nil {
		panic(err)
	}

	var storage commonstructure.Storage

	storage, err := storageondb.SetUp()
	storage.LoadEmojisList()

	if err != nil {
		panic(err)
	}

	fmt.Println("Ready to react!!1!")
	if err := handlers.StartServer(storage); err != nil {
		fmt.Println(err)
	}
}
