package main

import (
	"bytes"
	"encoding/json"
	"flag"
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

type channel struct {
	ID string `json:"id"`
}

type addReactionActionMessage struct {
	Timestamp string `json:"ts"`
}

type addReactionAction struct {
	ResponseURL string                   `json:"response_url"`
	Channel     channel                  `json:"channel"`
	Message     addReactionActionMessage `json:"message"`
}

var asUser *bool

const slackTokenEnv = "SLACK_TOKEN"
const slackOauthBotToken = "SLACK_OAUTH_BOT_TOKEN"
const slackOauthUserToken = "SLACK_OAUTH_USER_TOKEN"
const slackAddReactionURL = "https://slack.com/api/reactions.add"

func getOauthToken() string {
	if !*asUser {
		return os.Getenv(slackOauthBotToken)
	}
	return os.Getenv(slackOauthUserToken)
}

func handleURLVerification(data *bytes.Buffer, w http.ResponseWriter) {
	type urlVerification struct {
		Token     string
		Challenge string
	}

	var urlverification urlVerification
	json.Unmarshal(data.Bytes(), &urlverification)

	slackToken := os.Getenv(slackTokenEnv)

	if slackToken != urlverification.Token {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		return
	}

	w.Write([]byte(urlverification.Challenge))
}

func addReaction(reactionName string, timestamp string, channel string) {

	oauthToken := getOauthToken()

	resp := response{Token: oauthToken, Name: reactionName, Timestamp: timestamp, Channel: channel}
	client := http.Client{}

	marshalled, _ := json.Marshal(resp)
	writer := bytes.NewBuffer(marshalled)

	request, erro := http.NewRequest(http.MethodPost, slackAddReactionURL, writer)

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
		go addReaction("thumbsup", msg.Event.Ts, msg.Event.Channel)
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

}

func addReactionToMessage(payload *string) {
	var addReactionActionMessage addReactionAction
	json.Unmarshal([]byte(*payload), &addReactionActionMessage)

	addReaction("heart", addReactionActionMessage.Message.Timestamp, addReactionActionMessage.Channel.ID)
}

func handleActions(w http.ResponseWriter, req *http.Request) {
	type callbackID struct {
		CallbackID string `json:"callback_id"`
	}

	var callbackid callbackID
	payload := req.FormValue("payload")

	json.Unmarshal([]byte(payload), &callbackid)

	if callbackid.CallbackID == "add_reaction_to_message" {
		go addReactionToMessage(&payload)
	} else {
		fmt.Printf("Error: Unhandled action `%s`", callbackid.CallbackID)
	}

	w.Write([]byte("GnocchettiAlVapore"))
}

func main() {
	asUser = flag.Bool("u", false, "React as user")
	flag.Parse()

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Missing .env file")
	}

	fmt.Println("Ready to react!!1!")

	http.HandleFunc("/", handle)
	http.HandleFunc("/actions", handleActions)
	http.ListenAndServe(":8008", nil)

}
