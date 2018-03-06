package i18n

import "fmt"

// DefaultTranslator provides Translate and MustTranslate methods to translate messages.
//
// Translate and MustTranslate accept default translations that are used if no translation
// is found in the bundle. These default translations can be extracted from Go source code
// using the goi18n command.
type DefaultTranslator struct {
	*Translator

	// DefaultLangaugeTag is the language of default translations passed
	// to its Translate and MustTranslate methods.
	// It is used to determine how to pluralize default translations.
	DefaultLanguageTag string
}

// NewDefaultTranslator returns a new DefaultTranslator that looks up translations
// in the bundle according to the order of language tags found in prefs.
//
// It can parse languages from Accept-Language headers (RFC 2616),
// but it assumes weights are monotonically decreasing.
//
// Default translations are pluralized according to the rules for defaultLanguageTag.
func NewDefaultTranslator(bundle *Bundle, prefs, defaultLanguageTag string) *DefaultTranslator {
	t := NewTranslator(bundle, prefs)
	return &DefaultTranslator{
		Translator:         t,
		DefaultLanguageTag: defaultLanguageTag,
	}
}

// Translate returns the translation for the message.
// If no translation is found, it returns the default translation provided.
//
// Valid invocations have one of the following signatures:
//     Translate(id string, default string, templateData interface{})
//     Translate(singleMessage *SingleMessage, templateData interface{})
//     Translate(pluralMessage *PluralMessage, pluralCount interface{})
//     Translate(pluralMessage *PluralMessage, pluralCount, templateData interface{})
func (t *DefaultTranslator) Translate(idOrMessage interface{}, args ...interface{}) (string, error) {
	newInvocationError := func(reason string) error {
		return &invocationError{
			name:   "Translate",
			args:   append([]interface{}{idOrMessage}, args...),
			reason: reason,
		}
	}
	switch v := idOrMessage.(type) {
	case string:
		defaultTranslation, ok := indexOrNil(args, 0).(string)
		if !ok {
			return "", newInvocationError("expected second parameter to be a string")
		}
		templateData := indexOrNil(args, 1)
		return t.translateMessage(&SingleMessage{ID: v, Content: defaultTranslation}, 0, templateData)
	case *SingleMessage:
		templateData := indexOrNil(args, 0)
		return t.translateMessage(v, 0, templateData)
	case *PluralMessage:
		if argc := len(args); argc != 1 && argc != 2 {
			return "", newInvocationError("expected two three arguments")
		}
		pluralCount, templateData := parseArgs(args)
		return t.translateMessage(v, pluralCount, templateData)
	default:
		return "", newInvocationError("expected first argument to be of type string, *SingleMessage, or *PluralMessage")
	}
}

func (t *DefaultTranslator) translateMessage(message Message, pluralCount interface{}, templateData interface{}) (string, error) {
	operands, _ := newOperands(pluralCount)
	translated, err := t.Translator.translate(message.MessageID(), operands, templateData)
	if err != nil {
		return "", err
	}
	if translated != "" {
		return translated, nil
	}
	translation, err := message.Translation()
	if err != nil {
		return "", err
	}
	pluralForm := t.Translator.pluralForm(t.DefaultLanguageTag, operands)
	if pluralForm == Invalid {
		return "", fmt.Errorf("unable to pluralize %q because there no plural rule for %q", message.MessageID(), t.DefaultLanguageTag)
	}
	return translation.Translate(pluralForm, templateData), nil
}

// MustTranslate is similar to Translate, except it panics if an error happens.
func (t *DefaultTranslator) MustTranslate(idOrMessage interface{}, args ...interface{}) string {
	translated, err := t.Translate(idOrMessage, args...)
	if err != nil {
		if invocationErr, ok := err.(*invocationError); ok {
			invocationErr.name = "MustTranslate"
		}
		panic(err)
	}
	return translated
}

func indexOrNil(slice []interface{}, idx int) interface{} {
	if len(slice) <= idx {
		return nil
	}
	return slice[idx]
}

// TODO: unexport
type Message interface {
	MessageID() string
	Translation() (*Translation, error)
}

type SingleMessage struct {
	ID          string
	Description string
	Content     string
}

func (m *SingleMessage) MessageID() string {
	return m.ID
}

func (m *SingleMessage) Translation() (*Translation, error) {
	return NewTranslation(m.ID, map[string]string{
		"description": m.Description,
		"other":       m.Content,
	})
}

var _ = Message(&SingleMessage{})

type PluralMessage struct {
	ID          string
	Description string
	Zero        string
	One         string
	Two         string
	Few         string
	Many        string
	Other       string
}

func (m *PluralMessage) MessageID() string {
	return m.ID
}

func (m *PluralMessage) Translation() (*Translation, error) {
	return NewTranslation(m.ID, map[string]string{
		"description": m.Description,
		"zero":        m.Zero,
		"one":         m.One,
		"two":         m.Two,
		"few":         m.Few,
		"many":        m.Many,
		"other":       m.Other,
	})
}

var _ = Message(&PluralMessage{})
