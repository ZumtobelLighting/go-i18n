package i18n_test

import (
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/v2-beta/i18n"
)

func TestNewTranslator(t *testing.T) {
	tests := []struct {
		prefs string
		tags  []string
	}{
		{
			prefs: "aa-bb-cc-dd",
			tags:  []string{"aa-bb-cc-dd", "aa-bb-cc", "aa-bb", "aa"},
		},
		{
			prefs: "-aa-bb-cc-dd-",
			tags:  []string{"aa-bb-cc-dd", "aa-bb-cc", "aa-bb", "aa"},
		},
		{
			prefs: "x-aa-bb-cc-dd-x",
			tags:  []string{"aa-bb-cc-dd", "aa-bb-cc", "aa-bb", "aa"},
		},
		{
			prefs: "aa-bbb-cccc-ddddd",
			tags:  []string{"aa-bbb-cccc-ddddd", "aa-bbb-cccc", "aa-bbb", "aa"},
		},
		{
			prefs: "aa-bb, aa-cc",
			tags:  []string{"aa-bb", "aa", "aa-cc"},
		},
	}
	for _, test := range tests {
		translator := i18n.NewTranslator(nil, test.prefs)
		if !reflect.DeepEqual(translator.LanguageTags, test.tags) {
			t.Errorf("i18n.NewTranslator(nil, %q)\ngot  %#v\nwant %#v", test.prefs, translator.LanguageTags, test.tags)
		}
	}
}
