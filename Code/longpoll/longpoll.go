package longpoll

import (
	"errors"
	"github.com/go-vk-api/vk"
	"time"
)

type LongPoll struct {
	VKLongPoll *VKLongPoll

	TimeStart int64
}

type VKLongPoll struct {
	Client  *vk.Client
	GroupID int64
	Server  string
	TS      string
	Key     string
}

func GetLongPoll() *LongPoll {
	return &LongPoll{TimeStart: time.Now().Unix()}
}

func (longpoll *LongPoll) AddVKLongPoll(client *vk.Client, GroupID int64) error {

	response := map[string]string{}
	if err := client.CallMethod("groups.getLongPollServer",
		vk.RequestParams{"group_id": GroupID},
		&response); err != nil {
		return errors.New(VKGetLongPollServerError)
	}

	longpoll.VKLongPoll = &VKLongPoll{Client: client, Server: response["server"], TS: response["ts"], Key: response["key"], GroupID: GroupID}
	return nil

}
