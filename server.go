package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
	"io"
	"bytes"
	"strings"
	"strconv"
)

type Neighbor struct {
	Ip				string
	P2p_port	int
	User_port	int
}

type Config struct {
	P2p_port			int
	User_port			int
	Neighbor_list	[]Neighbor
	Target				string
}

const protocol = "tcp"
var nodeList []string
var config Config

func startServer() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Panic(err)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &config)

	ln, err := net.Listen(protocol, "0.0.0.0:" + strconv.Itoa(config.User_port))
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn)
	}
}

func sendData(addr string, data []byte) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func getUserPort(addr string) string {
	for _, neighbor := range config.Neighbor_list {
		if addr == neighbor.Ip {
			return strconv.Itoa(neighbor.User_port)
		}
	}
	return ""
}

func getP2pPort(addr string) string {
	for _, neighbor := range config.Neighbor_list {
		if addr == neighbor.Ip {
			return strconv.Itoa(neighbor.P2p_port)
		}
	}
	return ""
}

func handleGetBlockCount(conn net.Conn, rpc map[string]interface{}) {
	addr := conn.RemoteAddr().String()
	addr = strings.Split(addr, ":")[0]
	addr = addr + ":" + getUserPort(addr)

	var height int
	var bc Blockchain
	bc.loadFromFile()
	if len(bc.Blocks) > 0 {
		height = bc.Blocks[len(bc.Blocks) - 1].Height
	} else {
		height = 0
	}
	
	jsonData := map[string]interface{}{
		"error": 0,
		"result": height,
	}
	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Panic(err)
	}

	sendData(addr, data)
}

func handleConnection(conn net.Conn) {
	rpcBytes, _ := ioutil.ReadAll(conn)
	var rpc map[string]interface{}
	json.Unmarshal(rpcBytes, &rpc)

	// Test data
	// fmt.Println(rpc["method"], rpc["data"].(map[string]interface{})["height"])

	// IP
	//fmt.Println(conn.RemoteAddr().String())

	if rpc["method"] == "getBlockCount" {
		handleGetBlockCount(conn, rpc)
	}

	conn.Close()
}
