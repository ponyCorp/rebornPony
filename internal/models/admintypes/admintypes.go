package admintypes

type AdminType int
type levelMap map[int]AdminType

const (
	Owner AdminType = 100
	L1    AdminType = 1
	L2    AdminType = 2
	L3    AdminType = 3
	User  AdminType = 0
)

func getMap() levelMap {
	return map[int]AdminType{
		0:   User,
		1:   L1,
		2:   L2,
		3:   L3,
		100: Owner,
	}
}
func MaxLevel() AdminType {
	m := getMap()
	max := User
	for _, val := range m {
		if max < val && val != Owner {
			max = val
		}
	}
	return max
}
func New(level int) AdminType {
	mapLevels := getMap()
	levelType, ok := mapLevels[level]
	if !ok {
		return User
	}
	return levelType
}
func (t AdminType) Number() int {
	return int(t)
}

// IncreaseLevel
func (t AdminType) IncreaseLevel() AdminType {
	if t == MaxLevel() {
		return t
	}
	return New(t.Number() + 1)
}

// DecreaseLevel
func (t AdminType) DecreaseLevel() AdminType {
	if t == User {
		return t
	}
	return New(t.Number() - 1)
}
