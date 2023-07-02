package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StorageAccount struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Login      string
	Password   string
	Created    time.Time
	Characters []Character
}

type Storage struct {
	Accounts           map[string]StorageAccount
	AccountsWithLogins map[string]StorageAccount
	sync.RWMutex
	path string
}

func NewStorage(path string) (*Storage, error) {
	storage := &Storage{
		RWMutex:            sync.RWMutex{},
		Accounts:           make(map[string]StorageAccount),
		AccountsWithLogins: make(map[string]StorageAccount),
		path:               path,
	}

	var (
		buf []byte
		err error
	)

	buf, err = os.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		buf = []byte("{}")

		err := os.WriteFile(path, buf, 0777)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, storage)

	return storage, err
}

func (s *Storage) SaveAccount(sa *StorageAccount) error {
	defer s.Save()

	s.Lock()
	defer s.Unlock()

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("can't generate UUID: %v", err)
	}

	s.Accounts[newUUID.String()] = *sa
	s.AccountsWithLogins[sa.Login] = *sa

	return nil
}

func (s *Storage) GetAccount(id string) (StorageAccount, bool) {
	s.RLock()
	defer s.RUnlock()

	sa, ok := s.Accounts[id]

	return sa, ok
}

func (s *Storage) GetAccountByLogin(login string) (StorageAccount, bool) {
	s.RLock()
	defer s.RUnlock()

	sa, ok := s.AccountsWithLogins[login]

	return sa, ok
}

func (s *Storage) Save() error {
	s.Lock()
	defer s.Unlock()

	err := os.Remove(s.path)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, buf, 0777)
}
