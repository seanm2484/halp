# Download a file over HTTP

```golang
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()

out, err := os.Create(filepath)
if err != nil {
    return err
}
defer out.Close()

// copy body into file
_, err = io.Copy(out, resp.Body)
```