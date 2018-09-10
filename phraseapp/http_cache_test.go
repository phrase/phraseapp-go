package phraseapp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocaleDownloadCaching(t *testing.T) {
	var cached = false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		if url == "/v2/projects/1/locales/1/download" {
			if cached {
				etag := r.Header.Get("If-None-Match")
				if etag != "123" {
					t.Errorf("etag should be '123' but is: '%s'", etag)
				}
			}
			cached = true
		}
		w.Header().Set("Etag", "123")
		io.WriteString(w, "OK")
		return
	}))
	defer server.Close()

	client, _ := NewClient(Credentials{Host: server.URL}, false)
	// TODO don't use user cache dir for cache
	client.EnableCaching()
	client.LocaleDownload("1", "1", &LocaleDownloadParams{})
	client.LocaleDownload("1", "1", &LocaleDownloadParams{})
}
