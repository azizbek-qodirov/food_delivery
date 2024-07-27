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
        "/change-role/{id}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Changes the role of a user or admin. Only admins are allowed to use this function.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin-panel"
                ],
                "summary": "Change a user's role",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id or email of the user",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "id",
                            "email"
                        ],
                        "type": "string",
                        "description": "Search with",
                        "name": "data",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "admin",
                            "user"
                        ],
                        "type": "string",
                        "description": "New role of the user",
                        "name": "role",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User role updated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/confirm-registration": {
            "post": {
                "description": "Confirms a user's registration using the code sent to their email.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "registration"
                ],
                "summary": "Confirm registration with code",
                "parameters": [
                    {
                        "description": "Confirmation request",
                        "name": "confirmation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ConfirmRegistrationReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT tokens",
                        "schema": {
                            "$ref": "#/definitions/token.Tokens"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Incorrect verification code",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Verification code expired or email not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/forgot-password": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Sends a confirmation code to email recovery password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "password_recovery"
                ],
                "summary": "Forgot passwrod",
                "parameters": [
                    {
                        "description": "User login credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ForgotPasswordReq"
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
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Page not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Authenticate user with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "login"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User login credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT tokens",
                        "schema": {
                            "$ref": "#/definitions/token.Tokens"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid email or password",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get the profile of the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.GetProfileResp"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/recover-password": {
            "post": {
                "description": "Verifies the code and updates the password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "password_recovery"
                ],
                "summary": "Recover password (Use this one after sending verification code)",
                "parameters": [
                    {
                        "description": "Recover Password Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RecoverPasswordReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password successfully updated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Incorrect verification code",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Verification code expired or email not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error updating password",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Register a new user with email, username, and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "registration"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User registration request",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RegisterReqSwag"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "JWT tokens",
                        "schema": {
                            "$ref": "#/definitions/token.Tokens"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ConfirmRegistrationReq": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Confirmation code received via email",
                    "type": "string"
                },
                "email": {
                    "description": "User's email address",
                    "type": "string"
                }
            }
        },
        "models.ForgotPasswordReq": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "User's email address",
                    "type": "string"
                }
            }
        },
        "models.GetProfileResp": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "User's email address",
                    "type": "string"
                },
                "id": {
                    "description": "User's unique identifier",
                    "type": "string"
                },
                "is_confirmed": {
                    "description": "Add IsConfirmed to the model",
                    "type": "boolean"
                },
                "password": {
                    "description": "User's password",
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "models.LoginReq": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "User's email",
                    "type": "string"
                },
                "password": {
                    "description": "User's password",
                    "type": "string"
                }
            }
        },
        "models.RecoverPasswordReq": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "new_password": {
                    "type": "string"
                }
            }
        },
        "models.RegisterReqSwag": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "User's email address",
                    "type": "string"
                },
                "password": {
                    "description": "User's password",
                    "type": "string"
                }
            }
        },
        "token.Tokens": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
