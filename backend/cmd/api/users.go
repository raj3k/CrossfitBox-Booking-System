package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"time"

	"crossfitbox.booking.system/internal/cookies"
	"crossfitbox.booking.system/internal/data"
	"crossfitbox.booking.system/internal/tokens"
	"crossfitbox.booking.system/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serveErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationErrors(w, r, v.Errors)
		return
	}

	err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exist")
			app.failedValidationErrors(w, r, v.Errors)
		default:
			app.serveErrorResponse(w, r, err)
		}
		return
	}

	otp, err := tokens.GenerateOTP()
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	err = app.storeInRedis("activation_", otp.Hash, user.ID, app.config.tokenExpiration.duration)
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	now := time.Now()
	expiration := now.Add(app.config.tokenExpiration.duration)
	exact := expiration.Format(time.RFC1123)

	app.background(func() {
		data := map[string]interface{}{
			"token":       tokens.FormatOTP(otp.Secret),
			"firstName":   user.FirstName,
			"userID":      user.ID,
			"frontendURL": app.config.frontendURL,
			"expiration":  app.config.tokenExpiration.durationString,
			"exact":       exact,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
		app.logger.PrintInfo(fmt.Sprintf("Email sent to %s", user.ID), nil)
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	user, err := app.models.User.GetByEmail(input.Email, true)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serveErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serveErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	var userID = data.UserID{
		Id: user.ID,
	}

	var buf bytes.Buffer

	err = gob.NewEncoder(&buf).Encode(&userID)
	if err != nil {
		app.serveErrorResponse(w, r, err)
		return
	}

	session := buf.String()

	// Store session in redis
	err = app.storeInRedis("sessionid_", session, userID.Id, app.config.secret.sessionExpiration)
	if err != nil {
		app.logError(r, err)
	}

	cookie := http.Cookie{
		Name:     "sessionid",
		Value:    session,
		Path:     "/",
		MaxAge:   int(app.config.secret.sessionExpiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err = cookies.WriteEncrypted(w, cookie, app.config.secret.secretKey)
	if err != nil {
		app.serveErrorResponse(w, r, errors.New("something happened setting your cookie data"))
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
	}
}

func (app *application) currentUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, status, err := app.extractParamsFromSession(r)
	if err != nil {
		switch *status {
		case http.StatusUnauthorized:
			app.unauthorizedResponse(w, r, err)
		case http.StatusBadRequest:
			app.badRequestResponse(w, r, err)
		case http.StatusInternalServerError:
			app.serveErrorResponse(w, r, err)
		default:
			app.serveErrorResponse(w, r, errors.New("something happened and we could not fullfil your request at the moment"))
		}
		return
	}

	// Get session from redis
	_, err = app.getFromRedis(fmt.Sprintf("sessionid_%s", userID.Id))
	if err != nil {
		app.unauthorizedResponse(w, r, errors.New("you are not authorized to access this resource"))
		return
	}

	user, err := app.models.User.Get(userID.Id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serveErrorResponse(w, r, err)
		return
	}
}
