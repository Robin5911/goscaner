package scanner

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"strconv"
	"time"
)
func ListenICMPUnreachable(timeoutSecond int) ([]byte, error) {
	//监听ipv4的icmp报文
	c, _ := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	buf := make([]byte, 1500)
	//设置超时时间
	_ = c.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeoutSecond)))
	//读取会话中信息
	n, _, _ := c.ReadFrom(buf)
	//解析报文内容
	msg, err := icmp.ParseMessage(ipv4.ICMPTypeDestinationUnreachable.Protocol(), buf[:n])
	if err != nil {
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
	dstPort := strconv.Itoa(int(binary.BigEndian.Uint16(udpHeader[2:4])))

	// 打印 UDP 数据包
	//fmt.Printf("UDP header:\n")
	//fmt.Printf("SrcPort: %s\n", srcPort)
	//fmt.Printf("DstPort: %s\n", dstPort)
	result.Port = dstPort
	return result
}
func (t *Target) UDP() {
	var msg string
	addrAndPort := fmt.Sprintf("%s:%s", t.Ip, t.Port)
	conn, _ := net.Dial("udp", addrAndPort)

	// 发送数据
	data := []byte("Check Udp Port!")
	_, err := conn.Write(data)
	if err !=nil {
		msg = fmt.Sprintf("ERRO - Write data into udp conn %s:%s , err : %s\n",t.Ip,t.Port,err.Error())
		fmt.Println(msg)
		return
	}
	// 开启监听ICMP不可达包
	unreach, err := ListenICMPUnreachable(5)
	if err == nil {
		res := ParseUnreachUDP(unreach)
		if res.Ip == t.Ip && res.Port == t.Port {
			fmt.Printf("解析出来的icmp不可达包和参数传入的ip以及port均匹配%v", *t)
			return
		} else {
			fmt.Printf("虽收到icmp不可达包[%v:%v]，但和当前探测的IP和端口[%v:%v]不匹配判定UDP端口为打开状态", res.Ip, res.Port, t.Ip, t.Port)
			return
		}
	}

}
