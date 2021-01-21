# expense-tracker-api

api to support the expense tracker mobile app.

## technology

- [echo](https://github.com/labstack/echo)
- [mongo-driver](https://github.com/mongodb/mongo-go-driver)
- [air](https://github.com/cosmtrek/air)
- [swag](https://github.com/swaggo/swag/cmd/swag)

## run local developemnt environment

Need to install `air` & `swag` globally. Can be done by this 

```go
go get -u github.com/cosmtrek/air

go get -u github.com/swaggo/swag/cmd/swag
```

To spin up the development stack simply run 

```shell
make dev
```

