
# Receipt Processor Challenge

This project is a Go-based service to process receipts, assign a unique identifier to each receipt, and calculate points based on predefined criteria. The service exposes two main endpoints:
- `POST /receipts/process`: Accepts receipt data and returns a unique receipt ID.
- `GET /receipts/{id}/points`: Retrieves the points associated with the receipt ID.

## Table of Contents

- [Installation](#installation)
- [Running the Service](#running-the-service)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Sample Requests](#sample-requests)

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or higher)
- Git (optional, for cloning the repository)

### Steps

1. **Clone the Repository**

   Clone the repository using Git or download the source code:

   ```bash
   git clone https://github.com/fetch-rewards/receipt-processor-challenge.git
   cd receipt-processor-challenge
   ```

2. **Install Dependencies**

   No external dependencies are required for this project, as it uses Go’s standard library.

3. **Build the Project**

   Build the project by running:

   ```bash
   go build -o receipt-processor
   ```

This command will generate an executable called `receipt-processor` in the project directory.

## Running the Service

To start the service, run the following command in the project directory:

```bash
./receipt-processor
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### POST /receipts/process

Accepts a JSON payload of receipt data, generates a unique receipt ID, and returns the ID.

- **URL**: `/receipts/process`
- **Method**: `POST`
- **Content-Type**: `application/json`

**Request Payload** (Example):

```json
{
  "retailer": "Target",
  "purchaseDate": "2024-10-25",
  "purchaseTime": "13:13",
  "items": [
    { "shortDescription": "Pepsi", "price": 1.25 },
    { "shortDescription": "Bread", "price": 2.50 }
  ],
  "total": 3.75
}
```

**Response**:

```json
{
  "id": "12345"
}
```

### GET /receipts/{id}/points

Retrieves the points for the given receipt ID.

- **URL**: `/receipts/{id}/points`
- **Method**: `GET`
- **Path Parameter**: `id` (the receipt ID)

**Response**:

```json
{
  "points": 150
}
```

## Testing

You can test the API using `curl` commands or tools like [Postman](https://www.postman.com/) or [Insomnia](https://insomnia.rest/).

### Example Test Cases

1. **Submit a Receipt**:

   ```bash
   curl -X POST http://localhost:8080/receipts/process    -H "Content-Type: application/json"    -d '{
         "retailer": "Target",
         "purchaseDate": "2024-10-25",
         "purchaseTime": "13:13",
         "items": [
           { "shortDescription": "Pepsi", "price": 1.25 },
           { "shortDescription": "Bread", "price": 2.50 }
         ],
         "total": 3.75
       }'
   ```

   **Expected Response**:

   ```json
   {
     "id": "12345"
   }
   ```

2. **Retrieve Points for a Receipt**:

   Assuming the receipt ID is `12345`:

   ```bash
   curl -X GET http://localhost:8080/receipts/12345/points
   ```

   **Expected Response**:

   ```json
   {
     "points": 150
   }
   ```

## License

This project is for educational and demonstration purposes.
