# Web & Panel & API Project

This project consists of three main components:

- A web application (frontend)
- An API service (backend)
- An Angular panel application (admin dashboard)

## Project Structure

```
/api                 - Backend API service
    /asset          - API assets including Swagger documentation
    /cmd            - Application entry point
    /handler        - HTTP request handlers
    /internal       - Internal packages (auth, config, connections, etc.)
    /model          - Data models and entities
    /pkg            - Reusable packages
    /service        - Business logic and application services
    /view           - API documentation views

/web                - Frontend Web application
    /asset          - Frontend assets (CSS, JS, images)
    /cmd            - Web application entry point
    /handler        - HTTP request handlers
    /model          - Data models
    /service        - Service packages
    /view           - Frontend templates and components

/panel              - Angular Admin Dashboard
    /src            - Source files
        /app        - Application modules and components
            /guards      - Authentication and permission guards
            /interfaces - TypeScript interfaces
            /layout     - Layout components
            /modules    - Feature modules
            /services   - Application services
            /shared    - Shared modules and components
```

## Features

### API Service

- Modular and scalable architecture
- JWT Authentication
- Database integrations (PostgreSQL, Redis, Kafka)
- Swagger API documentation
- Input validation
- Comprehensive error handling
- Logging system
- Email functionality
- SQL query builder
- Caching mechanisms

### Web Application

- Modern frontend architecture
- Localization support
- Template rendering
- GraphQL integration
- CDN support
- Excel file handling
- Component-based structure
- Responsive design

### Angular Panel

- Modern Angular-based admin dashboard
- Component-based architecture
- Authentication and authorization
- Role-based access control
- Responsive layout system
- Service integration with API
- Interceptors for request/response handling
- TypeScript interfaces for type safety
- Shared modules and components
- Guard-protected routes

## Getting Started

1. **Clone the repository**:

   ```bash
   git clone https://github.com/mstgnz/starter-kit.git
   ```

2. **Setup API Service**:

   ```bash
   cd api
   cp .env.example .env    # Configure your environment variables
   go mod tidy            # Install dependencies
   make live              # Run the API service
   ```

3. **Setup Web Application**:

   ```bash
   cd web
   cp .env.example .env    # Configure your environment variables
   go mod tidy            # Install dependencies
   make live              # Run the web application
   ```

4. **Setup Angular Panel**:
   ```bash
   cd panel
   npm install              # Install dependencies
   npm start               # Run the development server
   ```

## Configuration

Both the API and Web components use `.env` files for configuration. Example files (`.env.example`) are provided in each directory. Make sure to configure these properly before running the applications.

## Development

- The API service runs on port specified in the API's `.env` file
- The Web application runs on port specified in the Web's `.env` file
- API documentation is available at `/view/swagger.html` in the API service

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
