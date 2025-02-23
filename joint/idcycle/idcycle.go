// Package implements an internal mechanism to communicate with an impress terminal.
package idcycle

import "sync"

type ID struct {
	nextID  int
	freeIDs []int
	mutex   sync.Mutex
}

func New() *ID { return &ID{} }

func (i *ID) Next() int {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if len(i.freeIDs) == 0 {
		i.nextID++
		return i.nextID
	}
	nextID := i.freeIDs[len(i.freeIDs)-1]
	i.freeIDs = i.freeIDs[:len(i.freeIDs)-1]
	return nextID
}

func (i *ID) Back(id int) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.freeIDs = append(i.freeIDs, id)
}
