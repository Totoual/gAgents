# Email Task Scheduler with gAgents

This is an example use case illustrating how to integrate the `gAgents` library for scheduling email tasks. The application consists of three primary files: `email.go`, `service.go`, and `main.go`.

## Description

- **email.go**: Defines the `EmailTask` type and its methods. An `EmailTask` represents an email that can be scheduled for sending at a specific time. The task fetches the sender's Gmail authentication code from the environment and utilizes the `gomail` package to send emails.
  
- **service.go**: Contains the gRPC server implementation for the `EmailService`. It exposes a method `CreateEmailTask` to accept email details and schedule them for delivery.
  
- **main.go**: Initializes and runs the main agent, registers the email service, and starts listening for incoming requests.

## Setup

1. **Environment Variables**: Ensure you have set up the `.env` file with your Gmail authentication code under the key `gmail_code`.

2. **Running the service**: Use the command `go run main.go` to initialize and start the agent.

## Usage

1. Use a gRPC client or any suitable mechanism to call the `CreateEmailTask` method from the `EmailService`. Provide the required email details, such as sender, receiver, subject, body, and scheduled time.

2. The service will schedule the email for sending based on the provided scheduled time.

## Notes

- The email task is scheduled only once, and after sending the email, the task stops itself.
- Make sure to handle the `gmail_code` securely and avoid exposing it in public repositories.

---
