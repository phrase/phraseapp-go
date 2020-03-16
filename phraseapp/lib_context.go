package phraseapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func (client *Client) buildRequestWithContext(ctx context.Context, method string, u *url.URL, body io.Reader, contentType string) (*http.Request, error) {
	if client.debug {
		fmt.Fprintln(os.Stderr, "Method:", method)
		fmt.Fprintln(os.Stderr, "URL:", u)

		if body != nil {
			bodyBytes, err := ioutil.ReadAll(body)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading body:", err.Error())
			}

			fmt.Fprintln(os.Stderr, "Body:", string(bodyBytes))
			body = bytes.NewReader(bodyBytes)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	return req, nil
}

func (client *Client) sendRequestWithContext(ctx context.Context, method, urlPath, contentType string, body io.Reader, expectedStatus int) (io.ReadCloser, error) {
	endpointURL, err := url.Parse(client.Credentials.Host + urlPath)
	if err != nil {
		return nil, err
	}

	req, err := client.buildRequestWithContext(ctx, method, endpointURL, body, contentType)
	if err != nil {
		return nil, err
	}

	resp, err := client.send(req, expectedStatus)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}


// Create a translation.
func (client *Client) TranslationCreateWithContext(ctx context.Context, project_id string, params *TranslationParams) (*TranslationDetails, error) {
	retVal := new(TranslationDetails)
	err := func() error {

		url := fmt.Sprintf("/v2/projects/%s/translations", url.QueryEscape(project_id))

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := client.sendRequestWithContext(ctx, "POST", url, "application/json", paramsBuf, 201)

		if err != nil {
			return err
		}
		defer rc.Close()

		var reader io.Reader
		if client.debug {
			reader = io.TeeReader(rc, os.Stderr)
		} else {
			reader = rc
		}

		return json.NewDecoder(reader).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing translation.
func (client *Client) TranslationUpdateWithContext(ctx context.Context, project_id, id string, params *TranslationUpdateParams) (*TranslationDetails, error) {
	retVal := new(TranslationDetails)
	err := func() error {

		url := fmt.Sprintf("/v2/projects/%s/translations/%s", url.QueryEscape(project_id), url.QueryEscape(id))

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := client.sendRequestWithContext(ctx, "PATCH", url, "application/json", paramsBuf, 200)

		if err != nil {
			return err
		}
		defer rc.Close()

		var reader io.Reader
		if client.debug {
			reader = io.TeeReader(rc, os.Stderr)
		} else {
			reader = rc
		}

		return json.NewDecoder(reader).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new key.
func (client *Client) KeyCreateWithContext(ctx context.Context, project_id string, params *TranslationKeyParams) (*TranslationKeyDetails, error) {
	retVal := new(TranslationKeyDetails)
	err := func() error {

		url := fmt.Sprintf("/v2/projects/%s/keys", url.QueryEscape(project_id))

		paramsBuf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(paramsBuf)
		ctype := writer.FormDataContentType()

		if params.Branch != nil {
			err := writer.WriteField("branch", *params.Branch)
			if err != nil {
				return err
			}
		}

		if params.DataType != nil {
			err := writer.WriteField("data_type", *params.DataType)
			if err != nil {
				return err
			}
		}

		if params.Description != nil {
			err := writer.WriteField("description", *params.Description)
			if err != nil {
				return err
			}
		}

		if params.LocalizedFormatKey != nil {
			err := writer.WriteField("localized_format_key", *params.LocalizedFormatKey)
			if err != nil {
				return err
			}
		}

		if params.LocalizedFormatString != nil {
			err := writer.WriteField("localized_format_string", *params.LocalizedFormatString)
			if err != nil {
				return err
			}
		}

		if params.MaxCharactersAllowed != nil {
			err := writer.WriteField("max_characters_allowed", strconv.FormatInt(*params.MaxCharactersAllowed, 10))
			if err != nil {
				return err
			}
		}

		if params.Name != nil {
			err := writer.WriteField("name", *params.Name)
			if err != nil {
				return err
			}
		}

		if params.NamePlural != nil {
			err := writer.WriteField("name_plural", *params.NamePlural)
			if err != nil {
				return err
			}
		}

		if params.OriginalFile != nil {
			err := writer.WriteField("original_file", *params.OriginalFile)
			if err != nil {
				return err
			}
		}

		if params.Plural != nil {
			err := writer.WriteField("plural", strconv.FormatBool(*params.Plural))
			if err != nil {
				return err
			}
		}

		if params.RemoveScreenshot != nil {
			err := writer.WriteField("remove_screenshot", strconv.FormatBool(*params.RemoveScreenshot))
			if err != nil {
				return err
			}
		}

		if params.Screenshot != nil {
			part, err := writer.CreateFormFile("screenshot", filepath.Base(*params.Screenshot))
			if err != nil {
				return err
			}
			file, err := os.Open(*params.Screenshot)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return err
			}
			err = file.Close()
			if err != nil {
				return err
			}
		}

		if params.Tags != nil {
			err := writer.WriteField("tags", *params.Tags)
			if err != nil {
				return err
			}
		}

		if params.Unformatted != nil {
			err := writer.WriteField("unformatted", strconv.FormatBool(*params.Unformatted))
			if err != nil {
				return err
			}
		}

		if params.XmlSpacePreserve != nil {
			err := writer.WriteField("xml_space_preserve", strconv.FormatBool(*params.XmlSpacePreserve))
			if err != nil {
				return err
			}
		}
		err := writer.WriteField("utf8", "✓")
		writer.Close()

		rc, err := client.sendRequestWithContext(ctx, "POST", url, ctype, paramsBuf, 201)

		if err != nil {
			return err
		}
		defer rc.Close()

		var reader io.Reader
		if client.debug {
			reader = io.TeeReader(rc, os.Stderr)
		} else {
			reader = rc
		}

		return json.NewDecoder(reader).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing key.
func (client *Client) KeyUpdateWithContext(ctx context.Context, project_id, id string, params *TranslationKeyParams) (*TranslationKeyDetails, error) {
	retVal := new(TranslationKeyDetails)
	err := func() error {

		url := fmt.Sprintf("/v2/projects/%s/keys/%s", url.QueryEscape(project_id), url.QueryEscape(id))

		paramsBuf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(paramsBuf)
		ctype := writer.FormDataContentType()

		if params.Branch != nil {
			err := writer.WriteField("branch", *params.Branch)
			if err != nil {
				return err
			}
		}

		if params.DataType != nil {
			err := writer.WriteField("data_type", *params.DataType)
			if err != nil {
				return err
			}
		}

		if params.Description != nil {
			err := writer.WriteField("description", *params.Description)
			if err != nil {
				return err
			}
		}

		if params.LocalizedFormatKey != nil {
			err := writer.WriteField("localized_format_key", *params.LocalizedFormatKey)
			if err != nil {
				return err
			}
		}

		if params.LocalizedFormatString != nil {
			err := writer.WriteField("localized_format_string", *params.LocalizedFormatString)
			if err != nil {
				return err
			}
		}

		if params.MaxCharactersAllowed != nil {
			err := writer.WriteField("max_characters_allowed", strconv.FormatInt(*params.MaxCharactersAllowed, 10))
			if err != nil {
				return err
			}
		}

		if params.Name != nil {
			err := writer.WriteField("name", *params.Name)
			if err != nil {
				return err
			}
		}

		if params.NamePlural != nil {
			err := writer.WriteField("name_plural", *params.NamePlural)
			if err != nil {
				return err
			}
		}

		if params.OriginalFile != nil {
			err := writer.WriteField("original_file", *params.OriginalFile)
			if err != nil {
				return err
			}
		}

		if params.Plural != nil {
			err := writer.WriteField("plural", strconv.FormatBool(*params.Plural))
			if err != nil {
				return err
			}
		}

		if params.RemoveScreenshot != nil {
			err := writer.WriteField("remove_screenshot", strconv.FormatBool(*params.RemoveScreenshot))
			if err != nil {
				return err
			}
		}

		if params.Screenshot != nil {
			part, err := writer.CreateFormFile("screenshot", filepath.Base(*params.Screenshot))
			if err != nil {
				return err
			}
			file, err := os.Open(*params.Screenshot)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return err
			}
			err = file.Close()
			if err != nil {
				return err
			}
		}

		if params.Tags != nil {
			err := writer.WriteField("tags", *params.Tags)
			if err != nil {
				return err
			}
		}

		if params.Unformatted != nil {
			err := writer.WriteField("unformatted", strconv.FormatBool(*params.Unformatted))
			if err != nil {
				return err
			}
		}

		if params.XmlSpacePreserve != nil {
			err := writer.WriteField("xml_space_preserve", strconv.FormatBool(*params.XmlSpacePreserve))
			if err != nil {
				return err
			}
		}
		err := writer.WriteField("utf8", "✓")
		writer.Close()

		rc, err := client.sendRequestWithContext(ctx, "PATCH", url, ctype, paramsBuf, 200)

		if err != nil {
			return err
		}
		defer rc.Close()

		var reader io.Reader
		if client.debug {
			reader = io.TeeReader(rc, os.Stderr)
		} else {
			reader = rc
		}

		return json.NewDecoder(reader).Decode(&retVal)

	}()
	return retVal, err
}