package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/romanzac/gorilla-feast/controller/dbhandler"
	"github.com/romanzac/gorilla-feast/domain/repository"
	"github.com/romanzac/gorilla-feast/infra/router"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// APIv1 implements APIv1 handlers
type APIv1 struct {
	UserRepo repository.UserRepository
}

// NewAPIv1 creates new API V1
func NewAPIv1(userRepo *dbhandler.DbUserRepo) *APIv1 {
	apiV1 := new(APIv1)
	apiV1.UserRepo = userRepo

	return apiV1
}

// PingPong to test API is alive
func (a *APIv1) PingPong(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Pong!\n")); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// SignupUser implements user registration with password encoding
func (a *APIv1) SignupUser(w http.ResponseWriter, r *http.Request) {

	acct := r.FormValue("acct")
	fullname := r.FormValue("fullname")
	pwd := r.FormValue("pwd")

	// Validate acct(username)
	reAcct := regexp.MustCompile("^([a-z_][a-z0-9_]{3,30})$")
	if !reAcct.MatchString(acct) {
		http.Error(w, "Acct is not valid username", http.StatusBadRequest)
		return
	}

	// Validate fullname
	reFullname := regexp.MustCompile("^([A-Z][a-z]{0,40}\\s{1,10}[A-Z][a-z]{0,49})$")
	if !reFullname.MatchString(fullname) {
		http.Error(w, "Fullname does not follow pattern: \"Jacky Yang\"", http.StatusBadRequest)
		return
	}

	// Validate password for length
	if len(pwd) < 8 {
		http.Error(w, "Password length is less than 8 characters", http.StatusBadRequest)
		return
	}

	err := a.UserRepo.Create(r.FormValue("acct"), r.FormValue("fullname"), r.FormValue("pwd"))
	if err != nil {
		http.Error(w, "Error adding user to database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode("User created successfully"); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// LoginFailures reports user login failures to websocket client
func (a *APIv1) LoginFailures(w http.ResponseWriter, r *http.Request) {

	// Upgrade connection to websocket
	c, err := router.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade to websocket connection has failed:", err)
		return
	}

	// Defer close connection without waiting
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
		}
	}(c)

	// Wait and send login failures to the client
	for {
		message := <-router.LoginFailuresCh
		err = c.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Write to websocket failed:", err)
			break
		}
	}
}

// Login with JWT generation
func (a *APIv1) Login(w http.ResponseWriter, r *http.Request) {

	acct := r.FormValue("acct")
	pwd := r.FormValue("pwd")

	// Validate acct(username)
	reAcct := regexp.MustCompile("^([a-z_][a-z0-9_]{3,30})$")
	if !reAcct.MatchString(acct) {
		http.Error(w, "Acct is not valid username", http.StatusBadRequest)
		return
	}

	// Validate password for length
	if len(pwd) < 8 {
		http.Error(w, "Password length is less than 8 characters", http.StatusBadRequest)
		return
	}

	token, err := a.UserRepo.Validate(acct, pwd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		go func() {
			router.LoginFailuresCh <- err.Error()
		}()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&token); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// Helper function to check and parse query string for sorting field and order direction
func (a *APIv1) validateSortQuery(sortBy string) (string, error) {

	splits := strings.Split(sortBy, ".")
	if len(splits) != 2 {
		return "", errors.New("unknown sortBy parameter, should be field.orderdirection")
	}

	field, order := splits[0], splits[1]
	if order != "desc" && order != "asc" {
		return "", errors.New("unknown ordering in sortBy parameter, should be asc or desc")
	}

	return fmt.Sprintf("%s %s", field, strings.ToUpper(order)), nil
}

// ListAllUsers sends all users with option to have result sorted and paginated
func (a *APIv1) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	var (
		limit     = -1         // default limit is unlimited
		offset    = -1         // default offset is no offset
		sortQuery = "acct ASC" // default ordering is ascending with acct field
		err       error
	)

	sortByQuery := r.URL.Query().Get("sortBy")
	if sortByQuery != "" {
		sortQuery, err = a.validateSortQuery(sortByQuery)
		if err != nil {
			http.Error(w, "SortBy parameter is invalid: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	limitQuery := r.URL.Query().Get("limit")
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil || limit < -1 {
			http.Error(w, "limit parameter is invalid number", http.StatusBadRequest)
			return
		}
	}

	offsetQuery := r.URL.Query().Get("offset")
	if offsetQuery != "" {
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil || offset < -1 {
			http.Error(w, "offset parameter is invalid number", http.StatusBadRequest)
			return
		}
	}

	users, err := a.UserRepo.Find("", "", sortQuery, limit, offset, true)
	if err != nil {
		http.Error(w, "Error to find all users"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// SearchUserbyFullname sends all users with fullname matching the query parameter in URL
func (a *APIv1) SearchUserbyFullname(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	fullname, ok := urlParams["fullname"]

	if !ok {
		http.Error(w, "Unknown error", http.StatusUnprocessableEntity)
		return
	}

	// Validate fullname
	reFullname := regexp.MustCompile("^([A-Z][a-z]{0,40}\\s{1,10}[A-Z][a-z]{0,49})$")
	if !reFullname.MatchString(fullname) {
		http.Error(w, "Fullname does not follow pattern: \"Jacky Yang\"", http.StatusBadRequest)
		return
	}

	users, err := a.UserRepo.Find("", fullname, "", 0, 0, true)
	if err != nil {
		http.Error(w, "DB query error to find all users with fullname: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// GetUserDetail sends for acct all database fields except password
func (a *APIv1) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	acct, ok := urlParams["acct"]

	if !ok {
		http.Error(w, "Unknown error", http.StatusUnprocessableEntity)
		return
	}

	// Validate acct(username)
	reAcct := regexp.MustCompile("^([a-z_][a-z0-9_]{3,30})$")
	if !reAcct.MatchString(acct) {
		http.Error(w, "Acct is not valid username", http.StatusBadRequest)
		return
	}

	users, err := a.UserRepo.Find(acct, "", "", 0, 0, false)
	if err != nil {
		http.Error(w, "DB query error to find user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// UpdateUser updates user with new password or fullname
func (a *APIv1) UpdateUser(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	acct, ok := urlParams["acct"]
	fullname := r.FormValue("fullname")
	pwd := r.FormValue("pwd")

	if !ok {
		http.Error(w, "Unknown error", http.StatusUnprocessableEntity)
		return
	}

	// Validate acct(username)
	reAcct := regexp.MustCompile("^([a-z_][a-z0-9_]{3,30})$")
	if !reAcct.MatchString(acct) {
		http.Error(w, "Acct is not valid username", http.StatusBadRequest)
		return
	}

	// Check if anything to change
	if fullname == "" && pwd == "" {
		http.Error(w, "Nothing to change", http.StatusBadRequest)
		return
	}

	// Validate fullname
	reFullname := regexp.MustCompile("^([A-Z][a-z]{0,40}\\s{1,10}[A-Z][a-z]{0,49})$")
	if fullname != "" && !reFullname.MatchString(fullname) {
		http.Error(w, "Fullname does not follow pattern: \"Jacky Yang\"", http.StatusBadRequest)
		return
	}

	// Validate password for length
	if pwd != "" && len(pwd) < 8 {
		http.Error(w, "Password length is less than 8 characters", http.StatusBadRequest)
		return
	}

	if err := a.UserRepo.Update(acct, fullname, pwd); err != nil {
		http.Error(w, "Error updating user in database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("User updated successfully"); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// DeleteUser removes user from database
func (a *APIv1) DeleteUser(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	claimedAcct := r.Header.Get("acct")
	acct, ok := urlParams["acct"]

	// Compare user performing delete with the user to be deleted
	if claimedAcct == acct {
		http.Error(w, "User cannot delete herself", http.StatusUnprocessableEntity)
		return
	}

	if !ok {
		http.Error(w, "Unknown error", http.StatusUnprocessableEntity)
		return
	}

	// Validate acct(username)
	reAcct := regexp.MustCompile("^([a-z_][a-z0-9_]{3,30})$")
	if !reAcct.MatchString(acct) {
		http.Error(w, "Acct is not valid username", http.StatusBadRequest)
		return
	}

	if err := a.UserRepo.Delete(acct); err != nil {
		http.Error(w, "Error deleting user from database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("User deleted successfully"); err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}
