package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const CICLESTOTEST = 100
const SECONDSTOWAIT = 5

func main() {
	for {
		showMenu()
		option := readOption()
		executeOption(option)
	}
}

func showMenu() {
	fmt.Println()
	fmt.Println(">>>Monitor de Sites<<<")
	fmt.Println()
	fmt.Println("1- Iniciar monitoramento")
	fmt.Println("2- Exibir logs")
	fmt.Println("0- Sair")
	fmt.Print("Escolha uma opção:")
}

func readOption() int {
	var option int
	fmt.Scan(&option)
	return option
}

func executeOption(option int) {
	switch option {
	case 1:
		startMonitoring()
	case 2:
		showLogs()
	case 0:
		fmt.Println("Saindo do programa...")
		os.Exit(0)
	default:
		fmt.Println("Opção inexistente")
		os.Exit(-1)
	}
}

func startMonitoring() {
	fmt.Println("Iniciando monitoramento...")

	urls := readURLFromFile()
	for i := 1; i <= CICLESTOTEST; i++ {
		fmt.Println("---Teste", i, "---")
		for _, url := range urls {
			runUrlTest(url)
		}
		fmt.Println()
		time.Sleep(SECONDSTOWAIT * time.Second)
	}

	fmt.Println("---Fim dos testes---")
}

func runUrlTest(url string) {
	response, err := http.Get(url)

	if err == nil {
		if response.StatusCode == 200 {
			fmt.Println("ONLINE:", url)
			writeLog(url, true, strconv.Itoa(response.StatusCode))
		} else {
			fmt.Println(">>>OFFLINE!<<<:", url, "STATUSCODE:", response.StatusCode)
			writeLog(url, false, strconv.Itoa(response.StatusCode))
		}
	} else {
		fmt.Println(">>>OFFLINE!<<<:", url, "ERRO:", err)
		writeLog(url, true, err.Error())
	}
}

func readURLFromFile() []string {
	var urls []string
	file, err := os.Open("urls.txt")
	if err != nil {
		fmt.Println("Erro:", err)
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		urls = append(urls, line)

		if err == io.EOF {
			break
		}
	}
	file.Close()

	fmt.Println()
	fmt.Println("Preparando para monitorar:", urls)
	fmt.Println()
	return urls
}

func showLogs() {
	file, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(file))
}

func writeLog(url string, isOnline bool, statusCode string) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	onlineStatus := "OFFLINE"
	if isOnline {
		onlineStatus = "ONLINE"
	}

	logLine := "[" + time.Now().Local().Format("02/01/2006 15:04:05") + "]:" + onlineStatus + "-" + statusCode + "-" + url
	file.WriteString(logLine + "\n")

	file.Close()
}
