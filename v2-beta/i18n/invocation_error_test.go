package i18n

import "testing"

func TestInvocationError(t *testing.T) {
	testCases := []struct {
		invocationError *invocationError
		err             string
	}{
		{
			invocationError: &invocationError{
				name:   "Foo",
				args:   []interface{}{},
				reason: "reason",
			},
			err: `invalid invocation Foo(); reason`,
		},
		{
			invocationError: &invocationError{
				name:   "Foo",
				args:   []interface{}{"arg1"},
				reason: "reason",
			},
			err: `invalid invocation Foo("arg1"); reason`,
		},
		{
			invocationError: &invocationError{
				name:   "Foo",
				args:   []interface{}{"arg1", "arg2"},
				reason: "reason",
			},
			err: `invalid invocation Foo("arg1", "arg2"); reason`,
		},
		{
			invocationError: &invocationError{
				name:   "Foo",
				args:   []interface{}{"arg1", "arg2", "arg3"},
				reason: "reason",
			},
			err: `invalid invocation Foo("arg1", "arg2", "arg3"); reason`,
		},
	}
	for _, testCase := range testCases {
		if actual := testCase.invocationError.Error(); actual != testCase.err {
			t.Errorf("\nwant %q\n got %q", testCase.err, actual)
		}
	}
}
