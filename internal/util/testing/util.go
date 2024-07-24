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
}

type Describable interface {
	Describe() string
}

type TestHandler func(*testing.T, Describable) error

func GetTestHandler[TestType Describable, GotType any](executeTest func(TestType) GotType, validateTest func(*testing.T, TestType, any), cleanup func()) TestHandler {
	return func(t *testing.T, tt Describable) error {
		if a, ok := tt.(TestType); ok {
			got := executeTest(a)
			validateTest(t, a, got)
			cleanup()
		} else {
			return fmt.Errorf("Something went wrong when asserting Describable as the generic type TestType.")
		}

		return nil
	}
}

func HandleTests[TestType Describable](t *testing.T, tests []TestType, testHandler TestHandler) {
	for _, tt := range tests {
		t.Run(tt.Describe(), func(t *testing.T) {
			if err := testHandler(t, tt); err != nil {
				panic(err)
			}
		})
	}
}

func AssertGotAndWantType[V any](t *testing.T, gotBeforeAssertion any, wantBeforeAssertion any) (V, V) {
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
