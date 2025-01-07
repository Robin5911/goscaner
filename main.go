package main

import (
	"fmt"
	"goscanner/pkg/scanner"
	"goscanner/pkg/utils"
	"os"
	"strconv"
	"strings"
)
func run(inputDstPortMin int, inputDstPortMax int, inputProtocol string,checkTimeout int, dstIp string) []int {
	//  切片，用于保存扫描的结果
	var openPorts []int
	t := scanner.Target{
		Ip: dstIp,
		PortMin: inputDstPortMin,
		PortMax: inputDstPortMax,
		Protocol: inputProtocol,
		TimeoutSecond: checkTimeout,
		ColorLog: true,
	}
	if strings.ToLower(inputProtocol) == "tcp" {
		openPorts = t.TCP()
	}
	if strings.ToLower(inputProtocol) == "udp" {
		openPorts = t.UDP()
	}
	return openPorts
}
func main(){
	if len(os.Args) < 4 {
		fmt.Println("ERR ARGS , etc. ./goscanner tcp 192.168.1.100/24 80-443")
		return
	}
	portRangeStr := os.Args[3]
	if !strings.Contains(portRangeStr,"-") {
		fmt.Println("ERR ARGS , etc. ./goscanner tcp 192.168.1.100/24 80-443")
		return
	}
	var msg string
	inputDstCidr := os.Args[2]
	inputDstPortMinStr := strings.Split(portRangeStr,"-")[0]
	inputDstPortMaxStr := strings.Split(portRangeStr,"-")[1]
	inputDstPortMin,_ := strconv.Atoi(inputDstPortMinStr)
	inputDstPortMax,_ := strconv.Atoi(inputDstPortMaxStr)
	inputProtocol := os.Args[1]
	checkTimeout := 3  //second
	dstIps,err := utils.Hosts(inputDstCidr)
	if err != nil {
		msg = fmt.Sprintf("ERRO - INPUT CIDR %s ERR : %s",inputDstCidr,err.Error())
		fmt.Println(msg)
		return
	}
	var withColor bool
	if len(os.Args) > 4 {
		withColor=true
	}
	for _,dstIp := range dstIps {
		openPorts := run(inputDstPortMin,inputDstPortMax,inputProtocol,checkTimeout,dstIp)
		for i:=inputDstPortMin;i<=inputDstPortMax;i++ {
			if utils.IsContainInt(openPorts,i) {
				msg = fmt.Sprintf("%s %s:%d is open!",inputProtocol,dstIp,i)
				if withColor {
					msg = utils.Green(msg)
				}
				fmt.Println(msg)
			}else{
				msg = fmt.Sprintf("%s %s:%d is closed!",inputProtocol,dstIp,i)
				if withColor {
					msg = utils.Red(msg)
				}
				fmt.Println(msg)
			}
		}
	}
}

