package swagger

// OpenAPISpec is the complete OpenAPI specification served as a constant
const OpenAPISpec = `{
  "openapi": "3.0.0",
  "info": {
    "title": "JValleyVerse API",
    "description": "Community-driven learning platform with gamification, certification, and showcase features",
    "version": "1.0.0",
    "contact": {
      "name": "JValleyVerse Team"
    },
    "license": {
      "name": "MIT"
    }
  },
  "servers": [
    {
      "url": "http://localhost:3000/api/v1",
      "description": "Development Server"
    }
  ],
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": {"type": "integer"},
          "email": {"type": "string"},
          "name": {"type": "string"},
          "role": {"type": "string", "enum": ["admin", "user"]},
          "points": {"type": "integer"},
          "level": {"type": "integer"}
        }
      }
    }
  },
  "paths": {
    "/auth/register": {
      "post": {
        "summary": "Register new user",
        "tags": ["Auth"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": {"type": "string"},
                  "password": {"type": "string"},
                  "name": {"type": "string"}
                }
              }
            }
          }
        },
        "responses": {
          "201": {"description": "User registered successfully"}
        }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "Login user",
        "tags": ["Auth"],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": {"type": "string"},
                  "password": {"type": "string"}
                }
              }
            }
          }
        },
        "responses": {
          "200": {"description": "Login successful"}
        }
      }
    },
    "/showcases": {
      "post": {
        "summary": "Create showcase",
        "tags": ["Showcases"],
        "security": [{"BearerAuth": []}],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "title": {"type": "string"},
                  "description": {"type": "string"},
                  "media_urls": {"type": "array", "items": {"type": "string"}},
                  "category_id": {"type": "integer"},
                  "visibility": {"type": "string"}
                }
              }
            }
          }
        },
        "responses": {
          "201": {"description": "Showcase created"}
        }
      }
    },
    "/showcases/{id}/like": {
      "post": {
        "summary": "Like a showcase",
        "tags": ["Showcases"],
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {"description": "Showcase liked"}
        }
      }
    },
    "/classes/{id}/complete": {
      "post": {
        "summary": "Complete a class",
        "tags": ["Classes"],
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {"description": "Class completed"}
        }
      }
    }
  }
}`
