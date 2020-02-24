This repo is a proof of concept for running an AWS Go runtime Lambda as a proxy in front of a gRPC micro-service.

If you dont work at caring you may use this repo as a reference for doing this yourself.

If you work at caring, you may pull down the call scoring service and start it locally in docker.

Then start this lambda up using serverless offline

```bash
$ make
env GOOS=linux go build -ldflags="-s -w" -o bin/gateway gateway/main.go

$ npm start

> grpc-lambda-poc@0.0.1 start /Users/andrew/work/grpc-lambda-poc
> sls offline --P 3001 -e SLS_DEBUG=*

Serverless: Starting Offline: dev/us-east-1.

Serverless: Routes for gateway:
Serverless: GET /gateway
Serverless: POST /{apiVersion}/functions/grpc-lambda-dev-gateway/invocations

Serverless: Offline [HTTP] listening on http://localhost:3001
Serverless: Enter "rp" to replay the last request
```

Then you can go to postman and make a body-less get request to `http://localhost:3001/gateway`. The call-scoring service will respond through the Lambda proxy!
