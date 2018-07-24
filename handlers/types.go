package handlers

type addReactionActionMessage struct {
	Timestamp string `json:"ts"`
}

type addReactionActionUser struct {
	ID   string
	Name string
}

type sendMessageToUser struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	AsUser  bool   `json:"as_user"`
}

type addReactionAction struct {
	ResponseURL string                   `json:"response_url"`
	Channel     channel                  `json:"channel"`
	Message     addReactionActionMessage `json:"message"`
	User        addReactionActionUser    `json:"user"`
}

type event struct {
	Type    string
	EventTs string `json:"event_ts"`
	User    string
	Item    string
	Ts      string
	Channel string
	Text    string `json:"text"`
}

type message struct {
	Token     string
	APIAppID  string `json:"api_app_id"`
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

type messageType struct {
	Token string
	Type  string
}

type urlVerification struct {
	Token     string
	Challenge string
}

type callbackID struct {
	CallbackID string `json:"callback_id"`
}

type accessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

type postMessageResponseMessage struct {
	BotID string `json:"bot_id"`
}

type postMessageResponse struct {
	Ok      bool
	Message postMessageResponseMessage
}
