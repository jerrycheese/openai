#!/bin/bash

mkdir -p output

cp config_prod.yaml output/config.yaml
cp scripts/bootstrap.sh output/

go build -o openaiBin ./
mv openaiBin output/

