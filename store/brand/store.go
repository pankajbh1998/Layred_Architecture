package brand

import (
	"catalog/errors"
	"catalog/model"
	storeInterface "catalog/store"
	"database/sql"
)

type store struct{
	Db *sql.DB
}

func New(db *sql.DB) storeInterface.Brand{
	return store{db}
}

func (s store)GetById(id int)(model.Brand,error){
	emptyBrand:=model.Brand{}
	result:=s.Db.QueryRow("Select Id,Name from Brand where Id = ?",id)
	var br model.Brand
	err:=result.Scan(&br.Id,&br.Name)
	if err != nil {
		return emptyBrand,errors.BrandDoesNotExist
	}
	return br,nil
}

func (s store)GetByName(name string)(model.Brand,error){
	emptyBrand:=model.Brand{}
	result:=s.Db.QueryRow("Select Id,Name from Brand where Name = ?",name)
	var br model.Brand
	err:=result.Scan(&br.Id,&br.Name)
	if err != nil {
		return emptyBrand,errors.BrandDoesNotExist
	}
	return br,nil
}

func (s store)CreateBrand(br model.Brand)(int,error){
	result,err:=s.Db.Exec("Insert into Brand (Name) values (?)",br.Name)
	if err != nil {
		return 0,errors.ThereIsSomeTechnicalIssue
	}
	num,_:=result.LastInsertId()
	return int(num),nil
}
//
//func (s store)UpdateBrand(br model.Brand)error{
//	_,err:=s.Db.Exec("Update Brand set Name=? where id=?",br.Name,br.Id)
//	if err != nil {
//		return errors.PleaseEnterValidData
//	}
//	return nil
//}
//
//func (s store)DeleteBrand(id int)error{
//	_,err:=s.Db.Exec("Delete from Brand where id=?",id)
//	if err != nil {
//		return errors.BrandDoesNotExist
//	}
//	return nil
//}
//
//
