Dev notes

Docker compose services:
- backend-go: :8080
- frontend: dev server :5173 (proxy to backend)
- backend-cpp: stub :9090

Run:
```
docker-compose up --build
```


