# URL Shortener
## About
Very simple URL shortener written with golang

## How can I run it?
### #1 Make sure the docker is installed
```bash
docker --version
```
##### If the docker is not installed, you can use this link
<a href="https://docs.docker.com/engine/install/">Install Docker</a>

### #2 Make sure the docker service is running
```bash
sudo systemctl start docker.service
```

### #3 Execute docker-compose build command
```bash
sudo docker-compose build --no-cache
```

### #4 Execute docker-compose up command
```bash
sudo docker-compose up --force-recreate
```

### #5 Make sure everything is OK! :wink:
##### Send request to 'http://localhost:9000/api/ping':
```bash
curl -X GET "http://localhost:9000/api/ping" -H "Accept: application/json"
```
##### Result sample:
```json
{
  "status": true,
  "message": "ok",
  "errors": null,
  "data": "pong"
}
```

<hr>

### URL section
#### #1 Create new short url
```bash
curl -X POST "http://localhost:9000/api/create-url" -H "Content-Type: application/json" -H "Accept: application/json" -d '{"original_url": "http://google.com"}'
```
##### Result sample:
```json
{
  "status": true,
  "message": "created",
  "errors": null,
  "data": {
    "original_url": "http://google.com",
    "short_url": "http://localhost:9000/OZcovl3I"
  }
}
```

#### #2 Use from the short url
##### Open 'short_url' value in the browser and you will be redirected to 'original_url'! :sunglasses:

#### #3 Get short url detail
```bash
curl -X GET "http://localhost:9000/api/{url_name}" -H "Accept: application/json"
```
##### Result sample:
```json
{
  "status": true,
  "message": "ok",
  "errors": null,
  "data": {
    "original_url": "http://google.com",
    "click": 3
  }
}
```

<hr>

## How can I test it?
### #1 Get into Docker container shell
```bash
sudo docker exec -it url_shortener_app bash
```

### #2 Execute artisan command
```bash
go clean -testcache
go test ./...
```