package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Drinnn/go-expert-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Notebook", 16000.0)
	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.Nil(t, err)

	var productFound *entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, productFound.ID, product.ID)
	assert.Equal(t, productFound.Name, product.Name)
	assert.Equal(t, productFound.Price, product.Price)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	for i := 1; i <= 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		if err != nil {
			t.Error(err)
		}
		db.Create(product)
	}

	productDB := NewProduct(db)

	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 4)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 24", products[3].Name)
}

func TestFindProductByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Notebook", 16000.0)
	if err != nil {
		t.Error(err)
	}
	err = db.Create(product).Error
	if err != nil {
		t.Error(err)
	}

	productDB := NewProduct(db)
	foundProduct, err := productDB.FindById(product.ID.String())

	assert.Nil(t, err)
	assert.Equal(t, product.ID, foundProduct.ID)
	assert.Equal(t, product.Name, foundProduct.Name)
	assert.Equal(t, product.Price, foundProduct.Price)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Notebook", 16000.0)
	if err != nil {
		t.Error(err)
	}
	err = db.Create(product).Error
	if err != nil {
		t.Error(err)
	}

	product.Name = "Smartphone"
	product.Price = 13000.0
	productDB := NewProduct(db)
	err = productDB.Update(product)
	assert.Nil(t, err)

	err = db.First(&product, "id = ?", product.ID).Error
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "Smartphone", product.Name)
	assert.Equal(t, 13000.0, product.Price)
}

func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Notebook", 16000.0)
	if err != nil {
		t.Error(err)
	}
	err = db.Create(product).Error
	if err != nil {
		t.Error(err)
	}

	productDB := NewProduct(db)
	err = productDB.Delete(product.ID.String())
	assert.Nil(t, err)

	err = db.First(&product, "id = ?", product.ID).Error
	assert.Error(t, err)
}
