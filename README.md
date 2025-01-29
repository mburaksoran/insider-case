# insider-case

# Job Management System

This project is a system that manages the active/inactive status of jobs based on external `/jobs/start` and `/jobs/stop` requests. When active, the system processes 2 messages every two minutes. In the background, a `fetcher` checks the status of existing jobs and processes them using a worker pool structure. The system is designed to run on multiple pods and uses `SELECT FOR UPDATE SKIP LOCKED` to prevent duplicate processing.

## Technologies Used

- **Go 1.23.5**: The project is written in Go (version 1.23.5).
- **PostgreSQL**: Used to store the status of jobs and messages.
- **Vault**: Used to securely fetch configuration details (e.g., database connection details) from the Docker Compose setup.
- **SQLC**: Used to manage database operations in a type-safe manner.
- **Grafana & Prometheus**: Used for monitoring system performance and visualizing metrics.
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
    - PostgreSQL
    - Prometheus & Grafana (For monitoring metrics)

2. **Configuration**:
    - Configuration details (e.g., database connection, worker count) are fetched from Vault, which is started alongside the project using Docker Compose.

3. **Running the Project**:
    - Start the project and its dependencies using Docker Compose:
      ```bash
      docker-compose up
      ```
    - The application will automatically fetch its configuration from Vault.

4. **API Endpoints**:
    - **Start Job**: `POST /jobs/start`
    - **Stop Job**: `POST /jobs/stop`

## Contributing

If you'd like to contribute to the project, please follow these steps:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push your branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

---

For more information or questions about the project, please open an [issue](https://github.com/yourusername/your-repo/issues).