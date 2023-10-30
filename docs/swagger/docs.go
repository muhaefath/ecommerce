// Package swagger GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package swagger

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
        "/product": {
            "post": {
                "description": "create a product",
                "tags": [
                    "Product"
                ],
                "summary": "create a product",
                "operationId": "v1-CreateProduct",
                "parameters": [
                    {
                        "description": "UpsertProduct",
                        "name": "UpsertProduct",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpsertProduct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/product/list": {
            "get": {
                "operationId": "v1-GetProductList",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.GetProductListResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/product/review": {
            "post": {
                "description": "create a product review",
                "tags": [
                    "Product"
                ],
                "summary": "create a product review",
                "operationId": "v1-CreateProductReview",
                "parameters": [
                    {
                        "description": "UpsertProductReview",
                        "name": "UpsertProductReview",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpsertProductReview"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/product/{product_id}": {
            "get": {
                "description": "get a product",
                "tags": [
                    "Product"
                ],
                "summary": "get a product",
                "operationId": "v1-GetDetailProduct",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            },
            "put": {
                "description": "update a product",
                "tags": [
                    "Product"
                ],
                "summary": "update a product",
                "operationId": "v1-UpdateProduct",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "UpsertProduct",
                        "name": "UpsertProduct",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpsertProduct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Product": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "etalase": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "rating": {
                    "type": "number"
                },
                "sku": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                },
                "weight": {
                    "type": "number"
                }
            }
        },
        "request.UpsertProduct": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "etalase": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "product_images": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "image_url": {
                                "type": "string"
                            },
                            "short_description": {
                                "type": "string"
                            }
                        }
                    }
                },
                "sku": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                },
                "weight": {
                    "type": "number"
                }
            }
        },
        "request.UpsertProductReview": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string"
                },
                "product_id": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                }
            }
        },
        "response.BaseResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "response.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "strconv.ParseInt: parsing \"a\": invalid syntax"
                },
                "status_code": {
                    "type": "integer",
                    "example": 400
                }
            }
        },
        "response.GetProductListResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Product"
                    }
                },
                "message": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
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
