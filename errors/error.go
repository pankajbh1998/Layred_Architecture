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
	InputIsNotInCorrectFormat	= constError("Input is not in Correct Format")
)

//type msg struct{
//	msgerr string
//}
//func (err msg)Error()string{
//	return err.msgerr
//}
//func Newwrap(msg1 []error)msg{
//	for msg21{
//
//	}
//	return msg{msg1}
//}