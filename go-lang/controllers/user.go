package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vaibhav96-stack/go-lang/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

type UserCont struct {
	session *mgo.Session
}

func ReplicaController(s *mgo.Session) *UserCont {
	return &UserCont{s}
}

func (uc UserCont) Getuser(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id){
		w.WriteHeader(404)
	}
	oid := bson.ObjectIdHex(id)
	 u := models.User{}

	 err := uc.session.DB("insta-data").C("u-details").FindId(oid).One(&u)
	 if err != nil {
		w.WriteHeader(404)
		return 
	 }
	 um, err := json.Marshal(u)
	 if err != nil { 
		fmt.Println(err)
	 }
	 w.Header().Set("Content-Type", "application/json")
	 w.WriteHeader(http.StatusOK)
	 fmt.Fprintf(w, "%s\n", um)
}





func (uc UserCont) Make_user(w http.ResponseWriter, r *http.Request , _ httprouter.Params){
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	u.Id = bson.NewObjectId()
	u.Password, _ = HashPassword(u.Password)
	println(u.Password)
	uc.session.DB("insta-data").C("u-details").Insert(u)
	um, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n",um)
}

type PostController struct {
	session *mgo.Session
}

func ReplicaPostController(s *mgo.Session) *PostController {
	return &PostController{s}
}

func (pc PostController) Get_post(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id){
		w.WriteHeader(http.StatusNotFound)
	}
	object_id := bson.ObjectIdHex(id)
	 pt := models.Post{}

	 err := pc.session.DB("insta-data").C("posts").FindId(object_id).One(&pt)
	 if err != nil {
		w.WriteHeader(404)
		return 
	 }
	 pm, err := json.Marshal(pt)
	 if err != nil { 
		fmt.Println(err)
	 }
	 w.Header().Set("Content-Type", "application/json")
	 w.WriteHeader(http.StatusOK)
	 fmt.Fprintf(w, "%s\n", pm)
}

func (pc PostController) Make_post(w http.ResponseWriter, r *http.Request , _ httprouter.Params){
	pt := models.Post{}
	json.NewDecoder(r.Body).Decode(&pt)
	pt.Id = bson.NewObjectId()
	pt.Time = time.Now()
	pc.session.DB("insta-data").C("posts").Insert(pt)
	pm, err := json.Marshal(pt)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n",pm)
}


func (pc PostController) Get_user_post(w http.ResponseWriter,r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	posts := make([]models.Post, 0, 2)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	err := pc.session.DB("insta-data").C("posts").Find(bson.M{"userid":id}).Skip((page-1)*2).Limit(2).All(&posts)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	pm, err := json.Marshal(posts)
	 if err != nil {
		fmt.Println(err)
	 }
	 w.WriteHeader(http.StatusOK)
	 fmt.Fprintf(w, "%s\n", pm)
}
