package apps

import "errors"

var (
	BadFileContentErr = errors.New("File content is nil or not in json format")
)
