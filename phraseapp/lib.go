package phraseapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Authorization struct {
	CreatedAt      time.Time `json:"created_at"`
	HashedToken    string    `json:"hashed_token"`
	Id             string    `json:"id"`
	Note           string    `json:"note"`
	Scopes         []string  `json:"scopes"`
	TokenLastEight string    `json:"token_last_eight"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AuthorizationWithToken struct {
	Authorization

	Token string `json:"token"`
}

type BlacklistedKey struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	CreatedAt time.Time    `json:"created_at"`
	Id        string       `json:"id"`
	Message   string       `json:"message"`
	UpdatedAt time.Time    `json:"updated_at"`
	User      *UserPreview `json:"user"`
}

type KeyPreview struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Plural bool   `json:"plural"`
}

type Locale struct {
	Code         string         `json:"code"`
	CreatedAt    time.Time      `json:"created_at"`
	Default      bool           `json:"default"`
	Id           string         `json:"id"`
	Main         bool           `json:"main"`
	Name         string         `json:"name"`
	Rtl          bool           `json:"rtl"`
	SourceLocale *LocalePreview `json:"source_locale"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type LocaleDetails struct {
	Locale

	Statistics *LocaleStatistics `json:"statistics"`
}

type LocaleFileImport struct {
	CreatedAt time.Time `json:"created_at"`
	Format    string    `json:"format"`
	Id        string    `json:"id"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LocaleFileImportWithSummary struct {
	LocaleFileImport

	Summary SummaryType `json:"summary"`
}

type LocalePreview struct {
	Code string `json:"code"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LocaleStatistics struct {
	KeysTotalCount              int64 `json:"keys_total_count"`
	KeysUntranslatedCount       int64 `json:"keys_untranslated_count"`
	MissingWordsCount           int64 `json:"missing_words_count"`
	TranslationsCompletedCount  int64 `json:"translations_completed_count"`
	TranslationsUnverifiedCount int64 `json:"translations_unverified_count"`
	UnverifiedWordsCount        int64 `json:"unverified_words_count"`
	WordsTotalCount             int64 `json:"words_total_count"`
}

type Project struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectDetails struct {
	Project

	SharesTranslationMemory bool `json:"shares_translation_memory"`
}

type StatisticsListItem struct {
	Locale     *LocalePreview `json:"locale"`
	Statistics StatisticsType `json:"statistics"`
}

type StatisticsType struct {
	KeysTotalCount              int64 `json:"keys_total_count"`
	KeysUntranslatedCount       int64 `json:"keys_untranslated_count"`
	TranslationsCompletedCount  int64 `json:"translations_completed_count"`
	TranslationsUnverifiedCount int64 `json:"translations_unverified_count"`
}

type Styleguide struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StyleguideDetails struct {
	Styleguide

	Audience           string `json:"audience"`
	Business           string `json:"business"`
	CompanyBranding    string `json:"company_branding"`
	Formatting         string `json:"formatting"`
	GlossaryTerms      string `json:"glossary_terms"`
	GrammarConsistency string `json:"grammar_consistency"`
	GrammaticalPerson  string `json:"grammatical_person"`
	LiteralTranslation string `json:"literal_translation"`
	OverallTone        string `json:"overall_tone"`
	PublicUrl          string `json:"public_url"`
	Samples            string `json:"samples"`
	TargetAudience     string `json:"target_audience"`
	VocabularyType     string `json:"vocabulary_type"`
}

type StyleguidePreview struct {
	Id        string `json:"id"`
	PublicUrl string `json:"public_url"`
}

type SummaryType struct {
	LocalesCreated         int64 `json:"locales_created"`
	TagsCreated            int64 `json:"tags_created"`
	TranslationKeysCreated int64 `json:"translation_keys_created"`
	TranslationsCreated    int64 `json:"translations_created"`
	TranslationsUpdated    int64 `json:"translations_updated"`
}

type Tag struct {
	CreatedAt time.Time `json:"created_at"`
	KeysCount int64     `json:"keys_count"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TagWithStats struct {
	Tag

	Statistics []*StatisticsListItem `json:"statistics"`
}

type Translation struct {
	Content      string         `json:"content"`
	CreatedAt    time.Time      `json:"created_at"`
	Excluded     bool           `json:"excluded"`
	Id           string         `json:"id"`
	Key          *KeyPreview    `json:"key"`
	Locale       *LocalePreview `json:"locale"`
	PluralSuffix string         `json:"plural_suffix"`
	Unverified   bool           `json:"unverified"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type TranslationDetails struct {
	Translation

	User      *UserPreview `json:"user"`
	WordCount int64        `json:"word_count"`
}

type TranslationKey struct {
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	NameHash    string    `json:"name_hash"`
	Plural      bool      `json:"plural"`
	Tags        []string  `json:"tags"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TranslationKeyDetails struct {
	TranslationKey

	CommentsCount        int64  `json:"comments_count"`
	DataType             string `json:"data_type"`
	FormatValueType      string `json:"format_value_type"`
	MaxCharactersAllowed int64  `json:"max_characters_allowed"`
	NamePlural           string `json:"name_plural"`
	OriginalFile         string `json:"original_file"`
	ScreenshotUrl        string `json:"screenshot_url"`
	Unformatted          bool   `json:"unformatted"`
	XmlSpacePreserve     bool   `json:"xml_space_preserve"`
}

type TranslationOrder struct {
	AmountInCents                    int64              `json:"amount_in_cents"`
	CreatedAt                        time.Time          `json:"created_at"`
	Currency                         string             `json:"currency"`
	Id                               string             `json:"id"`
	Lsp                              string             `json:"lsp"`
	Message                          string             `json:"message"`
	Priority                         bool               `json:"priority"`
	ProgressPercent                  int64              `json:"progress_percent"`
	Quality                          bool               `json:"quality"`
	SourceLocale                     *LocalePreview     `json:"source_locale"`
	State                            string             `json:"state"`
	Styleguide                       *StyleguidePreview `json:"styleguide"`
	Tag                              string             `json:"tag"`
	TargetLocales                    []*LocalePreview   `json:"target_locales"`
	TranslationType                  string             `json:"translation_type"`
	UnverifyTranslationsUponDelivery bool               `json:"unverify_translations_upon_delivery"`
	UpdatedAt                        time.Time          `json:"updated_at"`
}

type TranslationVersion struct {
	ChangedAt    time.Time      `json:"changed_at"`
	Content      string         `json:"content"`
	CreatedAt    time.Time      `json:"created_at"`
	Id           string         `json:"id"`
	Key          *KeyPreview    `json:"key"`
	Locale       *LocalePreview `json:"locale"`
	PluralSuffix string         `json:"plural_suffix"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type TranslationVersionWithUser struct {
	TranslationVersion

	User *UserPreview `json:"user"`
}

type User struct {
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

type UserPreview struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type AuthorizationParams struct {
	Note   string   `json:"note"`
	Scopes []string `json:"scopes,omitempty"`
}

type BlacklistedKeyParams struct {
	Name string `json:"name"`
}

type CommentParams struct {
	Message string `json:"message"`
}

type TranslationKeyParams struct {
	DataType             *string  `json:"data_type,omitempty"`
	Description          *string  `json:"description,omitempty"`
	FormatValueType      *string  `json:"format_value_type,omitempty"`
	MaxCharactersAllowed *int64   `json:"max_characters_allowed,omitempty"`
	Name                 string   `json:"name"`
	NamePlural           *string  `json:"name_plural,omitempty"`
	OriginalFile         *string  `json:"original_file,omitempty"`
	Plural               *bool    `json:"plural,omitempty"`
	RemoveScreenshot     *bool    `json:"remove_screenshot,omitempty"`
	Screenshot           *string  `json:"screenshot,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	Unformatted          *bool    `json:"unformatted,omitempty"`
	XmlSpacePreserve     *bool    `json:"xml_space_preserve,omitempty"`
}

type LocaleParams struct {
	Code           string  `json:"code"`
	Default        *bool   `json:"default,omitempty"`
	Main           *bool   `json:"main,omitempty"`
	Name           string  `json:"name"`
	Rtl            *bool   `json:"rtl,omitempty"`
	SourceLocaleId *string `json:"source_locale_id,omitempty"`
}

type TranslationOrderParams struct {
	Category                         string   `json:"category"`
	IncludeUntranslatedKeys          *bool    `json:"include_untranslated_keys,omitempty"`
	IncludeUnverifiedTranslations    *bool    `json:"include_unverified_translations,omitempty"`
	Lsp                              string   `json:"lsp"`
	Message                          *string  `json:"message,omitempty"`
	Priority                         *bool    `json:"priority,omitempty"`
	Quality                          *bool    `json:"quality,omitempty"`
	SourceLocaleId                   string   `json:"source_locale_id"`
	StyleguideId                     *string  `json:"styleguide_id,omitempty"`
	Tag                              *string  `json:"tag,omitempty"`
	TargetLocaleIds                  []string `json:"target_locale_ids"`
	TranslationType                  string   `json:"translation_type"`
	UnverifyTranslationsUponDelivery *bool    `json:"unverify_translations_upon_delivery,omitempty"`
}

type ProjectParams struct {
	Name                    string `json:"name"`
	SharesTranslationMemory *bool  `json:"shares_translation_memory,omitempty"`
}

type StyleguideParams struct {
	Audience           *string `json:"audience,omitempty"`
	Business           *string `json:"business,omitempty"`
	CompanyBranding    *string `json:"company_branding,omitempty"`
	Formatting         *string `json:"formatting,omitempty"`
	GlossaryTerms      *string `json:"glossary_terms,omitempty"`
	GrammarConsistency *string `json:"grammar_consistency,omitempty"`
	GrammaticalPerson  *string `json:"grammatical_person,omitempty"`
	LiteralTranslation *string `json:"literal_translation,omitempty"`
	OverallTone        *string `json:"overall_tone,omitempty"`
	Samples            *string `json:"samples,omitempty"`
	TargetAudience     *string `json:"target_audience,omitempty"`
	Title              string  `json:"title"`
	VocabularyType     *string `json:"vocabulary_type,omitempty"`
}

type TagParams struct {
	Name string `json:"name"`
}

type TranslationParams struct {
	Content      string  `json:"content"`
	Excluded     *bool   `json:"excluded,omitempty"`
	KeyId        string  `json:"key_id"`
	LocaleId     string  `json:"locale_id"`
	PluralSuffix *string `json:"plural_suffix,omitempty"`
	Unverified   *bool   `json:"unverified,omitempty"`
}

type LocaleFileImportParams struct {
	ConvertEmoji       *bool                   `json:"convert_emoji,omitempty"`
	File               string                  `json:"file"`
	Format             *string                 `json:"format,omitempty"`
	FormatOptions      *map[string]interface{} `json:"format_options,omitempty"`
	LocaleId           *string                 `json:"locale_id,omitempty"`
	SkipUnverification *bool                   `json:"skip_unverification,omitempty"`
	SkipUploadTags     *bool                   `json:"skip_upload_tags,omitempty"`
	Tags               []string                `json:"tags,omitempty"`
	UpdateTranslations *bool                   `json:"update_translations,omitempty"`
}

// Create a new authorization.
func AuthorizationCreate(params *AuthorizationParams) (*AuthorizationWithToken, error) {
	retVal := new(AuthorizationWithToken)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/authorizations")

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing authorization. Please note that this will revoke access for that token, so API calls using that token will stop working.
func AuthorizationDelete(id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/authorizations/%s", id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// List all your authorizations.
func AuthorizationList(page, perPage int) ([]*Authorization, error) {
	retVal := []*Authorization{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/authorizations")

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single authorization.
func AuthorizationShow(id string) (*Authorization, error) {
	retVal := new(Authorization)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/authorizations/%s", id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing authorization.
func AuthorizationUpdate(id string, params *AuthorizationParams) (*Authorization, error) {
	retVal := new(Authorization)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/authorizations/%s", id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("PATCH", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new blacklisted key.
func BlacklistKeyCreate(project_id string, params *BlacklistedKeyParams) (*BlacklistedKey, error) {
	retVal := new(BlacklistedKey)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/blacklisted_keys", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing blacklisted key.
func BlacklistKeyDelete(project_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/blacklisted_keys/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details on a single blacklisted key for a given project.
func BlacklistKeyShow(project_id, id string) (*BlacklistedKey, error) {
	retVal := new(BlacklistedKey)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/blacklisted_keys/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing blacklisted key.
func BlacklistKeyUpdate(project_id, id string, params *BlacklistedKeyParams) (*BlacklistedKey, error) {
	retVal := new(BlacklistedKey)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/blacklisted_keys/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("PATCH", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// List all blacklisted keys for the given project.
func BlacklistShow(project_id string, page, perPage int) ([]*BlacklistedKey, error) {
	retVal := []*BlacklistedKey{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/blacklisted_keys", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new comment for a key.
func CommentCreate(project_id, key_id string, params *CommentParams) (*Comment, error) {
	retVal := new(Comment)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments", project_id, key_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing comment.
func CommentDelete(project_id, key_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments/%s", project_id, key_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// List all comments for a key.
func CommentList(project_id, key_id string, page, perPage int) ([]*Comment, error) {
	retVal := []*Comment{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments", project_id, key_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Check if comment was marked as read. Returns 204 if read, 404 if unread.
func CommentMarkCheck(project_id, key_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments/%s/read", project_id, key_id, id)

		rc, err := sendRequest("GET", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Mark a comment as read
func CommentMarkRead(project_id, key_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments/%s/read", project_id, key_id, id)

		rc, err := sendRequest("PUT", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Mark a comment as unread
func CommentMarkUnread(project_id, key_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments/%s/read", project_id, key_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details on a single comment.
func CommentShow(project_id, key_id, id string) (*Comment, error) {
	retVal := new(Comment)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments/%s", project_id, key_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing comment.
func CommentUpdate(project_id, key_id, id string, params *CommentParams) (*Comment, error) {
	retVal := new(Comment)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments/%s", project_id, key_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("PATCH", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new key.
func KeyCreate(project_id string, params *TranslationKeyParams) (*TranslationKeyDetails, error) {
	retVal := new(TranslationKeyDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(paramsBuf)
		ctype := writer.FormDataContentType()

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

		if params.FormatValueType != nil {
			err := writer.WriteField("format_value_type", *params.FormatValueType)
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

		err := writer.WriteField("name", params.Name)
		if err != nil {
			return err
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

		for i := range params.Tags {
			err := writer.WriteField("tags[]", params.Tags[i])
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
		writer.Close()

		rc, err := sendRequest("POST", url, ctype, paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing key.
func KeyDelete(project_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

type KeyListParams struct {
	LocaleId   *string `json:"locale_id,omitempty"`
	Order      *string `json:"order,omitempty"`
	Sort       *string `json:"sort,omitempty"`
	Translated *bool   `json:"translated,omitempty"`
}

// List all keys for the given project.
func KeyList(project_id string, page, perPage int, params *KeyListParams) ([]*TranslationKey, error) {
	retVal := []*TranslationKey{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequestPaginated("GET", url, "application/json", paramsBuf, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single key for a given project.
func KeyShow(project_id, id string) (*TranslationKeyDetails, error) {
	retVal := new(TranslationKeyDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing key.
func KeyUpdate(project_id, id string, params *TranslationKeyParams) (*TranslationKeyDetails, error) {
	retVal := new(TranslationKeyDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(paramsBuf)
		ctype := writer.FormDataContentType()

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

		if params.FormatValueType != nil {
			err := writer.WriteField("format_value_type", *params.FormatValueType)
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

		err := writer.WriteField("name", params.Name)
		if err != nil {
			return err
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

		for i := range params.Tags {
			err := writer.WriteField("tags[]", params.Tags[i])
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
		writer.Close()

		rc, err := sendRequest("PATCH", url, ctype, paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new locale.
func LocaleCreate(project_id string, params *LocaleParams) (*LocaleDetails, error) {
	retVal := new(LocaleDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing locale.
func LocaleDelete(project_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

type LocaleDownloadParams struct {
	ConvertEmoji             *bool                   `json:"convert_emoji,omitempty"`
	Format                   string                  `json:"format"`
	FormatOptions            *map[string]interface{} `json:"format_options,omitempty"`
	IncludeEmptyTranslations *bool                   `json:"include_empty_translations,omitempty"`
	KeepNotranslateTags      *bool                   `json:"keep_notranslate_tags,omitempty"`
	TagId                    *string                 `json:"tag_id,omitempty"`
}

// Download a locale in a specific file format.
func LocaleDownload(project_id, id string, params *LocaleDownloadParams) ([]byte, error) {
	retVal := []byte{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales/%s/download", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("GET", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		retVal, err = ioutil.ReadAll(rc)
		return err

	}()
	return retVal, err
}

// List all locales for the given project.
func LocaleList(project_id string, page, perPage int) ([]*Locale, error) {
	retVal := []*Locale{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single locale for a given project.
func LocaleShow(project_id, id string) (*LocaleDetails, error) {
	retVal := new(LocaleDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing locale.
func LocaleUpdate(project_id, id string, params *LocaleParams) (*LocaleDetails, error) {
	retVal := new(LocaleDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("PATCH", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Confirm an existing order. Sends the order to the language service provider for processing. Please note that your access token must include the <code>orders.create</code> scope to confirm orders.
func OrderConfirm(project_id, id string) (*TranslationOrder, error) {
	retVal := new(TranslationOrder)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/orders/%s/confirm", project_id, id)

		rc, err := sendRequest("PATCH", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new order. Please note that your access token must include the <code>orders.create</code> scope to create orders.
func OrderCreate(project_id string, params *TranslationOrderParams) (*TranslationOrder, error) {
	retVal := new(TranslationOrder)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/orders", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Cancel an existing order. Must not yet be confirmed.
func OrderDelete(project_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/orders/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// List all orders for the given project.
func OrderList(project_id string, page, perPage int) ([]*TranslationOrder, error) {
	retVal := []*TranslationOrder{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/orders", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single order.
func OrderShow(project_id, id string) (*TranslationOrder, error) {
	retVal := new(TranslationOrder)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/orders/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new project.
func ProjectCreate(params *ProjectParams) (*ProjectDetails, error) {
	retVal := new(ProjectDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects")

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing project.
func ProjectDelete(id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s", id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// List all projects the current user has access to.
func ProjectList(page, perPage int) ([]*Project, error) {
	retVal := []*Project{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects")

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single project.
func ProjectShow(id string) (*ProjectDetails, error) {
	retVal := new(ProjectDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s", id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing project.
func ProjectUpdate(id string, params *ProjectParams) (*ProjectDetails, error) {
	retVal := new(ProjectDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s", id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("PATCH", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Show details for current User.
func ShowUser() (*User, error) {
	retVal := new(User)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/user")

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new style guide.
func StyleguideCreate(project_id string, params *StyleguideParams) (*StyleguideDetails, error) {
	retVal := new(StyleguideDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/styleguides", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing style guide.
func StyleguideDelete(project_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/styleguides/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// List all styleguides for the given project.
func StyleguideList(project_id string, page, perPage int) ([]*Styleguide, error) {
	retVal := []*Styleguide{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/styleguides", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single style guide.
func StyleguideShow(project_id, id string) (*StyleguideDetails, error) {
	retVal := new(StyleguideDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/styleguides/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Update an existing style guide.
func StyleguideUpdate(project_id, id string, params *StyleguideParams) (*StyleguideDetails, error) {
	retVal := new(StyleguideDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/styleguides/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("PATCH", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new tag.
func TagCreate(project_id string, params *TagParams) (*TagWithStats, error) {
	retVal := new(TagWithStats)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/tags", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Delete an existing tag.
func TagDelete(project_id, name string) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/tags/%s", project_id, name)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// List all tags for the given project.
func TagList(project_id string, page, perPage int) ([]*Tag, error) {
	retVal := []*Tag{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/tags", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details and progress information on a single tag for a given project.
func TagShow(project_id, name string) (*TagWithStats, error) {
	retVal := new(TagWithStats)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/tags/%s", project_id, name)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a translation.
func TranslationCreate(project_id string, params *TranslationParams) (*TranslationDetails, error) {
	retVal := new(TranslationDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("POST", url, "application/json", paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

type TranslationListAllParams struct {
	Order      *string    `json:"order,omitempty"`
	Since      *time.Time `json:"since,omitempty"`
	Sort       *string    `json:"sort,omitempty"`
	Unverified *bool      `json:"unverified,omitempty"`
}

// List translations for the given project.
func TranslationListAll(project_id string, page, perPage int, params *TranslationListAllParams) ([]*Translation, error) {
	retVal := []*Translation{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequestPaginated("GET", url, "application/json", paramsBuf, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

type TranslationListKeyParams struct {
	Order      *string    `json:"order,omitempty"`
	Since      *time.Time `json:"since,omitempty"`
	Sort       *string    `json:"sort,omitempty"`
	Unverified *bool      `json:"unverified,omitempty"`
}

// List translations for a specific key.
func TranslationListKey(project_id, key_id string, page, perPage int, params *TranslationListKeyParams) ([]*Translation, error) {
	retVal := []*Translation{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/translations", project_id, key_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequestPaginated("GET", url, "application/json", paramsBuf, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

type TranslationListLocaleParams struct {
	Order      *string    `json:"order,omitempty"`
	Since      *time.Time `json:"since,omitempty"`
	Sort       *string    `json:"sort,omitempty"`
	Unverified *bool      `json:"unverified,omitempty"`
}

// List translations for a specific locale.
func TranslationListLocale(project_id, locale_id string, page, perPage int, params *TranslationListLocaleParams) ([]*Translation, error) {
	retVal := []*Translation{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales/%s/translations", project_id, locale_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequestPaginated("GET", url, "application/json", paramsBuf, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single translation.
func TranslationShow(project_id, id string) (*TranslationDetails, error) {
	retVal := new(TranslationDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

type TranslationUpdateParams struct {
	Content      string  `json:"content"`
	Excluded     *bool   `json:"excluded,omitempty"`
	PluralSuffix *string `json:"plural_suffix,omitempty"`
	Unverified   *bool   `json:"unverified,omitempty"`
}

// Update an existing translation.
func TranslationUpdate(project_id, id string, params *TranslationUpdateParams) (*TranslationDetails, error) {
	retVal := new(TranslationDetails)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("PATCH", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Upload a new file to your project. This will extract all new content such as keys, translations, locales, tags etc. and store them in your project.
func UploadCreate(project_id string, params *LocaleFileImportParams) (*LocaleFileImportWithSummary, error) {
	retVal := new(LocaleFileImportWithSummary)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/uploads", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(paramsBuf)
		ctype := writer.FormDataContentType()

		if params.ConvertEmoji != nil {
			err := writer.WriteField("convert_emoji", strconv.FormatBool(*params.ConvertEmoji))
			if err != nil {
				return err
			}
		}

		part, err := writer.CreateFormFile("file", filepath.Base(params.File))
		if err != nil {
			return err
		}
		file, err := os.Open(params.File)
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

		if params.Format != nil {
			err := writer.WriteField("format", *params.Format)
			if err != nil {
				return err
			}
		}

		if params.FormatOptions != nil {
			for key, val := range *params.FormatOptions {
				err := writer.WriteField("format_options["+key+"]", val.(string))
				if err != nil {
					return err
				}
			}
		}

		if params.LocaleId != nil {
			err := writer.WriteField("locale_id", *params.LocaleId)
			if err != nil {
				return err
			}
		}

		if params.SkipUnverification != nil {
			err := writer.WriteField("skip_unverification", strconv.FormatBool(*params.SkipUnverification))
			if err != nil {
				return err
			}
		}

		if params.SkipUploadTags != nil {
			err := writer.WriteField("skip_upload_tags", strconv.FormatBool(*params.SkipUploadTags))
			if err != nil {
				return err
			}
		}

		for i := range params.Tags {
			err := writer.WriteField("tags[]", params.Tags[i])
			if err != nil {
				return err
			}
		}

		if params.UpdateTranslations != nil {
			err := writer.WriteField("update_translations", strconv.FormatBool(*params.UpdateTranslations))
			if err != nil {
				return err
			}
		}
		writer.Close()

		rc, err := sendRequest("POST", url, ctype, paramsBuf, 201)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// View details and summary for a single upload.
func UploadShow(project_id, id string) (*LocaleFileImportWithSummary, error) {
	retVal := new(LocaleFileImportWithSummary)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/uploads/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// List all versions for the given translation.
func VersionList(project_id, translation_id string, page, perPage int) ([]*TranslationVersion, error) {
	retVal := []*TranslationVersion{}
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations/%s/versions", project_id, translation_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get details on a single version.
func VersionShow(project_id, translation_id, id string) (*TranslationVersionWithUser, error) {
	retVal := new(TranslationVersionWithUser)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations/%s/versions/%s", project_id, translation_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}
