package main 

// Using a custom type is okay because context keys are of type "any"
// By using custom types we are available to avoid collisions between strings for commong keys 
// like "isAuthenticated"
type contextKey string 

const isAuthenticatedContextKey = contextKey("isAuthenticated")