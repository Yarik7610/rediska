package memory

type Storage struct {
	StringStorage *StringStorage
	ListStorage   *ListStorage
}

func NewStorage() *Storage {
	return &Storage{
		StringStorage: NewStringStorage(),
		ListStorage:   NewListStorage(),
	}
}

func (s *Storage) GetKeys() []string {
	allStorageKeys := make([]string, 0)
	allStorageKeys = append(allStorageKeys, s.StringStorage.GetKeys()...)
	allStorageKeys = append(allStorageKeys, s.ListStorage.GetKeys()...)
	return allStorageKeys
}

func (s *Storage) Del(key string) {
	if _, ok := s.StringStorage.Get(key); ok {
		s.StringStorage.Del(key)
	} else if _, ok := s.ListStorage.Get(key); ok {
		s.ListStorage.Del(key)
	}
}
