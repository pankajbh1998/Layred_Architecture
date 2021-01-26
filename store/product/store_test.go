package product

import (
	"catalog/model"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
)

func TestStoreGetByID(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input int
		output model.Product
		err error
	}{
		{1,model.Product{1,"Ref",model.Brand{1,""}},nil},
		{2,model.Product{},errors.New("Product does not exist")},
	}
	str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(str).AddRow(tc.output.Id,tc.output.Name,tc.output.Brand.Id)
		mock.ExpectQuery("Select Id,Name,BrandId from Product where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)
		result,err:=store.GetById(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Fatalf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.err,err)
			}
		} else if tc.output !=result {
			t.Fatalf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		log.Printf("Passed at %v",i+1)
	}
}


func TestStoreGetByName(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input string
		output model.Product
		err error
	}{
		{"Ref",model.Product{1,"Ref",model.Brand{1,""}},nil},
		{"",model.Product{},errors.New("Product does not exist")},
	}
	str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(str).AddRow(tc.output.Id,tc.output.Name,tc.output.Brand.Id)
		mock.ExpectQuery("Select Id,Name,BrandID from Product where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)
		result,err:=store.GetByName(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Fatalf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.err,err)
			}
		} else if tc.output !=result {
			t.Fatalf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		log.Printf("Passed at %v",i+1)
	}
}

func TestStoreCreateProduct(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Product
		output int
	}{
		{model.Product{1,"Ref",model.Brand{1,""}},1},
		{model.Product{2,"Washing",model.Brand{1,""}},2},
	}
	//str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Insert into Product*").WithArgs(tc.input.Name,tc.input.Brand.Id).WillReturnResult(sqlmock.NewResult(int64(tc.output),1))
		result:=store.CreateProduct(tc.input)
		if tc.output !=result {
			t.Fatalf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		log.Printf("Passed at %v",i+1)
	}
}

func TestStoreUpateProduct(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Product
		rowsAffected int64
		outputErr error
	}{
		{model.Product{1,"Ref",model.Brand{}},1,nil},
		{model.Product{2,"",model.Brand{1,""}},1,nil},
		{model.Product{3,"Ref",model.Brand{1,""}},0,errors.New("Id does not exist")},
	}
	//str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Update Product set*").WithArgs(tc.input.Id).WillReturnError(tc.outputErr).WillReturnResult(sqlmock.NewResult(0,tc.rowsAffected))
		err:=store.UpdateProduct(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Fatalf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		log.Printf("Passed at %v",i+1)
	}
}


func TestStoreDeleteProduct(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input int
		rowsAffected int64
		outputErr error
	}{
		{1,1,nil},
		{2,1,nil},
		{3,0,errors.New("Id does not exist")},
	}
	//str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Delete from Product*").WithArgs(tc.input).WillReturnError(tc.outputErr).WillReturnResult(sqlmock.NewResult(0,tc.rowsAffected))
		err:=store.DeleteProduct(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Fatalf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		log.Printf("Passed at %v",i+1)
	}
}