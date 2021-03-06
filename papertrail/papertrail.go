package papertrail

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
)

type Client struct {
	conn *tls.Conn
	pool *x509.CertPool
	url  string
}

func New(url string) (*Client, error) {
	c := &Client{url: url}
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(pemCerts); !ok {
		return nil, errors.New("failed to parse certs")
	}
	c.pool = pool
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Write(b []byte) (int, error) {
	s := fmt.Sprintf("%s\n", b)
	return c.writeAndRetry([]byte(s))
}

func (c *Client) connect() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	conn, err := tls.Dial("tcp", c.url, &tls.Config{RootCAs: c.pool})
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Client) writeAndRetry(b []byte) (int, error) {
	if c.conn != nil {
		if n, err := c.write(b); err == nil {
			return n, err
		}
	}
	if err := c.connect(); err != nil {
		return 0, err
	}
	return c.write(b)
}
