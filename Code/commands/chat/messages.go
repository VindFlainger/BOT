package chat

import (
	"fmt"
	"github.com/go-vk-api/vk"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

var (
	ERR_PREFIX          = "[ERROR]: "
	SendMessageID int64 = 0
)

const (
	T_CHAT = "chat_id"
	T_USER = "peer_id"
)

//
type Message struct {
	Target   string
	TargetID int
	Message  string
}

type AdvMessage struct {
	Target    string
	TargetID  int
	Message   []interface{}
	Separator string
}

type PhotoMessage struct {
	MessageForm *Message
	Photo       string
}

type AdvPhotoMessage struct {
	MessageForm *AdvMessage
	Photo       string
}

//	Interface for working with Message, AdvMessage, PhotoMessage and AdvPhotoMessage
type MessParams interface {
	MakeMessage() *vk.RequestParams
}

func getrandomid() int {
	rand.Seed(time.Now().UnixMicro())
	return rand.Intn(100000)
}

//	Conver coming message slice into human-readable form
//	Convert rules:
//	1) separator between all elements of message 			Exc: message[0]<separator>message[1]...
//	2) separator between all elements of slice element 		Exc: message[0][0]<separator>message[0][1]<separator>message[1]
//	3) sumbol '\n' between all map pairs 					Exc: message[0]["1"]\nmessage[0]["2"]<separator>message[1]
func sendformat(separator string, message ...interface{}) (strtosend string) {
	for _, elem := range message {
		switch reflect.TypeOf(elem).Kind() {
		case reflect.Int:
			strtosend += strconv.Itoa(elem.(int)) + separator
		case reflect.Float64:
			strtosend += strconv.FormatFloat(elem.(float64), 'f', 4, 64) + separator
		case reflect.String:
			strtosend += elem.(string) + separator
		case reflect.Slice:
			sl := reflect.ValueOf(elem)
			count := 0
			for sl.Len() > count {
				strtosend += fmt.Sprintf("%v%s", sl.Index(count), separator)
				count++
			}
		case reflect.Map:
			it := reflect.ValueOf(elem).MapRange()
			for it.Next() {
				strtosend += fmt.Sprintf("%v : %v \n", it.Key(), it.Value())
			}
			strtosend += separator
		}
	}
	return strtosend
}

func (vkclient *VKClient) sendmessage(vkparams *vk.RequestParams, sendMessageID *int64) error {
	return vkclient.Client.CallMethod("messages.send", *vkparams, sendMessageID)
}

func (advms *AdvMessage) MakeMessage() *vk.RequestParams {
	return &vk.RequestParams{
		advms.Target: advms.TargetID,
		"message":    sendformat(advms.Separator, advms.Message...),
		"random_id":  getrandomid(),
	}
}

func (ms *Message) MakeMessage() *vk.RequestParams {
	return &vk.RequestParams{
		ms.Target:   ms.TargetID,
		"message":   ms.Message,
		"random_id": getrandomid(),
	}
}

func (phms *PhotoMessage) MakeMessage() *vk.RequestParams {
	return &vk.RequestParams{
		phms.MessageForm.Target: phms.MessageForm.TargetID,
		"message":               phms.MessageForm.Message,
		"random_id":             getrandomid(),
		"attachment":            phms.Photo,
	}
}

func (advphms *AdvPhotoMessage) MakeMessage() *vk.RequestParams {
	return &vk.RequestParams{
		advphms.MessageForm.Target: advphms.MessageForm.TargetID,
		"message":                  sendformat(advphms.MessageForm.Separator, advphms.MessageForm.Message...),
		"random_id":                getrandomid(),
		"attachment":               advphms.Photo,
	}
}
