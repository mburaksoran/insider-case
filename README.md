# Job Management System

This project is a system that manages the active/inactive status of jobs based on external `/jobs/start` and `/jobs/stop` requests. When active, the system processes 2 messages every two minutes. In the background, a `fetcher`  checks the status of existing jobs and processes them using a worker pool structure. The system is designed to run on multiple pods and uses `SELECT FOR UPDATE SKIP LOCKED` to prevent duplicate processing.

## Technologies Used

- **Go 1.23.5**: The project is written in Go (version 1.23.5).
- **PostgreSQL**: Used to store the status of jobs and messages.
- **Vault**: Used for as parameter store.
- **SQLC**: Used to manage database operations in a type-safe manner.
- **Grafana & Prometheus**: Used for monitoring system performance and visualizing metrics.
- **Loki**: Used for monitoring logs.
- **Worker Pool**: Used to manage parallel processing and handle tasks.
- **Go Channels**: Used for communication between goroutines, such as job and task channels.

## How the System Works

1. **Job Control**:
   - A background job called `fetcher` checks the jobs in PostgreSQL every 10 seconds.
   - If a job's `interval` value + last trigger time is less than the current time (`time.Now`), the job is sent to a job channel.

2. **Job Processing**:
   - A Go routine listens to the job channel and converts incoming jobs into tasks.
   - Tasks are then sent to a task channel for processing by the worker pool.
   - The worker pool has a configurable number of workers (fetched from Vault). Each worker processes tasks by creating the appropriate handler via a `handler factory`.

3. **Multi-Pod Support**:
   - To support running on multiple pods, jobs are fetched using `SELECT FOR UPDATE SKIP LOCKED`, and their status is updated from `active` to `in_progress`.
   - Similarly, messages to be processed are ordered by `created ASC` and fetched using `SKIP LOCKED`. This ensures that different pods do not process the same jobs or messages, preventing duplicate work.

## Installation and Running

1. **Requirements**:
   - Go 1.23.5 or higher
   - Docker and Docker Compose (for running Vault and other services)

2. **Configuration**:
   - Configuration details (e.g., database connection, worker count) are fetched from Vault, which is started alongside the project using Docker Compose.

3. **Running the Project**:
   - Start the project and its dependencies using Docker Compose:

    ```bash
    docker-compose up
    ```

   - The application will automatically fetch its configuration from Vault.
   - You can access the api with 
   ```bash
   http://localhost:8080/swagger/index.html
    ```

4. **API Endpoints**:
   - **Start Job**: `PUT /jobs/start`
   - **Stop Job**: `PUT /jobs/stop`
   - **Listing Messages sent**: `GET /messages/sent`
   - **HealthCheck**: `GET /health-check`
   - **HealthAlive**: `GET /health-alive`

---

## Redis UI Access

The project includes a Redis UI to visualize the processed data. You can access the Redis UI with the following credentials:

- **Username**: `root`
- **Password**: `qwerty`
- **Access URL**: `http://localhost:6379`

You can view and inspect the processed data through this interface.

---

## Grafana & Metrics Monitoring

Grafana is used to visualize system metrics. You can access the Grafana UI with the following credentials:

- **Access URL**: `http://localhost:3000`
- **Username**: `admin`
- **Password**: `admin`

Within Grafana, you can monitor the system's performance metrics collected via Prometheus.

---

## Loki & Log Monitoring

Loki is used for collecting and visualizing system logs. You can access the Loki UI for log inspection with the following:

- **Access URL**: `http://localhost:3100` (default port)

---

## Vault UI Access

You can access the Vault UI for configuration details and secrets management with the following credentials:

- **Access URL**: `http://localhost:8200`
- **Token**: `root`

Once logged in, you can view and manage the secrets and configurations used by the system.

---
