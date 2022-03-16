package mainattrs

import (
	"PIE_BOT/Code/commands/chat"
	"time"
)

type StartMessage struct {
	Greeting   string
	Deployment string
	FirstInfo  []string
	Commands   string
}

//	DB Const
const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "gemger2003"
	DBNAME   = "pie_bot"
)

var (
	StartMess = &StartMessage{
		Greeting:   "Hello! I'm PIE_BOT",
		Deployment: "Please wait a bit for the configuration of your conversation",
		FirstInfo: []string{
			"The initial configuration of the bot has been successfully completed!",
			"PIE_BOT: Now you have full access to all the adaptations I can offer.",
			"PIE_BOT: For all questions and problems that have arisen, feel free to write to the community's private messages.",
			"PIE_BOT: Important information, I am a chatbot and I do not respond to all personal messages except commands for subscribing to news.",
			"PIE_BOT: You can subscribe to the news simply by writing to me in private messages +news or -news to unsubscribe.",
			"PIE_BOT: Have a good day"},
	}
)

func Sayhello(sm *StartMessage, chatID int, mh chan<- chat.MessParams) {
	mh <- &chat.Message{Message: sm.Greeting, TargetID: chatID, Target: chat.T_CHAT}
	time.Sleep(time.Second)
	mh <- &chat.Message{Message: sm.Deployment, TargetID: chatID, Target: chat.T_CHAT}
	time.Sleep(time.Second * 5)
	for _, mess := range sm.FirstInfo {
		mh <- &chat.Message{Message: mess, TargetID: chatID, Target: chat.T_CHAT}
		time.Sleep(time.Second * 5)
	}
	time.Sleep(time.Second * 5)
	mh <- &chat.Message{Message: sm.Commands, TargetID: chatID, Target: chat.T_CHAT}
}
