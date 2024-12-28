
# Chat Service Backend

This project is a chat service backend implemented in Go, designed for managing chat rooms and real-time communication. The backend supports the following features:

## Features
- **User Management**: Create, update, delete, and retrieve user information.
- **Room Management**: Create and manage chat rooms. Users can join or leave rooms.
- **Private Messaging**: Send and receive private messages between users.
- **Broadcast Messaging**: Broadcast messages to all members of a chat room.
- **Real-time Communication**: Support for Server-Sent Events (SSE) to deliver real-time messages.

## Requirements
- **Go**: Version 1.23.2.
- **Dependencies**: Ensure the following Go libraries are installed:
  - `github.com/gorilla/mux`

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/MuhammedAshifVnr/Chat-Service.git
   cd Chat-Service
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints
### User Endpoints
1. **Create User**
   - **POST** `/users`
   - **Body**:
     ```json
     {
       "display_name": "John Doe"
     }
     ```

2. **Get User**
   - **GET** `/users?id={userID}`

3. **Update User**
   - **PUT** `/users`
   - **Body**:
     ```json
     {
       "user_id": "12345",
       "display_name": "John Updated"
     }
     ```

4. **Delete User**
   - **DELETE** `/users?id={userID}`

### Room Endpoints
1. **Create Room**
   - **POST** `/rooms`
   - **Body**:
     ```json
     {
       "room_name": "General"
     }
     ```

2. **Join Room**
   - **POST** `/rooms/join`
   - **Body**:
     ```json
     {
       "room_id": "67890",
       "user_id": "12345",
       "display_name": "John Doe"
     }
     ```

3. **Leave Room**
   - **POST** `/rooms/leave`
   - **Body**:
     ```json
     {
       "room_id": "67890",
       "user_id": "12345"
     }
     ```

### Messaging Endpoints
1. **Broadcast Message**
   - **POST** `/messages/broadcast`
   - **Body**:
     ```json
     {
       "room_id": "67890",
       "user_id": "12345",
       "content": "Hello everyone!"
     }
     ```

2. **Send Private Message**
   - **POST** `/messages/private`
   - **Body**:
     ```json
     {
       "sender_id": "12345",
       "receiver_id": "67890",
       "content": "Hello, how are you?"
     }
     ```

3. **Subscribe to Messages (SSE)**
   - **GET** `/messages/subscribe?user_id={userID}`

## Project Structure
```
Chat-Service/
├── internal/
│   ├── core/        # Core business logic
│   ├── models/      # Data models
│   ├── utils/       # Utility functions
├── handlers/        # HTTP handlers for API endpoints
├── main.go          # Entry point of the application
└── README.md        # Project documentation
```

## Usage
1. Start the server using `go run main.go`.
2. Use tools like `Postman` or `curl` to interact with the API.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for improvements or feature additions.

## License
This project is licensed under the MIT License.
