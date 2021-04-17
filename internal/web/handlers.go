package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	bundle "github.com/bennycio/bundle/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

const (
	AwsS3Region = bundle.AwsS3Region
	AwsS3Bucket = bundle.AwsS3Bucket
	ReqFileType = bundle.RequiredFileType
)

func init() {
	tpl = template.Must(template.ParseGlob("../../assets/templates/*.gohtml"))
}

func HandleRoot(w http.ResponseWriter, req *http.Request) {

	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

}
func HandleSignup(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		r.ParseForm()

		authURL := "https://bundle.us.auth0.com/oauth/token"

		secret := os.Getenv("CLIENT_SECRET")

		values := url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {"22oXY4A0h9Rfbo3XEAn8Fbptx715dBe4"},
			"client_secret": {secret},
			"audience":      {"https://bundlemc.io/auth/users"},
		}

		authRes, err := http.PostForm(authURL, values)
		if err != nil {
			log.Fatal(err)
		}

		defer authRes.Body.Close()
		body, _ := ioutil.ReadAll(authRes.Body)

		auth := &bundle.Authorization{}

		err = json.Unmarshal(body, auth)
		if err != nil {
			log.Fatal(err)
		}

		newUser := bundle.User{
			r.FormValue("username"),
			r.FormValue("email"),
			r.FormValue("password"),
		}

		asJSON, _ := json.Marshal(newUser)

		br := bytes.NewReader(asJSON)

		createUser, err := http.NewRequest(http.MethodPost, "http://localhost:8070/users", br)
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
		log.Fatal(err)
	}

}

func BundleHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(AwsS3Region)})

	if req.Method == http.MethodGet {
		fmt.Println("GET")

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		name := req.FormValue("name")
		version := req.FormValue("version")

		author, err := bundle.GetPluginAuthor(name)
		if err != nil {
			panic(err)
		}

		buf := aws.NewWriteAtBuffer([]byte{})

		fn := filepath.Join(author, name, version, name+".jar")

		fmt.Println("FILENAME: " + fn)
		downloader := s3manager.NewDownloader(sess)
		_, err = downloader.Download(buf, &s3.GetObjectInput{
			Bucket: aws.String(AwsS3Bucket),
			Key:    aws.String(fn),
		})
		if err != nil {
			panic(err)
		}

		w.Write(buf.Bytes())
	}

	if req.Method == http.MethodPost {

		version := req.Header.Get("Project-Version")
		userJSON := req.Header.Get("User")
		name := req.Header.Get("Project-Name")

		validatedUser, err := bundle.ValidateAndReturnUser(userJSON)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		dbUser, err := bundle.GetUserFromDatabase(validatedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(
			"mongodb+srv://benny-bundle:thisismypassword1@bundle.mveuj.mongodb.net/users?retryWrites=true&w=majority",
		))
		defer cancel()
		defer func() {
			if err := client.Disconnect(ctx); err != nil {
				bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		// TODO figure out first if it is their plugin, then if it is update the version, if not tell them bye bye

		collection := client.Database("users").Collection("plugins")

		countPlugins, err := collection.CountDocuments(ctx, bson.D{{"$regex", primitive.Regex{"plugin", "i"}}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if countPlugins > 0 {
			err = errors.New("Plugin already exists with given name")
			bundle.WriteResponse(w, err.Error(), http.StatusConflict)
			return
		}

		_, err = collection.InsertOne(ctx, bson.D{{"plugin", name}, {"user", dbUser.Username}, {"version", version}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fp := filepath.Join(dbUser.Username, name, version, name+".jar")

		uploader := s3manager.NewUploader(sess)
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   req.Body,
			Bucket: aws.String("bundle-repository"),
			Key:    aws.String(fp),
		})
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		log.Println("Successfully uploaded to", result.Location)
	}

}

func UserHandlerFunc(w http.ResponseWriter, r *http.Request) {

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

		newUser, err := bundle.ValidateAndReturnUser(string(bb))

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		bcryptPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(
			"mongodb+srv://benny-bundle:thisismypassword1@bundle.mveuj.mongodb.net/users?retryWrites=true&w=majority",
		))
		defer cancel()
		defer func() {
			if err := client.Disconnect(ctx); err != nil {
				bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		collection := client.Database("users").Collection("users")

		countUserName, err := collection.CountDocuments(ctx, bson.D{{"username", newUser.Username}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if countUserName > 0 {
			err = errors.New("user already exists with given username")
			bundle.WriteResponse(w, err.Error(), http.StatusConflict)
			return
		}

		countEmail, err := collection.CountDocuments(ctx, bson.D{{"email", newUser.Email}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if countEmail > 0 {
			err = errors.New("user already exists with given email")
			bundle.WriteResponse(w, err.Error(), http.StatusConflict)
			return
		}

		_, err = collection.InsertOne(ctx, bson.D{{"username", newUser.Username}, {"email", newUser.Email}, {"password", string(bcryptPass)}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bundle.WriteResponse(w, "", http.StatusOK)
	}
}
