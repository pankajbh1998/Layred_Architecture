package product

import (
	"catalog/model"
	storeB "catalog/store/brand"
	"catalog/store/product"
)

type service struct{
	storePr product.Store
	storeBr storeB.Store
}

func New(pr product.Store,br storeB.Store)service{
	return service{pr,br}
}

func (s service)GetById(id int)(model.Product,error){
	prd,err:=s.storePr.GetById(id)
	if err != nil {
		return model.Product{},err
	}
	br,_:=s.storeBr.GetById(prd.Brand.Id)
	prd.Brand=br
	return prd,nil
}
func (s service)GetByName(name string)(model.Product,error){
	prd,err:=s.storePr.GetByName(name)
	if err != nil {
		return model.Product{},err
	}
	br,_:=s.storeBr.GetById(prd.Brand.Id)
	prd.Brand=br
	return prd,nil
}

func (s service)CreateProduct(pr model.Product)(model.Product,error){
	br,err:=s.storeBr.GetByName(pr.Brand.Name)
	if err != nil {
		num:=s.storeBr.CreateBrand(pr.Brand)
		br.Id=num
	}
	//log.Println(br)
	pr.Brand=br
	num:=s.storePr.CreateProduct(pr)
	//log.Println(pr,num)
	pr,err=s.GetById(num)
	if err != nil {
		return model.Product{},err
	}
	return pr,nil
}
func (s service)UpdateProduct(pr model.Product)(model.Product,error){
	brId:=0
	if pr.Brand.Name != "" {
		br,err:=s.storeBr.GetByName(pr.Brand.Name)
		//log.Println(br)
		if err !=nil {
			num:=s.storeBr.CreateBrand(pr.Brand)
			brId=num
		} else {
			brId=br.Id
		}
		//log.Println(brId)
	}
	pr.Brand.Id=brId
	err:=s.storePr.UpdateProduct(pr)
	if err != nil {
			return model.Product{},err
		}
	pr,err=s.storePr.GetById(pr.Id)
	if err != nil {
		return model.Product{},err
	}
	pr.Brand,err=s.storeBr.GetById(pr.Brand.Id)
	if err != nil {
		return model.Product{},err
	}
	return pr,nil
}

func (s service)DeleteProduct(id int)(model.Product,error){
	pr,err:=s.storePr.GetById(id)
	if err != nil {
		return model.Product{},err
	}
	err=s.storePr.DeleteProduct(id)
	if err != nil {
		return model.Product{},err
	}
	return pr,nil
}