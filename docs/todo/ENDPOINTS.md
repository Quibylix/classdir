# Endpoints

The system supports a variety of endpoints that can be sent from the client to the API to configure the presentation.

Unless otherwise specified, all commands are sent as JSON objects through a RESTful API.

## Authentication

Requests to the API require authentication via a JWT token. The token must be sent in an httpOnly, secure, and sameSite cookie.

### Auth

- [x] `POST /api/v1/auth/login`: Authenticates the user with the configured password.
  - Request Body:
  ```json
  {
      "password": "string" // The secret password from the .env file.
  }
  ```
  - Response: 204 No Content. Sets an httpOnly, secure, and sameSite cookie with the JWT token.

- [x] `POST /api/v1/auth/logout`: Logs out by clearing the authentication cookie.
  - Response: Clears the JWT cookie.

## Configuration Endpoints

### Presentation Management

- [x] `POST /api/v1/presentation`: Creates a new blank presentation.
  - Request Body:
  ```json
  {
      "id": "string", // The unique identifier for the presentation. UUIDv7 format is required.
      "title": "string" // The title of the presentation.
  }
  ```
  - Response: Returns the created presentation object.
  ```json
  {
      "data": {
          "id": "string",
          "title": "string",
          "slides": []
      }
  }
  ```

- [x] `GET /api/v1/presentation/:presentationId`: Retrieves the details of a specific presentation, including its slides.
  - Response: Returns the presentation object with its slides.
  ```json
  {
      "data": {
          "id": "string",
          "title": "string",
          "slide_order": ["string"],
          "slides": [
              {
                  "id": "string",
                  "content": "string"
              }
          ]
      }
  }
  ```

- [x] `PUT /api/v1/presentation/:presentationId`: Updates the title of the specified presentation.
  - Request Body:
  ```json
  {
      "title": "string" // The new title for the presentation.
  }
  ```
  - Response: Returns the updated presentation object.
  ```json
  {
      "data": {
          "id": "string",
          "title": "string",
          "slide_order": ["string"],
          "slides": [
              {
                  "id": "string",
                  "content": "string"
              }
          ]
      }
  }
  ```

- [x] `DELETE /api/v1/presentation/:presentationId`: Deletes the specified presentation and all its slides.
  - Response: Returns an http status code indicating success or failure.

### Slide Management

- [x] `POST /api/v1/presentation/:presentationId/slides`: Adds a new slide to the specified presentation.
  - Request Body:
  ```json
  {
      "id": "string", // The unique identifier for the slide. UUIDv7 format is required.
      "content": "string" // The content of the slide, which must be a sanitized and valid Reveal.js code.
  }
  ```
  - Response: Returns the new slide object.
  ```json
  {
      "data": {
          "id": "string",
          "content": "string"
      }
  }
  ```

- [x] `GET /api/v1/presentation/:presentationId/slides/:slideId`: Retrieves the details of a specific slide in the specified presentation.
  - Response: Returns the slide object.
  ```json
  {
      "data": {
          "id": "string",
          "content": "string"
      }
  }
  ```

- [x] `PUT /api/v1/presentation/:presentationId/slides/:slideId`: Updates the content of a specific slide in the specified presentation.
  - Request Body:
  ```json
  {
      "content": "string" // The new content for the slide, which must be a sanitized and valid Reveal.js code.
  }
  ```
  - Response: Returns the updated slide object.
  ```json
  {
      "data": {
          "id": "string",
          "content": "string"
      }
  }
  ```

- [x] `DELETE /api/v1/presentation/:presentationId/slides/:slideId`: Deletes a specific slide from the specified presentation.
  - Response: Returns an http status code indicating success or failure.

### Student Management

- [ ] `POST /api/v1/presentation/:presentationId/students`: Adds a new student to the specified presentation.
  - Request Body:
  ```json
  {
      "id": "string", // The unique identifier for the student. UUIDv7 format is required.
      "name": "string" // The name of the student.
  }
  ```
  - Response: Returns the new student object.
  ```json
  {
      "data": {
          "id": "string",
          "name": "string"
      }
  }
  ```

- [ ] `GET /api/v1/presentation/:presentationId/students`: Retrieves all students registered for the specified presentation.
  - Response: Returns an array of student objects.
  ```json
  {
      "data": [
          {
              "id": "string",
              "name": "string"
          }
      ]
  }
  ```

- [ ] `GET /api/v1/presentation/:presentationId/students/:studentId`: Retrieves the details of a specific student in the specified presentation.
  - Response: Returns the student object.
  ```json
  {
      "data": {
          "id": "string",
          "name": "string"
      }
  }
  ```

- [ ] `PUT /api/v1/presentation/:presentationId/students/:studentId`: Updates the name of a specific student in the specified presentation.
  - Request Body:
  ```json
  {
      "name": "string" // The new name for the student.
  }
  ```
  - Response: Returns the updated student object.
  ```json
  {
      "data": {
          "id": "string",
          "name": "string"
      }
  }
  ```

- [ ] `DELETE /api/v1/presentation/:presentationId/students/:studentId`: Deletes a specific student from the specified presentation.
  - Response: Returns an http status code indicating success or failure.
