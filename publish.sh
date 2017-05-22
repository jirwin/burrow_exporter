#!/bin/bash

VERSION=$(git describe --tags)

gox -ldflags "-X main.Version $VERSION -X" -output "dist/{{.OS}}_{{.ARCH}}/burrow-exporter"

rm -rf releases
mkdir releases

for i in src/* ; do
  if [ -d "$i" ]; then
   ARCH=$(basename "$i")
   zip releases/burrow-exporter-$VERSION-$ARCH.zip dist/$ARCH/burrow-exporter
  fi
done

ghr -t $GITHUB_TOKEN -u jirwin -r burrow-exporter --replace `git describe --tags` dist/
