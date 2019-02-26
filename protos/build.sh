#!/bin/ash

# Make sure we're starting from a clean slate
rm -rf /protos/python/*
rm -rf /protos/golang/*

# For each proto file, generate the protobufs for each
# supported language
for dir in "." "player" "processor" "provider"; do
  for file in /protos/$dir/*.proto; do
    filename=$(basename -- "$file")
    include="-I /protos -I /protos/$dir -I /google-protos/src"
  
    # Golang
    mkdir -p /protos/golang/$dir
    protoc $include $filename --go_out=plugins=grpc:/protos/golang/$dir
  
    # Python
    mkdir -p /protos/python
    protoc $include $filename --python_out=/protos/python
  done
done

# Build the core protos
#for file in /protos/*.proto; do
#  filename=$(basename -- "$file")
#  include="-I /protos -I /google-protos/src"
#
#  # Golang
#  mkdir -p /protos/golang/$dir
#  protoc $include $filename --go_out=plugins=grpc:/protos/golang
#  
#  # Python
#  mkdir -p /protos/python
#  protoc $include $filename --python_out=/protos/python
#done
