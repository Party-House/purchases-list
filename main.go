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
    Bought bool     `bson: "bought" json:"bought"`
}

func main() {
    session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
    if err != nil {
        panic(err)
    }
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", GetPurchaseList(session))
    router.HandleFunc("/add", PostPurchaceItem(session)).Methods("POST")
    router.HandleFunc("/{purchaseId}/bought", UpdatePurchase(session)).Methods("POST")
    log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func GetPurchaseList(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func (w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()
        c := session.DB(os.Getenv("MONGODB_DB")).C(os.Getenv("MONGODB_COLLECTION"))
        var result []Purchase
        c.Find(bson.M{"bought": false}).All(&result)
        json.NewEncoder(w).Encode(result)
    }
}

func PostPurchaceItem(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()
        c := session.DB(os.Getenv("MONGODB_DB")).C(os.Getenv("MONGODB_COLLECTION"))
        var purchase Purchase
        decoder := json.NewDecoder(r.Body)
        decoder.Decode(&purchase)
        c.Insert(&purchase)
        w.WriteHeader(http.StatusCreated)
    }
}

func UpdatePurchase(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()
        c := session.DB(os.Getenv("MONGODB_DB")).C(os.Getenv("MONGODB_COLLECTION"))
        vars := mux.Vars(r)
        purchaseId := vars["purchaseId"]
        c.UpdateId(
            bson.ObjectIdHex(purchaseId),
            bson.M{"$set": bson.M{"bought": true}})
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, bson.ObjectIdHex(purchaseId))
    }
}