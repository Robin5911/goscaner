package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func (t *Target) TCP() []int{
	var openPorts []int
	var wg sync.WaitGroup
	for p:= t.PortMin;p<=t.PortMax;p++ {
		wg.Add(1)
		port := p
		go func(p int) {
			defer wg.Done()
			addrAndPort := fmt.Sprintf("%s:%d",t.Ip,p)
			conn,err := net.DialTimeout("tcp",addrAndPort, time.Duration(t.TimeoutSecond) *time.Second)
			if err != nil {
				return
			}
			//  关闭连接
			conn.Close()
			openPorts= append(openPorts, p)
		}(port)
	}
	wg.Wait()
	return openPorts
}
