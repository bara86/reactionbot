package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

const (
	slackEphemeralUrl = "https://slack.com/api/chat.postEphemeral"

	authorizeButton = `{
		"token": "%v",
		"channel": "%v",
		"text": "We need your approval to post reaction with your name.",
		"user": "%v"
	}`
)

/*
"attachments": [
			{
				"fallback": "Approve posting a reaction with your name https://foo.bar",
				"callback_id": "mignatta",
				"actions": [
					{
						"type": "button",
						"text": "Approve",
						"name": "approve"
					}
				]
			}
		]
*/

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

	//addReaction("heart", info.Message.Timestamp, info.Channel.ID)
}

func postEphemeralMessage(info *addReactionAction) error {
	fmt.Println("Post ephemeral message")

	jsonMsg := fmt.Sprintf(authorizeButton, getOauthToken(), info.Channel.Name, info.User.Id)
	buf := bytes.NewBufferString(jsonMsg)
	_, err2 := postToSlack(slackEphemeralUrl, buf)
	if err2 != nil {
		return err2
	}
	return nil
}
