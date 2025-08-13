# Digital Wellbeing Backend

A Go backend API for the Digital Wellbeing application that helps users track and manage their social media usage.

## Features

- **User Authentication**: JWT-based auth with registration and login
- **Goal Management**: Set and track daily usage limits for apps
- **Check-ins**: Daily mood and screen time tracking
- **AI Chat**: Wellness guidance and support
- **Usage Analytics**: Track and analyze screen time patterns
- **Admin Panel**: User management for administrators

## Tech Stack

- **Framework**: Gin (Go HTTP framework)
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis
- Docker (optional)

### Environment Setup

1. Copy the environment file:
```bash
cp .env.example .env
```

2. Update the `.env` file with your configuration:
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=digital_wellbeing
JWT_SECRET=your-super-secret-jwt-key
```

### Running with Docker

```bash
docker-compose up -d
```

### Running Locally

1. Install dependencies:
```bash
go mod tidy
```

2. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user

### Users (Protected)
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile

### Goals (Protected)
- `GET /api/v1/goals` - Get user goals
- `POST /api/v1/goals` - Create new goal
- `PUT /api/v1/goals/:id` - Update goal
- `DELETE /api/v1/goals/:id` - Delete goal
- `GET /api/v1/goals/progress` - Get goal progress

### Check-ins (Protected)
- `POST /api/v1/checkins` - Create daily check-in
- `GET /api/v1/checkins` - Get check-in history
- `GET /api/v1/checkins/today` - Get today's check-in
- `GET /api/v1/checkins/stats` - Get mood statistics

### Chat (Protected)
- `POST /api/v1/chat/message` - Send message to AI
- `GET /api/v1/chat/history` - Get chat history

### Usage (Protected)
- `POST /api/v1/usage/log` - Log app usage
- `GET /api/v1/usage/daily` - Get daily usage
- `GET /api/v1/usage/weekly` - Get weekly usage
- `GET /api/v1/usage/stats` - Get usage statistics

### Admin (Protected, Admin only)
- `GET /api/v1/admin/users` - Get all users

## Database Schema

The application uses the following main entities:

- **Users**: User accounts and profiles
- **Goals**: Daily usage limits for apps
- **CheckIns**: Daily mood and screen time entries
- **ChatLogs**: AI chat conversation history
- **UsageLogs**: App usage tracking data

## Development

### Project Structure

```
backend/
├── internal/
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and setup
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   ├── routes/         # Route definitions
│   ├── services/       # Business logic
│   └── utils/          # Utility functions
├── main.go            # Application entry point
├── go.mod             # Go module file
└── docker-compose.yml # Docker setup
```

### Adding New Features

1. Define models in `internal/models/`
2. Create service logic in `internal/services/`
3. Add HTTP handlers in `internal/handlers/`
4. Register routes in `internal/routes/`
5. Add database migrations if needed

## Deployment

### Docker Production

```bash
docker build -t digital-wellbeing-backend .
docker run -p 8080:8080 --env-file .env digital-wellbeing-backend
```

### Environment Variables

Make sure to set these in production:

- `GIN_MODE=release`
- `JWT_SECRET` (use a strong, random secret)
- Database credentials
- Redis credentials
- `CORS_ORIGINS` (your frontend URLs)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.