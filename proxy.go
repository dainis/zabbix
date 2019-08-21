package zabbix

import (
	"github.com/AlekSi/reflector"
)

// https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/definitions
type Proxy struct {
	ProxyId    string        `json:"proxyid,omitempty"`
	Host       string        `json:"host"`
	Status     StatusType    `json:"status"`
	LastAccess TimestampType `json:"lastaccess"`

	// Fields below used only when creating proxies
	Interfaces HostInterfaces `json:"interfaces,omitempty"`

	// https://www.zabbix.com/documentation/3.0/manual/appendix/api/proxy/definitions
	Description    string `json:"description"`
	TlsConnect     int    `json:"tls_connect"`
	TlsAccept      int    `json:"tls_accept"`
	TlsPskIdentity string `json:"tls_psk_identity"`
	TlsPsk         string `json:"tls_psk"`
}

type Proxies []Proxy

// Wrapper for proxy.get: https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/get
func (api *API) ProxiesGet(params Params) (res Proxies, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("proxy.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets proxy by Id only if there is exactly 1 matching proxy.
func (api *API) ProxyGetById(id string) (res *Proxy, err error) {
	proxies, err := api.ProxiesGet(Params{"proxyids": id})
	if err != nil {
		return
	}

	if len(proxies) == 1 {
		res = &proxies[0]
	} else {
		e := ExpectedOneResult(len(proxies))
		err = &e
	}
	return
}

// Gets proxy by Proxy only if there is exactly 1 matching proxy.
func (api *API) ProxyGetByProxy(proxy string) (res *Proxy, err error) {
	proxies, err := api.ProxiesGet(Params{"filter": map[string]string{"proxy": proxy}})
	if err != nil {
		return
	}

	if len(proxies) == 1 {
		res = &proxies[0]
	} else {
		e := ExpectedOneResult(len(proxies))
		err = &e
	}
	return
}

// Wrapper for proxy.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/create
func (api *API) ProxiesCreate(proxies Proxies) (err error) {
	response, err := api.CallWithError("proxy.create", proxies)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	for i, id := range proxyids {
		proxies[i].ProxyId = id.(string)
	}
	return
}

func (api *API) ProxiesUpdate(proxies Proxies) (err error) {
	_, err = api.CallWithError("proxy.update", proxies)
	return
}

// Wrapper for proxy.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/delete
// Cleans ProxyId in all proxies elements if call succeed.
func (api *API) ProxiesDelete(proxies Proxies) (err error) {
	ids := make([]string, len(proxies))
	for i, proxy := range proxies {
		ids[i] = proxy.ProxyId
	}

	err = api.ProxiesDeleteByIds(ids)
	if err == nil {
		for i := range proxies {
			proxies[i].ProxyId = ""
		}
	}
	return
}

// Wrapper for proxy.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/delete
func (api *API) ProxiesDeleteByIds(ids []string) (err error) {
	proxyIds := make([]map[string]string, len(ids))
	for i, id := range ids {
		proxyIds[i] = map[string]string{"proxyid": id}
	}

	response, err := api.CallWithError("proxy.delete", proxyIds)
	if err != nil {
		// Zabbix 2.4 uses new syntax only
		if e, ok := err.(*Error); ok && e.Code == -32602 {
			response, err = api.CallWithError("proxy.delete", ids)
		}
	}
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	if len(ids) != len(proxyids) {
		err = &ExpectedMore{len(ids), len(proxyids)}
	}
	return
}
