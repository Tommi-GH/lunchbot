package listener

import (
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

var weekdays = map[string]string {
	"monday"    : "maanantai",
	"tuesday"   : "tiistai",
	"wednesday" : "keskiviikko",
	"thursday"  : "torstai",
	"friday"    : "perjantai"}


func init() {
	http.HandleFunc("/", handleMessage)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {

	if !strings.EqualFold(r.PostFormValue("token"), token) {
		http.Error(w, "Invalid token.", http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	w.Header().Set("content-type", "application/json")

	//escape problematic characters
	message := strings.Replace(strings.Replace(r.PostFormValue("text"), `"`, "´´", -1), "\\", "/", -1)

	//for making comparing easier
	message = strings.ToLower(message)

	//Get menu and create response of it
	resp, _ := createResponse(r, message)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorf(ctx, "Error encoding JSON: %s", err)
	}

}

//Creates the response for the initial POST-request. The response
//includes an ephemeral slack-message
func createResponse(r *http.Request, message string) (*slashResponse, bool) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	respMessage := ""
	v := url.Values{}
	v.Set("restaurant_id", "23")

	resp, err := client.PostForm("https://www.kanresta.fi/app/lunchlist/view/", v)
	if err != nil {
		log.Errorf(ctx, "Unable to get lunchsite: %s", err)
		respMessage = "Unable to get lunchlist"
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Errorf(ctx, "Unable to get lunchsite: %s", err)
		respMessage = "Unable to get lunchlist"
	}

	doc.Find(".lunchlist-day").Each(func(i int, s *goquery.Selection) {
		day := s.Children().First().Text()
		day = strings.TrimSpace(day)
		if checkWeekday(day,message) || len(message) == 0 {

			description := s.Find(".description").Text()
			description = strings.TrimSpace(description)

			respMessage = respMessage + day+"\n\n"+ description +"\n"
			respMessage = respMessage + "\n-----------------\n"
		}
	})

	if len(respMessage) == 0 {
		respMessage = "Unable to get lunchlist with given arguments. Please provide weekday in finnish, or no message for the whole week's menu"
	}

	resp.Body.Close()

	return &slashResponse{
		ResponseType: "ephemeral",
		Text:         respMessage,
	}, false
}

func checkWeekday (day string, message string) bool{

	if weekdays[message] != "" {
		message = weekdays[message]
	}

	return strings.EqualFold(strings.ToLower(day),message)

}