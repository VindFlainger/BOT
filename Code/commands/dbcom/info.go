package dbcom

type BanInfo struct {
	UserID    int
	Unbantime int
	Reason    string
}

type MytInfo struct {
	UserID    int
	Unmyttime int
	Reason    string
}

//	Reads all warns from DB
//	Returns: (map where key - userID and value - warn point, error)
//	Returns errors: DataBaseErr
func (chatdb *ChatDB) GetWarns(ChatID int) (map[int]int, error) {
	var UserID, count int
	warns := make(map[int]int)
	rows, err := chatdb.DB.Query(inforeqWarns, ChatID)

	if err != nil {
		return warns, DataBaseErr
	}

	for rows.Next() {
		if err = rows.Scan(&UserID, &count); err != nil {
			return warns, DataBaseErr
		}
		warns[UserID] = count
	}

	return warns, nil
}

//	Reads all bans from DB
//	Returns: (BanInfo slice, error)
//	Returns errors: DataBaseErr
func (chatdb *ChatDB) GetBans(chatID int) ([]*BanInfo, error) {
	var bans []*BanInfo

	rows, err := chatdb.DB.Query(inforeqBans, chatID)
	if err != nil {
		return bans, DataBaseErr
	}

	for rows.Next() {
		bi := &BanInfo{}
		if err := rows.Scan(&bi.Unbantime, &bi.Reason, &bi.UserID); err != nil {
			return bans, DataBaseErr
		}
		bans = append(bans, bi)
	}

	return bans, nil
}

//	Reads all myts from DB
//	Returns: (MytInfo slice, error)
//	Returns errors: DataBaseErr
func (chatdb *ChatDB) GetMyts(chatID int) ([]*MytInfo, error) {
	var myts []*MytInfo

	rows, err := chatdb.DB.Query(inforeqMyts, chatID)
	if err != nil {
		return myts, DataBaseErr
	}

	for rows.Next() {
		bi := &MytInfo{}
		if err := rows.Scan(&bi.Unmyttime, &bi.Reason, &bi.UserID); err != nil {
			return myts, DataBaseErr
		}
		myts = append(myts, bi)
	}

	return myts, nil
}
