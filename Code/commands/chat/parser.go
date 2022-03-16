package chat

import (
	"log"
	"strconv"
	"strings"
)

var timeshorts = map[string]int{"s": 1, "m": 60, "h": 3600, "d": 86400}

const (
	T_ignore  = -1
	T_string  = 0
	T_int     = 1
	T_float64 = 2
	T_UserID  = 3
	T_Time    = 4
)

//	Special string parser for vk messages
//	Checks and returns parsed args with types parsedescription
//	ST TYPES: T_ignore, T_string, T_int, T_float64
//	SPEC TYPES:
//	T_UserID -- for porsing domains like [id3211231|@nigwol]
//	T_Time -- for parsing time in format 1x where x represented in timeshorts
//	Returns errors: ARGSERROR, TIMEFORMATERROR,
//	Panics when its get unknown type
func ParseArgs(str []string, types ...int) (outslice []interface{}, error error) {
	if len(str) != len(types) {
		error = ARGSERROR
		return
	}
	for i, posarg := range str {
		switch types[i] {

		case T_int:
			val, err := strconv.ParseInt(posarg, 10, 64)
			if err != nil {
				error = ARGSERROR
				return
			}
			outslice = append(outslice, int(val))
		case T_float64:
			val, err := strconv.ParseFloat(posarg, 64)
			if err != nil {
				error = ARGSERROR
				return
			}
			outslice = append(outslice, val)

		case T_string:
			outslice = append(outslice, posarg)

		case T_UserID:
			UserID, err := ParseShortName(posarg)
			if err != nil {
				error = ARGSERROR
				return
			}
			outslice = append(outslice, UserID)
		case T_Time:
			var miss bool
			for ab, timemul := range timeshorts {
				if strings.HasSuffix(posarg, ab) {
					strtime := strings.TrimRight(posarg, ab)
					parsetime, err := strconv.ParseFloat(strtime, 64)
					if err != nil {
						error = TIMEFORMATERROR
						return
					}
					outslice = append(outslice, int(parsetime)*timemul)
					miss = true
					break
				}
			}
			if !miss {
				error = TIMEFORMATERROR
				return
			}
		case T_ignore:
			{
			}
		default:
			log.Fatal("Unknown type")
		}
	}
	return outslice, nil
}
