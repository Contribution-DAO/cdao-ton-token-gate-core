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

func (h *ServiceHandler) LinkTwitter(address string, twitterUsername string, twitterAvatar string, twitterAccessToken string, twitterAccessTokenSecret string) (*model.Wallet, error) {
	wallet := model.Wallet{
		ID: address,
	}

	if result := h.db.Model(&wallet).Updates(map[string]interface{}{
		"twitter_username":            twitterUsername,
		"twitter_avatar":              twitterAvatar,
		"twitter_access_token":        twitterAccessToken,
		"twitter_access_token_secret": twitterAccessTokenSecret,
	}); result.Error != nil {
		return nil, result.Error
	}

	return &wallet, nil
}

// generates a random nonce
func generateNonce() string {
	rand.Seed(time.Now().UnixNano())
	nonce := rand.Int63()
	return strconv.FormatInt(nonce, 10)
}

// generates a timestamp
func generateTimestamp() string {
	timestamp := time.Now().Unix()
	return strconv.FormatInt(timestamp, 10)
}

// generates the signature base string
func generateBaseString(consumerKey string, nonce string, accessToken string, timestamp string) string {
	// set up the parameters for the signature base string
	parameters := make(url.Values)
	parameters.Set("oauth_consumer_key", consumerKey)
	parameters.Set("oauth_nonce", nonce)
	parameters.Set("oauth_signature_method", "HMAC-SHA1")
	parameters.Set("oauth_timestamp", timestamp)
	parameters.Set("oauth_token", accessToken)
	parameters.Set("oauth_version", "1.0")

	// encode the parameters and sort them
	encodedParams := parameters.Encode()
	sortedParams := strings.ReplaceAll(encodedParams, "+", "%20")

	// construct the signature base string
	baseString := http.MethodGet + "&" + url.QueryEscape("<request-url>") + "&" + url.QueryEscape(sortedParams)

	return baseString
}

// generates the signature
func generateSignature(consumerSecret string, accessSecret string, baseString string) string {
	// construct the signing key
	signingKey := url.QueryEscape(consumerSecret) + "&" + url.QueryEscape(accessSecret)

	// calculate the HMAC-SHA1 signature
	hmac := hmac.New(sha1.New, []byte(signingKey))
	hmac.Write([]byte(baseString))
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))

	return signature
}

// generates the authorization header
func generateAuthorizationHeader(consumerKey string, accessToken string, nonce string, signature string, timestamp string) string {
	// set up the authorization header parameters
	parameters := make(url.Values)
	parameters.Set("oauth_consumer_key", consumerKey)
	parameters.Set("oauth_nonce", nonce)
	parameters.Set("oauth_signature", signature)
	parameters.Set("oauth_signature_method", "HMAC-SHA1")
	parameters.Set("oauth_timestamp", timestamp)
	parameters.Set("oauth_token", accessToken)
	parameters.Set("oauth_version", "1.0")

	// encode the parameters and sort them
	encodedParams := parameters.Encode()
	sortedParams := strings.ReplaceAll(encodedParams, "+", "%20")

	// construct the authorization header
	authorization := "OAuth " + sortedParams

	return authorization
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

	nonce, err := GenerateRandomString(16)
	if err != nil {
		return false, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Set the request headers
	req.Header.Set("Authorization", generateAuthorizationHeader(
		os.Getenv("TWITTER_API_KEY"),
		accessToken,
		nonce,
		generateSignature(os.Getenv("TWITTER_API_SECRET"), accessTokenSecret, generateBaseString(
			os.Getenv("TWITTER_API_KEY"),
			nonce,
			accessToken,
			timestamp,
		)),
		timestamp,
	))

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
