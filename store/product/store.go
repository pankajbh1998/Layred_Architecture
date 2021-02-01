package product

import (
	"catalog/errors"
	"catalog/model"
	storeInterface "catalog/store"
	"database/sql"
	"strconv"
)

type store struct{
	Db *sql.DB
}

func New(db *sql.DB) storeInterface.Product {
	return store{db}
}

func (s store)GetById(id int)(model.Product,error){
	emptyProduct:=model.Product{}
	result:=s.Db.QueryRow("Select Id,Name,BrandId from Product where Id = ?",id)
	var prd model.Product
	err:=result.Scan(&prd.Id,&prd.Name,&prd.Brand.Id)
	if err != nil {
		return emptyProduct,errors.ProductDoesNotExist
	}
	return prd,nil
}


func (s store)GetByName(name string)([]model.Product,error){
	emptyProduct:=[]model.Product(nil)
	result,err:=s.Db.Query("Select Id,Name,BrandID from Product where Name = ?",name)
	if err != nil {
		return emptyProduct,errors.ProductDoesNotExist
	}
	pr:=[]model.Product(nil)
	for result.Next() {
		var temp model.Product
		_ = result.Scan(&temp.Id, &temp.Name, &temp.Brand.Id)
		pr=append(pr,temp)
	}
	return pr,nil
}
func (s store)CreateProduct(pr model.Product)(int,error){
	result,err:=s.Db.Exec("Insert into Product (Name,BrandId) values (?,?)",pr.Name,pr.Brand.Id)
	if err != nil {
		return 0,errors.ThereIsSomeTechnicalIssue
	}
	num,_:=result.LastInsertId()
	return int(num),err
}

func (s store)UpdateProduct(pr model.Product)error{
	query:="Update Product set"
	flag:=false
	if pr.Name != ""{
		query+=" Name='"+pr.Name+"' "
		flag=true
	}
	if pr.Brand.Id>0 {
		if flag {
			query+=","
		}
		query+=" BrandId='"+strconv.Itoa(pr.Brand.Id) + "' "
	}
	query+="where id = ?"
	_,err:=s.Db.Exec(query,pr.Id )
	if err != nil {
		return errors.PleaseEnterValidData
	}
	return nil
}

func (s store)DeleteProduct(id int)error{
	result,err:=s.Db.Exec("Delete from Product where id=?",id)
	if err != nil {
		return errors.ThereIsSomeTechnicalIssue
	}
	num,_:=result.RowsAffected()
	if num==0{
		return errors.ProductDoesNotExist
	}
	return nil
}