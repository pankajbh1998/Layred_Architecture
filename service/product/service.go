package product

import (
	"catalog/model"
	serviceInterface "catalog/service"
	"catalog/store"
)

type service struct{
	storePr store.Product
	storeBr store.Brand
}

func New(ps store.Product,bs store.Brand) serviceInterface.Product {
	return service{ps,bs}
}

func (s service)GetById(id int)(model.Product,error){
	prd,err:=s.storePr.GetById(id)
	if err != nil {
		return prd,err
	}
	br,_:=s.storeBr.GetById(prd.Brand.Id)
	prd.Brand=br
	return prd,nil
}
func (s service)GetByName(name string)([]model.Product,error){
	prd,err:=s.storePr.GetByName(name)
	if err != nil{
		return []model.Product(nil),err
	}
	for i,pr:=range prd {
		pr.Brand,_=s.storeBr.GetById(pr.Brand.Id)
		prd[i].Brand=pr.Brand
	}
	return prd,nil
}

func (s service)CreateProduct(pr model.Product)(model.Product,error){
	br,err:=s.storeBr.GetByName(pr.Brand.Name)
	if err != nil {
		num,err:=s.storeBr.CreateBrand(pr.Brand)
		if err != nil {
			return model.Product{},err
		}
		br.Id=num
	}
	pr.Brand=br
	num,err:=s.storePr.CreateProduct(pr)
	if err != nil {
		return model.Product{},err
	}
	pr,err=s.storePr.GetById(num)
	if err != nil {
		return model.Product{},err
	}
	br,_=s.storeBr.GetById(pr.Brand.Id)
	pr.Brand=br
	return pr,nil
}
func (s service)UpdateProduct(pr model.Product)(model.Product,error){
	emptyProduct:=model.Product{}
	if pr.Brand.Name != "" {
		br,err:=s.storeBr.GetByName(pr.Brand.Name)
		if err !=nil {
			br.Id,err =s.storeBr.CreateBrand(pr.Brand)
			if err != nil {
				return emptyProduct, err
			}
		}
		pr.Brand.Id=br.Id
	}
	err:=s.storePr.UpdateProduct(pr)
	if err != nil {
			return emptyProduct,err
		}
	pr,err=s.storePr.GetById(pr.Id)
	if err != nil {
		return emptyProduct,err
	}
	pr.Brand,_=s.storeBr.GetById(pr.Brand.Id)
	return pr,nil
}

func (s service)DeleteProduct(id int)error{
	err:=s.storePr.DeleteProduct(id)
	if err != nil {
		return err
	}
	return nil
}