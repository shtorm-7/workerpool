package callbackfield

import (
	"sync"

	list "github.com/bahlo/generic-list-go"
)

type (
	waiter struct {
		rLock bool
		await chan struct{}
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
	await := m.awaitWaiter(false)
	m.mtx.Unlock()
	<-await
}

func (m *RWMutex) RLock() {
	m.mtx.Lock()
	if m.locks == 0 || (m.locks > 0 && m.waiters.Front() == nil) {
		m.locks++
		m.mtx.Unlock()
		return
	}
	await := m.awaitWaiter(true)
	m.mtx.Unlock()
	<-await
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
	await := m.awaitPriorityLockWaiter()
	m.locks--
	if m.locks == 0 {
		m.notifyWaiters()
	}
	m.mtx.Unlock()
	<-await
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

func (m *RWMutex) awaitWaiter(rLock bool) <-chan struct{} {
	waiter := waiter{
		rLock: rLock,
		await: make(chan struct{}),
	}
	if m.pivot == nil {
		m.pivot = m.waiters.PushBack(waiter)
	} else {
		m.waiters.PushBack(waiter)
	}
	return waiter.await
}

func (m *RWMutex) awaitPriorityLockWaiter() <-chan struct{} {
	waiter := waiter{
		await: make(chan struct{}),
	}
	if m.pivot == nil {
		m.waiters.PushBack(waiter)
	} else {
		m.waiters.InsertBefore(waiter, m.pivot)
	}
	return waiter.await
}

func (m *RWMutex) releaseWaiter(waiterElement *list.Element[waiter]) {
	if m.pivot == waiterElement {
		m.pivot = waiterElement.Next()
	}
	m.waiters.Remove(waiterElement)
	close(waiterElement.Value.await)
}
