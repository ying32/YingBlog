package blog

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// 新建一个类型，继承sessions.CookieStore
type TSessions struct {
	session  *sessions.Session
	handlder *Handlder
}

var (
	sessionHelper *TSessions
	sessionStore  *sessions.CookieStore
)

func (self *TSessions) save() {
	self.session.Save(self.handlder.Request, self.handlder.ResponseWriter)
}

// 新建cookie存储仓
func newCookieStore(keyPairs ...[]byte) *sessions.CookieStore {
	cs := &sessions.CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   0, // 只存为会话，关闭则失效
			HttpOnly: true,
		},
	}

	cs.MaxAge(cs.Options.MaxAge)
	return cs
}
