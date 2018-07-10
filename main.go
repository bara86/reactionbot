package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type event struct {
	Type    string
	EventTs string `json:"event_ts"`
	User    string
	Item    string
	Ts      string
	Channel string
}

type message struct {
	Token     string
	Challenge string
	Type      string
	EventID   string `json:"event_id"`
	Event     event
}

type response struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
}

type responsebello struct {
	Ok  bool
	Err string `json:"error"`
}

func handleURLVerification(data *bytes.Buffer, w http.ResponseWriter) {
	type urlVerification struct {
		Token     string
		Challenge string
	}

	var urlverification urlVerification
	json.Unmarshal(data.Bytes(), &urlverification)

	slackToken := os.Getenv("SLACK_TOKEN")

	if slackToken != urlverification.Token {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		return
	}

	w.Write([]byte(urlverification.Challenge))
}

func addReaction(msg *message) {

	oauthToken := os.Getenv("SLACK_OAUTH_TOKEN")

	resp := response{Token: oauthToken, Name: "thumbsup", Timestamp: msg.Event.Ts, Channel: msg.Event.Channel}
	client := http.Client{}

	marshalled, _ := json.Marshal(resp)
	writer := bytes.NewBuffer(marshalled)

	request, erro := http.NewRequest("POST", "https://slack.com/api/reactions.add", writer)

	if erro != nil {
		fmt.Println("Error creating request")
		return
	}

	// Add Authorization token
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oauthToken))
	request.Header.Add("Content-type", "application/json")

	clientResponse, clientError := client.Do(request)

	if clientError != nil {
		fmt.Println("Errore dal client")
	} else {
		var data bytes.Buffer
		data.Reset()
		data.ReadFrom(clientResponse.Body)
		var respbello responsebello
		json.Unmarshal(data.Bytes(), &respbello)

		fmt.Println("ClientResponse", respbello)
	}

}

func handleEvent(data *bytes.Buffer) {
	var msg message

	json.Unmarshal(data.Bytes(), &msg)

	if msg.Event.Type == "message" {
		go addReaction(&msg)
	}
}

func handle(w http.ResponseWriter, req *http.Request) {

	type messageType struct {
		Type string
	}

	var data bytes.Buffer

	var messagetype messageType
	data.ReadFrom(req.Body)
	json.Unmarshal(data.Bytes(), &messagetype)

	if messagetype.Type == "url_verification" {
		handleURLVerification(&data, w)
	} else {
		if messagetype.Type != "event_callback" {
			panic(fmt.Sprint("Not `event_callback`", messagetype.Type))
		}
		handleEvent(&data)
	}

	w.Write([]byte("GnocchettiAlVapore"))

	// var ch message

	// data.ReadFrom(req.Body)
	// fmt.Println("Challenge", req.Body)
	// json.Unmarshal(data.Bytes(), &ch)

	// fmt.Println("type", ch.Type, ch.EventID, ch.Event.Type)

	// resp := response{Token: "xoxp-2362083804-80850956512-395538191665-32b66e7481ebf4010d601397969d5528", Name: "thumbsup", Timestamp: ch.Event.Ts, Channel: ch.Event.Channel}

	// client := http.Client{}

	// marshalled, _ := json.Marshal(resp)
	// fmt.Println(string(marshalled))
	// writer := bytes.NewBuffer(marshalled)

	// request, erro := http.NewRequest("POST", "https://slack.com/api/reactions.add", writer)

	// if erro != nil {
	// 	fmt.Println("Error creating request")
	// 	return
	// }

	// // Add Authorization token
	// request.Header.Add("Authorization", "Bearer xoxp-2362083804-80850956512-395538191665-32b66e7481ebf4010d601397969d5528")
	// request.Header.Add("Content-type", "application/json")

	// clientResponse, clientError := client.Do(request)

	// if clientError != nil {
	// 	fmt.Println("Errore dal client")
	// } else {
	// 	data.Reset()
	// 	data.ReadFrom(clientResponse.Body)
	// 	var respbello responsebello
	// 	json.Unmarshal(data.Bytes(), &respbello)

	// 	fmt.Println("ClientResponse", respbello)
	// }

}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Missing dotenv file")
	}

	http.HandleFunc("/", handle)
	http.ListenAndServe(":8008", nil)
}
