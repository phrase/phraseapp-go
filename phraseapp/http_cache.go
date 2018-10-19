package phraseapp

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/peterbourgon/diskv"
)

type httpCacheClient struct {
	cache        *diskv.Diskv
	debug        bool
	cacheSizeMax int64
}

type cacheRecord struct {
	URL      string
	ETag     string
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
	Uncompressed     bool
	Trailer          http.Header
}

// CacheConfig contains the configuration for caching api requests on disk
type CacheConfig struct {
	CacheDir     string
	CacheSizeMax int64 // size in bytes
}

func newHTTPCacheClient(debug bool, config CacheConfig) (*httpCacheClient, error) {
	if config.CacheDir == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}
		config.CacheDir = cacheDir
	}

	if config.CacheSizeMax <= 0 {
		config.CacheSizeMax = 1024 * 1024 * 100 // 100MB
	}

	cachePath := filepath.Join(config.CacheDir, "phraseapp")
	err := os.MkdirAll(cachePath, 0755)
	if err != nil {
		return nil, err
	}

	cache := &httpCacheClient{
		cache: diskv.New(diskv.Options{
			BasePath: cachePath,
		}),
		cacheSizeMax: config.CacheSizeMax,
		debug:        debug,
	}
	return cache, nil
}

func (client *httpCacheClient) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "" && req.Method != "GET" {
		return http.DefaultTransport.RoundTrip(req)
	}

	url := req.URL.String()
	requestParams := requestParams(req)
	cachedResponse, err := client.readCache(url, requestParams)
	if err != nil {
		if err.Error() != "no cache entry" {
			return nil, err
		}
	} else {
		req.Header.Set("If-None-Match", cachedResponse.ETag)
	}

	rsp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusNotModified {
		if client.debug {
			log.Println("found cache and returning cached body")
		}
		cachedResponse.setCachedResponse(rsp)
		return rsp, nil
	}

	err = handleResponseStatus(rsp, 200)
	if err != nil {
		return rsp, err
	}

	cacheSize, err := dirSize(client.cache.BasePath)
	if err != nil {
		return nil, err
	}
	if cacheSize > client.cacheSizeMax {
		client.cache.EraseAll()
	}

	err = client.writeCache(rsp, requestParams, url)
	return rsp, err
}

func requestParams(req *http.Request) string {
	if req.Body != nil {
		body, err := req.GetBody()
		if err != nil {
			return ""
		}
		requestBody, err := ioutil.ReadAll(body)
		if err != nil {
			return ""
		}

		return string(requestBody)
	}

	return ""
}

func (client *httpCacheClient) readCache(url string, requestParams string) (*cacheRecord, error) {
	cache, err := client.cache.Read(md5sum(url + requestParams))
	if err != nil {
		if client.debug {
			log.Println("doing request without etag")
		}
		return nil, fmt.Errorf("no cache entry")
	}

	var cachedResponse *cacheRecord
	var buf bytes.Buffer
	buf.Write(cache)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(&cachedResponse)
	if err != nil {
		return nil, err
	}
	if client.debug {
		log.Printf("found etag %s for request\n", cachedResponse.ETag)
	}

	return cachedResponse, nil
}

func (client *httpCacheClient) writeCache(rsp *http.Response, requestParams string, url string) error {
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	rsp.Body = ioutil.NopCloser(bytes.NewReader(body))
	etag := rsp.Header.Get("Etag")
	cacheKey := md5sum(url + requestParams)
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(cacheRecord{
		URL:     url,
		ETag:    etag,
		Payload: body,
		Response: &httpResponse{
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
	err = client.cache.Write(cacheKey, buf.Bytes())
	return err
}

func (record *cacheRecord) setCachedResponse(rsp *http.Response) {
	rsp.Status = record.Response.Status
	rsp.StatusCode = record.Response.StatusCode
	rsp.Proto = record.Response.Proto
	rsp.ProtoMajor = record.Response.ProtoMajor
	rsp.ProtoMinor = record.Response.ProtoMinor
	rsp.Header = record.Response.Header
	rsp.ContentLength = record.Response.ContentLength
	rsp.TransferEncoding = record.Response.TransferEncoding
	rsp.Trailer = record.Response.Header
	rsp.Body = ioutil.NopCloser(bytes.NewReader(record.Payload))
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func md5sum(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
