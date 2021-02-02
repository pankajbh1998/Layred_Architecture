package errors

type constError string
func (err constError) Error()string{
	return string(err)
}
const (
	ThereIsSomeTechnicalIssue 	= constError("There is some Technical Issue")
	PleaseEnterValidData		= constError("Please Enter Some Valid Data")
	ProductDoesNotExist			= constError("Product Doesn't Exist")
	BrandDoesNotExist			= constError("Brand Doesn't Exist")
	IdCantBeZeroOrNegative		= constError("Id can't be zero or negative")
	PleaseEnterValidId			= constError("Please enter a valid numeric Id greater than Zero")
	InputIsNotInCorrectFormat	= constError("Input is  Incorrect Format")
)
//
//type ProductDoesNotExist struct{
//	msg interface{}
//}
//func (err ProductDoesNotExist)Error()string{
//	return fmt.Sprintf("Product with %v Doesn't exist",err.msg)
//}