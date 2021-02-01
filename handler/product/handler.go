package product

import (
	"catalog/errors"
	"catalog/model"
	"catalog/service/product"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type handler struct {
	service product.Service
}

func New(pr product.Service)handler{
	return handler{pr}
}

func validateId(id string)(int ,error) {
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
	numId,err:=validateId(id)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	result,err:=h.service.GetById(numId)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	post,_:=json.Marshal(result)
	_, _ = w.Write(post)
}

func (h handler)GetByName(w http.ResponseWriter,r* http.Request){
	name:=mux.Vars(r)["name"]
	result,err:=h.service.GetByName(name)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	post,_:=json.Marshal(result)
	_, _ = w.Write(post)
}

func (h handler)CreateProduct(w http.ResponseWriter,r* http.Request){
	pr:=model.Product{}
	err:=json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		http.Error(w,errors.InputIsNotInCorrectFormat.Error(),http.StatusBadRequest)
		return
	}
	result,err:=h.service.CreateProduct(pr)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	post,_:=json.Marshal(result)
	_, _ = w.Write(post)
}
func (h handler)UpdateProduct(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["id"]
	pr:=model.Product{}
	err:=json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		http.Error(w,errors.PleaseEnterValidData.Error(),http.StatusBadRequest)
		return
	}
	numId, err:= validateId(id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	pr.Id=numId
	result,err:=h.service.UpdateProduct(pr)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	post,_:=json.Marshal(result)
	_, _ = w.Write(post)
}
func (h handler)DeleteProduct(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["id"]
	numId,err:=validateId(id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	err=h.service.DeleteProduct(numId)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
