package phraseapp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReposHandler(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(index))
	defer s.Close()

	s.Po
	http.Get(url)
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
	return
}
