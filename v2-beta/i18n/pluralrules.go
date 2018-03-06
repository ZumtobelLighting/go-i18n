package i18n

import "strings"

// PluralRule defines the CLDR plural rules for a language.
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
// http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type PluralRule struct {
	PluralForms    map[PluralForm]struct{}
	PluralFormFunc func(*Operands) PluralForm
}

func normalizePluralRuleID(id string) string {
	id = strings.Replace(id, "_", "-", -1)
	id = strings.ToLower(id)
	return id
}

func addPluralRules(rules map[string]*PluralRule, ids []string, ps *PluralRule) {
	for _, id := range ids {
		id = normalizePluralRuleID(id)
		rules[id] = ps
	}
}

func newPluralFormSet(pluralForms ...PluralForm) map[PluralForm]struct{} {
	set := make(map[PluralForm]struct{}, len(pluralForms))
	for _, plural := range pluralForms {
		set[plural] = struct{}{}
	}
	return set
}

func intInRange(i, from, to int64) bool {
	return from <= i && i <= to
}

func intEqualsAny(i int64, any ...int64) bool {
	for _, a := range any {
		if i == a {
			return true
		}
	}
	return false
}
