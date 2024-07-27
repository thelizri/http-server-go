package testingutil

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type BasicTest struct {
	Description string
	Want        any
	Error       string
}

type Describable interface {
	String() string
}

type TestHandler func(*testing.T, Describable)

func GetTestHandler[TestType Describable, GotType any](
	executeTest func(*testing.T, TestType) GotType, validateTest func(*testing.T, TestType, any), cleanup func()) TestHandler {
	return func(t *testing.T, tt Describable) {
		if a, ok := tt.(TestType); ok {
			got := executeTest(t, a)
			validateTest(t, a, got)
			cleanup()
		} else {
			t.Error("Something went wrong when asserting Describable as the generic type TestType.")
		}
	}
}

func HandleTests[TestType Describable](t *testing.T, tests []TestType, testHandler TestHandler) {
	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			testHandler(t, tt)
		})
	}
}

func AssertGotAndWantType[V comparable](t *testing.T, gotBeforeAssertion any, wantBeforeAssertion any) (V, V) {
	got, gotOk := gotBeforeAssertion.(V)
	want, wantOk := wantBeforeAssertion.(V)

	validateAssertion := func(name string, gw any, ok bool) {
		if !ok {
			t.Errorf("%s: expected type '%s' but received '%s'", name, reflect.TypeOf((*V)(nil)), reflect.TypeOf(gw))
		}
	}

	validateAssertion("got", got, gotOk)
	validateAssertion("want", want, wantOk)

	return got, want
}

func ValidateResult(t *testing.T, errStr string, got any, want any) {
	if got != want {
		t.Errorf(errStr)
	}
}

func ValidateError(t *testing.T, testFunctionStr string, actualError error, expectedErrorStr string) {
	if actualError != nil && actualError.Error() != expectedErrorStr {
		t.Errorf("%s expected the error '%s' but got '%s'", testFunctionStr, expectedErrorStr, actualError.Error())
	}
}

func parseDescription(args []string) string {
	return parseArgs(args)
}

func ParseError(function interface{}, args []string, got string, want string) string {
	return fmt.Sprintf("%s%s = %s, want: %s", getFunctionName(function), parseArgs(args), got, want)
}

func parseArgs(args []string) string {
	var sb strings.Builder

	appendArgs := func() {
		for i, arg := range args {
			sb.WriteString(arg)

			if i != len(args)-1 {
				sb.WriteString(", ")
			}
		}
	}

	sb.WriteString("(")
	appendArgs()
	sb.WriteString(")")

	return sb.String()
}

func getFunctionName(i interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	filterPackagePath := func() {
		lastSlash := strings.LastIndex(fullName, "/")
		if lastSlash != -1 {
			fullName = fullName[lastSlash+1:]
		}
	}

	filterFmSuffix := func() {
		if strings.HasSuffix(fullName, "-fm") {
			fullName = strings.TrimSuffix(fullName, "-fm")
		}
	}

	filterPackagePath()
	filterFmSuffix()

	return fullName
}
