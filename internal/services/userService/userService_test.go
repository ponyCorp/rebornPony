package userservice

import (
	"strconv"
	"strings"
	"testing"

	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockDb struct {
	users map[string]*models.User
}

func (m *mockDb) GetOrCreateUser(chatId int64, userId int64, user *models.User) error {

	builder := strings.Builder{}
	chatStr := strconv.FormatInt(chatId, 10)
	userStr := strconv.FormatInt(userId, 10)
	builder.WriteString(chatStr)
	builder.WriteString("_")
	builder.WriteString(userStr)
	if _, ok := m.users[builder.String()]; ok {
		return nil
	}
	m.users[builder.String()] = user
	return nil
}

type UserDBAdapterMock struct {
	mock.Mock
	DB *mockDb
}

func (m *UserDBAdapterMock) GetOrCreateUser(chatId int64, userId int64) (*models.User, error) {
	args := m.Called(chatId, userId)
	user := args.Get(0).(*models.User)
	err := args.Error(1)
	return user, err
}

// IncreaseReputation(chatId int64, userId int64, inc int) (int, error)
func (m *UserDBAdapterMock) IncreaseReputation(chatId int64, userId int64, inc int) (int, error) {
	args := m.Called(chatId, userId, inc)
	reputation := args.Get(0).(int)
	err := args.Error(1)
	return reputation, err
}

// DecreaseReputation(chatId int64, userId int64, inc int) (int, error)
func (m *UserDBAdapterMock) DecreaseReputation(chatId int64, userId int64, inc int) (int, error) {
	args := m.Called(chatId, userId, inc)
	reputation := args.Get(0).(int)
	err := args.Error(1)
	return reputation, err
}

// SetReputation(chatId int64, userId int64, reputation int) (int, error)
func (m *UserDBAdapterMock) SetReputation(chatId int64, userId int64, rep int) (int, error) {
	args := m.Called(chatId, userId, rep)
	reputation := args.Get(0).(int)
	err := args.Error(1)
	return reputation, err
}

type MySuite struct {
	suite.Suite
	adapter     *UserDBAdapterMock
	userService *UserService
}

func (s *MySuite) SetupTest() {
	// Создайте мок для интерфейса userDBAdapter
	adapter := &UserDBAdapterMock{
		DB: &mockDb{
			users: map[string]*models.User{},
		},
	}
	s.adapter = adapter
	// Вызовите метод GetOrCreateUser
	adapter.On("GetOrCreateUser", int64(1), int64(1)).Return(&models.User{
		UserId:     1,
		ChatId:     1,
		Role:       "User",
		Reputation: 0,
		Level:      1,
	}, nil)
	adapter.On("GetOrCreateUser", int64(1), int64(2)).Return(&models.User{
		UserId:     2,
		ChatId:     1,
		Role:       "User",
		Reputation: 0,
		Level:      1,
	}, nil)

	//user 3 with reputation 5
	adapter.On("GetOrCreateUser", int64(1), int64(3)).Return(&models.User{
		UserId:     3,
		ChatId:     1,
		Role:       "User",
		Reputation: 5,
		Level:      1,
	}, nil)
	adapter.On("IncreaseReputation", int64(1), int64(2), 1).Return(1, nil)
	adapter.On("DecreaseReputation", int64(1), int64(3), 1).Return(4, nil)

	//user 4 with reputation 0
	adapter.On("GetOrCreateUser", int64(1), int64(4)).Return(&models.User{
		UserId:     4,
		ChatId:     1,
		Role:       "User",
		Reputation: 0,
		Level:      1,
	}, nil)
	//SetReputation user 4 with reputation 10
	adapter.On("SetReputation", int64(1), int64(4), 10).Return(10, nil)

	userService := New(adapter)
	s.userService = userService
}

// TestAdapter
func (s *MySuite) TestAdapter() {

	// Вызовите метод GetOrCreateUser
	user, err := s.userService.GetOrCreateUser(1, 1)
	s.NoError(err)
	s.Equal(int64(1), user.User.UserId)
	s.Equal(int64(1), user.User.ChatId)
	s.Equal("User", user.User.Role)
	s.Equal(0, user.User.Reputation)
	s.Equal(1, user.User.Level)
}

// test IncreaseReputation for user 2
func (s *MySuite) TestIncreaseReputation() {

	// Вызовите метод IncreaseReputation
	rep, err := s.userService.IncreaseReputation(1, 2, 1)
	s.NoError(err)
	s.Equal(1, rep)

}

// test DecreaseReputation for user 3
func (s *MySuite) TestDecreaseReputation() {
	rep, err := s.userService.DecreaseReputation(1, 3, 1)
	s.NoError(err)
	s.Equal(4, rep)
}

// test SetReputation for user 4 with reputation 10
func (s *MySuite) TestSetReputation() {
	rep, err := s.userService.SetReputation(1, 4, 10)
	s.NoError(err)
	s.Equal(10, rep)
}
func TestMySuite(t *testing.T) {
	suite.Run(t, new(MySuite))
}
