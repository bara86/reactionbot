package commonstructure

type Storage interface {
	LookupUserToken(id string) (bool, error)
	AddUserToken(id string, token string) error
	RemoveUserToken(id string) error
	GetUserToken(id string) (string, error)
	PopUserToken(id string) (string, error)
	LoadEmojisList() error
	LookupEmoji(name string) (bool, error)
	AddCustomEmojis(emojisList []string) error
	GetGroupsForUser(id string) []string
	GetEmojisForUserForGroup(userID string, groupName string) []string
}
