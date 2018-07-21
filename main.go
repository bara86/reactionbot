package main

import (
	"fmt"

	"reactionbot/environment"
	"reactionbot/handlers"
	"reactionbot/storageondb"
	"reactionbot/storageonfile"
)

func main() {

	if err := environment.LoadEnvironmentVariables(); err != nil {
		panic(err)
	}

	a, err := storageondb.SetUp()

	if err == nil {
		fmt.Println("connesso", a)
	} else {
		fmt.Println("non connesso", err)
	}
	err = a.Add("c", "b")
	if err != nil {
		fmt.Println(err)
	}
	// panic("dd")
	val, err := a.Lookup("c")
	fmt.Println("eeeee", val, err)

	token, err := a.Get("c")
	fmt.Println("get", val, token)

	err = a.Remove("c")
	fmt.Println("remove c", err)

	err = a.Remove("k")
	fmt.Println("remove k", err)

	storage, err := storageonfile.SetUp()
	if err != nil {
		panic(err)
	}

	fmt.Println("Ready to react!!1!")
	if err := handlers.StartServer(storage); err != nil {
		fmt.Println(err)
	}
}
