#!/bin/sh

# Post a new repo
curl --data '{"URL":"git@gitlab.com:NatoBoram/git-to-ipns.git"}' -X POST http://localhost:62458/api/repos/

# Get the list of repos
curl -X GET http://localhost:62458/api/repos/

# Get a single repo
curl -X GET http://localhost:62458/api/repos/git@gitlab.com:NatoBoram/git-to-ipns.git

# Delete a single repo
curl -X DELETE http://localhost:62458/api/repos/git@gitlab.com:NatoBoram/git-to-ipns.git
