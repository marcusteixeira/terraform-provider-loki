package mimir

import (
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/timeinterval"
	"github.com/prometheus/common/model"
	"net"
	"net/url"
	"time"
)

func expandHTTPConfigOAuth2(v interface{}) *oauth2 {
	var oauth2Conf *oauth2
	data := v.([]interface{})
	if len(data) != 0 && data[0] != nil {
		oauth2Conf = &oauth2{}
		cfg := data[0].(map[string]interface{})
		oauth2Conf.ClientID = cfg["client_id"].(string)
		oauth2Conf.ClientSecret = cfg["client_secret"].(string)
		oauth2Conf.TokenURL = cfg["token_url"].(string)
		oauth2Conf.Scopes = expandStringArray(cfg["scopes"].([]interface{}))
		oauth2Conf.EndpointParams = expandStringMap(cfg["endpoint_params"].(map[string]interface{}))
	}
	return oauth2Conf
}

func flattenHTTPConfigOAuth2(v *oauth2) []interface{} {
	oauth2Conf := make(map[string]interface{})
	if v != nil {
		oauth2Conf["client_id"] = v.ClientID
		oauth2Conf["client_secret"] = v.ClientSecret
		oauth2Conf["token_url"] = v.TokenURL
		oauth2Conf["scopes"] = v.Scopes
		oauth2Conf["endpoint_params"] = v.EndpointParams
	}
	return []interface{}{oauth2Conf}
}

func expandHTTPConfigBasicAuth(v interface{}) *basicAuth {
	var basicAuthConf *basicAuth
	data := v.([]interface{})
	if len(data) != 0 && data[0] != nil {
		basicAuthConf := &basicAuth{}
		cfg := data[0].(map[string]interface{})
		basicAuthConf.Username = cfg["username"].(string)
		basicAuthConf.Password = cfg["password"].(string)
	}
	return basicAuthConf
}

func flattenHTTPConfigBasicAuth(v *basicAuth) []interface{} {
	basicAuthConf := make(map[string]interface{})
	if v != nil {
		basicAuthConf["username"] = v.Username
		basicAuthConf["password"] = v.Password
	}
	return []interface{}{basicAuthConf}
}

func expandHTTPConfigAuthorization(v interface{}) *authorization {
	var authConf *authorization
	data := v.([]interface{})
	if len(data) != 0 && data[0] != nil {
		authConf := &authorization{}
		cfg := data[0].(map[string]interface{})
		authConf.Type = cfg["type"].(string)
		authConf.Credentials = cfg["credentials"].(string)
	}
	return authConf
}

func flattenHTTPConfigAuthorization(v *authorization) []interface{} {
	authConf := make(map[string]interface{})
	if v != nil {
		authConf["type"] = v.Type
		authConf["credentials"] = v.Credentials
	}
	return []interface{}{authConf}
}

func expandTLSConfig(v interface{}) *tlsConfig {
	var tlsConf *tlsConfig
	data := v.([]interface{})
	if len(data) != 0 && data[0] != nil {
		cfg := data[0].(map[string]interface{})
		tlsConf.ServerName = cfg["server_name"].(string)
		tlsConf.InsecureSkipVerify = cfg["insecure_skip_verify"].(bool)
	}
	return tlsConf
}

func flattenTLSConfig(v *tlsConfig) []interface{} {
	tlsConf := make(map[string]interface{})
	tlsConf["server_name"] = v.ServerName
	tlsConf["insecure_skip_verify"] = v.InsecureSkipVerify
	return []interface{}{tlsConf}
}

func expandHTTPConfig(v interface{}) *httpClientConfig {
	var httpConf *httpClientConfig
	data := v.([]interface{})
	if len(data) != 0 && data[0] != nil {
		httpConf = &httpClientConfig{}
		cfg := data[0].(map[string]interface{})
		httpConf.ProxyURL = cfg["proxy_url"].(string)
		httpConf.FollowRedirects = new(bool)
		*httpConf.FollowRedirects = cfg["follow_redirects"].(bool)
		httpConf.BearerToken = cfg["bearer_token"].(string)

		if len(cfg["authorization"].([]interface{})) > 0 {
			httpConf.Authorization = expandHTTPConfigAuthorization(cfg["authorization"].([]interface{}))
		}

		if len(cfg["basic_auth"].([]interface{})) > 0 {
			httpConf.BasicAuth = expandHTTPConfigBasicAuth(cfg["basic_auth"].([]interface{}))
		}

		if len(cfg["oauth2"].([]interface{})) > 0 {
			httpConf.OAuth2 = expandHTTPConfigOAuth2(cfg["oauth2"].([]interface{}))
		}

		if len(cfg["tls_config"].([]interface{})) > 0 {
			httpConf.TLSConfig = expandTLSConfig(cfg["tls_config"].([]interface{}))
		}
	}

	return httpConf
}

func flattenHTTPConfig(v *httpClientConfig) []interface{} {
	httpConf := make(map[string]interface{})

	if v != nil {
		httpConf["proxy_url"] = v.ProxyURL
		httpConf["bearer_token"] = v.BearerToken

		if v.FollowRedirects != nil {
			httpConf["follow_redirects"] = v.FollowRedirects
		}

		if v.BasicAuth != nil {
			httpConf["basic_auth"] = flattenHTTPConfigBasicAuth(v.BasicAuth)
		}
		if v.OAuth2 != nil {
			httpConf["oauth2"] = flattenHTTPConfigOAuth2(v.OAuth2)
		}
		if v.Authorization != nil {
			httpConf["authorization"] = flattenHTTPConfigAuthorization(v.Authorization)
		}

		if v.TLSConfig != nil {
			httpConf["tls_config"] = flattenTLSConfig(v.TLSConfig)
		}
	}
	return []interface{}{httpConf}
}

func expandGlobalConfig(v interface{}) *globalConfig {
	var globalConf *globalConfig
	data := v.([]interface{})
	if len(data) != 0 && data[0] != nil {
		globalConf = &globalConfig{}
		cfg := data[0].(map[string]interface{})
		resolveTimeout, _ := model.ParseDuration(cfg["resolve_timeout"].(string))
		globalConf.ResolveTimeout = new(model.Duration)
		*globalConf.ResolveTimeout = resolveTimeout

		pagerdutyURL, _ := url.Parse(cfg["pagerduty_url"].(string))
		if pagerdutyURL.String() != "" {
			globalConf.PagerdutyURL = &config.URL{pagerdutyURL}
		}

		slackAPIURL, _ := url.Parse(cfg["slack_api_url"].(string))
		if slackAPIURL.String() != "" {
			globalConf.SlackAPIURL = &config.URL{slackAPIURL}
		}

		globalConf.SMTPFrom = cfg["smtp_from"].(string)
		globalConf.SMTPHello = cfg["smtp_hello"].(string)

		var hp config.HostPort
		hp.Host, hp.Port, _ = net.SplitHostPort(cfg["smtp_smarthost"].(string))
		globalConf.SMTPSmarthost = hp

		globalConf.SMTPAuthUsername = cfg["smtp_auth_username"].(string)
		globalConf.SMTPAuthPassword = cfg["smtp_auth_password"].(string)
		globalConf.SMTPAuthSecret = cfg["smtp_auth_secret"].(string)
		globalConf.SMTPAuthIdentity = cfg["smtp_auth_identity"].(string)
		globalConf.SMTPRequireTLS = new(bool)
		*globalConf.SMTPRequireTLS = cfg["smtp_require_tls"].(bool)
		globalConf.HTTPConfig = expandHTTPConfig(cfg["http_config"].(interface{}))

	}
	return globalConf
}

func flattenGlobalConfig(v *globalConfig) []interface{} {
	globalConf := make(map[string]interface{})

	if v != nil {
		if v.ResolveTimeout != nil {
			globalConf["resolve_timeout"] = v.ResolveTimeout.String()
		}

		if v.PagerdutyURL != nil {
			globalConf["pagerduty_url"] = v.PagerdutyURL.URL.String()
		}

		if v.SlackAPIURL != nil {
			globalConf["slack_api_url"] = v.SlackAPIURL.URL.String()
		}

		if v.HTTPConfig != nil {
			globalConf["http_config"] = flattenHTTPConfig(v.HTTPConfig)
		}
		globalConf["smtp_from"] = v.SMTPFrom
		globalConf["smtp_hello"] = v.SMTPHello
		globalConf["smtp_smarthost"] = v.SMTPSmarthost.String()
		globalConf["smtp_auth_username"] = v.SMTPAuthUsername
		globalConf["smtp_auth_password"] = v.SMTPAuthPassword
		globalConf["smtp_auth_secret"] = v.SMTPAuthSecret
		globalConf["smtp_auth_identity"] = v.SMTPAuthIdentity

		if v.SMTPRequireTLS != nil {
			globalConf["smtp_require_tls"] = v.SMTPRequireTLS
		}
	}

	return []interface{}{globalConf}
}

func expandReceiverConfig(v []interface{}) []*receiver {
	var receiverConf []*receiver

	for _, v := range v {
		cfg := &receiver{}
		data := v.(map[string]interface{})

		if raw, ok := data["name"]; ok {
			cfg.Name = raw.(string)
		}
		if raw, ok := data["pagerduty_configs"]; ok {
			cfg.PagerdutyConfigs = expandPagerdutyConfig(raw.([]interface{}))
		}
		if raw, ok := data["email_configs"]; ok {
			cfg.EmailConfigs = expandEmailConfig(raw.([]interface{}))
		}
		if raw, ok := data["wechat_configs"]; ok {
			cfg.WeChatConfigs = expandWeChatConfig(raw.([]interface{}))
		}
		if raw, ok := data["webhook_configs"]; ok {
			cfg.WebhookConfigs = expandWebhookConfig(raw.([]interface{}))
		}
		if raw, ok := data["pushover_configs"]; ok {
			cfg.PushoverConfigs = expandPushoverConfig(raw.([]interface{}))
		}
		receiverConf = append(receiverConf, cfg)
	}
	return receiverConf
}

func flattenReceiverConfig(v []*receiver) []interface{} {
	var receiverConf []interface{}

	if v == nil {
		return receiverConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["name"] = v.Name
		cfg["pagerduty_configs"] = flattenPagerdutyConfig(v.PagerdutyConfigs)
		cfg["email_configs"] = flattenEmailConfig(v.EmailConfigs)
		cfg["wechat_configs"] = flattenWeChatConfig(v.WeChatConfigs)
		cfg["webhook_configs"] = flattenWebhookConfig(v.WebhookConfigs)
		cfg["pushover_configs"] = flattenPushoverConfig(v.PushoverConfigs)
		receiverConf = append(receiverConf, cfg)
	}
	return receiverConf
}

func expandWebhookConfig(v []interface{}) []*webhookConfig {
	var webhookConf []*webhookConfig

	for _, v := range v {
		cfg := &webhookConfig{}
		data := v.(map[string]interface{})
		if raw, ok := data["send_resolved"]; ok {
			cfg.VSendResolved = new(bool)
			*cfg.VSendResolved = raw.(bool)
		}
		if raw, ok := data["http_config"]; ok {
			cfg.HTTPConfig = expandHTTPConfig(raw.(interface{}))
		}
		if raw, ok := data["url"]; ok {
			cfg.URL = raw.(string)
		}
		if raw, ok := data["max_alerts"]; ok {
			cfg.MaxAlerts = int32(raw.(int))
		}
		webhookConf = append(webhookConf, cfg)
	}
	return webhookConf
}

func flattenWebhookConfig(v []*webhookConfig) []interface{} {
	var webhookConf []interface{}

	if v == nil {
		return webhookConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["send_resolved"] = v.VSendResolved
		if v.HTTPConfig != nil {
			cfg["http_config"] = flattenHTTPConfig(v.HTTPConfig)
		}
		cfg["url"] = v.URL
		cfg["max_alerts"] = v.MaxAlerts
		webhookConf = append(webhookConf, cfg)
	}
	return webhookConf
}

func expandWeChatConfig(v []interface{}) []*weChatConfig {
	var weChatConf []*weChatConfig

	for _, v := range v {
		cfg := &weChatConfig{}
		data := v.(map[string]interface{})
		if raw, ok := data["send_resolved"]; ok {
			cfg.VSendResolved = new(bool)
			*cfg.VSendResolved = raw.(bool)
		}
		if raw, ok := data["http_config"]; ok {
			cfg.HTTPConfig = expandHTTPConfig(raw.(interface{}))
		}
		if raw, ok := data["api_secret"]; ok {
			cfg.APISecret = raw.(string)
		}
		if raw, ok := data["api_url_url"]; ok {
			cfg.APIURL = raw.(string)
		}
		if raw, ok := data["corp_id"]; ok {
			cfg.CorpID = raw.(string)
		}
		if raw, ok := data["agent_id"]; ok {
			cfg.AgentID = raw.(string)
		}
		if raw, ok := data["to_user"]; ok {
			cfg.ToUser = raw.(string)
		}
		if raw, ok := data["to_party"]; ok {
			cfg.ToParty = raw.(string)
		}
		if raw, ok := data["to_tag"]; ok {
			cfg.ToTag = raw.(string)
		}
		if raw, ok := data["message"]; ok {
			cfg.Message = raw.(string)
		}
		if raw, ok := data["message_type"]; ok {
			cfg.MessageType = raw.(string)
		}
		weChatConf = append(weChatConf, cfg)
	}
	return weChatConf
}

func flattenWeChatConfig(v []*weChatConfig) []interface{} {
	var weChatConf []interface{}

	if v == nil {
		return weChatConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["send_resolved"] = v.VSendResolved
		if v.HTTPConfig != nil {
			cfg["http_config"] = flattenHTTPConfig(v.HTTPConfig)
		}
		cfg["api_secret"] = v.APISecret
		cfg["api_url"] = v.APIURL
		cfg["corp_id"] = v.CorpID
		cfg["agent_id"] = v.AgentID
		cfg["to_user"] = v.ToUser
		cfg["to_party"] = v.ToParty
		cfg["to_tag"] = v.ToTag
		cfg["message"] = v.Message
		cfg["message_type"] = v.MessageType
		weChatConf = append(weChatConf, cfg)
	}
	return weChatConf
}

func expandEmailConfig(v []interface{}) []*emailConfig {
	var emailConf []*emailConfig

	for _, v := range v {
		cfg := &emailConfig{}
		data := v.(map[string]interface{})

		if raw, ok := data["send_resolved"]; ok {
			cfg.VSendResolved = new(bool)
			*cfg.VSendResolved = raw.(bool)
		}
		if raw, ok := data["to"]; ok {
			cfg.To = raw.(string)
		}
		if raw, ok := data["from"]; ok {
			cfg.From = raw.(string)
		}
		if raw, ok := data["hello"]; ok {
			cfg.Hello = raw.(string)
		}
		if raw, ok := data["smarthost"]; ok {
			var hp config.HostPort
			hp.Host, hp.Port, _ = net.SplitHostPort(raw.(string))
			cfg.Smarthost = hp
		}

		if raw, ok := data["auth_username"]; ok {
			cfg.AuthUsername = raw.(string)
		}
		if raw, ok := data["auth_password"]; ok {
			cfg.AuthPassword = raw.(string)
		}
		if raw, ok := data["auth_secret"]; ok {
			cfg.AuthSecret = raw.(string)
		}
		if raw, ok := data["auth_identity"]; ok {
			cfg.AuthIdentity = raw.(string)
		}
		if raw, ok := data["headers"]; ok {
			cfg.Headers = expandStringMap(raw.(map[string]interface{}))
		}
		if raw, ok := data["html"]; ok {
			cfg.HTML = raw.(string)
		}
		if raw, ok := data["text"]; ok {
			cfg.Text = raw.(string)
		}
		if raw, ok := data["require_tls"]; ok {
			cfg.RequireTLS = new(bool)
			*cfg.RequireTLS = raw.(bool)
		}
		emailConf = append(emailConf, cfg)
	}
	return emailConf
}

func flattenEmailConfig(v []*emailConfig) []interface{} {
	var emailConf []interface{}

	if v == nil {
		return emailConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["send_resolved"] = v.VSendResolved
		cfg["to"] = v.To
		cfg["from"] = v.From
		cfg["hello"] = v.Hello
		cfg["smarthost"] = v.Smarthost.String()
		cfg["auth_username"] = v.AuthUsername
		cfg["auth_password"] = v.AuthPassword
		cfg["auth_secret"] = v.AuthSecret
		cfg["auth_identity"] = v.AuthIdentity
		cfg["headers"] = v.Headers
		cfg["html"] = v.HTML
		cfg["text"] = v.Text
		cfg["require_tls"] = v.RequireTLS
		emailConf = append(emailConf, cfg)
	}
	return emailConf
}

func expandPagerdutyConfigLinks(v []interface{}) []pagerdutyLink {
	var pagerdutyLinkConf []pagerdutyLink

	for _, v := range v {
		var cfg pagerdutyLink
		data := v.(map[string]interface{})

		if raw, ok := data["text"]; ok {
			cfg.Text = raw.(string)
		}
		if raw, ok := data["href"]; ok {
			cfg.Href = raw.(string)
		}
		pagerdutyLinkConf = append(pagerdutyLinkConf, cfg)
	}
	return pagerdutyLinkConf
}

func expandPagerdutyConfigImages(v []interface{}) []pagerdutyImage {
	var pagerdutyImageConf []pagerdutyImage

	for _, v := range v {
		var cfg pagerdutyImage
		data := v.(map[string]interface{})

		if raw, ok := data["src"]; ok {
			cfg.Src = raw.(string)
		}
		if raw, ok := data["alt"]; ok {
			cfg.Alt = raw.(string)
		}
		if raw, ok := data["href"]; ok {
			cfg.Href = raw.(string)
		}
		pagerdutyImageConf = append(pagerdutyImageConf, cfg)
	}
	return pagerdutyImageConf
}

func expandPagerdutyConfig(v []interface{}) []*pagerdutyConfig {
	var pagerdutyConf []*pagerdutyConfig

	for _, v := range v {
		cfg := &pagerdutyConfig{}
		data := v.(map[string]interface{})

		if raw, ok := data["send_resolved"]; ok {
			cfg.VSendResolved = new(bool)
			*cfg.VSendResolved = raw.(bool)
		}
		if raw, ok := data["routing_key"]; ok {
			cfg.RoutingKey = raw.(string)
		}
		if raw, ok := data["service_key"]; ok {
			cfg.ServiceKey = raw.(string)
		}
		if raw, ok := data["url"]; ok {
			cfg.URL = raw.(string)
		}
		if raw, ok := data["http_config"]; ok {
			cfg.HTTPConfig = expandHTTPConfig(raw.(interface{}))
		}
		if raw, ok := data["images"]; ok {
			cfg.Images = expandPagerdutyConfigImages(raw.([]interface{}))
		}
		if raw, ok := data["links"]; ok {
			cfg.Links = expandPagerdutyConfigLinks(raw.([]interface{}))
		}
		if raw, ok := data["client"]; ok {
			cfg.Client = raw.(string)
		}
		if raw, ok := data["client_url"]; ok {
			cfg.ClientURL = raw.(string)
		}
		if raw, ok := data["description"]; ok {
			cfg.Description = raw.(string)
		}
		if raw, ok := data["severity"]; ok {
			cfg.Severity = raw.(string)
		}
		if raw, ok := data["class"]; ok {
			cfg.Class = raw.(string)
		}
		if raw, ok := data["component"]; ok {
			cfg.Component = raw.(string)
		}
		if raw, ok := data["group"]; ok {
			cfg.Group = raw.(string)
		}
		if raw, ok := data["details"]; ok {
			cfg.Details = expandStringMap(raw.(map[string]interface{}))
		}
		pagerdutyConf = append(pagerdutyConf, cfg)
	}
	return pagerdutyConf
}

func flattenPagerdutyConfig(v []*pagerdutyConfig) []interface{} {
	var pagerdutyConf []interface{}

	if v == nil {
		return pagerdutyConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["send_resolved"] = v.VSendResolved
		cfg["service_key"] = v.ServiceKey
		cfg["routing_key"] = v.RoutingKey
		if v.HTTPConfig != nil {
			cfg["http_config"] = flattenHTTPConfig(v.HTTPConfig)
		}
		cfg["url"] = v.URL
		cfg["client"] = v.Client
		cfg["client_url"] = v.ClientURL
		cfg["description"] = v.Description
		cfg["severity"] = v.Severity
		cfg["class"] = v.Class
		cfg["component"] = v.Component
		cfg["group"] = v.Group
		cfg["details"] = v.Details
		pagerdutyConf = append(pagerdutyConf, cfg)
	}
	return pagerdutyConf
}

func expandPushoverConfig(v []interface{}) []*pushoverConfig {
	var pushoverConf []*pushoverConfig

	for _, v := range v {
		cfg := &pushoverConfig{}
		data := v.(map[string]interface{})

		if raw, ok := data["send_resolved"]; ok {
			cfg.VSendResolved = new(bool)
			*cfg.VSendResolved = raw.(bool)
		}
		if raw, ok := data["http_config"]; ok {
			cfg.HTTPConfig = expandHTTPConfig(raw.(interface{}))
		}
		if raw, ok := data["user_key"]; ok {
			cfg.UserKey = raw.(string)
		}
		if raw, ok := data["token"]; ok {
			cfg.Token = raw.(string)
		}
		if raw, ok := data["title"]; ok {
			cfg.Title = raw.(string)
		}
		if raw, ok := data["message"]; ok {
			cfg.Message = raw.(string)
		}
		if raw, ok := data["url"]; ok {
			cfg.URL = raw.(string)
		}
		if raw, ok := data["url_title"]; ok {
			cfg.URLTitle = raw.(string)
		}
		if raw, ok := data["sound"]; ok {
			cfg.Sound = raw.(string)
		}
		if raw, ok := data["priority"]; ok {
			cfg.Priority = raw.(string)
		}
		if raw, ok := data["retry"]; ok {
			retry, _ := time.ParseDuration(raw.(string))
			cfg.Retry = retry
		}
		if raw, ok := data["expire"]; ok {
			expire, _ := time.ParseDuration(raw.(string))
			cfg.Expire = expire
		}
		if raw, ok := data["html"]; ok {
			cfg.HTML = raw.(bool)
		}
		pushoverConf = append(pushoverConf, cfg)
	}
	return pushoverConf
}

func flattenPushoverConfig(v []*pushoverConfig) []interface{} {
	var pushoverConf []interface{}

	if v == nil {
		return pushoverConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["send_resolved"] = v.VSendResolved
		if v.HTTPConfig != nil {
			cfg["http_config"] = flattenHTTPConfig(v.HTTPConfig)
		}
		cfg["user_key"] = v.UserKey
		cfg["token"] = v.Token
		cfg["title"] = v.Title
		cfg["message"] = v.Message
		cfg["url"] = v.URL
		cfg["url_title"] = v.URLTitle
		cfg["sound"] = v.Sound
		cfg["priority"] = v.Priority
		cfg["retry"] = v.Retry.String()
		cfg["expire"] = v.Expire.String()
		cfg["html"] = v.HTML

		pushoverConf = append(pushoverConf, cfg)
	}
	return pushoverConf
}

func expandRouteConfig(v interface{}) *route {
	routeConf := &route{}
	data := v.([]interface{})
	if len(data) != 0 && data[0] != nil {
		cfg := data[0].(map[string]interface{})
		if raw, ok := cfg["receiver"]; ok {
			routeConf.Receiver = raw.(string)
		}
		if raw, ok := cfg["group_by"]; ok {
			routeConf.GroupByStr = expandStringArray(raw.([]interface{}))
		}
		if raw, ok := cfg["matchers"]; ok {
			routeConf.Matchers = expandStringArray(raw.([]interface{}))
		}
		if raw, ok := cfg["continue"]; ok {
			routeConf.Continue = raw.(bool)
		}
		if raw, ok := cfg["child_route"]; ok {
			var routes []*route
			for _, item := range raw.([]interface{}) {
				routes = append(routes, expandRouteConfig([]interface{}{item.(map[string]interface{})}))
			}
			routeConf.Routes = routes
		}
		if raw, ok := cfg["group_wait"]; ok {
			routeConf.GroupWait = raw.(string)
		}
		if raw, ok := cfg["group_interval"]; ok {
			routeConf.GroupInterval = raw.(string)
		}
		if raw, ok := cfg["repeat_interval"]; ok {
			routeConf.RepeatInterval = raw.(string)
		}
		if raw, ok := cfg["mute_time_intervals"]; ok {
			routeConf.MuteTimeIntervals = expandStringArray(raw.([]interface{}))
		}
		if raw, ok := cfg["active_time_intervals"]; ok {
			routeConf.ActiveTimeIntervals = expandStringArray(raw.([]interface{}))
		}
	}
	return routeConf
}

func flattenRouteConfig(v *route) []interface{} {
	routeConf := make(map[string]interface{})

	routeConf["receiver"] = v.Receiver

	if len(v.GroupByStr) > 0 {
		routeConf["group_by"] = v.GroupByStr
	}

	if len(v.Matchers) > 0 {
		routeConf["matchers"] = v.Matchers
	}

	if v.Routes != nil {
		var routes []interface{}
		for _, route := range v.Routes {
			routes = append(routes, flattenRouteConfig(route)[0])
		}
		routeConf["child_route"] = routes
	}
	routeConf["continue"] = v.Continue
	routeConf["group_wait"] = v.GroupWait
	routeConf["group_interval"] = v.GroupInterval
	routeConf["repeat_interval"] = v.RepeatInterval

	if len(v.MuteTimeIntervals) > 0 {
		routeConf["mute_time_intervals"] = v.MuteTimeIntervals
	}
	if len(v.ActiveTimeIntervals) > 0 {
		routeConf["active_time_intervals"] = v.ActiveTimeIntervals
	}

	return []interface{}{routeConf}
}

func expandInhibitRuleConfig(v []interface{}) []*inhibitRule {
	var inhibitRuleConf []*inhibitRule

	for _, v := range v {
		cfg := &inhibitRule{}
		data := v.(map[string]interface{})

		if raw, ok := data["source_matchers"]; ok {
			cfg.SourceMatchers = expandStringArray(raw.([]interface{}))
		}
		if raw, ok := data["target_matchers"]; ok {
			cfg.TargetMatchers = expandStringArray(raw.([]interface{}))
		}
		if raw, ok := data["equal"]; ok {
			cfg.Equal = expandStringArray(raw.([]interface{}))
		}
		inhibitRuleConf = append(inhibitRuleConf, cfg)
	}
	return inhibitRuleConf
}

func flattenInhibitRuleConfig(v []*inhibitRule) []interface{} {
	var inhibitRuleConf []interface{}

	if v == nil {
		return inhibitRuleConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["source_matchers"] = v.SourceMatchers
		cfg["target_matchers"] = v.TargetMatchers
		cfg["equal"] = v.Equal
		inhibitRuleConf = append(inhibitRuleConf, cfg)
	}
	return inhibitRuleConf
}

func expandMuteTimeIntervalConfig(v []interface{}) []*muteTimeInterval {
	var muteTimeIntervalConf []*muteTimeInterval

	for _, v := range v {
		cfg := &muteTimeInterval{}
		data := v.(map[string]interface{})

		if raw, ok := data["name"]; ok {
			cfg.Name = raw.(string)
		}
		if raw, ok := data["time_intervals"]; ok {
			cfg.TimeIntervals = expandTimeIntervalConfig(raw.([]interface{}))
		}
		muteTimeIntervalConf = append(muteTimeIntervalConf, cfg)
	}
	return muteTimeIntervalConf
}

func flattenMuteTimeIntervalConfig(v []*muteTimeInterval) []interface{} {
	var muteTimeIntervalConf []interface{}

	if v == nil {
		return muteTimeIntervalConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["name"] = v.Name
		cfg["time_intervals"] = flattenTimeIntervalConfig(v.TimeIntervals)
		muteTimeIntervalConf = append(muteTimeIntervalConf, cfg)
	}
	return muteTimeIntervalConf
}

func expandTimeIntervalConfig(v []interface{}) []timeinterval.TimeInterval {
	var timeIntervalConf []timeinterval.TimeInterval

	for _, v := range v {
		var cfg timeinterval.TimeInterval
		data := v.(map[string]interface{})

		if raw, ok := data["times"]; ok {
			cfg.Times = expandTimeRange(raw.([]interface{}))
		}
		if raw, ok := data["weekdays"]; ok {
			cfg.Weekdays = expandWeekdayRange(raw.([]interface{}))
		}
		if raw, ok := data["days_of_month"]; ok {
			cfg.DaysOfMonth = expandDayOfMonthRange(raw.([]interface{}))
		}
		if raw, ok := data["months"]; ok {
			cfg.Months = expandMonthRange(raw.([]interface{}))
		}
		if raw, ok := data["years"]; ok {
			cfg.Years = expandYearRange(raw.([]interface{}))
		}
		timeIntervalConf = append(timeIntervalConf, cfg)
	}
	return timeIntervalConf
}

func expandTimeRange(v []interface{}) []timeinterval.TimeRange {
	var timeRangeConf []timeinterval.TimeRange

	for _, v := range v {
		var cfg timeinterval.TimeRange
		data := v.(map[string]interface{})

		if raw, ok := data["start_minute"]; ok {
			cfg.StartMinute = raw.(int)
		}
		if raw, ok := data["end_minute"]; ok {
			cfg.EndMinute = raw.(int)
		}

		timeRangeConf = append(timeRangeConf, cfg)
	}
	return timeRangeConf
}

func flattenTimeRange(v []timeinterval.TimeRange) []interface{} {
	var timeRangeConf []interface{}

	if v == nil {
		return timeRangeConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["start_minute"] = v.StartMinute
		cfg["end_minute"] = v.EndMinute
		timeRangeConf = append(timeRangeConf, cfg)
	}
	return timeRangeConf
}

func expandWeekdayRange(v []interface{}) []timeinterval.WeekdayRange {
	var inclusiveRangeConf []timeinterval.WeekdayRange

	for _, v := range v {
		var cfg timeinterval.WeekdayRange
		data := v.(map[string]interface{})

		if raw, ok := data["begin"]; ok {
			cfg.Begin = raw.(int)
		}
		if raw, ok := data["end"]; ok {
			cfg.End = raw.(int)
		}

		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func flattenWeekdayRange(v []timeinterval.WeekdayRange) []interface{} {
	var inclusiveRangeConf []interface{}

	if v == nil {
		return inclusiveRangeConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["begin"] = v.Begin
		cfg["end"] = v.End
		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func expandDayOfMonthRange(v []interface{}) []timeinterval.DayOfMonthRange {
	var inclusiveRangeConf []timeinterval.DayOfMonthRange

	for _, v := range v {
		var cfg timeinterval.DayOfMonthRange
		data := v.(map[string]interface{})

		if raw, ok := data["begin"]; ok {
			cfg.Begin = raw.(int)
		}
		if raw, ok := data["end"]; ok {
			cfg.End = raw.(int)
		}

		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func flattenDayOfMonthRange(v []timeinterval.DayOfMonthRange) []interface{} {
	var inclusiveRangeConf []interface{}

	if v == nil {
		return inclusiveRangeConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["begin"] = v.Begin
		cfg["end"] = v.End
		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func expandMonthRange(v []interface{}) []timeinterval.MonthRange {
	var inclusiveRangeConf []timeinterval.MonthRange

	for _, v := range v {
		var cfg timeinterval.MonthRange
		data := v.(map[string]interface{})

		if raw, ok := data["begin"]; ok {
			cfg.Begin = raw.(int)
		}
		if raw, ok := data["end"]; ok {
			cfg.End = raw.(int)
		}

		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func flattenMonthRange(v []timeinterval.MonthRange) []interface{} {
	var inclusiveRangeConf []interface{}

	if v == nil {
		return inclusiveRangeConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["begin"] = v.Begin
		cfg["end"] = v.End
		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func expandYearRange(v []interface{}) []timeinterval.YearRange {
	var inclusiveRangeConf []timeinterval.YearRange

	for _, v := range v {
		var cfg timeinterval.YearRange
		data := v.(map[string]interface{})

		if raw, ok := data["begin"]; ok {
			cfg.Begin = raw.(int)
		}
		if raw, ok := data["end"]; ok {
			cfg.End = raw.(int)
		}

		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func flattenYearRange(v []timeinterval.YearRange) []interface{} {
	var inclusiveRangeConf []interface{}

	if v == nil {
		return inclusiveRangeConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["begin"] = v.Begin
		cfg["end"] = v.End
		inclusiveRangeConf = append(inclusiveRangeConf, cfg)
	}
	return inclusiveRangeConf
}

func flattenTimeIntervalConfig(v []timeinterval.TimeInterval) []interface{} {
	var timeIntervalConf []interface{}

	if v == nil {
		return timeIntervalConf
	}

	for _, v := range v {
		cfg := make(map[string]interface{})
		cfg["times"] = flattenTimeRange(v.Times)
		cfg["weekdays"] = flattenWeekdayRange(v.Weekdays)
		cfg["days_of_month"] = flattenDayOfMonthRange(v.DaysOfMonth)
		cfg["months"] = flattenMonthRange(v.Months)
		cfg["years"] = flattenYearRange(v.Years)
		timeIntervalConf = append(timeIntervalConf, cfg)
	}
	return timeIntervalConf
}
