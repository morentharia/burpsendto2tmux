# burpsendto2tmux
connect tmux and "Custom Send To" burp plugin

Specify how to run commands in terminal:

macos
```
/opt/X11/bin/xterm +hold -e /Users/mavostrykh/GOPATH/bin/burpsendto2tmux -path /usr/local/opt/tmux/bin/ -c %C
```

ubuntu
```
/usr/bin/xterm +hold -e /home/mor/GOPATH/bin/burpsendto2tmux -path /usr/local/bin -c %C
```
