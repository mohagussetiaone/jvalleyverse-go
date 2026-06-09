package swagger

// OpenAPISpec is the complete OpenAPI specification served as a constant
const OpenAPISpec = `{
  "components": {
    "schemas": {
      "Category": {
        "properties": {
          "created_at": {
            "format": "date-time",
            "type": "string"
          },
          "description": {
            "example": "Complete web development course",
            "type": "string"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "name": {
            "example": "Web Development",
            "type": "string"
          },
          "slug": {
            "example": "web-dev",
            "type": "string"
          },
          "updated_at": {
            "format": "date-time",
            "type": "string"
          }
        },
        "type": "object"
      },
      "Certificate": {
        "properties": {
          "class_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "code": {
            "type": "string"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "issued_at": {
            "format": "date-time",
            "type": "string"
          },
          "user_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "valid_until": {
            "format": "date-time",
            "type": "string"
          }
        },
        "type": "object"
      },
      "Class": {
        "properties": {
          "category_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "content": {
            "type": "string"
          },
          "created_at": {
            "format": "date-time",
            "type": "string"
          },
          "description": {
            "type": "string"
          },
          "difficulty": {
            "enum": [
              "beginner",
              "intermediate",
              "advanced"
            ],
            "type": "string"
          },
          "duration": {
            "type": "integer"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "is_completed": {
            "type": "boolean"
          },
          "project_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "thumbnail": {
            "format": "url",
            "type": "string"
          },
          "title": {
            "type": "string"
          },
          "updated_at": {
            "format": "date-time",
            "type": "string"
          }
        },
        "type": "object"
      },
      "Discussion": {
        "properties": {
          "category_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "class_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "content": {
            "type": "string"
          },
          "created_at": {
            "format": "date-time",
            "type": "string"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "replies_count": {
            "type": "integer"
          },
          "title": {
            "type": "string"
          },
          "updated_at": {
            "format": "date-time",
            "type": "string"
          },
          "user_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "user_name": {
            "type": "string"
          }
        },
        "type": "object"
      },
      "LeaderboardEntry": {
        "properties": {
          "badge": {
            "type": "string"
          },
          "level": {
            "type": "integer"
          },
          "name": {
            "type": "string"
          },
          "points": {
            "type": "integer"
          },
          "rank": {
            "type": "integer"
          },
          "user_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          }
        },
        "type": "object"
      },
      "Pagination": {
        "properties": {
          "limit": {
            "type": "integer"
          },
          "page": {
            "type": "integer"
          },
          "total": {
            "type": "integer"
          }
        },
        "type": "object"
      },
      "Project": {
        "properties": {
          "admin_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "category_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "created_at": {
            "format": "date-time",
            "type": "string"
          },
          "description": {
            "type": "string"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "thumbnail": {
            "format": "url",
            "type": "string"
          },
          "title": {
            "type": "string"
          },
          "visibility": {
            "enum": [
              "public",
              "private"
            ],
            "type": "string"
          }
        },
        "type": "object"
      },
      "Reply": {
        "properties": {
          "content": {
            "type": "string"
          },
          "created_at": {
            "format": "date-time",
            "type": "string"
          },
          "discussion_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "parent_reply_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "nullable": true,
            "type": "string"
          },
          "updated_at": {
            "format": "date-time",
            "type": "string"
          },
          "user_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "user_name": {
            "type": "string"
          }
        },
        "type": "object"
      },
      "Showcase": {
        "properties": {
          "category_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "created_at": {
            "format": "date-time",
            "type": "string"
          },
          "description": {
            "type": "string"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "is_liked": {
            "type": "boolean"
          },
          "likes_count": {
            "type": "integer"
          },
          "media_urls": {
            "items": {
              "type": "string"
            },
            "type": "array"
          },
          "title": {
            "type": "string"
          },
          "updated_at": {
            "format": "date-time",
            "type": "string"
          },
          "user_id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "user_name": {
            "type": "string"
          },
          "visibility": {
            "enum": [
              "public",
              "private",
              "friends_only"
            ],
            "type": "string"
          }
        },
        "type": "object"
      },
      "User": {
        "properties": {
          "avatar": {
            "format": "url",
            "type": "string"
          },
          "bio": {
            "type": "string"
          },
          "created_at": {
            "format": "date-time",
            "type": "string"
          },
          "email": {
            "format": "email",
            "type": "string"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "is_active": {
            "type": "boolean"
          },
          "level": {
            "type": "integer"
          },
          "name": {
            "type": "string"
          },
          "points": {
            "type": "integer"
          },
          "role": {
            "enum": [
              "admin",
              "user"
            ],
            "type": "string"
          },
          "total_points": {
            "type": "integer"
          },
          "updated_at": {
            "format": "date-time",
            "type": "string"
          }
        },
        "type": "object"
      },
      "UserPublic": {
        "properties": {
          "avatar": {
            "format": "url",
            "type": "string"
          },
          "bio": {
            "type": "string"
          },
          "certificate_count": {
            "type": "integer"
          },
          "id": {
            "example": "cmq6b6ehc000010uj9ewttltv",
            "type": "string"
          },
          "level": {
            "type": "integer"
          },
          "name": {
            "type": "string"
          },
          "points": {
            "type": "integer"
          },
          "showcase_count": {
            "type": "integer"
          }
        },
        "type": "object"
      }
    },
    "securitySchemes": {
      "BearerAuth": {
        "bearerFormat": "JWT",
        "scheme": "bearer",
        "type": "http"
      }
    }
  },
  "info": {
    "contact": {
      "name": "JValleyVerse Team"
    },
    "description": "Community-driven learning platform with gamification, certification, and showcase features",
    "license": {
      "name": "MIT"
    },
    "title": "JValleyVerse API",
    "version": "1.0.0"
  },
  "openapi": "3.0.0",
  "paths": {
    "/api/admin/categories": {
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "description": {
                    "example": "Complete web development course",
                    "type": "string"
                  },
                  "name": {
                    "example": "Web Development",
                    "type": "string"
                  },
                  "slug": {
                    "example": "web-dev",
                    "type": "string"
                  }
                },
                "required": [
                  "name",
                  "slug"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Category"
                }
              }
            },
            "description": "Category created"
          },
          "401": {
            "description": "Unauthorized"
          },
          "403": {
            "description": "Forbidden - Admin role required"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Create category (Admin only)",
        "tags": [
          "Admin - Categories"
        ]
      }
    },
    "/api/admin/categories/{id}": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Category deleted"
          },
          "401": {
            "description": "Unauthorized"
          },
          "403": {
            "description": "Forbidden - Admin role required"
          },
          "404": {
            "description": "Category not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Delete category (Admin only)",
        "tags": [
          "Admin - Categories"
        ]
      },
      "put": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "description": {
                    "example": "Complete web development course",
                    "type": "string"
                  },
                  "name": {
                    "example": "Web Development",
                    "type": "string"
                  },
                  "slug": {
                    "example": "web-dev",
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Category"
                }
              }
            },
            "description": "Category updated"
          },
          "401": {
            "description": "Unauthorized"
          },
          "403": {
            "description": "Forbidden - Admin role required"
          },
          "404": {
            "description": "Category not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Update category (Admin only)",
        "tags": [
          "Admin - Categories"
        ]
      }
    },
    "/api/admin/classes": {
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "category_id": {
                    "example": "cmq6b6ehc000010uj9ewttltv",
                    "type": "string"
                  },
                  "content": {
                    "type": "string"
                  },
                  "description": {
                    "type": "string"
                  },
                  "difficulty": {
                    "enum": [
                      "beginner",
                      "intermediate",
                      "advanced"
                    ],
                    "type": "string"
                  },
                  "duration": {
                    "type": "integer"
                  },
                  "project_id": {
                    "example": "cmq6b6ehc000010uj9ewttltv",
                    "type": "string"
                  },
                  "thumbnail": {
                    "format": "url",
                    "type": "string"
                  },
                  "title": {
                    "type": "string"
                  }
                },
                "required": [
                  "title",
                  "project_id",
                  "category_id"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "description": "Class created"
          },
          "403": {
            "description": "Forbidden"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Create class",
        "tags": [
          "Admin - Classes"
        ]
      }
    },
    "/api/admin/classes/{id}": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Class deleted"
          },
          "403": {
            "description": "Forbidden - Admin only"
          },
          "404": {
            "description": "Class not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Delete class",
        "tags": [
          "Admin - Classes"
        ]
      },
      "put": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "content": {
                    "type": "string"
                  },
                  "description": {
                    "type": "string"
                  },
                  "difficulty": {
                    "enum": [
                      "beginner",
                      "intermediate",
                      "advanced"
                    ],
                    "type": "string"
                  },
                  "duration": {
                    "type": "integer"
                  },
                  "title": {
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "Class updated"
          },
          "403": {
            "description": "Forbidden - Admin only"
          },
          "404": {
            "description": "Class not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Update class",
        "tags": [
          "Admin - Classes"
        ]
      }
    },
    "/api/admin/dashboard": {
      "get": {
        "responses": {
          "200": {
            "description": "Welcome to admin dashboard"
          },
          "403": {
            "description": "Forbidden - Admin only"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Admin dashboard",
        "tags": [
          "Admin"
        ]
      }
    },
    "/api/admin/projects": {
      "get": {
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "schema": {
              "type": "integer"
            }
          },
          {
            "in": "query",
            "name": "limit",
            "schema": {
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "properties": {
                    "data": {
                      "items": {
                        "$ref": "#/components/schemas/Project"
                      },
                      "type": "array"
                    },
                    "pagination": {
                      "$ref": "#/components/schemas/Pagination"
                    }
                  },
                  "type": "object"
                }
              }
            },
            "description": "Projects list"
          },
          "403": {
            "description": "Forbidden"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "List all projects (Admin view)",
        "tags": [
          "Admin - Projects"
        ]
      },
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "category_id": {
                    "example": "cmq6b6ehc000010uj9ewttltv",
                    "type": "string"
                  },
                  "description": {
                    "type": "string"
                  },
                  "thumbnail": {
                    "format": "url",
                    "type": "string"
                  },
                  "title": {
                    "type": "string"
                  },
                  "visibility": {
                    "enum": [
                      "public",
                      "private"
                    ],
                    "type": "string"
                  }
                },
                "required": [
                  "title",
                  "category_id"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Project"
                }
              }
            },
            "description": "Project created"
          },
          "403": {
            "description": "Forbidden - Admin only"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Create project (Admin only)",
        "tags": [
          "Admin - Projects"
        ]
      }
    },
    "/api/admin/projects/{id}": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Project deleted"
          },
          "403": {
            "description": "Forbidden"
          },
          "404": {
            "description": "Project not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Delete project",
        "tags": [
          "Admin - Projects"
        ]
      },
      "put": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "description": {
                    "type": "string"
                  },
                  "title": {
                    "type": "string"
                  },
                  "visibility": {
                    "enum": [
                      "public",
                      "private"
                    ],
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "Project updated"
          },
          "403": {
            "description": "Forbidden"
          },
          "404": {
            "description": "Project not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Update project",
        "tags": [
          "Admin - Projects"
        ]
      }
    },
    "/api/admin/users": {
      "get": {
        "parameters": [
          {
            "description": "Page number",
            "in": "query",
            "name": "page",
            "schema": {
              "default": 1,
              "type": "integer"
            }
          },
          {
            "description": "Items per page (max: 100)",
            "in": "query",
            "name": "limit",
            "schema": {
              "default": 20,
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "properties": {
                    "data": {
                      "items": {
                        "properties": {
                          "avatar": {
                            "format": "url",
                            "type": "string"
                          },
                          "created_at": {
                            "format": "date-time",
                            "type": "string"
                          },
                          "email": {
                            "format": "email",
                            "type": "string"
                          },
                          "id": {
                            "example": "cmq6b6ehc000010uj9ewttltv",
                            "type": "string"
                          },
                          "is_active": {
                            "type": "boolean"
                          },
                          "level": {
                            "maximum": 5,
                            "minimum": 1,
                            "type": "integer"
                          },
                          "name": {
                            "type": "string"
                          },
                          "role": {
                            "enum": [
                              "admin",
                              "user"
                            ],
                            "type": "string"
                          },
                          "total_points": {
                            "type": "integer"
                          }
                        },
                        "type": "object"
                      },
                      "type": "array"
                    },
                    "pagination": {
                      "$ref": "#/components/schemas/Pagination"
                    }
                  },
                  "type": "object"
                }
              }
            },
            "description": "Users list with pagination"
          },
          "401": {
            "description": "Unauthorized - Missing or invalid token"
          },
          "403": {
            "description": "Forbidden - Admin role required"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "List all users with pagination (Admin only)",
        "tags": [
          "Admin - Users"
        ]
      }
    },
    "/api/auth/login": {
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "email": {
                    "format": "email",
                    "type": "string"
                  },
                  "password": {
                    "type": "string"
                  }
                },
                "required": [
                  "email",
                  "password"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "properties": {
                    "access_token": {
                      "type": "string"
                    },
                    "refresh_token": {
                      "type": "string"
                    },
                    "user": {
                      "$ref": "#/components/schemas/User"
                    }
                  },
                  "type": "object"
                }
              }
            },
            "description": "Login successful"
          },
          "401": {
            "description": "Invalid credentials"
          }
        },
        "summary": "Login user",
        "tags": [
          "Authentication"
        ]
      }
    },
    "/api/auth/register": {
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "email": {
                    "format": "email",
                    "type": "string"
                  },
                  "name": {
                    "type": "string"
                  },
                  "password": {
                    "minLength": 6,
                    "type": "string"
                  }
                },
                "required": [
                  "email",
                  "password",
                  "name"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "User created successfully"
          },
          "400": {
            "description": "Invalid input"
          },
          "409": {
            "description": "Email already exists"
          }
        },
        "summary": "Register new user",
        "tags": [
          "Authentication"
        ]
      }
    },
    "/api/categories": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/Category"
                  },
                  "type": "array"
                }
              }
            },
            "description": "List of categories"
          }
        },
        "summary": "List all categories",
        "tags": [
          "Categories"
        ]
      }
    },
    "/api/categories/{category_id}/projects": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "category_id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "items": {
                    "$ref": "#/components/schemas/Project"
                  },
                  "type": "array"
                }
              }
            },
            "description": "List of projects in category"
          },
          "404": {
            "description": "Category not found"
          }
        },
        "summary": "List projects in a category",
        "tags": [
          "Categories"
        ]
      }
    },
    "/api/categories/{slug}": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "slug",
            "required": true,
            "schema": {
              "example": "web-dev",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Category"
                }
              }
            },
            "description": "Category details"
          },
          "404": {
            "description": "Category not found"
          }
        },
        "summary": "Get category detail by slug",
        "tags": [
          "Categories"
        ]
      }
    },
    "/api/certificates": {
      "get": {
        "responses": {
          "200": {
            "description": "User certificates"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Get user certificates",
        "tags": [
          "Certificates"
        ]
      }
    },
    "/api/certificates/{code}": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "code",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Certificate details"
          },
          "403": {
            "description": "Forbidden - Not your certificate"
          },
          "404": {
            "description": "Certificate not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "View certificate (Private - Owner only)",
        "tags": [
          "Certificates"
        ]
      }
    },
    "/api/classes/{id}": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Class details"
          },
          "404": {
            "description": "Class not found"
          }
        },
        "summary": "Get class details",
        "tags": [
          "Classes"
        ]
      }
    },
    "/api/classes/{id}/complete": {
      "post": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Class marked as completed"
          },
          "401": {
            "description": "Unauthorized"
          },
          "404": {
            "description": "Class not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Mark class as completed",
        "tags": [
          "Classes"
        ]
      }
    },
    "/api/discussions": {
      "get": {
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "schema": {
              "type": "integer"
            }
          },
          {
            "in": "query",
            "name": "class_id",
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          },
          {
            "in": "query",
            "name": "category_id",
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Discussions list"
          }
        },
        "summary": "List discussions",
        "tags": [
          "Discussions"
        ]
      },
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "category_id": {
                    "example": "cmq6b6ehc000010uj9ewttltv",
                    "type": "string"
                  },
                  "class_id": {
                    "example": "cmq6b6ehc000010uj9ewttltv",
                    "type": "string"
                  },
                  "content": {
                    "type": "string"
                  },
                  "title": {
                    "type": "string"
                  }
                },
                "required": [
                  "title",
                  "content"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "description": "Discussion created"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Create discussion",
        "tags": [
          "Discussions"
        ]
      }
    },
    "/api/discussions/{id}": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Discussion deleted"
          },
          "403": {
            "description": "Forbidden"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Delete discussion",
        "tags": [
          "Discussions"
        ]
      },
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Discussion details with replies"
          },
          "404": {
            "description": "Discussion not found"
          }
        },
        "summary": "Get discussion with threaded replies",
        "tags": [
          "Discussions"
        ]
      },
      "put": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "content": {
                    "type": "string"
                  },
                  "title": {
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "Discussion updated"
          },
          "403": {
            "description": "Forbidden"
          },
          "404": {
            "description": "Not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Update discussion",
        "tags": [
          "Discussions"
        ]
      }
    },
    "/api/discussions/{id}/replies": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Replies list"
          }
        },
        "summary": "Get discussion replies",
        "tags": [
          "Replies"
        ]
      },
      "post": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "content": {
                    "type": "string"
                  }
                },
                "required": [
                  "content"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "description": "Reply created"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Reply to discussion",
        "tags": [
          "Replies"
        ]
      }
    },
    "/api/health": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "properties": {
                    "status": {
                      "example": "ok",
                      "type": "string"
                    }
                  },
                  "type": "object"
                }
              }
            },
            "description": "Server is healthy"
          }
        },
        "summary": "Health check endpoint",
        "tags": [
          "Health"
        ]
      }
    },
    "/api/leaderboard": {
      "get": {
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "schema": {
              "default": 1,
              "type": "integer"
            }
          },
          {
            "in": "query",
            "name": "limit",
            "schema": {
              "default": 50,
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "properties": {
                    "data": {
                      "items": {
                        "$ref": "#/components/schemas/LeaderboardEntry"
                      },
                      "type": "array"
                    },
                    "pagination": {
                      "$ref": "#/components/schemas/Pagination"
                    }
                  },
                  "type": "object"
                }
              }
            },
            "description": "Leaderboard"
          }
        },
        "summary": "Get top users by points",
        "tags": [
          "Gamification"
        ]
      }
    },
    "/api/levels": {
      "get": {
        "responses": {
          "200": {
            "description": "Level information"
          }
        },
        "summary": "Get level progression info",
        "tags": [
          "Gamification"
        ]
      }
    },
    "/api/projects/{id}/classes": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          },
          {
            "in": "query",
            "name": "page",
            "schema": {
              "default": 1,
              "type": "integer"
            }
          },
          {
            "in": "query",
            "name": "limit",
            "schema": {
              "default": 20,
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "properties": {
                    "data": {
                      "items": {
                        "$ref": "#/components/schemas/Class"
                      },
                      "type": "array"
                    },
                    "pagination": {
                      "$ref": "#/components/schemas/Pagination"
                    }
                  },
                  "type": "object"
                }
              }
            },
            "description": "Classes list"
          },
          "404": {
            "description": "Project not found"
          }
        },
        "summary": "List classes in a project",
        "tags": [
          "Classes"
        ]
      }
    },
    "/api/replies/{id}": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Reply deleted"
          },
          "401": {
            "description": "Unauthorized"
          },
          "403": {
            "description": "Forbidden"
          },
          "404": {
            "description": "Reply not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Delete reply",
        "tags": [
          "Replies"
        ]
      },
      "put": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "content": {
                    "type": "string"
                  }
                },
                "required": [
                  "content"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "Reply updated"
          },
          "401": {
            "description": "Unauthorized"
          },
          "403": {
            "description": "Forbidden"
          },
          "404": {
            "description": "Reply not found"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Update reply",
        "tags": [
          "Replies"
        ]
      }
    },
    "/api/replies/{id}/replies": {
      "post": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "content": {
                    "type": "string"
                  }
                },
                "required": [
                  "content"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "description": "Nested reply created"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Reply to a reply (nested)",
        "tags": [
          "Replies"
        ]
      }
    },
    "/api/showcases": {
      "get": {
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "schema": {
              "type": "integer"
            }
          },
          {
            "in": "query",
            "name": "category_id",
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          },
          {
            "in": "query",
            "name": "sort",
            "schema": {
              "enum": [
                "newest",
                "trending",
                "most_liked"
              ],
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Showcases list"
          }
        },
        "summary": "List showcases",
        "tags": [
          "Showcases"
        ]
      },
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "category_id": {
                    "example": "cmq6b6ehc000010uj9ewttltv",
                    "type": "string"
                  },
                  "description": {
                    "type": "string"
                  },
                  "media_urls": {
                    "items": {
                      "format": "url",
                      "type": "string"
                    },
                    "type": "array"
                  },
                  "title": {
                    "type": "string"
                  },
                  "visibility": {
                    "enum": [
                      "public",
                      "private",
                      "friends_only"
                    ],
                    "type": "string"
                  }
                },
                "required": [
                  "title",
                  "category_id"
                ],
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "201": {
            "description": "Showcase created"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Create showcase",
        "tags": [
          "Showcases"
        ]
      }
    },
    "/api/showcases/{id}": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Showcase deleted"
          },
          "403": {
            "description": "Forbidden"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Delete showcase",
        "tags": [
          "Showcases"
        ]
      },
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Showcase details"
          },
          "404": {
            "description": "Showcase not found"
          }
        },
        "summary": "Get showcase details",
        "tags": [
          "Showcases"
        ]
      },
      "put": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "description": {
                    "type": "string"
                  },
                  "title": {
                    "type": "string"
                  },
                  "visibility": {
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "description": "Showcase updated"
          },
          "403": {
            "description": "Forbidden"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Update showcase (Owner only)",
        "tags": [
          "Showcases"
        ]
      }
    },
    "/api/showcases/{id}/like": {
      "delete": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "Like removed"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Unlike showcase",
        "tags": [
          "Showcases"
        ]
      },
      "post": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Showcase liked"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Like showcase",
        "tags": [
          "Showcases"
        ]
      }
    },
    "/api/users/me": {
      "get": {
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "User profile"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Get current user profile",
        "tags": [
          "Users"
        ]
      },
      "put": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "avatar": {
                    "format": "url",
                    "type": "string"
                  },
                  "bio": {
                    "type": "string"
                  },
                  "name": {
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "required": true
        },
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            },
            "description": "Profile updated"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Update user profile",
        "tags": [
          "Users"
        ]
      }
    },
    "/api/users/me/activity": {
      "get": {
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "schema": {
              "type": "integer"
            }
          },
          {
            "in": "query",
            "name": "limit",
            "schema": {
              "type": "integer"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Activity log"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "summary": "Get user activity log",
        "tags": [
          "Gamification"
        ]
      }
    },
    "/api/users/{id}": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/UserPublic"
                }
              }
            },
            "description": "User profile"
          },
          "404": {
            "description": "User not found"
          }
        },
        "summary": "Get public user profile",
        "tags": [
          "Users"
        ]
      }
    },
    "/api/users/{id}/points": {
      "get": {
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "schema": {
              "example": "cmq6b6ehc000010uj9ewttltv",
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "User points info"
          }
        },
        "summary": "Get user points and level",
        "tags": [
          "Gamification"
        ]
      }
    }
  },
  "servers": [
    {
      "description": "Local development",
      "url": "http://localhost:3000"
    },
    {
      "description": "Production",
      "url": "https://jvalleyverse.mohagussetiaone.my.id"
    }
  ],
  "tags": [
    {
      "description": "User authentication endpoints",
      "name": "Authentication"
    },
    {
      "description": "User profile management",
      "name": "Users"
    },
    {
      "description": "Admin dashboard and management",
      "name": "Admin"
    },
    {
      "description": "Admin project management",
      "name": "Admin - Projects"
    },
    {
      "description": "Admin class management",
      "name": "Admin - Classes"
    },
    {
      "description": "Public class access",
      "name": "Classes"
    },
    {
      "description": "Certificate endpoints (private access)",
      "name": "Certificates"
    },
    {
      "description": "Discussion management",
      "name": "Discussions"
    },
    {
      "description": "Discussion replies (nested)",
      "name": "Replies"
    },
    {
      "description": "User portfolio items",
      "name": "Showcases"
    },
    {
      "description": "Points, levels, leaderboard",
      "name": "Gamification"
    },
    {
      "description": "Server health check",
      "name": "Health"
    },
    {
      "description": "Category list and browsing",
      "name": "Categories"
    },
    {
      "description": "Admin category management",
      "name": "Admin - Categories"
    }
  ]
}`
