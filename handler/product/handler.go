package product

import (
	"catalog/errors"
	"catalog/model"
	"catalog/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type handler struct {
	service service.Product
}

func New(pr service.Product)handler{
	return handler{pr}
}

func (h handler)myPrint(w http.ResponseWriter, result interface{}) {
	w.Header().Set("content-type","application.json")
	res,ok:=result.(model.JsonPrint)
	if ok {
		w.WriteHeader(res.Code)
	}
	post,_:=json.Marshal(result)
	fmt.Fprintf(w, "%v",string(post))
}
func (h handler)validateId(id string)(int ,error) {
	for _,val:=range id {
		if '0'>val || val>'9' {
			return 0,errors.PleaseEnterValidId
		}
	}
	numId,_:=strconv.Atoi(id)
	if numId<=0 {
		return 0,errors.IdCantBeZeroOrNegative
	}
	return numId,nil
}
func (h handler)GetById(w http.ResponseWriter,r *http.Request){
	id:=mux.Vars(r)["id"]
	numId,err:=h.validateId(id)
	if err != nil{
		h.myPrint(w,model.JsonPrint{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	result,err:=h.service.GetById(numId)
	if err != nil {
		h.myPrint(w,model.JsonPrint{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	h.myPrint(w,result)
}

func (h handler)GetByName(w http.ResponseWriter,r* http.Request){
	name:=mux.Vars(r)["name"]
	result,err:=h.service.GetByName(name)
	if err != nil {
		h.myPrint(w,model.JsonPrint{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	h.myPrint(w,result)
}

func (h handler)CreateProduct(w http.ResponseWriter,r* http.Request){
	pr:=model.Product{}
	err:=json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		h.myPrint(w,model.JsonPrint{Code: http.StatusBadRequest, Message: errors.PleaseEnterValidData.Error()})
		return
	}
	result,err:=h.service.CreateProduct(pr)
	if err != nil {
		h.myPrint(w,model.JsonPrint{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	h.myPrint(w,result)
}
func (h handler)UpdateProduct(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["id"]
	pr:=model.Product{}
	err:=json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		h.myPrint(w,model.JsonPrint{Code: http.StatusBadRequest, Message: errors.PleaseEnterValidData.Error()})
		return
	}
	numId, err:= h.validateId(id)
	if err != nil {
		h.myPrint(w, model.JsonPrint{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	pr.Id=numId
	result,err:=h.service.UpdateProduct(pr)
	if err != nil {
		sentCode:=http.StatusBadRequest
		if err == errors.ThereIsSomeTechnicalIssue {
			sentCode=http.StatusInternalServerError
		}
		h.myPrint(w, model.JsonPrint{Code: sentCode, Message: err.Error()})
		return
	}
	h.myPrint(w,result)
}
func (h handler)DeleteProduct(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["id"]
	numId,err:=h.validateId(id)
	if err != nil {
		h.myPrint(w,model.JsonPrint{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	err=h.service.DeleteProduct(numId)
	if err != nil {
		sentCode:=http.StatusBadRequest
		if err == errors.ThereIsSomeTechnicalIssue {
			sentCode=http.StatusInternalServerError
		}
		h.myPrint(w, model.JsonPrint{Code: sentCode, Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
