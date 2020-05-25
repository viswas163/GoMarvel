package v1

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Hello : Says hello to the user
func Hello(w http.ResponseWriter, req *http.Request) {
	log.Println("/hello API called")
	fmt.Fprintf(w, "Hello Agent!")
	log.Println("/hello API call successful!")
}

// Headers : Shows the request headers info
func Headers(w http.ResponseWriter, r *http.Request) {
	log.Println("/headers API called")
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	log.Println("/headers API call successful!")
}

// Image : Displays an image in root dir
func Image(w http.ResponseWriter, r *http.Request) {
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
