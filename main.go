package main

import (
	"fmt"

	"reactionbot/environment"
	"reactionbot/handlers"
	"reactionbot/storage"
)

func main() {

	if err := environment.LoadEnvironmentVariables(); err != nil {
		panic(err)
	}

	fmt.Println("Ready to react!!1!")

	storage := storage.UserStorage{}
	if err := handlers.StartServer(&storage); err != nil {
		fmt.Println(err)
	}
}
