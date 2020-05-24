# go-ipapi
Unofficial lib for https://ip-api.com/. Thanks for the awesome API

[![Build Status](https://travis-ci.com/Z-M-Huang/go-ipapi.svg?branch=master)](https://travis-ci.com/Z-M-Huang/go-ipapi)

# Example usage
```
  import "github.com/Z-M-Huang/go-ipapi"

  resp, err := ipapi.Get(host)
  if err != nil {
    panic(err)
  }
  fmt.Println(resp.Query)
```