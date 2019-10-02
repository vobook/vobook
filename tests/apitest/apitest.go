package apitest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/vovainside/vobook/cmd/server/routes"
	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/logger"
	"github.com/vovainside/vobook/tests/assert"
	"github.com/vovainside/vobook/utils"
)

var (
	DB        orm.DB
	AuthUser  *models.User
	authToken string
	router    *gin.Engine
)

type Headers map[string]string

type Request struct {
	method       string
	Path         string
	Body         interface{}
	bodyReader   io.Reader
	Headers      Headers
	BindResponse interface{}
	AssertStatus int
	IsPublic     bool
}

func init() {
	// changing working dir to vobook/bin
	// just want to use main .config.yaml, thats why
	_ = os.Chdir("../../../bin")
	conf := config.Get()

	// override config
	conf.LogsFilePath = "tests.log"
	conf.Mail.Driver = "test"
	conf.App.Env = config.TestEnv

	logger.Setup()

	db := database.Conn()
	_, _ = db.Exec("SET AUTOCOMMIT TO OFF")
	tx, err := db.Conn().Begin()
	if err != nil {
		panic(err)
	}
	database.SetDB(tx)
	DB = database.ORM()

	router = gin.New()
	routes.Register(router)
}

// makeRequest makes request (Duh...)
func makeRequest(t *testing.T, r Request) *httptest.ResponseRecorder {
	req, err := http.NewRequest(
		r.method,
		r.Path,
		r.bodyReader,
	)
	assert.NotError(t, err)
	t.Log(fmt.Sprintf("[%d] %s %s", r.AssertStatus, r.method, r.Path))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if !r.IsPublic {
		if authToken == "" {
			//Login(t)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// TestRequest sends given request,
// asserts that status is correct
// binds response and returns response
func TestRequest(t *testing.T, req Request) *httptest.ResponseRecorder {
	var body []byte
	var err error

	if len(req.Headers) == 0 {
		req.Headers = Headers{}
	}

	if req.Body != nil {
		body, err = json.Marshal(req.Body)
		assert.NotError(t, err)
	}

	req.bodyReader = strings.NewReader(string(body))
	resp := makeRequest(t, req)
	if resp.Code != req.AssertStatus {
		t.Errorf(`%s "%s" responds with %d, expecting %d`, req.method, req.Path, resp.Code, req.AssertStatus)
		t.Errorf("Response:\n%s", resp.Body.String())
		if req.Body != nil {
			t.Errorf("Request:\n%+v", req.bodyReader)
		}
		t.FailNow()
	}

	err = json.Unmarshal(resp.Body.Bytes(), &req.BindResponse)
	if err != nil {
		t.Log(resp.Body.String())
		t.Fatal(err)
	}

	return resp
}

//func Login(t *testing.T) {
//	user, err := users.Fake()
//	assert.NotError(t, err)
//	LoginAs(t, user)
//}

//func User(t *testing.T) *users.User {
//	if AuthUser == nil {
//		Login(t)
//	}
//
//	return AuthUser
//}

//func LoginAs(t *testing.T, user *users.User) {
//	shouldUpdate := false
//	if user.Token == "" {
//		token, err := utils.UniqueToken(users.Table)
//		assert.NotError(t, err)
//		user.Token = token
//		shouldUpdate = true
//	}
//	if user.TokenExpiresAt.Before(time.Now()) {
//		user.TokenExpiresAt = time.Now().Add(time.Hour)
//		shouldUpdate = true
//	}
//
//	if shouldUpdate {
//		err := user.Update()
//		assert.NotError(t, err)
//	}
//
//	AuthUser = user
//	authToken = users.EncodeToken(user)
//}

func Logout() {
	AuthUser = nil
	authToken = ""
}

// POST is a wrapper for TestRequest
func POST(t *testing.T, req Request) *httptest.ResponseRecorder {
	req.method = "POST"
	return TestRequest(t, req)
}

// PUT is a wrapper for TestRequest
func PUT(t *testing.T, req Request) *httptest.ResponseRecorder {
	req.method = "PUT"
	return TestRequest(t, req)
}

// PATCH is a wrapper for TestRequest
func PATCH(t *testing.T, req Request) *httptest.ResponseRecorder {
	req.method = "PATCH"
	return TestRequest(t, req)
}

// DELETE is a wrapper for TestRequest
func DELETE(t *testing.T, req Request) *httptest.ResponseRecorder {
	req.method = "DELETE"
	return TestRequest(t, req)
}

// GET is a wrapper for TestRequest
func GET(t *testing.T, req Request) *httptest.ResponseRecorder {
	req.method = "GET"
	return TestRequest(t, req)
}

// Create makes "create" request
func Create(t *testing.T, path string, body, response interface{}) *httptest.ResponseRecorder {
	req := Request{
		Path:         path,
		Body:         body,
		BindResponse: response,
		AssertStatus: http.StatusCreated,
	}

	return POST(t, req)
}

// Update makes "update" request
func Update(t *testing.T, path string, body, response interface{}) *httptest.ResponseRecorder {
	req := Request{
		Path:         path,
		Body:         body,
		BindResponse: response,
		AssertStatus: http.StatusOK,
	}

	return PUT(t, req)
}

// Patch makes "patch" request
func Patch(t *testing.T, path string, body, response interface{}) *httptest.ResponseRecorder {
	req := Request{
		Path:         path,
		Body:         body,
		BindResponse: response,
		AssertStatus: http.StatusOK,
	}

	return PATCH(t, req)
}

// Delete makes "delete" request
func Delete(t *testing.T, path string, response interface{}) *httptest.ResponseRecorder {
	req := Request{
		Path:         path,
		BindResponse: response,
		AssertStatus: http.StatusOK,
	}

	return DELETE(t, req)
}

// View makes "get" request for a single entity
func View(t *testing.T, path string, response interface{}) *httptest.ResponseRecorder {
	req := Request{
		Path:         path,
		BindResponse: response,
		AssertStatus: http.StatusOK,
	}

	return GET(t, req)
}

// Search makes "get" request with query params
func Search(t *testing.T, path string, request, response interface{}) *httptest.ResponseRecorder {
	query, err := utils.StructToQueryString(request)
	if err != nil {
		t.Fatal(err.Error())
	}

	path = fmt.Sprintf(`%s?%s`, path, query)

	req := Request{
		Path:         path,
		BindResponse: response,
		AssertStatus: http.StatusOK,
	}

	return GET(t, req)
}
