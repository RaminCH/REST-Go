package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

//Article stuct ...
type Item struct {
	Id     int    `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
	Price  string `json:"price"`
}

type ErrorMessage struct {
	Message string `json:"Message"`
}

//Items - local DataBase
var Items []Item

//GET request for /items
func GetAllItems(w http.ResponseWriter, r *http.Request) {

	if len(Items) < 1 {
		w.WriteHeader(http.StatusNotFound) // Изменить статус код запроса на 404
		var erM = ErrorMessage{Message: "Error: No one item exists"}
		json.NewEncoder(w).Encode(erM)
	} else {
		json.NewEncoder(w).Encode(Items)
	}

}

//GET request for article with ID
func GetItemWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	find := false
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}

	for _, item := range Items {
		if item.Id == id {
			find = true
			json.NewEncoder(w).Encode(item)
		}
	}
	if !find {
		w.WriteHeader(http.StatusNotFound) // Изменить статус код запроса на 404
		var erM = ErrorMessage{Message: "Error: This item doesn't exist"}
		json.NewEncoder(w).Encode(erM)
	}
}

//PostNewArticle func for create new article
func PostNewItem(w http.ResponseWriter, r *http.Request) {
	// {
	// 	"Id" : "3",
	// 	"Title" : "Title from json POST method",
	// 	"Author" : "Me",
	// 	"Content" : "Content from json POST method"
	// }
	reqBody, _ := ioutil.ReadAll(r.Body)
	var item Item
	json.Unmarshal(reqBody, &item)    // Считываем все из тела зпроса в подготовленный пустой объект Article
	w.WriteHeader(http.StatusCreated) // Изменить статус код запроса на 201
	Items = append(Items, item)
	fmt.Println(Items)
	json.NewEncoder(w).Encode(item) //После добавления новой статьи возвращает добавленную
}

//PutExistsArticle ....
func PutExistsItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// СТРОКА
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}
	find := false

	for index, item := range Items {
		if item.Id == id {
			find = true
			reqBody, _ := ioutil.ReadAll(r.Body)
			w.WriteHeader(http.StatusAccepted)     // Изменяем статус код на 202
			json.Unmarshal(reqBody, &Items[index]) // перезаписываем всю информацию для статьи с Id
		}
	}

	if !find {
		w.WriteHeader(http.StatusNotFound) // Изменяем статус код на 404
		var erM = ErrorMessage{Message: "Item with this id doesn't exist"}
		json.NewEncoder(w).Encode(erM)
	}

}

//DeleterArticleWithId ...
func DeleterItemWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}

	find := false

	for index, item := range Items {
		if item.Id == id {
			find = true
			w.WriteHeader(http.StatusAccepted) // Изменить статус код на 202
			Items = append(Items[:index], Items[index+1:]...)
		}
	}
	if !find {
		w.WriteHeader(http.StatusNotFound) // Изменить статус код на 404
		var erM = ErrorMessage{Message: "Item with that id doesn't exist"}
		json.NewEncoder(w).Encode(erM)
	}

}

func main() {
	//Добавляю 2 статьи в свою базу
	Items = []Item{
		// Item{id: 1, item: "chair", amount: 20, price: "123.123"},
		// Item{id: 2, item: "table", amount: 40, price: "321.321"},
	}
	fmt.Println("REST API V2.0 worked....")
	//СОздаю свой маршрутизатор на основе либы mux
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/items", GetAllItems).Methods("GET")
	myRouter.HandleFunc("/item/{id}", GetItemWithId).Methods("GET")
	// //Создадим запрос на добавление новой статьи
	myRouter.HandleFunc("/item", PostNewItem).Methods("POST")

	// //Создадим запрос на удаление статьи (гарантировано существует)
	myRouter.HandleFunc("/item/{id}", DeleterItemWithId).Methods("DELETE")
	myRouter.HandleFunc("/item/{id}", PutExistsItem).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
