package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/inancgumus/screen"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type SystemStatus struct {
	CPULoad     float64 `json:"cpu_load"`
	MemoryUsage float64 `json:"memory_usage"`
	Hostname    string  `json:"hostname"`
}

func connectToServer() (string, *websocket.Conn, error) {
	// Connect to the WebSocket server
	serverHost := "localhost:8080"
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+serverHost+"/ws", nil)
	return serverHost, conn, err
}

func main() {
	var conn *websocket.Conn
	var err error

	screen.Clear()
	for {
		screen.MoveTopLeft()

		// Attempt to establish a connection to the server
		if conn == nil {
			_, conn, err = connectToServer()
			if err != nil {
				// log.Println("Error connecting to server:", err)
				fmt.Println("Error connecting to server:", err)
				time.Sleep(1 * time.Second)
				continue
			}
			defer conn.Close()
		}

		// Get CPU load
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Println("Error getting CPU load:", err)
			continue // Skip this iteration on error
		}
		cpuLoad := cpuPercent[0]

		// Get memory usage
		virtualMemory, err := mem.VirtualMemory()
		if err != nil {
			log.Println("Error getting memory usage:", err)
			continue // Skip this iteration on error
		}
		memoryUsage := virtualMemory.UsedPercent

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("Error getting hostname:", err)
			return
		}

		// Create a SystemStatus struct with CPU load and memory usage
		status := SystemStatus{
			CPULoad:     cpuLoad,
			MemoryUsage: memoryUsage,
			// Hostname:    hostname,
			Hostname: "test",
		}

		// Convert the struct to JSON
		jsonData, err := json.Marshal(status)
		if err != nil {
			log.Println("Error encoding JSON:", err)
			continue // Skip this iteration on error
		}

		// Send the JSON data to the server
		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Println("Error sending JSON:", err)
			conn.Close() // Close the current connection
			conn = nil   // Reset the connection
			continue
		}

		fmt.Printf("%-20s %-20s %-20s\n", "Hostname", "CPU Load", "Memory Usage")
		fmt.Printf("%-20s %-20f %-20f\n", hostname, cpuLoad, memoryUsage)

		time.Sleep(1 * time.Second)
	}
}
