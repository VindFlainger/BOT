package dbcom

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type ChatRemover interface {
	RemoveFromChat(chatID, userID int) error
}

//	Writes to the DB myt for received userID with current chatID
//	~If count in DB = MaxWarns >>> call Ban() with reason MaxWarnCount and Time = 86400
//	Returns: (current warn count, error)
//	Returns errors: dbcom.DataBaseErr
func (chatdb *ChatDB) Warn(chatID, userID, maxWarns int, ban ChatRemover) (int, error) {
	var count int

	err := chatdb.DB.QueryRow(fmt.Sprintf(warnreqIns, chatID, userID, chatID, userID)).Scan(&count)
	if err == nil {
		return count, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		err = chatdb.DB.QueryRow(fmt.Sprintf(warnreqUpd1, chatID, userID, maxWarns-1)).Scan(&count)

		switch err {
		case nil:
			return count, nil
		case sql.ErrNoRows:
			_, err = chatdb.DB.Exec(fmt.Sprintf(warnreqDel1, chatID, userID))
			if err != nil {
				return 0, DataBaseErr
			}

			err = chatdb.Ban(chatID, userID, 86400, "MaxWarnCount", ban)
			if err != nil {
				return 0, err
			}
			return maxWarns, nil

		default:
			return 0, DataBaseErr

		}

	}
	return 0, DataBaseErr

}

//	Removes from the DB myt of received userID with current chatID
//	~If it doesn't exist >>> returns NoMytErr
//	Returns: (current warns count, error)
//	Returns errors: dbcom.DataBaseErr, NoMytErr
func (chatdb *ChatDB) Unwarn(chatID, userID int) (int, error) {
	var count int64
	err := chatdb.DB.QueryRow(fmt.Sprintf(warnreqUpd2, chatID, userID)).Scan(&count)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, NoMytErr
	}

	if count == 0 {
		_, err := chatdb.DB.Exec(fmt.Sprintf(warnreqDel2, chatID, userID))

		if err != nil {
			return 0, DataBaseErr
		}
	}

	return int(count), nil

}

//	Writes to the DB ban for received userID with current chatID
//	~If it exist >>> rewrite without returning error
//	Returns errors: dbcom.DataBaseErr, RemoveFromChat() errors
func (chatdb *ChatDB) Ban(chatID, userID, ptime int, reason string, ban ChatRemover) error {
	var unbantime = int(time.Now().Unix()) + ptime

	if err := ban.RemoveFromChat(userID, chatID); err != nil {
		return err

	}

	_, err := chatdb.DB.Exec(banreqIns, chatID, unbantime, reason, userID)

	if err != nil {
		return DataBaseErr
	}

	_, err = chatdb.DB.Exec(banreqUpd, unbantime, reason, chatID, userID)
	if err != nil {
		return DataBaseErr
	}

	return nil
}

//	Removes from the DB ban of received userID with current chatID
//	~If it doesn't exist >>> returns NoBanErr
//	Returns errors:	DataBaseErr, NoBanErr
func (chatdb *ChatDB) UnBan(chatID, userID int) error {
	var unbantime int

	err := chatdb.DB.QueryRow(fmt.Sprintf(banreqDel, chatID, userID)).Scan(&unbantime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NoBanErr
		}
		return DataBaseErr
	}

	return nil
}

//	Writes to the DB myt for received userID with current chatID
//	~If it exist >>> rewrite without returning error
//	Returns errors: dbcom.DataBaseErr
func (chatdb *ChatDB) Myt(chatID, userID, Time int, reason string) error {
	var unmyttime = int(time.Now().Unix()) + Time

	_, err := chatdb.DB.Exec(mytreqIns, chatID, unmyttime, reason, userID)

	if err != nil {
		return DataBaseErr
	}

	_, err = chatdb.DB.Exec(fmt.Sprintf(mytreqUpd, unmyttime, reason, chatID, userID))
	if err != nil {
		return DataBaseErr
	}

	return nil

}

//	Removes from the DB myt of received userID with current chatID
//	~If it doesn't exist >>> returns NoMytErr
//	Returns errors:	DataBaseErr, NoMytErr
func (chatdb *ChatDB) UnMyt(chatID, userID int) error {
	var unmyttime int

	err := chatdb.DB.QueryRow(fmt.Sprintf(mytreqDel, chatID, userID)).Scan(&unmyttime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NoMytErr
		}
		return DataBaseErr
	}

	return nil
}
