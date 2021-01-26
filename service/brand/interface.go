package brand

import "catalog/model"

type Service interface {
	GetById(int)(model.Brand,error)
}

