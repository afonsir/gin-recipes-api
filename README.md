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
