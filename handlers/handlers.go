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
	slackEphemeralURL       = "https://slack.com/api/chat.postEphemeral"
	slackAddReactionURL     = "https://slack.com/api/reactions.add"
	slackOauthAccessURL     = "https://slack.com/api/oauth.access"
	slackOauthURL           = "slack.com/oauth/authorize"
	slackChatPostMessageURL = "https://slack.com/api/chat.postMessage"
	slackGetEmojisListURL   = "https://slack.com/api/emoji.list"

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

var dataStorage commonstructure.Storage

func StartServer(storage commonstructure.Storage) error {
	dataStorage = storage

	if err := loadCustomEmojis(); err != nil {
		return err
	}

	http.HandleFunc("/actions", handleActions)
	http.HandleFunc("/oauth", handleOauth)
	http.HandleFunc("/events", handleEvents)
	return http.ListenAndServe(fmt.Sprintf(":%s", environment.GetConnectionPort()), nil)
}

func loadCustomEmojis() error {
	resp, err := http.PostForm(slackGetEmojisListURL,
		url.Values{
			"token": {environment.GetOauthAccessToken()},
		})

	if err != nil {
		return err
	}

	var responseJSON map[string]interface{}
	unmarshallData(resp.Body, &responseJSON)

	emojis := responseJSON["emoji"].(map[string]interface{})
	var emojisList []string
	for k := range emojis {
		emojisList = append(emojisList, k)
	}

	return dataStorage.AddCustomEmojis(emojisList)
}

func handleEvents(w http.ResponseWriter, req *http.Request) {
	var messagetype messageType
	reader := unmarshallData(req.Body, &messagetype)

	if messagetype.Type == "url_verification" {
		handleURLVerification(w, reader)
	} else {
		if messagetype.Type != "event_callback" {
			panic(fmt.Sprint("Not `event_callback`", messagetype.Type))
		}
		handleEvent(reader)
	}

	w.Write([]byte("GnocchettiAlVapore"))
}

func unmarshallData(reader io.Reader, v interface{}) io.Reader {
	var data bytes.Buffer
	data.ReadFrom(reader)
	err := json.Unmarshal(data.Bytes(), &v)

	if err != nil {
		fmt.Println("Unable to unmarshal data for type", reflect.TypeOf(v).String())
	}
	return bytes.NewBuffer(data.Bytes())
}

func handleURLVerification(w http.ResponseWriter, reader io.Reader) {

	var urlverification urlVerification
	unmarshallData(reader, &urlverification)

	if environment.GetSlackToken() != urlverification.Token {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		return
	}

	w.Write([]byte(urlverification.Challenge))
}

func postToSlack(token string, url string, w io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, w)

	if err != nil {
		fmt.Println("Error creating request")
		return nil, err
	}

	// Add Authorization token
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("Content-type", "application/json")

	client := http.Client{}
	clientResponse, clientError := client.Do(request)

	if clientError != nil {
		fmt.Println("Errore dal client")
		return nil, clientError
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
	if parseMessage(msg) {
		return
	}

	if msg.Event.Type == "message" && msg.Event.User != environment.GetBotID() {
		fmt.Println("DDDDD", msg.Event.Text)
		sendMessageToUser(strings.ToUpper(msg.Event.Text), msg.Event.Channel)
	}
}

func sendMessageToUser(message string, channel string) {
	messageToUser := sendMessageToUserStruct{
		Token:   environment.GetSlackToken(),
		Channel: channel,
		Text:    message,
		AsUser:  true,
	}

	marshalled, _ := json.Marshal(messageToUser)
	stringBuffer := bytes.NewBuffer(marshalled)
	req, err := postToSlack(environment.GetOauthToken(), slackChatPostMessageURL, stringBuffer)

	var responseData clientResponseData
	unmarshallData(req.Body, &responseData)

	if !responseData.Ok {
		fmt.Println("handleEvent error", responseData.Err)
	}
	if err != nil {
		fmt.Println("Handle event error", err)
	}
}

func parseMessage(msg message) bool {
	text := msg.Event.Text

	if strings.HasPrefix(text, "list groups") {
		handleListGroups(msg)
		return true
	} else if strings.HasPrefix(text, "list emojis") {
		handleListEmojisForGroup(msg)
		return true
	}
	return false
}

func handleListEmojisForGroup(msg message) {
	split := strings.Split(msg.Event.Text, " ")
	group := split[len(split)-1]

	userGroups := dataStorage.GetGroupsForUser(msg.Event.User)
	found := false

	for _, userGroup := range userGroups {
		if userGroup == group {
			found = true
			break
		}
	}

	if !found {
		sendMessageToUser("Group not found", msg.Event.Channel)
		return
	}

	emojis := dataStorage.GetEmojisForUserForGroup(msg.Event.User, group)

	if len(emojis) == 0 {
		sendMessageToUser(fmt.Sprintf("No emojis for group %s", group), msg.Event.Channel)
		return
	}

	sendMessageToUser(strings.Join(emojis, " "), msg.Event.Channel)
}

func handleListGroups(msg message) {
	groupsList := dataStorage.GetGroupsForUser(msg.Event.User)

	sendMessageToUser(strings.Join(groupsList, ", "), msg.Event.Channel)
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
	token, err := dataStorage.GetUserToken(info.User.ID)

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

	err := dataStorage.AddUserToken(uuid, info.User.ID)
	if err != nil {
		return err
	}

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

	userID, err := dataStorage.PopUserToken(state)
	if err != nil {
		fmt.Println(err)
	}

	var accessTokenData accessToken
	unmarshallData(resp.Body, &accessTokenData)

	dataStorage.AddUserToken(userID, accessTokenData.AccessToken)
}
