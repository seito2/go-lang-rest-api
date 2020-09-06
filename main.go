package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"./db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/t0w4/toDoListBackend/controller"
)

type Article struct {
	ID       int    `json:id`
	Title    string `json:title`
	Author   string `json:author`
	PostDate string `json:year`
}

var articles []Article

func getArticles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(articles)
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("Get article is called")
	// get http://localhost:8000/books/hoge -> hoge を取得
	params := mux.Vars(r)
	log.Println(params) // map[id:1]

	// /What's params type?
	log.Println(reflect.TypeOf(params["id"])) // -> Get String

	// Convert Type from String -> Int
	// Not handling err -> _
	i, _ := strconv.Atoi(params["id"])

	// URL に指定した ID の情報を取得
	for _, article := range articles {
		if article.ID == i {
			json.NewEncoder(w).Encode(&article)
		}
	}
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("Add article is called")
	var article Article
	// json -> struct
	json.NewDecoder(r.Body).Decode(&article)
	fmt.Println("article: ", article)

	articles = append(articles, article)

	// struct -> json
	json.NewEncoder(w).Encode(articles)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("Update article is called")
	var article Article
	json.NewDecoder(r.Body).Decode(&article)

	for i, item := range articles {
		if item.ID == article.ID {
			articles[i] = article
		}
	}

	json.NewEncoder(w).Encode(article)
}

func removeArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("Remove article is called")
	params := mux.Vars(r)
	fmt.Println("params: ", params)

	id, _ := strconv.Atoi(params["id"])
	fmt.Println("id: ", id)

	fmt.Println("articles: ", articles)

	for i, item := range articles {
		if item.ID == id {
			articles = append(articles[:i], articles[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(articles)
}

func main() {
	// リクエストを裁くルーターを作成
	router := mux.NewRouter()

	// エンドポイント
	// router.HandleFunc("/articles", getArticles).Methods("GET")
	// router.HandleFunc("/articles/{id}", getArticle).Methods("GET")
	// router.HandleFunc("/articles", addArticle).Methods("POST")
	// router.HandleFunc("/articles", updateArticle).Methods("PUT")
	// router.HandleFunc("/articles/{id}", removeArticle).Methods("DELETE")

	// articles = append(articles,
	// 	Article{ID: 1, Title: "Article1", Author: "Gopher", PostDate: "2019/1/1"},
	// 	Article{ID: 2, Title: "Article2", Author: "Gopher", PostDate: "2019/2/2"},
	// 	Article{ID: 3, Title: "Article3", Author: "Gopher", PostDate: "2019/3/3"},
	// 	Article{ID: 4, Title: "Article4", Author: "Gopher", PostDate: "2019/4/4"},
	// 	Article{ID: 5, Title: "Article5", Author: "Gopher", PostDate: "2019/5/5"},
	// )

	dbConn, err := db.Init()
	if err != nil {
		log.Printf("db init failed: %v", err)
		os.Exit(1)
	}

	tc := &controller.TaskController{dbConn}
	router.HandleFunc("/tasks", tc.CreateTask).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/tasks", tc.GetTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{uuid}", tc.GetTask).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{uuid}", tc.PutTask).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{uuid}", tc.DeleteTask).Methods(http.MethodDelete, http.MethodOptions)
	log.Print(http.ListenAndServe("0.0.0.0:8080", router))
	os.Exit(1)

	// Start Server
	log.Println("Listen Server ....")
	// 異常があった場合、処理を停止したいため、log.Fatal で囲む
	log.Fatal(http.ListenAndServe(":8000", router))
}
