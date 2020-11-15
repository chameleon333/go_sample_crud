package session

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func Login(c *gin.Context, mail string) {

	//セッションにデータを格納する
	session := sessions.Default(c)
	session.Set("mail", mail)
	if err := session.Save(); err != nil {
		log.Fatal(err)
	}
}
