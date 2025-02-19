# Go Epoll Server

This Go project demonstrates the use of epoll (both edge-triggered and level-triggered modes) for efficient handling of multiple TCP connections in a server. The server listens for incoming connections and uses the epoll mechanism to handle them in an efficient, non-blocking way.

## Features
- TCP server using Go
- Implements both edge-triggered and level-triggered modes of epoll
- Uses `golang.org/x/sys/unix` for epoll functionality
- Designed for Linux-based systems (does not work on Windows)

## Requirements
- Linux (tested on Kali Linux in a VMware virtual machine)
- Go 1.20+

## How to Run
1. **Clone the repository**:
    ```bash
    git clone https://github.com/username/go-epoll.git
    cd go-epoll
    ```

2. **Build and run the program**:
    ```bash
    go build
    ./go-epoll
    ```

3. **Test using Netcat**:
    In another terminal, you can test the server by connecting to it with Netcat:
    ```bash
    nc localhost 8080
    ```

4. Type some text to simulate a client sending data to the server. The server will output the received data based on whether it's running in edge-triggered or level-triggered mode.

## Limitations
- **Windows Compatibility**: This project is designed for Linux systems and will **not work** on Windows because the `epoll` mechanism is specific to Linux. Windows uses a different I/O multiplexing system (e.g., IOCP) that is not compatible with `epoll`.
- **Virtualization**: The project was developed and tested in Kali Linux running in a VMware virtual machine, ensuring a compatible Linux environment.

## Epoll Overview
Epoll is a high-performance I/O event notification mechanism in Linux. It is designed to handle multiple file descriptors (such as network sockets) efficiently, without the need for blocking I/O operations. 

- **Edge-Triggered Mode (EPOLLET)**: The system notifies the program only when an event changes state (e.g., data becomes available on a socket).
- **Level-Triggered Mode (EPOLLIN)**: The system continues to notify the program as long as the event condition persists (e.g., data is available to read from a socket).

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
# Epoll-tcp
Swarch course A little assigment
