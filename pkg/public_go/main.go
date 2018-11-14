package main

import (
	"flag"
	"fmt"
	"context"
	"log"

	"github.com/soushin/bazel-multiprojects/pkg/common_go/util"
)


func main() {
	greet := flag.String("greet", "greet", "greet string")
	flag.Parse()

	greetUsecase, err := initializeGreetUsecase(context.Background(), *greet)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(util.Add(fmt.Sprintf("Hello World! %s", greetUsecase.Greet)))
}
