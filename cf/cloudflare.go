package cf

import (
	"fmt"
	"github.com/cloudflare-go"
)

type CloudflareAPI struct {
	API        *cloudflare.API
	ZoneID     string
	Host       string
	IPListName string
}

func NewCloudflareAPI(APIKey, zoneID, host, ipListName string) (*CloudflareAPI, error) {
	api, err := cloudflare.NewWithAPIToken(APIKey)
	if err != nil {
		return nil, err
	}

	return &CloudflareAPI{
		API:        api,
		ZoneID:     zoneID,
		Host:       host,
		IPListName: ipListName,
	}, nil
}

func (c *CloudflareAPI) AddIPToWhitelist(ip string) error {
	err := c.API.CreateFirewallAccessRule(c.ZoneID, cloudflare.FirewallAccessRule{
		Mode:  "whitelist",
		Notes: c.IPListName,
		Configuration: cloudflare.FirewallAccessRuleConfiguration{
			Value:    ip,
			Target:   "ip",
			ZoneName: c.Host,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *CloudflareAPI) UpdateIPWhitelist(oldIP, newIP string) error {
	rules, err := c.API.ListFirewallAccessRules(c.ZoneID, cloudflare.FirewallAccessRuleConfiguration{
		Value:    oldIP,
		Target:   "ip",
		ZoneName: c.Host,
	})
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Notes == c.IPListName {
			err := c.API.UpdateFirewallAccessRule(c.ZoneID, rule.ID, cloudflare.FirewallAccessRule{
				Mode:  "whitelist",
				Notes: c.IPListName,
				Configuration: cloudflare.FirewallAccessRuleConfiguration{
					Value:    newIP,
					Target:   "ip",
					ZoneName: c.Host,
				},
			})
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("no matching whitelist rule found for IP %s", oldIP)
}

func (c *CloudflareAPI) RemoveIPFromWhitelist(ip string) error {
	rules, err := c.API.ListFirewallAccessRules(c.ZoneID, cloudflare.FirewallAccessRuleConfiguration{
		Value:    ip,
		Target:   "ip",
		ZoneName: c.Host,
	})
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Notes == c.IPListName {
			err := c.API.DeleteFirewallAccessRule(c.ZoneID, rule.ID)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("no matching whitelist rule found for IP %s", ip)
}

// 添加 IP 到白名单
func addIPToFilter(api *cloudflare.API, zoneID, host, ip string) error {
	filter := cloudflare.Filter{
		Expression: "ip.src eq " + ip,
		Action:     "allow",
	}

	_, err := api.CreateFilter(zoneID, filter)
	if err != nil {
		return err
	}

	return nil
}

// 修改白名单 IP
func updateIPInFilter(api *cloudflare.API, zoneID, host, oldIP, newIP string) error {
	filters, err := api.ListFilters(zoneID)
	if err != nil {
		return err
	}

	for _, filter := range filters {
		if filter.Expression == "ip.src eq "+oldIP && filter.Action == "allow" {
			filter.Expression = "ip.src eq " + newIP

			err := api.UpdateFilter(zoneID, filter.ID, filter)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}

// 删除白名单 IP
func deleteIPFromFilter(api *cloudflare.API, zoneID, host, ip string) error {
	filters, err := api.ListFilters(zoneID)
	if err != nil {
		return err
	}

	for _, filter := range filters {
		if filter.Expression == "ip.src eq "+ip && filter.Action == "allow" {
			err := api.DeleteFilter(zoneID, filter.ID)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}

func (c *CloudflareAPI) UpdateIPRules(ip string) error {
	filters := cloudflare.Filters{
		Match: "ip",
		Value: ip,
	}

	rules := []cloudflare.Filter{
		{
			Expression: filters,
			Action:     "allow",
			Comment:    c.Comment,
		},
	}

	err := c.API.UpdateFirewallRules(c.ZoneID, rules)
	if err != nil {
		return err
	}

	return nil
}

func (c *CloudflareAPI) RemoveIPFromRules(ip string) error {
	filters := cloudflare.Filters{
		Match: "ip",
		Value: ip,
	}

	rules := []cloudflare.Filter{
		{
			Expression: filters,
			Action:     "allow",
			Comment:    c.Comment,
		},
	}

	err := c.API.DeleteFirewallRules(c.ZoneID, rules)
	if err != nil {
		return err
	}

	return nil
}
