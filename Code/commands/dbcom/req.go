package dbcom

const (
	warnreqUpd1 = "UPDATE warns " +
		" SET count = count +  1 " +
		" WHERE chat_id = %d AND user_id = %d AND count < %d" +
		" RETURNING count "

	warnreqUpd2 = " UPDATE warns " +
		" SET count = count-1 " +
		" WHERE chat_id = %d AND user_id = %d " +
		" RETURNING count"

	warnreqDel1 = "DELETE FROM warns " +
		" WHERE chat_id = %d AND user_id = %d"

	warnreqDel2 = " DELETE FROM warns " +
		" WHERE chat_id = %d AND user_id = %d"

	warnreqIns = "INSERT INTO warns(chat_id, user_id, count)" +
		" SELECT %d, %d, 1 " +
		" WHERE NOT EXISTS  " +
		" (SELECT user_id FROM warns " +
		" WHERE chat_id = %d AND user_id = %d" +
		" RETURNING count"

	banreqUpd = "UPDATE bans " +
		" SET unbantime = $1, reason = $2 " +
		" WHERE chat_id = $3 AND user_id = $4"

	banreqIns = "INSERT INTO bans" +
		" SELECT $1, $2, $3, $4" +
		" WHERE NOT EXISTS(SELECT user_id FROM bans WHERE chat_id = $1 AND user_id = $4)"

	banreqDel = "DELETE FROM bans " +
		" WHERE chat_id = %d AND user_id = %d " +
		" RETURNING unbantime"

	mytreqUpd = "UPDATE myts " +
		" SET unmyttime = %d, reason = '%s' " +
		" WHERE chat_id = %d AND user_id = %d"

	mytreqIns = "INSERT INTO myts" +
		" SELECT $1, $2, $3, $4" +
		" WHERE NOT EXISTS(SELECT user_id FROM myts WHERE chat_id = $1 AND user_id = $4) "

	mytreqDel = "DELETE FROM myts " +
		" WHERE chat_id = %d AND user_id = %d " +
		" RETURNING unmyttime"

	inforeqWarns = "SELECT user_id, count FROM warns " +
		" WHERE chat_id = $1"

	inforeqBans = "SELECT unbantime, reason, user_id FROM bans " +
		" WHERE chat_id = $1"

	inforeqMyts = "SELECT unmyttime, reason, user_id FROM myts " +
		" WHERE chat_id = $1"

	permreqRole = "SELECT role_id FROM chat_roles where chat_id = $1 AND user_id = $2"

	permreqDel = "DELETE FROM chat_roles WHERE chat_id = $1 AND user_id = $2"

	permreqIns = "INSERT  INTO chat_roles VALUES ($1, $2, $3)"

	permreqUpd = "UPDATE chat_roles SET role_id = $1 WHERE chat_id = $2 AND user_id = $3"

	permreqSel1 = "SELECT " +
		" EXISTS (SELECT role_id" +
		" FROM roles" +
		" WHERE role_id = $1" +
		" AND array_position(perms, $2::smallint) != 0)"

	permreqSel2 = "SELECT EXISTS(SELECT owner_id FROM chats WHERE chat_id = $1 AND owner_id = $2)"

	permreqSel3 = "SELECT name FROM roles WHERE role_id = $1"

	permreqSel4 = "SELECT role_id FROM roles" +
		" WHERE array_position(perms::integer[], $1) != 0" +
		" ORDER BY role_id LIMIT 1"

	permreqSel5 = "SELECT user_id, role_ID FROM chat_roles" +
		" WHERE chat_id = $1"

	static_sittingsreqUpd = "UPDATE static_settings" +
		" SET maxwarncount = $1, language = $2 " +
		" WHERE chat_id = $3"

	eventreqIns = "INSERT INTO dynamic_settings " +
		" SELECT $1, $2, $3::jsonb " +
		" WHERE NOT exists(SELECT FROM dynamic_settings WHERE chat_id = $1 AND event_id = $2)"

	eventreqUpd = "UPDATE dynamic_settings " +
		" SET params =$3::jsonb " +
		" WHERE chat_id =$1 AND event_id =$2 "

	eventreqDel = "DELETE FROM dynamic_settings " +
		" WHERE chat_id = $1 AND event_id = $2 "

	migratereqUpd = "UPDATE static_settings " +
		" SET maxwarncount = $1, language = $2 " +
		" WHERE chat_id = $3"
)
