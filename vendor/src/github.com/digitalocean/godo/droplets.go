package godo

import (
	"encoding/json"
	"fmt"
)

const dropletBasePath = "v2/droplets"

// DropletsService is an interface for interfacing with the droplet
// endpoints of the Digital Ocean API
// See: https://developers.digitalocean.com/documentation/v2#droplets
type DropletsService interface {
	List(*ListOptions) ([]Droplet, *Response, error)
	Get(int) (*Droplet, *Response, error)
	Create(*DropletCreateRequest) (*Droplet, *Response, error)
	Delete(int) (*Response, error)
	Kernels(int, *ListOptions) ([]Kernel, *Response, error)
	Snapshots(int, *ListOptions) ([]Image, *Response, error)
	Backups(int, *ListOptions) ([]Image, *Response, error)
	Actions(int, *ListOptions) ([]Action, *Response, error)
	Neighbors(int) ([]Droplet, *Response, error)
}

// DropletsServiceOp handles communication with the droplet related methods of the
// DigitalOcean API.
type DropletsServiceOp struct {
	client *Client
}

var _ DropletsService = &DropletsServiceOp{}

// Droplet represents a DigitalOcean Droplet
type Droplet struct {
	ID          int       `json:"id,float64,omitempty"`
	Name        string    `json:"name,omitempty"`
	Memory      int       `json:"memory,omitempty"`
	Vcpus       int       `json:"vcpus,omitempty"`
	Disk        int       `json:"disk,omitempty"`
	Region      *Region   `json:"region,omitempty"`
	Image       *Image    `json:"image,omitempty"`
	Size        *Size     `json:"size,omitempty"`
	SizeSlug    string    `json:"size_slug,omitempty"`
	BackupIDs   []int     `json:"backup_ids,omitempty"`
	SnapshotIDs []int     `json:"snapshot_ids,omitempty"`
	Locked      bool      `json:"locked,bool,omitempty"`
	Status      string    `json:"status,omitempty"`
	Networks    *Networks `json:"networks,omitempty"`
	ActionIDs   []int     `json:"action_ids,omitempty"`
	Created     string    `json:"created_at,omitempty"`
}

// Kernel object
type Kernel struct {
	ID      int    `json:"id,float64,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// Convert Droplet to a string
func (d Droplet) String() string {
	return Stringify(d)
}

// DropletRoot represents a Droplet root
type dropletRoot struct {
	Droplet *Droplet `json:"droplet"`
	Links   *Links   `json:"links,omitempty"`
}

type dropletsRoot struct {
	Droplets []Droplet `json:"droplets"`
	Links    *Links    `json:"links"`
}

type kernelsRoot struct {
	Kernels []Kernel `json:"kernels,omitempty"`
	Links   *Links   `json:"links"`
}

type snapshotsRoot struct {
	Snapshots []Image `json:"snapshots,omitempty"`
	Links     *Links  `json:"links"`
}

type backupsRoot struct {
	Backups []Image `json:"backups,omitempty"`
	Links   *Links  `json:"links"`
}

// DropletCreateImage identifies an image for the create request. It prefers slug over ID.
type DropletCreateImage struct {
	ID   int
	Slug string
}

// MarshalJSON returns either the slug or id of the image. It returns the id
// if the slug is empty.
func (d DropletCreateImage) MarshalJSON() ([]byte, error) {
	if d.Slug != "" {
		return json.Marshal(d.Slug)
	}

	return json.Marshal(d.ID)
}

// DropletCreateSSHKey identifies a SSH Key for the create request. It prefers fingerprint over ID.
type DropletCreateSSHKey struct {
	ID          int
	Fingerprint string
}

// MarshalJSON returns either the fingerprint or id of the ssh key. It returns
// the id if the fingerprint is empty.
func (d DropletCreateSSHKey) MarshalJSON() ([]byte, error) {
	if d.Fingerprint != "" {
		return json.Marshal(d.Fingerprint)
	}

	return json.Marshal(d.ID)
}

// DropletCreateRequest represents a request to create a droplet.
type DropletCreateRequest struct {
	Name              string                `json:"name"`
	Region            string                `json:"region"`
	Size              string                `json:"size"`
	Image             DropletCreateImage    `json:"image"`
	SSHKeys           []DropletCreateSSHKey `json:"ssh_keys"`
	Backups           bool                  `json:"backups"`
	IPv6              bool                  `json:"ipv6"`
	PrivateNetworking bool                  `json:"private_networking"`
	UserData          string                `json:"user_data,omitempty"`
}

func (d DropletCreateRequest) String() string {
	return Stringify(d)
}

// Networks represents the droplet's networks
type Networks struct {
	V4 []NetworkV4 `json:"v4,omitempty"`
	V6 []NetworkV6 `json:"v6,omitempty"`
}

// NetworkV4 represents a DigitalOcean IPv4 Network
type NetworkV4 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (n NetworkV4) String() string {
	return Stringify(n)
}

// NetworkV6 represents a DigitalOcean IPv6 network.
type NetworkV6 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   int    `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (n NetworkV6) String() string {
	return Stringify(n)
}

// List all droplets
func (s *DropletsServiceOp) List(opt *ListOptions) ([]Droplet, *Response, error) {
	path := dropletBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dropletsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Droplets, resp, err
}

// Get individual droplet
func (s *DropletsServiceOp) Get(dropletID int) (*Droplet, *Response, error) {
	path := fmt.Sprintf("%s/%d", dropletBasePath, dropletID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dropletRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Droplet, resp, err
}

// Create droplet
func (s *DropletsServiceOp) Create(createRequest *DropletCreateRequest) (*Droplet, *Response, error) {
	path := dropletBasePath

	req, err := s.client.NewRequest("POST", path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(dropletRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Droplet, resp, err
}

// Delete droplet
func (s *DropletsServiceOp) Delete(dropletID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", dropletBasePath, dropletID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// Kernels lists kernels available for a droplet.
func (s *DropletsServiceOp) Kernels(dropletID int, opt *ListOptions) ([]Kernel, *Response, error) {
	path := fmt.Sprintf("%s/%d/kernels", dropletBasePath, dropletID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(kernelsRoot)
	resp, err := s.client.Do(req, root)
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Kernels, resp, err
}

// Actions lists the actions for a droplet.
func (s *DropletsServiceOp) Actions(dropletID int, opt *ListOptions) ([]Action, *Response, error) {
	path := fmt.Sprintf("%s/%d/actions", dropletBasePath, dropletID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Actions, resp, err
}

// Backups lists the backups for a droplet.
func (s *DropletsServiceOp) Backups(dropletID int, opt *ListOptions) ([]Image, *Response, error) {
	path := fmt.Sprintf("%s/%d/backups", dropletBasePath, dropletID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Backups, resp, err
}

// Snapshots lists the snapshots available for a droplet.
func (s *DropletsServiceOp) Snapshots(dropletID int, opt *ListOptions) ([]Image, *Response, error) {
	path := fmt.Sprintf("%s/%d/snapshots", dropletBasePath, dropletID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(snapshotsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Snapshots, resp, err
}

// Neighbors lists the neighbors for a droplet.
func (s *DropletsServiceOp) Neighbors(dropletID int) ([]Droplet, *Response, error) {
	path := fmt.Sprintf("%s/%d/neighbors", dropletBasePath, dropletID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dropletsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Droplets, resp, err
}

func (s *DropletsServiceOp) dropletActionStatus(uri string) (string, error) {
	action, _, err := s.client.DropletActions.GetByURI(uri)

	if err != nil {
		return "", err
	}

	return action.Status, nil
}
