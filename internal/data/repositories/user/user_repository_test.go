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
	want        int
}

type user struct {
	username, password string
}

type anotherUser struct {
	username, password string
}

type singleUserTest struct {
	basicTest
	user
}

type twoUsersTest struct {
	basicTest
	user
	anotherUser
}

func (test singleUserTest) Description() string {
	return test.description
}

func (test twoUsersTest) Description() string {
	return test.description
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
				basicTest{
					"Successfully create user",
					1,
				},
				user{
					"daniel",
					"123456",
				},
			},

			{
				basicTest{
					"Password too short (less than 6 chars)",
					0,
				},
				user{
					"daniel",
					"12345",
				},
			},
		}

		twoUsersTests = []twoUsersTest{
			{
				basicTest{
					"Successfully create two users",
					2,
				},
				user{
					"daniel",
					"123456",
				},
				anotherUser{
					"karl",
					"123456",
				},
			},

			{
				basicTest{
					"Conflicting usernames",
					1,
				},
				user{
					"daniel",
					"123456",
				},
				anotherUser{
					"daniel",
					"123456",
				},
			},
		}

		return singleUserTests, twoUsersTests
	}

	setupTestHandlers := func() (func(*testing.T, singleUserTest), func(*testing.T, twoUsersTest)) {
		executeSingleUserTest := func(tt singleUserTest) int {
			testFunction(tt.username, tt.password)
			return repository.count()
		}

		validateSingleUserTest := func(t *testing.T, tt singleUserTest, got int) {
			const USERS_STR = " users"

			args := []string{tt.username, tt.password}
			gotStr, wantStr := strconv.Itoa(got)+USERS_STR, strconv.Itoa(tt.want)+USERS_STR

			if got != tt.want {
				t.Errorf(testingutil.ParseError(testFunction, args, gotStr, wantStr))
			}
		}

		executeTwoUsersTest := func(tt twoUsersTest) int {
			testFunction(tt.user.username, tt.user.password)
			testFunction(tt.anotherUser.username, tt.anotherUser.password)
			return repository.count()
		}

		validateTwoUsersTest := func(t *testing.T, tt twoUsersTest, got int) {
			if got != tt.want {
				t.Errorf("CreateUser(%s, %s) -> CreateUser(%s, %s) = %d, want: %d", tt.user.username, tt.user.password, tt.anotherUser.username, tt.anotherUser.password, got, tt.want)
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
