# AI Secretary Assistant

An intelligent microservices-based application that acts as a personal assistant, helping with task management, scheduling, and organization.

## Project Structure

```
ai-secretary-assistant/
├── client/         # React frontend
├── secretary/      # Python-based AI service
├── server/         # Go-based backend API
└── docker-compose.yml
```

## Services

### Client (Frontend)
- Built with React
- Provides user interface for interaction with the AI secretary
- Handles task visualization and user inputs

### Secretary (AI Service)
- Python-based service
- Handles AI processing and decision making
- Manages natural language processing and task understanding

### Server (Backend)
- Go-based REST API
- Handles business logic and data persistence
- Manages authentication and request processing

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Node.js (for local development)
- Python 3.x (for local development)
- Go 1.x (for local development)

### Installation

1. Clone the repository
```bash
git clone git@github.com:proSamik/ai-secretary-assistant.git
cd ai-secretary-assistant
```

2. Start the services using Docker Compose
```bash
docker-compose up
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Secretary Service: http://localhost:5000

## Development

### Client
```bash
cd client
npm install
npm start
```

### Secretary
```bash
cd secretary
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -r requirements.txt
python app.py
```

### Server
```bash
cd server
go mod download
go run main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
