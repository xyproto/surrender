#!/bin/sh
go build -mod=vendor
for x in circle empty rainforest_2c_opt rainforest_8c_opt xyproto; do
  echo -e -n "Rendering $x:\t\t"
  ./render ../../testdata/$x.svg $x.png || echo FAIL
done
