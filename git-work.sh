#!/bin/sh
cd C:\Users\Administrador\Desktop\api_management
git pull origin master -f
git add --all
git commit -am "Auto-committed on %date%"
git push  origin master -f
