package brand

import (
	"catalog/model"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
)

func TestStoreGetById(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input int
		output model.Brand
		err error
	}{
		{1,model.Brand{1,"LG"},nil},
		{2,model.Brand{},errors.New("Brand does not exist")},
	}
	str:=[]string{"id","name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(str).AddRow(tc.output.Id,tc.output.Name)
		mock.ExpectQuery("Select Id,Name from Brand where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)

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
		t.Fatalf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input string
		output model.Brand
		err error
	}{
		{"LG",model.Brand{1,"LG"},nil},
		{"",model.Brand{},errors.New("Brand does not exist")},
	}
	str:=[]string{"id","name"}
	for i,tc:=range testCases{
		row:=sqlmock.NewRows(str).AddRow(tc.output.Id,tc.output.Name)
		mock.ExpectQuery("Select Id,Name from Brand where*").WithArgs(tc.input).WillReturnError(tc.err).WillReturnRows(row)
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


func TestStoreCreateBrand(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Brand
		output int
	}{
		{model.Brand{0,"LG"},1},
		{model.Brand{0,"Hyundai"},2},
	}
	//str:=[]string{"id","name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Insert into Brand*").WithArgs(tc.input.Name).WillReturnResult(sqlmock.NewResult(int64(tc.output),1))
		result:=store.CreateBrand(tc.input)
		if tc.output !=result {
			t.Fatalf("Failed at : %v\nExpected Output : %v\nActual Output : %v",i+1,tc.output,result)
		}
		log.Printf("Passed at %v",i+1)
	}
}


func TestStoreUpdateBrand(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input model.Brand
		rowsAffected int64
		outputErr error
	}{
		{model.Brand{1,"LG"},1,nil},
		{model.Brand{2,"Hyundai"},0,errors.New("Id does not exist")},
	}
	//str:=[]string{"id","name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Update Brand*").WithArgs(tc.input.Name,tc.input.Id).WillReturnError(tc.outputErr).WillReturnResult(sqlmock.NewResult(0,tc.rowsAffected))
		err:=store.UpdateBrand(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Fatalf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		log.Printf("Passed at %v",i+1)
	}
}


func TestStoreDeleteBrand(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Cannot connect to fake db")
	}
	store:=New(fdb)
	testCases:=[]struct{
		input int
		rowsaffected int64
		outputErr error
	}{
		{1,1,nil},
		{2,0,errors.New("Id does not exist")},
	}
	//str:=[]string{"id","name"}
	for i,tc:=range testCases{
		mock.ExpectExec("Delete from Brand where*").WithArgs(tc.input).WillReturnError(tc.outputErr).WillReturnResult(sqlmock.NewResult(0,tc.rowsaffected))
		err:=store.DeleteBrand(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.outputErr) {
				t.Fatalf("Failed at : %v\nExpected Error : %v\nActual Error : %v",i+1,tc.outputErr,err)
			}
		}
		log.Printf("Passed at %v",i+1)
	}
}

