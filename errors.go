package main

import "errors"

var (
	errNoEntryFound        = errors.New("no entry found")
	errNoGroupFound        = errors.New("no group found")
	errPasswordMismatch    = errors.New("password and repetition do not match")
	errCredentialsMissing  = errors.New("key file or password has to be provided")
	errEntryTitleNotUnique = errors.New("entry title must be unique within a parent group")
	errGroupNameNotUnique  = errors.New("group name must be unique within a parent group")
)
