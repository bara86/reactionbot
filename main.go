package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

type clientResponseData struct {
	Ok  bool
	Err string `json:"error"`
}

type channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var asUser *bool

const slackTokenEnv = "SLACK_TOKEN"
const slackOauthBotToken = "SLACK_OAUTH_BOT_TOKEN"
const slackOauthUserToken = "SLACK_OAUTH_USER_TOKEN"
const slackAddReactionURL = "https://slack.com/api/reactions.add"
const connectionPort = "PORT"

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

func postToSlack(url string, w io.Reader) (*http.Response, error) {
	request, erro := http.NewRequest(http.MethodPost, url, w)

	if erro != nil {
		fmt.Println("Error creating request")
		return &http.Response{}, erro
	}

	// Add Authorization token
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getOauthToken()))
	request.Header.Add("Content-type", "application/json")

	client := http.Client{}
	clientResponse, clientError := client.Do(request)

	if clientError != nil {
		fmt.Println("Errore dal client")
	} else {
		var data bytes.Buffer
		var clientRespData clientResponseData

		data.ReadFrom(clientResponse.Body)
		json.Unmarshal(data.Bytes(), &clientRespData)

		if !clientRespData.Ok {
			fmt.Println("Error send HTTP post request to Slack:", clientRespData.Err)
		}
	}
	return clientResponse, clientError
}

func addReaction(reactionName string, timestamp string, channel string) {
	fmt.Println("addReaction method")
	resp := response{Token: getOauthToken(), Name: reactionName, Timestamp: timestamp, Channel: channel}
	marshalled, _ := json.Marshal(resp)
	stringBuffer := bytes.NewBuffer(marshalled)
	postToSlack(slackAddReactionURL, stringBuffer)
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

func checkEnvVariables(envVariables []string) []string {

	var missingVariables []string
	for _, envVariable := range envVariables {
		if _, ok := os.LookupEnv(envVariable); !ok {
			missingVariables = append(missingVariables, envVariable)
		}
	}
	return missingVariables
}

func main() {
	asUser = flag.Bool("u", false, "React as user")
	flag.Parse()

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Missing .env file, try to read env variables anyway")
	}

	if missingVariables := checkEnvVariables([]string{slackTokenEnv, slackOauthBotToken, slackOauthUserToken, connectionPort}); len(missingVariables) != 0 {
		panic(fmt.Sprintf("Missing env variables %v, can't continue", missingVariables))
	}

	port := os.Getenv(connectionPort)

	fmt.Println("Ready to react!!1!")

	http.HandleFunc("/", handle)
	http.HandleFunc("/actions", handleActions)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
