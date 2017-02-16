# phraseapp-go

Go library for PhraseApp API v2.

# Start using it

1. Download and install \
`go get github.com/phrase/phraseapp-go/phraseapp`

2. Import it in your code \
`import "github.com/phrase/phraseapp-go/phraseapp"`

# API examples
### Init client
```
credentials := phraseapp.Credentials{
	Host:  "https://api.phraseapp.com",
	Token: "access_token",
}
client := phraseapp.Client{
	Credentials: credentials,
}
```

### Create project
```
projectName := "project_name"
sharesTranslationMemory := true
projectParams := phraseapp.ProjectParams{
	Name: &projectName,
	SharesTranslationMemory: &sharesTranslationMemory,
}
project, err := client.ProjectCreate(&projectParams)
```

### Create locale
```
localeCode := "en-GB"
localeDetails := phraseapp.LocaleParams{
	Name: &localeCode,
	Code: &localeCode,
}
locale, err := client.LocaleCreate("project_id", &localeDetails)
```

### Create key
```
keyName := "key_name"
tags := "tag1, tag2"
keyParams := phraseapp.TranslationKeyParams{
	Name: &keyName,
	Tags: &tags,
}
key, err := client.KeyCreate("project_id", &keyParams)
```

### Create translation
```
keyID := "key_id"
localeID := "locale_id"
content := "my_content"
keyID := "key_id"
translationParams := phraseapp.TranslationParams{
    KeyID:    &keyID,
    LocaleID: &localeID,
    Content:  &content,
}
translation, err := client.TranslationCreate("project_id", &translationParams)
```

### Upload translation file
```
fileName := "file.json"
fileFormat := "simple_json"
updateTranslations := true
uploadParams := phraseapp.UploadParams{
	File:               &fileName,
	LocaleID:           &localeID,
	FileFormat:         &fileFormat,
	UpdateTranslations: &updateTranslations,
}
upload, err := client.UploadCreate("project_id", &uploadParams)
```

### Download locale as a file
```
fileFormat := "simple_json"
localeDownloadParams := phraseapp.LocaleDownloadParams{
	FileFormat: &fileFormat,
}
var localeData []byte
localeData, err := client.LocaleDownload("project_id", "locale_id", &localeDownloadParams)
ioutil.WriteFile("en.json", localeData, 0644)
```

### Query translations
```
translationsQuery := "tags:tag1,tag2"
translationSearchParams := phraseapp.TranslationsSearchParams{
	Q: &translationsQuery,
}
translations, err := client.TranslationsSearch("project_id", 1, 1000, &translationSearchParams)
```
More [query options](https://phraseapp.com/docs/api/v2/translations/#search)


For a more complete example the wiki contains an example how to [upload files as translations](https://github.com/phrase/phraseapp-go/wiki/Sync-local-files-to-PhraseApp) to PhraseApp.

## Contributing

This library is auto-generated from templates that run against a API specification file. Therefore we can not accept any pull requests in this repository. Please use the GitHub Issue Tracker to report bugs.
