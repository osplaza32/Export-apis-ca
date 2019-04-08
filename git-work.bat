@echo off
echo type "commit" or "update"
for /F "tokens=2" %%i in ('date /t') do set mydate=%%i
set mytime=%time%
set arg1=%1
set arg2=%2
@cd %arg1%
@git checkout -b %arg2%
@git add --all
@git commit -am "Auto-committed on %arg2% %mydate%:%mytime%"
@git push origin %arg2% -u