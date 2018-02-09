package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// UnmarshalFunc unmarshals data into v.
type UnmarshalFunc func(data []byte, v interface{}) error

// Bundle stores all translations and pluralization rules.
// Generally, your application should only need a single bundle
// that is initialized early in your application's lifecycle.
type Bundle struct {
	// Translations maps language tags to language ids to translations.
	Translations map[string]map[string]*Translation

	// PluralRules maps language tags to their plural rules.
	PluralRules map[string]*PluralRule

	// UnmarshalFuncs maps file formats to unmarshal functions.
	UnmarshalFuncs map[string]UnmarshalFunc
}

// NewBundle returns a new bundle that contains the
// CLDR plural rules and a json unmarshaler.
func NewBundle() *Bundle {
	return &Bundle{
		PluralRules: CLDRPluralRules(),
		UnmarshalFuncs: map[string]UnmarshalFunc{
			"json": json.Unmarshal,
		},
	}
}

// RegisterUnmarshalFunc registers an UnmarshalFunc for format.
func (b *Bundle) RegisterUnmarshalFunc(format string, unmarshalFunc UnmarshalFunc) {
	if b.UnmarshalFuncs == nil {
		b.UnmarshalFuncs = make(map[string]UnmarshalFunc)
	}
	b.UnmarshalFuncs[format] = unmarshalFunc
}

// LoadTranslationFile loads the bytes from path
// and then calls ParseTranslationFileBytes.
func (b *Bundle) LoadTranslationFile(path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return b.ParseTranslationFileBytes(buf, path)
}

// MustLoadTranslationFile is similar to LoadTranslationFile
// except it panics if an error happens.
func (b *Bundle) MustLoadTranslationFile(path string) {
	if err := b.LoadTranslationFile(path); err != nil {
		panic(err)
	}
}

// ParseTranslationFileBytes parses the bytes in buf to add translations to the bundle.
// It is useful for parsing translation files embedded with go-bindata.
//
// The format of the file is everything after the first ".", or the whole path if there is no ".".
//
// The language tag of path is the last match of LanguageTagRegex.
func (b *Bundle) ParseTranslationFileBytes(buf []byte, path string) error {
	if len(buf) == 0 {
		return nil
	}

	format := parseFormat(path)
	unmarshalFunc := b.UnmarshalFuncs[format]
	if unmarshalFunc == nil {
		return fmt.Errorf("no unmarshaler registered for %s", format)
	}

	var raw map[string]interface{}
	if err := unmarshalFunc(buf, &raw); err != nil {
		return err
	}

	var translations []*Translation
	for id, data := range raw {
		strdata := make(map[string]string)
		switch value := data.(type) {
		case string:
			strdata["other"] = value
		case map[string]interface{}:
			for k, v := range value {
				vstr, ok := v.(string)
				if !ok {
					return fmt.Errorf("expected [%s][%s][%s] to be a string but got %#v", path, id, k, v)
				}
				strdata[k] = vstr
			}
		case map[interface{}]interface{}:
			for k, v := range value {
				kstr, ok := k.(string)
				if !ok {
					return fmt.Errorf("[%s][%s] has a non-string key %#v", path, id, k)
				}
				vstr, ok := v.(string)
				if !ok {
					return fmt.Errorf("[%s][%s][%s] has a non-string value %#v", path, id, k, v)
				}
				strdata[kstr] = vstr
			}
		default:
			return fmt.Errorf("translation key %s in %s has invalid value: %#v", id, path, value)
		}
		t, err := NewTranslation(id, strdata)
		if err != nil {
			return err
		}
		translations = append(translations, t)
	}
	langTags := LanguageTagRegex.FindAllString(path, -1)
	langTag := langTags[len(langTags)-1]
	return b.AddTranslations(langTag, translations...)
}

func parseFormat(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return path
}

// MustParseTranslationFileBytes is similar to ParseTranslationFileBytes
// except it panics if an error happens.
func (b *Bundle) MustParseTranslationFileBytes(buf []byte, path string) {
	if err := b.ParseTranslationFileBytes(buf, path); err != nil {
		panic(err)
	}
}

// AddTranslations adds translations for a language.
// It is useful if your translations are in a format not supported by ParseTranslationFileBytes.
func (b *Bundle) AddTranslations(langTag string, translations ...*Translation) error {
	if b.PluralRules == nil {
		b.PluralRules = CLDRPluralRules()
	}
	pluralID := langTag
	for i, r := range langTag {
		if r == '-' || r == '_' {
			pluralID = langTag[:i]
			break
		}
	}
	pluralRule := b.PluralRules[pluralID]
	if pluralRule == nil {
		return fmt.Errorf("no plural rule registered for %s", pluralID)
	}
	b.PluralRules[langTag] = pluralRule
	if b.Translations == nil {
		b.Translations = make(map[string]map[string]*Translation)
	}
	if b.Translations[langTag] == nil {
		b.Translations[langTag] = make(map[string]*Translation)
	}
	for _, t := range translations {
		b.Translations[langTag][t.ID] = t
	}
	return nil
}
