# SSH Config File

Create the file
```
~/.ssh/config
```

The contents:
```
Host myRemoteHost
    HostName myRemoteHost.com
    User sean
    Port 22
    IdentityFile ~/.ssh/id_rsa
```
