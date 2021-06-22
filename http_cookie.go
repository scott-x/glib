/*
* @Author: scottxiong
* @Date:   2021-06-16 20:39:13
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-06-22 22:20:53
 */
package glib

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type Session struct {
	ID        string //session id
	User      string //user info
	LoginTime time.Time
	TTL       time.Duration
}

var (
	X_Session_ID = "X_Sess_Scott"
	session map[string]Session
	expires = time.Hour*12 //12 hours
)

//set cookie name
func SetCookieName(key string) {
	X_Session_ID = key
}

//set expired time
func SetExpireTime(h time.Duration) {
	expires = h
}

func init() {
	session = make(map[string]Session, 0)
}

func GetSession(r *http.Request) (Session, error) {
	cookie_sid, err := r.Cookie(X_Session_ID)
	if err != nil {
		log.Println("GetSession:", err)
		return Session{}, err
	}

	sid := cookie_sid.Value
	if sess, ok := session[sid]; ok { //map不一定可以取到值
		return sess, nil
	}
	for k,v := range sess{
		fmt.Printf("[k,v]:(%s,%v)\n",k,v)
	}

	return Session{}, errors.New("Session not exists")

}

func NewDefaultSession(user string) Session {
	id, _ := uuid()
	ss := Session{
		ID:        id,
		User:      user,
		LoginTime: time.Now(),
		TTL:       expires, //12 hours 过期
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
	id, _ := uuid()
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
		Secure:   false,
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
		log.Println(err)
		return false
	}

	sid := cookie_sid.Value
	log.Println("current sid:", sid)
	if len(sid) == 0 {
		return false
	}

	ok := isSessionExpired(sid)
	if ok {
		log.Println("current sid 过期了")
		RemoveFrontCookie(w, domain)
		delete(session, sid)
		log.Println("已删除过期的sid")
		return false
	}
	// log.Println("current sid 没过期")
	return true
}

//check whether session is expired or not?
func isSessionExpired(sid string) bool {
	if sess, ok := session[sid]; ok {
		return time.Now().After(sess.LoginTime.Add(sess.TTL)) //t1不变，t2是变量
	}
	//if can't get sess, that means use not login
	return true
}

func RemoveFrontCookie(w http.ResponseWriter, domain string) {
	log.Println("正在删除前端的cookie")
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
		Secure:   false,
	})
}
