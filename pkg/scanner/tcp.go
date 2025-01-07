package scanner

import (
	"fmt"
	"goscanner/pkg/utils"
	"net"
	"time"
)

func (t *Target) TCP() {
	var msg string
	addrAndPort := fmt.Sprintf("%s:%s",t.Ip,t.Port)
	conn,err := net.DialTimeout("tcp",addrAndPort, time.Duration(t.TimeoutSecond) *time.Second)
	if err != nil {
		msg = fmt.Sprintf("tcp %s:%s is closed!\n",t.Ip,t.Port)
		if t.ColorRead {
			msg = utils.Red(msg)
		}
		fmt.Println(msg)
	}else{
		msg = fmt.Sprintf("tcp %s:%s is open!\n",t.Ip,t.Port)
		if t.ColorRead {
			msg = utils.Green(msg)
		}
		fmt.Println(msg)
		defer conn.Close()
	}
}
