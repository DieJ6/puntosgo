package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/DieJ6/puntosgo/internal/di"
)

type ctxKey string

const ctxUser ctxKey = "auth_user"

type AuthUser struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Login       string   `json:"login"`
	Permissions []string `json:"permissions"`
	Enabled     bool     `json:"enabled"`
}

func isAdmin(u *AuthUser) bool {
	if u == nil || !u.Enabled {
		return false
	}
	for _, p := range u.Permissions {
		if strings.EqualFold(p, "admin") {
			return true
		}
	}
	return false
}

func isUser(u *AuthUser) bool {
	if u == nil || !u.Enabled {
		return false
	}
	// admin cuenta como usuario logueado
	if isAdmin(u) {
		return true
	}
	for _, p := range u.Permissions {
		if strings.EqualFold(p, "user") {
			return true
		}
	}
	return false
}

func extractAuthHeader(r *http.Request) (rawHeader string, tokenOnly string, err error) {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	if raw == "" {
		return "", "", errors.New("missing Authorization header")
	}

	// Si viene "Bearer xxx", sacamos el token.
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "bearer ") {
		tok := strings.TrimSpace(raw[7:])
		if tok == "" {
			return "", "", errors.New("empty bearer token")
		}
		return raw, tok, nil
	}

	// Si viene token plano (p.ej ObjectID), lo aceptamos también
	return raw, raw, nil
}

func RequireAuth(inj *di.Injector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rawAuth, tokenOnly, err := extractAuthHeader(r)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			client := &http.Client{Timeout: 5 * time.Second}

			// 1) Primer intento: forward tal cual vino del cliente
			u, ok := fetchCurrentUser(r.Context(), client, inj.AuthURL, rawAuth)
			if !ok {
				// 2) Si el cliente mandó Bearer y authgo espera token plano, reintentamos solo con el token
				if strings.HasPrefix(strings.ToLower(rawAuth), "bearer ") {
					u2, ok2 := fetchCurrentUser(r.Context(), client, inj.AuthURL, tokenOnly)
					if !ok2 {
						http.Error(w, "unauthorized", http.StatusUnauthorized)
						return
					}
					u = u2
				} else {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ctxUser, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func fetchCurrentUser(ctx context.Context, client *http.Client, url string, authHeader string) (*AuthUser, bool) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, false
	}

	var u AuthUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, false
	}
	return &u, true
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := r.Context().Value(ctxUser).(*AuthUser)
		if !isAdmin(u) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := r.Context().Value(ctxUser).(*AuthUser)
		if !isUser(u) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CurrentUser(r *http.Request) *AuthUser {
	u, _ := r.Context().Value(ctxUser).(*AuthUser)
	return u
}
