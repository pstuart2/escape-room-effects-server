#!/usr/bin/env bash

sudo python clockEncoder.py 17 18 "http://192.168.86.50:8080/hours" &
sudo python clockEncoder.py 5 6 "http://192.168.86.50:8080/minutes" &