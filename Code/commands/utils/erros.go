package utils

import "errors"

var (
	DownloadErr = errors.New("ERROR DURING DOWNLOADING FILE")
	EncodeErr   = errors.New("ERROR DURING READING DOWNLOADED DATA")
)
