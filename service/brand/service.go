package brand

import (
	"catalog/model"
	storeB "catalog/store/brand"
)

type service struct{
	store storeB.Store
}

func New(prd storeB.Store)service{
	return service{prd}
}


func (s service)GetById(id int)(model.Brand,error){
	br,err:=s.store.GetById(id)
	if err != nil {
		return model.Brand{},err
	}
	return br,nil
}
