package socketcast

import "sync"

type Auth struct {
	sync.RWMutex
	ID            uint64
	Token         string
	Authenticated bool
}

func (a *Auth) IsAuthenticated() bool {
	a.Lock()
	defer a.Unlock()
	return a.Authenticated
}
func (a *Auth) Authenticate() {
	a.Lock()
	defer a.Unlock()
	a.Authenticated = true
}

func (a *Auth) RemoveAuth() {
	a.Lock()
	defer a.Unlock()
	a.Authenticated = false
}

func (a *Auth) GetID() uint64 {
	a.Lock()
	defer a.Unlock()
	return a.ID
}
func (a *Auth) SetId(id uint64) {
	a.Lock()
	defer a.Unlock()
	a.ID = id
}

func (a *Auth) GetToken() string {
	a.Lock()
	defer a.Unlock()
	return a.Token
}
func (a *Auth) SetToken() string {
	a.Lock()
	defer a.Unlock()
	return a.Token
}
