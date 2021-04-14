package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
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

const (
	AWS_S3_REGION      = "us-east-1"
	AWS_S3_BUCKET      = "bundle-repository"
	REQUIRED_FILE_TYPE = "application/zip"
)

func BundleHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})

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
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(fn),
		})
		if err != nil {
			panic(err)
		}

		w.Write(buf.Bytes())
	}

	if req.Method == http.MethodPost {

		// fileBytes, err := io.ReadAll(req.Body)

		// if err != nil {
		// 	panic(err)
		// }

		// fileType := http.DetectContentType(fileBytes)

		// if fileType != REQUIRED_FILE_TYPE {
		// 	log.Fatal("File must be of " + REQUIRED_FILE_TYPE + " type")
		// }

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

		collection := client.Database("users").Collection("plugins")

		// Test this out for case sensitivity

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

		_, err = collection.InsertOne(ctx, bson.D{{"plugin", name}, {"user", dbUser.Username}})

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
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
