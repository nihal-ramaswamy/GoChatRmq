# Go Chat Application
Go chat application using Redis and RabbitMQ.

## Table of Contents
- [Installation](#installation)
- [Setup](#setup)
- [Running](#running)
- [Internal Working](#internal-working)
- [Testing](#testing)

## Installation 
Ensure you have Docker installed on your machine.
This service can run without docker but it is recommended to use docker for easy setup.
The below instructions will guide you on how to run the service using docker.

## Setup 
Copy the `.env.example` file to `.env` and fill in the required values.

## Running 
Start the server using the `docker-compose.yaml` command file.

## Internal Working 
- Once the user is logged in, they can send messages by hitting the HTTP endpoint documented inside [API](./internal/api/chat/) folder.
- This saves the message in the database and also sends the message to the RabbitMQ queue with ID as the receiver's ID.
- The receiver can either read messages from the DB or listen to the RabbitMQ queue for new messages.

## Testing 
Tests are written using [testContainer](https://testcontainers.com). Ensure docker is running before running the tests.
### Running Tests
First make the `.env.test` file based on the instructions inside the `.env.example` file. Then run the below command.
```bash
ENVIRONMENT=test go test ./... -v -cover
```

