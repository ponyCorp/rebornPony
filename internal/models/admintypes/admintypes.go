package admintypes

type EntityType int
type Level int

const (
	EUser EntityType = iota
	EModerator
	EOwner
)

type Entity interface {
	IncreaseLevel() Level
	DecreaseLevel() Level
	Level() Level
	Type() EntityType
}

// type User struct {
// 	level Level
// }

func ItoLevel(i int) Level {
	return Level(i)
}

// func NewUser() *User {
// 	return &User{}
// }
// func (u *User) IncreaseLevel() Level {
// 	u.level++
// 	return u.level
// }
// func (u *User) DecreaseLevel() Level {
// 	u.level--
// 	return u.level
// }
// func (u *User) Level() Level {
// 	return 0
// }
// func (u *User) Type() EntityType {
// 	return EUser
// }

type AdminType string

// type levelMap map[int]AdminType

const (
	User AdminType = "user"

	// 	U1 AdminType = 1
	// 	U2 AdminType = 2
	// 	U3 AdminType = 3
	Moderator AdminType = "moderator"
	// 	//Moderator level with 20+ level
	// 	M1 AdminType = 20
	// 	M2 AdminType = 21
	// 	M3 AdminType = 22
	// 	M4 AdminType = 23

	Owner AdminType = "owner"
)

// func GetEntityType(level AdminType) EntityType {
// 	if level == Owner {
// 		return EOwner
// 	}
// 	if level < 20 {
// 		return EUser
// 	}
// 	return EModerator
// }
// func userList() []AdminType {
// 	return []AdminType{
// 		User,
// 		U1,
// 		U2,
// 		U3,
// 	}
// }
// func moderatorList() []AdminType {
// 	return []AdminType{
// 		M1,
// 		M2,
// 		M3,
// 		M4,
// 	}
// }
// func getMap() levelMap {
// 	//list of levels
// 	list := []AdminType{
// 		Owner,
// 	}
// 	list = append(list, userList()...)
// 	list = append(list, moderatorList()...)
// 	mapList := make(map[int]AdminType)
// 	for _, v := range list {
// 		mapList[v.Number()] = v
// 	}

//		return mapList
//	}
//
//	func MaxLevel() AdminType {
//		m := getMap()
//		max := User
//		for _, val := range m {
//			if max < val && val != Owner {
//				max = val
//			}
//		}
//		return max
//	}
//
//	func New(level int) AdminType {
//		mapLevels := getMap()
//		levelType, ok := mapLevels[level]
//		if !ok {
//			return User
//		}
//		return levelType
//	}
// func (t AdminType) Number() int {
// 	return int(t)
// }

// // IncreaseLevel
// func (t AdminType) IncreaseLevel() AdminType {
// 	if t == MaxLevel() {
// 		return t
// 	}
// 	return New(t.Number() + 1)
// }

// // DecreaseLevel
// func (t AdminType) DecreaseLevel() AdminType {
// 	if t == User {
// 		return t
// 	}
// 	return New(t.Number() - 1)
// }
