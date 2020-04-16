package utils

import (
	"reflect"
  "strings"
)

const (
	//READ is used to address read operation
	READ = "read"
	//WRITE is is used to address write operation
	WRITE = "write"
)

//FilterRead is used to filter output to onlly output what client can see based on scopes
func FilterRead(data interface{}, scopesAllowed []string) {
	valueOf := reflect.Indirect(reflect.ValueOf(data))
	typeOf := valueOf.Type()
	fieldNum := valueOf.NumField()

	for i := 0; i < fieldNum; i++ {
		curField := typeOf.Field(i)
		scopeRequired := curField.Tag.Get("readScope")
		if !ScopeInAllowed(scopeRequired, scopesAllowed) {
			field := valueOf.Field(i)
			field.Set(reflect.Zero(field.Type()))
		}
	}
}

// strech converts a slice to bigger length by adding toAdd to it
func strech(scopeSlice []string, toAdd string, toLen int) []string {
	lenDiff := toLen - len(scopeSlice)
	for i := 0; i < lenDiff; i++ {
		scopeSlice = append(scopeSlice, toAdd)
	}
	return scopeSlice
}

// MatchScopes matches two scopes using Wildcard Scope Matching Strategy (asymetric)
func MatchScopes(scopeA, scopeB string) bool {
	scopeASplit := strings.Split(scopeA, ":")
	scopeBSplit := strings.Split(scopeB, ":")
	scopeALen := len(scopeASplit)
	scopeBLen := len(scopeBSplit)

	// If scopeBLen is smaller than scopeALen and last char of scopeB is not * return false
	if scopeBLen < scopeALen && scopeBSplit[scopeBLen-1] != "*" {
		return false
		// If scopeBLen is smaller than scopeALen and last char of scopeB is * stretch scopeB To Len Of ScopeA By Adding "*"
	} else if scopeBLen < scopeALen && scopeBSplit[scopeBLen-1] == "*" {
		scopeBSplit = strech(scopeBSplit, "*", scopeALen)
		// If scopeBLen is greater than scopeALen and last char of scopeA is not * return false
	} else if scopeBLen > scopeALen && scopeASplit[scopeALen-1] != "*" {
		return false
	}

	for i := 0; i < scopeALen; i++ {
		if !(scopeASplit[i] == scopeBSplit[i] || scopeBSplit[i] == "*") {
			return false
		}
	}

	return true
}

// ScopeInAllowed is used to check if scope is allowed based on allowed scopes list
func ScopeInAllowed(scope string, allowedScopes []string) bool {
	for _, allowedScope := range allowedScopes {
		if MatchScopes(scope, allowedScope) {
			return true
		}
	}
	return false
}
