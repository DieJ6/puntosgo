package rest

import (
	"encoding/json"
	"net/http"

	"github.com/DieJ6/puntosgo/internal/di"
	"github.com/DieJ6/puntosgo/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryHandlers struct {
	Inj *di.Injector
}

func (h CategoryHandlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Nombre             string   `json:"nombre"`
		ForKIdEquivalencia string   `json:"forK_id_equivalencia"`
		Prioridad          int      `json:"prioridad"`
		Articulos          []string `json:"articulos"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	eqID, err := primitive.ObjectIDFromHex(body.ForKIdEquivalencia)
	if err != nil {
		http.Error(w, "equivalencia inválida", 400)
		return
	}

	uc := usecases.CrearCategoriaUC{
		CategorySrv: h.Inj.CategorySrv,
	}

	cat, err := uc.Execute(usecases.CrearCategoriaInput{
		Nombre:             body.Nombre,
		ForKIdEquivalencia: eqID,
		Prioridad:          body.Prioridad,
		Articulos:          body.Articulos,
	})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(cat)
}

func (h CategoryHandlers) AddArticle(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IDCategoria string `json:"id_Categoria"`
		IDArticulo  string `json:"id_Article"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	catID, err := primitive.ObjectIDFromHex(body.IDCategoria)
	if err != nil {
		http.Error(w, "id de categoría inválido", 400)
		return
	}

	uc := usecases.AgregarArticuloUC{
		CategorySrv: h.Inj.CategorySrv,
	}

	err = uc.Execute(usecases.AgregarArticuloInput{
		IDCategoria: catID,
		IDArticulo:  body.IDArticulo,
	})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Artículo agregado correctamente",
	})
}

func (h CategoryHandlers) RemoveArticle(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IDCategoria string `json:"id_Categoria"`
		IDArticulo  string `json:"id_Article"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	catID, err := primitive.ObjectIDFromHex(body.IDCategoria)
	if err != nil {
		http.Error(w, "id de categoría inválido", 400)
		return
	}

	uc := usecases.QuitarArticuloUC{
		CategorySrv: h.Inj.CategorySrv,
	}

	err = uc.Execute(usecases.QuitarArticuloInput{
		IDCategoria: catID,
		IDArticulo:  body.IDArticulo,
	})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Artículo removido correctamente",
	})
}
