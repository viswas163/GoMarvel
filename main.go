package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello Agent!")
}

func headers(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func image(w http.ResponseWriter, r *http.Request) {
	fileName := "YuGarden.jpg"
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Open() :", errors.New("Error opening image file with name = "), fileName)
		return
	}

	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("ReadAll() :", errors.New("Error reading file. Check it once"))
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(content)
}

func main() {

	fmt.Println("Starting Server...")

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/image", image)

	fmt.Println("Started Server... OK!")

	http.ListenAndServe(":3001", nil)
}
