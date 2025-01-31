#!/bin/bash
a=0;
max=40
while [ true ]; do
    a=$((a+1))
    echo -n "random text $a "
    if [ $a -gt $max ]; then
	echo
	a=0
	max=$((1 + $RANDOM % 100))
    fi
    
done
