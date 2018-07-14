package handlers

type addReactionActionMessage struct {
	Timestamp string `json:"ts"`
}

type addReactionActionUser struct {
	ID   string
	Name string
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

type messageType struct {
	Type string
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
