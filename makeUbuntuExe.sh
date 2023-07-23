#!/usr/bin/bash

if [ -d "AbsenceHelperForUbuntu" ]; then
	echo "The folder AbsenceHelperForUbuntu already exists. Do you want to update it? (y/n)"
	read answer
	if [ "$answer" != "y" ]; then
		echo "Aborting script execution."
		exit 1
	fi
fi

if go build &&
	if [ ! -d "AbsenceHelperForUbuntu" ]; then
		mkdir AbsenceHelperForUbuntu
	fi
	mv -f AbsenceHelper AbsenceHelperForUbuntu &&
	cp -f config.json AbsenceHelperForUbuntu; then
	echo "Success!"
else
	echo "Failure!"
	exit 1
fi

