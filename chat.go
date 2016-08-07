package main

import (
	//"github.com/kataras/iris"
	//"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
	"log"
	"encoding/json"
	"fmt"
)


type Message struct {
	Username  	string `json:"username"`
	Message 	string `json:"message"`
}

type Login struct {
	NumUsers 	int `json:"numUsers"`
}

type UserJoinLeft struct {
	Username  	string `json:"username"`
	NumUsers 	int `json:"numUsers"`
}


type UserName struct {
	Username  	string `json:"username"`
}

type clientPage struct {
	Title string
	Host  string
}

func main() {


	iris.Static("/stuff", "./public", 1)

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Render("index.html", clientPage{"Client Page", ctx.HostString()})
	})

	iris.Config.Websocket.Endpoint = "/my_endpoint"
	//iris.Config().Websocket.Headers["Access-Control-Allow-Origin"] = "*"

	var myChatRoom = "room1"
	var users = make(map[string]string)


	/////

	var numUsers int
	numUsers = 0

	iris.Websocket.OnConnection(func(c iris.WebsocketConnection) {


		var addedUser = false
		log.Println("\nOn connection")

		c.Join(myChatRoom)

		c.On("add user", func(username string) {
			fmt.Printf("add userrr, %s\n",username)

			if (addedUser == false) {

				users[c.ID()] = username
				numUsers += 1
				addedUser = true

				users := Login{NumUsers:numUsers}
				users_json, _ := json.Marshal(users)

				c.Emit("login",users_json)


				info := UserJoinLeft{
					Username:username,
					NumUsers:numUsers,
				}

				infojson, _ := json.Marshal(info)

				c.To(iris.Broadcast).Emit("user joined",infojson)

			}


		})


		c.On("new message",func(message string){

			log.Println("new message")
			log.Printf("%+v",message)

			username_text := users[c.ID()]

			msgk := Message{Username:username_text,Message:message}
			msg_json, _ := json.Marshal(msgk)

			//c.To(myChatRoom).Emit("new message",msg_json)

			c.To(iris.Broadcast).Emit("new message",msg_json)

		})

		c.On("typing",func() {

			log.Println("typing")

			username_text := users[c.ID()]

			username := UserName{Username:username_text}
			username_json, _ := json.Marshal(username)

			//c.To(myChatRoom).Emit("typing",username_json)
			c.To(iris.Broadcast).Emit("typing",username_json)

		})

		c.On("stop typing",func() {

			log.Println("stop typing")

			//username_text := users[c.ID()]

			//username := UserName{Username:username_text}
			//username_json, _ := json.Marshal(username)

			//c.To(myChatRoom).Emit("stoptyping",username_json)
			//c.To(iris.Broadcast).Emit("stoptyping",username_json)

		})

		c.OnDisconnect(func() {

			if (addedUser == true) {

				numUsers -= 1

				username_text := users[c.ID()]

				info := UserJoinLeft{Username:username_text,NumUsers:numUsers}
				infojson, _ := json.Marshal(info)

				c.To(iris.Broadcast).Emit("user left",infojson)

			}

			log.Println("on disconnect")
			fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())

		})


		//c.On("chat", func(message string) {
		//	// to all except this connection ->
		//	//c.To(websocket.Broadcast).Emit("chat", "Message from: "+c.ID()+"-> "+message)
		//
		//	// to the client ->
		//	//c.Emit("chat", "Message from myself: "+message)
		//
		//	//send the message to the whole room,
		//	//all connections are inside this room will receive this message
		//	username := users[c.ID()]
		//
		//	c.To(myChatRoom).Emit("chat", username  + ": "+message)
		//})

	})

	iris.Listen(":8080")


}
