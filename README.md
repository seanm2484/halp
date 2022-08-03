# halp

`halp` is an interactive cheat sheet helper, from the terminal. 

https://asciinema.org/a/cngSkLVa4L4XoTdM1ut5x5H19

## Creating Cheatsheets

Cheatsheets are read from the directory `./cheatFiles`
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