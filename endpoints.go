package main

import (
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	_"github.com/lib/pq"
)

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request){
	/*
	params := mux.Vars(req)
	for _, item := range people{
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
	*/

	params := mux.Vars(req)
	rows, e := db.Query("INSERT INTO people(id, first_name, last_name) VALUES ($1, $2, $3)", 1, params["fname"], params["lname"])
	if e != nil {
		log.Fatal(e)
	}
	json.NewEncoder(w).Encode(rows)
}

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request){
	json.NewEncoder(w).Encode(people)
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	for index, item := range people {
		if item.ID == params["id"]{
			people = append(people[:index], people[index+1:]...)
			return
		}
	}
	json.NewEncoder(w).Encode(people)
}

//-------------------- Rewriting Server ----------------------------------------------------------

func ApiMeEndpoint(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	rows, e := db.Query("select * from users where github_id = $1", params["user"])
	json.NewEncoder(w).Encode(rows)
}

func ApiPutUserMethod(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var user UserUpdate
	_ = json.NewDecoder(req.Body).Decode(&user)
	user.id = params["id"]
	rows, e := db.Query("update users set linkedin_url = $1, twitter_url = $2, portfolio_url = $3, about_me = $4 where github_id = $5", user.linkedin_url, user.twitter_url, user.portfolio_url, user.about_me, user.id)
	json.NewEncoder(w).Encode(rows)
}

func ApiAllUsersEndpoint(w http.ResponseWriter, req *http.Request){
	rows, e := db.Query("select * from users")
	json.NewEncoder(w).Encode(rows)
}

func ApiChannelUsersByChannelID(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select * from channel_uers join users on channel_users.user_id = users.github_id where channel_users.channel_id = $1", params["channel_id"])
	json.NewEncoder(w).Encode(r)
}

func ApiCreateServer(w http.ResponseWriter, req *http.Request){
	var server_name string `json:"server_name,omitempty"`
	params := mux.Vars(req)
	_ = json.NewDecoder(req.Body).Decode(&server_name)
	r, e := db.Query("insert into server(admin_github_id, server_name) values ( $1, $2 ); select * from server where admin_github_id = $1", params["user"], server_name)
	json.NewEncoder(w).Encode(r)
}

func ApiServers(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select * from server where admin_github_id = $1", params["user"])
	json.NewEncoder(w).Encode(r)
}

func ApiDeleteServer(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("delete from server where admin_github_id = $1 and id = $2", params["user"], params["id"])
	json.NewEncoder(w).Encode(r)
}

func ApiDeleteChannel (w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select * from channel_users join channel on channel_users.channel_id = channel.id join server on channel.server_id = server.id where user_id = $1", params["user"])
	var flag bool = false
	for index, item := range r{
		if(r[index].admin_github_id == params["user"] && r[index].channel_id == params["id"]){
			flag = true 
			break
		}
	}
	if flag == true {
		row, err := db.Query("delete from channel where id = $1", params["id"])
		json.NewEncoder(w).Encode(row)
		return
	}
	json.NewEncoder(w).Encode(&ErroredResponse{ message: "Can't complete that response"})
}

func ApiCreateChannel(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var c ChannelCreaterData
	_ = json.NewDecoder(req.Body).Decode(&c)
	r, e := db.Query("insert into channel(server_id, channel_name, channel_color) values ( $1, $2, $3);", c.server_id, c.channel_name, c.channel_color)
	ro, er := db.Query("select id from channel where server_id = $1 and channel_name = $2;", c.server_id, c.channel_name)
	rows, err := db.Query("insert into channel_users(user_id, channel_id) values ( $1, $2 );", params["user"], ro[0].id)
	json.NewEncoder(w).Encode(&SuccessMessage{message: "Successfully created channel"})
}

func ApiMyServers(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select * from channel_users join channel on channel_users.channel_id = channel.id join server on channel.server_id = server.id where user_id = $1", params["user"])
	json.NewEncoder(w).Encode(r)
}

func ApiMyChannelsByServerId(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select * from channel where server_id = $1", params["id"] )
	json.NewEncoder(w).Encode(r)
}

func ApiMyChannels(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select * from channel join channel_users on channel.id = channel_users.channel_id where user_id = $1;", params["user"])
	json.NewEncoder(w).Encode(r)
}


func ApiMyServersAdmin (w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select * from server where admin_github_id = $1", params["user"])
	json.NewEncoder(w).Encode(r)
}

func ApiAddChannelUser(w http.ResponseWriter, req *http.Request){

}

func ApiChannelPermissions(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var cid ChannelId
	_ = json.NewDecoder(req.Body).Decode(&cid)
	r, e := db.Query("select * from channnel_users where user_id = $1 and channel_id = $2;", params["user"], cid.channelId)
	var myflag Flag
	if r[0] == true {
		json.NewEncoder(w).Encode(&myflag{ flag: true})
		return
	}
	json.NewEncoder(w).Encode(&myflag{ flag: false})
}

func ApiMessagesByChannelId(w http.ResponseWriter, req *http.Request){
	
}

func ApiMessages(w http.ResponseWriter, req *http.Request){
	
}



