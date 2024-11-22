package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func main() {
	// Command-line options
	device := flag.String("device", "", "Path to the serial device (e.g., /dev/ttyUSB0 or COM3)")
	listenAddr := flag.String("listen", "0.0.0.0:5000", "TCP listen address and port")
	baudRate := flag.Int("baud", 9600, "Baud rate for the serial device")
	allowedIPs := flag.String("allowed", "", "Comma-separated list of allowed source IPs (leave empty to allow all)")
	verbose := flag.Bool("verbose", false, "Enable verbose logging for IN/OUT data")

	flag.Parse()

	// Validate the device argument
	if *device == "" {
		fmt.Println("Available serial devices:")
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			log.Fatalf("Failed to list serial ports: %v", err)
		}
		for _, port := range ports {
			fmt.Printf("- %s (isUSB: %v)\n", port.Name, port.IsUSB)
		}
		log.Fatal("Please specify a serial device using the -device flag.")
	}

	// Parse allowed IPs into a map for quick lookup
	allowedIPSet := make(map[string]bool)
	if *allowedIPs != "" {
		for _, ip := range strings.Split(*allowedIPs, ",") {
			allowedIPSet[ip] = true
		}
	}

	// Open the serial device
	mode := &serial.Mode{
		BaudRate: *baudRate,
	}
	serialPort, err := serial.Open(*device, mode)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer serialPort.Close()

	log.Printf("Serial device %s opened at %d baud", *device, *baudRate)

	// Start the TCP server
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s and forwarding to %s at %d baud", *listenAddr, *device, *baudRate)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
		if len(allowedIPSet) > 0 && !allowedIPSet[clientIP] {
			log.Printf("Connection from disallowed IP: %s", clientIP)
			conn.Close()
			continue
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		// Handle client connection
		go handleConnection(conn, serialPort, *verbose)
	}
}

func handleConnection(conn net.Conn, serialPort serial.Port, verbose bool) {
	defer conn.Close()

	// Forward data from TCP to serial
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from TCP: %v", err)
				}
				break
			}
			if verbose {
				log.Printf("IN  (TCP->Serial): %q", buffer[:n])
			}
			_, err = serialPort.Write(buffer[:n])
			if err != nil {
				log.Printf("Error writing to serial: %v", err)
				break
			}
		}
	}()

	// Forward data from serial to TCP
	buffer := make([]byte, 1024)
	for {
		n, err := serialPort.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from serial: %v", err)
			}
			break
		}
		if verbose {
			log.Printf("OUT (Serial->TCP): %q", buffer[:n])
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Printf("Error writing to TCP: %v", err)
			break
		}
	}
}
