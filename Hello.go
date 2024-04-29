// WEBAPPS MONITOR

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"strconv"
)

const monitorCount = 6
// const delay = 5 * time.Second
const delay = 10
	
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
			fmt.Println("---<<< Program End!")
			os.Exit(0)
		default:
			fmt.Println("Command not recognized!")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	version := 1.1
	fmt.Println("Hello There!")
	fmt.Println("This program is on", version, "version.")
	fmt.Println("What would you like to do?")
	fmt.Println("")
}

func showOptions() {
	fmt.Println("1- Monitoring Start")
	fmt.Println("2- Show logs")
	fmt.Println("3- Program exit")
	fmt.Println("")
}

func readComand() int {
	var comand int
	fmt.Scan(&comand) // neste contexto o & é um ponteiro para a memória
	fmt.Println("Command selected:", comand)
	fmt.Println("")

	return comand
}

func startMonitoring() {
	fmt.Println("---<<< Monitoring >>>---")
	fmt.Println("")

	sites := getSitesFromFile()

	// First for is monitoring the URLs count times
	// and testing each site each 10 minutes
	for i := 0 ; i < monitorCount ; i++ {
		fmt.Println("---<<< Monitoring", i+1,"/", monitorCount, "next rout in", delay, "minutes...")
		for i, site := range sites {
			testSite(i, site)
		}
		time.Sleep(delay * time.Minute)
		fmt.Println("")
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

	file.WriteString("_online: " + strconv.FormatBool(status) + " - site:" + site + "\n")

	file.Close()
}

func launchError(errorMessage string, errorSistem error) {
	fmt.Println("---<<<", errorMessage, ">>>---")
	fmt.Println(errorSistem)
	fmt.Println("")
}