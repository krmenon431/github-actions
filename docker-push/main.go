package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("calling main function")
	fmt.Println("docker-file-name:", os.Getenv("INPUT_DOCKER_FILE_NAME"))
}
