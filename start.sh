#!/bin/sh

echo "INSERT A FRAMERATE AND CONFIRM WITH \"ENTER\" [8]"
read -s -t 30 FPS

if [ "$FPS" == "" ]
  then
  FPS=8
fi

go run snake.go -w $(tput cols) -h $(tput lines) -f $FPS
