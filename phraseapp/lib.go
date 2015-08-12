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

type AffectedCount struct {
	RecordsAffected int64 `json:"records_affected"`
}

type AffectedResources struct {
	RecordsAffected int64 `json:"records_affected"`
}

type Authorization struct {
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
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

type ExcludeRule struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Format struct {
	ApiName         string `json:"api_name"`
	DefaultEncoding string `json:"default_encoding"`
	DefaultFile     string `json:"default_file"`
	Description     string `json:"description"`
	Exportable      bool   `json:"exportable"`
	Extension       string `json:"extension"`
	Importable      bool   `json:"importable"`
	Name            string `json:"name"`
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
	PluralForms  []string       `json:"plural_forms"`
	Rtl          bool           `json:"rtl"`
	SourceLocale *LocalePreview `json:"source_locale"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type LocaleDetails struct {
	Locale

	Statistics *LocaleStatistics `json:"statistics"`
}

type LocaleFileImport struct {
	CreatedAt  time.Time `json:"created_at"`
	FileFormat string    `json:"file_format"`
	Id         string    `json:"id"`
	State      string    `json:"state"`
	UpdatedAt  time.Time `json:"updated_at"`
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
	CreatedAt  time.Time `json:"created_at"`
	Id         string    `json:"id"`
	MainFormat string    `json:"main_format"`
	Name       string    `json:"name"`
	UpdatedAt  time.Time `json:"updated_at"`
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
	Placeholders []string       `json:"placeholders"`
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
	DataType    string    `json:"data_type"`
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
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Note      string     `json:"note"`
	Scopes    []string   `json:"scopes,omitempty"`
}

func (params *AuthorizationParams) ApplyDefaults(defaults map[string]interface{}) (*AuthorizationParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(AuthorizationParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type CommentParams struct {
	Message string `json:"message"`
}

func (params *CommentParams) ApplyDefaults(defaults map[string]interface{}) (*CommentParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(CommentParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type ExcludeRuleParams struct {
	Name string `json:"name"`
}

func (params *ExcludeRuleParams) ApplyDefaults(defaults map[string]interface{}) (*ExcludeRuleParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(ExcludeRuleParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type TranslationKeyParams struct {
	DataType              *string `json:"data_type,omitempty"`
	Description           *string `json:"description,omitempty"`
	LocalizedFormatKey    *string `json:"localized_format_key,omitempty"`
	LocalizedFormatString *string `json:"localized_format_string,omitempty"`
	MaxCharactersAllowed  *int64  `json:"max_characters_allowed,omitempty"`
	Name                  string  `json:"name"`
	NamePlural            *string `json:"name_plural,omitempty"`
	OriginalFile          *string `json:"original_file,omitempty"`
	Plural                *bool   `json:"plural,omitempty"`
	RemoveScreenshot      *bool   `json:"remove_screenshot,omitempty"`
	Screenshot            *string `json:"screenshot,omitempty"`
	Tags                  *string `json:"tags,omitempty"`
	Unformatted           *bool   `json:"unformatted,omitempty"`
	XmlSpacePreserve      *bool   `json:"xml_space_preserve,omitempty"`
}

func (params *TranslationKeyParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationKeyParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationKeyParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type LocaleParams struct {
	Code           string  `json:"code"`
	Default        *bool   `json:"default,omitempty"`
	Main           *bool   `json:"main,omitempty"`
	Name           string  `json:"name"`
	Rtl            *bool   `json:"rtl,omitempty"`
	SourceLocaleId *string `json:"source_locale_id,omitempty"`
}

func (params *LocaleParams) ApplyDefaults(defaults map[string]interface{}) (*LocaleParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(LocaleParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
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

func (params *TranslationOrderParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationOrderParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationOrderParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type ProjectParams struct {
	Name                    string `json:"name"`
	SharesTranslationMemory *bool  `json:"shares_translation_memory,omitempty"`
}

func (params *ProjectParams) ApplyDefaults(defaults map[string]interface{}) (*ProjectParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(ProjectParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
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

func (params *StyleguideParams) ApplyDefaults(defaults map[string]interface{}) (*StyleguideParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(StyleguideParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type TagParams struct {
	Name string `json:"name"`
}

func (params *TagParams) ApplyDefaults(defaults map[string]interface{}) (*TagParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TagParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type TranslationParams struct {
	Content      string  `json:"content"`
	Excluded     *bool   `json:"excluded,omitempty"`
	KeyId        string  `json:"key_id"`
	LocaleId     string  `json:"locale_id"`
	PluralSuffix *string `json:"plural_suffix,omitempty"`
	Unverified   *bool   `json:"unverified,omitempty"`
}

func (params *TranslationParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

type LocaleFileImportParams struct {
	ConvertEmoji       *bool   `json:"convert_emoji,omitempty"`
	File               string  `json:"file"`
	FileFormat         *string `json:"file_format,omitempty"`
	LocaleId           *string `json:"locale_id,omitempty"`
	SkipUnverification *bool   `json:"skip_unverification,omitempty"`
	SkipUploadTags     *bool   `json:"skip_upload_tags,omitempty"`
	Tags               *string `json:"tags,omitempty"`
	UpdateTranslations *bool   `json:"update_translations,omitempty"`
}

func (params *LocaleFileImportParams) ApplyDefaults(defaults map[string]interface{}) (*LocaleFileImportParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(LocaleFileImportParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Create a new authorization.
func AuthorizationCreate(params *AuthorizationParams) (*AuthorizationWithToken, error) {
	retVal := new(AuthorizationWithToken)
	err := func() error {
		url := fmt.Sprintf("/v2/authorizations")

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

// Delete an existing authorization. API calls using that token will stop working.
func AuthorizationDelete(id string) error {

	err := func() error {
		url := fmt.Sprintf("/v2/authorizations/%s", id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details on a single authorization.
func AuthorizationShow(id string) (*Authorization, error) {
	retVal := new(Authorization)
	err := func() error {
		url := fmt.Sprintf("/v2/authorizations/%s", id)

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
		url := fmt.Sprintf("/v2/authorizations/%s", id)

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

// List all your authorizations.
func AuthorizationsList(page, perPage int) ([]*Authorization, error) {
	retVal := []*Authorization{}
	err := func() error {
		url := fmt.Sprintf("/v2/authorizations")

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
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments", project_id, key_id)

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
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments/%s", project_id, key_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Check if comment was marked as read. Returns 204 if read, 404 if unread.
func CommentMarkCheck(project_id, key_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments/%s/read", project_id, key_id, id)

		rc, err := sendRequest("GET", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Mark a comment as read.
func CommentMarkRead(project_id, key_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments/%s/read", project_id, key_id, id)

		rc, err := sendRequest("PATCH", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Mark a comment as unread.
func CommentMarkUnread(project_id, key_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments/%s/read", project_id, key_id, id)

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
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments/%s", project_id, key_id, id)

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
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments/%s", project_id, key_id, id)

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

// List all comments for a key.
func CommentsList(project_id, key_id string, page, perPage int) ([]*Comment, error) {
	retVal := []*Comment{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/comments", project_id, key_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new blacklisted key.
func ExcludeRuleCreate(project_id string, params *ExcludeRuleParams) (*BlacklistedKey, error) {
	retVal := new(BlacklistedKey)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/blacklisted_keys", project_id)

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
func ExcludeRuleDelete(project_id, id string) error {

	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/blacklisted_keys/%s", project_id, id)

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
func ExcludeRuleShow(project_id, id string) (*BlacklistedKey, error) {
	retVal := new(BlacklistedKey)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/blacklisted_keys/%s", project_id, id)

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
func ExcludeRuleUpdate(project_id, id string, params *ExcludeRuleParams) (*BlacklistedKey, error) {
	retVal := new(BlacklistedKey)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/blacklisted_keys/%s", project_id, id)

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
func ExcludeRulesIndex(project_id string, page, perPage int) ([]*BlacklistedKey, error) {
	retVal := []*BlacklistedKey{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/blacklisted_keys", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Get a handy list of all localization file formats supported in PhraseApp.
func FormatsList(page, perPage int) ([]*Format, error) {
	retVal := []*Format{}
	err := func() error {
		url := fmt.Sprintf("/v2/formats")

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
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
		url := fmt.Sprintf("/v2/projects/%s/keys", project_id)

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

		if params.Name != "" {
			err := writer.WriteField("name", params.Name)
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
		url := fmt.Sprintf("/v2/projects/%s/keys/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details on a single key for a given project.
func KeyShow(project_id, id string) (*TranslationKeyDetails, error) {
	retVal := new(TranslationKeyDetails)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/%s", project_id, id)

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
		url := fmt.Sprintf("/v2/projects/%s/keys/%s", project_id, id)

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

		if params.Name != "" {
			err := writer.WriteField("name", params.Name)
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

type KeysDeleteParams struct {
	LocaleId *string `json:"locale_id,omitempty"`
	Q        *string `json:"q,omitempty"`
}

func (params *KeysDeleteParams) ApplyDefaults(defaults map[string]interface{}) (*KeysDeleteParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(KeysDeleteParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Delete all keys matching query. Same constraints as list.
func KeysDelete(project_id string, params *KeysDeleteParams) (*AffectedResources, error) {
	retVal := new(AffectedResources)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("DELETE", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

type KeysListParams struct {
	LocaleId *string `json:"locale_id,omitempty"`
	Order    *string `json:"order,omitempty"`
	Q        *string `json:"q,omitempty"`
	Sort     *string `json:"sort,omitempty"`
}

func (params *KeysListParams) ApplyDefaults(defaults map[string]interface{}) (*KeysListParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(KeysListParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// List all keys for the given project. Alternatively you can POST requests to /search.
func KeysList(project_id string, page, perPage int, params *KeysListParams) ([]*TranslationKey, error) {
	retVal := []*TranslationKey{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys", project_id)

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

type KeysSearchParams struct {
	LocaleId *string `json:"locale_id,omitempty"`
	Order    *string `json:"order,omitempty"`
	Q        *string `json:"q,omitempty"`
	Sort     *string `json:"sort,omitempty"`
}

func (params *KeysSearchParams) ApplyDefaults(defaults map[string]interface{}) (*KeysSearchParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(KeysSearchParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Search keys for the given project matching query.
func KeysSearch(project_id string, page, perPage int, params *KeysSearchParams) ([]*TranslationKey, error) {
	retVal := []*TranslationKey{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/search", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequestPaginated("POST", url, "application/json", paramsBuf, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

type KeysTagParams struct {
	LocaleId *string `json:"locale_id,omitempty"`
	Q        *string `json:"q,omitempty"`
	Tags     string  `json:"tags"`
}

func (params *KeysTagParams) ApplyDefaults(defaults map[string]interface{}) (*KeysTagParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(KeysTagParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Tags all keys matching query. Same constraints as list.
func KeysTag(project_id string, params *KeysTagParams) (*AffectedResources, error) {
	retVal := new(AffectedResources)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/tag", project_id)

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

type KeysUntagParams struct {
	LocaleId *string `json:"locale_id,omitempty"`
	Q        *string `json:"q,omitempty"`
	Tags     string  `json:"tags"`
}

func (params *KeysUntagParams) ApplyDefaults(defaults map[string]interface{}) (*KeysUntagParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(KeysUntagParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Removes specified tags from keys matching query.
func KeysUntag(project_id string, params *KeysUntagParams) (*AffectedResources, error) {
	retVal := new(AffectedResources)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/tag", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequest("DELETE", url, "application/json", paramsBuf, 200)
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
		url := fmt.Sprintf("/v2/projects/%s/locales", project_id)

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
		url := fmt.Sprintf("/v2/projects/%s/locales/%s", project_id, id)

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
	FileFormat               string                  `json:"file_format"`
	FormatOptions            *map[string]interface{} `json:"format_options,omitempty"`
	IncludeEmptyTranslations *bool                   `json:"include_empty_translations,omitempty"`
	KeepNotranslateTags      *bool                   `json:"keep_notranslate_tags,omitempty"`
	Tag                      *string                 `json:"tag,omitempty"`
}

func (params *LocaleDownloadParams) ApplyDefaults(defaults map[string]interface{}) (*LocaleDownloadParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(LocaleDownloadParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Download a locale in a specific file format.
func LocaleDownload(project_id, id string, params *LocaleDownloadParams) ([]byte, error) {
	retVal := []byte{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/locales/%s/download", project_id, id)

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

// Get details on a single locale for a given project.
func LocaleShow(project_id, id string) (*LocaleDetails, error) {
	retVal := new(LocaleDetails)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/locales/%s", project_id, id)

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
		url := fmt.Sprintf("/v2/projects/%s/locales/%s", project_id, id)

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

// List all locales for the given project.
func LocalesList(project_id string, page, perPage int) ([]*Locale, error) {
	retVal := []*Locale{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/locales", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Confirm an existing order and send it to the provider for translation. Same constraints as for create.
func OrderConfirm(project_id, id string) (*TranslationOrder, error) {
	retVal := new(TranslationOrder)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/orders/%s/confirm", project_id, id)

		rc, err := sendRequest("PATCH", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// Create a new order. Access token scope must include <code>orders.create</code>.
func OrderCreate(project_id string, params *TranslationOrderParams) (*TranslationOrder, error) {
	retVal := new(TranslationOrder)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/orders", project_id)

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
		url := fmt.Sprintf("/v2/projects/%s/orders/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details on a single order.
func OrderShow(project_id, id string) (*TranslationOrder, error) {
	retVal := new(TranslationOrder)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/orders/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// List all orders for the given project.
func OrdersList(project_id string, page, perPage int) ([]*TranslationOrder, error) {
	retVal := []*TranslationOrder{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/orders", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
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
		url := fmt.Sprintf("/v2/projects")

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
		url := fmt.Sprintf("/v2/projects/%s", id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details on a single project.
func ProjectShow(id string) (*ProjectDetails, error) {
	retVal := new(ProjectDetails)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s", id)

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
		url := fmt.Sprintf("/v2/projects/%s", id)

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

// List all projects the current user has access to.
func ProjectsList(page, perPage int) ([]*Project, error) {
	retVal := []*Project{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects")

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
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
		url := fmt.Sprintf("/v2/user")

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
		url := fmt.Sprintf("/v2/projects/%s/styleguides", project_id)

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
		url := fmt.Sprintf("/v2/projects/%s/styleguides/%s", project_id, id)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details on a single style guide.
func StyleguideShow(project_id, id string) (*StyleguideDetails, error) {
	retVal := new(StyleguideDetails)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/styleguides/%s", project_id, id)

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
		url := fmt.Sprintf("/v2/projects/%s/styleguides/%s", project_id, id)

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

// List all styleguides for the given project.
func StyleguidesList(project_id string, page, perPage int) ([]*Styleguide, error) {
	retVal := []*Styleguide{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/styleguides", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
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
		url := fmt.Sprintf("/v2/projects/%s/tags", project_id)

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
		url := fmt.Sprintf("/v2/projects/%s/tags/%s", project_id, name)

		rc, err := sendRequest("DELETE", url, "", nil, 204)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
}

// Get details and progress information on a single tag for a given project.
func TagShow(project_id, name string) (*TagWithStats, error) {
	retVal := new(TagWithStats)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/tags/%s", project_id, name)

		rc, err := sendRequest("GET", url, "", nil, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

// List all tags for the given project.
func TagsList(project_id string, page, perPage int) ([]*Tag, error) {
	retVal := []*Tag{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/tags", project_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
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
		url := fmt.Sprintf("/v2/projects/%s/translations", project_id)

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

// Get details on a single translation.
func TranslationShow(project_id, id string) (*TranslationDetails, error) {
	retVal := new(TranslationDetails)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/%s", project_id, id)

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

func (params *TranslationUpdateParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationUpdateParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationUpdateParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Update an existing translation.
func TranslationUpdate(project_id, id string, params *TranslationUpdateParams) (*TranslationDetails, error) {
	retVal := new(TranslationDetails)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/%s", project_id, id)

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

type TranslationsByKeyParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsByKeyParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsByKeyParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsByKeyParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// List translations for a specific key.
func TranslationsByKey(project_id, key_id string, page, perPage int, params *TranslationsByKeyParams) ([]*Translation, error) {
	retVal := []*Translation{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/keys/%s/translations", project_id, key_id)

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

type TranslationsByLocaleParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsByLocaleParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsByLocaleParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsByLocaleParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// List translations for a specific locale.
func TranslationsByLocale(project_id, locale_id string, page, perPage int, params *TranslationsByLocaleParams) ([]*Translation, error) {
	retVal := []*Translation{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/locales/%s/translations", project_id, locale_id)

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

type TranslationsExcludeParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsExcludeParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsExcludeParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsExcludeParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Exclude translations matching query from locale export.
func TranslationsExclude(project_id string, params *TranslationsExcludeParams) (*AffectedCount, error) {
	retVal := new(AffectedCount)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/exclude", project_id)

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

type TranslationsIncludeParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsIncludeParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsIncludeParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsIncludeParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Include translations matching query in locale export.
func TranslationsInclude(project_id string, params *TranslationsIncludeParams) (*AffectedCount, error) {
	retVal := new(AffectedCount)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/include", project_id)

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

type TranslationsListParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsListParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsListParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsListParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// List translations for the given project. Alternatively, POST request to /search
func TranslationsList(project_id string, page, perPage int, params *TranslationsListParams) ([]*Translation, error) {
	retVal := []*Translation{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations", project_id)

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

type TranslationsSearchParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsSearchParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsSearchParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsSearchParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// List translations for the given project if you exceed GET request limitations on translations list.
func TranslationsSearch(project_id string, page, perPage int, params *TranslationsSearchParams) ([]*Translation, error) {
	retVal := []*Translation{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/search", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		rc, err := sendRequestPaginated("POST", url, "application/json", paramsBuf, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

type TranslationsUnverifyParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsUnverifyParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsUnverifyParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsUnverifyParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Mark translations matching query as unverified.
func TranslationsUnverify(project_id string, params *TranslationsUnverifyParams) (*AffectedCount, error) {
	retVal := new(AffectedCount)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/unverify", project_id)

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

type TranslationsVerifyParams struct {
	Order *string `json:"order,omitempty"`
	Q     *string `json:"q,omitempty"`
	Sort  *string `json:"sort,omitempty"`
}

func (params *TranslationsVerifyParams) ApplyDefaults(defaults map[string]interface{}) (*TranslationsVerifyParams, error) {
	str, err := json.Marshal(defaults)
	if err != nil {
		return params, err
	}
	defaultParams := new(TranslationsVerifyParams)
	err = json.Unmarshal(str, defaultParams)
	if err != nil {
		return params, err
	}

	return defaultParams, nil
}

// Verify translations matching query.
func TranslationsVerify(project_id string, params *TranslationsVerifyParams) (*AffectedCount, error) {
	retVal := new(AffectedCount)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/verify", project_id)

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

// Upload a new language file. Creates necessary resources in your project.
func UploadCreate(project_id string, params *LocaleFileImportParams) (*LocaleFileImportWithSummary, error) {
	retVal := new(LocaleFileImportWithSummary)
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/uploads", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(paramsBuf)
		ctype := writer.FormDataContentType()

		if params.ConvertEmoji != nil {
			err := writer.WriteField("convert_emoji", strconv.FormatBool(*params.ConvertEmoji))
			if err != nil {
				return err
			}
		}

		if params.File != "" {
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
		}

		if params.FileFormat != nil {
			err := writer.WriteField("file_format", *params.FileFormat)
			if err != nil {
				return err
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

		if params.Tags != nil {
			err := writer.WriteField("tags", *params.Tags)
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
		url := fmt.Sprintf("/v2/projects/%s/uploads/%s", project_id, id)

		rc, err := sendRequest("GET", url, "", nil, 200)
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
		url := fmt.Sprintf("/v2/projects/%s/translations/%s/versions/%s", project_id, translation_id, id)

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
func VersionsList(project_id, translation_id string, page, perPage int) ([]*TranslationVersion, error) {
	retVal := []*TranslationVersion{}
	err := func() error {
		url := fmt.Sprintf("/v2/projects/%s/translations/%s/versions", project_id, translation_id)

		rc, err := sendRequestPaginated("GET", url, "", nil, 200, page, perPage)
		if err != nil {
			return err
		}
		defer rc.Close()

		return json.NewDecoder(rc).Decode(&retVal)

	}()
	return retVal, err
}

func GetUserAgent() string {
	return "PhraseApp go (test)"
}
