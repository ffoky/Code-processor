package ram_storage

import (
	"container/list"
	"github.com/sirupsen/logrus"
	"http_server/domain"
	"sync"
	"time"
)

type sessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func (st *sessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	st.timeAccessed = time.Now()
	return nil
}

func (st *sessionStore) Get(key interface{}) interface{} {
	if v, ok := st.value[key]; ok {
		st.timeAccessed = time.Now()
		return v
	}
	return nil
}

func (st *sessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	st.timeAccessed = time.Now()
	return nil
}

func (st *sessionStore) SessionID() string {
	return st.sid
}

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewSessionProvider() *Provider {
	return &Provider{
		sessions: make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (p *Provider) SessionInit(sid string) (domain.Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	newsess := &sessionStore{
		sid:          sid,
		timeAccessed: time.Now(),
		value:        make(map[interface{}]interface{}),
	}
	element := p.list.PushFront(newsess)
	p.sessions[sid] = element
	return newsess, nil
}

func (p *Provider) SessionRead(sid string) (domain.Session, error) {
	logrus.WithField("sid", sid).Info("SessionRead: acquiring provider lock")
	p.lock.Lock()

	if element, ok := p.sessions[sid]; ok {
		logrus.WithField("sid", sid).Info("SessionRead: session found in cache")
		p.list.MoveToFront(element)
		sess := element.Value.(*sessionStore)
		sess.timeAccessed = time.Now()
		p.lock.Unlock()
		return sess, nil
	}

	p.lock.Unlock()

	logrus.WithField("sid", sid).Info("SessionRead: session not found, calling SessionInit")
	sess, err := p.SessionInit(sid)
	return sess, err
}

func (p *Provider) SessionDestroy(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if element, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		p.list.Remove(element)
	}
	return nil
}

func (p *Provider) SessionGC(maxlifetime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		element := p.list.Back()
		if element == nil {
			break
		}

		sess := element.Value.(*sessionStore)
		if time.Now().Unix()-sess.timeAccessed.Unix() > maxlifetime {
			p.list.Remove(element)
			delete(p.sessions, sess.sid)
		} else {
			break
		}
	}
}
