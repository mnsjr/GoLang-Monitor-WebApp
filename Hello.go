// WEBAPPS MONITOR

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

var monitorCount = 1
var delay = 1
const version = 1.1
	
func main(){
	showIntroduction()

	for {
		showOptions()
		comand := readComand()

		switch comand {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("---<<< Extrating Logs...")
		case 3:
			fmt.Println("---<<< Program End! >>>---")
			os.Exit(0)
		default:
			fmt.Println("Command not recognized!")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	fmt.Println("")
	fmt.Println("███   ███  ███    ██  ███████  ██████  ███████")
	fmt.Println("████ ████  ████   ██  ██           ██  ██   ██")
	fmt.Println("██ ███ ██  ██  ██ ██  ███████      ██  ██ ██")
	fmt.Println("██  █  ██  ██   ████       ██  ██  ██  ██  ██")
	fmt.Println("██     ██  ██    ███  ███████   ████   ██   ██ _v",version)
	fmt.Println("")
	fmt.Println("Hello There!")
	fmt.Println("")
}

func showOptions() {
	fmt.Println("What would you like to do?")
	fmt.Println("1- Monitoring")
	fmt.Println("2- Show logs history")
	fmt.Println("3- Program exit")
	fmt.Println("")
}

func readComand() int {
	var comand int
	fmt.Scan(&comand) // neste contexto o & é um ponteiro para a memória
	// fmt.Println("Command selected:", comand)
	fmt.Println("")
	return comand
}

func setTimesMonitoring(){
	var monitorTimes int
	fmt.Println("How many times do you want to monitor?")
	fmt.Scan(&monitorTimes)
	fmt.Println("")
	monitorCount = monitorTimes
}

func setDelayMonitoring(){
	var delayMinutes int
	fmt.Println("How many minutes between each monitoring?")
	fmt.Scan(&delayMinutes)
	fmt.Println("")
	delay = delayMinutes
}

func startMonitoring() {
	fmt.Println("---<<< Monitoring >>>---")
	fmt.Println("")

	setTimesMonitoring()
	setDelayMonitoring()

	fmt.Println("Monitoring Configs:", "run", monitorCount, "times", "with", delay, "minutes delay.")
	fmt.Println("")
	fmt.Println("---<<< Monitoring Start")
	fmt.Println("")

	sites := getSitesFromFile()

	// First for is monitoring the URLs count times
	// and testing each site each 10 minutes
	for i := 0 ; i < monitorCount ; i++ {
		fmt.Println("---<<< Monitoring", i+1,"/", monitorCount, "next rout in", delay, "minutes...")
		for i, site := range sites {
			testSite(i, site)
		}
		fmt.Println("")
		time.Sleep(time.Duration(delay) * time.Minute)
	}	
	fmt.Println("---<<< Monitoring completed! >>>---")
	fmt.Println("")
	
}

func getSitesFromFile() []string {
	var sites []string

	file, err := os.Open("sites.txt")
	if err != nil {
		launchError("Error while opening file!", err)	
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n') // ReadString irá ler até o delimitador indicado, poderia ser um regex
		line = strings.TrimSpace(line)  // TrimSpace is removing /n from end of the line
		sites = append(sites, line)
		
		if err == io.EOF {  // EOF = End Of File
			break
		}
	}

	file.Close()

	return sites
}

func testSite(i int, site string) {
	res, err := http.Get(site)

	if err != nil {
		launchError("Error while geting URL request!", err)	
	}

	if res.StatusCode == 200 {
		fmt.Println("---<<< Site", i + 1, "up!", site)
		writeLog(site, true)
	} else {
		fmt.Println("---<<< Site down!, StatusCode:", res.StatusCode, "URL:", site)
		writeLog(site, false)
	}
}

func writeLog(site string, status bool) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		launchError("Error while writing log!", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + "online: " + strconv.FormatBool(status) + " - site:" + site + "\n")

	file.Close()
}

func recoverLogs() {
	file, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		launchError("Error while reading log!", err)
	}

	fmt.Println((file))
}

func launchError(errorMessage string, errorSistem error) {
	fmt.Println("---<<<", errorMessage, ">>>---")
	fmt.Println(errorSistem)
	fmt.Println("")
}