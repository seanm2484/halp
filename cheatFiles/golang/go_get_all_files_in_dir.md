# Get all the files with ext in a directory

```golang
err := filepath.Walk(path,
    func(p string, info os.FileInfo, err error) error {
        if err != nil {
                return err
        }
        if filepath.Ext(p) != ".json" {
            return nil
        }
        fmt.Println(p, info.Size())
        loadJson(p)
        return nil
    })  
```