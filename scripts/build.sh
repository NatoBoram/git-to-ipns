#!/bin/sh

./scripts/ipfs.sh
hulk ./web/mustache/*.html --outputdir ./web/templates/
rice embed-go
