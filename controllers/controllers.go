package controllers

// the http response cache
var cache map[string]interface{}

func Init() {
  // set up the cache
  cache = make(map[string]interface{})
}
