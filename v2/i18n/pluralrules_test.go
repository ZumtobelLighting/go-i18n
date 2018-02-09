package i18n

import (
	"strconv"
	"strings"
	"testing"
)

type pluralFormTest struct {
	num        interface{}
	pluralForm PluralForm
}

func runTests(t *testing.T, pluralRuleID string, tests []pluralFormTest) {
	pluralRuleID = normalizePluralRuleID(pluralRuleID)
	pluralRules := CLDRPluralRules()
	if rule := pluralRules[pluralRuleID]; rule != nil {
		for _, test := range tests {
			if plural, err := rule.PluralForm(test.num); plural != test.pluralForm {
				t.Errorf("%s: PluralCategory(%#v) returned %s, %v; expected %s", pluralRuleID, test.num, plural, err, test.pluralForm)
			}
		}
	} else {
		t.Errorf("could not find plural rule for locale %s", pluralRuleID)
	}

}

func appendIntegerTests(tests []pluralFormTest, plural PluralForm, examples []string) []pluralFormTest {
	for _, ex := range expandExamples(examples) {
		i, err := strconv.ParseInt(ex, 10, 64)
		if err != nil {
			panic(err)
		}
		tests = append(tests, pluralFormTest{ex, plural}, pluralFormTest{i, plural})
	}
	return tests
}

func appendDecimalTests(tests []pluralFormTest, plural PluralForm, examples []string) []pluralFormTest {
	for _, ex := range expandExamples(examples) {
		tests = append(tests, pluralFormTest{ex, plural})
	}
	return tests
}

func expandExamples(examples []string) []string {
	var expanded []string
	for _, ex := range examples {
		if parts := strings.Split(ex, "~"); len(parts) == 2 {
			for ex := parts[0]; ; ex = increment(ex) {
				expanded = append(expanded, ex)
				if ex == parts[1] {
					break
				}
			}
		} else {
			expanded = append(expanded, ex)
		}
	}
	return expanded
}

func increment(dec string) string {
	runes := []rune(dec)
	carry := true
	for i := len(runes) - 1; carry && i >= 0; i-- {
		switch runes[i] {
		case '.':
			continue
		case '9':
			runes[i] = '0'
		default:
			runes[i]++
			carry = false
		}
	}
	if carry {
		runes = append([]rune{'1'}, runes...)
	}
	return string(runes)
}
