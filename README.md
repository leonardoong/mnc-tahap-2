# E-Wallet REST API

A simple e-wallet system with background transaction processing using Redis enqueue. Built with Golang and Gin Framework.

## Features
- User authentication (Sign Up, Login, JWT-based Auth)
- Top Up, Payment, and Transfer functionality
- Background transaction processing using Redis
- Transaction history

## Tech Stack
- Backend: Golang, Gin Framework
- Database: MySQL 8.0
- Task Queue: Redis
- Authentication: JWT (JSON Web Token)
- Containerization: Docker & Docker Compose

## Installation
### Prerequisites:
- Golang 1.20+
- Docker & Docker Compose
- MySQL 8.0
- Redis 7.0

### Steps

1. **Clone the Repository:**
    ```sh
    git clone https://github.com/leonardoong/simple-e-wallet.git
    cd simple-e-wallet
    ```
2. **Run the Application Using Docker Compose:**
    ```sh
    docker-compose up --build -d
    ```
3. **Check Running Containers:**
   ```sh
    docker ps
    ```
4. **Access the Application: The API server will be running at ```http://localhost:8080```**

## API Endpoints

| Method | Endpoint                 | Description              | Need Auth  |
|--------|--------------------------|--------------------------|------------|
| POST   | `/register`              | Register a new user      | No         |
| POST   | `/login`                 | Login and get token      | No         |
| PUT    | `/profile    `           | Update user profile      | Yes        |
| POST   | `/topup`                 | Top Up money             | Yes        |
| POST   | `/payment`               | Payment                  | Yes        |
| POST   | `/transfer`              | Transfer funds           | Yes        |
| GET    | `/topup/:topup_id`       | Top Up money             | Yes        |
| GET    | `/payment/:payment_id`   | Payment                  | Yes        |
| GET    | `/transfer/:transfer_id` | Transfer funds           | Yes        |
| GET    | `/transactions`          | Transaction history      | Yes        |

Also you can check in the postman collection.
