package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	pb "go-space-chat/proto/star"
	"io"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"

	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"

	filter "github.com/antlinker/go-dirtyfilter"

	"github.com/antlinker/go-dirtyfilter/store"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

/* ================================================ */

var db *sql.DB
var err error

// func signupPage(res http.ResponseWriter, req *http.Request) {
// 	if req.Method != "POST" {
// 		http.ServeFile(res, req, "signup.html")
// 		return
// 	}

// 	username := req.FormValue("username")
// 	password := req.FormValue("password")

// 	var user string

// 	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

// 	switch {
// 	case err == sql.ErrNoRows:
// 		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 		if err != nil {
// 			http.Error(res, "Server error, unable to create your account.", 500)
// 			return
// 		}

// 		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
// 		if err != nil {
// 			http.Error(res, "Server error, unable to create your account.", 500)
// 			return
// 		}

// 		res.Write([]byte("User created!"))
// 		return
// 	case err != nil:
// 		http.Error(res, "Server error, unable to create your account.", 500)
// 		return
// 	default:
// 		http.Redirect(res, req, "/", 301)
// 	}
// }

func loginPage(res http.ResponseWriter, req *http.Request) {
	log.Printf("loginPage.") //loginPage
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	log.Printf("username: %s", username) //loginPage
	var databaseUsername string
	var databasePassword string
	//SELECT * FROM users WHERE username ='"+param_username+"' and password ='"+md5(param_password)+"'
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	res.Write([]byte("Hello" + databaseUsername))

}

// func homePage(res http.ResponseWriter, req *http.Request) {
// 	http.ServeFile(res, req, "index.html")
// }

/* ================================================ */

// 客户端集合 //클라이언트 컬렉션입니다.
var clients = make(map[*websocket.Conn]*pb.BotStatusRequest)

var clientsMutex sync.RWMutex

var connMutex sync.RWMutex

// 消息缓冲通道 //메시지 버퍼 채널입니다.
var messages = make(chan *pb.BotStatusRequest, 100)

var socketAddr = flag.String("socketAddr", ":9000", "socket address")
var webAddr = flag.String("webAddr", ":8081", "http service address") //80

var upgrader = websocket.Upgrader{}

var words = []string{}

var filterManage *filter.DirtyManager

func main() {
	flag.Parse()

	filterManage = readWords()

	go http.ListenAndServe(*webAddr, http.FileServer(http.Dir("web_resource/dist/")))

	//log.Printf("web 服务启动成功 端口 %s", *webAddr)
	log.Printf("web 서비스가 성공적으로 시작된 포트입니다. %s", *webAddr)

	//####################################################################################
	//add main
	db, err = sql.Open("mysql", "")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	// http.HandleFunc("/signup", signupPage)
	// http.HandleFunc("/login", loginPage)
	// http.HandleFunc("/", homePage)
	// http.ListenAndServe(":8080", nil)
	//####################################################################################

	http.HandleFunc("/ws", echo)
	// 广播 //방송
	go boardcast()

	// pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//log.Printf("socket 服务启动端口 %s", *socketAddr)
	log.Printf("socket 서비스 시작 포트입니다. %s", *socketAddr)
	// 这里的 ListenAndServe 已经 a 开启了 goroutine 协程了
	// ListenAndServe는 이미 goroutine 코브를 시작했습니다.
	err := http.ListenAndServe(*socketAddr, nil)
	if err != nil {
		log.Fatalf("create error %v", err)
	}

}

// 这个 echo 是在 serve 协程里面运行的 //이 에코는 server 에서 실행됩니다.
func echo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Printf("POST로 받은 DATA입니다.")
		fmt.Print("POST\n")
		res, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%s\n", res)
		log.Printf("%s\n", res)
	}

	// 跨域 //도메인 간.
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	// 升级http为websocket //http를 websocket으로 업그레이드합니다.
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		//log.Printf("升级webcoket %v", err)
		log.Printf("webcoket를 업그레이드합니다. %v", err)
		w.Write([]byte(err.Error()))
	} else {
		go checkConn(c)
	}
}

func checkConn(c *websocket.Conn) {
	defer func() {
		c.Close()
	}()
	// 监听 //수신 대기.
	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			//log.Printf("read error 读取失败 ", err)
			log.Printf("read error 읽기에 실패했습니다. %v", err)
			messages <- &pb.BotStatusRequest{
				BotId: clients[c].BotId,
				// 广播关闭连接 //브로드캐스트가 연결을 닫습니다.
				Status: pb.BotStatusRequest_close,
			}

			clientsMutex.Lock()
			// 清除连接 //연결을 지웁니다.
			delete(clients, c)
			clientsMutex.Unlock()

			break
		}

		// 使用 protobuf 解析 //protobuf 구문 분석.
		pbr := &pb.BotStatusRequest{}
		err = proto.Unmarshal(message, pbr)
		if err != nil {
			log.Printf("proto 파싱 ​​실패 %v", err) //解析失败
			break
		}

		// 敏感词过滤 //민감한 단어 필터링.
		pbr.Msg = wordsFilter(pbr.Msg)
		pbr.Name = wordsFilter(pbr.Name)

		// 初始化链接的id //초기화 링크의 ID입니다.
		if clients[c] == nil {

			clientsMutex.Lock()
			clients[c] = &pb.BotStatusRequest{
				BotId:  pbr.GetBotId(),
				Status: pb.BotStatusRequest_connecting,
			}
			clientsMutex.Unlock()
		}

		messages <- pbr
	}
}

func boardcast() {
	// 始终读取 messages //항상 messages를 읽습니다.
	for msg := range messages {
		if msg.Msg != "" {
			log.Printf("%s : %s", msg.BotId+":"+msg.Name, msg.Msg)
		}
		// 读取到之后进行广播，启动协程，是为了立即处理下一条 msg
		// 읽기 후 브로드캐스트하고 코스타를 시작하는 것은 다음 msg를 즉시 처리하기 위한 것입니다.
		go func() {
			clientsMutex.RLock()
			defer clientsMutex.RUnlock()
			for cli := range clients {
				// protobuf 协议 //protobuf 프로토콜.
				if clients[cli].BotId == msg.BotId {
					continue
				}

				pbrp := &pb.BotStatusResponse{BotStatus: []*pb.BotStatusRequest{msg}}
				b, err := proto.Marshal(pbrp)
				if err != nil {
					log.Printf("proto marshal error %v", err)
					continue
				}

				// 二进制发送 //이진 전송.
				connMutex.Lock()
				err = cli.WriteMessage(websocket.BinaryMessage, b)
				connMutex.Unlock()
				if err != nil {
					log.Printf("%v", err)
				}
			}
		}()
	}
}

func readWords() *filter.DirtyManager {
	fi, err := os.Open("words/gg.txt")
	if err != nil {
		panic(err.Error())
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		words = append(words, string(a))
	}

	memStore, err := store.NewMemoryStore(store.MemoryConfig{
		DataSource: words,
	})
	if err != nil {
		panic(err)
	}
	return filter.NewDirtyManager(memStore)
}

func wordsFilter(filterText string) string {

	result, err := filterManage.Filter().Filter(filterText, '*', '@')
	if err != nil {
		panic(err)
	}

	if result != nil {
		for _, w := range result {
			filterText = strings.ReplaceAll(filterText, w, "*")
		}
	}

	return filterText
}

// func gfuncDbChk() {
// 	// sql.DB 객체 생성
// 	db, err := sql.Open("mysql", "root:pwd@tcp(127.0.0.1:3306)/citypay")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// 복수 Row를 갖는 SQL 쿼리
// 	var id int
// 	var name string
// 	rows, err := db.Query("SELECT id, username, nick,CTP_address FROM users where id >= ?", 2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

// 	for rows.Next() {
// 		err := rows.Scan(&id, &name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(id, name)
// 	}
// }
