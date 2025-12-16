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

	if body.IDArticulo == "" {
		http.Error(w, "id de producto inválido", 400)
		return
	}

	catID, err := primitive.ObjectIDFromHex(body.IDCategoria)
	if err != nil {
		http.Error(w, "id de categoría inválido", 400)
		return
	}

	// 1) La categoría debe existir
	targetCat, err := h.Inj.CategorySrv.GetByID(catID)
	if err != nil || targetCat == nil {
		http.Error(w, "la categoría no existe", 404)
		return
	}

	// 2) Buscar en todas las categorías si el artículo ya está asignado
	cats, err := h.Inj.CategoryRepo.FindAll()
	if err != nil {
		http.Error(w, "error buscando categorías", 500)
		return
	}

	for _, c := range cats {
		if c.ID == catID {
			continue // es la misma categoría destino
		}
		// ¿está el id del producto en esta categoría?
		for _, a := range c.Articulos {
			if a == body.IDArticulo {
				// Camino alternativo: ya está asignado → avisar (409)
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]any{
					"Message":               "El artículo ya está asignado a otra categoría. Quitalo primero o elegí otra categoría.",
					"categoria_existente":   c.ID.Hex(),
					"categoria_existente_nombre": c.Nombre,
				})
				return
			}
		}
	}

	// 3) Agregar
	if err := h.Inj.CategorySrv.AddArticulo(catID, body.IDArticulo); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"Message": "Artículo agregado correctamente",
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

	if body.IDArticulo == "" {
		http.Error(w, "id de producto inválido", 400)
		return
	}

	catID, err := primitive.ObjectIDFromHex(body.IDCategoria)
	if err != nil {
		http.Error(w, "id de categoría inválido", 400)
		return
	}

	// categoría debe existir
	cat, err := h.Inj.CategorySrv.GetByID(catID)
	if err != nil || cat == nil {
		http.Error(w, "la categoría no existe", 404)
		return
	}

	err = h.Inj.CategorySrv.RemoveArticulo(catID, body.IDArticulo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"Message": "Artículo removido correctamente",
	})
}

