#!/bin/sh
cd $1
git pull origin master
git add --all
NOW=$(date +"%m-%d-%Y %H-%M")
git commit -am "Auto-committed on $NOW"
git push  origin master
