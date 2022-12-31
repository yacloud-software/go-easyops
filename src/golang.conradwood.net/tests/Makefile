all:
	for thi in *; do if [ -f $$thi/Makefile ]; then ( cd $$thi && echo "============ Making `pwd` ===========" ; make ) || exit 10 ; fi; done
