/*
* @Author: scottxiong
* @Date:   2021-06-16 20:39:13
* @Last Modified by:   scottxiong
* @Last Modified time: 2021-06-19 17:40:15
 */
package glib

import (
	"net/http"
	"time"
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
}

func getSession(r *http.Request) Session {
	cookie_sid, err := r.Cookie(X_Session_ID)
	if err != nil {
		return nil
	}

	sid := cookie_sid.Value
	return session[sid]
}

func NewDefaultSession(user string) Session {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil
	}
	id := base64.URLEncoding.EncodeToString(b)
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
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil
	}
	id := base64.URLEncoding.EncodeToString(b)
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
func isUserLogin(w http.ResponseWriter, r *http.Request, domain string) bool {
	cookie_sid, err := r.Cookie(X_Session_ID)
	if err != nil {
		return false
	}

	sid := cookie_sid.Value
	log.Println("current sid:", sid)
	if len(sid) == 0 {
		return false
	}

	uname, ok := isSessionExpired(sid)
	if ok {
		// log.Println("current sid 过期了")
		removeFrontCookie(w, domain)
		return false
	}
	log.Println("current sid 没过期")
	return true
}

//check whether session is expired or not?
func isSessionExpired(sid string) bool {
	if sess, ok := m[sid]; ok {
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
