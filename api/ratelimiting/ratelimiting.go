package ratelimiting

import (
	"sync"
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/utils"
	"time"
	"github.com/qiangxue/fasthttp-routing"
)

type LimiterRate struct {
	Period time.Duration
	Limit int
}

type RateStatus struct {
	NumberRequests int
	ExpirationTime int64
}

var limiter config.GApiRateLimitingConfig 
var limits map[string]RateStatus

type Updater struct {
	NewRate RateStatus
	Service servicediscovery.Service
	ReqName string
}

var sd servicediscovery.ServiceDiscovery
var rateLimitingMutex sync.RWMutex

func InitRateLimiting() {
	limiter = config.GApiConfiguration.RateLimiting
	limits = make(map[string]RateStatus)
	sd = *servicediscovery.GetServiceDiscoveryObject()
}


func RateLimiting(c *routing.Context) error {
	if ! limiter.Active {
		return nil
	}

	rateLimitingMutex.Lock()
	
	currentRequestMetricName := GetRateLimitingMetricName(c, limiter)

	service := serviceForUri(c)

	IncrementRateLimiting(currentRequestMetricName, service)

	rateStatus := RateLimitingStatusForRequest(currentRequestMetricName)

	if IsRateLimitExceeded(rateStatus, service) {
		http.Response(c, `{"error":true, "msg": "Rate limiting exceeded."}`, 429, c.Request.URI().String())
		c.Abort()
		return nil
	}

	rate := RateStatus{NumberRequests: (rateStatus.NumberRequests + 1), ExpirationTime: rateStatus.ExpirationTime}

	limits[currentRequestMetricName] = rate
	
    rateLimitingMutex.Unlock()

	return nil
}

func IncrementRateLimiting(currentRequestMetricName string, service servicediscovery.Service) {
	if _, ok := limits[currentRequestMetricName]; ok == false {
		limits[currentRequestMetricName] = RateStatus{NumberRequests:1, ExpirationTime: RateLimitingExpirationTime(service)}
	}
}

func RateLimitingStatusForRequest(currentRequestMetricName string) RateStatus {
	currentNumberRequests := limits[currentRequestMetricName].NumberRequests
	currentExpirationTime := limits[currentRequestMetricName].ExpirationTime

	if currentExpirationTime < utils.CurrentTimeMilliseconds() {
		currentExpirationTime = RateLimitingExpirationTime(servicediscovery.Service{})
		currentNumberRequests = 0
	}

	return RateStatus{NumberRequests: currentNumberRequests, ExpirationTime: currentExpirationTime}
}

func IsRateLimitExceeded(rateStatus RateStatus, service servicediscovery.Service) bool {
	// If rate limit time expired
	if rateStatus.ExpirationTime < utils.CurrentTimeMilliseconds() {
		return true
	}

	// Check rate limit for custom service rate
	if service.RateLimit > 0 && rateStatus.NumberRequests > service.RateLimit {
		return true
	}

	// Check for general rate limit
	if rateStatus.NumberRequests > limiter.Limit {
		return true
	}
	return false
}

func RateLimitingExpirationTime(service servicediscovery.Service) int64 {
	if service.RateLimitExpirationTime > 0 {
		return utils.CurrentTimeMilliseconds() + (service.RateLimitExpirationTime * 60 * 1000) 
	}
	return utils.CurrentTimeMilliseconds() + (limiter.Period * 60 * 1000)
}