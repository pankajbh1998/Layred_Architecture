package product

import (
	"catalog/errors"
	"catalog/model"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestStoreGetById(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	var (
		testCases = []struct {
			input  int
			output model.Product
			err    error
		}{
			{input: 1, output: model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1}}},
			{input: 2, err: errors.ProductDoesNotExist},
		}
	)
	str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(str).AddRow(tc.output.Id,tc.output.Name,tc.output.Brand.Id)
		mock.ExpectQuery("Select Id,Name,BrandId from Product where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)
		result,err:=store.GetById(tc.input)
		if tc.err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.err,err)
			}
		} else if tc.output !=result {
			t.Errorf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		
	}
}


func TestStoreGetByName(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input string
		output []model.Product
		err error
	}{
		{"Ref",[]model.Product{{1,"Ref",model.Brand{1,""}}},nil},
		{"",[]model.Product(nil),errors.ProductDoesNotExist},
	}
	str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(nil).AddRow()
		if tc.output != nil {
			row = sqlmock.NewRows(str).AddRow(tc.output[0].Id, tc.output[0].Name, tc.output[0].Brand.Id)
		}
		mock.ExpectQuery("Select Id,Name,BrandID from Product where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)

		result,err:=store.GetByName(tc.input)

		if tc.err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.err,err)
			}
		} else if ! reflect.DeepEqual(tc.output , result ){
			t.Errorf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		
	}
}

func TestStoreCreateProduct(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Product
		output int
		err error
	}{
		{input: model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1}}, output: 1},
		{err: errors.ThereIsSomeTechnicalIssue},
		{input: model.Product{Id: 2, Name: "Washing", Brand: model.Brand{Id: 1}}, output: 2},
	}
	for i,tc:=range testCases{
		mock.ExpectExec("Insert into Product*").WithArgs(tc.input.Name,tc.input.Brand.Id).WillReturnResult(sqlmock.NewResult(int64(tc.output),1)).WillReturnError(tc.err)

		result,err:=store.CreateProduct(tc.input)
		if tc.err != nil {
			if !reflect.DeepEqual(tc.err , err ){
				t.Error(err)
			}
		}else if tc.output !=result {
			t.Errorf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		
	}

}

func TestStoreUpdateProduct(t *testing.T){

	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Product
		rowsAffected int64
		outputErr error
	}{
		{input: model.Product{Id: 1, Name: "Ref"}, rowsAffected: 1},
		{input: model.Product{Id: 2, Brand: model.Brand{Id: 1}}, rowsAffected: 1},
		{input: model.Product{Id: 3, Name: "Ref", Brand: model.Brand{Id: 1}}, outputErr: errors.PleaseEnterValidData},
	}
	//str:=[]string{"id","name","brand name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Update Product set*").WithArgs(tc.input.Id).WillReturnError(tc.outputErr).WillReturnResult(sqlmock.NewResult(0,tc.rowsAffected))

		err:=store.UpdateProduct(tc.input)
		if tc.outputErr != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		
	}
	//log.Printf("Passed UpdateProduct")
}


func TestStoreDeleteProduct(t *testing.T){
	//log.Printf("Testing DeleteProduct")
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to Mock DataBase")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input int
		rowsAffected int64
		err []error
		outputErr error
	}{
		{input: 1, rowsAffected: 1,err: []error{nil}},
		{input: 2, err: []error{errors.ProductDoesNotExist},outputErr: errors.ProductDoesNotExist},
		{input: 3, err: []error{nil , errors.ProductDoesNotExist},outputErr: errors.ProductDoesNotExist},
	}

	for i,tc:=range testCases{
		mock.ExpectExec("Delete from Product*").WithArgs(tc.input).WillReturnError(tc.err[0]).WillReturnResult(sqlmock.NewResult(0,tc.rowsAffected))
		err:=store.DeleteProduct(tc.input)
		if tc.outputErr != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		
	}
}