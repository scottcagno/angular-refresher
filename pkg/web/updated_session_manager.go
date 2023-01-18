package web

import (
	"context"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type SessionID [32]byte

func MakeSessionID() SessionID {
	var sid SessionID
	for i := range sid {
		sid[i] = letters[rand.Intn(len(letters))]
	}
	return sid
}

type SessionT struct {
	id      SessionID
	data    *sync.Map
	expires time.Time
}

func (s *SessionT) ID() SessionID {
	return s.id
}

func (s *SessionT) ExpiresIn() time.Duration {
	return time.Until(s.expires)
}

func (s *SessionT) IsExpired() bool {
	return time.Until(s.expires) < 1
}

func (s *SessionT) Set(k, v string) {
	s.data.Store(k, v)
}

func (s *SessionT) Get(k string) (string, bool) {
	v, ok := s.data.Load(k)
	return v.(string), ok
}

func (s *SessionT) Del(k string) {
	s.data.Delete(k)
}

var singleton sync.Once

var DefaultSessionStore *SessionStoreT

type SessionStoreT struct {
	storeKey string
	timeout  time.Duration
	sessions *sync.Map
	ticker   *time.Ticker
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewSessionStoreManager(key string, timeout time.Duration) *SessionStoreT {
	singleton.Do(func() {
		DefaultSessionStore = initSessionStore(key, timeout)
	})
	return DefaultSessionStore
}

const (
	minimumTimeout = 1 * time.Minute
	tickerDefault  = 10 * time.Second
)

func initSessionStore(key string, timeout time.Duration) *SessionStoreT {
	if timeout < minimumTimeout {
		timeout = minimumTimeout
	}
	ctx, cancel := context.WithCancel(context.Background())
	ss := &SessionStoreT{
		storeKey: key,
		timeout:  timeout,
		sessions: new(sync.Map),
		ticker:   time.NewTicker(tickerDefault),
		ctx:      ctx,
		cancel:   cancel,
	}
	go ss.cleanUpRoutine()
	runtime.SetFinalizer(ss, (*SessionStoreT).Close)
	return ss
}

func (ss *SessionStoreT) NewSession() *SessionT {
	return &SessionT{
		id:      MakeSessionID(),
		data:    new(sync.Map),
		expires: time.Now().Add(ss.timeout),
	}
}

// GetSession takes a SessionID and attempts to locate the
// matching *Session. If a matching *Session can be found
// it is returned along with a boolean indicating weather or
// not the session was found.
func (ss *SessionStoreT) GetSession(sid SessionID) (*SessionT, bool) {
	session, found := ss.sessions.Load(sid)
	if !found {
		return nil, false
	}
	return session.(*SessionT), true
}

func (ss *SessionStoreT) SaveSession(session *SessionT) {
	// If the session is nil, do nothing the checkForExpiredSessions
	// will handle any of the extra cleanup necessary.
	if session == nil {
		return
	}
	// If the session is expired, remove it and return.
	if session.IsExpired() {
		ss.sessions.Delete(session.id)
		return
	}
	// Otherwise, bump the expiry time and save it.
	session.expires = time.Now().Add(ss.timeout)
	ss.sessions.Store(session.id, session)
}

func (ss *SessionStoreT) cleanUpRoutine() {
	// When we receive a "tick", we should loop through the
	// sessions, checking to see if any of them are expired.
	// If we find any that are expired, we should remove them.
	for {
		select {
		case t := <-ss.ticker.C:
			// Clean up expired sessions
			log.Printf("Checking for expired sessions: %v\n", t)
			ss.sessions.Range(func(sid, session any) bool {
				if session.(*SessionT).IsExpired() {
					ss.sessions.Delete(sid)
				}
				return true
			})
		case <-ss.ctx.Done():
			ss.ticker.Stop()
			return
		}
	}
}

func (ss *SessionStoreT) Close() {
	log.Printf("*SessionStore.Close has been called.\n")
	// stop the ticker and free any other
	// resources.
	ss.ticker.Stop()
	ss.cancel()
}

func (ss *SessionStoreT) GetSessionCount() int {
	var sessionCount int
	ss.sessions.Range(func(k, v any) bool {
		if k != nil && v != nil {
			sessionCount++
		}
		return true
	})
	return sessionCount
}
