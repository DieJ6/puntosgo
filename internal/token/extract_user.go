package token

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/tuusuario/puntosgo/internal/env"
)

type authUserResponse struct {
	ID string `json:"id"`
}

func ExtractUserID(r *http.Request) (string, error) {

	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	req, _ := http.NewRequest("GET", env.Get().AuthURL, nil)
	req.Header.Set("Authorization", auth)

	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("error contacting authgo")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("token inválido")
	}

	var data authUserResponse
	json.NewDecoder(resp.Body).Decode(&data)

	if data.ID == "" {
		return "", errors.New("user id not found")
	}

	return data.ID, nil
}
