package web

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type SessionManager interface {
	// New should create and return a new session
	New() *Session

	// Get should return a cached session
	Get(r *http.Request) (*Session, bool)

	// Save should persist session to the underlying store
	// implementation. Passing a nil session erases it.
	Save(w http.ResponseWriter, r *http.Request, s *Session)
}

func AddTime(t time.Time, duration time.Duration) time.Time {
	return t.Add(duration)
}

type Session struct {
	id      string
	data    *sync.Map
	expires time.Time
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Has(k any) bool {
	_, ok := s.data.Load(k)
	return ok
}

func (s *Session) Set(k, v any) {
	s.data.Store(k, v)
}

func (s *Session) Get(k any) (any, bool) {
	v, ok := s.data.Load(k)
	return v, ok
}

func (s *Session) Del(k any) {
	s.data.Delete(k)
}

func (s *Session) ExpiresIn() int64 {
	return s.expires.Unix() - time.Now().Unix()
}

type AuthUser interface {
	Register(username, password, role string)
	Authenticate(username, password string) (*SystemUser, bool)
}

type SystemUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (s *Session) Register(username, password, role string) {
	s.data.Store(
		"current_user", &SystemUser{
			Username: username,
			Password: password,
			Role:     role,
		},
	)
}

func (s *Session) Authenticate(username, password string) (*SystemUser, bool) {
	v, ok := s.data.Load("current_user")
	if !ok {
		return nil, false
	}
	su, ok := v.(*SystemUser)
	if !ok {
		return nil, false
	}
	if su.Username != username {
		return nil, false
	}
	if su.Password != password {
		return nil, false
	}
	return su, true
}

func (s *Session) GetUser() (*SystemUser, bool) {
	user, found := s.data.Load("current_user")
	if !found {
		return nil, false
	}
	su, ok := user.(*SystemUser)
	if !ok {
		return nil, false
	}
	return su, true
}

// SessionStoreConfig is a configuration object for a
// session manager
type SessionStoreConfig struct {
	SessionID string        `json:"session_id"`       // SessionID is the global session id
	Domain    string        `json:"domain"`           // Domain is the domain to limit the session scope
	Timeout   time.Duration `json:"timeout_duration"` // Timeout is the max idle session time allowed
}

// defaultSessionConfig is pretty self-explanatory
var defaultConfig = &SessionStoreConfig{
	SessionID: "go_sess_id",
	Domain:    "localhost",
	Timeout:   time.Duration(30) * time.Minute,
}

// checkConfig checks the SessionConfig and sets
// and default values that need to be set
func checkSessionStoreConfig(conf *SessionStoreConfig) {
	if conf == nil {
		conf = defaultConfig
	}
	if conf.SessionID == "" {
		conf.SessionID = defaultConfig.SessionID
	}
	if conf.Domain == "" {
		conf.Domain = defaultConfig.Domain
	}
	if conf.Timeout == 0 {
		conf.Timeout = defaultConfig.Timeout
	}
}

var ssOnce sync.Once

var SessionStoreInstance *SessionStore

// SessionStore implements the session manager interface
// and is a basic session manager using cookies.
type SessionStore struct {
	*SessionStoreConfig
	sessions *sync.Map
}

// NewSessionStore takes a *SessionStoreConfig. It initializes and returns
// a *SessionStore. The sid (provided in the config) will be used as the key
// for all session cookies, and the timeout (provided in the config) is the
// maximum allowable idle session time before the session is expired. This
// function uses sync.Once to ensure that only one *SessionStore is alive.
// You may call NewSessionStore() multiple times throughout your application
// and the same instance will be returned.
func NewSessionStore(conf *SessionStoreConfig) *SessionStore {
	if conf == nil {
		conf = defaultConfig
	}
	checkSessionStoreConfig(conf)
	ssOnce.Do(
		func() {
			SessionStoreInstance = initSessionStoreInstance(conf)
		},
	)
	return SessionStoreInstance
}

// initSessionStoreInstance takes a *SessionStoreConfig. It initializes and
// returns a *SessionStore. The sid (provided in the config) will be used as
// the key for all session cookies, and the timeout (provided in the config)
// is the maximum allowable idle session time before the session is expired.
func initSessionStoreInstance(conf *SessionStoreConfig) *SessionStore {
	ss := &SessionStore{
		SessionStoreConfig: conf,
		sessions:           new(sync.Map),
	}
	go ss.gc()
	return ss
}

// New creates and returns a new session
func (ss *SessionStore) New() *Session {
	return &Session{
		id:      RandStringN(32), // create session id 32 chars long
		data:    new(sync.Map),
		expires: AddTime(time.Now(), ss.Timeout),
	}
}

// Get returns a cached session (if one exists)
func (ss *SessionStore) Get(r *http.Request) (*Session, bool) {
	c := getCookie(r, ss.SessionID)
	if c == nil {
		return nil, false
	}
	v, ok := ss.sessions.Load(c.Value)
	if !ok {
		return nil, false
	}
	return v.(*Session), true
}

// MustGet returns a cached session (if one exists), otherwise, it creates and
// returns a new session. It is guaranteed to return a session.
func (ss *SessionStore) MustGet(r *http.Request) (*Session, bool) {
	// First, attempt to get an existing session from the cookie.
	c := getCookie(r, ss.SessionID)
	if c == nil {
		// If no cookie was found, then we must create a new session,
		// and return false (indicating that we did not find an existing
		// one, but one had to be created.) Don't forget, this session
		// WILL NOT be saved unless the Save method is called.
		return &Session{
			id:      RandStringN(32), // create session id 32 chars long
			data:    new(sync.Map),
			expires: AddTime(time.Now(), ss.Timeout),
		}, false
	}
	// Otherwise, we have found a session, but we must check to ensure
	// it is not expired and still exists on the server side.
	v, found := ss.sessions.Load(c.Value)
	if !found {
		// If no session was found, then we must create a new session,
		// and return false (indicating that we did not find an existing
		// one, but one had to be created.) Don't forget, this session
		// WILL NOT be saved unless the Save method is called.
		return &Session{
			id:      RandStringN(32), // create session id 32 chars long
			data:    new(sync.Map),
			expires: AddTime(time.Now(), ss.Timeout),
		}, false
	}
	// Otherwise, we have successfully located an existing session,
	// so we can simply return it along with a true (indicating that
	// the session did in fact exist)
	return v.(*Session), true
}

// Save persists the provided session. If you would like to remove a session, simply
// pass it a nil session, and it will time the cookie out.
func (ss *SessionStore) Save(w http.ResponseWriter, r *http.Request, session *Session) {
	if session == nil {
		cook := getCookie(r, ss.SessionID)
		if cook == nil {
			return
		}
		ss.sessions.Delete(cook.Value)
		http.SetCookie(w, newCookie(ss.SessionID, cook.Value, ss.Domain, time.Now()))
		return
	}
	session.expires = AddTime(time.Now(), ss.Timeout)
	ss.sessions.Store(session.id, session)
	http.SetCookie(w, newCookie(ss.SessionID, session.id, ss.Domain, session.expires))
}

// String is the session store's stringer method
func (ss *SessionStore) String() string {
	var sessions []string
	ss.sessions.Range(
		func(id, sess interface{}) bool {
			sessions = append(sessions, id.(string))
			return true
		},
	)
	return strings.Join(sessions, "\n")
}

// gc is the session store "garbage collector" and
// cleans and disposes of expired sessions (server side)
func (ss *SessionStore) gc() {
	ss.sessions.Range(
		func(id, sess interface{}) bool {
			if sess.(*Session).ExpiresIn() < 0 {
				ss.sessions.Delete(id)
			}
			return true
		},
	)
	time.AfterFunc(ss.Timeout/2, func() { ss.gc() })
}

// newCookie is a helper that wraps the creation of a new
// cookie and returns a filled out *http.Cookie instance
func newCookie(name, value, domain string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:       URLEncode(name),
		Value:      Base64Encode(value),
		Path:       "/",
		Domain:     domain,
		Expires:    expires,
		RawExpires: "",
		MaxAge:     setMaxAge(expires),
		Secure:     false,                   // set to true, if using TLS (false otherwise)
		HttpOnly:   true,                    // protects against XSS attacks
		SameSite:   http.SameSiteStrictMode, // protects against CSRF attacks
		Raw:        "",
		Unparsed:   nil,
	}
}

// setMaxAge calculates the MaxAge (seconds) property
// for a cookie by using the expires time.Time value
// returning the max age in seconds
func setMaxAge(expires time.Time) int {
	max := int(expires.Unix()-time.Now().Unix()) - 1
	if max < 1 {
		return -1
	}
	return max
}

// getCookie is a helper function to help look up and
// return a *http.Cookie by name from the request. It
// returns nil if it no such cookie exists or if it has
// expired in the meantime.
func getCookie(r *http.Request, name string) *http.Cookie {
	c, err := r.Cookie(URLEncode(name))
	if err != nil || err == http.ErrNoCookie {
		return nil
	}
	c.Value = Base64Decode(c.Value)
	return c
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringN creates a random string N characters in length
func RandStringN(n int) string {
	var src = rand.NewSource(time.Now().UnixNano() + int64(rand.Uint64()))
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// Base64Encode takes a plaintext string and returns a base64 encoded string
func Base64Encode(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

// Base64Decode takes a base64 encoded string and returns a plaintext string
func Base64Decode(s string) string {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("cookie: base64 decoding failed %q", err))
	}
	return string(b)
}

// URLEncode takes a plaintext string and returns a URL encoded string
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// URLDecode takes a URL encoded string and returns a plaintext string
func URLDecode(s string) string {
	us, err := url.QueryUnescape(s)
	if err != nil {
		panic(fmt.Sprintf("cookie: query unescape failed %q", err))
	}
	return us
}
