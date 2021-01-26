package brand

import (
	"catalog/model"
	"database/sql"
	"errors"
)

type storage struct{
	Db *sql.DB
}

func New(db *sql.DB) Store{
	return storage{db}
}

func (s storage)GetById(id int)(model.Brand,error){
	result:=s.Db.QueryRow("Select Id,Name from Brand where Id = ?",id)
	var br model.Brand
	err:=result.Scan(&br.Id,&br.Name)
	if err != nil {
		return model.Brand{},errors.New("Brand does not exist")
	}
	//log.Println(br)
	return br,nil
}

func (s storage)GetByName(name string)(model.Brand,error){
	result:=s.Db.QueryRow("Select Id,Name from Brand where Name = ?",name)
	var br model.Brand
	err:=result.Scan(&br.Id,&br.Name)
	if err != nil {
		return model.Brand{},errors.New("Brand does not exist")
	}
	return br,nil
}

func (s storage)CreateBrand(br model.Brand)(int){
	result,_:=s.Db.Exec("Insert into Brand (Name) values (?)",br.Name)
	//if err != nil {
	//	return 0,errors.New("Datatype Mismatch")
	//}
	num,_:=result.LastInsertId()
	return int(num)
}

func (s storage)UpdateBrand(br model.Brand)(error){
	_,err:=s.Db.Exec("Update Brand set Name=? where id=?",br.Name,br.Id)
	////log.Println(result,err)
	if err != nil {
		return errors.New("Id does not exist")
	}
	return nil
}

func (s storage)DeleteBrand(id int)(error){
	_,err:=s.Db.Exec("Delete from Brand where id=?",id)
	if err != nil {
		return errors.New("Id does not exist")
	}
	return nil
}


