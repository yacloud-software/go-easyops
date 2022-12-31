all:
	for thi in *; do if [ -f $$thi/Makefile ]; then ( echo Making $$thi ; cd $$thi && make ) || exit 10 ; fi; done
