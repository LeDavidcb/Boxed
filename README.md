<img src="./assets/boxedlogo.png" alt="isolated" width="200"/>

# Boxed
### Boxed is a lightweight, very simple self-hosted file management system

![Go Version](https://img.shields.io/github/go-mod/go-version/David/Boxed)
![License](https://img.shields.io/github/license/David/Boxed)

_This project is a work-in-progress._

---
## Prerequisites

Before setting up the project, ensure you have the following dependencies installed:
- **Node.js & npm**
- **Go compiler** (v1.25+ recommended)
- **FFmpeg** (for thumbnail generation)
- **PostgreSQL**

Additionally, you must create a PostgreSQL database for the server to save persistent data (it is recommended to name it `Boxed`).

## Installation 

To clone and run this project locally, follow the steps below:

1. Clone this repository:
   ```bash
   git clone https://github.com/David/Boxed.git
   cd Boxed
   ```

2. Build the project:
   ```bash
   make build
   ```
   *Note: If a `.env` file is not found, the build process will automatically run the interactive environment setup tool (`envcli`).*

3. Run the project:
   ```bash
   make run
   ```

---

## Environment Variables

These variables are required for the application to function. They are automatically managed by `make` and `envcli`, but can be edited manually in the `.env` file:

| Variable | Description | Example |
| :--- | :--- | :--- |
| `DB_URL` | PostgreSQL connection string | `postgresql://user:pass@localhost:5432/Boxed` |
| `BACKEND_PORT` | Port the server will listen on | `8080` |
| `FOLDER_PATH` | Directory where files will be stored | `/home/user/uploads` |
| `JWT_SECRET` | Secret key for signing tokens | `your-super-secret-key` |

---

## Project Structure

```text
.
├── cmd/
│   ├── api/            # Main API server entry point
│   └── cli/            # CLI utilities (envcli)
├── internal/
│   ├── auth/           # Authentication logic (JWT, Controllers, Services)
│   ├── files/          # File management logic (Upload, Serve, Thumbnails)
│   ├── common/         # Shared types and utilities
│   └── router.go       # Route definitions
├── migrations/         # SQL migration files
├── repositories/       # Database access layer
├── assets/             # Images and static assets for README
└── Makefile            # Automation commands
```

---

## API Endpoints

### Authentication
All routes under `/api` require a valid JWT in the `Authorization` header (`Bearer <token>`).

| Method | Route | Description | Required Input |
| :--- | :--- | :--- | :--- |
| `GET` | `/auth/login` | Authenticate user | JSON: `email`, `password` |
| `GET` | `/auth/register` | Register new user | JSON: `nickname`, `email`, `password` |
| `GET` | `/auth/refresh` | Refresh JWT token | Header: `refresh-token` |

### File Management (Protected / Must provide JWT.)

| Method | Route | Description | Required Input |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/upload-file` | Upload a single file | Multipart field: `file` |
| `POST` | `/api/upload-files` | Upload multiple files | Multipart field: `files` |
| `GET` | `/api/get-file` | Get file metadata | Header: `uuid` |
| `GET` | `/api/get-files` | List all user files | None |
| `GET` | `/api/serve-file` | Download file content | Header: `uuid` |
| `GET` | `/api/serve-thumbnail` | Get file thumbnail | Header: `uuid` |
| `DELETE` | `/api/delete-file` | Delete a file | Header: `uuid` |

---

## Usage Examples (Curl)

### 1. Register a new user
```bash
curl -X GET http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"nickname": "user1", "email": "user1@example.com", "password": "password123"}'
```

### 2. Login
```bash
curl -X GET http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user1@example.com", "password": "password123"}'
```
*Response will include a JWT. Use this token for subsequent requests.*

### 3. Upload a file
```bash
curl -X POST http://localhost:8080/api/upload-file \
  -H "Authorization: Bearer <your-jwt-token>" \
  -F "file=@/path/to/your/file.txt"
```

---

## License 

This project uses a [LICENSE](LICENSE) file. Please check the details before using the code.

---

## Future Plans 

- [x] ~**Thumbnails**~
  ~Enhance file metadata retrieval to include thumbnail IDs for compatible file types (e.g., videos, images) when accessed through relevant API endpoints.~

- [ ] **Frontend**
  While this project is primarily a backend service, creating a default frontend would be valuable, similar to how the `Jellyfin` project operates. This would provide an out-of-the-box user interface for managing files. Additionally, with the [API documentation](#api-endpoints), developers should find it straightforward to create custom frontends that consume these endpoints.

- [x] ~**Refactor**~
  ~Improve the current project structure to be more intuitive and inviting for contributors.~

- [x] ~**Standardize API responses**~
  ~Define a consistent response format across the repository (e.g. success, error, metadata) to improve predictability.~
  Mostly done, just have to documentate errors.

- [x] ~**Create a CLI (Coommand line interface) tool for `.env` management**~
  ~Build a CLI utility to initialize and update `.env` files without manual editing.~

- [x] ~**Add a Makefile for easier server setup**~
  ~Provide common commands (build, run, test, etc) to facilitate startup.~

