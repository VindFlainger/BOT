package dbcom

import (
	"github.com/lib/pq"
)

//	Pre-installed events
const (
	E_NEWS = iota + 1
)

var EVENTS = []int{E_NEWS}

func (chatdb *ChatDB) AddEvent(chatID, eventID int, params []byte) error {
	res, err := chatdb.DB.Exec(eventreqIns, chatID, eventID, params)

	if err != nil {
		if err.(*pq.Error).Code == "22P02" {
			return BadFileContentErr
		}
		return DataBaseErr
	}

	if af, _ := res.RowsAffected(); af == 0 {
		_, err := chatdb.DB.Exec(eventreqUpd, chatID, eventID, params)
		if err != nil {
			return DataBaseErr
		}
	}
	return nil
}

func (chatdb *ChatDB) DelEvent(chatID, eventID int) error {
	res, err := chatdb.DB.Exec(eventreqDel, chatID, eventID)
	if err != nil {
		return DataBaseErr
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return NoEventErr
	}

	return nil
}
