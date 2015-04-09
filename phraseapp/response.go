package phraseapp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
