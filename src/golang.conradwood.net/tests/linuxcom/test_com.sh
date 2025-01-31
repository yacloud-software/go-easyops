#!/bin/bash
a=0;
/tmp/test-daemon -d $$ &
max=40
while [ true ]; do
    a=$((a+1))
    sleep 0.1
    CGROUP=`cat /proc/self/cgroup`
    echo  "random text $a in ${CGROUP} "
    if [ $a -gt $max ]; then
	echo
	a=0
	max=$((1 + $RANDOM % 100))
    fi
    
done
