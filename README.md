# go-ipapi
Unofficial lib for https://ip-api.com/. Thanks for the awesome API

# Example usage
```
  import "github.com/Z-M-Huang/go-ipapi"

  resp, err := ipapi.Get(host)
  if err != nil {
    panic(err)
  }
  fmt.Println(resp.Query)
```