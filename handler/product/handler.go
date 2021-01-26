package product

import (
	"catalog/model"
	"catalog/service/product"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type handler struct {
	service product.Service
}

func New(pr product.Service)handler{
	return handler{pr}
}

func (h handler)GetById(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["id"]
	for _,val:=range id {
		if '0'>val || val>'9' {
			http.Error(w,("Please enter a valid numeric Id greater than Zero"),http.StatusBadRequest)
			return
		}
	}
	numId,_:=strconv.Atoi(id)
	if(numId<=0) {
		http.Error(w,("Id can not be zero or negative"),http.StatusBadRequest)
		return
	}
	result,err:=h.service.GetById(numId)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	//log.Println(result)
	post,err:=json.Marshal(result)
	if err != nil {
		http.Error(w,"Input is not in correct format",http.StatusBadRequest)
		return
	}
	w.Write(post)
}

func (h handler)GetByName(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["name"]
	//id=strings.ToUpper(id)
	result,err:=h.service.GetByName(id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	////log.Println(result)
	post,_:=json.Marshal(result)
	w.Write(post)
}

func (h handler)CreateProduct(w http.ResponseWriter,r* http.Request){
	pr:=model.Product{}
	json.NewDecoder(r.Body).Decode(&pr)
	result,err:=h.service.CreateProduct(pr)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	post,_:=json.Marshal(result)
	w.Write(post)
}
func (h handler)UpdateProduct(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["id"]
	pr:=model.Product{}
	json.NewDecoder(r.Body).Decode(&pr)
	numId,_:=strconv.Atoi(id)
	pr.Id= numId
	result,err:=h.service.UpdateProduct(pr)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	post,_:=json.Marshal(result)
	w.Write(post)
}
func (h handler)DeleteProduct(w http.ResponseWriter,r* http.Request){
	id:=mux.Vars(r)["id"]
	numId,_:=strconv.Atoi(id)
	result,err:=h.service.DeleteProduct(numId)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	var temp []interface{}
	temp=append(temp,result)
	temp=append(temp," Product is Deleted Successfully")
	post,_:=json.Marshal(temp)
	w.Write(post)
}
