package main

/*
 * ToTo - Go language proxy server
 * Song Hyeon Sik <blogcin@naver.com> 2016
 * @from https://github.com/blogcin/ToTo
 */

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	bufferLength = 8192
	headerLine   = 30
)

// ProxyServer stores the configuration for ProxyServer and some functions.
type ProxyServer struct {
	port string
}

// ask port
func (ps *ProxyServer) askPort() int {
	port := 0

	fmt.Print("Port : ")
	_, _ = fmt.Scanf("%d", &port)

	return port
}

// initialize socket server
func (ps *ProxyServer) init(port int) net.Listener {
	ps.port = ":"
	ps.port += strconv.Itoa(port)

	server, err := net.Listen("tcp", ps.port)

	if server == nil {
		panic("init: port listening error : " + err.Error())
	}

	fmt.Println("server listening on", ps.port)
	return server
}

// accept client
func (ps *ProxyServer) acceptClient(server net.Listener) chan net.Conn {
	channel := make(chan net.Conn)

	go func() {
		for {
			client, err := server.Accept()
			if client == nil {
				fmt.Println("[ERROR] acceptClient: Couldn't accept : ", err.Error())
				continue
			}
			channel <- client
		}
	}()

	return channel
}

// connect host by socket
func (ps *ProxyServer) connectHost(client net.Conn) {
	HeaderInfo, Datas := ps.getData(client)

	if (Datas[0] == 0) || (HeaderInfo == "-1") {
		return
	}

	ps.logf("DEBUG", "parse header line: %s, body:\n%s", HeaderInfo, string(Datas))
	requestType, host, _, port := ps.parseHTTPHeaderMethod(HeaderInfo)

	if port == -1 {
		fmt.Println(client.LocalAddr(), client.RemoteAddr(), )
		// client.Write([]byte("parse header error"))
		_= client.Close()
		return
	}

	connectionHost, _ := net.Dial("tcp", host+":"+strconv.Itoa(port))

	if requestType == "CONNECT" {
		connectionHost.Write([]byte("HTTP/1.1 200 Connection established\n"))
	} else {
		connectionHost.Write(Datas)
		go func() {
			io.Copy(connectionHost, client)
		}()
		io.Copy(client, connectionHost)
	}
	client.Close()
	connectionHost.Close()
	return
}

// Get first line of http header
func (ps *ProxyServer) getData(client net.Conn) (string, []byte) {
	buffer := make([]byte, bufferLength)
	client.Read(buffer)

	return ps.splitHeader(buffer)[0], buffer
}

// parse HTTPHeader from packets
func (ps *ProxyServer) parseHTTPHeaderMethod(headerMethod string) (string, string, string, int) {
	var (
		requestType string
		host        string
		protocol    string
		port        int
	)
	errorCounter := false

	// ex: GET http://google.com/ HTTP/1.1
	temp := headerMethod[strings.Index(headerMethod, " ")+1:]
	protocol = temp[strings.Index(temp, " ")+1:]

	url := temp[:strings.Index(temp, " ")]

	i := strings.Index(url, "://")

	if i == -1 {
		ps.logf("ERROR", "Un-correct URL: %s", url)
		errorCounter = true
	} else {
		host = url[i+3 : len(url)-1]
	}

	if errorCounter == false {
		i = strings.Index(host, ":")
		if i == -1 {
			port = 80
		} else {
			port, _ = strconv.Atoi(host[i+1:])
			host = host[:i]
		}

		i = strings.Index(host, "/")
		if i != -1 {
			host = host[:i]
		}
		requestType = headerMethod[:strings.Index(headerMethod, " ")]

		return requestType, host, protocol, port
	}

	return "", "", "", -1
}

// split header each lines
func (ps *ProxyServer) splitHeader(bytearray []byte) []string {
	result := make([]string, headerLine)
	lineNumber := 0
	temp := false

	if bytearray[0] == 0 {
		fmt.Println("ps: splitHeader: Couldn't get http header, zero filter")
		result[0] = string("-1")
		return result
	}

	for index, element := range bytearray {
		if element == '\r' {
			if bytearray[index+1] == '\n' {
				temp = true
			}
		}

		if temp != true {
			result[lineNumber] += string(element)
		}

		if element == '\n' {
			temp = false
			lineNumber++
		}
	}

	return result
}

func (ps *ProxyServer) logf(level string, format string, v ...interface{}) {
	fmt.Printf("[" + strings.ToUpper(level) + "] " + format + "\n", v ...)
}

func main() {
	var port int

	flag.IntVar(&port, "port", 0, "listen port")
	flag.IntVar(&port, "p", 0, "listen port")
	flag.Parse()

	proxyServer := &ProxyServer{}

	if port == 0 {
		port = proxyServer.askPort()
	}

	server := proxyServer.init(port)
	defer server.Close()

	connections := proxyServer.acceptClient(server)

	for {
		go proxyServer.connectHost(<-connections)
	}
}
