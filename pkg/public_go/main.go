package main

import (
	"flag"
	"fmt"
	"context"
	"log"
	"net/http"

	"github.com/soushin/bazel-multiprojects/pkg/common_go/util"
)

var msg string

func main() {
	greet := flag.String("greet", "Hello", "greet message")
	flag.Parse()
	greetUsecase, err := initializeGreetUsecase(context.Background(), *greet)
	if err != nil {
		log.Fatalln(err)
	}
	msg = greetUsecase.Msg

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, util.Add(fmt.Sprintf("%s Go!", msg)))
}
