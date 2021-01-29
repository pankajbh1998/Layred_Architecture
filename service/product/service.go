package product

import (
	"catalog/model"
	"catalog/store/brand"
	"catalog/store/product"
)

type service struct{
	storePr product.Store
	storeBr brand.Store
}

func New(ps product.Store,bs brand.Store)Service{
	return service{ps,bs}
}

func (s service)GetById(id int)(model.Product,error){
	emptyProduct:=model.Product{}
	prd,err:=s.storePr.GetById(id)
	if err != nil {
		return emptyProduct,err
	}

	br,_:=s.storeBr.GetById(prd.Brand.Id)
	prd.Brand=br
	return prd,nil
}
func (s service)GetByName(name string)([]model.Product,error){
	emptyProduct:=[]model.Product(nil)
	prd,err:=s.storePr.GetByName(name)
	if err != nil{
		return emptyProduct,err
	}
	for i,pr:=range prd {
		pr.Brand,_=s.storeBr.GetById(pr.Brand.Id)
		prd[i].Brand=pr.Brand
	}
	return prd,nil
}

func (s service)CreateProduct(pr model.Product)(model.Product,error){
	emptyProduct:=model.Product{}
	br,err:=s.storeBr.GetByName(pr.Brand.Name)
	if err!=nil {
		br.Id,err=s.storeBr.CreateBrand(pr.Brand)
		if err != nil {
			return emptyProduct,err
		}
	}
	pr.Brand=br
	num,err:=s.storePr.CreateProduct(pr)
	if err != nil {
		return emptyProduct,err
	}
	pr,_=s.storePr.GetById(num)
	pr.Brand,_=s.storeBr.GetById(pr.Brand.Id)
	return pr,nil
}
func (s service)UpdateProduct(pr model.Product)(model.Product,error){
	br:=model.Brand{}
	if pr.Brand.Name != "" {
		var err error
		br,err=s.storeBr.GetByName(pr.Brand.Name)
		if err !=nil {
			br.Id,err=s.storeBr.CreateBrand(pr.Brand)
			if err != nil {
				return model.Product{}, err
			}
		}
	}
	pr.Brand=br
	err:=s.storePr.UpdateProduct(pr)
	if err != nil {
			return model.Product{},err
		}
	pr,_=s.storePr.GetById(pr.Id)
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