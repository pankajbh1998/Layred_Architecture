package main

import (
	productHandler "catalog/handler/product"
	productService "catalog/service/product"
	brandStore "catalog/store/brand"
	productStore "catalog/store/product"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main(){
	r:=mux.NewRouter()
	db,err:=sql.Open("mysql","Pankaj:Pankaj@123@tcp(127.0.0.1)/Company")
	if err != nil {
		log.Fatal(err)
		log.Fatalf("cannot connect with database")
	}
	storebr:=brandStore.New(db)
	storepr:=productStore.New(db)
	servicepr:=productService.New(storepr,storebr)
	handlerpr:=productHandler.New(servicepr)
	r.HandleFunc("/Product",handlerpr.GetByName).Methods("GET").Queries("name","{name}")
	r.HandleFunc("/Product/{id}",handlerpr.GetById).Methods("GET")
	r.HandleFunc("/Product",handlerpr.CreateProduct).Methods("POST")
	r.HandleFunc("/Product/{id}",handlerpr.UpdateProduct).Methods("PUT")
	r.HandleFunc("/Product/{id}",handlerpr.DeleteProduct).Methods("DELETE")
	http.ListenAndServe(":8080",r)
}
