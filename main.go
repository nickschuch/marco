package main
 
import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
)
 
func main() {
        remote, err := url.Parse("http://172.17.0.2:2368")
        if err != nil {
                panic(err)
        }
 
        proxy := httputil.NewSingleHostReverseProxy(remote)
        http.HandleFunc("/", handler(proxy))
        err = http.ListenAndServe(":80", nil)
        if err != nil {
                panic(err)
        }
}
 
func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                log.Println(r.URL)
                w.Header().Set("X-Ben", "Rad")
                p.ServeHTTP(w, r)
        }
}

