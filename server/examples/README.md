This folder contains example JSON requests used to manually test the API with CURL.

```bash
curl -X POST "http://localhost:8080/enc/new-key" -H "Content-Type: application/json" -d @server/examples/enc/new-key.json -sS | jq

curl -X POST "http://localhost:8080/enc/encrypt" -H "Content-Type: application/json" -d @server/examples/enc/encrypt.json -sS | jq
```
