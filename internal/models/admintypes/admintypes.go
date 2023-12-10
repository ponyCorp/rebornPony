package admintypes

type AdminType int

const (
	Owner AdminType = 100
	L1    AdminType = 1
	L2    AdminType = 2
	L3    AdminType = 3
	User  AdminType = 0
)

func New(level int) AdminType {
	levels := map[int]AdminType{
		100: Owner,
		1:   L1,
		2:   L2,
		3:   L3,
		0:   User,
	}
	levelType, ok := levels[level]
	if !ok {
		return User
	}
	return levelType
}
func (t AdminType) Number() int {
	return int(t)
}
