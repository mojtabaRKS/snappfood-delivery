# snappfood delivery project

### installation

```
$ make init && make up
$ go run cmd/main.go server
```

use `localhost:8065` with these routes :
```
POST /orders/report -> body: {"order_id" : 13}

GET /orders/proccess -> header : employee_id = 10

GET /vendors/report

``` 