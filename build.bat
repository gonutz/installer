cd pathed
call build.bat
cd ..

copy pathed\pathed.exe task\pathed.exe
cd task
call make_pathed_bin.bat
cd ..

del pathed\pathed.exe
del task\pathed.exe

go build

pause