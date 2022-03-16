package dbcom

func (chatdb *ChatDB) MigrateSettings(fromChatID, toChatID int) error {
	cp, err := chatdb.InputSettings(fromChatID)
	if err != nil {
		return err
	}

	_, err = chatdb.DB.Exec(migratereqUpd, cp.MaxWarn, cp.Lang, toChatID)
	if err != nil {
		return DataBaseErr
	}
	return nil
}

func (chatdb *ChatDB) MigrateRoles(fromChatID, toChatID int) error {
	roles, err := chatdb.GetRoles(fromChatID)
	if err != nil {
		return err
	}

	for _, role := range roles {
		chatdb.WriteRole(toChatID, role.UserID, role.RoleID)
	}

	return nil
}

//func (chatdb *ChatDB) MigratePunish(fromChatID, toChatID int) error{
//	bans, err := chatdb.GetBans(fromChatID)
//}
