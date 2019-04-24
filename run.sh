#!/bin/bash

echo -e "\x1b[30mUNKO\x1b[0m" | go run main.go out/black.png
echo -e "\x1b[31mUNKO\x1b[0m" | go run main.go out/red.png
echo -e "\x1b[32mUNKO\x1b[0m" | go run main.go out/green.png
echo -e "\x1b[33mUNKO\x1b[0m" | go run main.go out/yellow.png
echo -e "\x1b[34mUNKO\x1b[0m" | go run main.go out/blue.png
echo -e "\x1b[35mUNKO\x1b[0m" | go run main.go out/magenta.png
echo -e "\x1b[36mUNKO\x1b[0m" | go run main.go out/syan.png
echo -e "\x1b[37mUNKO\x1b[0m" | go run main.go out/white.png
