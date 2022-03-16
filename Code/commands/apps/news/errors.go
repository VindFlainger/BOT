package news

import "errors"

var (
	NoContentErr   = errors.New("EMPTY CONTENT BLOCK")
	NoSuchAttrsErr = errors.New("THERE ARE NO SUCH ATTRS")
	RequestErr     = errors.New("NO SUCH A URL ADDRESS")
)
