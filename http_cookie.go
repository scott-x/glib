/*
* @Author: scottxiong
* @Date:   2021-06-16 20:39:13
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-06-22 15:17:16
 */
package glib

import (
	"net/http"
	"time"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	X_Session_ID = "X-Session-ID"
)

type Session struct {
	ID        string //session id
	User      string //user info
	LoginTime time.Time
	TTL       time.Duration
}

var session map[string]Session

func init() {
	session = make(map[string]Session, 0)
	initInMain() //设置时区
}

func initInMain() {
    var cstZone = time.FixedZone("CST", 8*3600) // 东八
    time.Local = cstZone
}

func uuid() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetSession(r *http.Request) (Session,error) {
	cookie_sid, err := r.Cookie(X_Session_ID)
	if err != nil {
		return Session{},err
	}

	sid := cookie_sid.Value
	return session[sid],nil
}

func NewDefaultSession(user string) Session {
	id,_ := uuid()
	ss := Session{
		ID:        id,
		User:      user,
		LoginTime: time.Now(),
		TTL:       time.Second * 60 * 5, //5分钟过期
	}
	session[id] = ss
	return ss
}

// const (
//     Nanosecond  Duration = 1
//     Microsecond          = 1000 * Nanosecond
//     Millisecond          = 1000 * Microsecond
//     Second               = 1000 * Millisecond
//     Minute               = 60 * Second
//     Hour                 = 60 * Minute
// )

//generate session
func NewSession(user string, ttl time.Duration) Session {
	id,_ := uuid()
	ss := Session{
		ID:        id,
		User:      user,
		LoginTime: time.Now(),
		TTL:       ttl,
	}
	session[id] = ss
	return ss
}

//set sessionID to cookie
func SetCookie(w http.ResponseWriter, sess Session) {
	expire := time.Now().Add(sess.TTL)
	cookie := http.Cookie{
		Name:     X_Session_ID,
		Value:    sess.ID, //usually sessionId
		Expires:  expire,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

// type Cookie struct {
//     Name  string
//     Value string

//     Path       string    // optional
//     Domain     string    // optional
//     Expires    time.Time // optional
//     RawExpires string    // for reading cookies only

//     // MaxAge=0 means no 'Max-Age' attribute specified.
//     // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
//     // MaxAge>0 means Max-Age attribute present and given in seconds
//     MaxAge   int
//     Secure   bool
//     HttpOnly bool
//     SameSite SameSite // Go 1.11
//     Raw      string
//     Unparsed []string // Raw text of unparsed attribute-value pairs
// }

//check whether user login or not?
func IsUserLogin(w http.ResponseWriter, r *http.Request, domain string) bool {
	cookie_sid, err := r.Cookie(X_Session_ID)
	if err != nil {
		return false
	}

	sid := cookie_sid.Value
	// log.Println("current sid:", sid)
	if len(sid) == 0 {
		return false
	}

	ok := isSessionExpired(sid)
	if ok {
		// log.Println("current sid 过期了")
		removeFrontCookie(w, domain)
		return false
	}
	// log.Println("current sid 没过期")
	return true
}

//check whether session is expired or not?
func isSessionExpired(sid string) bool {
	if sess, ok := session[sid]; ok {
		return time.Now().After(sess.LoginTime.Add(sess.TTL)) //t1不变，t2是变量
	} else {
		return false
	}
}

func removeFrontCookie(w http.ResponseWriter, domain string) {
	if domain == "" {
		domain = "127.0.0.1"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     X_Session_ID,
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour), // Set expires for older versions of IE
		Domain:   domain,
		Path:     "/",
		HttpOnly: true,
	})
}
