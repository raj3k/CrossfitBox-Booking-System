package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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
