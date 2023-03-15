package services

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Contribution-DAO/cdao-ton-token-gate-core/model"
)

func (h *ServiceHandler) LinkTwitter(address string, twitterUserId string, twitterUsername string, twitterName string, twitterAvatar string, twitterAccessToken string, twitterAccessTokenSecret string) (*model.Wallet, error) {
	wallet := model.Wallet{
		ID: address,
	}

	if result := h.db.Model(&wallet).Updates(map[string]interface{}{
		"twitter_user_id":             twitterUserId,
		"twitter_username":            twitterUsername,
		"twitter_name":                twitterName,
		"twitter_avatar":              twitterAvatar,
		"twitter_access_token":        twitterAccessToken,
		"twitter_access_token_secret": twitterAccessTokenSecret,
	}); result.Error != nil {
		return nil, result.Error
	}

	return &wallet, nil
}

type OAuth1 struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

// Params being any key-value url query parameter pairs
func (auth OAuth1) BuildOAuth1Header(method, path string, params map[string]string) string {
	vals := url.Values{}
	vals.Add("oauth_nonce", generateNonce())
	vals.Add("oauth_consumer_key", auth.ConsumerKey)
	vals.Add("oauth_signature_method", "HMAC-SHA1")
	vals.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	vals.Add("oauth_token", auth.AccessToken)
	vals.Add("oauth_version", "1.0")

	for k, v := range params {
		vals.Add(k, v)
	}
	// net/url package QueryEscape escapes " " into "+", this replaces it with the percentage encoding of " "
	parameterString := strings.Replace(vals.Encode(), "+", "%20", -1)

	// Calculating Signature Base String and Signing Key
	signatureBase := strings.ToUpper(method) + "&" + url.QueryEscape(strings.Split(path, "?")[0]) + "&" + url.QueryEscape(parameterString)
	signingKey := url.QueryEscape(auth.ConsumerSecret) + "&" + url.QueryEscape(auth.AccessSecret)
	signature := calculateSignature(signatureBase, signingKey)

	return "OAuth oauth_consumer_key=\"" + url.QueryEscape(vals.Get("oauth_consumer_key")) + "\", oauth_nonce=\"" + url.QueryEscape(vals.Get("oauth_nonce")) +
		"\", oauth_signature=\"" + url.QueryEscape(signature) + "\", oauth_signature_method=\"" + url.QueryEscape(vals.Get("oauth_signature_method")) +
		"\", oauth_timestamp=\"" + url.QueryEscape(vals.Get("oauth_timestamp")) + "\", oauth_token=\"" + url.QueryEscape(vals.Get("oauth_token")) +
		"\", oauth_version=\"" + url.QueryEscape(vals.Get("oauth_version")) + "\""
}

func calculateSignature(base, key string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(base))
	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

func generateNonce() string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 48)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}
	return string(b)
}

func (h *ServiceHandler) VerifyTwitterFollow(targetUsername string, accessToken string, accessTokenSecret string) (bool, error) {
	// Set up the request URL and query parameters
	path := "https://api.twitter.com/1.1/friendships/lookup.json"
	params := url.Values{
		"screen_name": {targetUsername},
	}
	query := "?" + params.Encode()

	// Create a new request with the URL and query parameters
	req, err := http.NewRequest("GET", path+query, nil)
	if err != nil {
		return false, err
	}

	auth := OAuth1{
		ConsumerKey:    os.Getenv("TWITTER_API_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_API_SECRET"),
		AccessToken:    accessToken,
		AccessSecret:   accessTokenSecret,
	}

	authorization := auth.BuildOAuth1Header(http.MethodGet, path, map[string]string{
		"screen_name": targetUsername,
	})

	// Set the request headers
	req.Header.Set("Authorization", authorization)

	// fmt.Println(authorization)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// fmt.Println(string(body))

	// Unmarshal the JSON response into a struct
	var data []struct {
		Connections []string `json:"connections"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return false, nil
	}

	for _, connection := range data[0].Connections {
		if connection == "following" {
			return true, nil
		}
	}

	return false, nil

	// Print the response data
	// fmt.Println("Field1:", data.Field1)
	// fmt.Println("Field2:", data.Field2)
}
