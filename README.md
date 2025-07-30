# Pharmacy Claim Service

This project implements a pharmacy claim management service, enabling the submission and reversal of claims. It includes functionalities for data persistence, API authentication, and observability with Prometheus and Grafana.

---

## 🚀 Getting Started

### Prerequisites

Make sure you have the following software installed on your machine:

* **Go** (version 1.22 or higher)
* **Docker** and **Docker Compose**

### ⚙️ Running Locally (Without Docker)

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/diogocarasco/go-pharmacy-service.git](https://github.com/diogocarasco/go-pharmacy-service.git)
    cd go-pharmacy-service
    ```
2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Create the `.env` file:**
    Create a `.env` file in the project root with the following environment variables:
    Replace `hippotoken` with a token of your choice.
  
4.  **Run the service:**
    ```bash
    go run main.go
    ```
    The service will be available at `http://localhost:8080`.

---

## ✅ Running Tests Locally

To run the project's unit tests, execute the following command in the project root:

```bash
go test ./...
```

--- 

## 🐳 Running the Application with Docker Compose

The recommended way to run the application and its observability tools is by using Docker Compose.

1.  **Create the .env file:**
As explained in the "Running Locally" section, create a .env file in the project root with the necessary environment variables. Ensure AUTH_TOKEN and PORT are configured.

2.  **Build and bring up the services:**
    From the project root, run:
    ```bash
    docker-compose up --build
    ```
    This command will build the Docker images and start the containers for the pharmacy service, Prometheus, and Grafana.

---

## 📊 Accessing Observability Tools

With Docker Compose running, you can access the following interfaces:

**Swagger UI (API Documentation)**
* **URL:** `http://localhost:8080/swagger/index.html`
* Access this URL in your browser to view the interactive API documentation, which allows you to test endpoints.

**Prometheus**
* **URL:** `http://localhost:9090`
* Prometheus is the monitoring and alerting tool. You can explore the metrics exported by the service on this dashboard.

**Grafana**
* **URL:** `http://localhost:3000`
* **Default User:** `admin`
* **Default Password:** `admin`
* In Grafana, you can create dashboards to visualize the metrics collected by Prometheus. Make sure to configure Prometheus as a data source in Grafana if it's not already automatically configured (it usually is pre-configured via Docker Compose).

---

## `data` Directory Structure
The `data` directory is crucial for the service, as it stores persistent data and input files.
* `data/`
    * `pharmacy.db`: This is the SQLite database file where all pharmacy, claim, and revert information is stored.
    * `pharmacies.csv`: A CSV file containing the initial list of pharmacies that will be loaded into the database upon service startup.
    * `claim/`: A directory where claim files (e.g., in JSON or CSV format, depending on your internal loader implementation) can be placed to be loaded into the database.
    * `reversal/`: A directory where revert files (similar to claims, in JSON or CSV format) can be placed to be loaded into the database.

--- 

## How to Make an API Call (Example)
You can interact with the API using tools like `curl` or Postman/Insomnia. Remember to replace `hippotoken` with the token configured in your `.env`.

**Example: Submit a New Claim**
**Endpoint:** `POST /claim`
**Headers:**
* `Content-Type: application/json`
* `Authorization: Bearer hippotoken`

**Request Body (JSON):**
```json
{
    "ndc": "00002323401",
    "quantity": 5.5,
    "npi": "1234567890",
    "price": 75.25
}
```

**Example with** `curl`
```bash
curl -X POST \
  http://localhost:8080/claim \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer hippotoken' \
  -d '{
    "ndc": "00002323401",
    "quantity": 5.5,
    "npi": "1234567890",
    "price": 75.25
  }'
```


**Example: Reverse an Existing Claim**
**Endpoint:** `POST /reversal`
**Headers:**
* `Content-Type: application/json`
* `Authorization: Bearer hippotoken`

**Request Body (JSON):**
```json
{
    "claim_id": "09c8533e-27bc-4370-ad76-c2d656390782"
}
```

**Example with** `curl`
```bash
curl -X POST \
  http://localhost:8080/reversal \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer hippotoken' \
  -d '{
        "claim_id": "09c8533e-27bc-4370-ad76-c2d656390782"
  }'
```