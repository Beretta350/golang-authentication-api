# Golang Authentication API

![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)

A simple lightweight authentication API built in Golang designed to be a fast and easy start for other projects. It includes JWT token generation for secure authentication, bcrypt cryptography for password hashing, and MongoDB for storing user credentials. This API also tries to avoid CSRF and XSS attacks to ensure basic levels of security.

## Layout

The project focuses on enabling a fast and easy start. Thinking on that, nothing better than the [**golang standards project layout**](https://github.com/golang-standards/project-layout). A well known layout among Go users that encourages best practices like modularization and separation of concerns, which can improve code quality and scalability. 

So, if you have questions about the layout check the **link above**!

## Key Technologies

- [**GIN (Web Framework)**](https://github.com/gin-gonic/gin): The application leverages the Gin web framework to handle HTTP requests, routing, and middleware, ensuring fast and scalable API development.

- [**JWT (JSON Web Token)**](https://jwt.io/): JWT token creation, validation, and parsing to provide a secure and efficient way to manage authentication and authorization in the application.

- [**Bcrypt (Cryptography package)**](https://pkg.go.dev/golang.org/x/crypto): Securely hash and compare passwords providing a reliable method for storing and verifying passwords, ensuring that user credentials are protected against unauthorized access.

- [**Docker**](https://www.docker.com/): Docker is used to containerize the application, making it easy to manage dependencies and ensure consistency across different environments.

- [**MongoDB**](https://www.mongodb.com/): NoSQL database with a robust and efficient way to perform CRUD operations, query data, and manage transactions in the application.

## Get Started

To get started with the authentication api, follow these steps:

1. Clone the repository to your local machine.

2. Ensure you have Golang and Docker installed.

3. Build and run the application using Docker Compose. 
    - ```docker compose up -d --build```

4. Access the application's API endpoints from **localhost:8080** to interact with the authetication api.

## Usage

Use the API endpoints to manage your users effectively.

## Example API Endpoints
- `GET /user?id=<userId>`: Get user infos.
- `GET /refreshToken`: Refresh access token using refresh token.
- `POST /login`: Login with user credentials (returns JWT with user ID claim).
- `POST /save`: Create a new user.
- `PUT /update?id=<userId>`: Update user data.
- `DELETE /delete?id=<userId>`: Delete user.

**OBS: Passwords must be at least 8 characters long**

## Contributions

We welcome contributions to this To-Do List application backend. If you're interested in enhancing or extending its functionality, feel free to create pull requests or open issues on the repository.

Enjoy using this flexible To-Do List backend built with Golang and Hexagonal Architecture!
