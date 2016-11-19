package aegateway

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const(
	ConfigPath = "./gateway.yaml"
)

func init() {
	config := LoadConfig(ConfigPath)

	for _, c := range config.Routes {
		handleGatewayRequest(c)
	}

	http.HandleFunc("/", handler)
}

func handleGatewayRequest(c GatewayRoute) {
	http.HandleFunc(c.Pattern, func (w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		client := urlfetch.Client(ctx)

		dstUrl := strings.Replace(r.URL.Path, c.Pattern, c.Dest, 1)

		log.Debugf(ctx, "request to %s %s", r.Method, dstUrl)
		req, err := http.NewRequest(r.Method, dstUrl, r.Body)
		if err != nil {
			log.Warningf(ctx, "failed to make request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Warningf(ctx, "failed to process request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warningf(ctx, "failed to read upstream response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
		if _, err := w.Write(body); err != nil {
			log.Warningf(ctx, "failed to send response to downstream: %v", err)
		}
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	config := LoadConfig(ConfigPath)
	fmt.Fprintf(w, "404 Not Found %s %+v", r.URL, config)
}
