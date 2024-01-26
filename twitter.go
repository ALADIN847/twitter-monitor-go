package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type TwitterMonitor struct {
	MonitoredPath string
	RecentPath    string
	Token         string
	Cookies       string
	CSRF          string
}

type User struct {
	Name string
	ID   string
}

type TweetInfo struct {
	CreatedAt string `json:"created_at"`
	IDStr     string `json:"id_str"`
}

type Tweet struct {
	CreatedAt string `json:"created_at"`
	IDStr     string `json:"id_str"`
}

type Entities struct {
	UserMentions []UserMention `json:"user_mentions"`
}

type UserMention struct {
	ScreenName string `json:"screen_name"`
}

type ExtendedEntities struct {
	Media []Media `json:"media"`
}

type Media struct {
	Type      string    `json:"type"`
	VideoInfo VideoInfo `json:"video_info"`
	MediaURLs []string  `json:"media_url_https"`
}

type VideoInfo struct {
	Variants []VideoVariant `json:"variants"`
}

type VideoVariant struct {
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []byte)
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			delete(clients, conn)
			return
		}

		fmt.Printf("Received message: %s\n", message)

		broadcast <- message
	}
}
func sendToAllClients(message []byte) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Error sending message to client:", err)
		}
	}
}

func asyncTask() error {
	// Simulate some async task
	fmt.Println("Async task executed")
	twitterMonitor := TwitterMonitor{
		MonitoredPath: "./monitored.json",
		RecentPath:    "./recent.json",
		Token:         "Bearer AAAAAAAAAAAAAAAAAAAAAF7aAAAAAAAASCiRjWvh7R5wxaKkFp7MM%2BhYBqM%3DbQ0JPmjU9F6ZoMhDfI4uTNAaQuTDm2uO9x3WFVr2xBZ2nhjdP0",
		Cookies:       "...", // Your cookies here
		CSRF:          "...", // Your CSRF token here
	}
	users := []User{
		{
			Name: "LordOfSavings",
			ID:   "1327342638413082626",
		},
		{
			Name: "LordOfDiscounts",
			ID:   "1374247110313439236",
		},
		{
			Name: "Tracker_Deals",
			ID:   "1341108069812453377",
		},
	}

	for _, user := range users {
		fmt.Printf("Monitoring \nUser Name: %s\nUser ID: %s\n", user.Name, user.ID)
		data, err := twitterMonitor.monitorUser(user)
		if err != nil {
			continue
		}
		tweetInfo, exists := data["tweet"].(map[string]interface{})
		if !exists {

			continue
		}
		if len(tweetInfo) == 0 {

		} else {
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				continue
			}
			sendToAllClients(jsonData)

			fmt.Println("Send Event here")
		}
	}

	return nil
}

func main() {

	http.HandleFunc("/ws", wsHandler)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	interval := 1 // Interval in seconds
	go func() {
		for {
			err := asyncTask()
			if err != nil {
				fmt.Println("Async task error:", err)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()
	select {}

}

func addRecent(user string, jsonData map[string]interface{}) error {

	jsonFile := "./recent.json"
	existingData := make(map[string]interface{})

	// Read existing JSON data
	existingJSON, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	// Parse existing JSON data
	if err := json.Unmarshal(existingJSON, &existingData); err != nil {
		return err
	}

	// Update existing JSON data with new user and json_data
	existingData[user] = jsonData

	// Convert updated JSON data to bytes
	updatedJSON, err := json.Marshal(existingData)
	if err != nil {
		return err
	}

	// Write updated JSON data back to the file
	if err := ioutil.WriteFile(jsonFile, updatedJSON, 0644); err != nil {
		return err
	}

	return nil
}

func (tm *TwitterMonitor) getRecent(user User) (map[string]interface{}, error) {
	endpoint := tm.formEndpoint(user.ID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("authority", "twitter.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("cookie", "d_prefs=MToxLGNvbnNlbnRfdmVyc2lvbjoyLHRleHRfdmVyc2lvbjoxMDAw; tweetdeck_version=beta; guest_id=v1%3A169055361674488664; guest_id_marketing=v1%3A169055361674488664; guest_id_ads=v1%3A169055361674488664; kdt=WopVvg4X2LS78w1533FHQyRZ7foywJScxW5eiDR8; auth_token=f869476a35dd9627696a588ac85df8cd23c13145; ct0=0edd1df5459754ec6f39541544783376a64ce4e2c4f9e4533ba929ecb6b047693d9477dfcc4c05ebbc78e649e398519a3cde0f4be012241d3dc6433cfba630919122a5806993ccf84f198cf926084fc9; twid=u%3D1480016932728365061; personalization_id=\"v1_MILqp2Zi9QZitO7e8Y9oww==\"; mbox=session#c10eb03a740e43aebb001d648a53af20#1691684394|PC#c10eb03a740e43aebb001d648a53af20.37_0#1754927334; _ga_34PHSZMC42=GS1.1.1691682535.1.1.1691682598.0.0.0; lang=en; _ga=GA1.2.1643705879.1688688782; _gid=GA1.2.1495936914.1691981985; external_referer=padhuUp37zjgzgv1mFWxJ12Ozwit7owX|0|8e8t2xd8A2w%3D; guest_id=v1%3A169168330739615608; guest_id_ads=v1%3A169168330739615608; guest_id_marketing=v1%3A169168330739615608; personalization_id=\"v1_AUcAJyaOdzCjxFs0kcRuhw==\"")
	req.Header.Add("dnt", "1")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("referer", "https://twitter.com/Emmet_Finance")
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"115\", \"Chromium\";v=\"115\"")
	req.Header.Add("sec-ch-ua-mobile", "?1")
	req.Header.Add("sec-ch-ua-platform", "\"Android\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Mobile Safari/537.36")
	req.Header.Add("x-client-transaction-id", "D0wLLo0ivqxGlsbovJu6sShquwGibQodd2oRQJFUPzxRcQvVcF2Bk6Jya3om1Tb/q2AShA/no3WJcwZo7EO+JfbtXDzzDg")
	req.Header.Add("x-csrf-token", "0edd1df5459754ec6f39541544783376a64ce4e2c4f9e4533ba929ecb6b047693d9477dfcc4c05ebbc78e649e398519a3cde0f4be012241d3dc6433cfba630919122a5806993ccf84f198cf926084fc9")
	req.Header.Add("x-twitter-active-user", "yes")
	req.Header.Add("x-twitter-auth-type", "OAuth2Session")
	req.Header.Add("x-twitter-client-language", "en")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []map[string]interface{}
	err2 := json.NewDecoder(resp.Body).Decode(&data)
	if err2 != nil {
		fmt.Println("Error decoding JSON:", err)

	}

	if len(data) == 0 {
		fmt.Printf("tweet undefined")
		return nil, nil
	} else {
		result := map[string]interface{}{
			"user":  "123",
			"Tweet": data[0],
		}

		return result, nil
	}

	// for _, item := range data {
	// 	fmt.Println("Created At:", item["created_at"])
	// 	fmt.Println("ID:", item["id"])
	// 	fmt.Println("Full Text:", item["full_text"])

	// 	// Access other fields similarly...
	// }

}

func (tm *TwitterMonitor) formEndpoint(userName string) string {
	return fmt.Sprintf("https://api.twitter.com/1.1/statuses/user_timeline.json?count=40&include_rts=0&user_id=%s&cards_platform=Web-13&include_entities=1&include_user_entities=1&include_cards=1&send_error_codes=1&tweet_mode=extended&include_ext_alt_text=true&include_reply_count=true", userName)
}

func (tm *TwitterMonitor) request(url string, req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func parseInt(value interface{}) int {
	switch v := value.(type) {
	case string:
		parsedInt, err := strconv.Atoi(v)
		if err != nil {
			return 0 // Default value or error handling
		}
		return parsedInt
	case int:
		return v
	default:
		return 0 // Default value or error handling
	}
}

func getTweetID(tweet map[string]interface{}) (int64, error) {
	tweetData, exists := tweet["tweetData"].(map[string]interface{})
	if !exists {
		return 0, fmt.Errorf("Error: tweetData not found in tweet")
	}

	tweetDataID, exists := tweetData["id"].(string)
	if !exists {
		return 0, fmt.Errorf("Error: tweetData.id not found in tweet")
	}

	tweetIDInt, err := strconv.ParseInt(tweetDataID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Error converting tweetData.id to int64: %s", err)
	}

	return tweetIDInt, nil
}

func (tm *TwitterMonitor) monitorUser(user User) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["author"] = nil
	result["tweet"] = nil
	tweetInfo, err := tm.getRecent(user)

	if err != nil {
		return result, err
	}
	if tweetInfo["Tweet"] == nil {
		return result, nil
	}

	tweet := make(map[string]interface{})
	tweet["creationTime"] = tweetInfo["Tweet"].(map[string]interface{})["created_at"].(string)
	tweet["author"] = map[string]interface{}{
		"name":       tweetInfo["Tweet"].(map[string]interface{})["user"].(map[string]interface{})["screen_name"].(string),
		"profileURL": "https://twitter.com/" + tweetInfo["Tweet"].(map[string]interface{})["user"].(map[string]interface{})["screen_name"].(string),
		"avatarURL":  tweetInfo["Tweet"].(map[string]interface{})["user"].(map[string]interface{})["profile_image_url"].(string),
	}

	tweet["tweetData"] = map[string]interface{}{
		"id":    tweetInfo["Tweet"].(map[string]interface{})["id_str"].(string),
		"link":  "https://twitter.com/" + tweetInfo["Tweet"].(map[string]interface{})["user"].(map[string]interface{})["screen_name"].(string) + "/status/" + tweetInfo["Tweet"].(map[string]interface{})["id_str"].(string),
		"likes": tweetInfo["Tweet"].(map[string]interface{})["favorite_count"].(float64),
	}

	tweet["content"] = map[string]interface{}{
		"text":      tweetInfo["Tweet"].(map[string]interface{})["full_text"].(string),
		"imageURLs": nil,
		"video":     nil,
		"gif":       nil,
		"tags":      tweetInfo["Tweet"].(map[string]interface{})["entities"].(map[string]interface{})["user_mentions"],
	}
	content := make(map[string]interface{})
	content["text"] = tweetInfo["full_text"]
	content["imageURLs"] = nil
	content["video"] = nil
	content["gif"] = nil
	entities, entitiesExists := tweetInfo["entities"].(map[string]interface{})
	if entitiesExists {
		userMentions, userMentionsExists := entities["user_mentions"]
		if userMentionsExists {
			content["tags"] = userMentions
		}
	}
	tweet["content"] = content
	if tweetInfo["extended_entities"] != nil {
		tweetMedia := tweetInfo["extended_entities"].(map[string]interface{})["media"].([]interface{})
		for _, media := range tweetMedia {
			mediaData := media.(map[string]interface{})
			if mediaData["type"].(string) == "video" {
				for _, variant := range mediaData["video_info"].(map[string]interface{})["variants"].([]interface{}) {
					variantData := variant.(map[string]interface{})
					if variantData["content_type"].(string) == "video/mp4" {
						content["video"] = variantData["url"]
					}
				}
			} else if mediaData["type"].(string) == "image" || mediaData["type"].(string) == "photo" {
				content["imageURLs"] = append(content["imageURLs"].([]string), mediaData["media_url_https"].(string))
			} else if mediaData["type"].(string) == "animated_gif" {
				content["gif"] = mediaData["video_info"].(map[string]interface{})["variants"].([]interface{})[len(mediaData["video_info"].(map[string]interface{})["variants"].([]interface{}))-1].(map[string]interface{})["url"]
			}
		}
	}

	fileContent, err := ioutil.ReadFile(tm.RecentPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return result, nil

	}
	data := make(map[string]interface{})
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return result, nil
	}
	recentID, exists := data[user.Name].(map[string]interface{})
	if !exists {
		addRecent(user.Name, tweet)
	} else if exists {
		tweetID, err := getTweetID(recentID)
		if err != nil {
			fmt.Println(err)
			return result, nil
		}
		tweetIDL, err := getTweetID(tweet)
		if err != nil {
			fmt.Println(err)
			return result, nil
		}
		if tweetID < tweetIDL {
			addRecent(user.Name, tweet)
			result["author"] = user.Name
			result["tweet"] = tweet

			return result, nil
		} else {
			// We have to remove this
			// result["author"] = user.Name
			// result["tweet"] = tweet
			return result, nil
		}
	} else {
		return result, nil
	}
	return result, nil
}
