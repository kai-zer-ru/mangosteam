package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vincentserpoul/mangosteam"
)

func TestDoLogin(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getMockKOLoginDologin())
	}))
	defer ts.Close()

	client := http.Client{}
	username := "mangosteam"
	encryptedPassword := "123"
	rsatimestamp := "123"
	emailauthKeyedIn := ""
	captchaGID := ""
	captchaKeyedIn := ""
	twoFactorCode := ""
	mangosteam.BaseSteamWebURL = ts.URL

	_, err := DoLogin(
		&client,
		username,
		encryptedPassword,
		rsatimestamp,
		emailauthKeyedIn,
		captchaGID,
		captchaKeyedIn,
		twoFactorCode,
	)

	if err == nil {
		t.Errorf("Dologin returns no error when login failed")
	}

	return
}

func TestDoOKLogin(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{Name: "sessionid", Value: "1234f"})
		http.SetCookie(w, &http.Cookie{Name: "steamLogin", Value: "123"})
		http.SetCookie(w, &http.Cookie{Name: "steamLoginSecure", Value: "1234"})
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetMockOKLoginDologin())
	}))
	defer ts.Close()

	client := http.Client{}
	username := "mangosteam"
	encryptedPassword := "123"
	rsatimestamp := "123"
	emailauthKeyedIn := ""
	captchaGID := ""
	captchaKeyedIn := ""
	twoFactorCode := ""
	mangosteam.BaseSteamWebURL = ts.URL

	_, err := DoLogin(
		&client,
		username,
		encryptedPassword,
		rsatimestamp,
		emailauthKeyedIn,
		captchaGID,
		captchaKeyedIn,
		twoFactorCode,
	)

	if err != nil {
		t.Errorf("Dologin returns error %v when login successful", err)
	}

	return
}

func TestHttpNotOKLogin(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := http.Client{}
	username := "mangosteam"
	encryptedPassword := "123"
	rsatimestamp := "123"
	emailauthKeyedIn := ""
	captchaGID := ""
	captchaKeyedIn := ""
	twoFactorCode := ""
	mangosteam.BaseSteamWebURL = ts.URL

	_, err := DoLogin(
		&client,
		username,
		encryptedPassword,
		rsatimestamp,
		emailauthKeyedIn,
		captchaGID,
		captchaKeyedIn,
		twoFactorCode,
	)

	if err == nil {
		t.Errorf("Dologin returns no error when http status is not statusOK")
	}

	return
}

func TestKODoLoginForm(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// An error is returned if caused by client policy (such as CheckRedirect),
		// or if there was an HTTP protocol error. A non-2xx response doesn't cause an error.
		time.Sleep(200 * time.Millisecond)
	}))
	defer ts.Close()
	// Change timeOut to make an HTTP protocol error timeout
	ts.Config.WriteTimeout = 20 * time.Millisecond
	client := http.Client{}

	username := "mangosteam"
	encryptedPassword := "123"
	rsatimestamp := "123"
	emailauthKeyedIn := ""
	captchaGID := ""
	captchaKeyedIn := ""
	twoFactorCode := ""
	mangosteam.BaseSteamWebURL = ts.URL

	_, err := DoLogin(
		&client,
		username,
		encryptedPassword,
		rsatimestamp,
		emailauthKeyedIn,
		captchaGID,
		captchaKeyedIn,
		twoFactorCode,
	)

	if err == nil {
		t.Errorf("Dologin returns no error when DoLogin Form failed")
	}

	return
}

func TestEmailauthNeeded(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getMockEmailauthNeeded())
	}))
	defer ts.Close()

	client := http.Client{}
	username := "mangosteam"
	encryptedPassword := "123"
	rsatimestamp := "123"
	emailauthKeyedIn := ""
	captchaGID := ""
	captchaKeyedIn := ""
	twoFactorCode := ""
	mangosteam.BaseSteamWebURL = ts.URL

	_, err := DoLogin(
		&client,
		username,
		encryptedPassword,
		rsatimestamp,
		emailauthKeyedIn,
		captchaGID,
		captchaKeyedIn,
		twoFactorCode,
	)

	if err == nil {
		t.Errorf("AuthLogin should have Email Needed")
	}

	return
}

func TestRespBody(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getMockRespBody())
	}))
	defer ts.Close()

	client := http.Client{}
	username := "mangosteam"
	encryptedPassword := "123"
	rsatimestamp := "123"
	emailauthKeyedIn := ""
	captchaGID := ""
	captchaKeyedIn := ""
	twoFactorCode := ""
	mangosteam.BaseSteamWebURL = ts.URL

	_, err := DoLogin(
		&client,
		username,
		encryptedPassword,
		rsatimestamp,
		emailauthKeyedIn,
		captchaGID,
		captchaKeyedIn,
		twoFactorCode,
	)

	if err == nil {
		t.Errorf("loginBody returns no error")
	}

	return
}

func TestOKIsLoggedIn(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := http.Client{}
	mangosteam.BaseSteamWebURL = ts.URL

	isLoggedIn, _ := IsLoggedIn(&client)

	if !isLoggedIn {
		t.Errorf("isLoggedIn should return true in case of status unauthorized")
	}

}

func TestKOIsLoggedIn(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	client := http.Client{}
	mangosteam.BaseSteamWebURL = ts.URL

	isLoggedIn, _ := IsLoggedIn(&client)
	if isLoggedIn {
		t.Errorf("isLoggedIn should return false in case of status unauthorized")
	}

}

func TestTimeOutIsLoggedIn(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		time.Sleep(200 * time.Millisecond)
	}))
	defer ts.Close()

	ts.Config.WriteTimeout = 20 * time.Millisecond
	client := http.Client{}
	mangosteam.BaseSteamWebURL = ts.URL

	isLoggedIn, err := IsLoggedIn(&client)
	if isLoggedIn || err == nil {
		t.Errorf("isLoggedIn should return false and an error in case of http error")
	}

}
