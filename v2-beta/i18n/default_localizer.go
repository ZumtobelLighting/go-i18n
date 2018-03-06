package i18n

import "fmt"

// DefaultLocalizer provides Localize and MustLocalize methods that return localized messages.
//
// Localize and MustLocalize accept default messages that are used if no message
// is found in the bundle. These default messages can be extracted from Go source code
// using the goi18n command.
type DefaultLocalizer struct {
	*Localizer

	// DefaultLangaugeTag is the language tag of the plural rule
	// that is used when pluralizing default messages passed to Localize and MustLocalize.
	DefaultLanguageTag string
}

// NewDefaultLocalizer returns a new DefaultLocalizer that looks up messages
// in the bundle according to the order of language tags found in prefs.
//
// It can parse languages from Accept-Language headers (RFC 2616),
// but it assumes weights are monotonically decreasing.
//
// Default messages are pluralized according to the rule for defaultLanguageTag.
func NewDefaultLocalizer(bundle *Bundle, prefs, defaultLanguageTag string) *DefaultLocalizer {
	l := NewLocalizer(bundle, prefs)
	return &DefaultLocalizer{
		Localizer:          l,
		DefaultLanguageTag: defaultLanguageTag,
	}
}

// Localize returns the localized message.
// If no message is found in the bundle, it returns the provided default message.
//
// Valid calls to Localize have one of the following signatures:
//     Localize(id string, default string, templateData interface{})
//     Localize(message *Message, templateData interface{})
//     Localize(message *Message, pluralCount interface{})
//     Localize(message *Message, pluralCount, templateData interface{})
func (l *DefaultLocalizer) Localize(idOrMessage interface{}, args ...interface{}) (string, error) {
	newInvocationError := func(reason string) error {
		return &invocationError{
			name:   "Localize",
			args:   append([]interface{}{idOrMessage}, args...),
			reason: reason,
		}
	}
	switch v := idOrMessage.(type) {
	case string:
		defaultMessage, ok := indexOrNil(args, 0).(string)
		if !ok {
			return "", newInvocationError("expected second parameter to be a string")
		}
		templateData := indexOrNil(args, 1)
		return l.localize(&Message{ID: v, Other: defaultMessage}, 0, templateData)
	case *Message:
		pluralCount, templateData := parseArgs(args)
		return l.localize(v, pluralCount, templateData)
	default:
		return "", newInvocationError("expected first argument to be of type string, *SingleMessage, or *PluralMessage")
	}
}

func (l *DefaultLocalizer) localize(message *Message, pluralCount interface{}, templateData interface{}) (string, error) {
	operands, _ := newOperands(pluralCount)
	localized, err := l.Localizer.localize(message.ID, operands, templateData)
	if err != nil {
		return "", err
	}
	if localized != "" {
		return localized, nil
	}
	template, err := message.Template()
	if err != nil {
		return "", err
	}
	pluralForm := l.Localizer.pluralForm(l.DefaultLanguageTag, operands)
	if pluralForm == Invalid {
		return "", fmt.Errorf("unable to pluralize %q because there no plural rule for %q", message.ID, l.DefaultLanguageTag)
	}
	return template.Execute(pluralForm, templateData), nil
}

// MustLocalize is similar to Localize, except it panics if an error happens.
func (l *DefaultLocalizer) MustLocalize(idOrMessage interface{}, args ...interface{}) string {
	localized, err := l.Localize(idOrMessage, args...)
	if err != nil {
		if invocationErr, ok := err.(*invocationError); ok {
			invocationErr.name = "MustLocalize"
		}
		panic(err)
	}
	return localized
}

func indexOrNil(slice []interface{}, idx int) interface{} {
	if len(slice) <= idx {
		return nil
	}
	return slice[idx]
}
