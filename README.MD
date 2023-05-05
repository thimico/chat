# Go Chat Application

Author: Thiago Menezes - thg.mnzs@gmail.com

This is a simple WebSocket chat server implemented in Golang. It allows multiple clients to connect and send messages to a chat room. The server also supports a bot that fetches stock data from an API and sends the data to the chat room.

## Features

- WebSocket-based chat server
- Multiple chat rooms support
- Stock bot fetching stock data from an API
- Message broadcasting to all connected clients
- Message ordering by timestamps
- Display of the last 50 messages

## Requirements

- Golang 1.17 or later
- React (front)
## Installation

1. Clone the repository:

```sh
git clone https://github.com/yourusername/chat.git
cd chat
```

2. Install dependencies:

```sh
go mod download
```

3. Configure the application:

Set the following environment variables in your system:

```
# Server Configuration
PORT=8080

# Stock API Configuration
STOCK_API_BASE_URL=https://stooq.com/q/l/
STOCK_API_PARAMS=?s=%s&f=sd2t2ohlcv&h&e=csv
```

Replace the placeholder values with Stooq API URL, Port if needed.


4. Run the application:

```sh
go run main.go
```

The server will start listening on port 8080 by default. You can change the port by setting the `PORT` environment variable.

### 5. Connect to the server

You can use a WebSocket client to connect to the server at `ws://localhost:8080/ws`. Alternatively, you can use the provided React-based frontend to interact with the chat server.


## up the front-end on webapp (is require nodejs installed)

```sh
cd webapp
npm install && npm start
```

The chat server will start, and you can access the chatroom at `http://localhost:3000`.
put username, Room name and message...

if you want request the stock use this command, For example:


```
/stock=aapl.us
```

The stock quote bot will fetch the stock quote and display it in the chatroom.


## Running Tests

To run tests for the project, use the following command:

```
go test ./...
```

This command will run all tests in the project.



## Project Structure

1. On root have a main.go and main_test.go and Dockerfile and docker-compose.yml

2. Webapp folder contains the webapplication in react


## Author

Thiago Menezes <thg.mnzs@gmail.com>