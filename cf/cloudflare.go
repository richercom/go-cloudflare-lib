package cf

import (
	"context"
	"fmt"
	"github.com/cloudflare-go"
	"log"
)

type CloudflareAPI struct {
	API        *cloudflare.API
	ZoneID     string
	Host       string
	IPListName string
	ZoneRe     *cloudflare.ResourceContainer
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
		ZoneRe:     cloudflare.ZoneIdentifier(zoneID),
	}, nil
}

// GetCustomHostname query Custom Hostnames info by name
func (c *CloudflareAPI) GetCustomHostName(ctx context.Context, hostname string) (*cloudflare.CustomHostname, error) {
	hostnames, _, err := c.API.CustomHostnames(ctx, c.ZoneID, 0, cloudflare.CustomHostname{Hostname: hostname})
	if err != nil {
		return nil, err
	}

	if len(hostnames) == 0 {
		return nil, nil // 未找到 Custom Hostname
	}

	return &hostnames[0], nil
}

// AddCustomHostname add Custom Hostname to enterprise host record
func (c *CloudflareAPI) AddCustomHostname(ctx context.Context, hostname, origin string) error {
	ret, err := c.API.CreateCustomHostname(ctx, c.ZoneID, cloudflare.CustomHostname{
		Hostname:           hostname,
		CustomOriginServer: origin,
	})
	if err != nil {
		return err
	}
	fmt.Println("[AddCustomHostname] set success", ret)

	return nil
}

// DeleteCustomHostname delete Custom Hostname to enterprise host record
func (c *CloudflareAPI) DeleteCustomHostname(ctx context.Context, id string) error {
	err := c.API.DeleteCustomHostname(ctx, c.ZoneID, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CloudflareAPI) AddDNSRecord(ctx context.Context, domain, recordType, content string) error {
	record := cloudflare.CreateDNSRecordParams{
		Type:    recordType,
		Name:    domain,
		Content: content,
		TTL:     1,
	}

	_, err := c.API.CreateDNSRecord(ctx, c.ZoneRe, record)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c *CloudflareAPI) DeleteDNSRecord(ctx context.Context, recordID string) error {
	err := c.API.DeleteDNSRecord(ctx, c.ZoneRe, recordID)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c *CloudflareAPI) GetDNSRecords(ctx context.Context, domain string) ([]cloudflare.DNSRecord, error) {
	records, _, err := c.API.ListDNSRecords(ctx, c.ZoneRe, cloudflare.ListDNSRecordsParams{Name: domain})
	if err != nil {
		return nil, err
	}

	return records, nil
}
