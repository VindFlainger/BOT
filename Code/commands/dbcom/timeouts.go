package dbcom

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	TOID_Members = 1

	TO_Members = 60
)

//	Realize DataBase call and returns timeout before reaccessing the method
//	DB TABLE:
//  CREATE TABLE timeouts(
//    chat_id integer  not null
//        constraint timeouts_chat_id_fkey
//            references chats
//            on delete cascade,
//    command smallint not null,
//    timeout integer  not null);
func GetTimeOut(db *sql.DB, ChatID, timeout, TOID_ID int) int {
	var TO interface{}
	nowtime := time.Now().Unix() - 2
	row := db.QueryRow(fmt.Sprintf("SELECT timeout FROM timeouts WHERE chat_id = %d AND command = %d",
		ChatID,
		TOID_ID))
	if err := row.Scan(&TO); err != nil {
		return 0
	}

	if TO == nil {
		db.Exec(fmt.Sprintf("INSERT INTO timeouts VALUES(%d, %d, %d)", ChatID, TOID_ID, timeout+int(nowtime)))
		return 0
	} else {
		if TO.(int64) < nowtime {
			fmt.Println(db.Exec(fmt.Sprintf("UPDATE timeouts SET timeout = %d WHERE chat_id = %d AND command = %d ",
				timeout+int(nowtime),
				ChatID,
				TOID_Members)))
			return 0
		}
	}
	return int(TO.(int64) - nowtime)
}
