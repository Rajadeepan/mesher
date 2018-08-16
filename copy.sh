#!/bin/sh


rm mesher mesher.tar  image/mesher

docker rmi -f mesher:raj

go build .

cp mesher image/

docker build -t mesher:raj .

docker save mesher:raj > mesher.tar

scp -P 12238   mesher.tar  root@14.141.84.187:/root/rajdeepan
