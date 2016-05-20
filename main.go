package main

import (
	"log"
	"net/http"
    "text/template"
    "path/filepath"
    "sync"
    "flag"
    "os"
    "trace"
    
    "github.com/stretchr/gomniauth"
    "github.com/stretchr/gomniauth/providers/facebook"
    "github.com/stretchr/gomniauth/providers/github"
    "github.com/stretchr/gomniauth/providers/google"
    "github.com/stretchr/objx"
)

// templ represents a single template
type templateHandler struct {
    once        sync.Once
    filename    string
    templ       *template.Template
}
// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    
    data := map[string]interface{}{
        "Host": r.Host,
    }
    if authCookie, err := r.Cookie("auth"); err == nil {
        // fmt.Println(objx.MustFromBase64(authCookie.Value));
        data["UserData"] = objx.MustFromBase64(authCookie.Value)
    }
    
    t.templ.Execute(w, data)
} 

var addr = flag.String("host", ":8080", "The host of the application")

func main() {
    flag.Parse() // parse the flags
    
    // setup gomniauth
    gomniauth.SetSecurityKey("milkisgoodwhennotcurded")
    gomniauth.WithProviders(
        facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
        github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
        google.New("892920471098-83at5o0oq1pmhnamgi0qoj7khhqq4k4r.apps.googleusercontent.com",
            "dD8OYGAOs8pCyRhsn70HoHSm", "http://localhost:8080/auth/callback/google"),
    )
    
    // r := newRoom(UseAuthAvatar)
    // r := newRoom(UseGravatar)
    r := newRoom(UseFileSystemAvatar)
    r.tracer = trace.New(os.Stdout)
    // root
	//http.Handle("/", &templateHandler{filename: "chat.html"})
    http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
    http.Handle("/login", &templateHandler{filename: "login.html"})
    http.Handle("/upload", &templateHandler{filename: "upload.html"})
    http.HandleFunc("/auth/", loginHandler)
    http.Handle("/room", r)
    http.Handle("/avatars/", 
        http.StripPrefix("/avatars/", 
            http.FileServer(http.Dir("./avatars"))))
    http.HandleFunc("/uploader", uploaderHandler)
    http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request){
        http.SetCookie(w, &http.Cookie{
            Name:   "auth",
            Value:  "",
            Path:   "/",
            MaxAge: -1,
        })
        w.Header()["Location"] = []string{"/chat"}
        w.WriteHeader(http.StatusTemporaryRedirect)
    })
    // get the room going
    go r.run()
	// start the web server
    log.Println("Starting the web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
