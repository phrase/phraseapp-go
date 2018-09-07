package phraseapp

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logrus.New()
	if err := run(l); err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
	return
}

type CacheClient struct {
	Credentials  Credentials
	Client       *http.Client
	CacheDir     string
	contentCache map[string]*cacheRecord
	etagCache    map[string]string
}

func NewCacheClient() CacheClient {
	return CacheClient{}
}

func run(l *logrus.Logger) error {
	cl := http.Client{}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	cl.Transport = &CacheClient{
		CacheDir:     filepath.Join(cacheDir, "phraseapp"),
		etagCache:    map[string]string{},
		contentCache: map[string]*cacheRecord{},
	}

	for _, code := range []string{"en", "de"} {
		for i := 0; i < 2; i++ {
			rsp, err := cl.Get("https://phraseapp.com/api/v2/projects/1d8ae641902624df63ce6fbd64ff9549/locales/" + code + "/download?file_format=yml")
			if err != nil {
				return err
			}
			defer rsp.Body.Close()
			if rsp.Status[0] != '2' {
				b, _ := ioutil.ReadAll(rsp.Body)
				return errors.Errorf("got status %s but expected 2x. body=%s", rsp.Status, string(b))
			}
			rsp.Body = nil
			l.Printf("fake status: %s", rsp.Status)
		}
	}

	return nil
}

type cacheRecord struct {
	URL string
	// TODO: replace with a copy of http.Response with only the primitive types
	Response *http.Response
	Payload  []byte
}

func (c *CacheClient) RoundTrip(r *http.Request) (*http.Response, error) {
	url := r.URL.String()
	l := logrus.New().WithField("url", url)
	// TODO: use auth header to do this
	r.SetBasicAuth(c.Credentials.Token, "")
	var ok bool
	etag, ok := c.etagCache[url]
	var cachedResponse *cacheRecord
	if ok {
		l.Printf("using etag %s in request", etag)
		cachedResponse, ok = c.contentCache[etag+url]
		if ok {
			r.Header.Set("If-None-Match", etag)
		} else {
			l.Printf("found etag but no cached response")
		}
	} else {
		l.Printf("doing request without etag")
	}

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}
	rsp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	l.Printf("real status=%d", rsp.StatusCode)

	b, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode == http.StatusNotModified {
		l.Printf("found in cache returning cached body")
		rsp := cachedResponse.Response
		rsp.Body = ioutil.NopCloser(bytes.NewReader(cachedResponse.Payload))
		return rsp, nil
	} else if rsp.Status[0] != '2' {
		return nil, errors.Errorf("got status %s but expected 2x. body=%s", rsp.Status, string(b))
	}
	etag = rsp.Header.Get("Etag")
	c.etagCache[url] = etag
	c.contentCache[etag+url] = &cacheRecord{Payload: b, Response: rsp}
	rsp.Body = ioutil.NopCloser(bytes.NewReader(b))
	return rsp, nil
}
