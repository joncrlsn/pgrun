#!/bin/bash

appname=pgrun

if [[ -d bin-linux ]]; then
    GOOS=linux GOARCH=386 go build -o bin-linux32/${appname}
    echo "Built linux."
else
    echo "Skipping linux.  No bin-linux directory."
fi

if [[ -d bin-osx ]]; then
    GOOS=darwin GOARCH=386 go build -o bin-osx/${appname}
    echo "Built osx."
else
    echo "Skipping osx.  No bin-osx directory."
fi

if [[ -d bin-win ]]; then
    GOOS=windows GOARCH=386 go build -o bin-win/${appname}.exe
    echo "Built win."
else
    echo "Skipping win.  No bin-win directory."
fi
