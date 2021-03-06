package service

import "catalog/model"

type Product interface {
	GetById(int)(model.Product,error)
	GetByName(string)([]model.Product,error)
	CreateProduct(model.Product)(model.Product,error)
	UpdateProduct(model.Product)(model.Product,error)
	DeleteProduct(int)error
}

