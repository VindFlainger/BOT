package apps

import (
	"PIE_BOT/Code/commands/apps/news"
	"PIE_BOT/Code/commands/chat"
	"PIE_BOT/Code/commands/utils"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"sync"
	"time"
)

type AddEvent struct {
	ChatID  int
	EventID int
	Params  []byte
}

type StopEvent struct {
	ChatID  int
	EventID int
}

func HandleNews(vkclient *chat.VKClient, GroupID int, vkMessage chan chat.MessParams, addCH <-chan *AddEvent, stopCH <-chan *StopEvent) {
	newsSubs := &sync.Map{}
	for {
		select {
		case <-time.After(time.Second * 1):
			newsSubs.Range(func(key, value interface{}) bool {
				params := value.(*news.Params)
				pages, err := params.ParsePage()
				if err == nil {
					for _, page := range pages {
						err := func() error {
							photo, err := utils.DownloadFile(page.PhotoURL)
							if err != nil {
								return err
							}
							filename := fmt.Sprintf("templates\\%d_%d.jpg", time.Now().Unix(), rand.Int())
							file, err := os.OpenFile(filename, os.O_CREATE, fs.ModePerm)
							if err != nil {
								return err
							}
							if n, err := file.Write(photo); n == 0 && err != nil {
								return err
							}

							uplShort, err := vkclient.SaveImage(filename, GroupID)
							if err != nil {
								return err
							}
							vkMessage <- &chat.PhotoMessage{
								MessageForm: &chat.Message{
									Target:   chat.T_CHAT,
									TargetID: key.(int),
									Message:  fmt.Sprintf("%s\n%s", page.Text, page.Date)},
								Photo: uplShort,
							}
							return nil
						}()
						if err != nil {
							vkMessage <- &chat.Message{
								Target:   chat.T_CHAT,
								TargetID: key.(int),
								Message:  fmt.Sprintf("%s\n%s", page.Text, page.Date),
							}
						}
					}
				} else {
					if !params.Error {
						vkMessage <- &chat.Message{
							Target:   chat.T_CHAT,
							TargetID: key.(int),
							Message:  fmt.Sprintf(err.Error())}
						params.Error = true
					}
				}
				return true
			})

		case event := <-addCH:
			params, err := news.ParseJson(event.Params)
			if err != nil {
				vkMessage <- &chat.Message{Message: err.Error(), Target: chat.T_CHAT, TargetID: event.ChatID}
			}
			newsSubs.Store(event.ChatID, params)

		case stop := <-stopCH:
			newsSubs.Delete(stop.ChatID)
		}

	}
}
