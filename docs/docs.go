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
        "/claims": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Receives claim data and processes it",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "claims"
                ],
                "summary": "Submit a new claim",
                "parameters": [
                    {
                        "description": "Claim data to submit",
                        "name": "claim",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ClaimSubmissionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Claim submitted successfully",
                        "schema": {
                            "$ref": "#/definitions/models.Claim"
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/claims/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Returns the details of a specific claim by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "claims"
                ],
                "summary": "Get claim by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Claim ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Claim details",
                        "schema": {
                            "$ref": "#/definitions/models.Claim"
                        }
                    },
                    "400": {
                        "description": "Claim ID not provided"
                    },
                    "404": {
                        "description": "Claim not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Returns an \"ok\" status if the application is running.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Checks application health",
                "responses": {
                    "200": {
                        "description": "Status OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/reversal": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Reverts an already submitted claim and records the reversal",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "claims"
                ],
                "summary": "Reverse an existing claim",
                "parameters": [
                    {
                        "description": "Claim ID to be reverted",
                        "name": "reversal",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ClaimReversalRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reversal successfully recorded",
                        "schema": {
                            "$ref": "#/definitions/models.ClaimReversalResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request or claim already reverted/not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Claim": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "Unique ID of the claim (UUID)",
                    "type": "string"
                },
                "ndc": {
                    "description": "National Drug Code of the medication",
                    "type": "string"
                },
                "npi": {
                    "description": "National Provider Identifier of the pharmacy",
                    "type": "string"
                },
                "price": {
                    "description": "Price of the medication",
                    "type": "number"
                },
                "quantity": {
                    "description": "Quantity of the medication",
                    "type": "number"
                },
                "reverted": {
                    "description": "Indicates if the claim has been reverted",
                    "type": "boolean"
                },
                "timestamp": {
                    "description": "Date and time of claim submission",
                    "type": "string"
                }
            }
        },
        "models.ClaimReversalRequest": {
            "type": "object",
            "properties": {
                "claim_id": {
                    "description": "ID of the claim to be reverted",
                    "type": "string"
                }
            }
        },
        "models.ClaimReversalResponse": {
            "type": "object",
            "properties": {
                "claim_id": {
                    "description": "ID of the reverted claim",
                    "type": "string"
                },
                "status": {
                    "description": "Operation status (e.g., \"claim reversed\")",
                    "type": "string"
                }
            }
        },
        "models.ClaimSubmissionRequest": {
            "type": "object",
            "properties": {
                "ndc": {
                    "description": "National Drug Code of the medication",
                    "type": "string"
                },
                "npi": {
                    "description": "National Provider Identifier of the pharmacy",
                    "type": "string"
                },
                "price": {
                    "description": "Price of the medication",
                    "type": "number"
                },
                "quantity": {
                    "description": "Quantity of the medication",
                    "type": "number"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Pharmacy Claim Service API",
	Description:      "This is a sample service for managing pharmacy claims.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
