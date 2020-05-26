package v1

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/logrusorgru/aurora"

	"github.com/viswas163/MarvelousShipt/models"
)

var (
	// PROD_ENV : Switch for prod/dev env
	PROD_ENV = false
	// BaseURL : Marvel API base URL
	BaseURL = "https://gateway.marvel.com/v1/public/"

	// MarvelPublicAPIKeyEnvKey : The OS environment key to fetch the Marvel Developer Public API Key string
	MarvelPublicAPIKeyEnvKey = "MARVEL_PUBLIC_API_KEY"
	// MarvelPrivateAPIKeyEnvKey : The OS environment key to fetch the Marvel Developer Private API Key string
	MarvelPrivateAPIKeyEnvKey = "MARVEL_PRIVATE_API_KEY"
	// MarvelPublicAPIKey : The Marvel public API key
	MarvelPublicAPIKey = ""
	// MarvelPrivateAPIKey : The Marvel private API key
	MarvelPrivateAPIKey = ""
	// AuthClient : The global Authentication Client instance
	AuthClient models.AuthClient

	privateKeyFileName = "MarvelPrivateKey.txt"
	publicKeyFileName  = "MarvelPublicKey.txt"
)

// InitAuthClient : Initializes the Authentication client
func InitAuthClient() models.AuthClient {
	hasKey := ""

	// Get the public key from parent directory
	path, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working dir path : ", err)
	}
	pubPath := filepath.Join(filepath.Dir(path), publicKeyFileName)
	content, err := ioutil.ReadFile(pubPath)
	if err != nil {
		log.Println("Error reading public key file", err)
	}
	MarvelPublicAPIKey = string(content)

	// Get user input for private key
	for PROD_ENV && !strings.EqualFold(hasKey, "y") && !strings.EqualFold(hasKey, "n") && !strings.EqualFold(hasKey, "yes") && !strings.EqualFold(hasKey, "no") {
		fmt.Print("\nDo you have a Marvel Developer Private API Key? (y/n) : ")
		fmt.Scanln(&hasKey)
	}

	// Get private key from user
	if PROD_ENV && (strings.EqualFold(hasKey, "y") || strings.EqualFold(hasKey, "yes")) {
		userAPIKey := ""
		fmt.Print(aurora.Cyan("\nNote : The Private API Key is not stored anywhere"), "\nEnter the Private API Key please : ")
		fmt.Scanln(&userAPIKey)
		os.Setenv(MarvelPrivateAPIKeyEnvKey, userAPIKey)
	} else { // Get private key from file in parent directory
		privPath := filepath.Join(filepath.Dir(path), privateKeyFileName)
		content, err := ioutil.ReadFile(privPath)
		if err != nil {
			log.Println("Error reading private key file", err)
		}
		os.Setenv(MarvelPrivateAPIKeyEnvKey, string(content))
	}
	MarvelPrivateAPIKey = os.Getenv(MarvelPrivateAPIKeyEnvKey)

	client := models.AuthClient{
		PublicKey:  MarvelPublicAPIKey,
		PrivateKey: MarvelPrivateAPIKey,
	}
	AuthClient = client
	return client
}

// GetAuthenticator : Gets the authenticator using client params
func getAuthenticator() *models.Authenticator {

	ts := strconv.FormatInt(time.Now().Unix(), 10)

	authString := ts + AuthClient.PrivateKey + AuthClient.PublicKey

	hasher := md5.New()
	hasher.Write([]byte(authString))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return &models.Authenticator{
		Timestamp: ts,
		PublicKey: AuthClient.PublicKey,
		Hash:      hash,
	}
}

// RunAPI : Runs Authentication for request
func RunAPI(cassette string) ([]byte, error) {
	req, err := GetAuthRequest(cassette)
	if err != nil {
		return nil, err
	}
	response := executeRequest(req)
	checkResponseCode(http.StatusOK, response)
	body, err := checkResponseBody(response)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetAuthRequest : Returns the authentication request URL for the provided Resource
func GetAuthRequest(cassette string) (*http.Request, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("no caller information when determining file location")
		return &http.Request{}, errors.New("Cannot get runtime caller file name")
	}

	path := filepath.Join(path.Dir(filename), "fixtures")
	rec, err := recorder.New(path)
	if err != nil {
		fmt.Println("could not create recoorder with path : ", path)
		return &http.Request{}, err
	}
	// defer rec.Stop()

	// Get Authenticator using client params
	auth := getAuthenticator()

	// Create Base sling
	recHTTPClient := &http.Client{
		Transport: rec,
	}
	base := sling.New().Client(recHTTPClient).Base(BaseURL).Path(cassette)
	// Construct URL query for base sling
	base.QueryStruct(auth)

	// Create http request from base sling
	req, err := base.Request()
	if err != nil {
		fmt.Println("Error parsing request : ", err)
		return &http.Request{}, err
	}
	return req, nil
}

func executeRequest(req *http.Request) *http.Response {

	// Using http client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Print("Error performing request : ", err)
	}
	return response
}

func checkResponseCode(expected int, response *http.Response) {
	actual := response.StatusCode
	if expected != actual {
		fmt.Printf("Expected response code %d. Got %d\n", expected, actual)
		return
	}
}

func checkResponseBody(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body : ", err)
		return nil, err
	}
	return body, nil
}
