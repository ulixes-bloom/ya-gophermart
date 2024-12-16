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
        "/api/user/balance": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Получение текущего баланса пользователя",
                "operationId": "GetUserBalance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса"
                    },
                    "401": {
                        "description": "пользователь не авторизован"
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера"
                    }
                }
            }
        },
        "/api/user/balance/withdraw": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Запрос на списание средств",
                "operationId": "WithdrawFromUserBalance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "User Registration Information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_ulixes-bloom_ya-gophermart_internal_models.WithdrawalRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса"
                    },
                    "401": {
                        "description": "пользователь не авторизован"
                    },
                    "402": {
                        "description": "на счету недостаточно средств"
                    },
                    "422": {
                        "description": "неверный номер заказа"
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера"
                    }
                }
            }
        },
        "/api/user/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Аутентификация пользователя",
                "operationId": "AuthUser",
                "parameters": [
                    {
                        "description": "User Registration Information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_ulixes-bloom_ya-gophermart_internal_models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "пользователь успешно аутентифицирован"
                    },
                    "400": {
                        "description": "неверный формат запроса"
                    },
                    "401": {
                        "description": "неверная пара логин/пароль"
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера"
                    }
                }
            }
        },
        "/api/user/orders": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Получение списка загруженных номеров заказов",
                "operationId": "GetUserOrders",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса"
                    },
                    "204": {
                        "description": "нет данных для ответа"
                    },
                    "401": {
                        "description": "пользователь не авторизован"
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера"
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Загрузка номера заказа",
                "operationId": "RegisterUserOrder",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "User Registration Information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_ulixes-bloom_ya-gophermart_internal_models.OrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "номер заказа уже был загружен этим пользователем"
                    },
                    "202": {
                        "description": "новый номер заказа принят в обработку"
                    },
                    "400": {
                        "description": "неверный формат запроса"
                    },
                    "401": {
                        "description": "пользователь не аутентифицирован"
                    },
                    "409": {
                        "description": "номер заказа уже был загружен другим пользователем"
                    },
                    "422": {
                        "description": "неверный формат номера заказа"
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера"
                    }
                }
            }
        },
        "/api/user/register": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Регистрация пользователя",
                "operationId": "RegisterUser",
                "parameters": [
                    {
                        "description": "User Registration Information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_ulixes-bloom_ya-gophermart_internal_models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "пользователь успешно зарегистрирован и аутентифицирован"
                    },
                    "400": {
                        "description": "неверный формат запроса"
                    },
                    "409": {
                        "description": "логин уже занят"
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера"
                    }
                }
            }
        },
        "/api/user/withdrawals": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Получение информации о выводе средств",
                "operationId": "GetUserWithdrawals",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса"
                    },
                    "204": {
                        "description": "нет ни одного списания"
                    },
                    "401": {
                        "description": "пользователь не авторизован"
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера"
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_ulixes-bloom_ya-gophermart_internal_models.OrderRequest": {
            "type": "object",
            "properties": {
                "number": {
                    "type": "string"
                }
            }
        },
        "github_com_ulixes-bloom_ya-gophermart_internal_models.User": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "github_com_ulixes-bloom_ya-gophermart_internal_models.WithdrawalRequest": {
            "type": "object",
            "properties": {
                "order_number": {
                    "type": "string"
                },
                "sum": {
                    "type": "number"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
