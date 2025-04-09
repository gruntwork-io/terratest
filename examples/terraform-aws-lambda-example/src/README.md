# AWS Lambda Function Handler Source

The lambda executable `handler` was built using

``` shell
go get github.com/aws/aws-lambda-go/lambda
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap .
```
