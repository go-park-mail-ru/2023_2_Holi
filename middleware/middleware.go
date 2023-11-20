package middleware

import (
	"2023_2_Holi/domain/grpc/session"
	logs "2023_2_Holi/logger"
	"context"
	"net/http"
	"time"

	req_context "github.com/gorilla/context"

	"2023_2_Holi/domain"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	AuthClient session.AuthCheckerClient
	Token      *domain.HashToken
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
			headerCsrfToken := r.Header.Get("X-CSRF-TOKEN")
			cookieCsrfToken, err := r.Cookie("csrf-token")
			if err != nil {
				if err == http.ErrNoCookie {
					domain.WriteError(w, err.Error(), http.StatusUnauthorized)
					logs.LogError(logs.Logger, "middleware", "CSRFProtection", err, err.Error())
					return
				}

				domain.WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}

			validCSRFToken, err := m.Token.Check(headerCsrfToken, cookieCsrfToken.Value)
			if err != nil {
				domain.WriteError(w, err.Error(), http.StatusForbidden)
				return
			}
			if !validCSRFToken {
				domain.WriteError(w, "invalid CSRF token", http.StatusForbidden)
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
				domain.WriteError(w, err.Error(), http.StatusUnauthorized)
				return
			}

			domain.WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if c.Expires.After(time.Now()) {
			domain.WriteError(w, "cookie is expired", http.StatusUnauthorized)
		}
		sessionToken := c.Value
		userID, err := m.AuthClient.IsAuth(
			context.Background(),
			&session.Token{
				Token: sessionToken,
			})
		if err != nil {
			domain.WriteError(w, err.Error(), domain.GetHttpStatusCode(err))
			return
		}
		if userID.ID == "" {
			domain.WriteError(w, domain.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		req_context.Set(r, "userID", userID.ID)
		next.ServeHTTP(w, r)
	})
}

func InitMiddleware(authCl session.AuthCheckerClient, token *domain.HashToken) *Middleware {
	return &Middleware{
		AuthClient: authCl,
		Token:      token,
	}
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
