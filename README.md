# xonic
WIP: Stream gamepads from one linux box to another


## Note
I have to revisit the scope of this project because this is possible with two lines of bash given ds4drv

```sh
ds4drv --emulate-xpad-wireless --next-controller --emulate-xpad-wireless --hidraw
ssh username@gamepadhost cat /dev/input/event20 >/dev/input/event26
```


## Revisited goals
A netcat proxy with authentication for streaming gamepad input events

## Goals (deprecated)

1. Installation candidates for client + server
2. Configureable transport layer over (ssh first)
3. Compatible with ds4drv

## Credits
Based on work from https://yingtongli.me/blog/2019/12/01/input-over-ssh-2.html 

