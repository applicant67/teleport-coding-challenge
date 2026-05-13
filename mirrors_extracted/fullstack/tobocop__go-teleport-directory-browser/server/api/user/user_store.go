package user

type Store interface {
	GetUserPassword(username string) (bool, string)
}

type StaticUserStore struct {
	users map[string]string
}

func NewStaticUserStore(users map[string]string) Store {
	return &StaticUserStore{ users }
}

func (s *StaticUserStore) GetUserPassword(username string) (bool, string)  {
	val, ok := s.users[username]
	return ok, val
}
