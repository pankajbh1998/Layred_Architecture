package product

import (
	"bytes"
	"catalog/errors"
	"catalog/model"
	"catalog/service/product"
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
	ps := product.NewMockService(ctrl)
	ph := New(ps)
	testCases := []struct {
		input      string
		output     interface{}
		statusCode int
		err        error
	}{
		{input: "1", output: model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1, Name: "LG"}},statusCode: 200},
		{input: "2", output: []byte(errors.ProductDoesNotExist), statusCode: 400, err: errors.ProductDoesNotExist},
		{input: "0", output: []byte(errors.IdCantBeZeroOrNegative), statusCode: 400},
		{input: "abc1", output: []byte(errors.PleaseEnterValidId), statusCode: 400},
	}
	for i, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/product", nil)
		r = mux.SetURLVars(r, map[string]string{"id": tc.input})

		outputByte:=tc.output
		id,err:=strconv.Atoi(tc.input)
		if err==nil && id>0{
			pr:=model.Product{}
			_,ok:=tc.output.([]byte)
			if !ok {
				outputByte,_=json.Marshal(tc.output)
				pr=tc.output.(model.Product)
			}
			ps.EXPECT().GetById(id).Return(pr, tc.err)
		}

		ph.GetById(w, r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(bytes.TrimSpace(res) , outputByte) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(outputByte.([]byte)), string(res))
		}
	}
}

func TestGetByName(t *testing.T){
	ctrl:=gomock.NewController(t)
	servicePr:=product.NewMockService(ctrl)
	handlerePr := New(servicePr)
	testCases:=[]struct {
		input      string
		output     interface{}
		err        error
		statusCode int
	}{
		{
			input: "Mountain Dew",
			output: []model.Product{
						{Id:1,Name:"Mountain Dew",Brand: model.Brand{Name:"Pepsiso"}},
					},
			statusCode: 200,
		},
		{
			input:      "Coca Cola",
			output: 	[]byte(errors.ProductDoesNotExist),
			err:        errors.ProductDoesNotExist,
			statusCode: 400,
		},
	}
	for i,tc:=range testCases{
		outputByte:=tc.output
		sendData:=[]model.Product(nil)
		if tc.err == nil{
			outputByte,_=json.Marshal(tc.output)
			sendData=tc.output.([]model.Product)
		}
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("GET","/product",nil)
		r=mux.SetURLVars(r,map[string]string{"name": tc.input})
		servicePr.EXPECT().GetByName(tc.input).Return(sendData,tc.err)
		handlerePr.GetByName(w,r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(bytes.TrimSpace(res) , outputByte) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(outputByte.([]byte)), string(res))
		}
	}
}
func TestCreateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=product.NewMockService(ctrl)
	handlerPr:=New(ps)
	testCases:=[]struct{
		input model.Product
		output interface{}
		statusCode int
		err error
	}{
		{
			input: model.Product{Name:"R1",Brand: model.Brand{Name:"Realme"}},
			output: model.Product{Id: 1,Name:"R1",Brand: model.Brand{Name:"Realme"}},
			statusCode: 200,
		},
		{
			input: model.Product{},
			output: []byte(errors.PleaseEnterSomeData),
			statusCode: 400,
			err:errors.PleaseEnterSomeData,
		},
		{
			output: []byte(errors.InputIsNotInCorrectFormat),
			statusCode: 400,
			err:errors.InputIsNotInCorrectFormat,
		},
	}
	//json.Compact()
	for i,tc := range testCases{
		pr:=model.Product{}
		outputByte:=tc.output
		_,ok:=tc.output.([]byte)
		if !ok {
			outputByte, _ = json.Marshal(tc.output)
			pr = tc.output.(model.Product)
		}

		var data interface{}
		data=tc.input

		if tc.err == errors.InputIsNotInCorrectFormat{
			data=[]string{"name}"}
		} else {
			ps.EXPECT().CreateProduct(tc.input).Return(pr, tc.err)
		}

		inputByte,_:=json.Marshal(data)
		sendBody:=bytes.NewReader(inputByte)
		//r:=httptest.NewRequest("POST","/product",strings.NewReader(string(inputByte)))
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("POST","/product",sendBody)

		handlerPr.CreateProduct(w,r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(bytes.TrimSpace(res) , outputByte) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(outputByte.([]byte)), string(res))
		}
	}
}

func TestUpdateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=product.NewMockService(ctrl)
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
			output: []byte(errors.ProductDoesNotExist),
			statusCode: 400,
			err:errors.ProductDoesNotExist,
		},
		{
			id:"3",
			input: model.Product{},
			output: []byte(errors.PleaseEnterSomeData),
			statusCode: 400,
			err:errors.PleaseEnterSomeData,
		},
		{
			id:"abc",
			input: model.Product{},
			output: []byte(errors.PleaseEnterValidId),
			statusCode: 400,
			err:errors.PleaseEnterValidId,
		},
	}
	for i,tc := range testCases{
		inputByte,_:=json.Marshal(tc.input)
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("PUT","/product",bytes.NewReader(inputByte))
		r=mux.SetURLVars(r,map[string]string{"id":tc.id})
		numId,_:=strconv.Atoi(tc.id)
		tc.input.Id = numId


		sendData:=model.Product{}
		_,ok:=tc.output.([]byte)
		outputByte:=tc.output
		if !ok {
			var err error
			sendData=tc.output.(model.Product)
			outputByte,err=json.Marshal(tc.output)
			if err != nil {
				t.Fatalf(err.Error())
			}
		}
		if tc.err != errors.PleaseEnterValidId {
			ps.EXPECT().UpdateProduct(tc.input).Return(sendData, tc.err)
		}
		handlerPr.UpdateProduct(w, r)
		result := w.Result()
		res, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatalf(err.Error())
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if !reflect.DeepEqual(bytes.TrimSpace(res) , outputByte) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(outputByte.([]byte)), string(res))
		}
	}
}

func TestDeleteProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	servicePr:=product.NewMockService(ctrl)
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
			output: []byte(errors.ProductDoesNotExist),
			err:	[]error{nil,errors.ProductDoesNotExist},
			statusCode: 400,
		},
		{
			input:	"0",
			output: []byte(errors.IdCantBeZeroOrNegative),
			err:	[]error{errors.IdCantBeZeroOrNegative,nil},
			statusCode: 400,
		},
		{
			input:	"abc",
			output: []byte(errors.PleaseEnterValidId),
			err:	[]error{errors.PleaseEnterValidId,nil},
			statusCode: 400,
		},

	}
	for i,tc := range testCases{
		w:=httptest.NewRecorder()
		r:=httptest.NewRequest("DELETE","/product",nil)
		r=mux.SetURLVars(r,map[string]string{"id":tc.input})
		numId,err:=strconv.Atoi(tc.input)
		if tc.err[0] == nil {
			servicePr.EXPECT().DeleteProduct(numId).Return(tc.err[1])
		}
		handlerPr.DeleteProduct(w,r)
		result:=w.Result()
		res,err:=ioutil.ReadAll(result.Body)
		//log.Println(i,string(res))
		if err != nil {
			t.Error(err)
		} else if tc.statusCode != result.StatusCode {
			t.Errorf("Failed at %v Wrong Status Code\n", i+1)
		} else if tc.output != nil && !reflect.DeepEqual(bytes.TrimSpace(res) , tc.output) {
			t.Errorf("Failed at %v\nExpected Output : %v\nActual Output   : %v\n", i+1, string(tc.output.([]byte)), string(res))
		}
	}
}