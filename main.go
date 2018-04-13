package main

import (
	"github.com/go-ini/ini"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_"github.com/lib/pq"
	"database/sql"
)

type Person struct{
	ID string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname string `json:"lastname,omitempty"`
	Address *Address `json:"address,omitempty"`
}



type Address struct{
	City string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

//Server Types

type ChannelId struct {
	channelId string `json:"channelId,omitempty"`
}

type UserUpdate struct {
	id string `json:"id,omitempty"`
	linkedin_url string `json:"linkedin_url,omitempty"`
	twitter_url string `json:"twitter_url,omitempty"`
	portfolio_url string `json:"portfolio_url,omitempty"`
	about_me string `json:"about_me,omitempty"`
}

type ErroredResponse struct {
	message string `json:"error,omitempty"`
}

type ChannelCreaterData struct {
	server_id string `json:"server_id,omitempty"`
	channel_name string `json:"channel_name,omitempty"`
	channel_color string `json:"channel_color,omitempty"`
}

type SuccessMessage struct {
	message string `json:"success,omitempty"`
}

type Flag struct {
	flag bool `json:"success,omitempty"`
}



var cfg, err = ini.Load("my.ini")
var db, othererr = sql.Open("postgres", cfg.Section("database").Key("connection_string").String())
var people []Person



func main(){
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "Trevor", Lastname: "Scheurer", Address: &Address{City: "Dallas", State: "Texas"}})
	people = append(people, Person{ID: "2", Firstname: "Mario", Lastname: "Hoyos"})
	people = append(people, Person{ID: "3", Firstname: "Jonathan", Lastname: "Aquino", Address: &Address{City: "Dallas", State: "Texas"}})
	people = append(people, Person{ID: "4", Firstname: "Ryan", Lastname: "Daniels", Address: &Address{City: "Dallas", State: "Texas"}})	
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{fname}/{lname}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", DeletePersonEndpoint).Methods("DELETE")


	router.HandleFunc("/api/me/{user}", ApiMeEndpoint).Methods("GET")
	router.HandleFunc("/api/user/{user}", ApiPutUserMethod).Methods("PUT")
	router.HandleFunc("/api/allUsers", ApiAllUsersEndpoint).Methods("GET")
	router.HandleFunc("/api/channel/users/{channel_id}", ApiChannelUsersByChannelID).Methods("GET")
	router.HandleFunc("/api/create/server/{user}", ApiCreateServer).Methods("POST")
	router.HandleFunc("/api/servers/{user}", ApiServers).Methods("GET")
	router.HandleFunc("/api/server/{id}/{user}", ApiDeleteServer).Methods("DELETE")
	router.HandleFunc("/api/channel/{id}/{user}", ApiDeleteChannel).Methods("DELETE")
	router.HandleFunc("/api/create/channel/{user}", ApiCreateChannel).Methods("POST")
	router.HandleFunc("/api/myServers/{user}", ApiMyServers).Methods("GET")
	router.HandleFunc("/api/myChannelsByServer/{id}", ApiMyChannelsByServerId).Methods("GET")
	router.HandleFunc("/api/myChannels/{user}", ApiMyChannels).Methods("GET")
	router.HandleFunc("/api/myServers/admin/{user}", ApiMyServersAdmin).Methods("GET")
	router.HandleFunc("/api/add/channelUser", ApiAddChannelUser).Methods("POST")
	router.HandleFunc("/api/channel/permissions/{channel_id}/{user}", ApiChannelPermissions).Methods("GET")
	router.HandleFunc("/api/messages/{channelId}", ApiMessagesByChannelId).Methods("GET")
	router.HandleFunc("/api/messages/{user}", ApiMessages).Methods("POST")

	log.Fatal(http.ListenAndServe(":12345", router))
}