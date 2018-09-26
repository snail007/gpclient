#!/bin/bash
VER="1.0"
RELEASE="release-${VER}"
TRIMPATH1="/Users/snail/go/src/github.com/snail007"
TRIMPATH=$(dirname ~/go/src/github.com/snail007)/snail007
if [ -d "$TRIMPATH1" ];then
	TRIMPATH=$TRIMPATH1
fi
OPTS="-gcflags=-trimpath=$TRIMPATH -asmflags=-trimpath=$TRIMPATH"

rm -rf ${RELEASE}
mkdir ${RELEASE}

#linux
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-386.tar.gz" gpclient
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-amd64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm-v6.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=6 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm64-v6.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm-v7.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=7 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm64-v7.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm-v5.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=5 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm64-v5.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm64-v8.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-arm-v8.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-mips.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-mips64.tar.gz" gpclient
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-mips64le.tar.gz" gpclient
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -o gpclient $OPTS -ldflags "-s -w"  && tar zcfv "${RELEASE}/gpclient-linux-mipsle.tar.gz" gpclient
CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-mips-softfloat.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 GOMIPS=softfloat go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-mips64-softfloat.tar.gz" gpclient
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le GOMIPS=softfloat go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-mips64le-softfloat.tar.gz" gpclient
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-mipsle-softfloat.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-ppc64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-ppc64le.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-linux-s390x.tar.gz" gpclient 
#android
CGO_ENABLED=0 GOOS=android GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-android-386.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=android GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-android-amd64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=android GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-android-arm.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-android-arm64.tar.gz" gpclient 
#darwin
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-darwin-386.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-darwin-amd64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=darwin GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-darwin-arm.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-darwin-arm64.tar.gz" gpclient 
#dragonfly
CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-dragonfly-amd64.tar.gz" gpclient 
#freebsd
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-freebsd-386.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-freebsd-amd64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=freebsd GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-freebsd-arm.tar.gz" gpclient 
#nacl
CGO_ENABLED=0 GOOS=nacl GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-nacl-386.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=nacl GOARCH=amd64p32 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-nacl-amd64p32.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=nacl GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-nacl-arm.tar.gz" gpclient 
#netbsd
CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-netbsd-386.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-netbsd-amd64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=netbsd GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-netbsd-arm.tar.gz" gpclient 
#openbsd
CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-openbsd-386.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-openbsd-amd64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-openbsd-arm.tar.gz" gpclient 
#plan9
CGO_ENABLED=0 GOOS=plan9 GOARCH=386 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-plan9-386.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=plan9 GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-plan9-amd64.tar.gz" gpclient 
CGO_ENABLED=0 GOOS=plan9 GOARCH=arm go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-plan9-arm.tar.gz" gpclient 
#solaris
CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -o gpclient $OPTS -ldflags "-s -w" && tar zcfv "${RELEASE}/gpclient-solaris-amd64.tar.gz" gpclient 
#windows
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o gpclient-noconsole.exe
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o gpclient.exe && tar zcfv "${RELEASE}/gpclient-windows-386.tar.gz" gpclient.exe gpclient-noconsole.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gpclient-noconsole.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gpclient.exe && tar zcfv "${RELEASE}/gpclient-windows-amd64.tar.gz" gpclient.exe gpclient-noconsole.exe

rm -rf gpclient gpclient.exe gpclient-noconsole.exe

#todo
#1.release.sh        VER="xxx"
#2.main.go           APP_VERSION="xxx"
