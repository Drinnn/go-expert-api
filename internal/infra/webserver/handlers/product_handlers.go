package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Drinnn/go-expert-api/internal/dto"
	"github.com/Drinnn/go-expert-api/internal/entity"
	"github.com/Drinnn/go-expert-api/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, req *http.Request) {
	var productInput dto.CreateProductInput
	err := json.NewDecoder(req.Body).Decode(&productInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := entity.NewProduct(productInput.Name, productInput.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
