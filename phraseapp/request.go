package phraseapp

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var Debug bool

func EnableDebug() {
	Debug = true
}

func (client *Client) sendRequestPaginated(method, rawurl, ctype string, r io.Reader, status, page, perPage int) (io.ReadCloser, error) {
	endpointUrl := client.Credentials.Host + rawurl
	u, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Add("page", strconv.Itoa(page))
	query.Add("per_page", strconv.Itoa(perPage))

	u.RawQuery = query.Encode()

	if Debug {
		fmt.Fprintln(os.Stderr, method, u.String())
	}

	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}

	if ctype != "" {
		req.Header.Add("Content-Type", ctype)
	}

	resp, err := client.send(req, status)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (client *Client) sendRequest(method, url, ctype string, r io.Reader, status int) (io.ReadCloser, error) {
	endpointUrl := client.Credentials.Host + url
	if Debug {
		fmt.Fprintln(os.Stderr, method, url)
		bytes, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		}
		str := string(bytes)
		fmt.Fprintln(os.Stderr, str)
		r = strings.NewReader(str)
	}
	req, err := http.NewRequest(method, endpointUrl, r)
	if err != nil {
		return nil, err
	}

	if ctype != "" {
		req.Header.Add("Content-Type", ctype)
	}

	resp, err := client.send(req, status)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (client *Client) send(req *http.Request, status int) (*http.Response, error) {
	err := client.authenticate(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = handleResponseStatus(resp, status)
	if err != nil {
		resp.Body.Close()
	}
	return resp, err
}
