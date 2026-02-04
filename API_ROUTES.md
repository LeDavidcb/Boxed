# Boxed API Routes Documentation

This document provides detailed information about the routes exposed by the Boxed application, outlining the expected input, output, and functionality.

---

## Authentication Routes

### GET `/auth/login`
- **Description**: Logs the user into the application.
- **Input**: None.
- **Output**:
  - **Success**:
    ```json
    {
      "token": "exampleJWT",
      "expiresAt": "2026-02-05T12:00:00Z"
    }
    ```
  - **Failure**:
    ```json
    {
      "error": "Invalid credentials"
    }
    ```

### GET `/auth/register`
- **Description**: Handles new-user registration.
- **Input**: (Pending clarification in code).

### GET `/auth/refresh`
- **Description**: Refreshes an authentication token.
- **Input**: Authorization header with a valid token.
- **Output**:
  - **Success**:
    ```json
    {
      "token": "newExampleJWT",
      "expiresAt": "2026-02-12T12:00:00Z"
    }
    ```
  - **Failure**:
    ```json
    {
      "error": "Token expired or invalid"
    }
    ```

---

## File Management Routes

### POST `/api/upload-file`
- **Description**: Uploads a single file.
- **Input**:
  - **Headers**:
    - `Authorization`: Bearer <JWT>.
  - **Body**: Multipart form data with a field `file`.
- **Output**:
  - **Success** (`201 Created`): No content.
  - **Failure**:
    ```json
    {
      "error": "Invalid file input"
    }
    ```

### POST `/api/upload-files`
- **Description**: Allows users to bulk-upload multiple files.
- **Input**:
  - **Headers**:
    - `Authorization`: Bearer <JWT>.
  - **Body**: Multipart form data with a field `files` (array of files).
- **Output**:
  - **Success** (`201 Created`): No content.
  - **Failure**:
    ```json
    {
      "error": "Invalid form data"
    }
    ```

### GET `/api/get-file`
- **Description**: Retrieves a file’s metadata by UUID header field.
- **Input**:
  - **Headers**:
    - `uuid`: File identifier.
- **Output**:
  - **Success**:
    ```json
    {
      "id": "file-uuid",
      "name": "example.txt",
      "size": "14kb"
    }
    ```
  - **Failure**:
    ```json
    {
      "error": "File not found"
    }
    ```

### DELETE `/api/delete-file`
- **Description**: Deletes file and database metadata using the UUID.
- **Input**:
  - **Headers**:
    - `uuid`: File identifier.
- **Output**:
  - **Success**: No content (`204`).
  - **Failure**:
    ```json
    {
      "error": "File not found"
    }
    ```

### GET `/api/get-files`
- **Description**: Gets all user-associated file metadata.
- **Input**:
  - **Headers**:
    - `Authorization`: Bearer <JWT>.
- **Output**:
  - **Success**:
    ```json
    {
      "length": 5,
      "files": [
        {
          "id": "1",
          "name": "file1.txt",
          "size": "14kb"
        },
        {
          "id": "2",
          "name": "file2.png",
          "size": "138kb"
        }
      ]
    }
    ```

### GET `/api/serve-file`
- **Description**: Streams the requested file if the user is authorized.
- **Input**:
  - **Headers**:
    - `uuid`: File identifier.
- **Output**:
  - **Success**: Streams binary content.
  - **Failure**:
    ```json
    {
      "error": "Access forbidden or file not found"
    }
    ```

---

## Notes
- All routes require JWT-based authentication unless explicitly noted otherwise.
- Use valid UUIDs for file-related operations.

---