package dbcom

import (
	"database/sql"
	"errors"
)

const (
	P_NONE = iota
	P_MEMBERS
	P_WARN
	P_BAN
	P_MYT
	P_ROLE
)

const (
	R_USER = iota
	R_HELPER
	R_MODER
	R_STMODER
	R_ADMIN
)

type UserRole struct {
	RoleName string
	RoleID   int
	UserID   int
}

//	Reads RoleID for received ChatID and UserID from the DB
//	Returns: RoleID for UserID in ChatID, when RoleID = 0 - UserID with such ChatID not in chat_roles
//	Returns errors: DataBaseErr
func (chatdb *ChatDB) GetRole(chatID, userID int) (int, error) {
	var roleID int
	row := chatdb.DB.QueryRow(permreqRole, chatID, userID)
	err := row.Scan(&roleID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, DataBaseErr
	}

	return roleID, err
}

func (chatdb *ChatDB) GetRoles(chatID int) ([]*UserRole, error) {
	userRoles := []*UserRole{}
	rows, err := chatdb.DB.Query(permreqSel5, chatID)
	if err != nil {
		return userRoles, DataBaseErr
	}

	for rows.Next() {
		userRole := &UserRole{}
		err = rows.Scan(&userRole.UserID, &userRole.RoleID)
		if err == nil {
			userRoles = append(userRoles, userRole)
		}
	}
	return userRoles, nil
}

//	Write/Rewrite(TableRoleID != RoleID)/Delete(RoleID = 0) role from the DB by RoleID
//	Returns errors:
//	NoRoleErr when UserID with such ChatID not in chat_roles and you are trying to assign him to a R_USER
//	VALUERROR when RoleID is not in (R_USER, R_HELPER, R_MODER, R_MODER, R_STMODER, R_ADMIN)
//	DataBaseErr
func (chatdb *ChatDB) WriteRole(chatID, userID, roleID int) error {
	switch roleID {
	case R_USER:
		res, err := chatdb.DB.Exec(permreqDel, chatID, userID)
		if err != nil {
			return DataBaseErr
		}
		ra, _ := res.RowsAffected()
		if ra == 0 {
			return NoRoleErr
		}
		return nil

	case R_HELPER, R_MODER, R_STMODER, R_ADMIN:
		roleid, err := chatdb.GetRole(chatID, userID)
		if err != nil {
			return err
		}
		if roleid == 0 {
			_, err := chatdb.DB.Exec(permreqIns, chatID, userID, roleID)

			if err != nil {
				return DataBaseErr
			}
			return nil
		}

		if roleid != roleID {
			_, err := chatdb.DB.Exec(permreqUpd, roleID, chatID, userID)

			if err != nil {
				return DataBaseErr
			}
		}
		return nil
	}
	return RoleValueErr
}

func (chatdb *ChatDB) GetRoleName(roleID int) (string, error) {
	var rolename string
	if err := chatdb.DB.QueryRow(permreqSel3, roleID).Scan(&rolename); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rolename, FindRoleErr
		}
		return rolename, DataBaseErr
	}
	return rolename, nil
}

func (chatdb *ChatDB) GetLowestRoleID(permID int) (int, error) {
	var roleid int

	err := chatdb.DB.QueryRow(permreqSel4, permID).Scan(&roleid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return roleid, FindRoleErr
		}
		return roleid, DataBaseErr
	}
	return roleid, nil
}

//	Reads roles from the DB and compare CalledPermID with RoleID
//	Returns: ( RoleID access for received CalledPermID, error )
//	Return errors: DataBaseErr
func (chatdb *ChatDB) CheckRolePermissions(roleID, permID int) (bool, error) {
	var resp bool
	row := chatdb.DB.QueryRow(permreqSel1, roleID, permID)
	err := row.Scan(&resp)
	if err != nil {
		return resp, DataBaseErr
	}
	return resp, nil
}

//	Reads ownerID for received ChatID from the DB and compare UserID with OwnerID
//	Returns: ( IsOwner?, error )
//	Return errors: DataBaseErr
func (chatdb *ChatDB) CheckOwnerPermissions(chatID, userID int) (bool, error) {
	var resp bool
	row := chatdb.DB.QueryRow(permreqSel2, chatID, userID)
	err := row.Scan(&resp)
	if err != nil {
		return resp, DataBaseErr
	}
	return resp, nil

}

//	Wrapper over CheckOwnerPermissions() and CheckRolePermissions()
//	Reads roles and ownerID for received ChatID from the DB and returns all UserID permissions
//	Returns: ( User Access for this CalledPermID, error)
//	Return errors: CheckOwnerPermissions() and CheckRolePermissions() errors
func (chatdb *ChatDB) CheckAllPermissions(chatID, userID, permID int) (bool, error) {
	RoleID, err := chatdb.GetRole(chatID, userID)
	if err != nil {
		return false, err
	}
	access, err := chatdb.CheckRolePermissions(RoleID, permID)
	if err != nil {
		return false, err
	}

	isown, err := chatdb.CheckOwnerPermissions(chatID, userID)
	if err != nil {
		return false, err
	}

	return access || isown, nil
}
