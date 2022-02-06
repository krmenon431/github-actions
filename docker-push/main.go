package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"golang.org/x/net/context"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func main() {
	fmt.Println("calling main function")

	DOCKER_FILE_NAME := os.Getenv("INPUT_DOCKER_FILE_NAME")
	APP_ROOT_PATH := os.Getenv("app_root_path")
	IMAGE_NAME := os.Getenv("image_name")
	IMAGE_TAG := os.Getenv("image_tag")
	REGISTRY := os.Getenv("registry")
	REGISTRY_USERNAME := os.Getenv("registry_username")
	REGISTRY_PASSWORD := os.Getenv("registry_password")

	fmt.Println("docker-file-name:", DOCKER_FILE_NAME)
	IMAGE := REGISTRY_USERNAME + "/" + IMAGE_NAME + ":" + IMAGE_TAG

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// cli := client.NewEnvClient()
	buildErr := buildDockerImage(dockerClient, APP_ROOT_PATH, DOCKER_FILE_NAME, IMAGE)
	if buildErr != nil {
		fmt.Println(err.Error())
		return
	}

	pushErr := pushToRegistry(dockerClient, REGISTRY, REGISTRY_USERNAME, REGISTRY_PASSWORD, IMAGE)
	if pushErr != nil {
		fmt.Println(err.Error())
		return
	}

}

func buildDockerImage(dockerClient *client.Client, APP_ROOT_PATH, DOCKER_FILE_NAME, IMAGE string) error {

	fmt.Println("calling buildDockerImage")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions(APP_ROOT_PATH, &archive.TarOptions{})
	if err != nil {
		return err
	}

	imageBuildOpts := types.ImageBuildOptions{
		Dockerfile: DOCKER_FILE_NAME,
		Tags:       []string{IMAGE},
		Remove:     true,
	}
	imageBuildResponse, err := dockerClient.ImageBuild(ctx, tar, imageBuildOpts)
	if err != nil {
		return err
	}

	err = printResponse(imageBuildResponse.Body)
	if err != nil {
		return err
	}
	defer imageBuildResponse.Body.Close()

	return nil
}

func printResponse(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func pushToRegistry(dockerClient *client.Client, REGISTRY, REGISTRY_USERNAME, REGISTRY_PASSWORD, IMAGE string) error {

	fmt.Println("calling buildDockerImage")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	var authConfig = types.AuthConfig{
		Username:      REGISTRY_USERNAME,
		Password:      REGISTRY_PASSWORD,
		ServerAddress: REGISTRY,
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	imagePushResponse, err := dockerClient.ImagePush(ctx, IMAGE, opts)
	if err != nil {
		return err
	}

	printResponse(imagePushResponse)
	defer imagePushResponse.Close()

	return nil

}
