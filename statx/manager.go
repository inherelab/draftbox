package statx

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// StoreFace interface
type StoreFace interface {
	Del(key string) error
	Get(key string) interface{}
	Set(key string, val interface{}, ttl time.Duration) error
}

// Config definition
type Config struct {
	// IDField define. default is "id"
	IDField string
	// LoginUrl setting. default is "/login"
	LoginUrl  string
	LogoutUrl string
	// redirect to urls
	GuestToUrl  string
	LoggedToUrl string
	LogoutToUrl string
	// KeyPrefix for store user data
	KeyPrefix string
	Excepted  []string
}

// DefaultConfig create
func DefaultConfig() *Config {
	return &Config{
		IDField:   "id",
		LoginUrl:  "/login",
		LogoutUrl: "/logout",
		// redirect to urls
		GuestToUrl:  "/",
		LoggedToUrl: "/",
		LogoutToUrl: "/",
		KeyPrefix:   "user_auth_data",
		Excepted:    []string{"password"},
	}
}

// User object definition. store logged user identity data.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// extra info
	Ext map[string]interface{} `json:"ext"`
	// other
	LoginAt time.Time `json:"login_at"`
}

// SessionID generate
func (u *User) SessionID(prefix string) string {
	sid := genMd5(fmt.Sprintf("%s%s%v", u.ID, u.Name, u.Ext))
	if prefix != "" {
		return prefix + ":" + sid
	}

	return sid
}

// Manager definition
type Manager struct {
	cfg *Config
	//
	Caches map[string]*User
	Store  StoreFace
	// PermChecker permission checker
	PermChecker func(id string, perm string, params interface{}) bool
}

// NewWithConfig create manager
func NewWithConfig(fn func(c *Config)) *Manager {
	m := NewDefault()
	fn(m.cfg)
	return m
}

// NewManager a Manager
func NewManager() *Manager {
	return &Manager{}
}

// NewDefault new a default manager
func NewDefault() *Manager {
	return &Manager{
		cfg: DefaultConfig(),
	}
}

// Login for a user
func (m *Manager) Login(id, name string, ext map[string]interface{}) (sid string, err error) {
	user := &User{ID: id, Name: name, Ext: ext, LoginAt: time.Now()}

	// storage
	sid = user.SessionID(m.cfg.KeyPrefix)
	err = m.Store.Set(sid, user, 0)
	return
}

// Logout
func (m *Manager) Logout(sid string) error {
	return m.Store.Del(sid)
}

// User
func (m *Manager) User(sid string) *User {
	user := m.Store.Get(sid).(*User)

	return user
}

// Can check current user can access the permission
func (m *Manager) ID(uid string) string {
	user := m.Store.Get(m.cfg.SessionKey).(*User)

	return user.ID
}

// Can check current user can access the permission
func (m *Manager) Can(permission string, args interface{}) bool {
	return m.CanAccess(permission, args)
}

// CanAccess check current user can access the permission
func (m *Manager) CanAccess(permission string, args interface{}) bool {
	if m.PermChecker == nil {
		return false
	}

	return m.PermChecker(m.Store.Get("id").(string), perm, args)
}

// genMd5 generate md5 string and length is 32
func genMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}
