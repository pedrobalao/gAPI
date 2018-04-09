package config

var CONFIGS_LOCATION = "./configs/"

var SERVICE_DISCOVERY_CONFIG_FILE = "services.json"
var AUTHENTICATION_CONFIG_FILE = "oauth.json"

var POST = "POST"
var GET = "GET"
var DELETE = "DELETE"
var PUT = "PUT"
var PATCH = "PATCH"


func LoadConfigs(){
	LoadGApiConfig()
	LoadURLConstants()
}