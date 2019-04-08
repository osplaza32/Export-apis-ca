@echo off
echo type "commit" or "update"
@cd .\api_management
@git pull origin master -f
@git add --all
@git commit -am "Auto-committed on %date%"
@git push  origin master -f
