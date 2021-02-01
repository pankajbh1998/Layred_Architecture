package product

import (
	"bytes"
	"catalog/errors"
	"catalog/model"
	"catalog/service"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGetBYId(t *testing.T) {
	ctrl := gomock.NewController(t)
	ps := service.NewMockProduct(ctrl)
	ph := New(ps)
	testCases := []struct {
		input      		string
		output     		interface{}
		statusCode 		int
		expectedErr 	error
	}{
		{input: "1", output: model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1, Name: "LG"}},statusCode: 200},
		{input: "2", output: model.JsonPrint{Code: 400, Message: errors.ProductDoesNotExist.Error()}, statusCode: 400, expectedErr: errors.ProductDoesNotExist},
		{input: "0", output: model.JsonPrint{Code: 400, Message: errors.IdCantBeZeroOrNegative.Error()}, statusCode: 400},
		{input: "abc1", output: model.JsonPrint{Code: 400, Message: errors.PleaseEnterValidId.Error()}, statusCode: 400},
	}
	for i, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/product", nil)
		r = mux.SetURLVars(r, map[string]string{"id": tc.input})
		id,err:=strconv.Atoi(tc.input)
		if err==nil && id>0{
			pr,_:=tc.output.(model.Product)
			ps.EXPECT().GetById(id).Return(pr, tc.expectedErr)
		}

		ph.GetById(w, r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		expectedOutput,_:=json.Marshal(tc.output)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(res , expectedOutput) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(expectedOutput), string(res))
		}
	}
}

func TestGetByName(t *testing.T){
	ctrl:=gomock.NewController(t)
	servicePr:=service.NewMockProduct(ctrl)
	handlerPr := New(servicePr)
	testCases:=[]struct {
		input      		string
		output     		interface{}
		expectedErr     error
		statusCode 		int
	}{
		{
			input: "Mountain Dew",
			output: []model.Product{
						{Id:1,Name:"Mountain Dew",Brand: model.Brand{Name:"Pepsico"}},
					},
			statusCode: 200,
		},
		{
			input:      "Coca Cola",
			output:      model.JsonPrint{Code: 400, Message: errors.ProductDoesNotExist.Error()},
			expectedErr: errors.ProductDoesNotExist,
			statusCode: 400,
		},
	}
	for i,tc:=range testCases{
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("GET","/product",nil)
		r=mux.SetURLVars(r,map[string]string{"name": tc.input})
		pr,_:=tc.output.([]model.Product)
		servicePr.EXPECT().GetByName(tc.input).Return(pr, tc.expectedErr)
		handlerPr.GetByName(w,r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		expectedOutput,_:=json.Marshal(tc.output)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(res , expectedOutput) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(expectedOutput), string(res))
		}
	}
}
func TestCreateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=service.NewMockProduct(ctrl)
	handlerPr:=New(ps)
	testCases:=[]struct{
		input model.Product
		output interface{}
		statusCode int
		err error
	}{
		{
			input: 		model.Product{Name:"R1",Brand: model.Brand{Name:"Realme"}},
			output: 	model.Product{Id: 1,Name:"R1",Brand: model.Brand{Name:"Realme"}},
			statusCode: 200,
		},
		{
			output: 	model.JsonPrint{Code: 400, Message: errors.PleaseEnterValidData.Error()},
			statusCode: 400,
			err: 		errors.PleaseEnterValidData,
		},
		{
			output: 	model.JsonPrint{Code: 500, Message: errors.ThereIsSomeTechnicalIssue.Error()},
			statusCode: 500,
			err:		errors.ThereIsSomeTechnicalIssue,
		},
	}
	//json.Compact()
	for i,tc := range testCases{
		var data interface{}
		data=tc.input
		if tc.err == errors.PleaseEnterValidData{
			data=[]string{"name}"}
		} else {
			pr,_:=tc.output.(model.Product)
			ps.EXPECT().CreateProduct(tc.input).Return(pr, tc.err)
		}

		inputByte,_:=json.Marshal(data)
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("POST","/product",bytes.NewReader(inputByte))
		//r:=httptest.NewRequest("POST","/product",strings.NewReader(string(inputByte)))

		handlerPr.CreateProduct(w,r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		expectedOutput,_:=json.Marshal(tc.output)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(res , expectedOutput) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(expectedOutput), string(res))
		}
	}
}

func TestUpdateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=service.NewMockProduct(ctrl)
	handlerPr:=New(ps)
	testCases:=[]struct{
		id string
		input model.Product
		output interface{}
		statusCode int
		err error
	}{
		{
			id:"1",
			input: model.Product{Name:"",Brand: model.Brand{Name:"Realme"}},
			output: model.Product{Id: 1,Name:"R1",Brand: model.Brand{Name:"Realme"}},
			statusCode: 200,
		},
		{
			id:"2",
			input: model.Product{Name:"R2",Brand: model.Brand{Name:""}},
			output: model.JsonPrint{Code: 400, Message: errors.ProductDoesNotExist.Error()},
			statusCode: 400,
			err:errors.ProductDoesNotExist,
		},
		{
			id:"3",
			input: model.Product{},
			output: model.JsonPrint{Code: 400, Message: errors.PleaseEnterValidData.Error()},
			statusCode: 400,
			err:errors.PleaseEnterValidData,
		},
		{
			id:"4",
			input: model.Product{},
			output: model.JsonPrint{Code: 500, Message: errors.ThereIsSomeTechnicalIssue.Error()},
			statusCode: 500,
			err:errors.ThereIsSomeTechnicalIssue,
		},
		{
			id:"abc",
			input: model.Product{},
			output: model.JsonPrint{Code: 400, Message: errors.PleaseEnterValidId.Error()},
			statusCode: 400,
			err:errors.PleaseEnterValidId,
		},
	}
	for i,tc := range testCases{
		inputByte,_:=json.Marshal(tc.input)
		if tc.err==errors.PleaseEnterValidData {
			inputByte,_=json.Marshal("Name{")
		}
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("PUT","/product",bytes.NewReader(inputByte))
		r=mux.SetURLVars(r,map[string]string{"id":tc.id})
		tc.input.Id,_=strconv.Atoi(tc.id)
		if tc.input.Id>0 && tc.err != errors.PleaseEnterValidData{
			pr,_:=tc.output.(model.Product)
			ps.EXPECT().UpdateProduct(tc.input).Return(pr, tc.err)
		}
		handlerPr.UpdateProduct(w, r)
		result := w.Result()
		res, err := ioutil.ReadAll(result.Body)
		expectedOutput,_:=json.Marshal(tc.output)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(res , expectedOutput) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(expectedOutput), string(res))
		}
	}
}

func TestDeleteProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	servicePr:=service.NewMockProduct(ctrl)
	handlerPr:=New(servicePr)
	testCases:=[]struct{
		input string
		output interface{}
		err []error
		statusCode int
	}{
		{
			input:	"1",
			output: nil,
			err:	[]error{nil,nil},
			statusCode: 204,
		},
		{
			input:	"2",
			output: model.JsonPrint{Code: 400, Message: errors.ProductDoesNotExist.Error()},
			err:	[]error{nil,errors.ProductDoesNotExist},
			statusCode: 400,
		},
		{
			input:	"2",
			output: model.JsonPrint{Code: 500, Message: errors.ThereIsSomeTechnicalIssue.Error()},
			err:	[]error{nil,errors.ThereIsSomeTechnicalIssue},
			statusCode: 500,
		},
		{
			input:	"0",
			output: model.JsonPrint{Code: 400, Message: errors.IdCantBeZeroOrNegative.Error()},
			err:	[]error{errors.IdCantBeZeroOrNegative,nil},
			statusCode: 400,
		},
		{
			input:	"abc",
			output: model.JsonPrint{Code: 400, Message: errors.PleaseEnterValidId.Error()},
			err:	[]error{errors.PleaseEnterValidId,nil},
			statusCode: 400,
		},

	}
	for i,tc := range testCases{
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("DELETE","/product",nil)
		r=mux.SetURLVars(r,map[string]string{"id":tc.input})

		if tc.err[0] == nil {
			numId,_:=strconv.Atoi(tc.input)
			servicePr.EXPECT().DeleteProduct(numId).Return(tc.err[1])
		}
		handlerPr.DeleteProduct(w,r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		expectedOutput, _ := json.Marshal(tc.output)
		if err != nil {
			t.Error(err)
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if tc.output != nil && !reflect.DeepEqual(res , expectedOutput) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(expectedOutput), string(res))
		}
	}
}