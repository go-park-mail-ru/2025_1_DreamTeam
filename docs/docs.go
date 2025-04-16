// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/getCourseLesson": {
            "get": {
                "description": "Returns the lesson the user should take in the course",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courses"
                ],
                "summary": "Get lesson of a course for the user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Course ID",
                        "name": "courseId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "next lesson of the course",
                        "schema": {
                            "$ref": "#/definitions/response.LessonResponse"
                        }
                    },
                    "400": {
                        "description": "invalid course ID",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "not authorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "method not allowed",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/getCourseRoadmap": {
            "get": {
                "description": "Returns the roadmap of a course for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courses"
                ],
                "summary": "Get course roadmap",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Course ID",
                        "name": "courseId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Course roadmap",
                        "schema": {
                            "$ref": "#/definitions/response.CourseRoadmapResponse"
                        }
                    },
                    "400": {
                        "description": "invalid course ID",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "method not allowed",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/getCourses": {
            "get": {
                "description": "Retrieves a course by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courses"
                ],
                "summary": "Get course",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Course ID",
                        "name": "courseId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "course",
                        "schema": {
                            "$ref": "#/definitions/response.CourseResponse"
                        }
                    },
                    "405": {
                        "description": "method not allowed",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/getNextLesson": {
            "get": {
                "description": "Returns the next lesson the user should take based on current lesson and course",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courses"
                ],
                "summary": "Get next lesson in a course",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Course ID",
                        "name": "courseId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Current Lesson ID",
                        "name": "lessonId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "next lesson content",
                        "schema": {
                            "$ref": "#/definitions/response.LessonResponse"
                        }
                    },
                    "400": {
                        "description": "invalid course or lesson ID",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "not authorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "method not allowed",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/isAuthorized": {
            "get": {
                "description": "Returns user profile if authorized",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Check if user is authorized",
                "responses": {
                    "200": {
                        "description": "User profile",
                        "schema": {
                            "$ref": "#/definitions/response.UserProfileResponse"
                        }
                    },
                    "401": {
                        "description": "not authorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "Login user with the given email, password and send cookie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "User information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UserDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "200 OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "invalid request | password too short | invalid email",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "email or password incorrect",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "method not allowed",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/logout": {
            "post": {
                "description": "Logout user by deleting session cookie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Logout user",
                "responses": {
                    "200": {
                        "description": "200 OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/markLessonAsNotCompleted": {
            "post": {
                "description": "Marks the specified lesson as not completed for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courses"
                ],
                "summary": "Mark a lesson as not completed",
                "parameters": [
                    {
                        "description": "Lesson ID",
                        "name": "lessonId",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LessonIDRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "ivalid lesson ID",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "uethod not allowed",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "Register a new user with the given name, email, and password and send cookie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "200 OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "invalid request | missing required fields | password too short | invalid email",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "email exists",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "method not allowed",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/updateProfile": {
            "post": {
                "description": "Updates the profile photo of the authorized user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update user profile photo",
                "parameters": [
                    {
                        "description": "Updated user profile",
                        "name": "profile",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserProfile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "200 OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "invalid request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "not authorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CourseDTO": {
            "type": "object",
            "properties": {
                "creator_id": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "purchases_amount": {
                    "type": "integer"
                },
                "rating": {
                    "type": "number"
                },
                "src_image": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "time_to_pass": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "dto.CoursePartDTO": {
            "type": "object",
            "properties": {
                "buckets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.LessonBucketDTO"
                    }
                },
                "part_id": {
                    "type": "integer"
                },
                "part_title": {
                    "type": "string"
                }
            }
        },
        "dto.CourseRoadmapDTO": {
            "type": "object",
            "properties": {
                "parts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.CoursePartDTO"
                    }
                }
            }
        },
        "dto.LessonBucketDTO": {
            "type": "object",
            "properties": {
                "bucket_id": {
                    "type": "integer"
                },
                "bucket_title": {
                    "type": "string"
                },
                "lessons": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.LessonPointDTO"
                    }
                }
            }
        },
        "dto.LessonDTO": {
            "type": "object",
            "properties": {
                "header": {
                    "$ref": "#/definitions/dto.LessonDtoHeader"
                },
                "lesson_body": {
                    "$ref": "#/definitions/dto.LessonDtoBody"
                }
            }
        },
        "dto.LessonDtoBody": {
            "type": "object",
            "properties": {
                "blocks": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "body": {
                                "type": "string"
                            }
                        }
                    }
                },
                "footer": {
                    "type": "object",
                    "properties": {
                        "next_lesson_id": {
                            "type": "integer"
                        },
                        "previous_lesson_id": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "dto.LessonDtoHeader": {
            "type": "object",
            "properties": {
                "bucket": {
                    "type": "object",
                    "properties": {
                        "order": {
                            "type": "integer"
                        },
                        "title": {
                            "type": "string"
                        }
                    }
                },
                "course_id": {
                    "type": "integer"
                },
                "course_title": {
                    "type": "string"
                },
                "part": {
                    "type": "object",
                    "properties": {
                        "order": {
                            "type": "integer"
                        },
                        "title": {
                            "type": "string"
                        }
                    }
                },
                "points": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "is_done": {
                                "type": "boolean"
                            },
                            "lesson_id": {
                                "type": "integer"
                            },
                            "type": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "dto.LessonIDRequest": {
            "type": "object",
            "properties": {
                "lesson_id": {
                    "type": "integer"
                }
            }
        },
        "dto.LessonPointDTO": {
            "type": "object",
            "properties": {
                "is_done": {
                    "type": "boolean"
                },
                "lesson_id": {
                    "type": "integer"
                },
                "lesson_title": {
                    "type": "string"
                },
                "lesson_type": {
                    "type": "string"
                }
            }
        },
        "dto.UserDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "dto.UserProfileDTO": {
            "type": "object",
            "properties": {
                "avatar_src": {
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "hide_email": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "salt": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "models.UserProfile": {
            "type": "object",
            "properties": {
                "avatarSrc": {
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "hideEmail": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "response.BucketCoursesResponse": {
            "type": "object",
            "properties": {
                "bucket_courses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.CourseDTO"
                    }
                }
            }
        },
        "response.CourseResponse": {
            "type": "object",
            "properties": {
                "course": {
                    "$ref": "#/definitions/dto.CourseDTO"
                }
            }
        },
        "response.CourseRoadmapResponse": {
            "type": "object",
            "properties": {
                "course_roadmap": {
                    "$ref": "#/definitions/dto.CourseRoadmapDTO"
                }
            }
        },
        "response.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "response.LessonResponse": {
            "type": "object",
            "properties": {
                "lesson": {
                    "$ref": "#/definitions/dto.LessonDTO"
                }
            }
        },
        "response.UserProfileResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "$ref": "#/definitions/dto.UserProfileDTO"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
