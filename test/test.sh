#!/bin/bash

###########
# Run tests
#
bin/sdl -f ./test/test.config.yml write qa
sleep 1
bin/sdl -f ./test/test.config.yml read prod

##############
# Check result
#
if ! diff -q test/qa-foo test/prod-foo &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

if ! diff -q test/qa-bar test/prod-bar &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

