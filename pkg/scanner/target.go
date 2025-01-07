package scanner


type Target struct {
	Ip string
	PortMin int
	PortMax int
	Protocol string
	TimeoutSecond int
	ColorLog bool
}
