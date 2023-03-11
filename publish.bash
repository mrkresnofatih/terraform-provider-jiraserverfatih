#! /bin/bash
VersionId=`date +%Y.%m.%d.%H.%M`
git add . && git status && git commit -m "update tf provider" && git tag "v"$VersionId"-dev" && git push origin "v"$VersionId"-dev"