package brand

import (
	"catalog/model"
	"catalog/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestStoreGetById(t *testing.T){
	//log.Printf("Testing GetById")
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input int
		output model.Brand
		err error
	}{
		{1,model.Brand{1,"LG"},nil},
		{2,model.Brand{},errors.BrandDoesNotExist},
	}
	str:=[]string{"id","name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(str).AddRow(tc.output.Id,tc.output.Name)
		mock.ExpectQuery("Select Id,Name from Brand where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)

		result,err:=store.GetById(tc.input)
		if tc.err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.err,err)
			}
		} else if tc.output !=result {
			t.Errorf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		
	}
	//log.Printf("Passed GetById")
}


func TestStoreGetByName(t *testing.T){
	//log.Printf("Testing GetByName")
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input string
		output model.Brand
		err error
	}{
		{"LG",model.Brand{1,"LG"},nil},
		{"",model.Brand{},errors.BrandDoesNotExist},
	}
	str:=[]string{"id","name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(str).AddRow(tc.output.Id,tc.output.Name)
		mock.ExpectQuery("Select Id,Name from Brand where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)
		result,err:=store.GetByName(tc.input)
		if tc.err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.err,err)
			}
		} else if tc.output !=result {
			t.Errorf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		
	}
	//log.Printf("Testing GetByName")
}


func TestStoreCreateBrand(t *testing.T){
	//log.Printf("Testing CreateBrand")
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Brand
		output int
		err error
	}{
		{input: model.Brand{Name: "LG"}, output: 1},
		{input: model.Brand{Name: "Hyundai"}, output: 2},
		{err:errors.ThereIsSomeTechnicalIssue},
	}
	//str:=[]string{"id","name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Insert into Brand*").WithArgs(tc.input.Name).WillReturnResult(sqlmock.NewResult(int64(tc.output),1)).WillReturnError(tc.err)
		result,err:=store.CreateBrand(tc.input)
		if tc.err != nil {
			if !reflect.DeepEqual(tc.err , err) {
				t.Error(err)
			}
		} else if tc.output !=result {
			t.Errorf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		
	}
}


func TestStoreUpdateBrand(t *testing.T){
	//log.Printf("Testing UpdateBrand")
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Brand
		rowsAffected int64
		outputErr error
	}{
		{input: model.Brand{Id: 1, Name: "LG"}, rowsAffected: 1},
		{ outputErr: errors.PleaseEnterSomeData},
	}
	//str:=[]string{"id","name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Update Brand*").WithArgs(tc.input.Name,tc.input.Id).WillReturnError(tc.outputErr).WillReturnResult(sqlmock.NewResult(0,tc.rowsAffected))
		err:=store.UpdateBrand(tc.input)
		if tc.outputErr != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		
	}
	//log.Printf("Passed UpdateBrand")
}


func TestStoreDeleteBrand(t *testing.T){
	//log.Printf("Testing DeleteBrand")
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Errorf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input int
		rowsAffected int64
		outputErr error
	}{
		{1,1,nil},
		{2,0,errors.BrandDoesNotExist},
	}
	//str:=[]string{"id","name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Delete from Brand where*").WithArgs(tc.input).WillReturnError(tc.outputErr).WillReturnResult(sqlmock.NewResult(0,tc.rowsAffected))
		err:=store.DeleteBrand(tc.input)
		if tc.outputErr != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Errorf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		
	}
	//log.Printf("Passed DeleteBrand")
}

