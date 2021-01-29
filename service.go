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

func New(ps product.Store,bs storeB.Store)Service{
	return service{ps,bs}
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
func (s service)GetByName(name string)([]model.Product,error){
	prd,err:=s.storePr.GetByName(name)
	if err != nil{
		return []model.Product(nil),err
	}
	for i,pr:=range prd {
		//log.Println(pr)
		pr.Brand,_=s.storeBr.GetById(pr.Brand.Id)
		//log.Println(pr)
		prd[i].Brand=pr.Brand
	}
	return prd,nil
}

func (s service)CreateProduct(pr model.Product)(model.Product,error){
	br,err:=s.storeBr.GetByName(pr.Brand.Name)
	//log.Println(br)
	if err != nil {
		num,err:=s.storeBr.CreateBrand(pr.Brand)
		if err != nil {
			return model.Product{},err
		}
		br.Id=num
	}
	//log.Println(br)
	pr.Brand=br
	num,err:=s.storePr.CreateProduct(pr)
	if err != nil {
		return model.Product{},err
	}
	//log.Println(pr,num)
	pr,err=s.storePr.GetById(num)
	if err != nil {
		return model.Product{},err
	}
	br,err=s.storeBr.GetById(pr.Brand.Id)
	if err != nil {
		return model.Product{},err
	}
	pr.Brand=br
	return pr,nil
}
func (s service)UpdateProduct(pr model.Product)(model.Product,error){
	brId:=0
	if pr.Brand.Name != "" {
		br,err:=s.storeBr.GetByName(pr.Brand.Name)
		//log.Println(br)
		if err !=nil {
			num,err:=s.storeBr.CreateBrand(pr.Brand)
			if err != nil {
				return model.Product{},err
			}
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

func (s service)DeleteProduct(id int)(error){
	//pr,err:=s.storePr.GetById(id)
	//if err != nil {
	//	return model.Product{},err
	//}
	//br,err:=s.storeBr.GetById(pr.Brand.Id)
	//if err != nil {
	//	return model.Product{},err
	//}
	//pr.Brand=br
	err:=s.storePr.DeleteProduct(id)
	if err != nil {
		return err
	}
	return nil
}