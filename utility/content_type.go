package utility

import (
	"io"
	"net/http"
)

func DetectContentType(file io.ReadSeeker) string {

	first512 := make([]byte, 512)
	file.Read(first512)
	file.Seek(0, 0)

	return http.DetectContentType(first512)
}
