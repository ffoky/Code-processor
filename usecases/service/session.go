package service

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"http_server/domain"
	"io"
	"net/http"
	"net/url"
)

type SessionService struct {
	manager *domain.SessionManager
}

func NewSessionService(provider domain.SessionProvider, cookieName string, maxLifetime int64) *SessionService {
	return &SessionService{
		manager: &domain.SessionManager{
			Provider:    provider,
			CookieName:  cookieName,
			MaxLifetime: maxLifetime,
		},
	}
}

func (s *SessionService) sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (s *SessionService) Start(w http.ResponseWriter, r *http.Request) (domain.Session, error) {
	logrus.Info("SessionService: acquiring lock")
	s.manager.Lock.Lock()
	defer s.manager.Lock.Unlock()
	logrus.Info("SessionService: lock acquired")

	cookie, err := r.Cookie(s.manager.CookieName)
	if err != nil || cookie.Value == "" {
		logrus.Info("SessionService: no valid cookie, creating new session")

		sid := s.sessionID()
		logrus.WithField("sid", sid).Info("SessionService: generated session ID")

		session, err := s.manager.Provider.SessionInit(sid)
		if err != nil {
			logrus.WithError(err).Error("SessionService: SessionInit failed")
			return nil, err
		}

		http.SetCookie(w, &http.Cookie{
			Name:     s.manager.CookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(s.manager.MaxLifetime),
		})
		logrus.Info("SessionService: new session initialized and cookie set")
		return session, nil
	}

	sid, _ := url.QueryUnescape(cookie.Value)
	logrus.WithField("sid", sid).Info("SessionService: found cookie, reading session")
	session, err := s.manager.Provider.SessionRead(sid)
	if err != nil {
		logrus.WithError(err).Error("SessionService: SessionRead failed")
	}
	return session, err
}

func (s *SessionService) Get(sid string) (domain.Session, error) {
	logrus.WithField("sid", sid).Info("SessionService: GetByID called")
	s.manager.Lock.Lock()
	defer s.manager.Lock.Unlock()

	return s.manager.Provider.SessionRead(sid)
}
