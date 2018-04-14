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
	var server CreateServer
	params := mux.Vars(req)
	_ = json.NewDecoder(req.Body).Decode(&server)
	r, e := db.Query("insert into server(admin_github_id, server_name) values ( $1, $2 ); select * from server where admin_github_id = $1", params["user"], server.server_name)
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
	defer r.Close()
	for r.Next(){
		var admin_github_id string
		var channel_id string
		err = r.Scan(&admin_github_id, &channel_id)
		if (admin_github_id == params["user"] && channel_id == params["id"]){
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
	defer ro.Close()
	var id string
	ro.Scan(&id)
	rows, err := db.Query("insert into channel_users(user_id, channel_id) values ( $1, $2 );", params["user"], id)
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
	params := mux.Vars(req)
	var chanInstance ChannelUserCreater
	_ = json.NewDecoder(req.Body).Decode(&chanInstance)
	if params["user"] == chanInstance.github_id{
		var e ErroredResponse
		json.NewEncoder(w).Encode(&e{message: "Cannot invite yourself to a channel"})
		return
	}
	r, e := db.Query("select * from channel_users where user_id = $1 and channel_id = $2", params["user"], chanInstance.channel_id)
	defer r.Close()
	if r.Next(){
		ro, er := db.Query("select * from channel_users where user_id = $1 and channel_id = $2", chanInstance.github_id , chanInstance.channel_id)
		if ro.Next(){
			var myResponse ErroredResponse
			json.NewEncoder(w).Encode(&myResponse{message: "Cannot add a person more that once"})
		} else {
			rows, errors := db.Query("insert into channel_users(user_id,channel_id) values ( $1, $2 );", chanInstance.github_id, chanInstance.channel_id)
			var success SuccessMessage
			json.NewEncoder(w).Encode(&success{message: "User added to channel"})
			return
		}
	}
	var myFinalResponse ErroredResponse
	json.NewEncoder(w).Encode(&myFinalResponse{message: "You do not have permisssions, as you are not it the channel"})
}

func ApiChannelPermissions(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var cid ChannelId
	_ = json.NewDecoder(req.Body).Decode(&cid)
	r, e := db.Query("select * from channnel_users where user_id = $1 and channel_id = $2;", params["user"], cid.channelId)
	defer r.Close()
	
	flag := r.Next()
	if flag == true {
		json.NewEncoder(w).Encode(&Flag{ flag: true})
		return
	}
	json.NewEncoder(w).Encode(&Flag{ flag: false})
}

func ApiMessagesByChannelId(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	r, e := db.Query("select messages.id, user_id, channel_id, content, is_code, language, loom_embed, timestamp, github_nickname, picture  from messages join users on messages.user_id = users.github_id where channel_id = $1 order by messages.id desc;", params["channelId"])
	json.NewEncoder(w).Encode(r)
}

func ApiMessages(w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var message MessagePost
	_ = json.NewDecoder(req.Body).Decode(&message)
	r, e := db.Query("insert into messages(user_id, channel_id, content, is_code, language, loom_embed values ($1, $2, $3, $4, $5, $6);", params["user"], message.channel_id, message.content, message.is_code, message.language, message.loom_embed)
	
	json.NewEncoder(w).Encode(&SuccessMessage{message: "message stored"})
}



