package scanner

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"sync"
	"time"
)

func ListenICMPUnreachable(timeoutSecond int) ([]byte, error) {
	//监听ipv4的icmp报文
	c, _ := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	buf := make([]byte, 1500)
	//设置超时时间
	_ = c.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeoutSecond)))
	//读取会话中信息
	n, _, err := c.ReadFrom(buf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ERRO - happen unknown err : %s",err.Error()))
	}
	//解析报文内容
	msg, err := icmp.ParseMessage(ipv4.ICMPTypeDestinationUnreachable.Protocol(), buf[:n])
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	//如果报文类型为icmp不可达类型的报文则返回报文内容
	if msg.Type == ipv4.ICMPTypeDestinationUnreachable {
		body := msg.Body.(*icmp.DstUnreach)
		return body.Data, nil
	}
	//如果在会话中没有读取到任何内容则返回空
	return nil, fmt.Errorf("no ICMP Destination Unreachable message")
}
func ParseUnreachUDP(unreachData []byte) Target {
	//解析udp不可达的报文
	ipHeader, err := ipv4.ParseHeader(unreachData)
	if err != nil {
		fmt.Printf("Failed to parse IP header:%s", err.Error())
		return Target{}
	}
	//fmt.Println("IP header:")
	//fmt.Printf("Version: %d\n", ipHeader.Version)
	//fmt.Printf("ID: %d\n", ipHeader.ID)
	//fmt.Printf("Protocol: %d\n", ipHeader.Protocol)
	//fmt.Printf("Src: %s\n", ipHeader.Src.String())
	//fmt.Printf("Dst: %s\n", ipHeader.Dst.String())
	//创建一个目标的结构体实例
	var result Target
	result.Ip = ipHeader.Dst.String() //头部的目标地址则为我们探测的目标地址
	//创建一个切片存放数据，长度为ipv4报文包头的长度
	dataBytes := unreachData[ipv4.HeaderLen:]
	// 解析 UDP 数据包
	var udpHeader []byte
	//出去前面8个字节的ipv4头
	udpHeader = append(udpHeader, dataBytes[:8]...)
	//srcPort := strconv.Itoa(int(binary.BigEndian.Uint16(udpHeader[:2])))
	//目的端口的位置，实测得到
	dstPort := int(binary.BigEndian.Uint16(udpHeader[2:4]))

	// 打印 UDP 数据包
	//fmt.Printf("UDP header:\n")
	//fmt.Printf("SrcPort: %s\n", srcPort)
	//fmt.Printf("DstPort: %s\n", dstPort)
	result.PortMin = dstPort
	result.PortMax = dstPort
	return result
}
func (t *Target) UDP() []int {
	var openPorts []int
	var wg sync.WaitGroup
	for p:= t.PortMin;p<=t.PortMax;p++ {
		wg.Add(1)
		port := p
		go func(p int) {
			defer wg.Done()
			//var msg string
			addrAndPort := fmt.Sprintf("%s:%d", t.Ip, p)
			conn, _ := net.Dial("udp", addrAndPort)
			// 发送数据
			data := []byte("HelloCheck")
			_, err := conn.Write(data)
			if err !=nil {
				//msg = fmt.Sprintf("ERRO - Write data into udp conn %s:%d , err : %s\n",t.Ip,p,err.Error())
				//fmt.Println(msg)
				return
			}
			// 开启监听ICMP不可达包
			unreach, err := ListenICMPUnreachable(t.TimeoutSecond)
			if err != nil {
				return
			}
			res := ParseUnreachUDP(unreach)
			if res.Ip == t.Ip && res.PortMin == port {
				openPorts= append(openPorts, port)
			}
		}(port)
	}
	wg.Wait()

	return openPorts
}
