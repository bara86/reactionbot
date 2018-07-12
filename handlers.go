package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
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

var userTokens UserStorage

func addReactionToMessage(payload *string) {
	var info addReactionAction
	json.Unmarshal([]byte(*payload), &info)
	if !userTokens.Lookup(info.User.Id) {
		if postEphemeralMessage(&info) != nil {

		}
	}
	// addReaction("heart", info.Message.Timestamp, info.Channel.ID)

}

func postEphemeralMessage(info *addReactionAction) error {
	fmt.Println("Post ephemeral message")

	url := url.URL{Path: slackOauthURL, Scheme: "https"}
	q := url.Query()
	q.Add("client_id", getClientID())
	q.Add("scope", reactionsWriteScope)
	q.Add("redirect_uri", fmt.Sprintf("%v/oauth", getAppURL()))
	url.RawQuery = q.Encode()

	jsonMsg := fmt.Sprintf(authorizeButton, getOauthToken(), info.Channel.ID, info.User.Id, url.String())
	fmt.Println(jsonMsg)
	buf := bytes.NewBufferString(jsonMsg)
	_, err2 := postToSlack(slackEphemeralURL, buf)
	if err2 != nil {
		return err2
	}
	return nil
}
