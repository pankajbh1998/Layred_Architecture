package store

import "catalog/model"

type  Brand interface {
	GetById(int)(model.Brand,error)
	GetByName(string)(model.Brand,error)
	CreateBrand(model.Brand)(int,error)
	//UpdateBrand(model.Brand)error
	//DeleteBrand(int)error
}

type  Product interface {
	GetById(int)(model.Product,error)
	GetByName(string)([]model.Product,error)
	CreateProduct(model.Product)(int,error)
	UpdateProduct(model.Product)(error)
	DeleteProduct(int)error
}