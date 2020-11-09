package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const sampleRate = 90240000
const seconds = 5

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
	buffer := make([]byte, sampleRate*seconds)

	reader := bufio.NewReader(r.Body)

	i := 0
	for {

		i += 1
		_, err := reader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		filename := fmt.Sprintf("sample%d", i)

		audiofile, err := os.Create("./getaudio/"+filename + ".wav")
		if err != nil {
			panic(err)
		}

		binary.Write(audiofile, binary.LittleEndian, buffer)

		defer audiofile.Close()
	}
	defer r.Body.Close()
}

func main() {
	fmt.Println("Server is running....")
	//fmt.Println(time.Now().Format(time.RFC3339))
	filename := time.Now().Format(time.RFC3339)
	fmt.Println(filename)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/wsaudio", serveWs)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

