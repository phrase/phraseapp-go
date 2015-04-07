package phraseapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

type AuthorizationFull struct {
	Authorization
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

type BlacklistedKeyFull struct {
	BlacklistedKey
}

type Comment struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *UserType `json:"user"`
}

type CommentFull struct {
	Comment
}

type KeyType struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Plural bool   `json:"plural"`
}

type Locale struct {
	Code         string      `json:"code"`
	CreatedAt    time.Time   `json:"created_at"`
	Default      bool        `json:"default"`
	Id           string      `json:"id"`
	Main         bool        `json:"main"`
	Name         string      `json:"name"`
	Rtl          bool        `json:"rtl"`
	SourceLocale *LocaleType `json:"source_locale"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type LocaleFileImport struct {
	CreatedAt time.Time `json:"created_at"`
	Format    string    `json:"format"`
	Id        string    `json:"id"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LocaleFileImportFull struct {
	LocaleFileImport

	Summary SummaryType `json:"summary"`
}

type LocaleFull struct {
	Locale

	Statistics *LocaleStatisticsType `json:"statistics"`
}

type LocaleStatisticsType struct {
}

type LocaleType struct {
	Code string `json:"code"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectFull struct {
	Project

	SharesTranslationMemory bool `json:"shares_translation_memory"`
}

type StatisticsListItem struct {
	Locale     *LocaleType    `json:"locale"`
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

type StyleguideFull struct {
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

type StyleguideType struct {
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

type TagFull struct {
	Tag

	Statistics []*StatisticsListItem `json:"statistics"`
}

type Translation struct {
	Content      string      `json:"content"`
	CreatedAt    time.Time   `json:"created_at"`
	Excluded     bool        `json:"excluded"`
	Id           string      `json:"id"`
	Key          *KeyType    `json:"key"`
	Locale       *LocaleType `json:"locale"`
	PluralSuffix string      `json:"plural_suffix"`
	Unverified   bool        `json:"unverified"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type TranslationFull struct {
	Translation

	User      *UserType `json:"user"`
	WordCount int64     `json:"word_count"`
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

type TranslationKeyFull struct {
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
	AmountInCents                    int64           `json:"amount_in_cents"`
	CreatedAt                        time.Time       `json:"created_at"`
	Currency                         string          `json:"currency"`
	Id                               string          `json:"id"`
	Lsp                              string          `json:"lsp"`
	Message                          string          `json:"message"`
	Priority                         bool            `json:"priority"`
	ProgressPercent                  int64           `json:"progress_percent"`
	Quality                          bool            `json:"quality"`
	SourceLocale                     *LocaleType     `json:"source_locale"`
	State                            string          `json:"state"`
	Styleguide                       *StyleguideType `json:"styleguide"`
	Tag                              string          `json:"tag"`
	TargetLocales                    []*LocaleType   `json:"target_locales"`
	TranslationType                  string          `json:"translation_type"`
	UnverifyTranslationsUponDelivery bool            `json:"unverify_translations_upon_delivery"`
	UpdatedAt                        time.Time       `json:"updated_at"`
}

type TranslationOrderFull struct {
	TranslationOrder
}

type TranslationVersion struct {
	ChangedAt    time.Time   `json:"changed_at"`
	Content      string      `json:"content"`
	CreatedAt    time.Time   `json:"created_at"`
	Id           string      `json:"id"`
	Key          *KeyType    `json:"key"`
	Locale       *LocaleType `json:"locale"`
	PluralSuffix string      `json:"plural_suffix"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type TranslationVersionFull struct {
	TranslationVersion

	User *UserType `json:"user"`
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

type UserType struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type AuthorizationParams struct {
	Note   *string   `json:"note"`
	Scopes *[]string `json:"scopes,omitempty"`
}

func (p *AuthorizationParams) SetNote(val string) error {
	p.Note = &val
	return nil
}

func (p *AuthorizationParams) SetScopes(val string) error {
	a := strings.Split(val, ",")
	p.Scopes = &a
	return nil
}

func (params *AuthorizationParams) validate() error {
	if params.Note == nil || *params.Note == "" {
		return fmt.Errorf("Required parameter \"Note\" of \"AuthorizationParams\" not set")
	}

	return nil
}

type BlacklistedKeyParams struct {
	Name *string `json:"name"`
}

func (p *BlacklistedKeyParams) SetName(val string) error {
	p.Name = &val
	return nil
}

func (params *BlacklistedKeyParams) validate() error {
	if params.Name == nil || *params.Name == "" {
		return fmt.Errorf("Required parameter \"Name\" of \"BlacklistedKeyParams\" not set")
	}

	return nil
}

type CommentParams struct {
	Message *string `json:"message"`
}

func (p *CommentParams) SetMessage(val string) error {
	p.Message = &val
	return nil
}

func (params *CommentParams) validate() error {
	if params.Message == nil || *params.Message == "" {
		return fmt.Errorf("Required parameter \"Message\" of \"CommentParams\" not set")
	}

	return nil
}

type TranslationKeyParams struct {
	DataType             *string   `json:"data_type,omitempty"`
	Description          *string   `json:"description,omitempty"`
	FormatValueType      *string   `json:"format_value_type,omitempty"`
	MaxCharactersAllowed *int64    `json:"max_characters_allowed,omitempty"`
	Name                 *string   `json:"name"`
	NamePlural           *string   `json:"name_plural,omitempty"`
	OriginalFile         *string   `json:"original_file,omitempty"`
	Plural               *bool     `json:"plural,omitempty"`
	RemoveScreenshot     *bool     `json:"remove_screenshot,omitempty"`
	Screenshot           *string   `json:"screenshot,omitempty"`
	Tags                 *[]string `json:"tags,omitempty"`
	Unformatted          *bool     `json:"unformatted,omitempty"`
	XmlSpacePreserve     *bool     `json:"xml_space_preserve,omitempty"`
}

func (p *TranslationKeyParams) SetDataType(val string) error {
	p.DataType = &val
	return nil
}

func (p *TranslationKeyParams) SetDescription(val string) error {
	p.Description = &val
	return nil
}

func (p *TranslationKeyParams) SetFormatValueType(val string) error {
	p.FormatValueType = &val
	return nil
}

func (p *TranslationKeyParams) SetMaxCharactersAllowed(val string) error {
	i, err := strconv.ParseInt(val, 10, 64)
	p.MaxCharactersAllowed = &i
	return err
}

func (p *TranslationKeyParams) SetName(val string) error {
	p.Name = &val
	return nil
}

func (p *TranslationKeyParams) SetNamePlural(val string) error {
	p.NamePlural = &val
	return nil
}

func (p *TranslationKeyParams) SetOriginalFile(val string) error {
	p.OriginalFile = &val
	return nil
}

func (p *TranslationKeyParams) SetPlural(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Plural = &b
	return nil
}

func (p *TranslationKeyParams) SetRemoveScreenshot(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.RemoveScreenshot = &b
	return nil
}

func (p *TranslationKeyParams) SetScreenshot(val string) error {
	p.Screenshot = &val
	return nil
}

func (p *TranslationKeyParams) SetTags(val string) error {
	a := strings.Split(val, ",")
	p.Tags = &a
	return nil
}

func (p *TranslationKeyParams) SetUnformatted(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Unformatted = &b
	return nil
}

func (p *TranslationKeyParams) SetXmlSpacePreserve(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.XmlSpacePreserve = &b
	return nil
}

func (params *TranslationKeyParams) validate() error {
	if params.Name == nil || *params.Name == "" {
		return fmt.Errorf("Required parameter \"Name\" of \"TranslationKeyParams\" not set")
	}

	return nil
}

type LocaleParams struct {
	Code           *string `json:"code"`
	Default        *bool   `json:"default,omitempty"`
	Main           *bool   `json:"main,omitempty"`
	Name           *string `json:"name"`
	Rtl            *bool   `json:"rtl,omitempty"`
	SourceLocaleId *string `json:"source_locale_id,omitempty"`
}

func (p *LocaleParams) SetCode(val string) error {
	p.Code = &val
	return nil
}

func (p *LocaleParams) SetDefault(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Default = &b
	return nil
}

func (p *LocaleParams) SetMain(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Main = &b
	return nil
}

func (p *LocaleParams) SetName(val string) error {
	p.Name = &val
	return nil
}

func (p *LocaleParams) SetRtl(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Rtl = &b
	return nil
}

func (p *LocaleParams) SetSourceLocaleId(val string) error {
	p.SourceLocaleId = &val
	return nil
}

func (params *LocaleParams) validate() error {
	if params.Code == nil || *params.Code == "" {
		return fmt.Errorf("Required parameter \"Code\" of \"LocaleParams\" not set")
	}
	if params.Name == nil || *params.Name == "" {
		return fmt.Errorf("Required parameter \"Name\" of \"LocaleParams\" not set")
	}

	return nil
}

type TranslationOrderParams struct {
	Category                         *string   `json:"category"`
	IncludeUntranslatedKeys          *bool     `json:"include_untranslated_keys,omitempty"`
	IncludeUnverifiedTranslations    *bool     `json:"include_unverified_translations,omitempty"`
	Lsp                              *string   `json:"lsp"`
	Message                          *string   `json:"message,omitempty"`
	Priority                         *bool     `json:"priority,omitempty"`
	Quality                          *bool     `json:"quality,omitempty"`
	SourceLocaleId                   *string   `json:"source_locale_id"`
	StyleguideId                     *string   `json:"styleguide_id,omitempty"`
	Tag                              *string   `json:"tag,omitempty"`
	TargetLocaleIds                  *[]string `json:"target_locale_ids"`
	TranslationType                  *string   `json:"translation_type"`
	UnverifyTranslationsUponDelivery *bool     `json:"unverify_translations_upon_delivery,omitempty"`
}

func (p *TranslationOrderParams) SetCategory(val string) error {
	p.Category = &val
	return nil
}

func (p *TranslationOrderParams) SetIncludeUntranslatedKeys(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.IncludeUntranslatedKeys = &b
	return nil
}

func (p *TranslationOrderParams) SetIncludeUnverifiedTranslations(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.IncludeUnverifiedTranslations = &b
	return nil
}

func (p *TranslationOrderParams) SetLsp(val string) error {
	p.Lsp = &val
	return nil
}

func (p *TranslationOrderParams) SetMessage(val string) error {
	p.Message = &val
	return nil
}

func (p *TranslationOrderParams) SetPriority(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Priority = &b
	return nil
}

func (p *TranslationOrderParams) SetQuality(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Quality = &b
	return nil
}

func (p *TranslationOrderParams) SetSourceLocaleId(val string) error {
	p.SourceLocaleId = &val
	return nil
}

func (p *TranslationOrderParams) SetStyleguideId(val string) error {
	p.StyleguideId = &val
	return nil
}

func (p *TranslationOrderParams) SetTag(val string) error {
	p.Tag = &val
	return nil
}

func (p *TranslationOrderParams) SetTargetLocaleIds(val string) error {
	a := strings.Split(val, ",")
	p.TargetLocaleIds = &a
	return nil
}

func (p *TranslationOrderParams) SetTranslationType(val string) error {
	p.TranslationType = &val
	return nil
}

func (p *TranslationOrderParams) SetUnverifyTranslationsUponDelivery(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.UnverifyTranslationsUponDelivery = &b
	return nil
}

func (params *TranslationOrderParams) validate() error {
	if params.Category == nil || *params.Category == "" {
		return fmt.Errorf("Required parameter \"Category\" of \"TranslationOrderParams\" not set")
	}
	if params.Lsp == nil || *params.Lsp == "" {
		return fmt.Errorf("Required parameter \"Lsp\" of \"TranslationOrderParams\" not set")
	}
	if params.SourceLocaleId == nil {
		return fmt.Errorf("Required parameter \"SourceLocaleId\" of \"TranslationOrderParams\" not set")
	}
	if params.TargetLocaleIds == nil {
		return fmt.Errorf("Required parameter \"TargetLocaleIds\" of \"TranslationOrderParams\" not set")
	}
	if params.TranslationType == nil || *params.TranslationType == "" {
		return fmt.Errorf("Required parameter \"TranslationType\" of \"TranslationOrderParams\" not set")
	}

	return nil
}

type ProjectParams struct {
	Name                    *string `json:"name"`
	SharesTranslationMemory *bool   `json:"shares_translation_memory,omitempty"`
}

func (p *ProjectParams) SetName(val string) error {
	p.Name = &val
	return nil
}

func (p *ProjectParams) SetSharesTranslationMemory(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.SharesTranslationMemory = &b
	return nil
}

func (params *ProjectParams) validate() error {
	if params.Name == nil || *params.Name == "" {
		return fmt.Errorf("Required parameter \"Name\" of \"ProjectParams\" not set")
	}

	return nil
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
	Title              *string `json:"title"`
	VocabularyType     *string `json:"vocabulary_type,omitempty"`
}

func (p *StyleguideParams) SetAudience(val string) error {
	p.Audience = &val
	return nil
}

func (p *StyleguideParams) SetBusiness(val string) error {
	p.Business = &val
	return nil
}

func (p *StyleguideParams) SetCompanyBranding(val string) error {
	p.CompanyBranding = &val
	return nil
}

func (p *StyleguideParams) SetFormatting(val string) error {
	p.Formatting = &val
	return nil
}

func (p *StyleguideParams) SetGlossaryTerms(val string) error {
	p.GlossaryTerms = &val
	return nil
}

func (p *StyleguideParams) SetGrammarConsistency(val string) error {
	p.GrammarConsistency = &val
	return nil
}

func (p *StyleguideParams) SetGrammaticalPerson(val string) error {
	p.GrammaticalPerson = &val
	return nil
}

func (p *StyleguideParams) SetLiteralTranslation(val string) error {
	p.LiteralTranslation = &val
	return nil
}

func (p *StyleguideParams) SetOverallTone(val string) error {
	p.OverallTone = &val
	return nil
}

func (p *StyleguideParams) SetSamples(val string) error {
	p.Samples = &val
	return nil
}

func (p *StyleguideParams) SetTargetAudience(val string) error {
	p.TargetAudience = &val
	return nil
}

func (p *StyleguideParams) SetTitle(val string) error {
	p.Title = &val
	return nil
}

func (p *StyleguideParams) SetVocabularyType(val string) error {
	p.VocabularyType = &val
	return nil
}

func (params *StyleguideParams) validate() error {
	if params.Title == nil || *params.Title == "" {
		return fmt.Errorf("Required parameter \"Title\" of \"StyleguideParams\" not set")
	}

	return nil
}

type TagParams struct {
	Name *string `json:"name"`
}

func (p *TagParams) SetName(val string) error {
	p.Name = &val
	return nil
}

func (params *TagParams) validate() error {
	if params.Name == nil || *params.Name == "" {
		return fmt.Errorf("Required parameter \"Name\" of \"TagParams\" not set")
	}

	return nil
}

type TranslationParams struct {
	Content      *string `json:"content"`
	Excluded     *bool   `json:"excluded,omitempty"`
	KeyId        *string `json:"key_id"`
	LocaleId     *string `json:"locale_id"`
	PluralSuffix *string `json:"plural_suffix,omitempty"`
	Unverified   *bool   `json:"unverified,omitempty"`
}

func (p *TranslationParams) SetContent(val string) error {
	p.Content = &val
	return nil
}

func (p *TranslationParams) SetExcluded(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Excluded = &b
	return nil
}

func (p *TranslationParams) SetKeyId(val string) error {
	p.KeyId = &val
	return nil
}

func (p *TranslationParams) SetLocaleId(val string) error {
	p.LocaleId = &val
	return nil
}

func (p *TranslationParams) SetPluralSuffix(val string) error {
	p.PluralSuffix = &val
	return nil
}

func (p *TranslationParams) SetUnverified(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Unverified = &b
	return nil
}

func (params *TranslationParams) validate() error {
	if params.Content == nil || *params.Content == "" {
		return fmt.Errorf("Required parameter \"Content\" of \"TranslationParams\" not set")
	}
	if params.KeyId == nil {
		return fmt.Errorf("Required parameter \"KeyId\" of \"TranslationParams\" not set")
	}
	if params.LocaleId == nil {
		return fmt.Errorf("Required parameter \"LocaleId\" of \"TranslationParams\" not set")
	}

	return nil
}

type LocaleFileImportParams struct {
	ConvertEmoji       *bool                   `json:"convert_emoji,omitempty"`
	File               *string                 `json:"file"`
	Format             *string                 `json:"format,omitempty"`
	FormatOptions      *map[string]interface{} `json:"format_options,omitempty"`
	LocaleId           *string                 `json:"locale_id,omitempty"`
	SkipUnverification *bool                   `json:"skip_unverification,omitempty"`
	SkipUploadTags     *bool                   `json:"skip_upload_tags,omitempty"`
	Tags               *[]string               `json:"tags,omitempty"`
	UpdateTranslations *bool                   `json:"update_translations,omitempty"`
}

func (p *LocaleFileImportParams) SetConvertEmoji(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.ConvertEmoji = &b
	return nil
}

func (p *LocaleFileImportParams) SetFile(val string) error {
	p.File = &val
	return nil
}

func (p *LocaleFileImportParams) SetFormat(val string) error {
	p.Format = &val
	return nil
}

func (p *LocaleFileImportParams) SetFormatOptions(val string) error {
	h := map[string]interface{}{}
	err := json.Unmarshal([]byte(val), &h)
	if err != nil {
		return err
	}
	p.FormatOptions = &h
	return nil
}

func (p *LocaleFileImportParams) SetLocaleId(val string) error {
	p.LocaleId = &val
	return nil
}

func (p *LocaleFileImportParams) SetSkipUnverification(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.SkipUnverification = &b
	return nil
}

func (p *LocaleFileImportParams) SetSkipUploadTags(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.SkipUploadTags = &b
	return nil
}

func (p *LocaleFileImportParams) SetTags(val string) error {
	a := strings.Split(val, ",")
	p.Tags = &a
	return nil
}

func (p *LocaleFileImportParams) SetUpdateTranslations(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.UpdateTranslations = &b
	return nil
}

func (params *LocaleFileImportParams) validate() error {
	if params.File == nil {
		return fmt.Errorf("Required parameter \"File\" of \"LocaleFileImportParams\" not set")
	}

	return nil
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

		err = params.validate()
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
func AuthorizationShow(id string) (*AuthorizationFull, error) {
	retVal := new(AuthorizationFull)
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
func AuthorizationUpdate(id string, params *AuthorizationParams) (*AuthorizationFull, error) {
	retVal := new(AuthorizationFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/authorizations/%s", id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func BlacklistKeyCreate(project_id string, params *BlacklistedKeyParams) (*BlacklistedKeyFull, error) {
	retVal := new(BlacklistedKeyFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/blacklisted_keys", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func BlacklistKeyShow(project_id, id string) (*BlacklistedKeyFull, error) {
	retVal := new(BlacklistedKeyFull)
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
func BlacklistKeyUpdate(project_id, id string, params *BlacklistedKeyParams) (*BlacklistedKeyFull, error) {
	retVal := new(BlacklistedKeyFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/blacklisted_keys/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func CommentCreate(project_id, key_id string, params *CommentParams) (*CommentFull, error) {
	retVal := new(CommentFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments", project_id, key_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func CommentShow(project_id, key_id, id string) (*CommentFull, error) {
	retVal := new(CommentFull)
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
func CommentUpdate(project_id, key_id, id string, params *CommentParams) (*CommentFull, error) {
	retVal := new(CommentFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/keys/%s/comments/%s", project_id, key_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func KeyCreate(project_id string, params *TranslationKeyParams) (*TranslationKeyFull, error) {
	retVal := new(TranslationKeyFull)
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
			for i := range *params.Tags {
				err := writer.WriteField("tags[]", (*params.Tags)[i])
				if err != nil {
					return err
				}
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

func (p *KeyListParams) SetLocaleId(val string) error {
	p.LocaleId = &val
	return nil
}

func (p *KeyListParams) SetOrder(val string) error {
	p.Order = &val
	return nil
}

func (p *KeyListParams) SetSort(val string) error {
	p.Sort = &val
	return nil
}

func (p *KeyListParams) SetTranslated(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Translated = &b
	return nil
}

func (params *KeyListParams) validate() error {

	return nil
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

		err = params.validate()
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
func KeyShow(project_id, id string) (*TranslationKeyFull, error) {
	retVal := new(TranslationKeyFull)
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
func KeyUpdate(project_id, id string, params *TranslationKeyParams) (*TranslationKeyFull, error) {
	retVal := new(TranslationKeyFull)
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
			for i := range *params.Tags {
				err := writer.WriteField("tags[]", (*params.Tags)[i])
				if err != nil {
					return err
				}
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
func LocaleCreate(project_id string, params *LocaleParams) (*LocaleFull, error) {
	retVal := new(LocaleFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
	ConvertEmoji             *bool   `json:"convert_emoji,omitempty"`
	Format                   *string `json:"format"`
	IncludeEmptyTranslations *bool   `json:"include_empty_translations,omitempty"`
	KeepNotranslateTags      *bool   `json:"keep_notranslate_tags,omitempty"`
	TagId                    *string `json:"tag_id,omitempty"`
}

func (p *LocaleDownloadParams) SetConvertEmoji(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.ConvertEmoji = &b
	return nil
}

func (p *LocaleDownloadParams) SetFormat(val string) error {
	p.Format = &val
	return nil
}

func (p *LocaleDownloadParams) SetIncludeEmptyTranslations(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.IncludeEmptyTranslations = &b
	return nil
}

func (p *LocaleDownloadParams) SetKeepNotranslateTags(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.KeepNotranslateTags = &b
	return nil
}

func (p *LocaleDownloadParams) SetTagId(val string) error {
	p.TagId = &val
	return nil
}

func (params *LocaleDownloadParams) validate() error {
	if params.Format == nil || *params.Format == "" {
		return fmt.Errorf("Required parameter \"Format\" of \"LocaleDownloadParams\" not set")
	}

	return nil
}

// Download a locale in a specific file format.
func LocaleDownload(project_id, id string, params *LocaleDownloadParams) error {

	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales/%s/download", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
		if err != nil {
			return err
		}

		rc, err := sendRequest("GET", url, "application/json", paramsBuf, 200)
		if err != nil {
			return err
		}
		defer rc.Close()

		return nil
	}()
	return err
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
func LocaleShow(project_id, id string) (*LocaleFull, error) {
	retVal := new(LocaleFull)
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
func LocaleUpdate(project_id, id string, params *LocaleParams) (*LocaleFull, error) {
	retVal := new(LocaleFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/locales/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func OrderConfirm(project_id, id string) (*TranslationOrderFull, error) {
	retVal := new(TranslationOrderFull)
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
func OrderCreate(project_id string, params *TranslationOrderParams) (*TranslationOrderFull, error) {
	retVal := new(TranslationOrderFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/orders", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func ProjectCreate(params *ProjectParams) (*ProjectFull, error) {
	retVal := new(ProjectFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects")

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func ProjectShow(id string) (*ProjectFull, error) {
	retVal := new(ProjectFull)
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
func ProjectUpdate(id string, params *ProjectParams) (*ProjectFull, error) {
	retVal := new(ProjectFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s", id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func StyleguideCreate(project_id string, params *StyleguideParams) (*StyleguideFull, error) {
	retVal := new(StyleguideFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/styleguides", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func StyleguideShow(project_id, id string) (*StyleguideFull, error) {
	retVal := new(StyleguideFull)
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
func StyleguideUpdate(project_id, id string, params *StyleguideParams) (*StyleguideFull, error) {
	retVal := new(StyleguideFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/styleguides/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func TagCreate(project_id string, params *TagParams) (*TagFull, error) {
	retVal := new(TagFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/tags", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func TagShow(project_id, name string) (*TagFull, error) {
	retVal := new(TagFull)
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
func TranslationCreate(project_id string, params *TranslationParams) (*TranslationFull, error) {
	retVal := new(TranslationFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations", project_id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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

func (p *TranslationListAllParams) SetOrder(val string) error {
	p.Order = &val
	return nil
}

func (p *TranslationListAllParams) SetSince(val string) error {
	t, err := time.Parse(time.RubyDate, val)
	p.Since = &t
	return err
}

func (p *TranslationListAllParams) SetSort(val string) error {
	p.Sort = &val
	return nil
}

func (p *TranslationListAllParams) SetUnverified(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Unverified = &b
	return nil
}

func (params *TranslationListAllParams) validate() error {

	return nil
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

		err = params.validate()
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

func (p *TranslationListKeyParams) SetOrder(val string) error {
	p.Order = &val
	return nil
}

func (p *TranslationListKeyParams) SetSince(val string) error {
	t, err := time.Parse(time.RubyDate, val)
	p.Since = &t
	return err
}

func (p *TranslationListKeyParams) SetSort(val string) error {
	p.Sort = &val
	return nil
}

func (p *TranslationListKeyParams) SetUnverified(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Unverified = &b
	return nil
}

func (params *TranslationListKeyParams) validate() error {

	return nil
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

		err = params.validate()
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

func (p *TranslationListLocaleParams) SetOrder(val string) error {
	p.Order = &val
	return nil
}

func (p *TranslationListLocaleParams) SetSince(val string) error {
	t, err := time.Parse(time.RubyDate, val)
	p.Since = &t
	return err
}

func (p *TranslationListLocaleParams) SetSort(val string) error {
	p.Sort = &val
	return nil
}

func (p *TranslationListLocaleParams) SetUnverified(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Unverified = &b
	return nil
}

func (params *TranslationListLocaleParams) validate() error {

	return nil
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

		err = params.validate()
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
func TranslationShow(project_id, id string) (*TranslationFull, error) {
	retVal := new(TranslationFull)
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
	Content      *string `json:"content"`
	Excluded     *bool   `json:"excluded,omitempty"`
	PluralSuffix *string `json:"plural_suffix,omitempty"`
	Unverified   *bool   `json:"unverified,omitempty"`
}

func (p *TranslationUpdateParams) SetContent(val string) error {
	p.Content = &val
	return nil
}

func (p *TranslationUpdateParams) SetExcluded(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Excluded = &b
	return nil
}

func (p *TranslationUpdateParams) SetPluralSuffix(val string) error {
	p.PluralSuffix = &val
	return nil
}

func (p *TranslationUpdateParams) SetUnverified(val string) error {
	var b bool
	if val == "true" {
		b = true
	} else if val == "false" { // ignore
	} else {
		return fmt.Errorf("invalid value %q", val)
	}
	p.Unverified = &b
	return nil
}

func (params *TranslationUpdateParams) validate() error {
	if params.Content == nil || *params.Content == "" {
		return fmt.Errorf("Required parameter \"Content\" of \"TranslationUpdateParams\" not set")
	}

	return nil
}

// Update an existing translation.
func TranslationUpdate(project_id, id string, params *TranslationUpdateParams) (*TranslationFull, error) {
	retVal := new(TranslationFull)
	err := func() error {
		url := fmt.Sprintf("https://api.phraseapp.com/v2/projects/%s/translations/%s", project_id, id)

		paramsBuf := bytes.NewBuffer(nil)
		err := json.NewEncoder(paramsBuf).Encode(&params)
		if err != nil {
			return err
		}

		err = params.validate()
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
func UploadCreate(project_id string, params *LocaleFileImportParams) (*LocaleFileImportFull, error) {
	retVal := new(LocaleFileImportFull)
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
		if params.File != nil {
			part, err := writer.CreateFormFile("file", filepath.Base(*params.File))
			if err != nil {
				return err
			}
			file, err := os.Open(*params.File)
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
		if params.Tags != nil {
			for i := range *params.Tags {
				err := writer.WriteField("tags[]", (*params.Tags)[i])
				if err != nil {
					return err
				}
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
func UploadShow(project_id, id string) (*LocaleFileImportFull, error) {
	retVal := new(LocaleFileImportFull)
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
func VersionShow(project_id, translation_id, id string) (*TranslationVersionFull, error) {
	retVal := new(TranslationVersionFull)
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
