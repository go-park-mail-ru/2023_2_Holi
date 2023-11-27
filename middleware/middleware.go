package middleware

import (
	"2023_2_Holi/domain/grpc/session"
	logs "2023_2_Holi/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	req_context "github.com/gorilla/context"

	"2023_2_Holi/domain"

	"github.com/sirupsen/logrus"
)

type Middleware struct {
	AuthClient session.AuthCheckerClient
	Token      *domain.HashToken
	metrics    *Metrics
}

func InitMiddleware(authCl session.AuthCheckerClient, token *domain.HashToken) *Middleware {
	return &Middleware{
		AuthClient: authCl,
		Token:      token,
		metrics:    NewMetrics(),
	}
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
				if errors.Is(err, http.ErrNoCookie) {
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
			if errors.Is(err, http.ErrNoCookie) {
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
		fmt.Println("-----------")
		fmt.Println(userID)
		fmt.Println("-----------")

		req_context.Set(r, "userID", userID.ID)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		m.metrics.workTime.WithLabelValues(strconv.Itoa(lrw.statusCode), r.URL.Path).Observe(float64(time.Since(start).Milliseconds()))
		m.metrics.hits.WithLabelValues(strconv.Itoa(lrw.statusCode), r.URL.Path).Inc()
	})
}

type AccessLogger struct {
	LogrusLogger *logrus.Logger
}

func (ac *AccessLogger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		ac.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Header.Get("Request-ID"),
			"work_time":   time.Since(start),
			"status":      lrw.statusCode,
		}).Info(r.URL.Path)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
