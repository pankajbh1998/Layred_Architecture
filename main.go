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
	db,_:=sql.Open("mysql","Pankaj:Pankaj@123@tcp(127.0.0.1)/Company")
	err:=db.Ping()
	if err != nil {
		log.Fatalf("cannot connect with database")
	}

	storeBr:=brandStore.New(db)
	storePr:=productStore.New(db)
	servicePr:=productService.New(storePr,storeBr)
	handlerPr:=productHandler.New(servicePr)
	r.HandleFunc("/product",handlerPr.GetByName).Methods("GET").Queries("name","{name}")
	r.HandleFunc("/product",handlerPr.GetById).Methods("GET").Queries("id","{id}")
	r.HandleFunc("/product",handlerPr.CreateProduct).Methods("POST")
	r.HandleFunc("/product",handlerPr.UpdateProduct).Methods("PUT").Queries("id","{id}")
	r.HandleFunc("/product",handlerPr.DeleteProduct).Methods("DELETE").Queries("id","{id}")
	_=http.ListenAndServe(":8080",r)
}
