## Run 
```
go run .
``` 

```
go run . -j "{\"a\":[{\"abc\":\"123\",\"x\":\"hello\"},{\"abc\":\"456\",\"x\":\"world\"},{\"abc\":\"456\",\"x\":\"test\"}]}" -p "$.a[?(@.abc=='456')].x" -o replace -v "\"TEST\""
```

## Build
```
go build -o ./.output/jpbp .
```