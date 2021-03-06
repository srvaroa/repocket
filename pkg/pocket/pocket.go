package pocket

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// https://getpocket.com/developer/docs/v3/retrieve
type retrieveResponse struct {
	Status int
	List   map[string]Article
}

// https://getpocket.com/developer/docs/v3/send
type sendResponse struct {
	ActionResults bool
	Status        int
}

// https://getpocket.com/developer/docs/v3/retrieve
type Article struct {
	ItemId        string `json:"item_id"`
	ResolvedId    string `json:"resolved_id"`
	GivenUrl      string `json:"given_url"`
	GivenTitle    string `json:"given_title"`
	Favorite      bool   `json:"favorite"`
	Status        int    `json:"status"` // 0, 1, 2 where 1 = archived, 2 = to delete
	ResolvedTitle string `json:"resolved_title"`
	ResolvedUrl   string `json:"resolved_url"`
	Excerpt       string `json:"excerpt"`
	IsArticle     bool   `json:"is_article"`
	HasVideo      int    `json:"has_video"` // 0, 1, 2 where 1 = has videos, 2 = is a video
	HasImage      int    `json:"has_image"` // 0, 1, 2 where 1 = has images, 2 = is an image
	WordCount     int    `json:"word_count"`
	Tags          string `json:"tags"`    // actually another object but I care not now
	Authors       string `json:"authors"` // actually another object but I care not now
	Images        string `json:"images"`  // actually another object but I care not now
	Videos        string `json:"videos"`  // actually another object but I care not now
}

const apiUrl = "https://getpocket.com/v3"

const (
	STATE_ALL     = "all"
	STATE_UNREAD  = "unread"
	STATE_ARCHIVE = "archive"
)

const (
	ACTION_DELETE   = "delete"
	ACTION_FAVORITE = "favorite"
	ACTION_ARCHIVE  = "archive"
)

func QueryFavourites(accessToken, consumerKey string) map[string]Article {
	return query(accessToken, consumerKey, STATE_ALL, 1, 0)
}

func QueryNewest(accessToken, consumerKey string, count int) map[string]Article {
	return query(accessToken, consumerKey, STATE_UNREAD, 0, count)
}

func QueryUnread(accessToken, consumerKey string) map[string]Article {
	return query(accessToken, consumerKey, STATE_UNREAD, 0, 0)
}

func query(accessToken, consumerKey, state string, favourites, count int) map[string]Article {

	payload := map[string]interface{}{
		"access_token": accessToken,
		"consumer_key": consumerKey,
		"favorite":     favourites,
		"detailType":   "complete",
		"sort":         "newest",
		"state":        state,
	}

	if count > 0 {
		payload["count"] = count
	}

	data, _ := json.Marshal(payload)

	res, err := http.Post(apiUrl+"/get",
		"application/json",
		strings.NewReader(string(data)))

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Unable to retrieve items", err)
	}

	var retrieved retrieveResponse

	json.Unmarshal(body, &retrieved)

	return retrieved.List

}

func Archive(accessToken, consumerKey string, itemIds []string) bool {
	return action(ACTION_ARCHIVE, accessToken, consumerKey, itemIds)
}

func Delete(accessToken, consumerKey string, itemIds []string) bool {
	return action(ACTION_DELETE, accessToken, consumerKey, itemIds)
}

func Fav(accessToken, consumerKey string, itemIds []string) bool {
	return action(ACTION_FAVORITE, accessToken, consumerKey, itemIds)
}

func action(action, accessToken, consumerKey string, itemIds []string) bool {

	timestamp := time.Now().UTC()
	var actions []map[string]interface{}

	if len(itemIds) <= 0 {
		return true
	}

	for _, itemId := range itemIds {
		actions = append(actions, map[string]interface{}{
			"action":    action,
			"item_id":   itemId,
			"timestamp": timestamp,
		})
	}

	payload := map[string]interface{}{
		"access_token": accessToken,
		"consumer_key": consumerKey,
		"actions":      actions,
	}

	data, _ := json.Marshal(payload)

	res, err := http.Post(apiUrl+"/send",
		"application/json",
		strings.NewReader(string(data)))

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Unable to delete items", err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("Error deleting items %v \n", res.Header)
	}

	var deleted sendResponse

	json.Unmarshal(body, &deleted)

	return deleted.ActionResults

}

// Authorize returns the token for the given consumer key by firing
// GetPocket's auth process.
func Authorize(consumerKey string) (string, error) {

	log.Printf("Fetching token for consumer key: %s", consumerKey)

	log.Print("Initiating OAuth process..")
	res, err := http.PostForm(apiUrl+"/oauth/request", url.Values{
		"consumer_key": {consumerKey},
		"redirect_uri": {"localhost"},
	})

	if err != nil {
		log.Fatal("Unable to authorise", err)
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Unable to retrieve code", err)
		return "", err
	}

	code := strings.Split(string(body), "=")[1]

	log.Print("Authorizing application...")
	log.Printf("Browse to this URL, you may ignore errors: "+
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=localhost"+
		"\n\nPress enter when done",
		code)

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	res, err = http.PostForm(apiUrl+"/oauth/authorize", url.Values{
		"consumer_key": {consumerKey},
		"code":         {code},
	})

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Unable to retrieve token", err)
		return "", err
	}

	theBody := string(body)
	parts := strings.Split(theBody, "&")
	if len(parts) != 2 {
		log.Fatalf("Unexpected final autorization response "+
			"expecting access_token=<token>&username=<username> but got %s",
			theBody)
	}
	log.Printf("Authorized as %s", strings.Split(string(parts[1]), "=")[1])
	return strings.Split(string(parts[0]), "=")[1], nil
}
