# Go packages

- Install packages:

```bash
go mod tidy
```

# Documentation

- Generate Swagger documentation:

```bash
swagger generate spec --output ./swagger.json
```

- Serve Swagger documentation:

```bash
swagger serve --flavor swagger ./swagger.json
```

# Start App

```bash
MONGODB_URI='mongodb://<USER>:<PASSWORD>@localhost:27017/test?authSource=admin' MONGODB_DATABASE=demo go run main.go
```
