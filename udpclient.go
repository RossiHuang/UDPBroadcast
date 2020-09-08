package main

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	Type    string
	Message string
	IP      string
	buffer  []byte
	conn    *net.UDPConn
	Addr    *net.UDPAddr
}

func IniClient(_type string, msg string, IP string, conn *net.UDPConn) *Client {

	s := Client{Message: msg, conn: conn, IP: IP, buffer: make([]byte, 2048)}
	s.Type = "udp"

	if _type == "local" {

		var err error
		s.Addr, err = net.ResolveUDPAddr(s.Type, s.IP)
		if err != nil {
			//panic(err)
			fmt.Printf("Some error %v", err)
		}

		fmt.Printf("Client IP : %v\n", *s.Addr) //{192.168.17.89 1234 }
		//*s.conn
		// Make a connection
		s.conn, err = net.ListenUDP(s.Type, s.Addr)
		//checkError(err)
		if err != nil {
			fmt.Printf("Some error %v", err)
			//panic(err)
			//return
		}
	} else {

		var err error
		s.Addr, _ = net.ResolveUDPAddr(s.Type, s.IP)
		if err != nil {
			//panic(err)
			fmt.Printf("Some error %v", err)
		}
	}

	return &s
}

//check channel close or not
func IsClosed(ch <-chan []byte) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func (c *Client) write() {

	//println("%v", c.Addr)
	println("%v", c.Message)

	_, err := c.conn.WriteToUDP([]byte(c.Message), c.Addr)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
}

func (c *Client) read() {
	// fmt.Printf("c.conn %v\n", c.conn)
	// fmt.Printf("c.Message %v\n", c.Message)
	// fmt.Printf("c.Addr %v\n", c.Addr)
	var err error

	//timer := time.AfterFunc(1 * time.Second) //设置时间周期
	//bufChan := make(chan []byte, 2048)
	bufChan := make(chan []byte)
	go func() {
		for {
			//buf := []byte{}
			buf := make([]byte, 2048)
			_, _, err = c.conn.ReadFrom(buf)
			if err != nil {
				//fmt.Printf("Some error %v\n", err)
				//panic(err)
			}

			//檢查channel 是否被關閉
			if !IsClosed(bufChan) {
				bufChan <- buf
			}
		}
	}()

	counter := 0
	for {
		//if IsClosed(bufChan) == false {
		select {

		case a := <-bufChan:
			//_, _, err := c.conn.ReadFrom(c.buffer)
			//fmt.Println("%v", a)
			//fmt.Println("test3")
			c.buffer = a
			if err == nil {
				fmt.Printf("%s\n", c.buffer)

				//關閉channel
				if !IsClosed(bufChan) {
					close(bufChan)
					fmt.Println(IsClosed(bufChan))
				}

				//defer c.conn.Close()
				return
			} else {
				fmt.Printf("Some error %v\n", err)
			}
		default:

			time.Sleep(1 * time.Second)
			counter++
			if counter == 3 {
				return
			}
		}
		//}
	}

}

//
func IniMap(m map[string]Client) {

	m["local"] = *IniClient("local", "", ":4321", nil)
	// Resolving Address

	m["server1"] = *IniClient("server", "Hi UDP Server :1234, How are you doing?", "255.255.255.255:1234", m["local"].conn)

	m["server2"] = *IniClient("server", "Hi UDP Server :1235, How are you doing?", "255.255.255.255:1235", m["local"].conn)
}

func main() {

	timer := time.NewTicker(time.Second) //设置时间周期

	m := make(map[string]Client)

	for {
		select {
		case t := <-timer.C:
			fmt.Println(t)

			//m = make(map[string]Client)
			if m["local"].conn != nil {

				//force to close
				m["local"].conn.Close()
				//time.Sleep(30 * time.Second)
			}

			IniMap(m)
			//使用range()方法從timer.C管道中取出數據
			for k, v := range m {
				if k != "local" {
					v.write()
					v.read()
					//time.Sleep(1 * time.Second)
				}
			}
		}
	}

	//defer m["local"].conn.Close()
}
