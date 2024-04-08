LEGAL_ASSISTANT (WORKSPACE)
|-- cmd
|   `-- legal_assistant_chat_pe..
|       `-- main.go
|-- internal
|   |-- app
|   |   |-- handlers
|   |   |   |-- errors.go
|   |   |   `-- router.go
|   |   `-- http
|   |       `-- server.go
|   |-- domain
|   |   |-- daos
|   |   |-- dtos
|   |   |-- models
|   |   |-- repositories
|   |   `-- services
|   `-- pkg
|       |-- db
|       |   |-- mongo.go
|       |   `-- services.go
|       |-- helpers
|       |   `-- nullable.go
|       `-- utils
|           |-- auth
|           |-- env
|           |   |-- env_tool.go
|           |   |-- inflector
|           |   `-- validator
|           `-- mock
|-- deployments
|   |-- docker-compose.yml
|   `-- k8s
|-- docs
|   `-- README.md
|-- pkg
|   `-- conf
|       `-- http_conf.go
|-- .gitignore
|-- Dockerfile
|-- go.mod
`-- go.sum