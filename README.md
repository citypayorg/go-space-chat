##  孤独 Lonely

![d2139b33a9868d1f17a471201d1272371588868902.jpg](https://cdn.jsdelivr.net/gh/sunshinev/remote_pics/d2139b33a9868d1f17a471201d1272371588868902.jpg)


E:\dev\chat\go-space-chat
cd E:
## Demo

http://chat.osinger.com/

## 特色
1. 支持性别修改、并且有颜色替换
2. 支持敏感词过滤
3. 支持姓名修改
1. 성별 수정이 지원되며 색상 대체가 있습니다.
2. 민감한 단어 필터링을 지원합니다.
3. 이름 변경이 지원됩니다.

## 介绍

通过canvas 2d来模拟了3D的视觉效果。
캔버스 2d를 통해 3D의 시각적 효과를 시뮬레이션했습니다.

并且在该项目中使用了protobuf来进行前端和后端的通讯协议，这一点非常方便！
그리고 프론트 엔드와 백 엔드 통신 프로토콜에 대한 프로젝트에 protobuf를 사용하는 것이 편리합니다!

## 操作

1. 项目使用传统`WASD`按键来控制上下左右
2. 眼睛可以跟随鼠标的位置进行转动
3. 按下`space` 空格可以输入消息，按下回车发送消息
4. 左上角按钮可以输入名称，点击空白处名称生效
1. 프로젝트는 기존의 'WASD' 버튼을 사용하여 위쪽과 아래쪽을 제어합니다.
2. 눈은 마우스의 위치를 따라 회전 할 수 있습니다.
3. 'space' 공백을 눌러 메시지를 입력하고 Enter 키를 눌러 메시지를 보냅니다.
4. 왼쪽 상단 버튼은 이름을 입력하고 빈 공간을 클릭하여 적용 할 수 있습니다.

## 运行 운영

```$xslt
go run main.go
```

该命令会启动web-server作为静态服务，默认80端口，如果需要修改端口，用下面的命令
이 명령은 웹 서버를 정적 서비스, 기본 80 포트로 시작하고 포트를 수정해야 하는 경우 다음 명령을 사용합니다.
```
go run main.go -web_server 8081
```

项目启动默认websocket服务端口为9000端口，如果需要修改
프로젝트 시작 기본 웹소켓 서비스 포트는 필요한 경우 9000포트입니다.
```
go run main.go -socket_server 9001
```
注意：如果修改websocket端口，同时需要修改js里面的socket端口
참고: websocket 포트를 수정하는 경우 js 내의 소켓 포트도 수정해야 합니다.

## 技术工具
프런트 엔드. 
前端 Vue+canvas+websocket+protobuf

백 엔드입니다. 
后端 Golang+websocket+protobuf+goroutine

## 有意思的难点
> 这里列举几个在实现过程中，遇到的很有意思的问题

1. 如何实现无限画布？
2. 如何实现游戏状态同步？
> 구현 과정에서 겪는 흥미로운 문제 몇 가지를 나열합니다.

1. 무한 캔버스를 구현하려면 어떻게 해야 합니까?
2. 게임 상태를 동기화하려면 어떻게 해야 합니까?



## 相关链接

[Canvas 基本用法](https://developer.mozilla.org/zh-CN/docs/Web/API/Canvas_API/Tutorial/Basic_usage)

[Protobuf Guide](https://developers.google.com/protocol-buffers/docs/proto3)

[Vue.js](https://cn.vuejs.org/index.html)



게임에 오신 것을 환영합니다.

개념은 EVE 게임과 채팅방에서 가져온 것이지만 게임 코드는 완전히 완전히 구현됩니다.

방법:

1. W A S D는 상하좌우로 수행됩니다.

2. 공백은 채팅 상자를 열고 메시지를 보내기 위해 다시 전화합니다.

3. 왼쪽 상단 모서리에서 닉네임을 수정하고 빈 공간을 클릭하여 성공적으로 수정했습니다.

저자 GIT:https://github.com/sunshinev/go-space-chat.

프런트 엔드 Vue+canvas+websocket+protobuf.

백 엔드 Golang+websocket+protobuf+goroutine.

