package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	db "user_admin/DB"
	models "user_admin/Models"
	auth "user_admin/helpers"
)

type errors struct {
	UsernameError string
	PasswordError string
}

type SignupError struct {
	InvalidUsername string
	InvalidPhone    string
	InvalidPassword string
	InvalidSignUp   string
	InvalidFullName string
}

type home struct {
	Username2 string
}

var InvalidSignData SignupError
var errorV errors
var h home
var username string
var sessions = make(map[string]string)

var c http.Cookie

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		w.Header().Set("Cache-Control", "no-store")

		var username1 string
		var password1 string
		if err := r.ParseForm(); err != nil {
			fmt.Println("error here", err)
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username)
		fmt.Println(password)

		if username == "sudhin" && password == "sudhin" {
			errorV.UsernameError = ""
			errorV.PasswordError = ""
			tokenString, err := auth.GenerateJWT(username, "sudhin")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error creating token:%v", err)
				return
			}
			fmt.Printf("Token: %s\n", tokenString)
			cookie := http.Cookie{
				Name:     "jwt_admin_token",
				Value:    tokenString,
				Expires:  time.Now().Add(24 * time.Hour), //Set cookie expiration time
				HttpOnly: true,                           //Cookie accessible only by HTTP requests
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		var count int64
		db.DB.Model(&models.User{}).Where("user_name = ?", username).Count(&count)
		if count != 0 {
			fmt.Println("username exists")
			db.DB.Model(&models.User{}).Where("user_name = ?", username).Pluck("user_name", &username1)
			db.DB.Model(&models.User{}).Where("user_name = ?", username).Pluck("password", &password1)
			fmt.Println("username1", username1)
			fmt.Println("password1", password1)
		}
		fmt.Println("------", username1)
		fmt.Println("-----", password1)

		if username1 == username && password1 == password && username != "" && password != "" {
			errorV.UsernameError = ""
			errorV.PasswordError = ""
			tokenString, err := auth.GenerateJWT(username, "user")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error creating token %v", err)
				return
			}
			fmt.Printf("Token: %s\n", tokenString)
			cookie := http.Cookie{
				Name:     "jwt_token",
				Value:    tokenString,
				Expires:  time.Now().Add(24 * time.Hour), // Set cookie expiration time
				HttpOnly: true,                           //Cookie accessible only by HTTP requests
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else if username1 != username && password1 == password {
			errorV.UsernameError = "Invalid username"
			errorV.PasswordError = ""
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else if password1 != password && username1 == username {
			errorV.PasswordError = "Invalid password"
			errorV.UsernameError = ""
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			errorV.UsernameError = "Invalid username"
			errorV.PasswordError = "Invalid password"
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
	fmt.Println("success in roothandler")
	fmt.Println(username)

	_, err := r.Cookie("jwt_token")
	if err == nil { //means cookie here
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	a, err := r.Cookie("jwt_admin_token")
	if err == nil { // means cookie here
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
	fmt.Println(a)

	// Pass the claims to the template
	tmpl, err := template.ParseFiles("Templates/index.html")
	if err != nil {
		log.Fatalf("error %v", err)
	}
	tmpl.ExecuteTemplate(w, "index.html", errorV)
	fmt.Println("out of rootHandler")
}

// func SignupHandler(w http.ResponseWriter, r *http.Request) {

// 	tmpl, err := template.ParseFiles("Templates/signup.html")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tmpl.ExecuteTemplate(w, "signup.html", nil)
// 	if err := r.ParseForm(); err != nil {
// 		fmt.Println("error here", err)
// 		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Println("hi")
// 	models.A1.UserName = r.FormValue("formName")
// 	fmt.Println("hello")
// 	models.A1.Email = r.FormValue("formEmail")
// 	models.A1.PhoneNumber = r.FormValue("formPhoneNumber")
// 	models.A1.Password = r.FormValue("formPassword")
// 	models.A1.Gender = r.FormValue("gender")
// 	fmt.Println("hi hello")
// 	query := "INSERT INTO users (full_name,user_name, email, phone_number, password, gender) VALUES ($1,$2,$3,$4,$5,$6)"
// 	fmt.Println("hello hi")
// 	db.Db.Exec(query, models.A1.FullName, models.A1.UserName, models.A1.Email, models.A1.PhoneNumber, models.A1.Password, models.A1.Gender)
// 	fmt.Println("error here")
// 	if err != nil {
// 		fmt.Println("err: please print", err)
// 		tmpl, err := template.ParseFiles("Templates/signup.html")
// 		if err != nil {
// 			log.Fatalf("error %v", err)
// 		}
// 		tmpl.ExecuteTemplate(w, "signup.html", nil)
// 		return
// 	}

// }
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	errorV.UsernameError = ""
	errorV.PasswordError = ""
	fmt.Println("Success in signupHandler")

	if r.Method == http.MethodPost {
		w.Header().Set("Cache-Control", "no-store")
		// Parse form data
		if err := r.ParseForm(); err != nil {
			fmt.Println("error here", err)
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}

		var signup models.User // Extract form values
		signup.FullName = r.FormValue("formName")
		signup.UserName = r.FormValue("formUsername")
		signup.Email = r.FormValue("formEmail")
		signup.PhoneNumber = r.FormValue("formPhoneNumber")
		signup.Password = r.FormValue("formPassword")
		signup.Gender = r.FormValue("gender")
		fmt.Println(signup)
		//Validate form data
		if signup.FullName == "" || signup.UserName == "" || signup.Email == "" || signup.PhoneNumber == "" || signup.Password == "" || signup.Gender == "" {
			InvalidSignData.InvalidSignUp = "All fields are required"
			tmpl, err := template.ParseFiles("Templates/signup.html")
			if err != nil {
				log.Fatal(err)
			}
			tmpl.ExecuteTemplate(w, "signup.html", InvalidSignData)
			return
		}

		if signup.Password != r.FormValue("formConfirmPassword") {
			InvalidSignData.InvalidSignUp = ""
			InvalidSignData.InvalidUsername = ""
			InvalidSignData.InvalidPhone = ""
			InvalidSignData.InvalidPassword = "two password must match"
			tmpl, err := template.ParseFiles("Templates/signup.html")
			if err != nil {
				log.Fatal(err)
			}
			tmpl.ExecuteTemplate(w, "signup.html", nil)
			return
		}
		fmt.Println("length of phone", len(signup.PhoneNumber))
		if len(signup.PhoneNumber) != 10 {
			InvalidSignData.InvalidSignUp = ""
			InvalidSignData.InvalidUsername = ""
			InvalidSignData.InvalidPassword = ""
			InvalidSignData.InvalidPhone = "phone number must be 10 digits"
			tmpl, err := template.ParseFiles("Templates/signup.html")
			if err != nil {
				log.Fatal(err)
			}
			tmpl.ExecuteTemplate(w, "signup.html", InvalidSignData)
			return
		}
		var count int64
		db.DB.Model(&models.User{}).Where("user_name = ?", signup.UserName).Count(&count)
		if count != 0 {
			InvalidSignData.InvalidUsername = "Already registered Username"
			InvalidSignData.InvalidSignUp = ""
			tmpl, err := template.ParseFiles("Templates/signup.html")
			if err != nil {
				log.Fatal(err)
			}
			tmpl.ExecuteTemplate(w, "signup.html", InvalidSignData)
			return
		}
		// fmt.Println(models.A1.FullName)
		// fmt.Println(models.A1.UserName)
		// fmt.Println(models.A1.Email)
		// fmt.Println(models.A1.PhoneNumber)
		// fmt.Println(models.A1.Password)
		// fmt.Println(models.A1.Gender)
		// Insert into database
		fmt.Println(signup)
		// query := "INSERT INTO users (full_name, user_name, email, phone_number, password, gender) VALUES ($1, $2, $3, $4, $5, $6)"
		result := db.DB.Save(&signup)
		if result.Error != nil {
			fmt.Println("Error executing query:", result.Error)
			http.Error(w, "Failed to insert data into database", http.StatusInternalServerError)
			return
		}

		// Redirect to a success page or login page after successful registration
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	_, err := r.Cookie("jwt_token")
	if err == nil { // cookie here
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	a, err := r.Cookie("jwt_admin_token")
	if err == nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
	fmt.Println(a)

	// If GET method, simply render the signup page
	tmpl, err := template.ParseFiles("Templates/signup.html")
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Println("----------", h)
	tmpl.ExecuteTemplate(w, "signup.html", nil)
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	fmt.Println("success in homehandler")
	fmt.Println(username)

	_, err := r.Cookie("jwt_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	h.Username2 = username

	tmpl, err := template.ParseFiles("Templates/home.html")
	if err != nil {
		log.Fatalf("error %v", err)
	}
	fmt.Println("----------", h)
	tmpl.ExecuteTemplate(w, "home.html", h)
	fmt.Println("out of homehandler")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		c = http.Cookie{Name: "jwt_token", Value: "", Expires: time.Now().AddDate(0, 0, -1), MaxAge: -1}

		http.SetCookie(w, &c)

		delete(sessions, username)

		fmt.Println("in logoutHandler")
		fmt.Printf("%v\n", sessions)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "404 not found", http.StatusNotFound)
	}
}
