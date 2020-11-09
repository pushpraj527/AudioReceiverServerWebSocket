package main

import (

	"fmt"
	"net/http"
//	"flag"
	"log"
//	"github.com/gorilla/websocket"
	"io/ioutil"

)


const sampleRate = 44100
const seconds = 4


func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Responce from Server, Connected  to server"))
}

func serveWs(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Connected to Audio Server"))
	/*ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()*/
	//buffer := make([]float32, sampleRate * seconds)
	bodyByte, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	fmt.Printf(string(bodyByte)) //Only for test

}


func main() {
	fmt.Println("Server is running....")

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/wsaudio", serveWs)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
