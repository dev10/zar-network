## Zar 설치하기

이 가이드는 `zard`와 `zarcli`를 엔트리포인트를 시스템에 설치하는 방법을 설명합니다. `zard`와 `zarcli`가 설치된 서버를 통해 [풀노드](./join-testnet.md#run-a-full-node) 또는 [밸리데이터로](./validators/validator-setup.md)써 최신 테스트넷에 참가하실 수 있습니다.

### Go 설치하기

공식 [Go 문서](https://golang.org/doc/install)를 따라서 `go`를 설치하십시오. `$GOPATH`, 그리고 `$PATH`의 환경을 꼭 세팅하세요. 예시: 

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

::: tip
코스모스 SDK를 운영하기 위해서는 **Go 1.11.ㅎ+** 이상 버전이 필요합니다.
:::

### 바이너리 설치하기

다음은 최신 Zar 버전을 설치하는 것입니다. 예시에서는 최신 스테이블 릴리즈가 포함되어있는 `master` 브랜치를 이용해 진행됩니다. 필요에 따라 `git checkout`을 통해 [최신 릴리즈](https://github.com/cosmos/cosmos-sdk/releases)가 설치되어있는지 확인하세요.

```bash
git clone -b <latest-release-tag> https://github.com/zar-network/zar-network
cd zar && make install
```

> *참고*: 여기에서 문제가 발생한다면, Go의 최신 스테이블 버전이 설치되어있는지 확인하십시오.

위 절차를 따라하시면 `zard`와 `zarcli` 바이너리가 설치될 것입니다. 설치가 잘 되어있는지 확인하십시오:


```bash
$ zard version
$ zarcli version
```

### 다음 절차

축하합니다! 이제 [퍼블릭 테스트넷](./join-testnet.md)에 참가하시거나 or [프라이빗 테스트넷](./deploy-testnet.md)을 운영하실 수 있습니다.
