# Twitter Monitor Project

This project is a Twitter monitoring tool that allows you to track and receive real-time updates on tweets from specified Twitter users. It utilizes a WebSocket server to enable real-time communication with connected clients and provides a continuous monitoring mechanism to check for new tweets from predefined Twitter accounts.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [WebSocket Endpoint](#websocket-endpoint)
- [Configuration](#configuration)
- [Dependencies](#dependencies)

## Installation

To install and run the Twitter Monitor, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/twitter-monitor.git

2. Navigate to the project directory:

   cd twitter-monitor

3. Build and run the project:

   go build
   ./twitter-monitor

## Usage

   After running the application, the WebSocket server will be accessible at ws://localhost:8080/ws. Clients can connect to this endpoint to receive real-time    updates on monitored Twitter accounts.
   
   The Twitter Monitor continuously checks for new tweets from specified users and broadcasts them to all connected clients.

## WebSocket Endpoint

   The WebSocket endpoint for connecting to the Twitter Monitor is:

   ws://localhost:8080/ws

# Configuration

   The monitoring settings and Twitter account information are configured in the `TwitterMonitor` struct within the `main.go` file. Update the following fields with your specific values:
   
   - `MonitoredPath`: Path to the file where recent tweets are stored.
   - `RecentPath`: Path to the file containing the recently monitored tweets.
   - `Token`: Twitter API Bearer token for authentication.
   - `Cookies`: Your Twitter account cookies.
   - `CSRF`: Your Twitter account CSRF token.
   
   Add or modify Twitter users in the `users` slice, specifying their `Name` (screen name) and `ID` (user ID).

## Dependencies

   The project uses the following third-party packages:

   - [github.com/gorilla/websocket](https://github.com/gorilla/websocket): Package for WebSocket functionality.

   You can install these dependencies using the following:

   ```bash
   go get -u github.com/gorilla/websocket

