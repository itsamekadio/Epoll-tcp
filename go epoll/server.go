package main

import (
    "fmt"
    "net"
    "os"
    "golang.org/x/sys/unix"
)

const (
    serverAddress = "localhost:8080"
)

func handleEdgeTriggered(connFd int) {
    epollFd, err := unix.EpollCreate1(0)
    if err != nil {
        fmt.Println("Error creating epoll:", err)
        return
    }

    event := unix.EpollEvent{Events: unix.EPOLLIN | unix.EPOLLET, Fd: int32(connFd)}
    err = unix.EpollCtl(epollFd, unix.EPOLL_CTL_ADD, connFd, &event)
    if err != nil {
        fmt.Println("Error adding fd to epoll:", err)
        return
    }

    events := make([]unix.EpollEvent, 10)
    for {
        n, err := unix.EpollWait(epollFd, events, -1)
        if err != nil {
            fmt.Println("Error in EpollWait:", err)
            return
        }

        for i := 0; i < n; i++ {
            if events[i].Events&unix.EPOLLIN != 0 {
                buffer := make([]byte, 1024)
                n, err := unix.Read(connFd, buffer)
                if err != nil {
                    fmt.Println("Error reading from connection:", err)
                    return
                }
                if n > 0 {
                    fmt.Printf("Edge Triggered Mode: Received: %s\n", string(buffer[:n]))
                }
            }
        }
    }
}

func handleLevelTriggered(connFd int) {
    epollFd, err := unix.EpollCreate1(0)
    if err != nil {
        fmt.Println("Error creating epoll:", err)
        return
    }

    event := unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(connFd)}
    err = unix.EpollCtl(epollFd, unix.EPOLL_CTL_ADD, connFd, &event)
    if err != nil {
        fmt.Println("Error adding fd to epoll:", err)
        return
    }

    events := make([]unix.EpollEvent, 10)
    for {
        n, err := unix.EpollWait(epollFd, events, -1)
        if err != nil {
            fmt.Println("Error in EpollWait:", err)
            return
        }

        for i := 0; i < n; i++ {
            if events[i].Events&unix.EPOLLIN != 0 {
                buffer := make([]byte, 1024)
                n, err := unix.Read(connFd, buffer)
                if err != nil {
                    fmt.Println("Error reading from connection:", err)
                    return
                }
                if n > 0 {
                    fmt.Printf("Level Triggered Mode: Received: %s\n", string(buffer[:n]))
                }
            }
        }
    }
}

func main() {
    listener, err := net.Listen("tcp", serverAddress)
    if err != nil {
        fmt.Println("Error starting server:", err)
        os.Exit(1)
    }
    defer listener.Close()

    conn, err := listener.Accept()
    if err != nil {
        fmt.Println("Error accepting connection:", err)
        return
    }
    defer conn.Close()

    // Get the file descriptor
    tcpConn, ok := conn.(*net.TCPConn)
    if !ok {
        fmt.Println("Error: failed to convert connection to TCPConn")
        return
    }

    file, err := tcpConn.File()
    if err != nil {
        fmt.Println("Error getting file from TCP connection:", err)
        return
    }
    connFd := int(file.Fd())  // Now we correctly handle the file descriptor

    // Start handling in separate goroutines
    go handleEdgeTriggered(connFd)
    go handleLevelTriggered(connFd)

    // Block to keep the server running
    select {}
}
