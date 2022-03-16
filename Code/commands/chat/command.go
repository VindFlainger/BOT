package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-vk-api/vk"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type VKClient struct {
	Client *vk.Client
}

type ParseST struct {
	Len   int
	PType []reflect.Kind
}

func ParseShortName(shortname string) (UserID int, err error) {
	ind := strings.Index(shortname, "|")
	if ind <= 3 {
		err = BADSHORTNAMEFORMATERROR
		return
	}
	id64form, err := strconv.ParseInt(shortname[3:ind], 10, 32)
	if err != nil {
		err = BADSHORTNAMEFORMATERROR
		return
	}
	UserID = int(id64form)
	return
}

func VkUpload(url string, img io.Reader) (server int, photo, hash string, err error) {
	type UploadResponse struct {
		Server int    `json:"server"`
		Photo  string `json:"photo"`
		Hash   string `json:"hash"`
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, img); err != nil {
		return
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	uplRes := UploadResponse{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&uplRes)
	if err != nil {
		return
	}
	defer res.Body.Close()

	server = uplRes.Server
	photo = uplRes.Photo
	hash = uplRes.Hash
	return
}

func (vkclient *VKClient) SendMessage(params MessParams) error {
	return vkclient.sendmessage(params.MakeMessage(), &SendMessageID)
}

func (vkclient *VKClient) SaveImage(name string, groupID int) (string, error) {
	var resp interface{}
	vkclient.Client.CallMethod("photos.getMessagesUploadServer", vk.RequestParams{"group_id": groupID}, &resp)
	upload_url := resp.(map[string]interface{})["upload_url"].(string)
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	server, photo, hash, err := VkUpload(upload_url, file)
	if err != nil {
		return "", err
	}
	err = vkclient.Client.CallMethod("photos.saveMessagesPhoto", vk.RequestParams{"photo": photo, "server": server, "hash": hash}, &resp)
	if err != nil {
		return "", err
	}
	id := int(resp.([]interface{})[0].(map[string]interface{})["id"].(float64))
	owner_id := int(resp.([]interface{})[0].(map[string]interface{})["owner_id"].(float64))

	return fmt.Sprintf("photo%d_%d", owner_id, id), nil
}

//	Method Deprecated
func (vkclient *VKClient) GetSaveImage(name string) (string, error) {
	savefile, err := os.OpenFile(name, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return "", IMAGEFILEERROR
	}
	stat, _ := savefile.Stat()
	buf := make([]byte, stat.Size())
	savefile.Read(buf)
	for _, line := range strings.Split(string(buf), "\n") {
		if strings.Contains(line, name) {
			return strings.TrimRight(strings.Split(line, "\t")[1], "\n"), nil
		}
	}
	return "", NOSAVEIMAGEERROR
}

func (vkclient *VKClient) RemoveFromChat(UserID int, ChatID int) error {
	var resp interface{}

	err := vkclient.Client.CallMethod("messages.removeChatUser", vk.RequestParams{"chat_id": ChatID, "member_id": UserID}, &resp)
	if err != nil {
		return err
	}
	return nil
}
