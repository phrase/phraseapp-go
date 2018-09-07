package phraseapp

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/peterbourgon/diskv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CacheClient struct {
	Client
	Credentials      Credentials
	CacheDir         string
	CacheSizeMax     uint64
	contentCacheDisk *diskv.Diskv
	etagCacheDisk    *diskv.Diskv
}

func NewCacheClient() (*CacheClient, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	var cacheSizeMax uint64 = 1024 * 1024 * 100 // 100MB
	client := &CacheClient{
		CacheDir:     filepath.Join(cacheDir, "phraseapp"),
		CacheSizeMax: cacheSizeMax,
		contentCacheDisk: diskv.New(diskv.Options{
			BasePath:     cacheDir,
			CacheSizeMax: cacheSizeMax,
		}),
	}
	client.Transport = client
	return client, nil
}

type cacheRecord struct {
	URL      string
	Response *httpResponse
	Payload  []byte
}

// httpResponse is a serializable copy of a http.Response
type httpResponse struct {
	Status           string
	StatusCode       int
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Uncompressed     bool
	Trailer          http.Header
}

func (c *CacheClient) RoundTrip(r *http.Request) (*http.Response, error) {
	url := r.URL.String()
	l := logrus.New().WithField("url", url)
	c.authenticate(r)
	etagValue, err := c.etagCacheDisk.Read(md5sum(url))
	var cachedResponse *cacheRecord
	if err != nil {
		l.Printf("doing request without etag")
	} else {
		etag := string(etagValue)
		l.Printf("using etag %s in request", etag)
		cache, err := c.contentCacheDisk.Read(md5sum(etag + url))
		if err != nil {
			l.Printf("found etag but no cached response")
		} else {
			r.Header.Set("If-None-Match", etag)
			var buf bytes.Buffer
			buf.Write(cache)
			decoder := gob.NewDecoder(&buf)
			err = decoder.Decode(&cachedResponse)
		}
	}

	rsp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	l.Printf("real status=%d", rsp.StatusCode)
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode == http.StatusNotModified {
		l.Printf("found in cache returning cached body")
		rsp.Status = cachedResponse.Response.Status
		rsp.StatusCode = cachedResponse.Response.StatusCode
		rsp.Proto = cachedResponse.Response.Proto
		rsp.ProtoMajor = cachedResponse.Response.ProtoMajor
		rsp.ProtoMinor = cachedResponse.Response.ProtoMinor
		rsp.Header = cachedResponse.Response.Header
		rsp.ContentLength = cachedResponse.Response.ContentLength
		rsp.TransferEncoding = cachedResponse.Response.TransferEncoding
		rsp.Trailer = cachedResponse.Response.Header
		rsp.Body = ioutil.NopCloser(bytes.NewReader(cachedResponse.Payload))
		return rsp, nil
	} else if rsp.Status[0] != '2' {
		return nil, errors.Errorf("got status %s but expected 2x. body=%s", rsp.Status, string(b))
	}

	etag := rsp.Header.Get("Etag")
	etagCacheKey := md5sum(url)
	err = c.etagCacheDisk.Write(etagCacheKey, []byte(etag))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(cacheRecord{Payload: b, Response: &httpResponse{
		Status:           rsp.Status,
		StatusCode:       rsp.StatusCode,
		Proto:            rsp.Proto,
		ProtoMajor:       rsp.ProtoMajor,
		ProtoMinor:       rsp.ProtoMinor,
		Header:           rsp.Header,
		ContentLength:    rsp.ContentLength,
		TransferEncoding: rsp.TransferEncoding,
		Trailer:          rsp.Header,
	}})
	contentCacheKey := md5sum(etag + url)
	err = c.contentCacheDisk.Write(contentCacheKey, buf.Bytes())
	if err != nil {
		return nil, err
	}

	rsp.Body = ioutil.NopCloser(bytes.NewReader(b))
	return rsp, nil
}

func md5sum(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
