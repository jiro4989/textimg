#!/bin/bash

make build
echo -e "\x1b[30mUNKO\x1b[0m" | ./bin/txtimg out/black.png
echo -e "\x1b[31mUNKO\x1b[0m" | ./bin/txtimg out/red.png
echo -e "\x1b[32mUNKO\x1b[0m" | ./bin/txtimg out/green.png
echo -e "\x1b[33mUNKO\x1b[0m" | ./bin/txtimg out/yellow.png
echo -e "\x1b[34mUNKO\x1b[0m" | ./bin/txtimg out/blue.png
echo -e "\x1b[35mUNKO\x1b[0m" | ./bin/txtimg out/magenta.png
echo -e "\x1b[36mUNKO\x1b[0m" | ./bin/txtimg out/syan.png
echo -e "\x1b[37mUNKO\x1b[0m" | ./bin/txtimg out/white.png
echo Test | grep --color=always Te | ./bin/txtimg out/grep.png
