package userrepository

import (
	"fmt"
	"http-server/internal/models"
	testingutil "http-server/internal/util/testing"
	"os"
	"testing"
)

var repository UserRepository

// For TestCreateUser
// type user struct {
// 	id                 int
// 	username, password string
// }

type singleUserTest struct {
	testingutil.BasicTest
	models.User
}

type twoUsersTest struct {
	testingutil.BasicTest
	users []models.User
}

// For TestGetUserById
type twoIdsTest struct {
	testingutil.BasicTest
	ids []int
}

var (
	USER = models.User{
		Id:       1,
		Username: "daniel",
		Password: "123456",
	}
	ANOTHER_USER = models.User{
		Id:       2,
		Username: "karl",
		Password: "123456",
	}
)

func (test singleUserTest) String() string {
	return test.Description
}

func (test twoUsersTest) String() string {
	return test.Description
}

func (test twoIdsTest) String() string {
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
					Error:       CREATE_USER_PASSWORD_TOO_SHORT_ERR,
				},
				models.User{
					Id:       USER.Id,
					Username: USER.Username,
					Password: USER.Password[:len(USER.Password)-1],
				},
			},
		}

		twoUsersTests = []twoUsersTest{
			{
				testingutil.BasicTest{
					Description: "Successfully create two users",
					Want:        2,
				},
				[]models.User{
					USER,
					ANOTHER_USER,
				},
			},

			{
				testingutil.BasicTest{
					Description: "Conflicting usernames",
					Want:        1,
					Error:       CREATE_USER_USERNAME_TAKEN_ERR,
				},
				[]models.User{
					USER,
					{
						Id:       USER.Id,
						Username: USER.Username,
						Password: USER.Password,
					},
				},
			},
		}

		return singleUserTests, twoUsersTests
	}

	setupTestHandlers := func() (testingutil.TestHandler, testingutil.TestHandler) {
		const TEST_FUNCTION = "CreateUser"

		executeSingleUserTest := func(t *testing.T, tt singleUserTest) int {
			err := testFunction(tt.Username, tt.Password)
			testingutil.ValidateError(t, TEST_FUNCTION, err, tt.Error)
			return repository.count()
		}

		validateSingleUserTest := func(t *testing.T, tt singleUserTest, gotBeforeAssertion any) {
			got, want := testingutil.AssertGotAndWantType[int](t, gotBeforeAssertion, tt.Want)
			err := fmt.Sprintf(
				"%s(%s, %s) -> repository.count() = %d, want: %d",
				TEST_FUNCTION, tt.Username, tt.Password, got, want)

			testingutil.ValidateResult(t, err, got, want)
		}

		executeTwoUsersTest := func(t *testing.T, tt twoUsersTest) int {
			testFunction(tt.users[0].Username, tt.users[0].Password)
			err := testFunction(tt.users[1].Username, tt.users[1].Password)

			testingutil.ValidateError(t, TEST_FUNCTION, err, tt.Error)

			return repository.count()
		}

		validateTwoUsersTest := func(t *testing.T, tt twoUsersTest, gotBeforeAssertion any) {
			got, want := testingutil.AssertGotAndWantType[int](t, gotBeforeAssertion, tt.Want)
			err := fmt.Sprintf(
				"%s(%s, %s) -> %s(%s, %s) -> repository.count() = %d, want: %d",
				TEST_FUNCTION, tt.users[0].Username, tt.users[0].Password,
				TEST_FUNCTION, tt.users[1].Username, tt.users[1].Password,
				got, want)

			testingutil.ValidateResult(t, err, got, want)
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
	testFunction := repository.GetUserById
	createUser := repository.CreateUser

	setupTests := func() []twoIdsTest {
		tests := []twoIdsTest{
			{
				testingutil.BasicTest{
					Description: "Gets correct user by comparing username",
					Want:        []models.User{USER, ANOTHER_USER},
				},
				[]int{
					1,
					2,
				},
			},
			{
				testingutil.BasicTest{
					Description: "Throws error if ID does not exist",
					Want:        nil,
					Error:       GET_USER_BY_ID_ERR,
				},
				[]int{
					0,
					3,
				},
			},
			{
				testingutil.BasicTest{
					Description: "Is idempotent",
					Want:        []models.User{USER, USER},
				},
				[]int{
					1,
					1,
				},
			},
		}

		return tests
	}

	setupTestHandler := func() testingutil.TestHandler {
		const TEST_FUNCTION = "GetUserById"

		executeTests := func(t *testing.T, tt twoIdsTest) []models.User {
			createUser(USER.Username, USER.Password)
			createUser(ANOTHER_USER.Username, ANOTHER_USER.Password)
			res := make([]models.User, 0, 2)

			for _, id := range tt.ids {
				user, err := testFunction(id)

				testingutil.ValidateError(t, TEST_FUNCTION, err, tt.Error)

				if user != nil {
					res = append(res, *user)
				}
			}

			return res
		}

		validateTests := func(t *testing.T, tt twoIdsTest, gotsBeforeAssertion any) {
			const CREATE_USER = "CreateUser"

			gots, gotsOk := gotsBeforeAssertion.([]models.User)
			wants, wantsOk := tt.Want.([]models.User)

			if gotsOk && wantsOk {
				for i, got := range gots {
					err := fmt.Sprintf(
						"%s(%s, %s) -> %s(%s, %s) -> %s(%d) = %s, want: %s",
						CREATE_USER, USER.Username, USER.Password,
						CREATE_USER, USER.Username, USER.Password,
						TEST_FUNCTION, tt.ids[i], got.String(), wants[i].Username,
					)

					testingutil.ValidateResult(t, err, got, wants[i])
				}
			}
		}

		return testingutil.GetTestHandler(executeTests, validateTests, cleanup)
	}

	tests := setupTests()
	testHandler := setupTestHandler()

	testingutil.HandleTests(t, tests, testHandler)
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
