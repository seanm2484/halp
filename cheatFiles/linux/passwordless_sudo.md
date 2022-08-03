# Passwordless sudo

To implement passwordless sudo for a specific user, you must modify the `/etc/sudoers` file. To do this,
use the builtin editor `visudo` (which must be run with `sudo`)

Go to the bottom of this file and append the following line:
```
user ALL=(ALL) NOPASSWD:ALL
```

This line must be present __after__ the line:
```
@includedir /etc/sudoers.d
```