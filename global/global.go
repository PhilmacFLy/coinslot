package global

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/securecookie"

	"golang.org/x/crypto/scrypt"
)

var Execdir string
var Conf Config

const saltSize = 32

var sc = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func BuildMessage(template string, message string) string {
	message = html.EscapeString(message)
	return strings.Replace(template, "$MESSAGE$", message, -1)
}

func SetCookie(w http.ResponseWriter, u string) error {

	value := map[string]string{
		"name": u,
	}

	encoded, err := sc.Encode("coinslot", value)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  "coinslot",
		Value: encoded,
		Path:  "/",
	}
	cookie.Expires = time.Now().Add(10 * 365 * 24 * time.Hour)

	http.SetCookie(w, cookie)
	return nil
}

func GetCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("coinslot")

	if err != nil {
		return "", err
	}
	value := make(map[string]string)
	err = sc.Decode("coinslot", cookie.Value, &value)

	if err != nil {
		return "", err
	}
	return value["name"], nil

}

func RemoveCookie(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)

	cookieMonster := &http.Cookie{
		Name:    "coinslot",
		Expires: expire,
		Value:   "",
	}

	// http://golang.org/pkg/net/http/#SetCookie
	// add Set-Cookie header
	http.SetCookie(w, cookieMonster)
}

type DBConnection struct {
	Driver     string
	Connection string
}

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	Execdir = dir + "/"
	Conf.load()
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type Config struct {
	Port       int
	Connection DBConnection
}

func (c *Config) load() error {
	filename := Execdir + "config.json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePasswordHash(p, s string) (string, error) {
	dk, err := scrypt.Key([]byte(p), []byte(s), 16384, 12, 1, 32)
	if err != nil {
		return "", err
	}
	h := hex.EncodeToString(dk)
	return h, err
}

func GenerateSalt() (string, error) {
	buf := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, buf)
	return string(buf), err
}
