#!/bin/bash
go vet . 2>/tmp/govet.txt >/tmp/govet.txt
RESULT=$?
if [ "$(cat /tmp/govet.txt | grep "go: downloading" | wc -l)" -ne 0 ]; then
  echo -e "::group::\e[0;34mð¥ Dependency downloads\e[0m"
  cat /tmp/govet.txt | grep "go: downloading" | sed -e "s/go: downloading /   ð¦ /"
  echo "::endgroup::"
fi
if [ $RESULT -ne 0 ]; then
  echo -e "::group::\e[0;31mâ go vet found problems\e[0m"
  cat /tmp/govet.txt | grep -v "go: downloading" | sed -e 's/# /   ð¦ /g' -e "s#\./#      ð $(echo -ne '\e[0;1;97m')#g" -e "s/: /:$(echo -ne '\e[0m') /g"
  echo "::endgroup::"
else
  echo -e "\e[0;32mâ go vet found no problems\e[0m"
fi
exit $RESULT