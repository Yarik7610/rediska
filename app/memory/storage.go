package memory

type baseStorage interface {
	Keys() []string
	Has(key string) bool
	Del(key string)
}

type MultiTypeStorage struct {
	storages map[string]baseStorage
}

const (
	TYPE_STRING = "string"
	TYPE_LIST   = "list"
	TYPE_STREAM = "stream"
	TYPE_NONE   = "none"
)

func NewMultiTypeStorage() *MultiTypeStorage {
	return &MultiTypeStorage{
		storages: map[string]baseStorage{
			TYPE_STRING: NewStringStorage(),
			TYPE_LIST:   NewListStorage(),
			TYPE_STREAM: NewStreamStorage(),
		},
	}
}

func (s *MultiTypeStorage) Keys() []string {
	allStorageKeys := make([]string, 0)
	allStorageKeys = append(allStorageKeys, s.StringStorage().Keys()...)
	allStorageKeys = append(allStorageKeys, s.ListStorage().Keys()...)
	allStorageKeys = append(allStorageKeys, s.StreamStorage().Keys()...)
	return allStorageKeys
}

func (s *MultiTypeStorage) Del(key string) {
	if s.StringStorage().Has(key) {
		s.StringStorage().Del(key)
	} else if s.ListStorage().Has(key) {
		s.ListStorage().Del(key)
	} else if s.StreamStorage().Has(key) {
		s.StreamStorage().Del(key)
	}
}

func (s *MultiTypeStorage) KeyExistsWithOtherType(key string, allowedType string) bool {
	for storageKey, storage := range s.storages {
		if storageKey == allowedType {
			continue
		}
		if storage.Has(key) {
			return true
		}
	}
	return false
}

func (s *MultiTypeStorage) StringStorage() StringStorage {
	return s.storages[TYPE_STRING].(StringStorage)
}

func (s *MultiTypeStorage) ListStorage() ListStorage {
	return s.storages[TYPE_LIST].(ListStorage)
}

func (s *MultiTypeStorage) StreamStorage() StreamStorage {
	return s.storages[TYPE_STREAM].(StreamStorage)
}
