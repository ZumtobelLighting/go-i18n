package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// UnmarshalFunc unmarshals data into v.
type UnmarshalFunc func(data []byte, v interface{}) error

// Bundle stores all messages and pluralization rules.
// Generally, your application should only need a single bundle
// that is initialized early in your application's lifecycle.
type Bundle struct {
	// MessageTemplates maps language tags to language ids to message templates.
	MessageTemplates map[string]map[string]*MessageTemplate

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

// LoadMessageFile loads the bytes from path
// and then calls ParseMessageFileBytes.
func (b *Bundle) LoadMessageFile(path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return b.ParseMessageFileBytes(buf, path)
}

// MustLoadMessageFile is similar to LoadTranslationFile
// except it panics if an error happens.
func (b *Bundle) MustLoadMessageFile(path string) {
	if err := b.LoadMessageFile(path); err != nil {
		panic(err)
	}
}

// ParseMessageFileBytes parses the bytes in buf to add translations to the bundle.
// It is useful for parsing translation files embedded with go-bindata.
//
// The format of the file is everything after the last ".".
//
// The language tag of path is the last match of LanguageTagRegex, excluding everything after the last ".".
func (b *Bundle) ParseMessageFileBytes(buf []byte, path string) error {
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

	var messageTemplates []*MessageTemplate
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
		t, err := NewMessageTemplate(id, strdata)
		if err != nil {
			return err
		}
		messageTemplates = append(messageTemplates, t)
	}
	pathNoFormat := path[:len(path)-len(format)]
	langTags := LanguageTagRegex.FindAllString(pathNoFormat, -1)
	if len(langTags) == 0 {
		return fmt.Errorf("no language tag found in path: %s", path)
	}
	langTag := langTags[len(langTags)-1]
	return b.AddMessageTemplates(langTag, messageTemplates...)
}

func parseFormat(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}

// MustParseMessageFileBytes is similar to ParseMessageFileBytes
// except it panics if an error happens.
func (b *Bundle) MustParseMessageFileBytes(buf []byte, path string) {
	if err := b.ParseMessageFileBytes(buf, path); err != nil {
		panic(err)
	}
}

// AddMessageTemplates adds message templates for a language.
// It is useful if your messages are in a format not supported by ParseMessageFileBytes.
func (b *Bundle) AddMessageTemplates(langTag string, messageTemplates ...*MessageTemplate) error {
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
	if b.MessageTemplates == nil {
		b.MessageTemplates = make(map[string]map[string]*MessageTemplate)
	}
	if b.MessageTemplates[langTag] == nil {
		b.MessageTemplates[langTag] = make(map[string]*MessageTemplate)
	}
	for _, t := range messageTemplates {
		b.MessageTemplates[langTag][t.ID] = t
	}
	return nil
}
