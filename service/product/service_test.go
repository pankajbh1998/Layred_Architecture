package product

import (
	"catalog/errors"
	"catalog/model"
	"catalog/store"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestGetById(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=store.NewMockProduct(ctrl)
	bs:=store.NewMockBrand(ctrl)
	servicePr:=New(ps,bs)
	testCases:=[]struct{
		input int
		intermediateOutput model.Product
		output model.Product
		err error
	}{
		{
			input:              1,
			intermediateOutput: model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1}},
			output:             model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1, Name: "LG"}},
		},
		{
			input: 2,
			err:   errors.ProductDoesNotExist,
		},
		{
			input:              3,
			intermediateOutput: model.Product{Id: 3, Name: "Washing", Brand: model.Brand{Id: 2}},
			output:             model.Product{Id: 3, Name: "Washing", Brand: model.Brand{Id: 2, Name: "Hyundai"}},
		},
	}
	for i,tc:=range testCases{
		ps.EXPECT().GetById(tc.input).Return(tc.intermediateOutput,tc.err)
		if tc.err == nil {
			bs.EXPECT().GetById(tc.output.Brand.Id).Return(tc.output.Brand, tc.err)
		}
		result,err:=servicePr.GetById(tc.input)
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
	ps:=store.NewMockProduct(ctrl)
	bs:=store.NewMockBrand(ctrl)
	servicePr:=New(ps,bs)
	testCases:=[]struct{
		input string
		intermediateOutput []model.Product
		output []model.Product
		expectedErr error
	}{
		{
			input:              "Ref",
			intermediateOutput: []model.Product{{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1}}},
			output:             []model.Product{{Id: 1, Name: "Ref", Brand: model.Brand{Id: 1, Name: "LG"}}},
		},
		{
			input:              "Wash",
			intermediateOutput: []model.Product(nil),
			output:             []model.Product(nil),
			expectedErr:        errors.ProductDoesNotExist,
		},
		{
			input:              "Washing",
			intermediateOutput: []model.Product{{Id: 3, Name: "Washing", Brand: model.Brand{Id: 2}}},
			output:             []model.Product{{Id: 3, Name: "Washing", Brand: model.Brand{Id: 2, Name: "Hyundai"}}},
		},
	}
	for i,tc:=range testCases{
		ps.EXPECT().GetByName(tc.input).Return(tc.intermediateOutput,tc.expectedErr)
		if tc.expectedErr == nil {
			bs.EXPECT().GetById(tc.output[0].Brand.Id).Return(tc.output[0].Brand, tc.expectedErr)
		}
		result,err:=servicePr.GetByName(tc.input)
		if tc.expectedErr != nil {
			if !reflect.DeepEqual(err , tc.expectedErr){
				t.Error(err)
			}
		} else if !reflect.DeepEqual(result ,tc.output) {
			t.Errorf("Failed at %v\n Expected Output %v\n Actual Output %v\n",i+1,tc.output,result)
		}
	}
}

func TestCreateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	pr:=store.NewMockProduct(ctrl)
	br:=store.NewMockBrand(ctrl)
	servicePr:=New(pr,br)
	testCases:=[]struct{
		input model.Product
		intermediateOutput model.Product
		output model.Product
		err []error
		expectedErr error
	}{
		{
			input:              model.Product{Name: "Ref", Brand: model.Brand{Id: 2,Name: "LG"}},
			intermediateOutput: model.Product{ Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			output:             model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			err:				[]error{nil,nil,nil,nil},
			expectedErr: 		nil,
		},
		{
			input:              model.Product{Name: "Ref", Brand: model.Brand{Name: "LG"}},
			intermediateOutput: model.Product{Name: "Ref", Brand: model.Brand{Id: 2,Name: "LG"}},
			output:             model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			err:                []error{errors.BrandDoesNotExist,nil,nil,nil},
			expectedErr: 		nil,
		},
		{
			input:              model.Product{Name: "Ref", Brand: model.Brand{Name: "LG"}},
			intermediateOutput: model.Product{Name: "Ref", Brand: model.Brand{Id: 2,Name: "LG"}},
			output:             model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			err:                []error{errors.BrandDoesNotExist,errors.ThereIsSomeTechnicalIssue,nil,nil},
			expectedErr: 		errors.ThereIsSomeTechnicalIssue,
		},
		{
			input:              model.Product{Name: "Ref", Brand: model.Brand{Name: "LG"}},
			intermediateOutput: model.Product{Name: "Ref", Brand: model.Brand{Id: 2,Name: "LG"}},
			output:             model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			err:                []error{errors.BrandDoesNotExist,nil,errors.ThereIsSomeTechnicalIssue,nil},
			expectedErr: 		errors.ThereIsSomeTechnicalIssue,
		},
		{
			input:              model.Product{Name: "Ref", Brand: model.Brand{Name: "LG"}},
			intermediateOutput: model.Product{Name: "Ref", Brand: model.Brand{Id: 2,Name: "LG"}},
			output:             model.Product{Id: 1, Name: "Ref", Brand: model.Brand{Id: 2, Name: "LG"}},
			err:                []error{errors.BrandDoesNotExist,nil,nil,errors.ThereIsSomeTechnicalIssue},
			expectedErr: 		errors.ThereIsSomeTechnicalIssue,
		},
	}
	for i,tc:= range testCases {
		br.EXPECT().GetByName(tc.input.Brand.Name).Return(tc.output.Brand,tc.err[0])
		if tc.err[0] != nil {
			br.EXPECT().CreateBrand(tc.input.Brand).Return(tc.intermediateOutput.Brand.Id,tc.err[1])
		}
		if tc.err[1] == nil {
			pr.EXPECT().CreateProduct(tc.intermediateOutput).Return(tc.output.Id, tc.err[2])
			if tc.err[2]== nil {
				pr.EXPECT().GetById(tc.output.Id).Return(tc.output, tc.err[3])
				if tc.err[3]== nil {
					br.EXPECT().GetById(tc.output.Brand.Id).Return(tc.output.Brand, nil)
				}
			}
		}
		result,err:=servicePr.CreateProduct(tc.input)
		if tc.expectedErr != nil {
			if !reflect.DeepEqual(err , tc.expectedErr){
				t.Error(err)
			}
		} else if result !=tc.output {
			t.Errorf("Failed at %v\n Expected Output %v\n Actual Output %v\n",i+1,tc.output,result)
		}
	}

}

func TestUpdateProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=store.NewMockProduct(ctrl)
	bs:=store.NewMockBrand(ctrl)
	servicePr:=New(ps,bs)
	testCases:=[]struct{
		input model.Product
		output []model.Product
		expectedOutput model.Product
		err []error
		expectedErr error
	}{
		{
			input:  		model.Product{Id: 1,Brand: model.Brand{Name:"Oppo"}},
			output:  		[]model.Product{
								{Id: 1,Brand: model.Brand{Id: 2,Name:"Oppo"}},
								{Id: 1,Name:"R1",Brand: model.Brand{Id: 2}},
							},
			expectedOutput: model.Product{Id: 1,Name:"R1",Brand: model.Brand{Id: 2,Name:"Oppo"}},
			err:			[]error{nil,nil,nil,nil},
		},
		{
			input:  		model.Product{Id: 1,Name:"R2",Brand: model.Brand{Name:"Oppo"}},
			output:  		[]model.Product{
							{Id: 1,Name:"R2",Brand: model.Brand{Id:2,Name:"Oppo"}},
							{Id: 1, Name: "R2", Brand: model.Brand{Id: 2}},
							},
			expectedOutput: model.Product{Id: 1,Name:"R2",Brand: model.Brand{Id: 2,Name:"Oppo"}},
			err:			[]error{errors.BrandDoesNotExist,nil,nil,nil},
			expectedErr: 	nil,
		},
		{
			input:  		model.Product{Id: 1,Name:"R2",Brand: model.Brand{Name:"Oppo"}},
			output:  		[]model.Product{
							{Id:1},
							{Id: 1, Name: "R2",Brand: model.Brand{Id:2,Name:"Oppo"}},
							},
			err:			[]error{errors.BrandDoesNotExist,errors.ThereIsSomeTechnicalIssue,nil,nil},
			expectedErr: 	errors.ThereIsSomeTechnicalIssue,
		},
		{
			input:  		model.Product{Id: 1},
			err:			[]error{nil,nil,errors.PleaseEnterValidData},
			output:			[]model.Product{{Id:1},{Id:1}},
			expectedErr: 	errors.PleaseEnterValidData,
		},
		{
			input:  		model.Product{Id: 1,Brand: model.Brand{Name:"Oppo"}},
			output:  		[]model.Product{
							{Id: 1,Brand: model.Brand{Id: 2,Name:"Oppo"}},
							{Id: 1,Name:"R1",Brand: model.Brand{Id: 2}},
							},
			err:			[]error{nil,nil,nil,errors.ProductDoesNotExist},
			expectedErr: 	errors.ProductDoesNotExist,
		},
	}
	for i,tc:=range testCases {
		if tc.input.Brand.Name != "" {
			bs.EXPECT().GetByName(tc.input.Brand.Name).Return(tc.output[0].Brand,tc.err[0])
			if tc.err[0] != nil {
				bs.EXPECT().CreateBrand(tc.input.Brand).Return(tc.output[0].Brand.Id,tc.err[1])
			}
		}
		if tc.err[1] == nil {
			ps.EXPECT().UpdateProduct(tc.output[0]).Return(tc.err[2])
			if tc.err[2] == nil {
				ps.EXPECT().GetById(tc.input.Id).Return(tc.output[1], tc.err[3])
				if tc.err[3] == nil {
					bs.EXPECT().GetById(tc.expectedOutput.Brand.Id).Return(tc.expectedOutput.Brand, nil)
				}
			}
		}

		result,err:=servicePr.UpdateProduct(tc.input)
		if tc.expectedErr != nil {
			if !reflect.DeepEqual(err , tc.expectedErr){
				t.Error(err)
			}
		} else if result !=tc.expectedOutput {
			t.Errorf("Failed at %v\n Expected Output :%v\n Actual Output : %v\n",i+1,tc.expectedOutput,result)
		}
	}
}

func TestDeleteProduct(t *testing.T){
	ctrl:=gomock.NewController(t)
	ps:=store.NewMockProduct(ctrl)
	bs:=store.NewMockBrand(ctrl)
	servicePr:=New(ps,bs)
	testCases:=[]struct {
		input  int
		expectedErr error
	}{
		{
			input: 1,
			expectedErr: errors.ProductDoesNotExist,
		},
		{
			input: 2,
			expectedErr: nil,
		},
	}
	for _,tc:= range testCases{
		ps.EXPECT().DeleteProduct(tc.input).Return(tc.expectedErr)
		err:=servicePr.DeleteProduct(tc.input)
		if tc.expectedErr != nil {
			if !reflect.DeepEqual(err , tc.expectedErr){
				t.Error(err)
			}
		}
	}
}
