#! /bin/bash
VersionId=`date +%Y%m%d%H%M%S`
touch release.txt && echo "v1.0."$VersionId"-dev" > release.txt
git add . && git status && git commit -m "update tf provider to version v1.0."$VersionId"-dev" && git tag "v1.0."$VersionId"-dev" && git push origin "v1.0."$VersionId"-dev"