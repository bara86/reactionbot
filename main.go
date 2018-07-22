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

	fmt.Println("sto per fare add")
	err = storage.AddUser("pippo", "pluto")
	if err != nil {
		panic(err)
	}

	found, err := storage.LookupUser("casa")
	fmt.Println("looking for casa", found, err)

	found, err = storage.LookupUser("pippo")
	fmt.Println("looking for pippo", found, err)

	err = storage.RemoveUser("casa")
	fmt.Println(err)

	fmt.Println(storage.RemoveUser("pippo"))
	return

	fmt.Println("Ready to react!!1!")
	if err := handlers.StartServer(storage); err != nil {
		fmt.Println(err)
	}
}
