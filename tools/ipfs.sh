#!/bin/sh

# Bootstrap
ipfs add --wrap-with-directory=true --chunker=rabin --cid-version=1 https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css.map
ipfs add --wrap-with-directory=true --chunker=rabin --cid-version=1 https://code.jquery.com/jquery-3.3.1.slim.min.js
ipfs add --wrap-with-directory=true --chunker=rabin --cid-version=1 https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js.map
ipfs add --wrap-with-directory=true --chunker=rabin --cid-version=1 https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js.map

# Hogan.JS
ipfs add --wrap-with-directory=true --chunker=rabin --cid-version=1 https://twitter.github.io/hogan.js/builds/3.0.1/hogan-3.0.1.js
