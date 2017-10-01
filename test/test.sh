#!/bin/bash

###########
# Run tests
#
bin/sdl -f ./test/test.config.yml write tree-write
sleep 1
bin/sdl -f ./test/test.config.yml read tree-read
sleep 1
bin/sdl -f ./test/test.config.yml write qa
sleep 1
bin/sdl -b testB -f ./test/test.config.yml read prod


zip ./test/zipTest.zip ./test/zipTest
bin/sdl -f ./test/test.config.yml write zipFileWrite
sleep 1
bin/sdl -f ./test/test.config.yml read zipFileRead
sleep 1
unzip -o ./test/tmpRead/zipTestRead.zip -d ./test/tmpRead/

##############
# Check result
#

### File contents
if ! diff -q test/zipTest test/tmpRead/test/zipTest &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

if ! diff -q test/qa-foo test/create/prod-foo &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

if ! diff -q test/qa-bar test/create2/many/levels/prod-bar &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

if ! diff -q test/createTree/a/a test/secretTree/a/a &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

if ! diff -q test/createTree/b/c test/secretTree/b/c &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

if ! diff -q test/createTree/d test/secretTree/d &>/dev/null; then
   echo "Test failed!"
   exit 1
fi

### File mods
if [ $(stat -c "%a" test/create/many/levels/prod-fooo) -ne 600 ]; then
    echo "test/create/many/levels/prod-fooo is $(stat -c "%a" test/create/many/levels/prod-fooo)"
    exit 1
fi

if [ $(stat -c "%a" test/create/prod-foo) -ne 740 ]; then
    echo "test/create/prod-foo is $(stat -c "%a" test/create/many/levels/prod-fooo)"
    exit 1
fi
