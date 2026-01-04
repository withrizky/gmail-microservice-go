
# High-Performance Gmail Microservice (Go)

A high-throughput, non-blocking email delivery service built with **Golang**. This service solves the daily sending limit problem of standard Gmail accounts by implementing a smart **Round-Robin Account Rotation** strategy.

It utilizes the **Worker Pool pattern** and **Buffered Channels** to handle thousands of concurrent email requests efficiently without relying on external message brokers like RabbitMQ or Redis.

## 🏗 Architecture

The system uses an **In-Memory Event Driven** architecture. HTTP requests are immediately acknowledged and pushed into a buffered channel. Background workers (Goroutines) pick up jobs, automatically select the next available Gmail account from the pool, and execute the SMTP transaction.
```mermaid
graph LR
    User[Client / API] -- POST /send --> Gin[Gin HTTP Server]
    Gin -- Non-blocking --> Channel[Buffered Channel (RAM)]
    Channel --> Dispatcher[Dispatcher]
    Dispatcher --> Worker1[Worker 1]
    Dispatcher --> Worker2[Worker 2]
    Dispatcher --> WorkerN[Worker N...]
    Worker1 -- HTTP Request --> WAHA[WAHA Engine]
    Worker2 -- HTTP Request --> WAHA
    WAHA -- Send Message --> WhatsApp[WhatsApp Server]

```

## Key Features

* Smart Account Rotation: Automatically distributes outgoing emails across multiple configured Gmail accounts. If Account A sends email #1, Account B will send email #2. This bypasses the 500 emails/day limit per account.

* Ultra Fast Response: The API returns 202 Accepted immediately. Clients do not wait for the slow SMTP process.

* In-Memory Worker Pool: Replaces complex brokers with Go's native Channels and Goroutines for lower latency and simpler infrastructure.

* Graceful Shutdown: Ensures all active workers finish their current sending tasks before the server stops (prevents data loss).

* Clean Architecture: Structured using standard Go project layout (cmd, internal, pkg).

## 📂 Folder Structure

```
gmail_microservice/
├── cmd/
│   └── server/
│       └── main.go       # Application Entry Point
├── internal/
│   ├── model/            # Data Structures (Payloads)
│   ├── mailer/           # SMTP Logic & Header Formatting
│   └── worker/           # Worker Pool & Rotation Logic
├── .env                  # Configuration File
├── go.mod                # Go Modules
└── README.md             # Documentation

```

## 🛠 Prerequisites

* **Go** (version 1.18+)
* **Gmail App Passwords** : You must generate "App Passwords" for each Gmail account used. Do not use your login password.

* * Go to Google Account > Security > 2-Step Verification > App passwords.


## 🚀 Installation & Setup

1. **Clone the repository**
```bash
git clone https://github.com/withrizky/gmail-microservice-go.git
cd gmail-microservice-go

```


2. **Install Dependencies**
```bash
go mod tidy

```


3. **Environment Configuration**
Create a `.env` file in the root directory:
```env
PORT=8082

# Configure your accounts here (JSON Array)
GMAIL_ACCOUNTS=[{"email":"account1@gmail.com","password":"app_password_1"},{"email":"account2@gmail.com","password":"app_password_2"}]

```


4. **Run the Server**
```bash
go run cmd/server/main.go

```



## 📡 API Documentation

### Send Message

Sends a message to the processing queue.

* **URL**: `/send-email`
* **Method**: `POST`
* **Content-Type**: `application/json`

**Request Body:**

```json
{
    "to": "client@example.com",
    "subject": "Invoice #1023",
    "message": "Dear Customer, please find attached your invoice."
}

```

**Response (Success):**

```json
{
    "status": "queued",
    "message": "Pesan masuk antrean"
}

```

Status Code: `202 Accepted*`

**Response (Error):**

```json
{
    "error": "Payload invalid"
}

```

*Status Code: `400 Bad Request*`

## 📈 Performance Strategy

This service is optimized for high concurrency:

1. **Buffered Channels**: Can hold up to **5,000** (configurable) pending email in RAM.
2. **Concurrency**: Spawns **20** (configurable) concurrent workers. This means 20 messages are processed in parallel every millisecond.
3. **No I/O Blocking**: The HTTP handler does not wait for the SMTP server to respond.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
