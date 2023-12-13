package userservice

import (
	"fmt"
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
	fn := args.Get(0).(func(chatId int64, userId int64, user *models.User) (*models.User, error))
	user, err := fn(chatId, userId, &models.User{})
	if err == nil {
		// Пользователь уже существует в базе данных
		return user, nil
	}

	// Пользователь не существует в базе данных
	user = &models.User{
		UserId:     userId,
		ChatId:     chatId,
		Role:       "User",
		Reputation: 0,
		Level:      1,
	}

	// Создаем нового пользователя
	err = m.DB.GetOrCreateUser(chatId, userId, user)

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
	adapter.On("GetOrCreateUser", mock.Anything, mock.Anything, mock.Anything).Return(func(chatId int64, userId int64, user *models.User) (*models.User, error) {
		builder := strings.Builder{}
		chatStr := strconv.FormatInt(chatId, 10)
		userStr := strconv.FormatInt(userId, 10)
		builder.WriteString(chatStr)
		builder.WriteString("_")
		builder.WriteString(userStr)
		if _, ok := adapter.DB.users[builder.String()]; ok {
			return user, nil
		}
		adapter.DB.users[builder.String()] = user
		return user, nil
	})
	adapter.On("IncreaseReputation", mock.Anything, mock.Anything, mock.Anything).Return(func(chatId int64, userId int64, inc int) (int, error) {
		builder := strings.Builder{}
		chatStr := strconv.FormatInt(chatId, 10)
		userStr := strconv.FormatInt(userId, 10)
		builder.WriteString(chatStr)
		builder.WriteString("_")
		builder.WriteString(userStr)
		user, ok := adapter.DB.users[builder.String()]
		if !ok {
			return 0, fmt.Errorf("User not found")
		}
		user.Reputation += inc
		return user.Reputation, nil
	})
	adapter.On("DecreaseReputation", mock.Anything, mock.Anything, mock.Anything).Return(func(chatId int64, userId int64, inc int) (int, error) {
		builder := strings.Builder{}
		chatStr := strconv.FormatInt(chatId, 10)
		userStr := strconv.FormatInt(userId, 10)
		builder.WriteString(chatStr)
		builder.WriteString("_")
		builder.WriteString(userStr)
		user, ok := adapter.DB.users[builder.String()]
		if !ok {
			return 0, fmt.Errorf("User not found")
		}
		user.Reputation -= inc
		return user.Reputation, nil
	})
	adapter.On("SetReputation", mock.Anything, mock.Anything, mock.Anything).Return(func(chatId int64, userId int64, rep int) (int, error) {
		builder := strings.Builder{}
		chatStr := strconv.FormatInt(chatId, 10)
		userStr := strconv.FormatInt(userId, 10)
		builder.WriteString(chatStr)
		builder.WriteString("_")
		builder.WriteString(userStr)
		user, ok := adapter.DB.users[builder.String()]
		if !ok {
			return 0, fmt.Errorf("User not found")
		}
		user.Reputation = rep
		return user.Reputation, nil
	})
	userService := New(adapter)
	s.userService = userService
}

// TestAdapter
func (s *MySuite) TestAdapter() {

	// Вызовите метод GetOrCreateUser
	user, err := s.userService.GetOrCreateUser(1, 1)
	s.NoError(err)
	s.Equal(1, user.User.UserId)
	s.Equal(1, user.User.ChatId)
	s.Equal("User", user.User.Role)
	s.Equal(0, user.User.Reputation)
	s.Equal(1, user.User.Level)
}
func TestMySuite(t *testing.T) {
	suite.Run(t, new(MySuite))
}
