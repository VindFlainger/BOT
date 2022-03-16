package chat

import (
	"fmt"
	"github.com/go-vk-api/vk"
	"reflect"
	"strconv"
)

type Users struct {
	Count     int
	UsersIDs  []int
	OwnerID   int
	AdminsIDs []int
}

//	Realize "messages.getConversationMembers" request for given ChatID
//	Returns *Users where all chat members are in UsersIDs(including admin and owner), all chat admins in AdminsIDs
//	!OwnerID field contains OwnerID if it was in chat during API call
func (vklient *VKClient) GetChatUsers(ChatID int) (*Users, error) {
	var resp interface{}
	users := &Users{}

	if err := vklient.Client.CallMethod("messages.getConversationMembers",
		vk.RequestParams{
			"peer_id": ChatID}, &resp); err != nil {
		return users, err
	}

	refresp := reflect.ValueOf(resp)
	items := refresp.MapIndex(reflect.ValueOf("items")).Interface()
	users.Count = int(refresp.MapIndex(reflect.ValueOf("count")).Interface().(float64))

	for _, user := range items.([]interface{}) {
		urval := reflect.ValueOf(user)
		userid := int(urval.MapIndex(reflect.ValueOf("member_id")).Interface().(float64))
		users.UsersIDs = append(users.UsersIDs, userid)
		if !(urval.MapIndex(reflect.ValueOf("is_admin")).Kind() == reflect.Invalid) {
			users.AdminsIDs = append(users.AdminsIDs, userid)
		}
		if !(urval.MapIndex(reflect.ValueOf("is_owner")).Kind() == reflect.Invalid) {
			users.OwnerID = userid
		}
	}

	return users, nil
}

//	Wrapper over GetChatUsers()
//	Method searching UserID matches in Users.AdminsIDs
//	It true when UserID in Users.AdminsIDs and false when not or GetChatUsers() returned error
func (vkclient *VKClient) IsAdmin(ChatID, UserID int) bool {
	users, err := vkclient.GetChatUsers(ChatID)
	if err != nil {
		return false
	}

	for _, adid := range users.AdminsIDs {
		if adid == UserID {
			return true
		}
	}
	return false

}

//	Wrapper over GetChatUsers()
//	Method compare UserID with Users.OwnerID
//	It true when UserID = Users.OwnerID and false when not or GetChatUsers() returned error
func (vkclient *VKClient) IsOwner(ChatID, UserID int) bool {
	users, err := vkclient.GetChatUsers(ChatID)
	if err != nil || users.OwnerID != UserID {
		return false
	}
	return true
}

//	Realize "users.get" request for given UserID
//	Returns string combination "firstname surname" or string represention of UserID if API call returns error
func (vkclient *VKClient) GetName(UserID int) string {
	var resp interface{}
	if err := vkclient.Client.CallMethod("users.get", vk.RequestParams{"user_ids": UserID}, &resp); err != nil {
		return strconv.Itoa(UserID)
	}
	usercontent := resp.([]interface{})[0].(map[string]interface{})
	username := usercontent["first_name"].(string) + " " + usercontent["last_name"].(string)
	return username
}

//	Wrapper over GetName(), realize GetName() on all UserIDs IDs
//	Returns slice with UserIDs capacity of usernames in format "firstname surname"
//	Pay attention for running time, if take 1300 ms for 20 IDs
func (vkclient *VKClient) GetNames(UserIDs []int) []string {
	usernames := make([]string, 0, len(UserIDs))
	for _, id := range UserIDs {
		usernames = append(usernames, vkclient.GetName(id))
	}
	return usernames
}

//	Realize "users.get" request for given UserID
//	Returns string combination "[UserID|domain]"
//	Pay attention for running time, if take 1300 ms for 20 IDs
func (vkclient *VKClient) GetDomain(UserID int) string {
	var resp interface{}
	if err := vkclient.Client.CallMethod("users.get", vk.RequestParams{"user_ids": UserID, "fields": "domain"}, &resp); err != nil {
		return strconv.Itoa(UserID)
	}
	domain := reflect.ValueOf(resp.([]interface{})[0]).MapIndex(reflect.ValueOf("domain")).Interface().(string)
	callabledomain := fmt.Sprintf("[id%s|%s]", strconv.Itoa(UserID), domain)
	return callabledomain
}

//	Wrapper over GetDomain(), realize GetDomain() on all UserIDs IDs
//	Returns slice with UserIDs capacity of domains in format "[UserID|domain]"
func (vkclient *VKClient) GetDomains(UserIDs []int) []string {
	domains := make([]string, 0, len(UserIDs))
	for _, id := range UserIDs {
		domains = append(domains, vkclient.GetDomain(id))
	}
	return domains
}
