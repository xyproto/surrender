#!/bin/sh
go build -mod=vendor && ./render ../../testdata/xyproto.svg output.png
