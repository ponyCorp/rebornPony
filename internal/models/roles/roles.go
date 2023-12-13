package roles

type Role string

const (
	Owner Role = "owner"
	Admin Role = "admin"
	User  Role = "user"
)

func (r Role) String() string {
	return string(r)
}
