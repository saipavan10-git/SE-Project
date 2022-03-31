package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

var results string
var i string = "\"\""
var keyword int
var code int

func (app *application) getOneVaccine(w http.ResponseWriter, r *http.Request) {
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()
	params := httprouter.ParamsFromContext(r.Context())
	fmt.Println(params)
	id := params.ByName("id")

	var p1 []models.Vaccine
	db.Where("name = ?", id).Find(&p1)
	app.writeJSON(w, http.StatusOK, p1, "vaccines")
	log.Println(p1)

}
func (app *application) getOneVaccineProcess() ([]models.Vaccine, error) {
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()

	var p5 []models.Vaccine
	err := db.Find(&p5).Error
	fmt.Println(p5)
	if err != nil {
		return nil, err
	}
	return p5, nil
}

func (app *application) getAllVaccines(w http.ResponseWriter, r *http.Request) {
	p5, err := app.getAllVaccinesProcess()
	if err != nil {
		app.errorJSON(w, err)
	}
	app.writeJSON(w, http.StatusOK, p5, "vaccines")
}

func (app *application) getAllVaccinesProcess() ([]models.Vaccine, error) {
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()

	var p5 []models.Vaccine
	err := db.Find(&p5).Error
	fmt.Println(p5)
	if err != nil {
		return nil, err
	}
	return p5, nil
}

func (app *application) getBooking(w http.ResponseWriter, r *http.Request) {
	code = 0
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()
	dbAppoint, _ := gorm.Open("sqlite3", "./appointment.db")
	defer dbAppoint.Close()
	dbAppoint.AutoMigrate(&models.UserAppoint{})
	time.Sleep(3 * time.Second)
	var vaccine models.Vaccine
	log.Println("2")

	err := json.NewDecoder(r.Body).Decode(&vaccine)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println(vaccine)

	if i != "\"\"" {
		db.Model(&vaccine).Update("available", 0)
		log.Println("Vaccine with ID ", vaccine.ID, " is booked.")
		log.Println("Detailed Info:", vaccine.Name, ",", vaccine.VaccineNum, "-dose.")
		log.Println("Availability: ", vaccine.Available)

		trimmedString := strings.Trim(i, "\"")

		n := true
		var num int
		var checkUser models.UserAppoint
		for n {
			min := 100000
			max := 999999
			num = rand.Intn(max-min) + min

			dbAppoint.Where("code = ?", num).First(&checkUser)

			if checkUser.Code == 0 {
				n = false
			}
		}

		appointment := models.UserAppoint{
			Email: trimmedString,
			ID:    vaccine.ID,
			Code:  num,
		}
		code = num

		dbAppoint.Create(&appointment)
		i = "\"\""
		app.writeJSON(w, 200, &appointment, "Great")
	} else {
		app.errorJSON(w, errors.New("no user"))
	}
}

func (app *application) recordSignup(w http.ResponseWriter, r *http.Request) {
	db, _ := gorm.Open("sqlite3", "./user.db")
	defer db.Close()
	db.AutoMigrate(&models.User{})
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	encryptedUser := models.User{
		Email:    user.Email,
		Password: string(encryptedPassword),
		Fname:    user.Fname,
		Lname:    user.Lname,
	}

	db.Create(&encryptedUser)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	db, _ := gorm.Open("sqlite3", "./user.db")
	defer db.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var findUser models.User
	db.Where("email = ?", user.Email).Take(&findUser)

	var empty models.User
	log.Println(findUser)

	if findUser == empty {
		app.errorJSON(w, errors.New("user not found"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(user.Password)); err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	} else {

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.Email,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		})

		token, _ := claims.SignedString([]byte(SecretKey))

		if err != nil {
			app.errorJSON(w, errors.New("could not log in"))
			return
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    token,
				Expires:  time.Now().Add(time.Hour * 24),
				HttpOnly: true,
			})
			app.writeJSON(w, http.StatusOK, token, "token")

		}
	}

}

func (app *application) user(w http.ResponseWriter, r *http.Request) {
	db, _ := gorm.Open("sqlite3", "./user.db")
	defer db.Close()
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			// app.errorJSON(w, err)
			// w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value

	tkn, err := jwt.ParseWithClaims(tknStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := tkn.Claims.(*jwt.StandardClaims)
	var user models.User

	db.Where("email = ?", claims.Issuer).Take(&user)
	log.Println("123123123", user)
	i = claims.Issuer
	app.writeJSON(w, http.StatusOK, user, "message")
}

func (app *application) searchRecord(w http.ResponseWriter, r *http.Request) {
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()

	var result models.SearchVaccine
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	results = result.Result

}

func (app *application) searchResult(w http.ResponseWriter, r *http.Request) {
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()

	var result models.SearchVaccine
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var p1 []models.Vaccine
	db.Where("name = ?", results).Find(&p1)
	app.writeJSON(w, http.StatusOK, p1, "vaccines")
	log.Println(p1)

}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {

	i = "\"\""

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	log.Println(i)
	//return
}

func (app *application) receiveFront(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if string(body) != "\"\"" {
		i = string(body)
		log.Println(i)
		log.Println("1")
	} else {
		app.errorJSON(w, errors.New("please sign in"))
	}
}

func (app *application) getAppoint(w http.ResponseWriter, r *http.Request) {

	dbAppoint, _ := gorm.Open("sqlite3", "./appointment.db")
	defer dbAppoint.Close()
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()
	var user models.UserAppoint

	time.Sleep(1 * time.Second)
	if i != "\"\"" {
		dbAppoint.Where("email = ?", i).Take(&user)
		var vaccine models.Vaccine
		db.Where("id = ?", user.ID).Take(&vaccine)
		log.Println("Hey", vaccine)
		app.writeJSON(w, 200, vaccine, "message")

	} else {
		app.errorJSON(w, errors.New("no appointment"))

	}

}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

	db, _ := gorm.Open("sqlite3", "./user.db")
	defer db.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println(user.Birthdate)
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	encryptedUser := models.User{
		Email:     user.Email,
		Password:  string(encryptedPassword),
		Fname:     user.Fname,
		Lname:     user.Lname,
		Birthdate: user.Birthdate,
		SSN:       user.SSN,
	}

	log.Println(encryptedUser)
	db.Model(&user).Where("email = ?", encryptedUser.Email).Update("fname", encryptedUser.Fname)
	db.Model(&user).Where("email = ?", encryptedUser.Email).Update("lname", encryptedUser.Lname)
	db.Model(&user).Where("email = ?", encryptedUser.Email).Update("birthdate", encryptedUser.Birthdate)
	db.Model(&user).Where("email = ?", encryptedUser.Email).Update("ssn", encryptedUser.SSN)

	if len(user.Password) != 0 {
		db.Model(&user).Where("email = ?", encryptedUser.Email).Update("password", encryptedUser.Password)
	}
}

func (app *application) deleteBooking(w http.ResponseWriter, r *http.Request) {
	dbAppoint, _ := gorm.Open("sqlite3", "./appointment.db")
	defer dbAppoint.Close()
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var user models.UserAppoint

	log.Println(string(body))
	AllList := strings.Fields(string(body))

	log.Println(AllList[0])
	log.Println(AllList[1])
	dbAppoint.Where("email = ?", AllList[0]).Delete(&user)

	var vaccine models.Vaccine
	db.Where("id = ?", AllList[1]).Find(&vaccine)
	db.Model(&vaccine).Update("available", 1)
}

func (app *application) searchCode(w http.ResponseWriter, r *http.Request) {
	dbAppoint, _ := gorm.Open("sqlite3", "./appointment.db")
	defer dbAppoint.Close()
	db, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer db.Close()

	var i models.Search

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println(i.Value)
	result, _ := strconv.Atoi(i.Value)
	log.Println(result)
	keyword = result
}

func (app *application) displayCert(w http.ResponseWriter, r *http.Request) {
	dbAppoint, _ := gorm.Open("sqlite3", "./appointment.db")
	defer dbAppoint.Close()
	dbUser, _ := gorm.Open("sqlite3", "./user.db")
	defer dbUser.Close()
	dbVaccine, _ := gorm.Open("sqlite3", "./vaccine.db")
	defer dbVaccine.Close()

	var appoint models.UserAppoint
	var user models.User
	var vaccine models.Vaccine

	dbAppoint.Where("code = ?", keyword).First(&appoint)
	log.Printf("Code is %d", keyword)
	dbVaccine.Where("id = ?", appoint.ID).First(&vaccine)
	dbUser.Where("email = ?", appoint.Email).First(&user)

	certificate := models.Cert{
		Email:      appoint.Email,
		Fname:      user.Fname,
		Lname:      user.Lname,
		Birthdate:  user.Birthdate,
		SSN:        user.SSN,
		Code:       appoint.Code,
		Name:       vaccine.Name,
		VaccineNum: vaccine.VaccineNum,
		State:      vaccine.State,
		ZipCode:    vaccine.ZipCode,
	}

	app.writeJSON(w, 200, &certificate, "Great")

}

func (app *application) code(w http.ResponseWriter, r *http.Request) {

	newCode := fmt.Sprintf("%d", code)
	log.Println(newCode)
	appoint := models.Search{
		Value: newCode,
	}

	app.writeJSON(w, 200, &appoint, "Great")

}
