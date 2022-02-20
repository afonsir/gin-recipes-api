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
JWT_SECRET='<TOKEN>' MONGODB_URI='mongodb://<USER>:<PASSWORD>@localhost:27017/test?authSource=admin' MONGODB_DATABASE=demo go run *.go
```

# Healthcheck

```bash
docker container inspect --format "{{json .State.Health }}" redis | jq '.Log[].Output'
```

# Performance Benchmark

```bash
# -n num of requests
# -c num of concurrent requests
# -g file name (ex. with-cache.data or without-cache.data)
ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipes
```
