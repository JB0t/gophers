#!/bin/bash

_rm_hpn_path(){
    h_out_path="${GOPATH:-HOME}/.h_out"
    dir_path=$(cat $h_out_path 2>/dev/null)
    if [[ ! -z "$(readlink -e "$dir_path" 2>/dev/null)" ]]; then
        cd "$dir_path"
        rm "$h_out_path"
        return 0
    fi
}

hpn(){
    if [ ! command -v harpoon &> /dev/null ]; then
        echo "harpoon could not be found. Please build and install it."
        return 1
    fi
    if [[ ! -z "$(readlink -e "$1" 2>/dev/null)" ]]; then
        harpoon -a "$1"
        return 0
    fi
    if [[ "$1" =~ ^[0-9]+$ ]]; then
        harpoon -i "$1"
        _rm_hpn_path
        return 0
    else
        harpoon "$@"
        _rm_hpn_path
        return 0
    fi
    echo "Nothing happened.."
}

