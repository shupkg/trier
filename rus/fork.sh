#!/usr/bin/env sh

set -e

gitAddr="https://github.com/sirupsen/logrus.git"
oldPkg="github.com/sirupsen/logrus"
newPkg="github.com/shupkg/trier/rus"
oldName="logrus"
newName="rus"

cd "$(dirname "$0")"

getAuthor(){
  cd tmp
  git log --abbrev-commit --max-count=1 | head -n 2 | tail -n 1 | sed -e 's/Author: \(.*\)/\1/'
  cd ..
}

getDate(){
  cd tmp
  git log --abbrev-commit --max-count=1 | head -n 3 | tail -n 1 | sed -e 's/Date: \(.*\)/\1/'
  cd ..
}

getVersion() {
  cd tmp
  git describe --tags --long --always --dirty
  cd ..
}

find . -depth 1 ! -name fork.sh -exec rm -rf {} \;

if [ -d ./tmp ]; then
  rm -rf tmp
fi

git clone --depth=1 ${gitAddr} tmp

find ./tmp -name \*.go -exec sed -i.bak "s?${oldPkg}?${newPkg}?g;s?${oldName}?${newName}?g" {} \;
find ./tmp -type d ! -path \*.\* ! -name tmp ! -name example -exec sh -c 'mkdir -p ${1:6};echo mkdir -p ${1:6}' sh {} \;
find ./tmp -type f -name \*.go ! -name \*_test.go -exec sh -c 'mv $1 ${1:6}; echo mv $1 ${1:6}' sh {} \;

fVersion=$(getVersion)
fAuthor=$(getAuthor)
fDate=$(getDate)
rm -rf tmp

go build -v

cat >fork.txt <<EOF
fork time:
	$(date -R)
fork from:
	${gitAddr}
fork version:
    ${fVersion} ${fAuthor} ${fDate}
EOF
