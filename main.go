package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gamexg/proxyclient"
	"github.com/panjf2000/ants/v2"
)

var (
	file        string
	thread      int
	sucessFile  string
	proxyEnable = false
	proxyURI    string
)
var mutex sync.Mutex
var sucess_file *os.File

func init() {
	flag.StringVar(&file, "f", "ip.txt", "扫描文件")
	flag.IntVar(&thread, "t", 100, "线程")
	flag.StringVar(&sucessFile, "o", "success.txt", "输出文件")
	flag.StringVar(&proxyURI, "p", "", "代理，支持http和socks5")
}

func main() {
	start := time.Now()
	flag.Parse()
	if proxyURI != "" {
		proxyEnable = true
	}
	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("文件不存在")
		return
	}
	_, err = os.Stat(sucessFile)
	if os.IsNotExist(err) {
		ioutil.WriteFile(sucessFile, []byte{}, 0666)
	}
	sucess_file, err = os.OpenFile(sucessFile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		if err != nil {
			fmt.Println(sucessFile, "文件不存在")
			return
		}
	}
	defer sucess_file.Close()
	data := strings.Split(string(f), "\r\n")

	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(thread, func(i interface{}) {
		ipaddr := i.([]string)
		num, _ := strconv.Atoi(ipaddr[0])
		frp(num+1, ipaddr[1])
		wg.Done()
	})
	defer p.Release()
	for num, i := range data {
		l := []string{strconv.Itoa(num), i}
		wg.Add(1)
		_ = p.Invoke(l)
	}
	wg.Wait()
	cost := time.Since(start)
	fmt.Printf("总共使用了%s", cost)
}
func frp(num int, ipaddr string) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	conn, err := tcp(ipaddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	str1, _ := hex.DecodeString("000100010000000100000000")
	conn.Write(str1)
	str2, _ := hex.DecodeString("000000000000000100000094")
	conn.Write(str2)
	str3, _ := hex.DecodeString("6f000000000000008b7b2276657273696f6e223a22302e34362e30222c226f73223a2277696e646f7773222c2261726368223a22616d643634222c2270726976696c6567655f6b6579223a226130326161643733333736656435383962373239363738333262343239373166222c2274696d657374616d70223a313637343330313239372c22706f6f6c5f636f756e74223a317d")
	conn.Write(str3)
	fmt.Println(num, "发送成功")
	buff := make([]byte, 512)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Read(buff)
	_, err = conn.Read(buff)
	_, err = conn.Read(buff)
	conn.Close()

	i := strings.Index(string(buff), "{")
	j := strings.Index(string(buff), "}")
	// if i == -1 || j == -1 {
	// 	return
	// }
	jsonn := string(buff)[i : j+1]

	jsonMap := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(jsonn), &jsonMap)
	if _, ok := jsonMap["error"]; ok {
		if jsonMap["error"] != "" {
			fmt.Println(num, "不存在未授权")
			return
		}
	}
	mutex.Lock()
	write := bufio.NewWriter(sucess_file)
	write.WriteString(ipaddr + "\n")
	write.Flush()
	mutex.Unlock()
	fmt.Println(num, "[+]", ipaddr, "存在未授权")
}

func tcp(ipaddr string) (net.Conn, error) {
	if proxyEnable == true {
		p, err := proxyclient.NewProxyClient(proxyURI)
		if err != nil {
			return nil, err
		}
		conn, err := p.DialTimeout("tcp", ipaddr, 5*time.Second)
		return conn, err
	} else {
		conn, err := net.DialTimeout("tcp", ipaddr, 5*time.Second)
		return conn, err
	}
}
