package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reactionbot/commonstructure"
	"reflect"
	"strings"

	"reactionbot/environment"

	"github.com/satori/go.uuid"
)

const (
	slackEphemeralURL   = "https://slack.com/api/chat.postEphemeral"
	slackAddReactionURL = "https://slack.com/api/reactions.add"
	slackOauthAccessURL = "https://slack.com/api/oauth.access"
	slackOauthURL       = "slack.com/oauth/authorize"

	authorizeButton = `{
		"token": "%v",
		"channel": "%v",
		"text": "We need your approval to post reaction with your name.",
		"user": "%v",
		"attachments": [
			{
				"fallback": "Approve posting a reaction with your name https://foo.bar",
				"callback_id": "mignatta",
				"actions": [
					{
						"type": "button",
						"text": "Approve",
						"name": "approve",
						"url": "%v"
					}
				]
			}
		]
	}`

	reactionsWriteScope = "reactions:write"
)

var userTokens commonstructure.Storage
var temporaryTokens commonstructure.Storage

func StartServer(storage commonstructure.Storage) error {
	userTokens = storage
	temporaryTokens = storage

	http.HandleFunc("/", handle)
	http.HandleFunc("/actions", handleActions)
	http.HandleFunc("/oauth", handleOauth)
	return http.ListenAndServe(fmt.Sprintf(":%s", environment.GetConnectionPort()), nil)
}

func handle(w http.ResponseWriter, req *http.Request) {

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

func unmarshallData(reader io.Reader, v interface{}) {
	var data bytes.Buffer
	data.ReadFrom(reader)
	err := json.Unmarshal(data.Bytes(), &v)

	if err != nil {
		fmt.Println("Unable to unmarshal data for type", reflect.TypeOf(v).String())
	}
}

func handleURLVerification(w http.ResponseWriter, req *http.Request) {

	var urlverification urlVerification
	unmarshallData(req.Body, &urlverification)

	if environment.GetSlackToken() != urlverification.Token {
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
		go addReaction(environment.GetOauthToken(), "thumbsup", msg.Event.Ts, msg.Event.Channel)
	}
}

func handleActions(w http.ResponseWriter, req *http.Request) {

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

func addReactionToMessage(payload *string) {

	var info addReactionAction
	json.Unmarshal([]byte(*payload), &info)
	token, err := userTokens.Get(info.User.ID)

	if err != nil {
		if postEphemeralMessage(&info) != nil {

		}
	} else {
		addReaction(token, "heart", info.Message.Timestamp, info.Channel.ID)
	}

}

func postEphemeralMessage(info *addReactionAction) error {
	fmt.Println("Post ephemeral message")

	url := url.URL{Path: slackOauthURL, Scheme: "https"}
	uuid := uuid.Must(uuid.NewV4()).String()

	q := url.Query()
	q.Add("client_id", environment.GetClientID())
	q.Add("scope", reactionsWriteScope)
	q.Add("redirect_uri", fmt.Sprintf("%v/oauth", environment.GetAppURL()))
	q.Add("state", uuid)
	url.RawQuery = q.Encode()

	temporaryTokens.Add(uuid, info.User.ID)

	jsonMsg := fmt.Sprintf(authorizeButton, environment.GetOauthToken(), info.Channel.ID, info.User.ID, url.String())
	fmt.Println(jsonMsg)
	buf := bytes.NewBufferString(jsonMsg)
	_, err2 := postToSlack(environment.GetOauthToken(), slackEphemeralURL, buf)
	if err2 != nil {
		return err2
	}
	return nil
}

func handleOauth(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handleOauth", req.URL.Query())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<body onload=\"window.open(window.location.href, '_self', ''); window.close()\"></body>")

	query := req.URL.Query()
	state := string(query["state"][0])

	resp, _ := http.PostForm(slackOauthAccessURL,
		url.Values{
			"client_id":     {environment.GetClientID()},
			"client_secret": {environment.GetClientSecret()},
			"code":          {string(query["code"][0])},
			"redirect_url":  {fmt.Sprintf("%v/oauth", environment.GetAppURL())},
		})
	fmt.Println(resp.Body)

	userID, err := temporaryTokens.Pop(state)
	if err != nil {
		fmt.Println(err)
	}

	var accessTokenData accessToken
	unmarshallData(resp.Body, &accessTokenData)

	userTokens.Add(userID, accessTokenData.AccessToken)
}
