package userrepository

import (
	testingutil "http-server/internal/util/testing"
	"os"
	"strconv"
	"testing"
)

var repository UserRepository

type basicTest struct {
	description string
	want        any
}

type user struct {
	username, password string
}

type anotherUser struct {
	username, password string
}

type singleUserTest struct {
	testingutil.BasicTest
	user
}

type twoUsersTest struct {
	testingutil.BasicTest
	user
	anotherUser
}

var (
	USER = user{
		"daniel",
		"123456",
	}
	ANOTHER_USER = anotherUser{
		"karl",
		"123456",
	}
)

func (test singleUserTest) Describe() string {
	return test.Description
}

func (test twoUsersTest) Describe() string {
	return test.Description
}

func TestMain(m *testing.M) {
	afterAll := beforeAll()
	code := m.Run()
	afterAll(code)
}

func TestCreateUser(t *testing.T) {
	testFunction := repository.CreateUser

	setupTests := func() (singleUserTests []singleUserTest, twoUsersTests []twoUsersTest) {
		singleUserTests = []singleUserTest{
			{
				testingutil.BasicTest{
					Description: "Successfully create user",
					Want:        1,
				},
				USER,
			},

			{
				testingutil.BasicTest{
					Description: "Password too short (less than 6 chars)",
					Want:        0,
				},
				user{
					USER.username,
					USER.password[:len(USER.password)-1],
				},
			},
		}

		twoUsersTests = []twoUsersTest{
			{
				testingutil.BasicTest{
					Description: "Successfully create two users",
					Want:        2,
				},
				USER,
				ANOTHER_USER,
			},

			{
				testingutil.BasicTest{
					Description: "Conflicting usernames",
					Want:        1,
				},
				USER,
				anotherUser{
					USER.username,
					USER.password,
				},
			},
		}

		return singleUserTests, twoUsersTests
	}

	setupTestHandlers := func() (testingutil.TestHandler, testingutil.TestHandler) {
		const TEST_FUNCTION_STR = "CreateUser"
		const USERS_STR = " users"

		executeSingleUserTest := func(tt singleUserTest) int {
			testFunction(tt.username, tt.password)
			return repository.count()
		}

		validateSingleUserTest := func(t *testing.T, tt singleUserTest, gotBeforeAssertion any) {
			got, want := testingutil.AssertGotAndWantType[int](t, gotBeforeAssertion, tt.Want)
			gotStr, wantStr := strconv.Itoa(got)+USERS_STR, strconv.Itoa(want)+USERS_STR

			if got != tt.Want {
				t.Errorf("%s(%s, %s) -> repository.count() = %s, want: %s",
					TEST_FUNCTION_STR, tt.username, tt.password,
					gotStr, wantStr)
			}
		}

		executeTwoUsersTest := func(tt twoUsersTest) int {
			testFunction(tt.user.username, tt.user.password)
			testFunction(tt.anotherUser.username, tt.anotherUser.password)
			return repository.count()
		}

		validateTwoUsersTest := func(t *testing.T, tt twoUsersTest, gotBeforeAssertion any) {
			got, want := testingutil.AssertGotAndWantType[int](t, gotBeforeAssertion, tt.Want)
			gotStr, wantStr := strconv.Itoa(got)+USERS_STR, strconv.Itoa(want)+USERS_STR

			if got != tt.Want {
				t.Errorf("%s(%s, %s) -> %s(%s, %s) -> repository.count() = %s, want: %s",
					TEST_FUNCTION_STR, tt.user.username, tt.user.password,
					TEST_FUNCTION_STR, tt.anotherUser.username, tt.anotherUser.password,
					gotStr, wantStr)
			}
		}

		singleUserTestHandler := testingutil.GetTestHandler(executeSingleUserTest, validateSingleUserTest, cleanup)
		twoUsersTestHandler := testingutil.GetTestHandler(executeTwoUsersTest, validateTwoUsersTest, cleanup)
		return singleUserTestHandler, twoUsersTestHandler
	}

	singleUserTests, twoUsersTests := setupTests()
	singleUserTestHandler, twoUsersTestHandler := setupTestHandlers()
	testingutil.HandleTests(t, singleUserTests, singleUserTestHandler)
	testingutil.HandleTests(t, twoUsersTests, twoUsersTestHandler)
}

func TestGetUserById(t *testing.T) {

}

func cleanup() {
	repository.deleteAll()
}

func beforeAll() func(int) {
	repository = NewUserRepository()

	return func(code int) {
		repository = nil
		os.Exit(code)
	}
}
