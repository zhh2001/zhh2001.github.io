package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type redisError string

func (e redisError) Error() string { return string(e) }

func sendRequest(w *bufio.Writer, args ...string) error {
	if _, err := fmt.Fprintf(w, "*%d\r\n", len(args)); err != nil {
		return err
	}
	for _, arg := range args {
		// len(string) 是字节数，参数可以含 NUL、CR、LF
		if _, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(arg), arg); err != nil {
			return err
		}
	}
	return w.Flush()
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return "", fmt.Errorf("invalid RESP line ending")
	}
	return line[:len(line)-2], nil
}

// err 表示传输或格式错误，服务端错误保留为 redisError 值。
// 数组中的错误也要完整读取，不能提前返回并留下未消费的元素。
func readResponse(r *bufio.Reader, depth int) (any, error) {
	if depth > 64 {
		return nil, fmt.Errorf("RESP nesting exceeds demo limit")
	}
	prefix, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if !strings.ContainsRune("+-:$*", rune(prefix)) {
		return nil, fmt.Errorf("unsupported RESP prefix %q", prefix)
	}
	line, err := readLine(r)
	if err != nil {
		return nil, err
	}
	switch prefix {
	case '+':
		return line, nil
	case '-':
		return redisError(line), nil
	case ':':
		return strconv.ParseInt(line, 10, 64)
	}

	n, err := strconv.ParseInt(line, 10, 64)
	if err != nil || n < -1 {
		return nil, fmt.Errorf("invalid RESP length %q", line)
	}
	if n == -1 {
		return nil, nil
	}
	if prefix == '$' {
		// 教学示例限制单个字符串为 64 MiB，区别于服务端配置
		if n > 64<<20 {
			return nil, fmt.Errorf("bulk string exceeds demo limit")
		}
		data := make([]byte, int(n)+2)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, err
		}
		if data[n] != '\r' || data[n+1] != '\n' {
			return nil, fmt.Errorf("invalid bulk string ending")
		}
		return string(data[:n]), nil
	}
	if n > 1<<20 {
		return nil, fmt.Errorf("array exceeds demo limit")
	}
	values := make([]any, int(n))
	for i := range values {
		values[i], err = readResponse(r, depth+1)
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

func main() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		panic(err)
	}
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	commands := [][]string{
		{"SET", "resp-demo:name", "zhh\x00\r\n"},
		{"SET", "resp-demo:empty", ""},
		{"DEL", "resp-demo:missing"},
		{"MGET", "resp-demo:name", "resp-demo:empty", "resp-demo:missing"},
		{"LPOP", "resp-demo:name"}, // 演示服务端错误
		{"PING"},                   // 错误之后仍能读取下一条响应
	}
	for _, args := range commands {
		if err := sendRequest(w, args...); err != nil {
			panic(err)
		}
		value, err := readResponse(r, 0)
		if err != nil {
			panic(err)
		}
		if serverErr, ok := value.(redisError); ok {
			fmt.Printf("%s error: %s\n", args[0], serverErr)
			continue
		}
		fmt.Printf("%s: %#v\n", args[0], value)
	}
}
