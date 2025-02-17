# 📡 ARP Network Scanner with Google Cloud Logging

This project is a **Golang-based network scanner** that listens for **ARP replies** to detect devices on your local network. It captures **IP and MAC addresses** of devices and logs them to **Google Cloud Logging** for monitoring purposes.

---

## 🛠️ **How It Works**

1. The script listens for **ARP replies** on a specified network interface.
2. It captures packets using **gopacket** and **pcap**.
3. If an **ARP reply** is detected, the script extracts:
   - **IP Address**
   - **MAC Address**
4. The detected device's information is **logged to Google Cloud Logging**.

---

## 🚀 **Installation & Setup**

### **Prerequisites**

1. **Install Go (Golang)**: [Download Here](https://go.dev/dl/)
2. **Npcap (Windows Only)**: [Download Here](https://npcap.com)
3. **Enable Google Cloud Logging API**:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Enable **Cloud Logging API** for your project
4. **Create a Service Account Key**:
   - Go to `IAM & Admin > Service Accounts`
   - Create a new service account
   - Assign the `Logging Writer` role
   - Download the JSON key file
   - Set the key in your environment:
     ```powershell
     $env:GOOGLE_APPLICATION_CREDENTIALS="C:\path\to\your-key.json"
     ```
     ```bash
     export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your-key.json"
     ```

---

## 🏃 **Running the Scanner**

### **Step 1: Find Your Network Interface**
To list available interfaces, run:

```powershell
windump -D   # Windows (Requires WinDump)
```
```bash
ifconfig      # Linux/macOS
ip link show  # Linux Alternative
```

Example output:
```
1. \Device\NPF_{12345678-ABCD-4321-9876-ABCDEF123456} (Ethernet)
2. \Device\NPF_{23456789-BCDE-5432-8765-BCDEF234567} (Wi-Fi)
```
Pick the correct one for your connection type.

### **Step 2: Run the Scanner**

```powershell
go run main.go "\Device\NPF_{YOUR_INTERFACE}"  # Windows
```
```bash
go run main.go "eth0"  # Linux/macOS
```

Example output:
```
Listening for ARP replies...
Captured an ARP packet...
📡 Device Found! IP: 192.168.1.2 | MAC: 00:1A:2B:3C:4D:5E
✅ Sent to Google Cloud Logging!
```

---

## 📝 **Code Breakdown**

### **1️⃣ Packet Capturing**
The script opens a network interface to listen for ARP replies:
```go
handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
```

### **2️⃣ Filtering ARP Packets**
It checks if the captured packet is an ARP reply:
```go
if arpPacket.Operation == 2 {
    senderIP := net.IP(arpPacket.SourceProtAddress).String()
    senderMAC := net.HardwareAddr(arpPacket.SourceHwAddress).String()
```

### **3️⃣ Sending Logs to Google Cloud**
It logs device details to Google Cloud Logging:
```go
logger.Log(logging.Entry{
    Timestamp: time.Now(),
    Severity:  logging.Info,
    Payload:   fmt.Sprintf("Discovered Device -> IP: %s, MAC: %s", ip, mac),
})
```

---

## 🛠️ **Troubleshooting**

### **1. Error: `No network interface found`**
✔ **Fix**: Ensure **Npcap** is installed and running.
```powershell
Get-Service -Name npcap
Start-Service -Name npcap
```

### **2. Error: `Failed to create logging client: could not find default credentials`**
✔ **Fix**: Set your **Google Cloud credentials** properly.
```powershell
$env:GOOGLE_APPLICATION_CREDENTIALS="C:\path\to\your-key.json"
```

### **3. Device Not Showing Up?**
✔ **Fix**: Try reconnecting the device to the network.
✔ **Fix**: Use `arp -a` to manually check devices.
