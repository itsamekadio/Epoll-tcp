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
# Conceptual Distinction: Level Triggered (LT) vs. Edge Triggered (ET)

At a systems design level, **Level Triggered (LT)** and **Edge Triggered (ET)** represent two distinct models of event notification semantics in I/O multiplexing. These modes define how the operating system communicates I/O readiness to an application — particularly in the context of network servers, high-concurrency applications, or any system performing asynchronous I/O.

- **Level Triggered mode** follows the principle of **state persistence**. The kernel continuously notifies the application as long as the resource (e.g., a socket) remains in a ready state. This allows the application multiple opportunities to process events, even if it misses or delays handling them.

- **Edge Triggered mode**, by contrast, is **state-transition driven**. The kernel emits an event **only once** when a resource changes from not-ready to ready. If the application fails to process the data in that moment, it receives no further notification until another state change occurs.

This core difference directly impacts how developers must design the I/O loop and buffer handling, and has serious implications for system **performance**, **fault tolerance**, and **scalability**.

---

## Software Architecture Implications

###  Level Triggered (LT)

From an architectural standpoint, **LT favors fault tolerance and maintainability**. It supports a **fail-soft** approach, where transient I/O handling errors (e.g., not reading the full buffer) are **recoverable**. The system continues notifying the application until the resource is no longer ready, allowing **multiple recovery opportunities**.

LT aligns well with design goals such as:

- Robustness  
-  Developer friendliness  
-  Observability  
-  Operational safety  

However, this comes with a **performance trade-off**:

-  Higher syscall overhead, as the kernel keeps notifying the app even if it's already aware of the readiness state.  
-  Increased context switching, particularly under high load.  
-  Spurious wake-ups, where the application is notified unnecessarily.  
-  Lower scalability, since the system generates more events per connection.  

LT is well-suited for:

-  Systems with **moderate concurrency** (e.g., hundreds to low thousands of connections).  
-  Applications where **correctness and fault tolerance** outweigh latency or throughput.  
-  Prototyping, academic learning, and early-stage architecture development.  
-  Teams that want to avoid the complexity of aggressive event-loop control.  

---

### Edge Triggered (ET)

In contrast, **ET aligns with high-performance, event-driven architectures**. It emphasizes **throughput**, **resource efficiency**, and **horizontal scalability** — essential in domains like real-time messaging, HTTP proxies, or any system operating at web scale.

ET emits notifications **only when a new event occurs**, minimizing the number of user-kernel transitions. This makes it an excellent choice for reducing syscall volume and CPU usage — especially when handling **thousands or millions** of sockets concurrently.

But this performance gain comes at a cost:

-  **Non-blocking I/O is mandatory**, and you must loop until `EAGAIN` / `EWOULDBLOCK` to drain buffers.  
-  If data isn’t fully read when the event fires, the app won’t be re-notified — leading to **data starvation** or **stalled connections**.  
-  Bugs in event handling can silently kill parts of the I/O path — this makes ET **less fault-tolerant** and harder to debug.  

ET is architecturally preferred when:

-  The system is built for **massive concurrency** and **low latency**.  
-  The developer team has the **expertise** to build robust, non-blocking I/O handling.  
-  You need **tight control** over performance bottlenecks and kernel interaction.  
-  Minimizing system call frequency and CPU overhead is critical to system health.  

---

##  When to Use Which?

This is the heart of the matter — and often the most misunderstood.

###  Use **Level Triggered (LT)** if:

- Your system must **survive developer mistakes** in event handling.  
- You're working on a **small to medium-scale application** (e.g., an internal tool, API server, academic project).  
- **Correctness and resilience** are more important than absolute performance.  
- You want to **minimize debugging complexity** and avoid "invisible failures."  
- You don’t need to squeeze every drop of throughput out of the system.  

> **In short:** LT is the safer, more fault-tolerant option, and is often “good enough” unless you know you’re going to hit high concurrency and performance limits.

---

###  Use **Edge Triggered (ET)** if:

- You are building a **high-throughput, latency-sensitive system** (e.g., a web server, real-time broker, or load balancer).  
- **Scalability** is critical — you expect to handle **tens of thousands or millions** of concurrent connections.  
- You want to **minimize syscalls and CPU load** per connection.  
- Your I/O handling is **mature, rigorously tested**, and fully non-blocking.  
- You can tolerate a **less forgiving environment** in exchange for raw efficiency.  

> **In short:** ET is the high-performance, high-risk option. You must design defensively and implement proper read/write loops — but you gain significant performance and scalability benefits when done correctly.

---

##  Performance vs. Fault Tolerance: A Design Trade-Off

Choosing between LT and ET is a classic **architecture trade-off**:

- **ET** favors **performance, scalability, and resource efficiency** — but sacrifices **fault tolerance and simplicity**.  
- **LT** favors **resilience, debuggability, and developer safety** — but may **underperform** in large-scale, high-traffic environments.  

It mirrors other system design trade-offs:

- **Consistency vs. Availability**  
- **Efficiency vs. Maintainability**  
- **Speed vs. Simplicity**

If you’re building a **fault-tolerant general-purpose service**, **LT** will likely meet your needs with fewer pitfalls. But if you're scaling up and **performance bottlenecks** matter more than operational safety, **ET** becomes the necessary — though unforgiving — choice.

---

##  Summary

- **LT** is architecturally **safer but less efficient**.  
  Ideal for **small to medium workloads**, and where **correctness and ease of implementation** matter more than raw performance.

- **ET** is highly **efficient but less forgiving**.  
  Ideal for **systems at scale**, where performance must be **maximized**, and developers can enforce **strict discipline** around I/O.

---

> **Bottom line:**  
> **Use LT** when you want **fault tolerance** and **developer simplicity**.  
> **Use ET** when you want **extreme performance** and can **manage the complexity**.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
# Epoll-tcp
Swarch course A little assigment

