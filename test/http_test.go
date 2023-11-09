package test

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {

	fmt.Println("please visit http://127.0.0.1:123435")
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		s := fmt.Sprintf("你好，世界！-- Time:%s", time.Now().String())
		w.Write([]byte(s))
		fmt.Println(w, "%v\n", s)
		log.Printf("%v\n", s)
	})
	if err := http.ListenAndServe(":12345", nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}

}
