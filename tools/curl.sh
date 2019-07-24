#!/bin/sh

curl --data '{"URL":"git@gitlab.com:NatoBoram/git-to-ipfs.git"}' -X POST http://localhost:62458/api/add/
