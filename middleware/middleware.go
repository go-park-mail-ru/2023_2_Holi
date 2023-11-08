package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"2023_2_Holi/domain"
)

type Middleware struct {
	AuthUsecase domain.AuthUsecase
	Token       *domain.HashToken
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, User-Agent, X-CSRF-TOKEN")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (m *Middleware) CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			csrfToken := r.Header.Get("X-CSRF-TOKEN")
			validCSRFToken, err := m.Token.Check(uuid.NewString(), csrfToken)
			if err != nil || !validCSRFToken {
				http.Error(w, `{"err":"invalid CSRF token"}`, http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusUnauthorized)
				return
			}

			http.Error(w, `{"err":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		if c.Expires.After(time.Now()) {
			http.Error(w, `{"err":"cookie is expired"}`, http.StatusUnauthorized)
		}
		sessionToken := c.Value
		exists, err := m.AuthUsecase.IsAuth(sessionToken)
		if err != nil {
			http.Error(w, `{"err":"`+domain.ErrInternalServerError.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, `{"err":"`+domain.ErrUnauthorized.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func InitMiddleware(authUsecase domain.AuthUsecase) *Middleware {
	return &Middleware{AuthUsecase: authUsecase}
}

type AccessLogger struct {
	LogrusLogger *logrus.Logger
}

func (ac *AccessLogger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		ac.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}
