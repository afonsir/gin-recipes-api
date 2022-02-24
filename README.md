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

- JWT auth mechanism:

```bash
AUTH_MECHANISM='JWT' JWT_SECRET='<TOKEN>' MONGODB_URI='mongodb://<USER>:<PASSWORD>@localhost:27017/test?authSource=admin' MONGODB_DATABASE=demo go run *.go
```

- Cookie auth mechanism:

```bash
AUTH_MECHANISM='COOKIE' MONGODB_URI='mongodb://<USER>:<PASSWORD>@localhost:27017/test?authSource=admin' MONGODB_DATABASE=demo go run *.go
```

- Auth0 auth mechanism:

```bash
AUTH_MECHANISM='AUTH0' AUTH0_DOMAIN='<AUTH0_DOMAIN>' AUTH0_API_IDENTIFIER='<AUTH0_API_ID>' MONGODB_URI='mongodb://<USER>:<PASSWORD>@localhost:27017/test?authSource=admin' MONGODB_DATABASE=demo go run *.go
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

# Security

- Generate self-signed certificate:

```bash
openssl req \
  -x509 \
  -nodes \
  -days 365 \
  -newkey rsa:2048 \
  -keyout certs/localhost.key \
  -out certs/localhost.crt
```