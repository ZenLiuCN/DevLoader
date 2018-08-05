@echo off
rsrc.exe -manifest ico.manifest -o main.syso -ico ico.ico
if "%1"=="" (go build -ldflags "-s -w -H windowsgui" ./) else ( go build -ldflags "-H windowsgui" ./)
