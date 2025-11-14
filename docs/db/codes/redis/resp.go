package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
)

func sendRequest(args ...string) {
	_, err := writer.WriteString("*" + strconv.Itoa(len(args)) + "\r\n")
	if err != nil {
		panic(err)
	}
	for _, arg := range args {
		_, err = writer.WriteString("$" + strconv.Itoa(len(arg)) + "\r\n")
		if err != nil {
			panic(err)
		}
		_, err = writer.WriteString(arg + "\r\n")
		if err != nil {
			panic(err)
		}
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

func handleResponse() any {
	prefix, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	switch prefix {
	case '*':
		return nil
	case '+':
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		return string(line)
	case '-':
		return nil
	case ':':
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		n, err := strconv.Atoi(string(line))
		if err != nil {
			panic(err)
		}
		return n
	case '$':
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		n, err := strconv.Atoi(string(line))
		if err != nil {
			panic(err)
		}
		builder := strings.Builder{}
		for i := 0; i < n; i++ {
			line, err := reader.ReadByte()
			if err != nil {
				panic(err)
			}
			builder.WriteByte(line)
		}
		_, _, err = reader.ReadLine()
		if err != nil {
			panic(err)
		}
		return builder.String()
	default:
		panic("错误数据格式")
	}
}

func main() {
	// 1. 建立连接
	var err error
	conn, err = net.Dial("tcp", "10.120.100.27:6379")
	if err != nil {
		panic(err)
	}

	// 5. 释放连接
	defer func() {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	// 2. 获取输出流和输入流
	reader = bufio.NewReader(conn)
	writer = bufio.NewWriter(conn)

	// 3. 发出请求
	sendRequest("SET", "name", "张三")

	// 4. 解析响应
	res := handleResponse()
	fmt.Printf("%#v\n", res) // "OK"

	// 测试
	sendRequest("GET", "name")
	res = handleResponse()
	fmt.Printf("%#v\n", res) // "张三"

	sendRequest("SET", "age", "20")
	res = handleResponse()
	fmt.Printf("%#v\n", res) // "OK"

	sendRequest("GET", "age")
	res = handleResponse()
	fmt.Printf("%#v\n", res) // "20"

	sendRequest("OBJECT", "ENCODING", "age")
	res = handleResponse()
	fmt.Printf("%#v\n", res) // "int"
}