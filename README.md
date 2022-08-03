# halp

`halp` is an interactive cheat sheet helper, from the terminal. 

https://asciinema.org/a/cngSkLVa4L4XoTdM1ut5x5H19

## Creating Cheatsheets

Cheatsheets are read from the directory `~/.halp/cheatFiles`
They are YAML files, with an easy to understand syntax. An example is shown below:

```yaml
- description: send a http request and write the output to a file
  command: curl -X _method_ _url_ -o _filename_
  variables:
    - method
    - url
    - filename
- description: send a get http request and follow redirects
  command: curl -L _url_
  variables:
    - url
```
Variables are put in the command with underscores, then those variables are listed in the `variables` list. Do not mess that up or *shrug* your term may explode? Idk.

## Rendering Markdown Files

You can now render markdown files in retrieved from `halp`. To do this, create a cheatsheet file like before, but use the `file` keyword to point to the path:

```yaml
- description: "[golang] get all files in a directory"
  file: "golang/go_get_all_files_in_dir.md"
```

Then in your markdown file you can do the usual markdown stuff,

```markdown
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
```