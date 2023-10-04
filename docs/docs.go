// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Alex Chinaev",
            "url": "https://vk.com/l.chinaev",
            "email": "ax.chinaev@yandex.ru"
        },
        "license": {
            "name": "AS IS (NO WARRANTY)"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/auth/login": {
            "post": {
                "description": "create user session and put it into cookie",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "login user",
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "json"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/logout": {
            "post": {
                "description": "delete current session and nullify cookie",
                "tags": [
                    "auth"
                ],
                "summary": "logout user",
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "json"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/register": {
            "post": {
                "description": "add new user to db and return it id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "register user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "json"
                        }
                    }
                }
            }
        },
        "/api/v1/films/genre/{genre}": {
            "get": {
                "description": "Get a list of films based on the specified genre.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "films"
                ],
                "summary": "Get films by genre",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The genre of the films you want to retrieve.",
                        "name": "genre",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Film"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "json"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "json"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Film": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "previewPath": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1E",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Netfilx API",
	Description:      "API of the nelfix project by holi",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
