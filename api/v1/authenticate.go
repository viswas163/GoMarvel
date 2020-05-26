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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/logrusorgru/aurora"

	"github.com/viswas163/MarvelousShipt/models"
)

var (
	// ProdEnv : Switch for prod/dev env
	ProdEnv = true
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
	// Limit : Limits the response query results
	Limit = 100

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
	pubPath := filepath.Join(path, publicKeyFileName)
	content, err := ioutil.ReadFile(pubPath)
	if err != nil {
		log.Println("Error reading public key file", err)
	}
	MarvelPublicAPIKey = string(content)

	// Get user input for private key
	for ProdEnv && !strings.EqualFold(hasKey, "y") && !strings.EqualFold(hasKey, "n") && !strings.EqualFold(hasKey, "yes") && !strings.EqualFold(hasKey, "no") {
		fmt.Print("\nDo you have a Marvel Developer Private API Key? (y/n) : ")
		fmt.Scanln(&hasKey)
	}

	// Get private key from user
	if ProdEnv && (strings.EqualFold(hasKey, "y") || strings.EqualFold(hasKey, "yes")) {
		userAPIKey := ""
		fmt.Print(aurora.Cyan("\nNote : The Private API Key is not stored anywhere"), "\nEnter the Private API Key please : ")
		fmt.Scanln(&userAPIKey)
		os.Setenv(MarvelPrivateAPIKeyEnvKey, userAPIKey)
	} else { // Get private key from file in parent directory
		privPath := filepath.Join(path, privateKeyFileName)
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
func getAuthenticator() (*models.Authenticator, error) {

	ts := strconv.FormatInt(time.Now().Unix(), 10)

	authString := ts + AuthClient.PrivateKey + AuthClient.PublicKey

	hasher := md5.New()
	if _, err := hasher.Write([]byte(authString)); err != nil {
		fmt.Println("getAuthenticator() : Error hashing byte array of authstring : ", authString, "\n Error : ", err)
		return nil, err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	return &models.Authenticator{
		Timestamp: ts,
		PublicKey: AuthClient.PublicKey,
		Hash:      hash,
	}, nil
}

// GetAuthRequest : Returns the authentication request URL for the provided Resource
func GetAuthRequest(resource string) (*sling.Sling, error) {
	var netClient = http.DefaultClient

	// Get Authenticator using client params
	auth, err := getAuthenticator()
	if err != nil {
		return nil, err
	}

	base := sling.New().Client(netClient).Base(BaseURL).Path(resource)
	// Construct URL query for base sling
	base.QueryStruct(auth)

	return base, nil
}

func RunAPIWithParam(resource string, offset int, params []string) ([]byte, error) {
	return RunAPI(resource, offset, params)
}

func RunAPIWithoutParam(resource string, offset int) ([]byte, error) {
	return RunAPI(resource, offset, []string{})
}

// RunAPI : Runs Authentication for request
func RunAPI(resource string, offset int, param []string) ([]byte, error) {
	req, err := GetAuthRequest(resource)
	if err != nil {
		return nil, err
	}

	params := &models.CharacterParams{Limit: Limit, Offset: offset}
	response, err := executeRequest(req, resource, param, params)
	if err != nil {
		return nil, err
	}

	checkResponseCode(http.StatusOK, response)
	body, err := checkResponseBody(response)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func executeRequest(req *sling.Sling, pathURL string, param []string, params interface{}) (*http.Response, error) {

	req = req.New().Get(pathURL)
	if len(param) > 0 {
		str := pathURL
		for _, p := range param {
			str += "/" + p
		}
		// fmt.Println(str)
		req.Path(str)
	}
	req = req.QueryStruct(params)

	r, _ := req.Request()
	// response := &http.Response{}
	// fmt.Println(r)

	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		fmt.Println("executeRequest() : Error in client response : ", err)
		return nil, err
	}
	return response, nil
}

func checkResponseCode(expected int, response *http.Response) {
	actual := response.StatusCode
	if expected != actual {
		fmt.Printf("Expected response code %d. Got %d\n", expected, actual)
		return
	}
}

func checkResponseBody(response *http.Response) ([]byte, error) {
	if response == nil || response.Body == nil {
		return nil, errors.New("Empty response")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body : ", err)
		return nil, err
	}

	return body, nil
}
