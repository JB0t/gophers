# Harpoon

Shameless bash-fs re-creation of ThePrimagen's NeoVim plugin [Harpoon](https://github.com/ThePrimeagen/harpoon)

## Setup

```bash
# Install the base script
go install .
# Add the bash wrapper funcs to your bash profile
echo "source $PWD/harpoon.sh >> ~/.bashrc
source ~/.bashrc
```

## Usage

```bash
# list current dirs
hpn
# Add a directory
hpn -a path/to/dir
# Go to a directory in the list
hpn -i $index
# Remove a directory
hpn -r $index
# Clear list
hpn -c
```
