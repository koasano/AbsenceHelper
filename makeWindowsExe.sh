#!/usr/bin/bash

# Check if AbsenceHelperForWindows folder exists
if [ -d "AbsenceHelperForWindows" ]; then
	echo "The folder AbsenceHelperForWindows already exists. Do you want to update it? (y/n)"
	read answer
	if [ "$answer" != "y" ]; then
		echo "Aborting script execution."
		exit 1
	fi
fi

# Try to build the Windows executable and copy the configuration file
if env GOOS=windows GOARCH=amd64 go build -o AbsenceHelper.exe &&
	if [ ! -d "AbsenceHelperForWindows" ]; then
		mkdir AbsenceHelperForWindows
	fi
	mv -f AbsenceHelper.exe AbsenceHelperForWindows &&
	cp -f config.json AbsenceHelperForWindows; then
	echo "Success!"
else
	echo "Failure!"
	exit 1
fi

