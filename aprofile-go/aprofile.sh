#!/bin/bash

_rm_aprofile_path(){
    a_out_path="${GOPATH:-HOME}/.a_out"
    if [[ ! -f "$a_out_path" ]]; then
        echo "aprofile output file not found."
        return 1
    fi
    export AWS_PROFILE=$(cat $a_out_path)
    echo "AWS_PROFILE set to $(cat $a_out_path)"
    rm "$a_out_path"
}

aprofile(){
    if [ ! command -v aprofile-go &> /dev/null ]; then
        echo "aprofile-go could not be found. Please build and install it."
        return 1
    fi
    aprofile-go
    _rm_aprofile_path
    return 0
}

