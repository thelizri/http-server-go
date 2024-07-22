package testingutil

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type Describable interface {
	Description() string
}

func GetTestHandler[TestType Describable, GotType any](executeTest func(TestType) GotType, validateTest func(*testing.T, TestType, GotType), cleanup func()) func(*testing.T, TestType) {
	return func(t *testing.T, tt TestType) {
		got := executeTest(tt)
		validateTest(t, tt, got)
		cleanup()
	}
}

func HandleTests[TestType Describable](t *testing.T, tests []TestType, testHandler func(*testing.T, TestType)) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf(tt.Description(), i), func(t *testing.T) {
			testHandler(t, tt)
		})
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
