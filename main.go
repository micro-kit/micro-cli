package main

import (
	"log"
	"os"

	"github.com/micro-kit/micro-cli/program"
)

var (
	VERSION  string // 程序版本
	GIT_HASH string // git hash
)

//https://github.com/hashicorp/consul/blob/master/main.go
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// 应用实例
	p := program.New(VERSION, GIT_HASH)
	os.Exit(p.Run())
}
