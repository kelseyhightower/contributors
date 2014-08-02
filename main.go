package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"
)

var (
	client *http.Client
	pool   *x509.CertPool
)

var html = template.Must(template.New("html").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Contributors</title>
  </head>
  <body>
    <h1>Contributors</h1>
    <form action="/">
      <label for="owner">owner:</label>
      <input type="text" name="owner" id="owner"/>
      <label for="repo">repo:</label>
      <input type="text" name="repo" id="repo"/>
      <input type="submit" value="Submit">
    </form>
    <p>{{.ErrorMessage}}</p>
    <table>{{range .Contributors}}
      <tr>
        <td>{{.Author.Login}}</td>
        <td><img width=100 height=100 src="{{.Author.AvatarUrl}}"><td>
      </tr>{{end}}
    </table>
  </body>
</html>
`))

type Result struct {
	Contributors Contributors
	ErrorMessage string
}

type Contributors []Contributor

type Contributor struct {
	Author Author `json:"author"`
}

type Author struct {
	AvatarUrl string `json:"avatar_url"`
	Login     string `json:"login"`
}

func init() {
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}
}

func contributors(w http.ResponseWriter, req *http.Request) {
	var result Result
	owner := req.FormValue("owner")
	repo := req.FormValue("repo")
	if owner != "" && repo != "" {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/stats/contributors", owner, repo)
		resp, err := client.Get(url)
		if err != nil {
			result.ErrorMessage = err.Error()
			goto L
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				result.ErrorMessage = err.Error()
				goto L
			}
			err = json.Unmarshal(data, &result.Contributors)
			if err != nil {
				result.ErrorMessage = err.Error()
				goto L
			}
		}
		if resp.StatusCode == http.StatusNotFound {
			result.ErrorMessage = fmt.Sprintf("%s/%s not found.", owner, repo)
		}
	}
L:
	err := html.Execute(w, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	hostPort := net.JoinHostPort("0.0.0.0", os.Getenv("PORT"))
	http.HandleFunc("/", contributors)
	log.Println("Starting contributors app")
	log.Printf("Listening on %s\n", hostPort)
	if err := http.ListenAndServe(hostPort, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
