package rest

import (
	"encoding/json"
	"net/http"

	"github.com/DieJ6/puntosgo/internal/di"
	"github.com/DieJ6/puntosgo/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PointsHandlers struct {
	Inj *di.Injector
}

func (h PointsHandlers) GetPoints(w http.ResponseWriter, r *http.Request) {
	u, _ := r.Context().Value(ctxUser).(*AuthUser)
	if u == nil || u.ID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uc := usecases.ConsultarPuntosUC{
		SaldoSrv: h.Inj.SaldoSrv,
		MvSrv:    h.Inj.MvSrv,
	}

	result, err := uc.Execute(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}