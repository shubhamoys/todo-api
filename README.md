# Todo API

A RESTful API for task management built with Go. This project serves as a template for building scalable Go applications with MongoDB.

## Features

- 🔐 JWT Authentication
- 👤 Role-based access control
- 📝 CRUD operations for tasks
- 🔍 Advanced filtering and searching
- 📄 Pagination support
- 🔗 Field population (similar to JOIN operations)
- 📊 Sorting and field selection

## Prerequisites

- Go 1.20 or higher
- MongoDB 6.0 or higher
- Make (optional, for using Makefile commands)

## Environment Setup

1. Clone the repository:

```bash
git clone https://github.com/shubhamoys/todo-api
cd todo-api
```

2. Copy paste the .env.example file, rename it to .env and make the necessary changes.

3. Install dependencies:

```bash
go mod download
```

## Running the Application

### Development

```bash
go run cmd/main.go
```

### Seed Super Admin data

```bash
go run cmd/seed_superadmin/seed_superadmin.go
```

### Using Make (if available)

```bash
# Run the application
make run

# Build the application
make build

# Run tests
make test
```

## API Documentation

### Authentication

- POST `/users/register` - Register a new user
- POST `/users/login` - Login and get JWT token

### Tasks

- GET `/tasks` - List all tasks (supports pagination, filtering)
- POST `/tasks` - Create a new task
- PUT `/tasks/:id` - Update a task
- DELETE `/tasks/:id` - Delete a task

### Query Parameters

- `page`: Page number for pagination
- `limit`: Items per page
- `sort`: Sort field and order (e.g., "mrc" for most recent created)
- `fields`: Select specific fields
- `populate`: Populate related fields (e.g., "user")
- `status`: Filter by status
- `search`: Search in task name and description

## Project Structure

```
todo-api/
├── cmd/                        # Application entry point
│   ├── main.go                # Main application file
│   └── seed_superadmin/       # Admin user seeder
├── config/                    # Configuration files
│   └── app_config.go          # App configuration
├── constants/                 # Global constants
│   ├── error_constant.go      # Error message constants
│   └── user_roles.go         # User role definitions
├── db/                       # Database connection
│   └── mongodb.go            # MongoDB connection setup
├── internal/                 # Internal packages
│   ├── todos/               # Task-related packages
│   │   ├── dtos/           # Data Transfer Objects
│   │   ├── models/         # Task models
│   │   ├── task_inputs/    # Input validation structs
│   │   ├── tasks_controller/ # Task HTTP handlers
│   │   └── tasks_service/  # Task business logic
│   └── users/              # User-related packages
│       ├── dtos/          # User DTOs
│       ├── models/        # User models
│       ├── user_inputs/   # Input validation
│       └── users_controller/ # User HTTP handlers
│       └── users_service/ # User business logic
├── pkg/                    # Reusable packages
│   ├── middleware/        # HTTP middlewares
│   ├── mongodb/          # MongoDB utilities
│   │   └── pipeline_builder.go # Query builder
│   └── utils/            # Helper functions
│       ├── logger.go     # Logging utility
│       └── response.go   # HTTP response helpers
├── routes/               # Route definitions
│   └── route_groups.go   # API route setup
├── .env                  # Environment variables
├── .env.example         # Environment template
├── .gitignore           # Git ignore rules
├── Makefile             # Build automation
├── go.mod               # Go module file
└── README.md            # Project documentation
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
