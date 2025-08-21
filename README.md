# GoChat 🚀

GoChat is a simple, real-time chat application built with **Go** and **WebSockets**. It provides a basic yet powerful example of how to handle real-time communication between a server and multiple clients.

## Features ✨

* **Real-time Messaging:** Users can send and receive messages instantly without needing to refresh the page.
* **WebSocket Technology:** Leverages the `gorilla/websocket` package for efficient, full-duplex communication.
* **Simple UI:** A straightforward, clean user interface makes it easy to use.
* **Concurrent Handling:** The Go server is designed to handle multiple client connections concurrently, ensuring good performance and scalability.

---

## Getting Started 🏁

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

* Go (version 1.18 or higher)
* A web browser

### Installation

1.  Clone the repository:
    ```bash
    git clone [https://github.com/your-username/gochat.git](https://github.com/your-username/gochat.git)
    cd gochat
    ```
2.  Install the dependencies:
    ```bash
    go mod tidy
    ```
3.  Run the server:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

---

## How it Works 🧠

The application uses a **single Go routine** to manage a broadcast channel for all incoming messages. When a client connects via a WebSocket, a new Go routine is spawned to handle that specific connection.

* **Client to Server:** Messages sent from a client are received by the server-side Go routine.
* **Server to All Clients:** The server then broadcasts the message to the main channel, which in turn sends the message to all connected clients.

This architecture ensures that all clients stay synchronized in real time.

---

## Technologies Used 🛠️

* **Go:** The primary language for the backend server.
* **`gorilla/websocket`:** A popular Go package for implementing WebSockets.
* **HTML/CSS/JavaScript:** Used for the front-end user interface.