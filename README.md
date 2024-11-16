# Chess Engine Microservice Architecture

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?logo=docker&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-Enabled-blue)

## **Overview**

The **Chess Engine Microservice** is designed to handle chess engine
computations in a highly concurrent and efficient manner. It uses a
**worker pool architecture** to manage multiple chess engine instances
and processes chess move calculations. The service integrates with the
frontend and backend systems, providing a scalable and modular way to
execute chess-related computations.

---

## **Role of the Microservice**

The Chess Engine Microservice is responsible for:

1. **Chess Move Computation**:

   - Receives requests with chess board positions (in FEN format).
   - Computes the best move using chess engine backends (e.g., Stockfish, Argo).

2. **Worker Pool Management**:

   - Maintains a pool of chess engine instances for concurrent processing.
   - Routes requests to the appropriate engine type.

3. **Efficient Resource Usage**:

   - Ensures that resources (chess engine instances) are reused and not overwhelmed.
   - Manages request queuing and results delivery efficiently.

4. **Communication Middleware**:
   - Exposes APIs over gRPC for backend communication.
   - Relays results to the client frontend via a WebSocket.

---

## **Architecture**

### **1. High Level Architecture**

- **Frontend - Webscoket Client (React)**:

  - Acts as the chess game client.
  - Communicates with the backend server over WebSocket to send user moves and
    receive engine responses.

- **Backend Server - Websocket Server/gRPC Client (Go)**:

  - Acts as the gateway between the frontend and the chess engine microservice.
  - Handles WebSocket connections with the frontend.
  - Forwards chess move requests to the microservice via gRPC.

- **Chess Engine Microservice - gRPC Server (Go)**:
  - The core computational service.
  - Processes requests via worker pools for specific engine types
    (e.g., Stockfish, Argo (engine developed as a side project by me)).
  - Sends results back to the backend server.

---

### **2. Chess Engine Microservice Workflow**

1. **Request Flow**:

   - The backend server sends a chess move request (with FEN and engine type)
     to the microservice via gRPC.

2. **Worker Pool**:

   - The microservice routes the request to the appropriate worker pool
     based on the engine type.
   - Each worker pool maintains a set of chess engine instances to handle
     requests concurrently.

3. **Result Processing**:

   - Each worker computes the best move for the given FEN position.
   - Results are stored in a thread-safe sync.Map and returned to the backend server.

4. **Response Flow**:
   - The backend server forwards the result to the frontend via WebSocket.

---

## **Detailed Architecture Diagram**

The architecture is represented in the following components and interactions:

![Screenshot](./images/microservice-diagram.png)

---

## **Advantages**

1. **Scalability**:

   - The microservice architecture supports horizontal scaling, allowing additional
     workers or engine types to be added easily.

2. **Concurrency**:

   - Efficiently handles concurrent requests using goroutines and worker pools.

3. **Modularity**:

   - Decouples the frontend and backend computation logic,
     simplifying maintenance and upgrades.

4. **Performance**:
   - Uses gRPC and worker pools to minimize latency and maximize throughput.

---

## Setup and Installation

### Prerequisites

- **Go**: The microservice is written in Go, so you'll need Go installed.
- **gRPC**: Ensure that gRPC is installed and properly set up.
- **Docker** (optional): To containerize the microservice for deployment.

### Installation Steps

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/chess-engine-microservice.git
   cd chess-engine-microservice
   ```

2. Install dependencies:

   ```sh
   go mod download
   ```

3. Build the microservice:

   ```sh
   go build -o chess_engine_service
   ```

4. Run the microservice:

   ```sh
   ./chess_engine_service
   ```

## Configuration

The microservice can be configured through environment variables:

- **`WORKER_POOL_SIZE`**: The number of workers in each engine pool.
- **`ENGINE_TYPE`**: Specify the type of engines to use (e.g., "stockfish").
- **`GRPC_PORT`**: The port on which the gRPC server will listen.

## Usage

- **Submit a Move Request**: Send a gRPC request to the microservice with the FEN
  string and engine type.
- **Receive Computed Move**: The microservice will compute the best move and send
  it back as a gRPC response.

## Example gRPC Request

```proto
todo: Add example protobuf request and response here
```

## **Future Enhancements**

1. **Dynamic Pool Management**:

   - Automatically scale worker pools based on request load.

2. **Caching**:

   - Implement a caching layer for frequently requested positions.

3. **Monitoring**:

   - Add telemetry and logging for real-time performance tracking.

4. **Additional Engines**:
   - Extend support for other chess engines or computation models.

---

## **Contributing**

Feel free to fork the repository and submit pull requests. All contributions are
welcome to improve the functionality, scalability, and usability of the microservice.

## **Conclusion**

The Chess Engine Microservice provides a robust, scalable, and efficient architecture
for managing chess move computations. Its modular design and focus on concurrency
make it well-suited for real-time chess applications.
