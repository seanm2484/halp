# Goroutines

```golang
// threads go into this @sem variable. So we can have 500 threads
// otherwise we hit disk i/o issues
sem := make(chan struct{}, 500)
var wg sync.WaitGroup
wg.Add(1024 * 1024)
// walk the directories recursive
err := filepath.Walk(path,
    func(p string, info os.FileInfo, err error) error {
        if err != nil {
                return err
        }
        if filepath.Ext(p) != ".json" {
                return nil
        }
        // get the JSON data in the file. This will contain URLs to
        // all the PoCs for that CVE
        jsonData, err := loadJson(p)
        if err != nil {
                return err
        }

        // now that we have the PoC URLs, download each PoC
        for _, j := range jsonData {
            go func(p string, url string) {
                // if there are 500 goroutines running, this will block
                // waiting for 1 to finish
                sem <- struct{}{}
                // once this goroutine finishes, empty the buffer by one
                // so the next process may start
                defer func() { <-sem }()
                // this must be defrred after a read from @sem
                defer wg.Done()

                downloadFile(filepath.Base(p), url+"/archive/master.zip")
            }(p, j.HTMLURL)
        }

        return nil
    })
if err != nil {
        log.Println(err)
}
wg.Wait()
close(sem)
```