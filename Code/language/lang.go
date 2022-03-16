package language

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Lang struct {
	Ping     string `json:"ping"`
	Repeat   string `json:"repeat"`
	Members  string `json:"members"`
	Devote   string `json:"devote"`
	Warn     string `json:"warn"`
	Unwarn   string `json:"unwarn"`
	Ban      string `json:"ban"`
	UnBan    string `json:"unban"`
	Myt      string `json:"myt"`
	UnMyt    string `json:"unmyt"`
	MytList  string `json:"mytlist"`
	WarnList string `json:"warnlist"`
	BanList  string `json:"banlist"`
	AddEvent string `json:"addevent"`
	DelEvent string `json:"delevent"`
}

type FeedBack struct {
	PingFB    string `json:"ping_fb"`
	WarnFB    string `json:"warn_fb"`
	UnwarnFB  string `json:"unwarn_fb"`
	UnBanFB   string `json:"unban_fb"`
	UnMytFB   string `json:"unmyt_fb"`
	TimeOutFB string `json:"timeout_fb"`
	BanFB     string `json:"ban_fb"`
	MytFB     string `json:"myt_fb"`
	DevoteFB  string `json:"devote_fb"`

	BotNotAdmERR      string `json:"botnotadm_err"`
	UnknownERR        string `json:"unknown_err"`
	NoSuchUserErr     string `json:"nosuchuser_err"`
	TimeFormatErr     string `json:"timeformat_err"`
	IncorrectAgrsErr  string `json:"argsformat_err"`
	BadFileContentErr string `json:"badfilecontent_err"`

	AccessWarn     string `json:"access_warn"`
	NoPermWarn     string `json:"noperm_warn"`
	NoWarnsWarn    string `json:"nowarns_warn"`
	NoBansWarn     string `json:"nobans_warn"`
	NoMytsWarn     string `json:"nomyts_warn"`
	NoDefEventWarn string `json:"nodefevent_warn"`

	WarnsHat string `json:"warns_hat"`
	BansHat  string `json:"bans_hat"`
	MytsHat  string `json:"myts_hat"`

	BansForm string `json:"bans_form"`
	MytsForm string `json:"myts_form"`
}

func (lang *Lang) unmarshalJSON(jsoncont []byte) error {
	prelang := Lang{
		//TODO: Create default fields for Lang
	}

	if err := json.Unmarshal(jsoncont, &prelang); err != nil {
		return ErrInvalidConfContent
	}

	*lang = prelang
	return nil
}

func (fb *FeedBack) unmarshalJSON(jsoncont []byte) error {
	prefb := FeedBack{
		//TODO: Create default fields for FeedBack
	}

	if err := json.Unmarshal(jsoncont, &prefb); err != nil {
		return err
	}

	*fb = prefb
	return nil
}

//	Configure new lang configs from received dir(located in Code\language\langconfs\)
//	Returns errors: ErrNoConfDir, ErrNoConfFiles, FeedBack.unmarshalJSON() and Lang.unmarshalJSON() errors
func ConfigureLang(cfname string) (*Lang, *FeedBack, error) {
	currentpath := fmt.Sprintf("Code\\language\\langconfs\\%s", cfname)
	lang := new(Lang)
	fb := new(FeedBack)

	err := func() error {
		_, err := os.ReadDir(currentpath)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return ErrNoConfDir
			}
			return err
		}

		langconf, err := os.ReadFile(currentpath + "\\lang.json")
		fbconf, err := os.ReadFile(currentpath + "\\feedback.json")

		if err != nil {
			return ErrNoConfFiles
		}

		err = lang.unmarshalJSON(langconf)
		if err != nil {
			return err
		}
		err = fb.unmarshalJSON(fbconf)
		if err != nil {
			return err
		}

		return nil
	}()

	return lang, fb, err
}
