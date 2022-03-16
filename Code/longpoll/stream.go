package longpoll

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type VKEvents struct {
	Type  uint8        `json:"type"`
	TS    string       `json:"ts"`
	Event []*VKContent `json:"updates"`
}

var vkevent = &VKEvents{}

func (ve *VKEvents) Reset() {
	*ve = *vkevent
}

type VKContent struct {
	Object *VKObject `json:"object"`
	Type   string    `json:"type"`
}

type VKObject struct {
	Message *VKMessage `json:"message"`
}

type VKMessage struct {
	FromID      int              `json:"from_id"`
	ChatID      int              `json:"peer_id"`
	Text        string           `json:"text"`
	Attachments []*VKAttachments `json:"attachments"`
	Reply       *VKReply         `json:"reply_message"`
	Action      *VKAction        `json:"action"`
	MarkID      int              `json:"id"`
}

type VKAttachments struct {
	Type    string     `json:"type"`
	Photo   *VKPhoto   `json:"photo"`
	Sticker *VKSticker `json:"sticker"`
	Voice   *VKVoice   `json:"audio_message"`
	Doc     *VkDoc     `json:"doc"`
}

type VKPhoto struct {
	Content   []*VKURL `json:"sizes"`
	AccessKey string   `json:"access_key"`
	ID        int      `json:"id"`
	OwnerID   int      `json:"owner_id"`
}

type VKURL struct {
	URL string `json:"url"`
}

type VKSticker struct {
	ProductID  int             `json:"product_id"`
	StickerID  int             `json:"sticker_id"`
	StickerIMG *VKStickerImage `json:"images"`
}

type VKStickerImage struct {
	URL string `json:"url"`
}

type VKVoice struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	Duration  uint16 `json:"duration"`
	AccessKey string `json:"access_key"`
	URL       string `json:"link_mp3"`
}

type VKReply struct {
	FromID         int    `json:"from_id"`
	Text           string `json:"text"`
	ConversationID int    `json:"conversation_message_id"`
}

type VKAction struct {
	Type   string `json:"type"`
	UserID int    `json:"member_id"`
	Text   string `json:"text"`
}

type VkDoc struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Size  int    `json:"size"`
}

type Stream struct {
	Updates <-chan interface{}
	Errors  <-chan error
}

func (lp *LongPoll) StartStreaming(stream *Stream) {
	updatesch := make(chan interface{}, 100)
	errorsch := make(chan error, 100)
	stream.Updates = updatesch
	stream.Errors = errorsch

	if lp.VKLongPoll != nil {
		go func() {
			respparse := &VKEvents{}
			for {
				respparse.Reset()
				url := fmt.Sprintf("https://%s?act=a_check&key=%s&ts=%s&wait=20&mode=2&version=2", strings.Replace(lp.VKLongPoll.Server, "https://", "", 10),
					lp.VKLongPoll.Key, lp.VKLongPoll.TS)
				response, err := http.Get(url)
				if err != nil || response.StatusCode != 200 {
					errorsch <- errors.New(RESPONSEEROR)
					time.Sleep(time.Second * 10)
					continue
				}

				json.NewDecoder(response.Body).Decode(respparse)
				lp.VKLongPoll.TS = respparse.TS

				if respparse.Type != 0 {
					lp.AddVKLongPoll(lp.VKLongPoll.Client, lp.VKLongPoll.GroupID)
					continue
				}

				for _, event := range respparse.Event {
					if event.Object.Message != nil && event.Object.Message.MarkID == 0 {
						updatesch <- event.Object
					}
				}
			}
		}()
	}

}

func (lp *LongPoll) GetStream() *Stream {
	stream := &Stream{}
	lp.StartStreaming(stream)
	return stream
}
