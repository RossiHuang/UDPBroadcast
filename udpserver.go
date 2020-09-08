package main

import (
	"fmt"
	"net"
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	//fmt.Printf("Client123 IP : %v", *addr) //{192.168.17.89 1234 }

	n, err := conn.WriteToUDP([]byte("From server port 1234 : Hello I got your mesage with UDP"), addr)
	//fmt.Println("Error", err)
	fmt.Printf("n = %d\n", n)
	fmt.Printf("Client IP : %v", *addr) //{192.168.17.89 1234 }
	if err != nil {
		//fmt.Printf("Couldn't send response %v", err)
		fmt.Println("Couldn't send response ", err)
	}
}

// Start :
func main() {

	//make a slice of byte of length 2048
	buffer := make([]byte, 2048)
	// addr := net.UDPAddr{
	// 	Port: 1234,
	// 	//IP:   net.ParseIP("127.0.0.1"),
	// 	//Server reviced broadcast 不用指定IP
	// 	//IP: net.ParseIP("255.255.255.255"),
	// }
	//SerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	SerAddr, err := net.ResolveUDPAddr("udp", ":1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	ServerConn, err := net.ListenUDP("udp", SerAddr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, udpaddr, err := ServerConn.ReadFromUDP(buffer)
		fmt.Printf("%v\n", udpaddr)
		fmt.Printf("Read a message from -- %v : %s \n", udpaddr, buffer)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		//Goroutines
		go sendResponse(ServerConn, udpaddr)
	}

}
