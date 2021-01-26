package brand

import "catalog/model"

type  Store interface {
	GetById(int)(model.Brand,error)
	GetByName(string)(model.Brand,error)
	CreateBrand(model.Brand)(int)
	UpdateBrand(model.Brand)error
	DeleteBrand(int)(error)
}