package phraseapp

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/peterbourgon/diskv"
)

type httpCacheClient struct {
	contentCache *diskv.Diskv
	etagCache    *diskv.Diskv
	debug        bool
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
	Uncompressed     bool
	Trailer          http.Header
}

// CacheConfig contains the configuration for caching api requests on disk
type CacheConfig struct {
	CacheDir     string
	CacheSizeMax uint64
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
		var cacheSizeMax uint64 = 1024 * 1024 * 100 // 100MB
		config.CacheSizeMax = cacheSizeMax
	}

	cachePath := filepath.Join(config.CacheDir, "phraseapp")
	cache := &httpCacheClient{
		contentCache: diskv.New(diskv.Options{
			BasePath:     cachePath,
			CacheSizeMax: config.CacheSizeMax,
		}),
		etagCache: diskv.New(diskv.Options{
			BasePath:     cachePath,
			CacheSizeMax: config.CacheSizeMax,
		}),
		debug: debug,
	}
	return cache, nil
}

func (client *httpCacheClient) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "" && req.Method != "GET" {
		rsp, err := http.DefaultTransport.RoundTrip(req)
		return rsp, err
	}

	url := req.URL.String()
	cachedResponse, err := client.getCache(req, url)
	if err != nil {
		return nil, err
	}

	rsp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	rsp.Body.Close()
	if rsp.StatusCode == http.StatusNotModified {
		if client.debug {
			log.Println("found in cache and returning cached body")
		}
		cachedResponse.setCachedResponse(rsp)
		return rsp, nil
	} else {
		rsp.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	err = handleResponseStatus(rsp, 200)
	if err != nil {
		return rsp, err
	}

	etag := rsp.Header.Get("Etag")
	etagCacheKey := md5sum(url)
	err = client.etagCache.Write(etagCacheKey, []byte(etag))
	if err != nil {
		return nil, err
	}

	contentCacheKey := md5sum(etag + url)
	encodedCache := cachedResponse.encode(rsp, url, body)
	err = client.contentCache.Write(contentCacheKey, encodedCache)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (client *httpCacheClient) getCache(req *http.Request, url string) (*cacheRecord, error) {
	var cachedResponse *cacheRecord
	etagResult, err := client.etagCache.Read(md5sum(url))
	if err != nil {
		if client.debug {
			log.Println("doing request without etag")
		}
	} else {
		etag := string(etagResult)
		if client.debug {
			log.Printf("using etag %s in request\n", etag)
		}
		cache, err := client.contentCache.Read(md5sum(etag + url))
		if err != nil {
			if client.debug {
				log.Println("found no cached response for etag")
			}
		} else {
			req.Header.Set("If-None-Match", etag)
			var buf bytes.Buffer
			buf.Write(cache)
			decoder := gob.NewDecoder(&buf)
			err = decoder.Decode(&cachedResponse)
			if err != nil {
				return nil, err
			}
		}
	}

	return cachedResponse, nil
}

func (record *cacheRecord) encode(rsp *http.Response, url string, body []byte) []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(cacheRecord{URL: url, Payload: body, Response: &httpResponse{
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

	return buf.Bytes()
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

func md5sum(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
