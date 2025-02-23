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
        "/secret/{key}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Gets a secret by key path",
                "tags": [
                    "secret"
                ],
                "summary": "Retrieve a secret",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Secret key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Secret retrieved",
                        "schema": {
                            "$ref": "#/definitions/models.GetSecretResponse"
                        }
                    },
                    "400": {
                        "description": "Missing key path",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Stores a secret with a key path",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "secret"
                ],
                "summary": "Store a secret",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Secret key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Secret data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.StoreSecretRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Secret stored",
                        "schema": {
                            "$ref": "#/definitions/models.StoreSecretResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Deletes a secret by key path",
                "tags": [
                    "secret"
                ],
                "summary": "Delete a secret",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Secret key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Missing key path",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/token": {
            "get": {
                "description": "Generates a short-lived token for secret operations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token"
                ],
                "summary": "Generate a token",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Token TTL in seconds",
                        "name": "ttl",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Generated token",
                        "schema": {
                            "$ref": "#/definitions/models.IssueTokenResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid TTL",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Deletes all stored secrets for the authenticated token",
                "tags": [
                    "token"
                ],
                "summary": "Delete all secrets for a token",
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/token/valid": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Validates if a token is still active",
                "tags": [
                    "token"
                ],
                "summary": "Validate a token",
                "responses": {
                    "200": {
                        "description": "Token is valid",
                        "schema": {
                            "$ref": "#/definitions/models.TokenValidationResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid or expired token",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ErrorResponse": {
            "description": "API error response format",
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "request_id": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "models.GetSecretResponse": {
            "description": "Get secret response format",
            "type": "object",
            "properties": {
                "ttl": {
                    "type": "integer"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "models.IssueTokenResponse": {
            "description": "Issue token response format",
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "ttl": {
                    "type": "integer"
                }
            }
        },
        "models.StoreSecretRequest": {
            "description": "Store secret request format",
            "type": "object",
            "required": [
                "value"
            ],
            "properties": {
                "value": {
                    "type": "string"
                }
            }
        },
        "models.StoreSecretResponse": {
            "description": "Store secret response format",
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "ttl": {
                    "type": "integer"
                }
            }
        },
        "models.TokenValidationResponse": {
            "description": "Token validation response format",
            "type": "object",
            "properties": {
                "valid": {
                    "type": "boolean"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type \"Bearer {your_token}\" into the field below",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "localhost:8888",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Go Secrets API",
	Description:      "A simple API for managing secrets using Redis.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
