package dbcom

import "errors"

var (
	DataBaseErr       = errors.New("Critical error during DB request")
	RoleValueErr      = errors.New("Value error")
	NoRoleErr         = errors.New("This user has no roles")
	NoBanErr          = errors.New("User with this id & chat_id can't be find")
	NoMytErr          = errors.New("User with this id & chat_id can't be find")
	NoEventErr        = errors.New("Event with this id & chat_id can't be find")
	FindRoleErr       = errors.New("There are no roles with received args")
	TimeOverFlow      = errors.New("Time should be in int64 range (UNIX)")
	BadFileContentErr = errors.New("File content is nil or not in json format")
)
