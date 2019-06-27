package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/n3wscott/graph/pkg/graph"
)

var once sync.Once
var t *template.Template

func getQueryParam(r *http.Request, key string) string {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return ""
	}
	return keys[0]
}

func (c *Controller) RootHandler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		t, _ = template.ParseFiles(
			c.root+"/templates/index.html",
			c.root+"/templates/main.html",
		)
	})

	c.GraphHTML(w, r)
}

// TODO: support just fetching the graph image

var defaultFormat = "svg"    // or png
var defaultFocus = "trigger" // or png

func (c *Controller) GraphHTML(w http.ResponseWriter, r *http.Request) {
	format := getQueryParam(r, "format")
	if format == "" {
		format = defaultFormat
	}

	focus := getQueryParam(r, "focus")
	if focus == "" {
		focus = defaultFocus
	}

	var dotGraph string

	switch focus {
	case "sub", "subs", "subscription", "subscriptions":
		dotGraph = graph.ForSubscriptions(c.client, c.namespace)
	case "broker", "trigger", "triggers":
		fallthrough
	default:
		dotGraph = graph.ForTriggers(c.client, c.namespace)
	}

	file, err := dotToImage(format, []byte(dotGraph))
	if err != nil {
		log.Printf("dotToImage error %s", err)
		return
	}
	img, err := ioutil.ReadFile(file)

	defer os.Remove(file) // clean up

	var data map[string]interface{}
	if format == "svg" {
		data = map[string]interface{}{
			"Image":  img,
			"Format": format,
		}
	} else {
		data = map[string]interface{}{
			"Image":  base64.StdEncoding.EncodeToString(img),
			"Format": fmt.Sprintf("image/%s;base64", format),
		}
	}
	_ = t.Execute(w, data)
}

var dot string

func dotToImage(format string, b []byte) (string, error) {
	if dot == "" {
		var err error
		dot, err = exec.LookPath("dot")
		if err != nil {
			log.Fatalln("unable to find program 'dot', please install it or check your PATH")
		}
	}

	var img = filepath.Join(os.TempDir(), fmt.Sprintf("graph.%s", format))

	cmd := exec.Command(dot, fmt.Sprintf("-T%s", format), "-o", img)
	cmd.Stdin = bytes.NewBuffer(b)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return img, nil
}
