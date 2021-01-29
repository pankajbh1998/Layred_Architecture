package product

import (
	model "catalog/model"
	"catalog/store/brand"
	"catalog/store/product"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestGetById(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=product.NewMockStore(ctrl)
	bs:=brand.NewMockStore(ctrl)
	servicepr:=New(ps,bs)
	testCases:=[]struct{
		input int
		intermediateOutput model.Product
		output model.Product
		err error
	}{
		{
			1,
			model.Product{1,"Ref",model.Brand{1,""}},
			model.Product{1,"Ref",model.Brand{1,"LG"}},
			nil,
		},
		{
			2,
			model.Product{},
			model.Product{},
			errors.New("Product does not exist"),
		},
		{
			3,
			model.Product{3,"Washing", model.Brand{2,""}},
			model.Product{3,"Washing", model.Brand{2,"Hyundai"}},
			nil,
		},
	}
	for i,tc:=range testCases{
		ps.EXPECT().GetById(tc.input).Return(tc.intermediateOutput,tc.err)
		if tc.err == nil {
			bs.EXPECT().GetById(tc.output.Brand.Id).Return(tc.output.Brand, tc.err)
		}
		result,err:=servicepr.GetById(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Errorf("Failed at %v\n Error %v\n",i+1,tc.output)
			}
		} else if result !=tc.output {
			t.Errorf("Failed at %v\n Expected Output %v\n Actual Output %v\n",i+1,tc.output,result)
		}
	}
}

func TestGetByName(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=product.NewMockStore(ctrl)
	bs:=brand.NewMockStore(ctrl)
	servicepr:=New(ps,bs)
	testCases:=[]struct{
		input string
		intermediateOutput []model.Product
		output []model.Product
		err error
	}{
		{
			"Ref",
			[]model.Product{{1,"Ref",model.Brand{1,""}}},
			[]model.Product{{1,"Ref",model.Brand{1,"LG"}}},
			nil,
		},
		{
			"Wash",
			[]model.Product(nil),
			[]model.Product(nil),
			errors.New("Product does not exist"),
		},
		{
			"Washing",
			[]model.Product{{3,"Washing", model.Brand{2,""}}},
			[]model.Product{{3,"Washing", model.Brand{2,"Hyundai"}}},
			nil,
		},
	}
	for i,tc:=range testCases{
		ps.EXPECT().GetByName(tc.input).Return(tc.intermediateOutput,tc.err)
		if tc.err == nil {
			bs.EXPECT().GetById(tc.output[0].Brand.Id).Return(tc.output[0].Brand, tc.err)
		}
		result,err:=servicepr.GetByName(tc.input)
		if err != nil {
			if !reflect.DeepEqual(err,tc.err) {
				t.Errorf("Failed at %v\n Error %v\n",i+1,tc.output)
			}
		} else if !reflect.DeepEqual(result ,tc.output) {
			t.Errorf("Failed at %v\n Expected Output %v\n Actual Output %v\n",i+1,tc.output,result)
		}
	}
}

func TestCreateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	br:=brand.NewMockStore(ctrl)
	pr:=product.NewMockStore(ctrl)
	servicepr:=New(pr,br)
	testCases:=[]struct{
		input model.Product
		intermediateOutput model.Product
		output model.Product
		err error
	}{
		{
			model.Product{Name: "Ref", Brand: model.Brand{Id:2,Name: "LG"}},
			model.Product{ Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			nil,
		},
		{
			model.Product{Name: "Ref", Brand: model.Brand{Name: "LG"}},
			model.Product{Name: "Ref", Brand: model.Brand{Id:2,Name: "LG"}},
			model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			errors.New("Brand does not exist"),
		},
	}
	for i,tc:= range testCases {
		br.EXPECT().GetByName(tc.input.Brand.Name).Return(tc.output.Brand,tc.err)
		if tc.err != nil {
			br.EXPECT().CreateBrand(tc.input.Brand).Return(tc.intermediateOutput.Brand.Id,nil)
		}
		pr.EXPECT().CreateProduct(tc.intermediateOutput).Return(tc.output.Id,nil)
		pr.EXPECT().GetById(tc.output.Id).Return(tc.output,nil)
		br.EXPECT().GetById(tc.output.Brand.Id).Return(tc.output.Brand,nil)
		result,err:=servicepr.CreateProduct(tc.input)
		if err != nil {
			if err != tc.err{
				t.Error(err)
			}
		} else if result !=tc.output {
			t.Errorf("Failed at %v\n Expected Output %v\n Actual Output %v\n",i+1,tc.output,result)
		}
	}

}

func TestUpdateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=product.NewMockStore(ctrl)
	bs:=brand.NewMockStore(ctrl)
	servicepr:=New(ps,bs)
	testCases:=[]struct{
		input model.Product
		outputS1 model.Product
		outputS2 model.Product
		output model.Product
		err error
	}{
		{
			input:  model.Product{Id: 1,Brand: model.Brand{Name:"Oppo"}},
			outputS1:  model.Product{Id: 1,Brand: model.Brand{Id: 2,Name:"Oppo"}},
			outputS2:  model.Product{Id: 1,Name:"R1",Brand: model.Brand{Id: 2}},
			output: model.Product{Id: 1,Name:"R1",Brand: model.Brand{Id: 2,Name:"Oppo"}},
		},
		{
			input:  model.Product{Id: 1,Name:"R2",Brand: model.Brand{}},
			outputS1:  model.Product{Id: 1,Name:"R2",Brand: model.Brand{}},
			outputS2:  model.Product{Id: 1,Name:"R2",Brand: model.Brand{Id: 2}},
			output: model.Product{Id: 1,Name:"R2",Brand: model.Brand{Id: 2,Name:"Oppo"}},
		},
	}
	for i,tc:=range testCases {
		if tc.input.Brand.Name != "" {
			bs.EXPECT().GetByName(tc.input.Brand.Name).Return(tc.output.Brand,nil)
		}
		ps.EXPECT().UpdateProduct(tc.outputS1).Return(nil)
		ps.EXPECT().GetById(tc.input.Id).Return(tc.outputS2,nil)
		bs.EXPECT().GetById(tc.output.Brand.Id).Return(tc.output.Brand,nil)
		result,err:=servicepr.UpdateProduct(tc.input)
		if err != nil {
			if err != tc.err{
				t.Error(err)
			}
		} else if result !=tc.output {
			t.Errorf("Failed at %v\n Expected Output :%v\n Actual Output : %v\n",i+1,tc.output,result)
		}
	}
}

func TestDeleteProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=product.NewMockStore(ctrl)
	bs:=brand.NewMockStore(ctrl)
	servicepr:=New(ps,bs)
	testCases:=[]struct {
		input  int
		err error
	}{
		{
			input: 1,err:errors.New("Id does not exist"),
		},
		{
			input: 2,
		},
	}
	for _,tc:= range testCases{
		ps.EXPECT().DeleteProduct(tc.input).Return(tc.err)
		err:=servicepr.DeleteProduct(tc.input)
		if tc.err != nil {
			if err != tc.err{
				t.Error(err)
			}
		}
	}
}

func TestAll(t *testing.T){
	TestGetById(t)
	TestGetByName(t)
	TestCreateProduct(t)
	TestUpdateProduct(t)
	TestDeleteProduct(t)
}