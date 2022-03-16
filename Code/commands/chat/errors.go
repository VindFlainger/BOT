package chat

import "errors"

var (
	BADSHORTNAMEFORMATERROR = errors.New("Error during parsing user shortname")
	SENDINGMESSAGEERROR     = errors.New("Error during sending message")
	IMAGEFILEERROR          = errors.New("Error during reading photo file")
	NOSAVEIMAGEERROR        = errors.New("There are no such saves images")
	INCORRECTKEYERROR       = errors.New("Check function call keys")
	ARGSERROR               = errors.New("Incorrect args inputed")
	TIMEFORMATERROR         = errors.New("Incorrect time format")
	NOACCESS                = errors.New("Bot don't have perms to do this")
)
