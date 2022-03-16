package news

import (
	"encoding/json"
	"github.com/anaskhan96/soup"
	"strings"
)

type Params struct {
	URL             string   `json:"url"`
	Content         []string `json:"content"`
	Text            []string `json:"text"`
	Date            []string `json:"update_time"`
	PhotoURL        []string `json:"photo_url"`
	MoreInfo        []string `json:"more_info"`
	LastParsedTexts []string `json:"last_parsed_texts"`
	Error           bool     `json:"errors"`
	FileServer      string   `json:"file_server"`
}

type ParsedPage struct {
	Text     string
	PhotoURL string
	Date     string
}

func ParseJson(data []byte) (*Params, error) {
	params := &Params{}
	if err := json.Unmarshal(data, params); err != nil {
		return params, err
	}
	params.LastParsedTexts = make([]string, 15)
	return params, nil
}

func (params *Params) MarshParams() ([]byte, error) {
	marshed, err := json.Marshal(params)
	if err != nil {
		return marshed, err
	}
	return marshed, nil
}

//	/attrs src
func (params *Params) ParsePage() ([]*ParsedPage, error) {
	pp := []*ParsedPage{}
	resp, err := soup.Get(params.URL)
	if err != nil {
		return pp, RequestErr
	}

	root := soup.HTMLParse(resp)

	rootall := root.FindAll(params.Content...)
	if len(rootall) == 0 {
		return pp, NoContentErr
	}

	for _, root := range rootall {
		cp := &ParsedPage{}
		textBlock := root.Find(params.Text...)
		if textBlock.Error != nil {
			continue
		}
		if cp.Text = textBlock.Text(); cp.Text == "" {
			continue
		}

		if !func() bool {
			for _, lpt := range params.LastParsedTexts {
				if lpt == cp.Text {
					return false
				}

			}
			return true
		}() {
			continue
		}

		if len(params.PhotoURL) != 0 {
			attrs := func() int {
				for i, args := range params.PhotoURL {
					if strings.HasPrefix(args, "/attrs:") {
						return i
					}
				}
				return -1
			}()

			var photoBlock soup.Root
			if attrs > 0 {
				photoBlock = root.Find(params.PhotoURL[:attrs]...)
			} else {
				photoBlock = root.Find(params.PhotoURL...)
			}

			if photoBlock.Error == nil {
				url, ex := photoBlock.Attrs()[params.PhotoURL[attrs][strings.Index(params.PhotoURL[attrs], ":")+1:]]
				if !ex {
					return pp, NoSuchAttrsErr
				}

				cp.PhotoURL = func() (str string) {
					if strings.HasPrefix(url, "/") {
						return params.FileServer + url
					}
					if strings.Index(url, "url") > -1 {
						fi := strings.IndexAny(url, "(")
						si := strings.Index(url, ")")
						return url[fi+2 : si-1]
					}
					return ""
				}()
			}
		}

		if len(params.Date) != 0 {
			dateBlock := root.Find(params.Date...)
			if dateBlock.Error == nil {
				cp.Date = dateBlock.Text()
			}
		}

		params.LastParsedTexts = append([]string{cp.Text}, params.LastParsedTexts[:len(params.LastParsedTexts)-1]...)
		pp = append(pp, cp)
	}
	if len(pp) > 1 {
		return pp[:1], err
	}
	return pp, nil
}
