#!/bin/bash

# Usage
# ./test_mem.sh mem_test.py
# ./test_mem.sh mem_test_2.py
#
# This very simple script copies the mem test python scripts into the appropriate place and then
# executes the script. Note that you need to docker-exec into the running container and then look at the mem usage:
#
# docker exec -it goodbc-python-test bash
# ps aux | grep mem_test

./test.sh /bin/bash -c "\"ls .. && cp ../$1 . && python $1\""
