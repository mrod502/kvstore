package kvstore

import (
	"sync"
	"time"
)

// -------------- Exported --------------- //

//ByteObj - an oject that returns its data as bytes
type ByteObj interface {
	GetData() []byte
}

//Item - a data container
type Item struct {
	Data   ByteObj
	Expiry int64
}

//GetData - call the underlying item's GetData() method
func (i Item) GetData() []byte {
	return i.Data.GetData()
}

//Store - a thread-safe map with expiration
type Store struct {
	vals map[string]Item
	mux  *sync.RWMutex
}

//Get a value
func (s Store) Get(k string) Item {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.vals[k]
}

//GetValue - return the underlying object stored
func (s Store) GetValue(k string) ByteObj {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.vals[k].Data
}

//Set a value
func (s Store) Set(k string, v ByteObj, exp int64) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.vals[k] = Item{Data: v, Expiry: exp}

}

//Delete a value
func (s Store) Delete(k string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.vals, k)

}

//NewStore constructs a store and starts the janitor if enabled
func NewStore(enableJanitor bool) (s *Store) {
	s = &Store{vals: make(map[string]Item), mux: new(sync.RWMutex)}
	if enableJanitor {
		go s.janitor()
	}
	return s
}

// -------------- End Exported --------------- //

func (s Store) getKeys() (out []string) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	out = make([]string, len(s.vals))
	var i int
	for k := range s.vals {
		out[i] = k
	}
	return out
}

//janitor - check expired
func (s *Store) janitor() {
	for {
		time.Sleep(time.Second)
		now := time.Now().Unix()

		for _, k := range s.getKeys() {
			if tt := s.Get(k).Expiry; tt < now && tt != 0 {
				s.mux.Lock()
				delete(s.vals, k)
				s.mux.Unlock()
			}
		}
	}
}
