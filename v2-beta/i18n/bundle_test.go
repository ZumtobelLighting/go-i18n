package i18n_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2-beta/i18n"
	yaml "gopkg.in/yaml.v2"
)

var simpleMessageTemplate = i18n.MustNewMessageTemplate("simple", map[string]string{
	"other": "simple translation",
})

var detailMessageTemplate = i18n.MustNewMessageTemplate("detail", map[string]string{
	"description": "detail description",
	"other":       "detail translation",
})

var everythingMessageTemplate = i18n.MustNewMessageTemplate("everything", map[string]string{
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
	bundle.MustParseMessageFileBytes([]byte(`{
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

	expectMessageTemplate(t, bundle, "en-US", "simple", simpleMessageTemplate)
	expectMessageTemplate(t, bundle, "en-US", "detail", detailMessageTemplate)
	expectMessageTemplate(t, bundle, "en-US", "everything", everythingMessageTemplate)
}

func TestYAML(t *testing.T) {
	var bundle i18n.Bundle
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
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

	expectMessageTemplate(t, bundle, "en-US", "simple", simpleMessageTemplate)
	expectMessageTemplate(t, bundle, "en-US", "detail", detailMessageTemplate)
	expectMessageTemplate(t, bundle, "en-US", "everything", everythingMessageTemplate)
}

func TestTOML(t *testing.T) {
	var bundle i18n.Bundle
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
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

	expectMessageTemplate(t, bundle, "en-US", "simple", simpleMessageTemplate)
	expectMessageTemplate(t, bundle, "en-US", "detail", detailMessageTemplate)
	expectMessageTemplate(t, bundle, "en-US", "everything", everythingMessageTemplate)
}

func expectMessageTemplate(t *testing.T, bundle i18n.Bundle, langTag, messageID string, expected *i18n.MessageTemplate) {
	actual := bundle.MessageTemplates[langTag][messageID]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("bundle.MessageTemplates[%q][%q] = %#v; want %#v", langTag, messageID, actual, expected)
	}
}
