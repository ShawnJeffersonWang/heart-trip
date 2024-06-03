# HeartTrip - Travel Homestay App

**Collaborative Project (Backend)**

## Overview

HeartTrip is a travel homestay app designed using a microservices architecture. It provides essential features for travel accommodations, including user authentication, online booking, feed streams, likes, follows, and more.

## Technology Stack

- **Backend Framework**: go-zero
- **Database**: MySQL
- **Cache**: Redis
- **Message Queue**: Kafka
- **Search Engine**: Elasticsearch
- **Containerization**: Docker
- **API Gateway**: Nginx

## Key Features and Responsibilities

- **Microservices Architecture**: Constructed a stable backend microservices architecture utilizing go-zero. The development environment is managed with docker-compose, direct connections, and modd hot-reload configurations.
- **Domain-Driven Design**: Applied domain-driven design principles to segment the project into four distinct microservices, ensuring business logic is decoupled and responsibilities are clearly defined.
- **Efficient Data Management**: Leveraged MySQL and Redis for efficient data storage and caching, addressing cache consistency issues to boost system response speed and performance.
- **Distributed Session Management**: Implemented distributed session sharing using Redis, controlled access to shared resources with distributed locks, and utilized GEO for geolocation and queries.
- **Decoupled Messaging**: Employed a factory and strategy pattern for decoupling homestay notifications, and used an adapter pattern for OSS integration.
- **Containerization**: Deployed middleware components using Docker, with data mounted to ensure critical data is isolated and managed effectively.
- **Automated Deployment**: Solved the challenge of multi-machine deployment by using Jenkins and Shell scripts for automated deployment across multiple servers.
- **Environment Setup**: Independently set up the entire project environment and installed dependencies on cloud servers from scratch.

## Getting Started

### Prerequisites

- Docker
- Docker Compose
- go-zero
- MySQL
- Redis
- Kafka
- Elasticsearch
- Jenkins

### Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/ShawnJeffersonWang/HeartTrip.git
    ```

2. **Navigate to the project directory:**

    ```sh
    cd HeartTrip
    ```

3. **Set up the development environment:**

    ```sh
    docker-compose up -d
    ```

4. **Configure and run the application:**

    Ensure all necessary services (MySQL, Redis, Kafka, Elasticsearch) are configured and running.

## Contributing

We welcome contributions from the community. Please fork the repository and submit pull requests.

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details.

## Contact

For any inquiries or feedback, please open an issue on the GitHub repository.
