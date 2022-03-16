package language

import "errors"

var (
	//	You should create:
	//	... -language
	//		      -langconfs
	//		            -{confdir} <--- this dir called error
	ErrNoConfDir = errors.New("No such config dir")

	//	You should create two config files
	//	... -longconfs
	//	        -{confdir}
	//               -feedback.json <--- this file called error
	//               -lang.json     <--- this file called error
	ErrNoConfFiles = errors.New("No required config files in received config dir")

	ErrInvalidConfContent = errors.New("Error during reading json content from file")
)
