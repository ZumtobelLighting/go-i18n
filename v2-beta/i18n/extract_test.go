package i18n_test

import (
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/v2-beta/i18n"
)

func TestExtract2(t *testing.T) {

	tests := []struct {
		name         string
		file         string
		translations []*i18n.MessageTemplate
		err          error
	}{
		{
			name:         "no translations",
			file:         `package main`,
			translations: []*i18n.MessageTemplate{},
			err:          nil,
		},
		{
			name: "exhaustive plural translation",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2-beta/i18n"

			func main() {
				_ := &i18n.Message{
					ID:          "Plural ID",
					Description: "Plural description",
					Zero:        "Zero translation",
					One:         "One translation",
					Two:         "Two translation",
					Few:         "Few translation",
					Many:        "Many translation",
					Other:       "Other translation",
				}
			}
			`,
			translations: []*i18n.MessageTemplate{},
			err:          nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			translations, err := i18n.ExtractTranslations([]byte(test.file))
			if err != test.err {
				t.Errorf("i18n.ExtractTranslations(%q)\nexpected error: %q\n     got error: %q", test.file, test.err, err)
			}
			if !reflect.DeepEqual(translations, test.translations) {
				t.Errorf("i18n.ExtractTranslations(%q)\nexpected: %#v\n     got: %#v", test.file, test.translations, translations)
			}
		})
	}
}
