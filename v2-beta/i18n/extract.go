package i18n

import (
	"fmt"
	"go/parser"
	"go/token"
)

type TranslatorExperiment struct {
	Bundle       *Bundle
	LanguageTags []string

	// DefaultLangaugeTag is the language that default translations are in.
	DefaultLanguageTag string
}

// func (t *TranslatorExperiment) Translate(id string) *Message {
// 	return &Message{ID: id}
// }

// type Message struct {
// 	ID                 string
// 	DefaultTranslation string
// 	TemplateData       interface{}
// }

// func (m *Message) Default(translation string) *Message {
// 	m.DefaultTranslation = translation
// 	return m
// }

// func (m *Message) Data(data interface{}) *Message {
// 	m.TemplateData = data
// 	return m
// }

// func (m *Message) String() string {
// 	return m.DefaultTranslation
// }

// ExtractTranslations extracts translations from the bytes of a Go source file.
func ExtractTranslations(buf []byte) ([]*Translation, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", buf, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v\n", file)

	translations := []*Translation{}
	return translations, nil
}
