// WEBAPPS MONITOR

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"log"
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
			fmt.Println("---<<< Cleaning Old Logs >>>---")
			fmt.Println("")
			fmt.Println("---<<< Dropping logs...")
			loading()
			deleteOldLogs()
		case 2:
			fmt.Println("---<<< Monitoring >>>---")
			fmt.Println("")
			startMonitoring()
		case 3:
			fmt.Println("---<<< Show logs history >>>---")
			fmt.Println("")
			fmt.Println("---<<< Extrating Logs...")
			loading()
			recoverLogs()
		case 0:
			exitProgram()
		default:
			fmt.Println("Command not recognized!")
			os.Exit(-1)
		 }
	}
} 

func exitProgram() {
	fmt.Println(" ███  █ ██ ███   ███  █ ██ ███")
	fmt.Println(" █  █ █ █  █     █  █ █ █  █  ")
	fmt.Println(" ███  ██   ██    ███  ██   ██ ")
	fmt.Println(" █  █ █    █     █  █ █    █  ")
	fmt.Println(" ███  ██   ███   ███  ██   ███ see you!.")
	fmt.Println("")
	os.Exit(0)
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
	fmt.Println("1- Delete old logs")
	fmt.Println("2- Monitoring")
	fmt.Println("3- Show logs history")
	fmt.Println("0- Program exit")
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
		fmt.Println("---<<< Monitoring", i+1,"/", monitorCount, "next round in", delay, "minutes...")
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
		fmt.Println("---<<< Site", i + 1, "down!, StatusCode:", res.StatusCode, "URL:", site)
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
	file, err := os.ReadFile("logs.txt")

	if err != nil {
		launchError("Error while reading log!", err)
	}
	fmt.Println("")
	fmt.Println(string(file))
}

func deleteOldLogs() {
	filePath := "logs.txt"
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()


	// ---<<< Get the current date and subtract 2 days >>>---
	now := time.Now()
	period := now.AddDate(0, 0, -2)

	// ---<<< Subtract 2 months from the current date to get the deadline >>>---
	// now := time.Now()
	// period = now.AddDate(0, -2, 0)


	// Create a new temporary file to store valid lines
	tempFilePath := "logs_temp.txt"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Fatalf("Error creating temp file: %v", err)
	}
	defer tempFile.Close()

	// Create a scanner to read the log file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dateStr := strings.Fields(line)[0] // Extract the date part of the line
		date, err := time.Parse("02/01/2006", dateStr) // Parse date
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			continue // Ignore the line if the date cannot be parsed
		}
		
		// Compare the log date with the deadline (period)
		if date.After(period) {
			// If the log date is later than period, write to the temporary file
			_, err := tempFile.WriteString(line + "\n")
			if err != nil {
				log.Fatalf("Error writing temp file: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading logs file: %v", err)
	}

	if err := file.Close(); err != nil {
		log.Fatalf("Error closing original file: %v", err)
	}

	if err := os.Rename(tempFilePath, filePath); err != nil {
		log.Fatalf("Error renaming temp file: %v", err)
	}
	fmt.Println("")
	fmt.Println("---<<< Old logs deleted successfully! >>>---")
	fmt.Println("")
}

func loading(){
	// Open the loading messages file
	filePath := "loading_messages.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	// Read loading messages from the file
	messages := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		messages = append(messages, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}
	
	for _, message := range messages {
		fmt.Printf("\r%s", message) // Clear the previous line
		time.Sleep(1300 * time.Millisecond)
	}
	fmt.Println("")
}

func launchError(errorMessage string, errorSistem error) {
	fmt.Println("---<<<", errorMessage, ">>>---")
	fmt.Println(errorSistem)
	fmt.Println("")
}