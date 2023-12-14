package userservice

import (
	"sync"

	"github.com/ponyCorp/rebornPony/internal/models"
)

type group struct {
	Users map[int64]*userMutex
	gmx   sync.RWMutex
}
type userMutex struct {
	umx  sync.RWMutex
	User *models.User
}
type userDBAdapter interface {
	GetOrCreateUser(chatId int64, userId int64) (*models.User, error)
	IncreaseReputation(chatId int64, userId int64, inc int) (int, error)
	DecreaseReputation(chatId int64, userId int64, inc int) (int, error)
	SetReputation(chatId int64, userId int64, reputation int) (int, error)
}

func (g *group) Get(userId int64) (*userMutex, bool) {
	defer g.gmx.RUnlock()
	g.gmx.RLock()
	user, ok := g.Users[userId]

	// defer user.umx.RUnlock()
	// user.umx.RLock()

	return user, ok
}
func (g *group) Set(userId int64, user *models.User) *userMutex {
	defer g.gmx.Unlock()
	g.gmx.Lock()
	g.Users[userId] = &userMutex{
		umx:  sync.RWMutex{},
		User: user,
	}
	return g.Users[userId]
}

type UserService struct {
	userRepo userDBAdapter

	//[chatId]map[userId]*models.User
	cache map[int64]*group
	mx    sync.RWMutex
}

func newGroup() *group {
	return &group{
		Users: make(map[int64]*userMutex),
		gmx:   sync.RWMutex{},
	}
}
func New(userAdapter userDBAdapter) *UserService {
	return &UserService{
		userRepo: userAdapter,
		cache:    make(map[int64]*group),
	}
}

// GetOrCreateGroup(chatId int64) *group
func (u *UserService) GetOrCreateGroup(chatId int64) *group {
	defer u.mx.Unlock()
	u.mx.Lock()
	chat, ok := u.cache[chatId]
	if !ok {
		chat = newGroup()
		u.cache[chatId] = chat
	}
	return chat
}
func (u *UserService) GetOrCreateUser(chatId int64, userId int64) (*userMutex, error) {
	chat := u.GetOrCreateGroup(chatId)
	user, ok := chat.Get(userId)
	if !ok {
		usr, err := u.userRepo.GetOrCreateUser(chatId, userId)
		if err != nil {
			return nil, err
		}
		user = chat.Set(userId, usr)
	}
	return user, nil
}

// IncreaseReputation(chatId int64, userId int64, inc int) error
func (u *UserService) IncreaseReputation(chatId int64, userId int64, inc int) (int, error) {
	user, err := u.GetOrCreateUser(chatId, userId)
	if err != nil {
		return 0, err
	}
	defer user.umx.Unlock()
	user.umx.Lock()
	reputation, err := u.userRepo.IncreaseReputation(chatId, userId, inc)
	if err != nil {
		return 0, err
	}
	user.User.Reputation = reputation
	return reputation, nil
}

// DecreaseReputation(chatId int64, userId int64, inc int) error
func (u *UserService) DecreaseReputation(chatId int64, userId int64, inc int) (int, error) {
	user, err := u.GetOrCreateUser(chatId, userId)
	if err != nil {
		return 0, err
	}
	defer user.umx.Unlock()
	user.umx.Lock()
	reputation, err := u.userRepo.DecreaseReputation(chatId, userId, inc)
	if err != nil {
		return 0, err
	}
	user.User.Reputation = reputation
	return reputation, nil
}

// SetReputation(chatId int64, userId int64, reputation int) error
func (u *UserService) SetReputation(chatId int64, userId int64, rep int) (int, error) {
	user, err := u.GetOrCreateUser(chatId, userId)
	if err != nil {
		return 0, err
	}
	defer user.umx.Unlock()
	user.umx.Lock()
	reputation, err := u.userRepo.SetReputation(chatId, userId, rep)
	if err != nil {
		return 0, err
	}
	user.User.Reputation = reputation
	return reputation, nil
}
