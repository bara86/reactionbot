package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/satori/go.uuid"
)

const (
	slackEphemeralURL = "https://slack.com/api/chat.postEphemeral"

	slackOauthURL = "slack.com/oauth/authorize"

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

type addReactionActionMessage struct {
	Timestamp string `json:"ts"`
}

type addReactionActionUser struct {
	Id   string
	Name string
}

type addReactionAction struct {
	ResponseURL string                   `json:"response_url"`
	Channel     channel                  `json:"channel"`
	Message     addReactionActionMessage `json:"message"`
	User        addReactionActionUser    `json:"user"`
}

type UserStorage struct {
	keys sync.Map
}

func (u *UserStorage) Lookup(user string) bool {
	_, ok := u.keys.Load(user)
	return ok
}

func (u *UserStorage) Add(user string, token string) {
	u.keys.Store(user, token)
}

func (u *UserStorage) Remove(user string) {
	u.keys.Delete(user)
}

func (u *UserStorage) Get(user string) (string, error) {
	value, ok := u.keys.Load(user)
	if !ok {
		return "", fmt.Errorf("No user %s", user)
	}

	return value.(string), nil
}

func (u *UserStorage) Pop(user string) (string, error) {
	value, err := u.Get(user)
	if err == nil {
		u.Remove(user)
	}
	return value, err
}

var userTokens UserStorage
var temporaryTokens UserStorage

func addReactionToMessage(payload *string) {
	var info addReactionAction
	json.Unmarshal([]byte(*payload), &info)
	token, err := userTokens.Get(info.User.Id)

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
	q.Add("client_id", getClientID())
	q.Add("scope", reactionsWriteScope)
	q.Add("redirect_uri", fmt.Sprintf("%v/oauth", getAppURL()))
	q.Add("state", uuid)
	url.RawQuery = q.Encode()

	temporaryTokens.Add(uuid, info.User.Id)

	jsonMsg := fmt.Sprintf(authorizeButton, getOauthToken(), info.Channel.ID, info.User.Id, url.String())
	fmt.Println(jsonMsg)
	buf := bytes.NewBufferString(jsonMsg)
	_, err2 := postToSlack(getOauthToken(), slackEphemeralURL, buf)
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
			"client_id":     {getClientID()},
			"client_secret": {getClientSecret()},
			"code":          {string(query["code"][0])},
			"redirect_url":  {fmt.Sprintf("%v/oauth", getAppURL())},
		})
	fmt.Println(resp.Body)

	type accessToken struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
	}

	userID, err := temporaryTokens.Pop(state)
	if err != nil {
		fmt.Println(err)
	}

	var accessTokenData accessToken
	unmarshallData(resp.Body, &accessTokenData)

	userTokens.Add(userID, accessTokenData.AccessToken)
}
