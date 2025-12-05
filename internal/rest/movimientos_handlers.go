package rest

import (
	"encoding/json"
	"net/http"

	"github.com/tuusuario/puntosgo/internal/di"
	"github.com/tuusuario/puntosgo/internal/token"
	"github.com/tuusuario/puntosgo/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovHandlers struct {
	Inj *di.Injector
}

func (h MovHandlers) GetMovements(w http.ResponseWriter, r *http.Request) {

	userID, err := token.ExtractUserID(r)
	if err != nil {
		http.Error(w, "token inválido", 401)
		return
	}

	uid, _ := primitive.ObjectIDFromHex(userID)

	uc := usecases.ConsultarMovimientosUC{
		MvSrv: h.Inj.MvSrv,
	}

	movs, err := uc.Execute(uid)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(movs)
}
