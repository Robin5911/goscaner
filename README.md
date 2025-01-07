## goscanner
![Static Badge](https://img.shields.io/badge/os-Linux/Unix-green)
![Static Badge](https://img.shields.io/badge/go-1.21-blue)

#### 特点:
- **功能: 支持扫描CIDR全网段**
- **功能: 支持tcp/udp双协议**
- **性能: 秒级扫描,速度快**


#### 用法

```bigquery

#./goscanner tcp 192.168.1.100/32 80-443

#查询结果红绿颜色显示
#./goscanner udp 192.168.1.0/32 701-703 color
```
#### 效果:
```
udp 192.168.1.107:701 is closed!
udp 192.168.1.107:702 is closed!
udp 192.168.1.107:703 is open!
```
