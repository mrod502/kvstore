package kvstore

import (
	"sync"
	"time"
)

//ByteObj - an oject that returns its data as bytes
type ByteObj interface {
	GetData() []byte
}
type Item struct {
	Data   ByteObj
	Expiry time.Time
}

func (i Item) GetData() []byte {
	return i.Data.GetData()
}

//Store - a thread-safe map with expiration
type Store struct {
	vals map[string]Item
	mux  *sync.RWMutex
}

//Get a value's data
func (s Store) Get(k string) Item {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.vals[k]
}

//Set a value
func (s Store) Set(k string, v ByteObj, exp time.Time) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.vals[k] = Item{Data: v, Expiry: exp}

}

//Set a value
func (s Store) Delete(k string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.vals, k)

}

func NewStore() (s *Store) {
	s = &Store{vals: make(map[string]Item), mux: new(sync.RWMutex)}
	go s.janitor()
	return s
}

func (s *Store) janitor() {

	for {
		time.Sleep(time.Second)
		now := time.Now()
		for k := range s.vals {
			if s.Get(k).Expiry < now {
				s.Delete(k)
			}
		}
	}
}
