#!/bin/bash
mkdir src
cd src
touch main.lua
echo 'print("Hello, from lua")' > main.lua
touch script.c
echo "#include <stdio>" > script.c
cd ..
mkdir st
cd st
touch store.mlist
echo "1 ml" > store.mlist
echo '"store data"' >> store.mlist
touch more.mlist
echo "1 ml" > more.mlist
echo '"more data"' >> more.mlist
cd src
