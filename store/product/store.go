package product

import (
	"catalog/model"
	"database/sql"
	"errors"
	"strconv"
)

type storage struct{
	Db *sql.DB
}

func New(db *sql.DB) Store {
	return storage{db}
}

func (s storage)GetById(id int)(model.Product,error){
	result:=s.Db.QueryRow("Select Id,Name,BrandId from Product where Id = ?",id)
	var prd model.Product
	err:=result.Scan(&prd.Id,&prd.Name,&prd.Brand.Id)
	if err != nil {
		return model.Product{},errors.New("Product does not exist")
	}
	return prd,nil
}


func (s storage)GetByName(name string)(model.Product,error){
	//log.Println(name)
	result:=s.Db.QueryRow("Select Id,Name,BrandID from Product where Name = ?",name)
	var pr model.Product
	err:=result.Scan(&pr.Id,&pr.Name,&pr.Brand.Id)
	//log.Println(pr)
	if err != nil {
		return model.Product{},errors.New("Product does not exist")
	}
	return pr,nil
}
func (s storage)CreateProduct(pr model.Product)int{
	result,_:=s.Db.Exec("Insert into Product (Name,BrandId) values (?,?)",pr.Name,pr.Brand.Id)
	//if err != nil {
	//	return model.Product{},errors.New("Datatype Mismatch")
	//}
	num,_:=result.LastInsertId()
	return int(num)
}

func (s storage)UpdateProduct(pr model.Product)(error){
	query:="Update Product set"
	flag:=false
	if pr.Name != ""{
		query+=" Name='"
		query+=pr.Name
		flag=true
	}
	if pr.Brand.Id>0 {
		if flag {
			query+="',"
		}
		query+=" BrandId='"
		query+=strconv.Itoa(pr.Brand.Id)
	}
	query+="' where id = ?"
	//log.Println(query,pr.Id)
	_,err:=s.Db.Exec(query,pr.Id )
	//log.Println(err)
	if err != nil {
		return errors.New("Id does not exist")
	}
	return nil
}

func (s storage)DeleteProduct(id int)(error){
	_,err:=s.Db.Exec("Delete from Product where id=?",id)
	if err != nil {
		return errors.New("Id does not exist")
	}
	return nil
}