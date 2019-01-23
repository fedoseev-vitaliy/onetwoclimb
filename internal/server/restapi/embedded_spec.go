// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "OneTwoClimb API",
    "version": "1.0"
  },
  "basePath": "/api/v1.0",
  "paths": {
    "/colors": {
      "get": {
        "summary": "get board colors",
        "operationId": "boardColors",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "colors"
              ],
              "properties": {
                "colors": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Color"
                  }
                }
              }
            }
          },
          "500": {
            "description": "General server error. Error codes:\n  - 4 Server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Color": {
      "description": "color object",
      "type": "object",
      "properties": {
        "hex": {
          "description": "color in hex",
          "type": "string"
        },
        "id": {
          "description": "color id",
          "type": "integer"
        },
        "name": {
          "description": "item name",
          "type": "string"
        },
        "pinCode": {
          "description": "pin code name",
          "type": "string"
        }
      }
    },
    "Error": {
      "type": "object",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "description": "internal status code",
          "type": "integer"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "OneTwoClimb API",
    "version": "1.0"
  },
  "basePath": "/api/v1.0",
  "paths": {
    "/colors": {
      "get": {
        "summary": "get board colors",
        "operationId": "boardColors",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "colors"
              ],
              "properties": {
                "colors": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Color"
                  }
                }
              }
            }
          },
          "500": {
            "description": "General server error. Error codes:\n  - 4 Server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Color": {
      "description": "color object",
      "type": "object",
      "properties": {
        "hex": {
          "description": "color in hex",
          "type": "string"
        },
        "id": {
          "description": "color id",
          "type": "integer"
        },
        "name": {
          "description": "item name",
          "type": "string"
        },
        "pinCode": {
          "description": "pin code name",
          "type": "string"
        }
      }
    },
    "Error": {
      "type": "object",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "description": "internal status code",
          "type": "integer"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
}
