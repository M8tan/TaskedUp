@echo off
setlocal enabledelayedexpansion
set /p Desc=Enter message:
cd .\
git add .
git commit -m "%Desc%"
git push
set "line="
for /L %%i in (1,1,10) do (
    set "line=!line!*"
    echo !line!
)
echo Done!
timeout /t 2 /nobreak>nul
exit