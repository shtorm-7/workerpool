package callbackfield

import (
	"sync"

	list "github.com/bahlo/generic-list-go"
)

type (
	waiter struct {
		rLock bool
		mtx   *sync.Mutex
	}

	RWMutex struct {
		waiters list.List[waiter]
		pivot   *list.Element[waiter]

		locks int

		mtx sync.Mutex
	}
)

func (m *RWMutex) Lock() {
	m.mtx.Lock()
	if m.locks == 0 {
		m.locks = -1
		m.mtx.Unlock()
		return
	}
	mtx := m.createWaiter(false)
	m.mtx.Unlock()
	mtx.Lock()
}

func (m *RWMutex) RLock() {
	m.mtx.Lock()
	if m.locks == 0 || (m.locks > 0 && m.waiters.Front() == nil) {
		m.locks++
		m.mtx.Unlock()
		return
	}
	mtx := m.createWaiter(true)
	m.mtx.Unlock()
	mtx.Lock()
}

func (m *RWMutex) Unlock() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.locks != -1 {
		panic("Unlock of unlocked mutex")
	}
	m.locks = 0
	m.notifyWaiters()
}

func (m *RWMutex) RUnlock() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.locks < 1 {
		panic("RUnlock of unlocked mutex")
	}
	m.locks--
	if m.locks == 0 {
		m.notifyWaiters()
	}
}

func (m *RWMutex) UnlockToRLock() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.locks != -1 {
		panic("Unlock of unlocked mutex")
	}
	m.locks = 1
	m.notifyWaiters()
}

func (m *RWMutex) RUnlockToLock() {
	m.mtx.Lock()
	if m.locks < 1 {
		panic("RUnlock of unlocked mutex")
	}
	if m.locks == 1 && (m.pivot == nil || m.pivot.Prev() == nil) {
		m.locks = -1
		m.mtx.Unlock()
		return
	}
	mtx := m.createPriorityLockWaiter()
	m.locks--
	if m.locks == 0 {
		m.notifyWaiters()
	}
	m.mtx.Unlock()
	mtx.Lock()
}

func (m *RWMutex) notifyWaiters() {
	for {
		waiterElement := m.waiters.Front()
		if waiterElement == nil {
			return
		}
		if waiterElement.Value.rLock {
			if m.locks >= 0 {
				m.locks++
				m.releaseWaiter(waiterElement)
			} else {
				return
			}
		} else {
			if m.locks == 0 {
				m.locks = -1
				m.releaseWaiter(waiterElement)
			}
			return
		}
	}
}

func (m *RWMutex) createWaiter(rLock bool) *sync.Mutex {
	waiter := waiter{
		rLock: rLock,
		mtx:   new(sync.Mutex),
	}
	if m.pivot == nil {
		m.pivot = m.waiters.PushBack(waiter)
	} else {
		m.waiters.PushBack(waiter)
	}
	waiter.mtx.Lock()
	return waiter.mtx
}

func (m *RWMutex) createPriorityLockWaiter() *sync.Mutex {
	waiter := waiter{
		mtx: new(sync.Mutex),
	}
	if m.pivot == nil {
		m.waiters.PushBack(waiter)
	} else {
		m.waiters.InsertBefore(waiter, m.pivot)
	}
	waiter.mtx.Lock()
	return waiter.mtx
}

func (m *RWMutex) releaseWaiter(waiterElement *list.Element[waiter]) {
	if m.pivot == waiterElement {
		m.pivot = waiterElement.Next()
	}
	m.waiters.Remove(waiterElement)
	waiterElement.Value.mtx.Unlock()
}
