#!/bin/sh

go get -u -v github.com/GeertJohan/go.rice
go get -u -v github.com/GeertJohan/go.rice/rice

sudo npm i -g hogan.js
sudo chown --changes --recursive $USER:$USER $HOME
