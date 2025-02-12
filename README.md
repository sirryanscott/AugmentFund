# AugmentFund

## Overview
AugmentFund is a small golang application to manage shares for a fund (company) between different owners. This application is a service that is interfaced via a REST API. Start the program by running the command below from the root of the project

```go run main.go```

A postman collection can be found here that can be used to interact with the application

```https://www.postman.com/sirryanscott/workspace/augment-fund/collection/5904405-d4dfc752-0bbf-43d4-a359-646b281ebf03?action=share&creator=5904405```

## Assumptions
* When a Fund is created, the total shares remain with the fund initally.
* When creating a transfer, a `fromOwnerId = 0` will indicate that shares are being transferred from the fund to an owner
* When creating a transfer, a `toOwnerId = 0` will indicate that shares are being transferred from the owner back to the fund
* Transferring between owners is not allowed
* Only available shares can be transferred
* If an owner transfers all of their available shares, the owner is removed from the cap table
* If a transfer happens to an existing user that does not have any shares, then the owner is created for the fund
* Users will also keep track of the funds in which they have shares
* A transfer history is kept for each transfer and is returned sorted by date descending

## Testing
A series of unit tests can be found and executed by running this code from the root of the project:

```go test ./...```

## Other Thoughts and Potential Future Enhancements
This project can be further enhanced with these additional features (list not exhaustive):
* proxy server for authentication
* concurrency: this small application wouldn't necessarily benefit from concurrency, but as load increases it would help with effeciency and speed
* logging: logging could be improved for better readibility into various error states
* the database is built around an interface to allow for a better storage solution (other than a file)
* API and input validation can be enhanced