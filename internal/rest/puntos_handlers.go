package rest

import (
	"encoding/json"
	"net/http"

	"github.com/DieJ6/puntosgo/internal/di"
	"github.com/DieJ6/puntosgo/internal/token"
	"github.com/DieJ6/puntosgo/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PointsHandlers struct {
	Inj *di.Injector
}

func (h PointsHandlers) GetPoints(w http.ResponseWriter, r *http.Request) {
	// Extraer userID del JWT
	userID, err := token.ExtractUserID(r)
	if err != nil {
		http.Error(w, "token inválido", 401)
		return
	}

	uid, _ := primitive.ObjectIDFromHex(userID)

	uc := usecases.ConsultarPuntosUC{
		SaldoSrv: h.Inj.SaldoSrv,
		MvSrv:    h.Inj.MvSrv,
	}

	result, err := uc.Execute(uid)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(result)
}
