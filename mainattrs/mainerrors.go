package mainattrs

import (
	"PIE_BOT/Code/commands/dbcom"
	"fmt"
)

type AccessErr struct {
	CalledPermID int
	Text         string
	DB           *dbcom.ChatDB
}

func NewAccessErr(permID int, text string, db *dbcom.ChatDB) *AccessErr {
	return &AccessErr{CalledPermID: permID, Text: text, DB: db}
}

func (ae *AccessErr) Error() string {
	needrole, err := ae.NeedRole(ae.DB)
	if err != nil {
		return fmt.Sprintf(err.Error())
	}
	return fmt.Sprintf(ae.Text, needrole)
}

func (ae *AccessErr) NeedRole(db *dbcom.ChatDB) (string, error) {
	var rolename string

	minrole, err := db.GetLowestRoleID(ae.CalledPermID)
	if err != nil {
		return rolename, err
	}

	rolename, err = db.GetRoleName(minrole)
	if err != nil {
		return rolename, err
	}

	return rolename, nil

}
