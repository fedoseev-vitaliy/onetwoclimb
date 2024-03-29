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
        "produces": [
          "application/json"
        ],
        "summary": "get board colors",
        "operationId": "getBoardColors",
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
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "summary": "post board colors",
        "operationId": "postBoardColors",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "type": "object",
              "$ref": "#/definitions/Color"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "500": {
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/colors/{colorId}": {
      "delete": {
        "summary": "delete borad color",
        "operationId": "delBoardColor",
        "parameters": [
          {
            "$ref": "#/parameters/colorId"
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "500": {
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/download": {
      "get": {
        "produces": [
          "image/png,image/jpeg"
        ],
        "summary": "Download image",
        "operationId": "downloadFile",
        "parameters": [
          {
            "$ref": "#/parameters/id"
          }
        ],
        "responses": {
          "200": {
            "description": "download file OK",
            "schema": {
              "type": "string",
              "format": "byte"
            }
          },
          "400": {
            "description": "Bad Argument",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "404": {
            "description": "File not found",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/upload": {
      "post": {
        "description": "The file can't be larger than **5MB**",
        "consumes": [
          "multipart/form-data"
        ],
        "summary": "Upload image",
        "operationId": "uploadFile",
        "parameters": [
          {
            "type": "file",
            "description": "The file to upload",
            "name": "file",
            "in": "formData"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "id"
              ],
              "properties": {
                "id": {
                  "description": "image id",
                  "type": "string"
                }
              }
            }
          },
          "500": {
            "description": "General server error",
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
          "description": "color Id",
          "type": "integer",
          "format": "int32"
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
  },
  "parameters": {
    "colorId": {
      "type": "integer",
      "format": "int32",
      "description": "color ID",
      "name": "colorId",
      "in": "path",
      "required": true
    },
    "id": {
      "type": "string",
      "description": "image id",
      "name": "id",
      "in": "query",
      "required": true
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
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
        "produces": [
          "application/json"
        ],
        "summary": "get board colors",
        "operationId": "getBoardColors",
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
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "summary": "post board colors",
        "operationId": "postBoardColors",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "type": "object",
              "$ref": "#/definitions/Color"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "500": {
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/colors/{colorId}": {
      "delete": {
        "summary": "delete borad color",
        "operationId": "delBoardColor",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "description": "color ID",
            "name": "colorId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "500": {
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/download": {
      "get": {
        "produces": [
          "image/png,image/jpeg"
        ],
        "summary": "Download image",
        "operationId": "downloadFile",
        "parameters": [
          {
            "type": "string",
            "description": "image id",
            "name": "id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "download file OK",
            "schema": {
              "type": "string",
              "format": "byte"
            }
          },
          "400": {
            "description": "Bad Argument",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "404": {
            "description": "File not found",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "General server error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/upload": {
      "post": {
        "description": "The file can't be larger than **5MB**",
        "consumes": [
          "multipart/form-data"
        ],
        "summary": "Upload image",
        "operationId": "uploadFile",
        "parameters": [
          {
            "type": "file",
            "description": "The file to upload",
            "name": "file",
            "in": "formData"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "id"
              ],
              "properties": {
                "id": {
                  "description": "image id",
                  "type": "string"
                }
              }
            }
          },
          "500": {
            "description": "General server error",
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
          "description": "color Id",
          "type": "integer",
          "format": "int32"
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
  },
  "parameters": {
    "colorId": {
      "type": "integer",
      "format": "int32",
      "description": "color ID",
      "name": "colorId",
      "in": "path",
      "required": true
    },
    "id": {
      "type": "string",
      "description": "image id",
      "name": "id",
      "in": "query",
      "required": true
    }
  }
}`))
}
