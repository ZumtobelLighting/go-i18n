package i18n

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

// LanguageTagRegex matches language tags like en, en-US, and zh-Hans-CN.
// Language tags are case-insensitive.
var LanguageTagRegex = regexp.MustCompile(`[a-zA-Z]{2,}([\-_][a-zA-Z]{2,})*`)

// Localizer provides Localize and MustLocalize methods that return localized messages.
//
// Use DefaultLocalizer to define default messages in Go source code that
// can be extracted using the goi18n command.
type Localizer struct {

	// Bundle contains the messages that can be returned by the Localizer.
	Bundle *Bundle

	// LanguageTags is the list of language tags that the Localizer checks
	// in order when localizing a message.
	LanguageTags []string
}

// NewLocalizer returns a new Localizer that looks up messages
// in the bundle according to the order of language tags found in prefs.
//
// It can parse languages from Accept-Language headers (RFC 2616),
// but it assumes weights are monotonically decreasing.
func NewLocalizer(bundle *Bundle, prefs string) *Localizer {
	localizer := &Localizer{
		Bundle:       bundle,
		LanguageTags: []string{},
	}

	langTags := LanguageTagRegex.FindAllString(prefs, -1)
	var tags []string
	for _, langTag := range langTags {
		tags = append(tags, expandTag(langTag)...)
	}
	localizer.LanguageTags = dedupe(tags)
	return localizer
}

func expandTag(tag string) []string {
	tag = strings.TrimSpace(tag)
	tag = strings.ToLower(tag)
	tags := []string{tag}
	for i := len(tag) - 1; i > 0; {
		r, size := utf8.DecodeLastRuneInString(tag[:i])
		i -= size
		if r == '-' || r == '_' {
			tags = append(tags, tag[:i])
		}
	}
	return tags
}

func dedupe(strs []string) []string {
	found := make(map[string]struct{}, len(strs))
	deduped := make([]string, 0, len(strs))
	for _, str := range strs {
		if _, ok := found[str]; !ok {
			found[str] = struct{}{}
			deduped = append(deduped, str)
		}
	}
	return deduped
}

// Localize returns the localized message.
// If no message is found in the bundle, it returns the empty string.
//
// Valid calls to Localize have one of the following signatures:
//     Localize(id string)
//     Localize(id string, templateData interface{})
//     Localize(id string, pluralCount interface{})
//     Localize(id string, pluralCount, templateData interface{})
func (l *Localizer) Localize(id string, args ...interface{}) (string, error) {
	pluralCount, templateData := parseArgs(args)
	operands, _ := newOperands(pluralCount)
	return l.localize(id, operands, templateData)
}

func (l *Localizer) localize(id string, operands *Operands, templateData interface{}) (string, error) {
	for _, langTag := range l.LanguageTags {
		templates := l.Bundle.MessageTemplates[langTag]
		if templates == nil {
			continue
		}
		template := templates[id]
		if template == nil {
			continue
		}
		pluralForm := l.pluralForm(langTag, operands)
		if pluralForm == Invalid {
			return "", fmt.Errorf("unable to pluralize %q because there no plural rule for %q", id, langTag)
		}
		if localized := template.Execute(pluralForm, templateData); localized != "" {
			return localized, nil
		}
	}
	return "", nil
}

func (l *Localizer) pluralForm(langTag string, operands *Operands) PluralForm {
	if operands == nil {
		return Other
	}
	pluralRule := l.Bundle.PluralRules[langTag]
	if pluralRule == nil {
		return Invalid
	}
	return pluralRule.PluralFormFunc(operands)
}

// MustLocalize is similar to Localize, except it panics if an error happens.
func (l *Localizer) MustLocalize(id string, args ...interface{}) string {
	localized, err := l.Localize(id, args...)
	if err != nil {
		panic(err)
	}
	return localized
}

func parseArgs(args []interface{}) (count interface{}, data interface{}) {
	if argc := len(args); argc > 0 {
		if isNumber(args[0]) {
			count = args[0]
			if argc > 1 {
				data = args[1]
			}
		} else {
			data = args[0]
		}
	}

	if count != nil {
		if data == nil {
			data = map[string]interface{}{"Count": count}
		} else {
			dataMap := toMap(data)
			dataMap["Count"] = count
			data = dataMap
		}
	} else {
		dataMap := toMap(data)
		if c, ok := dataMap["Count"]; ok {
			count = c
		}
	}
	return
}

func isNumber(n interface{}) bool {
	switch n.(type) {
	case int, int8, int16, int32, int64, string:
		return true
	}
	return false
}

func toMap(input interface{}) map[string]interface{} {
	if data, ok := input.(map[string]interface{}); ok {
		return data
	}
	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Ptr:
		return toMap(v.Elem().Interface())
	case reflect.Struct:
		return structToMap(v)
	default:
		return nil
	}
}

// Converts the top level of a struct to a map[string]interface{}.
// Code inspired by github.com/fatih/structs.
func structToMap(v reflect.Value) map[string]interface{} {
	out := make(map[string]interface{})
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			// skip unexported field
			continue
		}
		out[field.Name] = v.FieldByName(field.Name).Interface()
	}
	return out
}
