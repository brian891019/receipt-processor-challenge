package model

import "errors"

var ErrNotFound = errors.New("not found")
var ErrEmptyField = errors.New("receipt field shouldn't be empty")
var ErrValidateReceipt = errors.New("error validating receipt")
var ErrCalculatePoint = errors.New("error calculating points")
var ErrInvalidTotal = errors.New("invalid total")
var ErrInvalidItemPrice = errors.New("invalid item Price")
var ErrInvalidPurchaseTime = errors.New("invalid purchaseTime")
var ErrInvalidPurchaseDate = errors.New("invalid PurchaseDate")
