// Code generated by go-swagger; DO NOT EDIT.

package v1

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// PodDNSConfig PodDNSConfig defines the DNS parameters of a pod in addition to those generated from DNSPolicy.
//
// swagger:model PodDNSConfig
type PodDNSConfig struct {

	// A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed.
	Nameservers []string `json:"nameservers,omitempty"`

	// A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy.
	Options []*PodDNSConfigOption `json:"options,omitempty"`

	// A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed.
	Searches []string `json:"searches,omitempty"`
}
