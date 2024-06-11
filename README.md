
### Description

This project implements an HTTP API server designed to aggregate social media posts data from the Upfluence SSE API. The server listens on port 8080 and accepts only HTTP GET requests on the /analysis endpoint. The API provides statistical analysis based on the data received from the SSE stream.

Solution Overview

Architecture:

Although it could have been written in a simpler way, for future expandability reasons I went with a more modular design.

- The sseclient package is responsible for interacting with the SSE stream.
- The aggregator package handles data processing and analysis.
- The handler package provides HTTP endpoint logic.
- The server package sets up and starts the HTTP server.
- The types package defines reusable types and validation logic.
- The testingTools package provides variables and functions to facilitate testing

For easier testing, interfaces are used to facilitate testing with mocks.

Concurrency: Utilizes Go routines and channels for concurrent data processing and HTTP request handling.
### Running the project

Dependencies: Go 1.22 or later

Installation

1. Clone the repository: 
    ```
    git clone https://github.com/AlexisMontagne/junior-backend-engineer-challenge.git
    ```

2. Install the dependencies:
    ```
    cd upfcc
    go mod tidy
    ```

3. Running the server
    ```
    cd cmd/server
    go run main.go
    ```

### Example Usage
To analyze posts for a duration of 30 seconds based on the number of likes:

    curl "localhost:8080/analysis?duration=30s&dimension=likes"


### Trade-offs and Considerations

In terms of testing, due to time constraints, the project currently lacks extensive test coverage. However, there are several areas where additional tests could be beneficial. For example, it would be valuable to verify that requests are handled correctly for the specified duration. Additionally, testing the ability to handle multiple requests concurrently would be beneficial. End-to-end or integration tests could also be added to ensure the overall functionality of the system. Regrouping of mock implementations could be done also.

If this project were to be deployed in a production environment, there are several considerations and adjustments that would need to be made. Firstly, the absence of authentication, authorization, and encryption features should be addressed to ensure the security of the system. Deployment environment questions, such as the target platform and infrastructure, should also be taken into account. For scalability, multiple instances of the server could be deployed, and a load balancer could be used to distribute incoming requests across these instances. To improve performance and reliability, message queues could be implemented to handle requests more efficiently. 

To streamline the development process, implementing a CI/CD pipeline using tools like GitLab CI/CD or GitHub Actions would be beneficial. This would automate testing and deployment processes, ensuring that changes are thoroughly tested and deployed consistently. Additionally, a more robust logging system could be implemented to provide better visibility into the server's operations. Lastly, adding API versioning support would allow for future changes without breaking existing clients, ensuring backward compatibility.
