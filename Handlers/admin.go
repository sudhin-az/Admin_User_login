package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	db "user_admin/DB"
	models "user_admin/Models"
)

type PageData struct {
	PassError InvalidAdminSignupError
	PassTable models.User
}

type PageSearchData struct {
	UserAdminList []models.User
	SearchError   string
}

type InvalidAdminSignupError struct {
	InvalidUsername string
	InvalidPhone    string
	InvalidPassword string
	InvalidSignUp   string
	InvalidFullName string
}

var InvalidAdminData InvalidAdminSignupError

func Admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	if r.Method == http.MethodPost {
		_, err := r.Cookie("jwt_admin_token")
		if err != nil { //no cookie
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	fmt.Println("success in admin")
	fmt.Println("success in admin post")
	var UserAdminList []models.User
	db.DB.Find(&UserAdminList)
	data := PageSearchData{
		UserAdminList: UserAdminList,
	}

	if len(data.UserAdminList) == 0 {
		data.SearchError = "No users found"
	}

	tmp, err := template.ParseFiles("Templates/admin.html")
	if err != nil {
		fmt.Printf("error parsing template file: %v\n", err)
		return
	}
	err = tmp.Execute(w, data)
	if err != nil {
		log.Fatalf("error executing template: %v", err)
	}
}

func AdminAddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	fmt.Println("success in adminAddUser")

	if r.Method == http.MethodPost {
		fmt.Println("succes in adminadduser post")
		// parse form data
		if err := r.ParseForm(); err != nil {
			fmt.Println("error here", err)
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}

		adminSignup := models.User{
			FullName:    r.FormValue("formName"),
			UserName:    r.FormValue("formUsername"),
			Email:       r.FormValue("formEmail"),
			PhoneNumber: r.FormValue("formPhoneNumber"),
			Password:    r.FormValue("formPassword"),
			Gender:      r.FormValue("gender"),
		}
		fmt.Println("Parsed form data:", adminSignup)

		//Validate form data
		if adminSignup.FullName == "" {
			InvalidSignData.InvalidSignUp = ""
			InvalidSignData.InvalidUsername = ""
			InvalidSignData.InvalidPassword = ""
			InvalidSignData.InvalidFullName = "Full name should not be empty"
			fmt.Println("admin name us")
			tmp, err := template.ParseFiles("Templates/adminAddUserSignUp.html")
			if err != nil {
				log.Fatalf("error %v", err)
			}
			// fmt.Println("---------", h)
			tmp.ExecuteTemplate(w, "adminAddUserSignUp.html", InvalidAdminData)
			return
		}
		if adminSignup.UserName == "" {
			InvalidSignData.InvalidSignUp = ""
			InvalidSignData.InvalidPhone = ""
			InvalidSignData.InvalidPassword = ""
			InvalidSignData.InvalidUsername = "username should not be empty"
			fmt.Println("username")
			tmp, err := template.ParseFiles("Templates/adminAddUserSignUp.html")
			if err != nil {
				log.Fatalf("error %v", err)
			}
			// fmt.Println("---------", h)
			tmp.ExecuteTemplate(w, "adminAddUserSignUp.html", InvalidAdminData)
			return
		}
		if adminSignup.Password == "" {
			InvalidSignData.InvalidFullName = ""
			InvalidSignData.InvalidUsername = ""
			InvalidSignData.InvalidPhone = ""
			InvalidSignData.InvalidPassword = "password should not be empty"
			fmt.Println("password")
			tmp, err := template.ParseFiles("Templates/adminAddUserSignUp.html")
			if err != nil {
				log.Fatalf("error %v", err)
			}
			// fmt.Println("---------", h)
			tmp.ExecuteTemplate(w, "adminAddUserSignUp.html", InvalidAdminData)
			return
		}
		fmt.Println("length of phone", len(adminSignup.PhoneNumber))
		if len(adminSignup.PhoneNumber) != 10 {
			InvalidSignData.InvalidSignUp = ""
			InvalidSignData.InvalidUsername = ""
			InvalidSignData.InvalidPassword = ""
			InvalidSignData.InvalidPhone = "phone number must be 10 digits"
			fmt.Println(len(adminSignup.PhoneNumber))
			tmpl, err := template.ParseFiles("Templates/adminAddUserSignup.html")
			if err != nil {
				log.Fatal(err)
			}
			err = tmpl.ExecuteTemplate(w, "adminAddUserSignup.html", InvalidAdminData)
			if err != nil {
				fmt.Println("--------", err)
			}
			return
		}
		var count int64
		db.DB.Model(&models.User{}).Where("user_name = ?", adminSignup.UserName).Count(&count)
		if count != 0 {
			InvalidAdminData.InvalidUsername = "Already registered Username"
			InvalidAdminData.InvalidSignUp = ""
			tmpl, err := template.ParseFiles("Templates/adminAddUserSignup.html")
			if err != nil {
				log.Fatal(err)
			}
			tmpl.ExecuteTemplate(w, "adminAddUserSignup.html", InvalidAdminData)
			return
		}
		// querry := "INSERT INTO users (full_name, user_name, email, phone_number, password, gender) VALUES ($1,$2,$3,$4,$5,,$6)"
		// err := db.DB.Exec(querry, adminSignup.FullName, adminSignup.UserName, adminSignup.Email, adminSignup.PhoneNumber, adminSignup.Password, adminSignup.Gender).Error
		result := db.DB.Create(&adminSignup)
		if result.Error != nil {
			fmt.Println("Error executing query:", result.Error)
			http.Error(w, "Failed to insert data into database", http.StatusInternalServerError)
			return
		}
		fmt.Println("User added successfully, redirecting to /admin")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles("Templates/adminAddUserSignup.html")
	if err != nil {
		log.Fatalf("error %v", err)
	}
	tmpl.ExecuteTemplate(w, "adminAddUserSignup.html", nil)
	fmt.Println("hiii")
}

var StoreUsername string

func AdminUserUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("error here", err)
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}
	numValues := len(r.Form)
	fmt.Println("form value number", numValues)

	fmt.Println("helllooooooo")

	if r.Method == http.MethodPost && numValues == 1 {
		var PassData PageData
		var user models.User
		fmt.Println("success in adminUpdateUser post")

		fmt.Println(r.FormValue("usingNameToUpdate"))
		InvalidSignData.InvalidUsername = ""
		InvalidSignData.InvalidSignUp = ""

		fmt.Println("-------", r.FormValue("usingNameToUpdate"))

		if err := db.DB.Raw("SELECT * FROM users WHERE user_name = ?", r.FormValue("usingNameToUpdate")).Scan(&user).Error; err != nil {
			fmt.Println("Error", err)
			return
		}

		PassData = PageData{
			PassError: InvalidAdminSignupError(InvalidSignData),
			PassTable: user,
		}
		fmt.Println("not working", PassData.PassTable.FullName)
		tmpl, err := template.ParseFiles("Templates/adminUpdateUserSignUp.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, PassData); err != nil {
			fmt.Println("-----", err)
			return
		}
	}
	if r.Method == http.MethodPost && numValues > 1 {
		var PassData PageData
		var user models.User
		fmt.Println("2nd post")
		var signup models.User
		signup.FullName = r.FormValue("formName")
		signup.UserName = r.FormValue("formUsername")
		fmt.Println(signup.FullName)
		signup.Email = r.FormValue("formEmail")
		signup.PhoneNumber = r.FormValue("formPhoneNumber")
		signup.Password = r.FormValue("formPassword")
		signup.Gender = r.FormValue("gender")

		if signup.FullName == "" || signup.Email == "" || signup.PhoneNumber == "" || signup.Password == "" || signup.Gender == "" {
			InvalidSignData.InvalidSignUp = "Invalid data in sign up"
			InvalidSignData.InvalidUsername = "Invalid data in sign up"
			InvalidSignData.InvalidPhone = ""
			fmt.Println("does nil condition work")
			if err := db.DB.Raw("SELECT * FROM users WHERE user_name = ?", signup.UserName).Scan(&user).Error; err != nil {
				fmt.Println("Error:", err)
				return
			}
			PassData = PageData{
				PassError: InvalidAdminSignupError(InvalidSignData),
				PassTable: user,
			}
			fmt.Println("---", PassData.PassTable.FullName)
			fmt.Println("----", PassData.PassError.InvalidSignUp)
			tmpl, err := template.ParseFiles("Templates/adminUpdateUserSignUp.html")
			if err != nil {
				log.Fatalf("error %v", err)
			}
			err = tmpl.ExecuteTemplate(w, "adminUpdateUserSignUp.html", PassData)
			if err != nil {
				fmt.Println("----------", err)
				return
			}
			return
		}
		if len(signup.PhoneNumber) != 10 {
			InvalidSignData.InvalidSignUp = ""
			InvalidSignData.InvalidUsername = ""
			InvalidSignData.InvalidPhone = "phone number must be 10 digits"
			fmt.Println("does nil condition work")
			if err := db.DB.Raw("SELECT * FROM users WHERE user_name = ?", signup.UserName).Scan(&user).Error; err != nil {
				fmt.Println("Error:", err)
				return
			}
			PassData = PageData{
				PassError: InvalidAdminSignupError(InvalidSignData),
				PassTable: user,
			}
			fmt.Println("---", PassData.PassTable.FullName)
			fmt.Println("----", PassData.PassError.InvalidSignUp)
			tmp, err := template.ParseFiles("Templates/adinUpdateUserSignUp.html")
			if err != nil {
				log.Fatalf("error %v", err)
			}
			// fmt.Println("---------", h)
			err = tmp.ExecuteTemplate(w, "adminUpdateUserSignUp.html", PassData)
			if err != nil {
				fmt.Println("-------", err)
				return
			}
			return
		}
		fmt.Println("hiiii")
		fmt.Println("-------", signup.UserName)

		results := db.DB.Model(&models.User{}).Where("user_name = ?", signup.UserName).Updates(map[string]interface{}{
			"full_name":    signup.FullName,
			"email":        signup.Email,
			"phone_number": signup.PhoneNumber,
			"password":     signup.Password,
			"gender":       signup.Gender,
		})
		//fmt.Println("query executed")
		if results.Error != nil {
			fmt.Println("does query have error")
			fmt.Println("err:", results.Error)
		}
		fmt.Println("does it reach after error")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
	if r.Method == "GET" {
		fmt.Println("gee")
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}

func AdminSearchUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("success in adminSearchUser")

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Println("error here", err)
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}
		var userNaming []models.User
		query := db.DB.Where("user_name LIKE ?", "%"+r.FormValue("usernaming")+"%")
		query.Find(&userNaming)
		data := PageSearchData{
			UserAdminList: userNaming,
		}
		if len(userNaming) == 0 {
			data.SearchError = "No users found"
		}

		tmpl, err := template.ParseFiles("Templates/admin.html")
		if err != nil {
			fmt.Println("error in search parsing", err)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Println("error in search execute", err)
			return
		}
	}
}
func AdminUserDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method == http.MethodPost {
		fmt.Println("Success in adminDeleteUser")
		if err := r.ParseForm(); err != nil {
			fmt.Println("error here", err)
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}
		fmt.Println(r.FormValue("usingNameToDelete"))
		db.DB.Where("user_name", r.FormValue("usingNameToDelete")).Delete(&models.User{})
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
}
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("success in adminLogout")
	if r.Method == http.MethodPost {
		fmt.Println("inside post of admin logout")

		c = http.Cookie{Name: "jwt_admin_token", Value: "", Expires: time.Now().AddDate(0, 0, -1), MaxAge: -1}

		http.SetCookie(w, &c)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	fmt.Println("out of admin Logout")
}
