package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

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

type clientResponseData struct {
	Ok  bool
	Err string `json:"error"`
}

type channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const slackAddReactionURL = "https://slack.com/api/reactions.add"
const slackOauthAccessURL = "https://slack.com/api/oauth.access"

func unmarshallData(reader io.Reader, v interface{}) {
	var data bytes.Buffer
	data.ReadFrom(reader)
	err := json.Unmarshal(data.Bytes(), &v)

	if err != nil {
		fmt.Println("Unable to unmarshal data for type", reflect.TypeOf(v).String())
	}
}

func handleURLVerification(w http.ResponseWriter, req *http.Request) {
	type urlVerification struct {
		Token     string
		Challenge string
	}

	var urlverification urlVerification
	unmarshallData(req.Body, &urlverification)

	if getSlackToken() != urlverification.Token {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		return
	}

	w.Write([]byte(urlverification.Challenge))
}

func postToSlack(token string, url string, w io.Reader) (*http.Response, error) {
	request, erro := http.NewRequest(http.MethodPost, url, w)

	if erro != nil {
		fmt.Println("Error creating request")
		return &http.Response{}, erro
	}

	// Add Authorization token
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Content-type", "application/json")

	client := http.Client{}
	clientResponse, clientError := client.Do(request)

	if clientError != nil {
		fmt.Println("Errore dal client")
	} else {
		var clientRespData clientResponseData

		unmarshallData(clientResponse.Body, &clientRespData)

		if !clientRespData.Ok {
			fmt.Println("Error send HTTP post request to Slack:", clientRespData.Err)
		}
	}
	return clientResponse, clientError
}

func addReaction(token string, reactionName string, timestamp string, channel string) {
	fmt.Println("addReaction method")
	resp := response{Token: token, Name: reactionName, Timestamp: timestamp, Channel: channel}
	marshalled, _ := json.Marshal(resp)
	stringBuffer := bytes.NewBuffer(marshalled)
	postToSlack(token, slackAddReactionURL, stringBuffer)
}

func handleEvent(data io.Reader) {
	var msg message

	unmarshallData(data, &msg)

	if msg.Event.Type == "message" {
		go addReaction(getOauthToken(), "thumbsup", msg.Event.Ts, msg.Event.Channel)
	}
}

func handle(w http.ResponseWriter, req *http.Request) {

	type messageType struct {
		Type string
	}

	var messagetype messageType
	unmarshallData(req.Body, &messagetype)

	if messagetype.Type == "url_verification" {
		handleURLVerification(w, req)
	} else {
		if messagetype.Type != "event_callback" {
			panic(fmt.Sprint("Not `event_callback`", messagetype.Type))
		}
		handleEvent(req.Body)
	}

	w.Write([]byte("GnocchettiAlVapore"))

}

func handleActions(w http.ResponseWriter, req *http.Request) {
	type callbackID struct {
		CallbackID string `json:"callback_id"`
	}

	var callbackid callbackID
	payload := req.FormValue("payload")
	unmarshallData(strings.NewReader(payload), &callbackid)

	if callbackid.CallbackID == "add_reaction_to_message" {
		go addReactionToMessage(&payload)
	} else {
		fmt.Printf("Error: Unhandled action `%s`", callbackid.CallbackID)
	}

	w.Write([]byte("GnocchettiAlVapore"))
}

func main() {

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Missing .env file, try to read env variables anyway")
	}

	if missingVariables := checkEnvVariables(); len(missingVariables) != 0 {
		panic(fmt.Sprintf("Missing env variables %v, can't continue", missingVariables))
	}

	fmt.Println("Ready to react!!1!")

	http.HandleFunc("/", handle)
	http.HandleFunc("/actions", handleActions)
	http.HandleFunc("/oauth", handleOauth)
	http.ListenAndServe(fmt.Sprintf(":%s", getConnectionPort()), nil)

}
