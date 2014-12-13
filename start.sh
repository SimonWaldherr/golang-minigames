#!/bin/sh

echo "INSERT A FRAMERATE AND CONFIRM WITH \"ENTER\" [8]"
read -t 30 FPS

if [ "x$FPS" = "x" ]
  then
  FPS=8
fi

stty cbreak min 1
stty -echo

go run snake.go -w $(tput cols) -h $(tput lines) -f $FPS

stty echo
