package main

import (
	"fmt"
	"reactionbot/storageonfile"

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

	saveOnFile, err := environment.GetSaveOnFile()
	if err != nil {
		panic(fmt.Sprintf("Wrong format for \"SaveOnFile\" env variable %v", err))
	}

	if saveOnFile {
		storage, err = storageonfile.SetUp()
	} else {
		storage, err = storageondb.SetUp()
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("Ready to react!!1!")
	if err := handlers.StartServer(storage); err != nil {
		fmt.Println(err)
	}
}
