package main

import (
	"goscanner/pkg/scanner"
)



func main(){
	t := scanner.Target{
		Ip: "10.16.15.77",
		Port: "6688",
		TimeoutSecond: 3,
		ColorRead: true,
	}
	//t.TCP()
	t.UDP()
}

