package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/i-coder-robot/go-grpc-http-rest-microservice-todo/conf"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	address := flag.String("server", "http://localhost:"+conf.HttpPort, "HTTP gateway url, e.g. http://localhost:"+conf.HttpPort)
	flag.Parse()

	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	var body string

	// Call Create
	resp, err := http.Post(*address+"/v1/todo", "application/json", strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
			"toDo": {
				"title":"title (%s)",
				"description":"description (%s)",
				"reminder":"%s"
			}
		}
	`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("Create失败: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("Create response body读取失败: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// parse ID of created ToDo
	var created struct {
		API string `json:"api"`
		ID  string `json:"id"`
	}

	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("JSON response of Create method  反序列化失败: %v", err)
		fmt.Println("error:", err)
	}

	// Call Read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID))
	if err != nil {
		log.Fatalf("Read method失败: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("Read response body失败: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Update
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID),
		strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
			"toDo": {
				"title":"title (%s) + updated",
				"description":"description (%s) + updated",
				"reminder":"%s"
			}
		}
	`, pfx, pfx, pfx)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Update method失败: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("Update response body失败: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call ReadAll
	resp, err = http.Get(*address + "/v1/todo/all")
	if err != nil {
		log.Fatalf("ReadAll method失败: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("ReadAll response body失败: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Delete
	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Delete method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("Delete response body失败: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

}
