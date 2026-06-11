// File: internal/swagger/swagger.go
package swagger

// OpenAPISpec is the complete OpenAPI specification served as a constant
const OpenAPISpec = `
{
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
      "url": "http://localhost:3000",
      "description": "Local development"
    },
    {
      "url": "https://jvalleyverse.mohagussetiaone.my.id",
      "description": "Production"
    }
  ],
  "tags": [
    { "name": "Authentication", "description": "User authentication endpoints" },
    { "name": "Users", "description": "User profile management" },
    { "name": "Admin", "description": "Admin dashboard and management" },
    { "name": "Admin - Projects", "description": "Admin project management" },
    { "name": "Admin - Classes", "description": "Admin class management" },
    { "name": "Admin - Categories", "description": "Admin category management" },
    { "name": "Admin - Phases", "description": "Admin phase management" },
    { "name": "Admin - Users", "description": "Admin user management" },
    { "name": "Classes", "description": "Public class access" },
    { "name": "Certificates", "description": "Certificate endpoints (private access)" },
    { "name": "Discussions", "description": "Discussion management" },
    { "name": "Replies", "description": "Discussion replies (nested)" },
    { "name": "Showcases", "description": "User portfolio items" },
    { "name": "Gamification", "description": "Points, levels, leaderboard" },
    { "name": "Health", "description": "Server health check" },
    { "name": "Categories", "description": "Category list and browsing" },
    { "name": "Phases", "description": "Public phase access within projects" },
    { "name": "Projects", "description": "Public project listing and details" }
  ],
  "security": [
    { "BearerAuth": [] }
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
          "id": { "type": "string", "example": "" },
          "name": { "type": "string" },
          "email": { "type": "string", "format": "email" },
          "avatar": { "type": "string", "format": "url" },
          "bio": { "type": "string" },
          "role": { "type": "string", "enum": ["admin", "user"] },
          "level": { "type": "integer" },
          "points": { "type": "integer" },
          "total_points": { "type": "integer" },
          "is_active": { "type": "boolean" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "UserPublic": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "name": { "type": "string" },
          "avatar": { "type": "string", "format": "url" },
          "bio": { "type": "string" },
          "level": { "type": "integer" },
          "points": { "type": "integer" },
          "certificate_count": { "type": "integer" },
          "showcase_count": { "type": "integer" }
        }
      },
      "Category": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "name": { "type": "string", "example": "Web Development" },
          "slug": { "type": "string", "example": "web-dev" },
          "description": { "type": "string", "example": "Complete web development course" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "Project": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "title": { "type": "string" },
          "description": { "type": "string" },
          "thumbnail": { "type": "string", "format": "url" },
          "visibility": { "type": "string", "enum": ["public", "private"] },
          "category_id": { "type": "string", "example": "" },
          "admin_id": { "type": "string", "example": "" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "Phase": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "title": { "type": "string" },
          "description": { "type": "string" },
          "project_id": { "type": "string", "example": "" },
          "order_index": { "type": "integer" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "Class": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "title": { "type": "string" },
          "description": { "type": "string" },
          "content": { "type": "string" },
          "thumbnail": { "type": "string", "format": "url" },
          "duration": { "type": "integer" },
          "difficulty": { "type": "string", "enum": ["beginner", "intermediate", "advanced"] },
          "project_id": { "type": "string", "example": "" },
          "category_id": { "type": "string", "example": "" },
          "is_completed": { "type": "boolean" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "Certificate": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "code": { "type": "string" },
          "user_id": { "type": "string", "example": "" },
          "class_id": { "type": "string", "example": "" },
          "issued_at": { "type": "string", "format": "date-time" },
          "valid_until": { "type": "string", "format": "date-time" }
        }
      },
      "Discussion": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "title": { "type": "string" },
          "content": { "type": "string" },
          "user_id": { "type": "string", "example": "" },
          "user_name": { "type": "string" },
          "class_id": { "type": "string", "example": "" },
          "category_id": { "type": "string", "example": "" },
          "replies_count": { "type": "integer" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "Reply": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "content": { "type": "string" },
          "user_id": { "type": "string", "example": "" },
          "user_name": { "type": "string" },
          "discussion_id": { "type": "string", "example": "" },
          "parent_reply_id": { "type": "string", "example": "", "nullable": true },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "Showcase": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "example": "" },
          "title": { "type": "string" },
          "description": { "type": "string" },
          "media_urls": { "type": "array", "items": { "type": "string" } },
          "user_id": { "type": "string", "example": "" },
          "user_name": { "type": "string" },
          "category_id": { "type": "string", "example": "" },
          "visibility": { "type": "string", "enum": ["public", "private", "friends_only"] },
          "likes_count": { "type": "integer" },
          "is_liked": { "type": "boolean" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "LeaderboardEntry": {
        "type": "object",
        "properties": {
          "rank": { "type": "integer" },
          "user_id": { "type": "string", "example": "" },
          "name": { "type": "string" },
          "points": { "type": "integer" },
          "level": { "type": "integer" },
          "badge": { "type": "string" }
        }
      },
      "Pagination": {
        "type": "object",
        "properties": {
          "page": { "type": "integer" },
          "limit": { "type": "integer" },
          "total": { "type": "integer" }
        }
      }
    }
  },
  "paths": {
    "/api/auth/login": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Login user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["email", "password"],
                "properties": {
                  "email": { "type": "string", "format": "email" },
                  "password": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login successful",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "access_token": { "type": "string" },
                    "refresh_token": { "type": "string" },
                    "user": { "$ref": "#/components/schemas/User" }
                  }
                }
              }
            }
          },
          "401": { "description": "Invalid credentials" }
        }
      }
    },
    "/api/auth/register": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Register new user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["email", "password", "name"],
                "properties": {
                  "email": { "type": "string", "format": "email" },
                  "password": { "type": "string", "minLength": 6 },
                  "name": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "User created successfully",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/User" }
              }
            }
          },
          "400": { "description": "Invalid input" },
          "409": { "description": "Email already exists" }
        }
      }
    },
    "/api/categories": {
      "get": {
        "tags": ["Categories"],
        "summary": "List all categories",
        "responses": {
          "200": {
            "description": "List of categories",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": { "$ref": "#/components/schemas/Category" }
                }
              }
            }
          }
        }
      }
    },
    "/api/categories/{category_id}/projects": {
      "get": {
        "tags": ["Categories"],
        "summary": "List projects in a category",
        "parameters": [
          {
            "name": "category_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": {
            "description": "List of projects in category",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": { "$ref": "#/components/schemas/Project" }
                }
              }
            }
          },
          "404": { "description": "Category not found" }
        }
      }
    },
    "/api/categories/{slug}": {
      "get": {
        "tags": ["Categories"],
        "summary": "Get category detail by slug",
        "parameters": [
          {
            "name": "slug",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "web-dev" }
          }
        ],
        "responses": {
          "200": {
            "description": "Category details",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Category" }
              }
            }
          },
          "404": { "description": "Category not found" }
        }
      }
    },
    "/api/certificates": {
      "get": {
        "tags": ["Certificates"],
        "summary": "Get user certificates",
        "security": [{ "BearerAuth": [] }],
        "responses": {
          "200": { "description": "User certificates" },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/certificates/{code}": {
      "get": {
        "tags": ["Certificates"],
        "summary": "View certificate (Private - Owner only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "code",
            "in": "path",
            "required": true,
            "schema": { "type": "string" }
          }
        ],
        "responses": {
          "200": { "description": "Certificate details" },
          "403": { "description": "Forbidden - Not your certificate" },
          "404": { "description": "Certificate not found" }
        }
      }
    },
    "/api/classes/{id}": {
      "get": {
        "tags": ["Classes"],
        "summary": "Get class by ID",
        "description": "Retrieve a single class by its ID. Public access.",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": {
            "description": "Class details",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Class" }
              }
            }
          },
          "404": { "description": "Class not found" }
        }
      }
    },
    "/api/classes/{id}/complete": {
      "post": {
        "tags": ["Classes"],
        "summary": "Mark class as completed",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "Class marked as completed" },
          "401": { "description": "Unauthorized" },
          "404": { "description": "Class not found" }
        }
      }
    },
    "/api/discussions": {
      "get": {
        "tags": ["Discussions"],
        "summary": "List discussions",
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer" } },
          { "name": "class_id", "in": "query", "schema": { "type": "string", "example": "" } },
          { "name": "category_id", "in": "query", "schema": { "type": "string", "example": "" } }
        ],
        "responses": {
          "200": { "description": "Discussions list" }
        }
      },
      "post": {
        "tags": ["Discussions"],
        "summary": "Create discussion",
        "security": [{ "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["title", "content"],
                "properties": {
                  "title": { "type": "string" },
                  "content": { "type": "string" },
                  "class_id": { "type": "string", "example": "" },
                  "category_id": { "type": "string", "example": "" }
                }
              }
            }
          }
        },
        "responses": {
          "201": { "description": "Discussion created" },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/discussions/{id}": {
      "get": {
        "tags": ["Discussions"],
        "summary": "Get discussion with threaded replies",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "Discussion details with replies" },
          "404": { "description": "Discussion not found" }
        }
      },
      "put": {
        "tags": ["Discussions"],
        "summary": "Update discussion",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "title": { "type": "string" },
                  "content": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Discussion updated" },
          "403": { "description": "Forbidden" },
          "404": { "description": "Not found" }
        }
      },
      "delete": {
        "tags": ["Discussions"],
        "summary": "Delete discussion",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "204": { "description": "Discussion deleted" },
          "403": { "description": "Forbidden" }
        }
      }
    },
    "/api/discussions/{id}/replies": {
      "get": {
        "tags": ["Replies"],
        "summary": "Get discussion replies",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "Replies list" }
        }
      },
      "post": {
        "tags": ["Replies"],
        "summary": "Reply to discussion",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["content"],
                "properties": {
                  "content": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "201": { "description": "Reply created" },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/health": {
      "get": {
        "tags": ["Health"],
        "summary": "Health check endpoint",
        "responses": {
          "200": {
            "description": "Server is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": { "type": "string", "example": "ok" }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/leaderboard": {
      "get": {
        "tags": ["Gamification"],
        "summary": "Get top users by points",
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 50 } }
        ],
        "responses": {
          "200": {
            "description": "Leaderboard",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/LeaderboardEntry" }
                    },
                    "pagination": { "$ref": "#/components/schemas/Pagination" }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/levels": {
      "get": {
        "tags": ["Gamification"],
        "summary": "Get level progression info",
        "responses": {
          "200": { "description": "Level information" }
        }
      }
    },
    "/api/projects": {
      "get": {
        "tags": ["Projects"],
        "summary": "List public projects",
        "description": "Get a paginated list of public projects.",
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 20 } }
        ],
        "responses": {
          "200": {
            "description": "List of projects with pagination",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Project" }
                    },
                    "pagination": { "$ref": "#/components/schemas/Pagination" }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/projects/{project_id}": {
      "get": {
        "tags": ["Projects"],
        "summary": "Get project with phases and classes",
        "parameters": [
          {
            "name": "project_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "Project details with phases" },
          "404": { "description": "Project not found" }
        }
      }
    },
    "/api/projects/{project_id}/classes": {
      "get": {
        "tags": ["Classes"],
        "summary": "List classes in a project",
        "parameters": [
          {
            "name": "project_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          },
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 20 } }
        ],
        "responses": {
          "200": {
            "description": "Classes list",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Class" }
                    },
                    "pagination": { "$ref": "#/components/schemas/Pagination" }
                  }
                }
              }
            }
          },
          "404": { "description": "Project not found" }
        }
      }
    },
    "/api/projects/{project_id}/classes/{slug}": {
      "get": {
        "tags": ["Classes"],
        "summary": "Get class by project and slug",
        "parameters": [
          {
            "name": "project_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          },
          {
            "name": "slug",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "intro-to-go" }
          }
        ],
        "responses": {
          "200": { "description": "Class details" },
          "404": { "description": "Class not found" }
        }
      }
    },
    "/api/projects/{project_id}/phases": {
      "get": {
        "tags": ["Phases"],
        "summary": "List phases by project",
        "parameters": [
          {
            "name": "project_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": {
            "description": "Phases list",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Phase" }
                    }
                  }
                }
              }
            }
          },
          "404": { "description": "Project not found" }
        }
      }
    },
    "/api/projects/{project_id}/phases/{phase_id}": {
      "get": {
        "tags": ["Phases"],
        "summary": "Get phase with classes",
        "parameters": [
          {
            "name": "project_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          },
          {
            "name": "phase_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": {
            "description": "Phase details",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Phase" }
              }
            }
          },
          "404": { "description": "Phase not found" }
        }
      }
    },
    "/api/projects/{project_id}/phases/{phase_id}/classes": {
      "get": {
        "tags": ["Classes"],
        "summary": "List classes in a phase",
        "parameters": [
          {
            "name": "project_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          },
          {
            "name": "phase_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": {
            "description": "Classes list",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "total": { "type": "integer" },
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Class" }
                    }
                  }
                }
              }
            }
          },
          "404": { "description": "Phase not found" }
        }
      }
    },
    "/api/replies/{id}": {
      "put": {
        "tags": ["Replies"],
        "summary": "Update reply",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["content"],
                "properties": {
                  "content": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Reply updated" },
          "401": { "description": "Unauthorized" },
          "403": { "description": "Forbidden" },
          "404": { "description": "Reply not found" }
        }
      },
      "delete": {
        "tags": ["Replies"],
        "summary": "Delete reply",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "204": { "description": "Reply deleted" },
          "401": { "description": "Unauthorized" },
          "403": { "description": "Forbidden" },
          "404": { "description": "Reply not found" }
        }
      }
    },
    "/api/replies/{id}/replies": {
      "post": {
        "tags": ["Replies"],
        "summary": "Reply to a reply (nested)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["content"],
                "properties": {
                  "content": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "201": { "description": "Nested reply created" },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/showcases": {
      "get": {
        "tags": ["Showcases"],
        "summary": "List showcases",
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer" } },
          { "name": "category_id", "in": "query", "schema": { "type": "string", "example": "" } },
          { "name": "sort", "in": "query", "schema": { "type": "string", "enum": ["newest", "trending", "most_liked"] } }
        ],
        "responses": {
          "200": { "description": "Showcases list" }
        }
      },
      "post": {
        "tags": ["Showcases"],
        "summary": "Create showcase",
        "security": [{ "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["title", "category_id"],
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "media_urls": { "type": "array", "items": { "type": "string", "format": "url" } },
                  "category_id": { "type": "string", "example": "" },
                  "visibility": { "type": "string", "enum": ["public", "private", "friends_only"] }
                }
              }
            }
          }
        },
        "responses": {
          "201": { "description": "Showcase created" },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/showcases/{id}": {
      "get": {
        "tags": ["Showcases"],
        "summary": "Get showcase details",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "Showcase details" },
          "404": { "description": "Showcase not found" }
        }
      },
      "put": {
        "tags": ["Showcases"],
        "summary": "Update showcase (Owner only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "visibility": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Showcase updated" },
          "403": { "description": "Forbidden" }
        }
      },
      "delete": {
        "tags": ["Showcases"],
        "summary": "Delete showcase",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "204": { "description": "Showcase deleted" },
          "403": { "description": "Forbidden" }
        }
      }
    },
    "/api/showcases/{id}/like": {
      "post": {
        "tags": ["Showcases"],
        "summary": "Like showcase",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "Showcase liked" },
          "401": { "description": "Unauthorized" }
        }
      },
      "delete": {
        "tags": ["Showcases"],
        "summary": "Unlike showcase",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "204": { "description": "Like removed" },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/users/me": {
      "get": {
        "tags": ["Users"],
        "summary": "Get current user profile",
        "security": [{ "BearerAuth": [] }],
        "responses": {
          "200": {
            "description": "User profile",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/User" }
              }
            }
          },
          "401": { "description": "Unauthorized" }
        }
      },
      "put": {
        "tags": ["Users"],
        "summary": "Update user profile",
        "security": [{ "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "name": { "type": "string" },
                  "bio": { "type": "string" },
                  "avatar": { "type": "string", "format": "url" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Profile updated",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/User" }
              }
            }
          },
          "401": { "description": "Unauthorized" }
        }
      }
    },
    "/api/users/me/activity": {
      "get": {
        "tags": ["Gamification"],
        "summary": "Get user activity log",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer" } },
          { "name": "limit", "in": "query", "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": { "description": "Activity log" }
        }
      }
    },
    "/api/users/{id}": {
      "get": {
        "tags": ["Users"],
        "summary": "Get public user profile",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": {
            "description": "User profile",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserPublic" }
              }
            }
          },
          "404": { "description": "User not found" }
        }
      }
    },
    "/api/users/{id}/points": {
      "get": {
        "tags": ["Gamification"],
        "summary": "Get user points and level",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "User points info" }
        }
      }
    },
    "/api/admin/categories": {
      "post": {
        "tags": ["Admin - Categories"],
        "summary": "Create category (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["name", "slug"],
                "properties": {
                  "name": { "type": "string", "example": "Web Development" },
                  "slug": { "type": "string", "example": "web-dev" },
                  "description": { "type": "string", "example": "Complete web development course" }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Category created",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Category" }
              }
            }
          },
          "401": { "description": "Unauthorized" },
          "403": { "description": "Forbidden - Admin role required" }
        }
      }
    },
    "/api/admin/categories/{id}": {
      "put": {
        "tags": ["Admin - Categories"],
        "summary": "Update category (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "name": { "type": "string", "example": "Web Development" },
                  "slug": { "type": "string", "example": "web-dev" },
                  "description": { "type": "string", "example": "Complete web development course" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Category updated",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Category" }
              }
            }
          },
          "401": { "description": "Unauthorized" },
          "403": { "description": "Forbidden - Admin role required" },
          "404": { "description": "Category not found" }
        }
      },
      "delete": {
        "tags": ["Admin - Categories"],
        "summary": "Delete category (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "204": { "description": "Category deleted" },
          "401": { "description": "Unauthorized" },
          "403": { "description": "Forbidden - Admin role required" },
          "404": { "description": "Category not found" }
        }
      }
    },
    "/api/admin/classes": {
      "post": {
        "tags": ["Admin - Classes"],
        "summary": "Create class",
        "security": [{ "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["title", "project_id", "category_id"],
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "content": { "type": "string" },
                  "thumbnail": { "type": "string", "format": "url" },
                  "duration": { "type": "integer" },
                  "difficulty": { "type": "string", "enum": ["beginner", "intermediate", "advanced"] },
                  "project_id": { "type": "string", "example": "" },
                  "category_id": { "type": "string", "example": "" }
                }
              }
            }
          }
        },
        "responses": {
          "201": { "description": "Class created" },
          "403": { "description": "Forbidden" }
        }
      }
    },
    "/api/admin/classes/{id}": {
      "put": {
        "tags": ["Admin - Classes"],
        "summary": "Update class",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "content": { "type": "string" },
                  "duration": { "type": "integer" },
                  "difficulty": { "type": "string", "enum": ["beginner", "intermediate", "advanced"] }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Class updated" },
          "403": { "description": "Forbidden - Admin only" },
          "404": { "description": "Class not found" }
        }
      },
      "delete": {
        "tags": ["Admin - Classes"],
        "summary": "Delete class",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "204": { "description": "Class deleted" },
          "403": { "description": "Forbidden - Admin only" },
          "404": { "description": "Class not found" }
        }
      }
    },
    "/api/admin/classes/{id}/details": {
      "post": {
        "tags": ["Admin - Classes"],
        "summary": "Create class detail",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "title": {
                    "type": "string"
                  },
                  "content": {
                    "type": "string"
                  },
                  "order_index": {
                    "type": "integer"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Class detail created"
          },
          "403": {
            "description": "Forbidden - Admin only"
          },
          "404": {
            "description": "Class not found"
          }
        }
      }
  },
    "/api/admin/dashboard": {
      "get": {
        "tags": ["Admin"],
        "summary": "Admin dashboard",
        "security": [{ "BearerAuth": [] }],
        "responses": {
          "200": { "description": "Welcome to admin dashboard" },
          "403": { "description": "Forbidden - Admin only" }
        }
      }
    },
    "/api/admin/phases/{phase_id}": {
      "put": {
        "tags": ["Admin - Phases"],
        "summary": "Update phase (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "phase_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "order_index": { "type": "integer" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Phase updated",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Phase" }
              }
            }
          },
          "403": { "description": "Forbidden - Admin only" },
          "404": { "description": "Phase not found" }
        }
      },
      "delete": {
        "tags": ["Admin - Phases"],
        "summary": "Delete phase (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "phase_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "200": { "description": "Phase deleted" },
          "403": { "description": "Forbidden - Admin only" },
          "404": { "description": "Phase not found" }
        }
      }
    },
    "/api/admin/projects": {
      "get": {
        "tags": ["Admin - Projects"],
        "summary": "List all projects (Admin view)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer" } },
          { "name": "limit", "in": "query", "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": {
            "description": "Projects list",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/Project" }
                    },
                    "pagination": { "$ref": "#/components/schemas/Pagination" }
                  }
                }
              }
            }
          },
          "403": { "description": "Forbidden" }
        }
      },
      "post": {
        "tags": ["Admin - Projects"],
        "summary": "Create project (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["title", "category_id"],
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "thumbnail": { "type": "string", "format": "url" },
                  "category_id": { "type": "string", "example": "" },
                  "visibility": { "type": "string", "enum": ["public", "private"] }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Project created",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Project" }
              }
            }
          },
          "403": { "description": "Forbidden - Admin only" }
        }
      }
    },
    "/api/admin/projects/{id}": {
      "put": {
        "tags": ["Admin - Projects"],
        "summary": "Update project",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "visibility": { "type": "string", "enum": ["public", "private"] }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Project updated" },
          "403": { "description": "Forbidden" },
          "404": { "description": "Project not found" }
        }
      },
      "delete": {
        "tags": ["Admin - Projects"],
        "summary": "Delete project",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "responses": {
          "204": { "description": "Project deleted" },
          "403": { "description": "Forbidden" },
          "404": { "description": "Project not found" }
        }
      }
    },
    "/api/admin/projects/{project_id}/phases": {
      "post": {
        "tags": ["Admin - Phases"],
        "summary": "Create phase (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "project_id",
            "in": "path",
            "required": true,
            "schema": { "type": "string", "example": "" }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["title"],
                "properties": {
                  "title": { "type": "string" },
                  "description": { "type": "string" },
                  "order_index": { "type": "integer" }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Phase created",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Phase" }
              }
            }
          },
          "403": { "description": "Forbidden - Admin only" },
          "404": { "description": "Project not found" }
        }
      }
    },
    "/api/admin/users": {
      "get": {
        "tags": ["Admin - Users"],
        "summary": "List all users with pagination (Admin only)",
        "security": [{ "BearerAuth": [] }],
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "description": "Page number",
            "schema": { "type": "integer", "default": 1 }
          },
          {
            "name": "limit",
            "in": "query",
            "description": "Items per page (max: 100)",
            "schema": { "type": "integer", "default": 20 }
          }
        ],
        "responses": {
          "200": {
            "description": "Users list with pagination",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/User" }
                    },
                    "pagination": { "$ref": "#/components/schemas/Pagination" }
                  }
                }
              }
            }
          },
          "401": { "description": "Unauthorized - Missing or invalid token" },
          "403": { "description": "Forbidden - Admin role required" }
        }
      }
    }
  }
}
`
