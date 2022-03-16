package dbcom

import (
	"database/sql"
	"fmt"
)

const (
	LANG_RU = 1
	LANG_EN = 2
	LANG_BU = 3
	ST_DEF  = 1
	ST_VIP  = 2
	ST_PREM = 3
	ST_DEL  = 4
	ST_SUP  = 5
)

type ChatDB struct {
	DB    *sql.DB
	Count int // load factor
}

type ChatParams struct {
	Status     int
	MaxWarn    int
	Lang       int
	ChatEvents []int
}

func (chatdb *ChatDB) AddNewChat(ChatID, OwnerID int, params *ChatParams) (bool, error) {
	var isexist bool
	res, err := chatdb.DB.Exec(
		fmt.Sprintf("INSERT INTO chats SELECT %d , %d, %d WHERE NOT EXISTS( SELECT chat_id FROM chats WHERE chat_id = %d)",
			ChatID,
			OwnerID,
			params.Status,
			ChatID))
	if err != nil {
		return false, err
	}
	rowsA, _ := res.RowsAffected()

	res, err = chatdb.DB.Exec(
		fmt.Sprintf(
			"INSERT INTO static_settings SELECT %d , %d, %d WHERE NOT EXISTS( SELECT chat_id FROM static_settings WHERE chat_id = %d)",
			ChatID,
			params.MaxWarn,
			params.Lang,
			ChatID))
	if err != nil {
		return false, err
	}

	rowsB, _ := res.RowsAffected()
	if rowsB == 0 && rowsA == 0 {
		isexist = true
	}

	return isexist, nil
}

func (chatdb *ChatDB) InputSettings(ChatID int) (*ChatParams, error) {
	var maxwncount, lang, status interface{}
	row := chatdb.DB.QueryRow(fmt.Sprintf("SELECT status FROM chats WHERE chat_id = %d", ChatID))
	if row.Scan(&status) != nil {
		return &ChatParams{}, DataBaseErr
	}

	row = chatdb.DB.QueryRow(fmt.Sprintf("SELECT maxwarncount, language FROM static_settings WHERE chat_id = %d", ChatID))
	if row.Scan(&maxwncount, &lang) != nil {
		return &ChatParams{}, DataBaseErr
	}
	return &ChatParams{MaxWarn: int(maxwncount.(int64)), Lang: int(lang.(int64)), Status: int(status.(int64))}, nil
}

func (chatdb *ChatDB) ChangeSettings(ChatID int, params *ChatParams) error {
	_, err := chatdb.DB.Exec(static_sittingsreqUpd,
		params.MaxWarn,
		params.Lang,
		ChatID)
	if err != nil {
		return DataBaseErr
	}
	return nil
}
