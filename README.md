# Deployed site
You can access this API hosted on Render at: https://onecvtest.onrender.com/

Please note that as I have deployed the API using the free web service version on Render, the API may sometimes load slowly on the initial request.

# Steps to running API server locally

## 1. Clone the Repository

Into a new folder,

HTTPS: Clone with `git clone https://github.com/garylow2001/OneCV_Tech_Test.git`

SSH: Clone with `git clone git@github.com:garylow2001/OneCV_Tech_Test.git`

## 2. Install Go:
Ensure that you have Go installed on your system. You can download and install it from the official Go website: https://go.dev/dl/.

## 3. Install Dependencies:
Navigate to the directory where you cloned the repository.
Install the project dependencies using the Go module system:

`go mod tidy`

## 4. Run the server:
Ensure that you are on the root of the api directory (where main.go exists), then run:

`go run main.go`

The api should now be running and listening for incoming requests on the port 3000.

Disclaimer about the database: Usually I would store the port and database url in a .env as a secret so that it is not publicly accessible. However for the purpose of submission I have left it in the code.

## 5. Testing:
To run the unit tests, first cd into tests

`cd tests`

Next, run

`go test`

It should show "PASS" and "ok" at the end

To run individual test files, make sure to include `helpers.go` file, for example:

`go test register_test.go helpers.go`