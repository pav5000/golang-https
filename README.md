# golang-https
Easy way to add https to your go application with automatic Let's Encrypt certificate getting and updating

## Example

The simplest way to start your https server looks like this:
```go
srv := https.New("./data", "my@email.com", "my.site.com")

err := srv.ListenHTTPS(":443", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
}))
if err != nil {
    log.Fatal(err)
}
```

## Additional functions

- `ListenHTTPRedirect` starts server which redirects all http requests to the corresponding https url.
  + `http://some.site.com/path/to/handler?param1=value1` -> `https://some.site.com/path/to/handler?param1=value1`
- `GetHTTPSServer` returns https server ready to be started. You should use this function if you want to manually tweak some params of the server.
