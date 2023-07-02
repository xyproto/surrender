#!/bin/sh
go build -mod=vendor && ./render ../../img/xyproto.svg output.png
