@echo off
echo type "commit" or "update"
set arg1=%1
set arg2=%2
@cd %arg1%
for /F "tokens=2" %%i in ('date /t') do set mydate=%%i
set mytime=%time%
@git add . && git commit -m "Auto-committed on %arg2% %mydate%:%mytime%"
@git push -u origin %arg2%  -f
