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
			ops, err := newOperands(test.num)
			if err != nil {
				t.Errorf("%s: newOperands(%d) errored with %s", pluralRuleID, test.num, err)
				break
			}
			if pluralForm := rule.PluralFormFunc(ops); pluralForm != test.pluralForm {
				t.Errorf("%s: PluralFormFunc(%#v) returned %q, %v; expected %q", pluralRuleID, ops, pluralForm, err, test.pluralForm)
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