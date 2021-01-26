package product

import "catalog/model"

type  Store interface {
	GetById(int)(model.Product,error)
	GetByName(string)(model.Product,error)
	CreateProduct(model.Product)(int)
	UpdateProduct(model.Product)(error)
	DeleteProduct(int)(error)
}
