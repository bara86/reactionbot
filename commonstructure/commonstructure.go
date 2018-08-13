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
	AddGroupForUser(idUser string, groupName string) error
	AddEmojiForGroupForUser(emojiName string, groupName string, idUser string) error
	RemoveEmojiFromGroupForUser(emojiName string, groupName string, idUser string) error
	LookupForUserGroup(userID string, groupName string) (bool, error)
	RemoveGroupForUser(userID string, groupName string) error
}
