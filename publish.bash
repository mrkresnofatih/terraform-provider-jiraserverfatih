#! /bin/bash
VersionId=`date +%Y%m%d%H%M%S`
git add . && git status && git commit -m "update tf provider to version v"$VersionId"-dev" && git push origin dev && git tag "v1.0."$VersionId"-dev" && git push origin "v1.0."$VersionId"-dev"