# README

## Prerequisites

- Docker must be installed on the test environment.
- Golang must be installed on the test environment.

## Setup

To set up the environment, run the following command to start the PostgreSQL server container:

```bash
./run_postgres_server.sh
```

## Installation

To build the source code, execute the following command:

```bash
./build.sh
```

After running this script, an executable file named `transaction-system` will be generated in the current folder.

## Running the Application

To start the transaction system, run the following command:

```bash
./transaction-system
```

The transaction system will now be running on port `8080`. You can use Postman or any other HTTP client tools to test the requests.

## Examples of client requests

### Create account

```bash
curl --location 'localhost:8080/accounts' \
--header 'Content-Type: application/json' \
--data '{
    "account_id": 1,
    "initial_balance": 100
}'
```

### Get account detail

```bash
curl --location 'localhost:8080/accounts/1'
```

### Make transaction

```bash
curl --location 'localhost:8080/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "source_account_id": 1,
    "destination_account_id": 2,
    "amount": 50
}'
```

## Optional

The succesfull transactions are recorded in Postgres table: **transaction_db.transactions**. You can use SQL client to query them.
