package main

import (
	"PIE_BOT/Code/commands/apps"
	"PIE_BOT/Code/commands/chat"
	"PIE_BOT/Code/commands/dbcom"
	"PIE_BOT/Code/commands/utils"
	"PIE_BOT/Code/language"
	"PIE_BOT/Code/longpoll"
	"PIE_BOT/mainattrs"
	"database/sql"
	"fmt"
	"github.com/go-vk-api/vk"
	_ "github.com/lib/pq"

	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func messageserver(vkclient *chat.VKClient, mh <-chan chat.MessParams) {
	for message := range mh {
		err := vkclient.SendMessage(message)
		if err != nil {
			log.Println(err)
		}
	}
}

func mainworker(vkclient *chat.VKClient, db *sql.DB, chsettings *sync.Map, vkobjchan <-chan *longpoll.VKObject) {
	var cd = &dbcom.ChatDB{DB: db, Count: 1}
	var err error

	//	Handlers channels
	vkMessage := make(chan chat.MessParams, 1000)
	stopCH := make(chan *apps.StopEvent, 1000)
	addCH := make(chan *apps.AddEvent, 1000)

	// Starting Handlers
	go messageserver(vkclient, vkMessage)
	go apps.HandleNews(vkclient, 202265867, vkMessage, addCH, stopCH)

	ENL, ENFB, err := language.ConfigureLang("en")
	if err != nil {
		panic("Critical: Default configuration error")
	}

	RUL, RUFB, err := language.ConfigureLang("ru")
	if err != nil {
		panic("Critical: Default configuration error")
	}

	lang := map[int]*language.Lang{
		1: RUL,
		2: ENL,
	}
	fb := map[int]*language.FeedBack{
		1: RUFB,
		2: ENFB,
	}

	for vkobj := range vkobjchan {
		vkobj.Message.ChatID -= 2000000000
		err = func() error {
			if _, ok := chsettings.Load(vkobj.Message.ChatID); !ok {

				params, err := cd.InputSettings(vkobj.Message.ChatID)

				if err != nil {
					users, err := vkclient.GetChatUsers(vkobj.Message.ChatID)
					if err != nil {
						return err
					}
					params = &dbcom.ChatParams{Status: dbcom.ST_DEF, MaxWarn: 5, Lang: dbcom.LANG_EN}
					_, err = cd.AddNewChat(vkobj.Message.ChatID, users.OwnerID, params)
					if err != nil {
						return err
					}

					go mainattrs.Sayhello(mainattrs.StartMess, vkobj.Message.ChatID, vkMessage)

				}
				chsettings.Store(vkobj.Message.ChatID, params)
			}
			return nil
		}()

		v, _ := chsettings.Load(vkobj.Message.ChatID)
		st := v.(*dbcom.ChatParams)

		mess := strings.Split(vkobj.Message.Text, " ")
		if strings.HasPrefix(mess[0], "-") {
			switch mess[0] {

			case lang[st.Lang].Ping: // -ping
				vkMessage <- &chat.Message{Message: fb[st.Lang].PingFB, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

			case lang[st.Lang].Repeat: // -repeat {...var}
				vkMessage <- &chat.Message{Message: strings.Repeat(strings.Join(mess[1:], " ")+"\n", 10), TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

			case lang[st.Lang].Members: // -members
				err = func() error {
					if tout := dbcom.GetTimeOut(db, vkobj.Message.ChatID, dbcom.TO_Members, dbcom.TOID_Members); tout != 0 {
						vkMessage <- &chat.Message{Message: fmt.Sprintf(fb[st.Lang].TimeOutFB, tout), TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
						return nil
					}

					users, err := vkclient.GetChatUsers(vkobj.Message.ChatID + 2000000000)
					if err != nil {
						return err
					}

					var memlist string
					for _, usid := range users.UsersIDs {
						memlist += fmt.Sprintf("\n%s : %s", vkclient.GetName(usid), vkclient.GetDomain(usid))
						for _, adid := range users.AdminsIDs {
							if usid == adid {
								memlist += "\tAdmin"
							}
						}
						if usid == users.OwnerID {
							memlist += "\tOwner"
						}
					}
					vkMessage <- &chat.Message{Message: memlist, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

					return nil
				}()

			case lang[st.Lang].Devote: // -devote @UserID group
				err = func() error {
					access, err := cd.CheckAllPermissions(vkobj.Message.ChatID, vkobj.Message.FromID, dbcom.P_ROLE)
					if err != nil {
						return err
					}

					if !access {
						return mainattrs.NewAccessErr(dbcom.P_ROLE, fb[st.Lang].AccessWarn, cd)
					}

					args, err := chat.ParseArgs(mess[:], chat.T_ignore, chat.T_UserID, chat.T_string)
					if err != nil {
						return err
					}

					var role_id int
					var user_id = args[0].(int)
					var groupname = args[1].(string)

					switch groupname {
					case "admin":
						role_id = dbcom.R_ADMIN
					case "stmoder":
						role_id = dbcom.R_STMODER
					case "moder":
						role_id = dbcom.R_MODER
					case "helper":
						role_id = dbcom.R_HELPER
					case "user":
						role_id = dbcom.R_USER
					default:
						return dbcom.RoleValueErr
					}

					if err := cd.WriteRole(vkobj.Message.ChatID, user_id, role_id); err != nil {
						return err
					}

					vkMessage <- &chat.Message{Message: fmt.Sprintf(fb[st.Lang].DevoteFB, vkclient.GetName(user_id)), Target: chat.T_CHAT, TargetID: vkobj.Message.ChatID}
					return nil

				}()

			case lang[st.Lang].Warn: //	-warn @UserID
				err = func() error {
					access, err := cd.CheckAllPermissions(vkobj.Message.ChatID, vkobj.Message.FromID, dbcom.P_WARN)
					if err != nil {
						return err
					}

					if !access {
						return mainattrs.NewAccessErr(dbcom.P_WARN, fb[st.Lang].AccessWarn, cd)
					}

					args, err := chat.ParseArgs(mess[:], chat.T_ignore, chat.T_UserID)
					if err != nil {
						return err
					}

					count, err := cd.Warn(vkobj.Message.ChatID, args[0].(int), st.MaxWarn, vkclient)
					if err != nil {
						return err
					}

					if count < st.MaxWarn {
						vkMessage <- &chat.Message{
							Message:  fmt.Sprintf(fb[st.Lang].WarnFB, vkclient.GetDomain(vkobj.Message.FromID), count, st.MaxWarn),
							TargetID: vkobj.Message.ChatID,
							Target:   chat.T_CHAT}
					} else {
						vkclient.SendMessage(&chat.Message{
							Message:  "Доделать",
							TargetID: vkobj.Message.ChatID,
							Target:   chat.T_CHAT})
					}

					return nil
				}()

			case lang[st.Lang].Unwarn: // - unwarn @Userid
				err = func() error {
					access, err := cd.CheckAllPermissions(vkobj.Message.ChatID, vkobj.Message.FromID, dbcom.P_WARN)
					if err != nil {
						return err
					}

					if !access {
						mainattrs.NewAccessErr(dbcom.P_WARN, fb[st.Lang].AccessWarn, cd)
					}

					args, err := chat.ParseArgs(mess[:], chat.T_ignore, chat.T_UserID)
					if err != nil {
						return err
					}

					count, err := cd.Unwarn(vkobj.Message.ChatID, args[0].(int))
					if err != nil {
						return err
					}

					if err != nil {
						return err
					}

					if count == -1 {
						vkMessage <- &chat.Message{
							Message:  fb[st.Lang].NoWarnsWarn,
							TargetID: vkobj.Message.ChatID,
							Target:   chat.T_CHAT}
					} else {
						vkMessage <- &chat.Message{
							Message: fmt.Sprintf(fb[st.Lang].UnwarnFB,
								vkclient.GetDomain(vkobj.Message.FromID),
								count,
								st.MaxWarn),
							TargetID: vkobj.Message.ChatID,
							Target:   chat.T_CHAT}
					}

					return nil
				}()

			case lang[st.Lang].Ban: // -ban @UserID Time reason
				err = func() error {
					access, err := cd.CheckAllPermissions(vkobj.Message.ChatID, vkobj.Message.FromID, dbcom.P_BAN)
					if err != nil {
						return err
					}

					if !access {
						return mainattrs.NewAccessErr(dbcom.P_BAN, fb[st.Lang].AccessWarn, cd)
					}

					args, err := chat.ParseArgs(mess, chat.T_ignore, chat.T_UserID, chat.T_Time, chat.T_string)
					if err != nil {
						return err
					}

					userid := args[0].(int)
					bantime := args[1].(int)
					reason := args[2].(string)

					err = cd.Ban(vkobj.Message.ChatID, userid, bantime, reason, vkclient)
					if err != nil {
						return err
					}

					vkMessage <- &chat.Message{
						Message: fmt.Sprintf(fb[st.Lang].BanFB,
							vkclient.GetName(userid),
							time.Unix(time.Now().Unix()+int64(bantime), 0).Format("15:04 02.01.06"),
							reason),
						TargetID: vkobj.Message.ChatID,
						Target:   chat.T_CHAT}

					return nil
				}()

			case lang[st.Lang].UnBan: //-unban @UserID
				err = func() error {
					access, err := cd.CheckAllPermissions(vkobj.Message.ChatID, vkobj.Message.FromID, dbcom.P_BAN)
					if err != nil {
						return err
					}

					if !access {
						return mainattrs.NewAccessErr(dbcom.P_BAN, fb[st.Lang].AccessWarn, cd)
					}

					args, err := chat.ParseArgs(mess, chat.T_ignore, chat.T_UserID)
					if err != nil {
						return err
					}

					userID := args[0].(int)

					err = cd.UnBan(vkobj.Message.ChatID, userID)
					if err != nil {
						return err
					}

					vkMessage <- &chat.Message{
						Message:  fmt.Sprintf(fb[st.Lang].UnBanFB, vkclient.GetName(userID)),
						TargetID: vkobj.Message.ChatID,
						Target:   chat.T_CHAT}

					return nil
				}()

			case lang[st.Lang].Myt: //myt @UserID Time reason
				err = func() error {
					access, err := cd.CheckAllPermissions(vkobj.Message.ChatID, vkobj.Message.FromID, dbcom.P_MYT)
					if err != nil {
						return err
					}

					if !access {
						return mainattrs.NewAccessErr(dbcom.P_MYT, fb[st.Lang].AccessWarn, cd)
					}

					args, err := chat.ParseArgs(mess, chat.T_ignore, chat.T_UserID, chat.T_Time, chat.T_string)
					if err != nil {
						return err
					}

					userid := args[0].(int)
					myttime := args[1].(int)
					reason := args[2].(string)

					err = cd.Myt(vkobj.Message.ChatID, userid, myttime, reason)
					if err != nil {
						return err
					}

					vkMessage <- &chat.Message{
						Message: fmt.Sprintf(fb[st.Lang].MytFB,
							vkclient.GetName(userid),
							time.Unix(time.Now().Unix()+int64(myttime), 0).Format("15:04 02.01.06"),
							reason),
						TargetID: vkobj.Message.ChatID,
						Target:   chat.T_CHAT}

					return nil
				}()

			case lang[st.Lang].UnMyt: // -unmyt
				err = func() error {
					access, err := cd.CheckAllPermissions(vkobj.Message.ChatID, vkobj.Message.FromID, dbcom.P_MYT)
					if err != nil {
						return err
					}

					if !access {
						return mainattrs.NewAccessErr(dbcom.P_MYT, fb[st.Lang].AccessWarn, cd)
					}

					args, err := chat.ParseArgs(mess, chat.T_ignore, chat.T_UserID)
					if err != nil {
						return err
					}

					userID := args[0].(int)

					err = cd.UnMyt(vkobj.Message.ChatID, userID)
					if err != nil {
						return err
					}

					vkMessage <- &chat.Message{
						Message:  fmt.Sprintf(fb[st.Lang].UnMytFB, vkclient.GetName(userID)),
						TargetID: vkobj.Message.ChatID,
						Target:   chat.T_CHAT}

					return nil
				}()

			case lang[st.Lang].WarnList: // -warns
				err = func() error {
					warns, err := cd.GetWarns(vkobj.Message.ChatID)
					if err != nil {
						return err
					}

					hrwarns := map[string]int{}
					for key, val := range warns {
						hrwarns[vkclient.GetName(key)] = val
					}

					vkclient.SendMessage(&chat.AdvMessage{
						Message:   []interface{}{fb[st.Lang].WarnsHat, hrwarns},
						Separator: "\n",
						TargetID:  vkobj.Message.ChatID,
						Target:    chat.T_CHAT})

					return nil
				}()

			case lang[st.Lang].BanList: // -bans
				err = func() error {
					bans, err := cd.GetBans(vkobj.Message.ChatID)
					if err != nil {
						return err
					}

					var hrbans []string
					for _, ban := range bans {
						hrbans = append(hrbans, fmt.Sprintf(fb[st.Lang].BansForm,
							vkclient.GetName(ban.UserID),
							time.Unix(int64(ban.Unbantime), 0).Format("15:04 02.01.06"),
							ban.Reason))
					}

					vkclient.SendMessage(&chat.AdvMessage{
						Message:   []interface{}{fb[st.Lang].BansHat, hrbans},
						Separator: "\n",
						TargetID:  vkobj.Message.ChatID,
						Target:    chat.T_CHAT})
					return nil
				}()

			case lang[st.Lang].MytList:
				err = func() error {
					myts, err := cd.GetMyts(vkobj.Message.ChatID)
					if err != nil {
						return err
					}

					var hrmyts []string
					for _, myt := range myts {
						hrmyts = append(hrmyts,
							fmt.Sprintf(fb[st.Lang].MytsForm,
								vkclient.GetName(myt.UserID),
								time.Unix(int64(myt.Unmyttime), 0).Format("15:04 02.01.06"),
								myt.Reason))
					}

					vkclient.SendMessage(&chat.AdvMessage{
						Message:   []interface{}{fb[st.Lang].MytsHat, hrmyts},
						Separator: "\n",
						TargetID:  vkobj.Message.ChatID,
						Target:    chat.T_CHAT})
					return nil
				}()

			case lang[st.Lang].AddEvent: //-addevent ID file
				err = func() error {
					args, err := chat.ParseArgs(mess, chat.T_ignore, chat.T_int)
					if err != nil {
						return err
					}

					event_id := args[0].(int)
					if !func() bool {
						for _, event := range dbcom.EVENTS {
							if event_id == event {
								return true
							}
						}
						return false
					}() {
						vkMessage <- &chat.Message{Message: fb[st.Lang].NoDefEventWarn, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
						return nil
					}

					for _, attach := range vkobj.Message.Attachments {
						if attach.Type == "doc" && strings.HasSuffix(attach.Doc.Title, ".json") {
							buf, err := utils.DownloadFile(attach.Doc.URL)
							if err != nil {
								return err
							}

							if err := cd.AddEvent(vkobj.Message.ChatID, event_id, buf); err != nil {
								return err
							}
							addCH <- &apps.AddEvent{vkobj.Message.ChatID, event_id, buf}

						}
					}

					return nil
				}()

			case "-русский":
				st.Lang = dbcom.LANG_RU
				chsettings.Store(vkobj.Message.ChatID, st)
				err = func() error {
					if cd.ChangeSettings(vkobj.Message.ChatID, st) != nil {
						return err
					}
					vkMessage <- &chat.Message{Message: "О! Вы не из англии?!", TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
					return nil
				}()

			case "-english":
				st.Lang = dbcom.LANG_EN
				chsettings.Store(vkobj.Message.ChatID, st)
				err = func() error {
					if cd.ChangeSettings(vkobj.Message.ChatID, st) != nil {
						return err
					}
					vkMessage <- &chat.Message{Message: "English language!", TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
					return nil
				}()
			}
		}

		if err != nil {
			switch err {

			case chat.TIMEFORMATERROR:
				vkMessage <- &chat.Message{Message: fb[st.Lang].TimeFormatErr, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

			case dbcom.NoBanErr:
				vkMessage <- &chat.Message{Message: fb[st.Lang].NoBansWarn, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

			case dbcom.NoMytErr:
				vkMessage <- &chat.Message{Message: fb[st.Lang].NoMytsWarn, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

			case chat.ARGSERROR:
				vkMessage <- &chat.Message{Message: fb[st.Lang].IncorrectAgrsErr, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

			case dbcom.BadFileContentErr:
				vkMessage <- &chat.Message{Message: fb[st.Lang].BadFileContentErr, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}

			default:
				if noaccess, ok := err.(*mainattrs.AccessErr); ok {
					vkMessage <- &chat.Message{Message: noaccess.Error(), TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
				}

				if vkerr, ok := err.(*vk.MethodError); ok {
					switch vkerr.Code {
					case 15:
						if !vkclient.IsAdmin(vkobj.Message.ChatID+2000000000, -202265867) {
							vkMessage <- &chat.Message{Message: fb[st.Lang].BotNotAdmERR, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
						} else {
							vkMessage <- &chat.Message{Message: fb[st.Lang].AccessWarn, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
						}
					case 935:
						vkMessage <- &chat.Message{Message: fb[st.Lang].NoSuchUserErr, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
					}

				} else {
					vkMessage <- &chat.Message{Message: fb[st.Lang].UnknownERR, TargetID: vkobj.Message.ChatID, Target: chat.T_CHAT}
					log.Println(err)
				}
			}
			err = nil
		}
	}
}

func main() {
	client, err := vk.NewClientWithOptions(vk.WithToken(
		"1f183c054754c05c1d595faabbbf0001eef36a99265568b78edd0c6f0b30837a1b4eb98d1a0ecacf5e7fe"),
		vk.WithHTTPClient(&http.Client{Timeout: time.Second * 30}))
	if err != nil {
		panic("Critical: Error during creating NewClient")
	}

	lp := longpoll.GetLongPoll()
	lp.AddVKLongPoll(client, 202265867)
	steam := lp.GetStream()

	postconf := fmt.Sprintf("host= %s port =%d user= %s password= %s dbname= %s sslmode=disable", mainattrs.HOST, mainattrs.PORT, mainattrs.USER, mainattrs.PASSWORD, mainattrs.DBNAME)
	db, err := sql.Open("postgres", postconf)
	vkclient := &chat.VKClient{Client: client}
	vkobjectch := make(chan *longpoll.VKObject, 1000)

	chsettings := &sync.Map{}

	go mainworker(vkclient, db, chsettings, vkobjectch)

	for {
		select {
		case newupdate := <-steam.Updates:
			switch content := newupdate.(type) {
			case *longpoll.VKObject:
				vkobjectch <- content
			}
		case newerror := <-steam.Errors:
			log.Println(newerror)

		}
	}

}
