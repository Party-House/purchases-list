package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "os"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/mux"
)

type Purchase struct {
    ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
    Items []string  `bson:"Items" json:"items"`
    Bought bool     `json:"bought"`
}

func setUpConnection() *mgo.Collection{
    session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
    if err != nil {
        panic(err)
    }

    c := session.DB(os.Getenv("MONGODB_DB")).C(os.Getenv("MONGODB_COLLECTION"))
    return c
}

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", GetPurchaseList)
    router.HandleFunc("/todos", PostPurchaceItem)
    router.HandleFunc("/todos/{todoId}", UpdatePurchase)
    log.Fatal(http.ListenAndServe(":8080", router))
}

func GetPurchaseList(w http.ResponseWriter, r *http.Request) {
    var result []Purchase
    c := setUpConnection()
    c.Find(bson.M{}).All(&result)
    json.NewEncoder(w).Encode(result)
}

func PostPurchaceItem(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Todo Index!")
}

func UpdatePurchase(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    fmt.Fprintln(w, "Todo show:", todoId)
}