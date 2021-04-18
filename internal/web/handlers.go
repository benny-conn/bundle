package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	bundle "github.com/bennycio/bundle/internal"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

const ReqFileType = bundle.RequiredFileType

var (
	AwsS3Region = bundle.AwsS3Region
	AwsS3Bucket = bundle.AwsS3Bucket
	MongoURL    = bundle.MongoURL
)

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}

func RootHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

}
func SignupHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == http.MethodPost {

		auth, err := bundle.GetAuthToken()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusForbidden)
			return
		}

		newUser := bundle.User{
			r.FormValue("username"),
			r.FormValue("email"),
			r.FormValue("password"),
		}

		asJSON, _ := json.Marshal(newUser)

		br := bytes.NewReader(asJSON)

		createUser, err := http.NewRequest(http.MethodPost, "http://localhost:8080/users", br)
		if err != nil {
			log.Fatal(err)
		}

		createUser.Header.Set("authorization", "Bearer "+auth.Token)

		createRes, err := http.DefaultClient.Do(createUser)
		if err != nil {
			log.Fatal(err)
		}

		defer createRes.Body.Close()
		createResBody, _ := ioutil.ReadAll(createRes.Body)
		fmt.Println(string(createResBody))

	}

	err := tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	if err != nil {
		bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func BundleHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(AwsS3Region)})

	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		version := r.FormValue("version")

		plugin, err := bundle.GetPluginByName(name)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf := aws.NewWriteAtBuffer([]byte{})

		fn := filepath.Join(plugin.User, plugin.Plugin, version, plugin.Plugin+".jar")

		downloader := s3manager.NewDownloader(sess)
		_, err = downloader.Download(buf, &s3.GetObjectInput{
			Bucket: aws.String(AwsS3Bucket),
			Key:    aws.String(fn),
		})
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		bundle.WriteResponse(w, string(buf.Bytes()), http.StatusOK)
	}

	if r.Method == http.MethodPost {

		version := r.Header.Get("Plugin-Version")
		userJSON := r.Header.Get("User")
		pluginName := r.Header.Get("Plugin-Name")

		validatedUser, err := bundle.ValidateAndReturnUser(userJSON)

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := bundle.GetUserFromDatabase(validatedUser)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		session, err := bundle.GetMongoSession()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer session.Cancel()

		if err != nil {
			panic(err)
		}

		collection := session.Client.Database("users").Collection("plugins")
		decodedPluginResult := &bundle.Plugin{}

		err = collection.FindOne(session.Ctx, bson.D{{"plugin", bundle.NewCaseInsensitiveRegex(pluginName)}}).Decode(decodedPluginResult)

		if err == nil {
			isUserPluginAuthor := strings.EqualFold(decodedPluginResult.User, validatedUser.Username)

			if isUserPluginAuthor {
				_, err = collection.UpdateOne(session.Ctx, bson.D{{"plugin", bundle.NewCaseInsensitiveRegex(pluginName)}}, bson.D{{"plugin", decodedPluginResult}, {"user", dbUser.Username}, {"version", version}})
				if err != nil {
					bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				err = errors.New("you are not permitted to edit this plugin")
				bundle.WriteResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
		} else {
			_, err = collection.InsertOne(session.Ctx, bson.D{{"plugin", pluginName}, {"user", dbUser.Username}, {"version", version}})
		}

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fp := filepath.Join(dbUser.Username, decodedPluginResult.Plugin, version, decodedPluginResult.Plugin+".jar")

		uploader := s3manager.NewUploader(sess)
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   r.Body,
			Bucket: aws.String("bundle-repository"),
			Key:    aws.String(fp),
		})
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		fmt.Println("Successfully uploaded to", result.Location)
	}

}

func UserHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == http.MethodPost {

		authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
		token := authHeaderParts[1]

		hasScope := bundle.CheckScope("create:bundleuser", token)

		if !hasScope {
			bundle.WriteResponse(w, "Insufficient Permission", http.StatusForbidden)
			return
		}

		bb, err := io.ReadAll(r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		validatedUser, err := bundle.ValidateAndReturnUser(string(bb))

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		bcryptPass, err := bcrypt.GenerateFromPassword([]byte(validatedUser.Password), bcrypt.DefaultCost)

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := bundle.GetMongoSession()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer session.Cancel()

		collection := session.Client.Database("users").Collection("users")

		countUserName, err := collection.CountDocuments(session.Ctx, bson.D{{"username", bundle.NewCaseInsensitiveRegex(validatedUser.Username)}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if countUserName > 0 {
			err = errors.New("user already exists with given username")
			bundle.WriteResponse(w, err.Error(), http.StatusConflict)
			return
		}

		countEmail, err := collection.CountDocuments(session.Ctx, bson.D{{"email", bundle.NewCaseInsensitiveRegex(validatedUser.Email)}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if countEmail > 0 {
			err = errors.New("user already exists with given email")
			bundle.WriteResponse(w, err.Error(), http.StatusConflict)
			return
		}

		_, err = collection.InsertOne(session.Ctx, bson.D{{"username", validatedUser.Username}, {"email", validatedUser.Email}, {"password", string(bcryptPass)}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bundle.WriteResponse(w, "", http.StatusOK)
	}
}

func PluginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("plugin")

		plugin, err := bundle.GetPluginByName(pluginName)
		if err != nil {
			panic(err)
		}

		asJSON, err := json.Marshal(plugin)
		if err != nil {
			panic(err)
		}

		bundle.WriteResponse(w, string(asJSON), http.StatusOK)

	}
}
