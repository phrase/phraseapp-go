package phraseapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func sendRequestPaginated(method, rawurl, ctype string, r io.Reader, status, page, perPage int) (io.ReadCloser, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Add("page", strconv.Itoa(page))
	query.Add("per_page", strconv.Itoa(perPage))

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}

	if ctype != "" {
		req.Header.Add("Content-Type", ctype)
	}

	resp, err := send(req, status)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func sendRequest(method, url, ctype string, r io.Reader, status int) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	if ctype != "" {
		req.Header.Add("Content-Type", ctype)
	}

	resp, err := send(req, status)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func send(req *http.Request, status int) (*http.Response, error) {
	err := authenticate(req)
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

func handleResponseStatus(resp *http.Response, expectedStatus int) error {
	switch resp.StatusCode {
	case expectedStatus:
		return nil
	case 400:
		e := new(ErrorResponse)
		err := json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return err
		}
		return e
	case 404:
		return fmt.Errorf("not found")
	case 422:
		e := new(ValidationErrorResponse)
		err := json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return err
		}
		return e
	case 429:
		e, err := NewRateLimitError(resp)
		if err != nil {
			return err
		}
		return e
	default:
		return fmt.Errorf("unexpected status code (%d) received; expected %d", resp.StatusCode, expectedStatus)
	}
}

type ErrorResponse struct {
	Message string
}

func (err *ErrorResponse) Error() string {
	return err.Message
}

type ValidationErrorResponse struct {
	ErrorResponse

	Errors []ValidationErrorMessage
}

func (err *ValidationErrorResponse) Error() string {
	msgs := make([]string, len(err.Errors))
	for i := range err.Errors {
		msgs[i] = err.Errors[i].String()
	}
	return fmt.Sprintf("%s\n%s", err.Message, strings.Join(msgs, "\n"))
}

type ValidationErrorMessage struct {
	Resource string
	Field    string
	Message  string
}

func (msg *ValidationErrorMessage) String() string {
	return fmt.Sprintf("\t[%s:%s] %s", msg.Resource, msg.Field, msg.Message)
}

type RateLimitingError struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

func NewRateLimitError(resp *http.Response) (*RateLimitingError, error) {
	var err error
	re := new(RateLimitingError)

	limit := resp.Header.Get("X-Rate-Limit-Limit")
	re.Limit, err = strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	remaining := resp.Header.Get("X-Rate-Limit-Remaining")
	re.Remaining, err = strconv.Atoi(remaining)
	if err != nil {
		return nil, err
	}

	reset := resp.Header.Get("X-Rate-Limit-Reset")
	sinceEpoch, err := strconv.ParseInt(reset, 10, 64)
	if err != nil {
		return nil, err
	}
	re.Reset = time.Unix(sinceEpoch, 0)

	return re, nil
}

func (rle *RateLimitingError) Error() string {
	return fmt.Sprintf("Rate limit exceeded: from %d requests %d are remaning (reset in %d seconds)", rle.Limit, rle.Remaining, int64(rle.Reset.Sub(time.Now()).Seconds()))
}
