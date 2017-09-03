package listener

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func init() {
	http.HandleFunc("/", handleMessage)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	w.Header().Set("content-type", "application/json")

	//escape problematic characters
	message := strings.Replace(strings.Replace(r.PostFormValue("text"), `"`, "´´", -1), "\\", "/", -1)

	//If the request is a valid report, do the following steps,
	//else return appropriate error-message
	resp, _ := createResponse(r, message)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorf(ctx, "Error encoding JSON: %s", err)
	}

}

//Creates the response for the initial POST-request. The response
//includes a slack-message
func createResponse(r *http.Request, message string) (*slashResponse, bool) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	respMessage := ""
	resp, err := client.Get("http://eatwork.fi/tilat/panuntie/")

	if err != nil {
		log.Errorf(ctx, "Unable to get lunchsite: %s", err)
		respMessage = "Unable to get lunchlist"
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		log.Errorf(ctx, "Unable to get lunchsite: %s", err)
		respMessage = "Unable to get lunchlist"
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && strings.HasPrefix(n.Data, "Lounaslista") {
			respMessage = n.Data + "\n"
		} else if n.Type == html.TextNode && strings.HasPrefix(n.Data, "Maanantai") {
			respMessage = respMessage + n.Data
			respMessage = respMessage + n.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
		} else if n.Type == html.TextNode && strings.HasPrefix(n.Data, "Tiistai") {
			respMessage = respMessage + n.Data
			respMessage = respMessage + n.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
		} else if n.Type == html.TextNode && strings.HasPrefix(n.Data, "Keskiviikko") {
			respMessage = respMessage + n.Data
			respMessage = respMessage + n.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
		} else if n.Type == html.TextNode && strings.HasPrefix(n.Data, "Torstai") {
			respMessage = respMessage + n.Data
			respMessage = respMessage + n.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
		} else if n.Type == html.TextNode && strings.HasPrefix(n.Data, "Perjantai") {
			respMessage = respMessage + n.Data
			respMessage = respMessage + n.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
			respMessage = respMessage + n.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	resp.Body.Close()

	return &slashResponse{
		ResponseType: "in_channel",
		Text:         respMessage,
	}, false
}
