package session

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionInfo struct {
	Mail interface{}
}

var LoginInfo SessionInfo

func Login(c *gin.Context, mail string) {

	//セッションにデータを格納する
	session := sessions.Default(c)
	session.Set("Mail", mail)
	if err := session.Save(); err != nil {
		log.Fatal(err)
	}
}

func Logout(c *gin.Context) {

	//セッションからデータを破棄する
	session := sessions.Default(c)
	log.Println("セッション取得")
	session.Clear()
	log.Println("クリア処理")
	if err := session.Save(); err != nil {
		log.Fatal(err)
	}
}

func SessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("SessionCheck")
		session := sessions.Default(c)
		LoginInfo.Mail = session.Get("Mail")
		fmt.Println(LoginInfo.Mail)
		// セッションが無い場合、ログインフォームを出す
		if LoginInfo.Mail == nil {
			log.Println("ログインしていません")
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort()
		} else {
			c.Set("Mail", LoginInfo.Mail)
			c.Next()
		}
		log.Println("ログインチェック終了")
	}
}
