package utils

import (
	"io"
	"net/http"
)

func DownloadFile(url string) ([]byte, error) {
	var content []byte

	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != 200 {
		return content, DownloadErr
	}

	if content, err = io.ReadAll(resp.Body); err != nil {
		return content, EncodeErr
	}

	defer resp.Body.Close()

	return content, nil
}
