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
		http.Error(w, "Equivalencia inválida", 400)
		return
	}

	// Validar existencia
	eq, err := h.Inj.EquivSrv.GetByID(eqID)
	if err != nil || eq == nil {
		http.Error(w, "La equivalencia no existe", http.StatusBadRequest)
		return
	}

	uc := usecases.CrearCategoriaUC{ CategorySrv: h.Inj.CategorySrv }

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

	// 1) parse body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.IDArticulo == "" {
		http.Error(w, "Id de producto inválido", http.StatusBadRequest)
		return
	}

	// 2) Validar existencia de articulo
	authHeader := r.Header.Get("Authorization")

	exists, err := h.Inj.Catalog.Exists(body.IDArticulo, authHeader)
	if err != nil {
		http.Error(w, "Error consultando catálogo", http.StatusBadGateway)
		return
	}
	if !exists {
		http.Error(w, "El artículo no existe", http.StatusBadRequest)
		return
	}

	// 3) validar ObjectID categoría
	catID, err := primitive.ObjectIDFromHex(body.IDCategoria)
	if err != nil {
		http.Error(w, "Id de categoría inválido", http.StatusBadRequest)
		return
	}

	// 4) La categoría debe existir (precondición)
	targetCat, err := h.Inj.CategorySrv.GetByID(catID)
	if err != nil || targetCat == nil {
		http.Error(w, "La categoría no existe", http.StatusNotFound)
		return
	}

	// 5) UC: agrega pero primero valida que no esté asignado a otra categoría
	uc := usecases.AgregarArticuloUC{
		CategorySrv: h.Inj.CategorySrv,
	}

	err = uc.Execute(usecases.AgregarArticuloInput{
		IDCategoria: catID,
		IDArticulo:  body.IDArticulo,
	})
	if err != nil {
		if err == usecases.ErrArticuloYaAsignado {
			http.Error(w, err.Error(), http.StatusConflict) // 409
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	if body.IDArticulo == "" {
		http.Error(w, "Id de producto inválido", 400)
		return
	}

	catID, err := primitive.ObjectIDFromHex(body.IDCategoria)
	if err != nil {
		http.Error(w, "Id de categoría inválido", 400)
		return
	}

	// categoría debe existir
	cat, err := h.Inj.CategorySrv.GetByID(catID)
	if err != nil || cat == nil {
		http.Error(w, "La categoría no existe", 404)
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

