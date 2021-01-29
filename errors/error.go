package errors

type constError string
func (err constError) Error()string{
	return string(err)
}
var (
	ThereIsSomeTechnicalIssue 	= constError("There is some Technical Issue")
	PleaseEnterSomeData			= constError("Please Enter Some Data")
	ProductDoesNotExist			= constError("Product Doesn't Exist")
	BrandDoesNotExist			= constError("Brand Doesn't Exist")
	IdCantBeZeroOrNegative		= constError("Id can't be zero or negative")
	PleaseEnterValidId			= constError("Please enter a valid numeric Id greater than Zero")
	InputIsNotInCorrectFormat	= constError("Input is not in Correct Format")
)
