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
	"bytes"
	"mime/multipart"
	"io/ioutil"
	"path/filepath"
)

const sampleRate = 88125
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

		fmt.Printf("\n%T\n",filename)

		audiofile, err := os.Create(filename + ".wav")
		if err != nil {
			panic(err)
		}

		binary.Write(audiofile, binary.LittleEndian, buffer)

		fmt.Println("yha tak sb thik hai")

		audiofile.Close()

		sendAudioEngine(filename)

		fmt.Println("yha tak sb thik hai 2")

		defer audiofile.Close()
	}
	defer r.Body.Close()
}

func sendAudioEngine(filename string) {

	url := "http://localhost:8080/detectSentiment"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("./getaudio/"+filename+".wav")
	defer file.Close()
	part1,
	errFile1 := writer.CreateFormFile("audio_file",filepath.Base("./getaudio/"+filename+".wav"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}


	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
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

