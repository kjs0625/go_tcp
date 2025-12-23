# Go TCP Server & Client 

## Prerequisites
- **Go** (Version 1.20 or higher recommended)
- **Git**

## Installation
```Bash
go mod init tcp
go mod tidy```

## 📡 Protocol Specification

The communication is based on a custom binary protocol with **Little Endian** byte order.

### Packet Structure

| Field | Type | Size | Description |
| :--- | :--- | :--- | :--- |
| **Packet Size** | `uint32` | 4 bytes | Length of the payload |
| **Packet Type** | `uint16` | 2 bytes | ID of the service/command |
| **Payload** | `[]byte` | Variable | Actual data (JSON, Text, etc.) |

**Total Header Size:** 6 bytes
