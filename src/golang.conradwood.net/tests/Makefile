all:
	for thi in *; do if [ -f $$thi/Makefile ]; then ( cd $$thi && make ) || exit 10 ; fi; done
