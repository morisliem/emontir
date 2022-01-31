package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type (
	UserBaseModel struct {
		ID          string         `db:"id"`
		Name        string         `db:"name"`
		Email       string         `db:"email"`
		Password    string         `db:"password"`
		Address     sql.NullString `db:"address"`
		PhoneNumber sql.NullString `db:"phone_num"`
		IsActive    bool           `db:"is_active"`
	}

	RegisterUser struct {
		ID       string `db:"id"`
		Name     string `db:"name"`
		Email    string `db:"email"`
		Password string `db:"password"`
	}

	LoginUser struct {
		Email    string `db:"email"`
		Password string `db:"password"`
	}

	UserLocation struct {
		ID            string    `db:"id"`
		Label         string    `db:"label"`
		Address       string    `db:"address"`
		AddressDetail string    `db:"address_detail"`
		PhoneNumber   string    `db:"phone_num"`
		RecipientName string    `db:"recipient_name"`
		Latitude      string    `db:"latitude"`
		Longitude     string    `db:"longitude"`
		CreatedAt     time.Time `db:"created_at"`
	}
)

type User interface {
	RegisterUser(ctx context.Context, param *RegisterUser) error
	IsEmailUsed(ctx context.Context, email string) (bool, error)
	ActivateEmail(ctx context.Context, email string) error
	GetUserByEmail(ctx context.Context, email string) (*UserBaseModel, error)
	GetUserCurrentLocation(ctx context.Context, userID string) (*UserLocation, error)
	GetListOfUserLocation(ctx context.Context, userID string) ([]UserLocation, error)
	AddUserLocation(ctx context.Context, userID string, param *UserLocation) error
	StoreFCMKey(ctx context.Context, userID, fmc string) error
	GetFCMKey(ctx context.Context, userID string) (string, error)
	GetUserIDNOrderIDByInvoiceID(ctx context.Context, invoiceID string) (string, string, error)
	GetUserIDByOrderID(ctx context.Context, orderID string) (string, error)
	// SetEmailTimeLimit(email, userID string) error
	// IsEmailStillValid(email string) error
}

type user struct {
	db *sqlx.DB
	// redis   *redis.Client
	queries map[string]*sqlx.Stmt
}

func NewUser(db *sqlx.DB) User {
	user := new(user)
	user.db = db
	// user.redis = redis
	user.queries = make(map[string]*sqlx.Stmt, len(userQueries))
	for k, v := range userQueries {
		stmt, err := db.Preparex(v)
		if err != nil {
			log.Fatal().Msg("error : " + err.Error() + "\nuser : " + v)
		}
		user.queries[k] = stmt
	}
	return user
}

var (
	userActivateEmail    = "ActivateEmail"
	userActivateEmailSQL = `UPDATE "users" SET is_active = $2 WHERE email = $1`

	userSetNewUser    = "AddNewUser"
	userSetNewUserSQL = `INSERT INTO "users" (id, name, email, password, is_active) VALUES ($1,$2,$3,$4,$5)`

	userGetUserByEmail    = "GetUserByEmail"
	userGetUserByEmailSQL = `SELECT "id", "is_active", "password" from "users" WHERE email = $1`

	getUserIDByInvoiceID    = "getUserByInvoiceID"
	getUserIDByInvoiceIDSQL = `SELECT "user_id", "id" from "orders" WHERE "invoice_id" = $1`

	GetUserIDByOrderID    = "getUserByOrderID"
	getUserIDByOrderIDSQL = `SELECT "user_id" from "orders" WHERE "id" = $1`

	userIsEmailUsed    = "IsEmailUsed"
	userIsEmailUsedSQL = `SELECT "email" FROM "users" WHERE email = $1`

	setFCMKey    = "setFCMKey"
	setFCMKeySQL = `UPDATE "users" SET "fcm_key" = $2 WHERE "id" = $1`

	getFCMKey    = "getFCMKey"
	getFCMKeySQL = `SELECT "fcm_key" FROM "users" WHERE "id" = $1`

	getUserLocation          = "getUserLocation"
	userLocField1            = `"id", "label", "address", "address_detail", `
	userLocField2            = `"phone_num", "recipient_name", "latitude", "longitude"`
	getUserLocFields         = userLocField1 + userLocField2
	getUserLocationCondition = `WHERE "user_id" = $1 ORDER BY "created_at" DESC`
	getUserLocationSQL       = `SELECT ` + getUserLocFields + ` FROM "user_addresses" ` + getUserLocationCondition

	setUserLocation       = "setUserLocation"
	setUserLocField1      = `"id", "user_id", "label", "address", "address_detail", "phone_num",`
	setUserLocField2      = `"recipient_name", "created_at", "latitude", "longitude"`
	setUserLocationFields = `(` + setUserLocField1 + setUserLocField2 + `)`
	setUserLocValue       = `VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	setUserLocationSQL    = `INSERT INTO "user_addresses" ` + setUserLocationFields + setUserLocValue

	userQueries = map[string]string{
		userSetNewUser:       userSetNewUserSQL,
		userIsEmailUsed:      userIsEmailUsedSQL,
		userActivateEmail:    userActivateEmailSQL,
		userGetUserByEmail:   userGetUserByEmailSQL,
		getUserLocation:      getUserLocationSQL,
		setUserLocation:      setUserLocationSQL,
		setFCMKey:            setFCMKeySQL,
		getFCMKey:            getFCMKeySQL,
		getUserIDByInvoiceID: getUserIDByInvoiceIDSQL,
		getReviewByOrderID:   getUserIDByOrderIDSQL,
	}
)

func (c *user) RegisterUser(ctx context.Context, param *RegisterUser) error {
	_, err := c.queries[userSetNewUser].ExecContext(ctx, param.ID, param.Name, param.Email, param.Password, false)
	if err != nil {
		return err
	}
	return nil
}

func (c *user) IsEmailUsed(ctx context.Context, email string) (bool, error) {
	var tmp string
	if err := c.queries[userIsEmailUsed].QueryRow(email).Scan(&tmp); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *user) ActivateEmail(ctx context.Context, email string) error {
	_, err := c.queries[userActivateEmail].ExecContext(ctx, email, true)

	if err != nil {
		return err
	}

	return nil
}

func (c *user) GetUserByEmail(ctx context.Context, email string) (*UserBaseModel, error) {
	var result UserBaseModel
	if err := c.queries[userGetUserByEmail].GetContext(ctx, &result, email); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *user) GetUserCurrentLocation(ctx context.Context, userID string) (*UserLocation, error) {
	var userLoc UserLocation
	if err := c.queries[getUserLocation].GetContext(ctx, &userLoc, userID); err != nil {
		return nil, err
	}
	return &userLoc, nil
}

func (c *user) AddUserLocation(ctx context.Context, userID string, param *UserLocation) error {
	// nolint(gosec) // false positive
	_, err := c.queries[setUserLocation].ExecContext(ctx, param.ID, userID, param.Label, param.Address, param.AddressDetail, param.PhoneNumber, param.RecipientName, time.Now(), param.Latitude, param.Longitude)
	if err != nil {
		return err
	}
	return nil
}

func (c *user) GetListOfUserLocation(ctx context.Context, userID string) ([]UserLocation, error) {
	var result []UserLocation
	err := c.queries[getUserLocation].SelectContext(ctx, &result, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *user) StoreFCMKey(ctx context.Context, userID, fcm string) error {
	_, err := c.queries[setFCMKey].ExecContext(ctx, userID, fcm)
	if err != nil {
		return err
	}
	return nil
}

func (c *user) GetFCMKey(ctx context.Context, userID string) (string, error) {
	var fcmKey sql.NullString
	err := c.queries[getFCMKey].QueryRowContext(ctx, userID).Scan(&fcmKey)
	if err != nil {
		return "", err
	}
	return fcmKey.String, nil
}

func (c *user) GetUserIDNOrderIDByInvoiceID(ctx context.Context, invoiceID string) (string, string, error) {
	var userID, orderID sql.NullString
	err := c.queries[getUserIDByInvoiceID].QueryRowContext(ctx, invoiceID).Scan(&userID, &orderID)
	if err != nil {
		return "", "", err
	}
	return userID.String, orderID.String, nil
}

func (c *user) GetUserIDByOrderID(ctx context.Context, orderID string) (string, error) {
	var userID sql.NullString
	err := c.queries[getUserIDByInvoiceID].QueryRowContext(ctx, orderID).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID.String, nil
}

// func (c *user) SetEmailTimeLimit(email, userID string) error {
// 	err := c.redis.Set(email, userID, time.Minute*1).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *user) IsEmailStillValid(email string) error {
// 	_, err := c.redis.Get(email).Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
