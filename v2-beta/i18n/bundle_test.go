package i18n_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2-beta/i18n"
	yaml "gopkg.in/yaml.v2"
)

var simpleTranslation = i18n.MustNewTranslation("simple", map[string]string{
	"other": "simple translation",
})

var detailTranslation = i18n.MustNewTranslation("detail", map[string]string{
	"description": "detail description",
	"other":       "detail translation",
})

var everythingTranslation = i18n.MustNewTranslation("everything", map[string]string{
	"description": "everything description",
	"zero":        "zero translation",
	"one":         "one translation",
	"two":         "two translation",
	"few":         "few translation",
	"many":        "many translation",
	"other":       "other translation",
})

func TestJSON(t *testing.T) {
	var bundle i18n.Bundle
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`{
	"simple": "simple translation",
	"detail": {
		"description": "detail description",
		"other": "detail translation"
	},
	"everything": {
		"description": "everything description",
		"zero": "zero translation",
		"one": "one translation",
		"two": "two translation",
		"few": "few translation",
		"many": "many translation",
		"other": "other translation"
	}
}`), "en-US.json")

	expectTranslation(t, bundle, "en-US", "simple", simpleTranslation)
	expectTranslation(t, bundle, "en-US", "detail", detailTranslation)
	expectTranslation(t, bundle, "en-US", "everything", everythingTranslation)
}

func TestYAML(t *testing.T) {
	var bundle i18n.Bundle
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`
# Comment
simple: simple translation

# Comment
detail:
  description: detail description 
  other: detail translation

# Comment
everything:
  description: everything description
  zero: zero translation
  one: one translation
  two: two translation
  few: few translation
  many: many translation
  other: other translation
`), "en-US.yaml")

	expectTranslation(t, bundle, "en-US", "simple", simpleTranslation)
	expectTranslation(t, bundle, "en-US", "detail", detailTranslation)
	expectTranslation(t, bundle, "en-US", "everything", everythingTranslation)
}

func TestTOML(t *testing.T) {
	var bundle i18n.Bundle
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`
# Comment
simple = "simple translation"

# Comment
[detail]
description = "detail description"
other = "detail translation"

# Comment
[everything]
description = "everything description"
zero = "zero translation"
one = "one translation"
two = "two translation"
few = "few translation"
many = "many translation"
other = "other translation"
`), "en-US.toml")

	expectTranslation(t, bundle, "en-US", "simple", simpleTranslation)
	expectTranslation(t, bundle, "en-US", "detail", detailTranslation)
	expectTranslation(t, bundle, "en-US", "everything", everythingTranslation)
}

func expectTranslation(t *testing.T, bundle i18n.Bundle, langTag, translationID string, expected *i18n.Translation) {
	actual := bundle.Translations[langTag][translationID]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("bundle.Translations[%q][%q] = %#v; want %#v", langTag, translationID, actual, expected)
	}
}
