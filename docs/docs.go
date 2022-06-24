// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Rinat Almakhov",
            "url": "https://github.com/Gasoid/"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/albums/{sourceName}/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "returns albums",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Albums",
                "parameters": [
                    {
                        "type": "string",
                        "description": "source name",
                        "name": "sourceName",
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
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/download-album/{albumID}/{sourceName}/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "download all photos of particular album, returns destination of your photos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "download photos of album",
                "parameters": [
                    {
                        "type": "string",
                        "description": "source name",
                        "name": "sourceName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "album ID",
                        "name": "albumID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "directory where photos will be stored",
                        "name": "dir",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/download-all-albums/{sourceName}/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "download all photos of all albums, returns destination of your photos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "download photos of albums",
                "parameters": [
                    {
                        "type": "string",
                        "description": "source name",
                        "name": "sourceName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "directory where photos will be stored",
                        "name": "dir",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/sources/": {
            "get": {
                "description": "returns sources",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Sources",
                "responses": {
                    "200": {
                        "description": "sources",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "api_key",
            "in": "query"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "PhotoDumper",
	Description:      "app downloads photos from vk.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
