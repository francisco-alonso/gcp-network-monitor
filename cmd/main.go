/*
This script is a network scanner that listens for ARP replies (Address Resolution Protocol). It captures devices on your local network and logs their IP addresses and MAC addresses to Google Cloud Logging
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const projectID = "auto-local-network-monitor"
const logName = "home-network-scan"

// Function to send logs to Google Cloud
func sendToCloudLogging(ip, mac string) {
	// Create a logging client
	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}
	defer client.Close()

	logger := client.Logger(logName)
	logger.Log(logging.Entry{
		Timestamp: time.Now(),
		Severity:  logging.Info,
		Payload:   fmt.Sprintf("Discovered Device -> IP: %s, MAC: %s", ip, mac),
	})

	fmt.Println("✅ Sent to Google Cloud Logging!")
}

func main() {
    if len(os.Args) < 2 {
		fmt.Println("❌ Error: Please provide a network interface name.")
		fmt.Println("🔹 Usage: go run main.go <interface>")
		return
	}
	device := os.Args[1]
	// Open network interface for packet capture
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal("Error opening device: ", err)
	}
	defer handle.Close()

	fmt.Println("Listening for ARP replies...")

	// Set up a packet source to read packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Process packets
	for packet := range packetSource.Packets() {    
        // Check if it's an ARP packet
        arpLayer := packet.Layer(layers.LayerTypeARP)
        if arpLayer != nil {
            arpPacket, _ := arpLayer.(*layers.ARP)
    
            fmt.Printf("ARP Packet - Operation: %d, Sender IP: %s, Target IP: %s\n",
                arpPacket.Operation,
                net.IP(arpPacket.SourceProtAddress).String(),
                net.IP(arpPacket.DstProtAddress).String(),
            )
    
            // Check if it's an ARP reply (Operation = 2)
            if arpPacket.Operation == 2 {
                senderIP := net.IP(arpPacket.SourceProtAddress).String()
                senderMAC := net.HardwareAddr(arpPacket.SourceHwAddress).String()
    
                fmt.Printf("📡 Device Found! IP: %s | MAC: %s\n", senderIP, senderMAC)
                
                sendToCloudLogging(senderIP, senderMAC)
            }
        }
    }
    
}