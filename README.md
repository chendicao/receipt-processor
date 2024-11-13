
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
   git clone https://github.com/chendicao/receipt-processor.git
   cd receipt-processor
   ```

2. **Initialize Go Modules**

    Ensure that the Go modules are initialized (this should already be done if you're using Go 1.16 or higher). If the `go.mod` file is not already present, you can initialize it with:

   ```bash
   go mod init github.com/chendicao/receipt-processor
   ```
3. **Install Dependencies**
   The project uses the following dependencies:
   - **Gorilla Mux**: A powerful URL router and dispatcher for Go.
   - **UUID**: A package to generate UUIDs.
   - **Go Playground Validator**: For data validation.

   The `go.mod` file will automatically track these dependencies.
4. **Verify Dependencies Installation**
   ```bash
   go run main.go
   ```

3. **Build the Project**

   Build the project by running:

   ```bash
   go build -o receipt-processor
   ```

   This command will generate an executable called `receipt-processor` in the project directory.

   Alternatively, you can run the service directly without building the executable:

   ```bash
   go run main.go
   ```


## Running the Service

To start the service, run the following command in the project directory:

```bash
./receipt-processor
```

The server will start on `http://localhost:8000` by default.

## API Endpoints

### POST /receipts/process

Accepts a JSON payload of receipt data, generates a unique receipt ID, and returns the ID.

- **URL**: `/receipts/process`
- **Method**: `POST`
- **Content-Type**: `application/json`

Receipt ID Generation: A unique identifier (UUID) is generated for each receipt using the code `receiptID := uuid.New().String()`. This ensures that each receipt has a distinct ID.

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
  "id": "cc2cb204-eacb-4689-a58e-6f4f41945299"
}
```

**Invalid Receipt Example**:

If the receipt data is invalid (e.g., missing required fields or incorrect price format), you will receive an error.

**Request Payload (Invalid)**:

```json
{
  "retailer": "Target",
  "purchaseDate": "2024-10-25",
  "purchaseTime": "13:13",
  "items": [
    { "shortDescription": "Pepsi", "price": "1.2" },  // Invalid price format
    { "shortDescription": "Bread", "price": "invalid" }  // Invalid price
  ],
  "total": 3.75
}
```

**Expected Response (400 Bad Request)**:

```json
{
  "error": "Invalid receipt"
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

**Receipt Not Found Example**:

If the receipt ID does not exist or is invalid, you will receive a `404 Not Found` error.

**Request (Non-existent Receipt ID)**:

```bash
curl -X GET http://localhost:8000/receipts/invalid-receipt-id/points
```

**Expected Response (404 Not Found)**:

```json
{
  "error": "No receipt found for that id"
}
```

## Schema Validation

The service uses **go-playground/validator** to validate the structure of the receipt data.

- **Required Fields**: All required fields must be provided in the receipt JSON (e.g., `retailer`, `purchaseDate`, `items`).
- **Field Formats**:
  - `purchaseDate`: Must follow the `YYYY-MM-DD` format.
  - `purchaseTime`: Must follow the `HH:MM` format.
  - `price`: Must be a valid floating-point number (e.g., `1.25`).
  - `total`: Must be a valid floating-point number (e.g., `3.75`).
  
If any of the fields are missing or incorrectly formatted, the service will return a `400 Bad Request` with a detailed error message indicating which field caused the validation failure.

### Example Test Cases

1. **Submit a Valid Receipt**:

   ```bash
   curl -X POST http://localhost:8000/receipts/process    -H "Content-Type: application/json"    -d '{
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
        "id": "cc2cb204-eacb-4689-a58e-6f4f41945299"
      }
      ```
 

2. **Retrieve Points for a Receipt**:

   Assuming the receipt ID is `cc2cb204-eacb-4689-a58e-6f4f41945299`:

   ```bash
   curl -X GET http://localhost:8000/receipts/cc2cb204-eacb-4689-a58e-6f4f41945299/points
   ```

   **Expected Response**:

   ```json
   {
     "points": 47
   }
   ```

3. **Invalid Receipt Example (400 Bad Request)**:

   ```bash
   curl -X POST http://localhost:8000/receipts/process    -H "Content-Type: application/json"    -d '{
         "retailer": "Target",
         "purchaseDate": "2024-10-25",
         "purchaseTime": "13:13",
         "items": [
           { "shortDescription": "Pepsi", "price": "1.2" },
           { "shortDescription": "Bread", "price": "invalid" }
         ],
         "total": 3.75
       }'
   ```

   **Expected Response**:

   ```json
    {
      "error": "Invalid receipt structure: Key: 'Item.Price' Error:Field validation for 'Price' failed on the 'required' tag"
    }
   ```

4. **Receipt Not Found Example (404 Not Found)**:

   ```bash
   curl -X GET http://localhost:8000/receipts/invalid-receipt-id/points
   ```

   **Expected Response**:

   ```json
   {
     "error": "No receipt found for that id"
   }
   ```

