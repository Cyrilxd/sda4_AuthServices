
# SDA4 Authentication Services Group B

Welcome to the SDA4 Authentication Services! This repository contains a Go-based authentication service. Follow the steps below to set up and run the application.

---

## Prerequisites

### 1. Install Go
- Follow the official Go installation guide: [Install Go](https://go.dev/doc/install)
- Verify your Go environment by checking the following:
  ```bash
  $ echo $GOPATH
  $ echo $GOROOT
  ```
    - **`GOROOT`**: Points to the directory where Go is installed.
    - **`GOPATH`**: Your working directory for Go projects.

### 2. Install Docker
- Follow the official Docker installation guide: [Install Docker Desktop](https://docs.docker.com/desktop/)

> **Note:** After completing these steps, restart your computer to ensure all environment variables are correctly set up.

---

## Run the Application

1. Clone the repository and navigate to the project directory:
   ```bash
   cd sda4_Authservices
   ```

2. Start the application using Docker Compose:
   ```bash
   docker compose up
   ```

---

## Testing the Application

### 1. Import Request Collection
- Import the file `TEST-DATA.yaml` into your **Insomnia** application (or any API testing tool of your choice).

### 2. Make API Requests
- Use the imported collection to test the application's endpoints.

---

## Notes
- Ensure Docker is running before starting the application.
- For any issues, consult the official Go and Docker documentation linked above.

Happy Testing! 🚀
