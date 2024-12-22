# Go Project Template

This repository provides a well-structured and modular foundation for building Go applications. It's designed to help you quickly get started with a new Go project by providing a predefined directory layout, basic configurations, and sample code.

## Features

- **Modular Structure**: A clean and organized directory layout that scales with your project.
- **Ready-to-Use Configuration**: Includes configurations for common tools and frameworks.
- **Sample Code**: Provides basic examples for common tasks like authentication, database connections, and validation.
- **Best Practices**: Follows Go best practices in terms of project structure and code organization.

## Directory Layout
This repository is a starter kit designed to provide a foundational structure for your project. Feel free to modify and adapt the directory layout and files according to your specific needs and project requirements. The provided structure is intended to serve as a basis, and adjustments can be made to better fit your use case.


```
/asset
    /css    - Contains CSS files for styling the frontend of the application.
    /img    - Stores image assets used in the application.
    /js     - Holds JavaScript files for adding interactivity to the frontend.
    /lang   - Localization
    swagger.yaml - API documentation written in the OpenAPI Specification (formerly Swagger).
    Note: This directory is typically used for frontend projects. If your project is an API-only project, this directory can be removed.
/cmd
    main.go - The entry point for the application, where the main function resides.
/handler
    - Handles HTTP requests, mapping them to corresponding services or business logic.
/internal
    /auth
        bcrypt.go   - Handles password hashing and verification using bcrypt.
        jwt.go      - Manages JWT (JSON Web Token) creation, parsing, and validation.
    /config
        catch.go    - Central error handling
        conf.go     - Centralizes and manages application-wide dependencies and services.
        routes.go   - Generate route list with localization
    /conn
        db.go       - Manages database connections and queries.
        kafka.go    - Handles Kafka messaging system connections and operations.
        redis.go    - Manages Redis connections and operations.
    /load
        excel.go    - Exelize package
        render.go   - Render templ template 
        sql.go      - Loads and manages SQL queries or scripts.
    /localization
        localization.go - Localization for use in template
    /logger
        logger.go   - Provides logging functionalities for the application.
    /response
        json.go     - Formats and handles JSON responses.
        response.go - Defines structures and functions for handling HTTP responses.
    /validate
        validate.go - Implements input validation using libraries like go-playground/validator.
/model
    - Contains the data models and entities that represent the structure of the data used in the application.
/pkg
    /mstgnz
        cache.go        - Implements caching mechanisms, potentially using technologies like Redis.
        gobuilder.go    - With SQL builder  you can write your sql queries in a readable way.
        mail.go         - Manages email sending functionalities.
/service
    - Contains business logic and application services that orchestrate operations between handlers, models, and other packages.
/view
    /component  - Stores reusable frontend components (HTML, templates, etc.).
    /page       - Holds the page-specific templates or views for the application.
    swagger.html - The HTML file that displays the Swagger UI for API documentation.
    Note: This directory is typically used for frontend projects or projects with user interfaces. If your project is an API-only project, this directory can be removed.
.dockerignore - Lists files and directories to ignore in the Docker context during image build.
.env            - Contains environment variables for application configuration.
.env.example    - An example .env file, showing required environment variables without sensitive data.
.gitignore      - Specifies files and directories to be ignored by Git version control.
dockerfile      - Script with instructions to build a Docker image for the application.
go.mod          - Defines the module path and lists the dependencies of the Go project.
go.sum          - Records the checksums of the dependencies listed in go.mod.
LICENSE         - The license under which the project is distributed.
makefile        - Contains rules to automate tasks such as building, testing, and running the application.
query.sql       - SQL query files that can be used for raw queries.
README.md       - The main documentation file that provides an overview of the project and instructions for setup and usage.
```

## Getting Started

1. **Clone the repository**:
    ```bash
    git clone https://github.com/mstgnz/starter-kit.git
    ```
2. **Rename the directory** to your project name:
    ```bash
    mv starter-kit my-new-project
    ```
3. **Initialize a new Go module**:
    ```bash
    cd my-new-project
    go mod init github.com/mstgnz/my-new-project
    ```
4. **Install dependencies**:
    ```bash
    go mod tidy
    ```
5. **Start developing**: Customize the template to suit your project's needs.
6. **Run project**:
    ```bash
    make live
    ```

## How to Use

- **Adding a new package**: Create a new directory under `/pkg` and place your code there. Follow the modular approach for clear organization.
- **Configuration**: Adjust the settings in the `/internal/config` package to fit your environment (e.g., database credentials, API keys).

## Contribution

If you have any suggestions or improvements, feel free to open an issue or submit a pull request. Contributions are welcome!

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.