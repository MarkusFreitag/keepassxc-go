package keepassxc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/kevinburke/nacl/scalarmult"

	"github.com/MarkusFreitag/keepassxc-go/internal"
	"github.com/MarkusFreitag/keepassxc-go/pkg/keystore"
)

const APPLICATIONNAME = "keepassxc-go"

var (
	ErrUnspecifiedSocketPath = errors.New("unspecified socket path")
	ErrInvalidPeerKey        = errors.New("invalid peer key")
	ErrNotImplemented        = errors.New("not implemented yet")
)

type Client struct {
	socketPath      string
	applicationName string
	socket          *net.UnixConn

	privateKey nacl.Key
	publicKey  nacl.Key
	peerKey    nacl.Key

	id string

	associatedName string
	associatedKey  nacl.Key
}

type ClientOption func(*Client)

func WithApplicationName(name string) ClientOption {
	return func(client *Client) {
		client.applicationName = name
	}
}

func NewClient(socketPath, assoName string, assoKey nacl.Key, options ...ClientOption) *Client {
	if assoKey == nil || len(assoKey) == 0 {
		assoKey = nacl.NewKey()
	}

	client := &Client{
		socketPath:      socketPath,
		applicationName: APPLICATIONNAME,

		privateKey: nacl.NewKey(),

		associatedName: assoName,
		associatedKey:  assoKey,
	}
	client.publicKey = scalarmult.Base(client.privateKey)

	for _, option := range options {
		option(client)
	}

	client.id = client.applicationName + internal.NaclNonceToB64(nacl.NewNonce())

	return client
}

func (c *Client) encryptMessage(msg Message) ([]byte, error) {
	if len(c.peerKey) == 0 {
		return nil, ErrInvalidPeerKey
	}
	msgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return box.EasySeal(msgData, c.peerKey, c.privateKey), nil
}

func (c *Client) decryptResponse(encryptedMsg []byte) ([]byte, error) {
	if len(c.peerKey) == 0 {
		return nil, ErrInvalidPeerKey
	}
	return box.EasyOpen(encryptedMsg, c.peerKey, c.privateKey)
}

func (c *Client) sendMessage(msg Message, encrypted bool) (Response, error) {
	if encrypted {
		encryptedMsg, err := c.encryptMessage(msg)
		if err != nil {
			return nil, err
		}
		action := msg["action"]
		msg = Message{
			"action":  action,
			"message": base64.StdEncoding.EncodeToString(encryptedMsg[nacl.NonceSize:]),
			"nonce":   base64.StdEncoding.EncodeToString(encryptedMsg[:nacl.NonceSize]),
		}
	} else {
		msg["nonce"] = internal.NaclNonceToB64(nacl.NewNonce())
	}
	msg["clientID"] = c.id

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	_, err = c.socket.Write(data)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)
	count, err := c.socket.Read(buf)
	if err != nil {
		return nil, err
	}
	buf = buf[0:count]

	var resp Response
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return nil, err
	}

	if err, ok := resp["error"]; ok {
		return nil, fmt.Errorf("%v %s", resp["errorCode"], err.(string))
	}

	if encrypted {
		decoded, err := base64.StdEncoding.DecodeString(resp["nonce"].(string) + resp["message"].(string))
		if err != nil {
			return nil, err
		}
		decryptedMsg, err := c.decryptResponse(decoded)
		if err != nil {
			return nil, err
		}
		var msg map[string]interface{}
		err = json.Unmarshal(decryptedMsg, &msg)
		if err != nil {
			return nil, err
		}
		resp["message"] = msg
	}

	return resp, err
}

func (c *Client) GetAssociatedProfile() (string, string) {
	return c.associatedName, internal.NaclKeyToB64(c.associatedKey)
}

func (c *Client) Connect() error {
	if c.socketPath == "" {
		return ErrUnspecifiedSocketPath
	}

	var err error
	c.socket, err = net.DialUnix("unix", nil, &net.UnixAddr{Name: c.socketPath, Net: "unix"})
	return err
}

func (c *Client) Disconnect() error {
	if c.socket != nil {
		return c.socket.Close()
	}
	return nil
}

func (c *Client) ChangePublicKeys() error {
	resp, err := c.sendMessage(Message{
		"action":    "change-public-keys",
		"publicKey": internal.NaclKeyToB64(c.publicKey),
	}, false)
	if err != nil {
		return err
	}
	if peerKey, ok := resp["publicKey"]; ok {
		c.peerKey = internal.B64ToNaclKey(peerKey.(string))
		return nil
	}
	return errors.New("change-public-keys failed")
}

func (c *Client) GetDatabaseHash() (string, error) {
	resp, err := c.sendMessage(Message{
		"action": "get-databasehash",
	}, true)
	if err != nil {
		return "", err
	}
	if v, ok := resp["message"]; ok {
		if msg, ok := v.(map[string]interface{}); ok {
			if hash, ok := msg["hash"]; ok {
				return hash.(string), nil
			}
		}
	}
	return "", errors.New("get-databasehash failed")
}

func (c *Client) Associate() error {
	resp, err := c.sendMessage(Message{
		"action": "associate",
		"key":    internal.NaclKeyToB64(c.publicKey),
		"idKey":  internal.NaclKeyToB64(c.associatedKey),
	}, true)
	if err != nil {
		return err
	}
	if v, ok := resp["message"]; ok {
		if msg, ok := v.(map[string]interface{}); ok {
			if id, ok := msg["id"]; ok {
				c.associatedName = id.(string)
				return nil
			}
		}
	}
	return errors.New("associate failed")
}

func (c *Client) TestAssociate() error {
	_, err := c.sendMessage(Message{
		"action": "test-associate",
		"key":    internal.NaclKeyToB64(c.associatedKey),
		"id":     c.associatedName,
	}, true)
	return err
}

func (c *Client) GeneratePassword() (*Entry, error) {
	return nil, ErrNotImplemented
}

func (c *Client) GetLogins(url string) ([]*Entry, error) {
	msg := Message{
		"action": "get-logins",
		"url":    url,
		"keys": []map[string]string{
			{
				"id":  c.associatedName,
				"key": internal.NaclKeyToB64(c.associatedKey),
			},
		},
	}
	resp, err := c.sendMessage(msg, true)
	if err != nil {
		return nil, err
	}

	return resp.entries()
}

func (c *Client) SetLogin() error {
	return ErrNotImplemented
}

func (c *Client) LockDatabase() error {
	return ErrNotImplemented
}

func (c *Client) GetDatabaseGroups() ([]*DBGroup, error) {
	return nil, ErrNotImplemented
}

func (c *Client) CreateDatabaseGroup(name string) (string, string, error) {
	return "", "", ErrNotImplemented
}

func (c *Client) GetTOTP(uuid string) (string, error) {
	return "", ErrNotImplemented
}

func DefaultClient() (*Client, error) {
	store, err := keystore.Load()
	if err != nil {
		return nil, err
	}

	profile, err := store.DefaultProfile()
	if err != nil {
		return nil, err
	}

	socketPath, err := SocketPath()
	if err != nil {
		return nil, err
	}

	client := NewClient(
		socketPath,
		profile.Name,
		profile.NaclKey(),
	)

	if err := client.Connect(); err != nil {
		return nil, err
	}

	if err := client.ChangePublicKeys(); err != nil {
		return nil, err
	}

	if key := profile.NaclKey(); key == nil {
		if err := client.Associate(); err != nil {
			return nil, err
		}

		profile.Name, profile.Key = client.GetAssociatedProfile()

		if err := store.Add(profile); err != nil {
			return nil, err
		}

		if err := store.Save(); err != nil {
			return nil, err
		}
	} else {
		if err := client.TestAssociate(); err != nil {
			return nil, err
		}
	}

	return client, nil
}
