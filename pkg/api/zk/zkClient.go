package zk

import (
	"atlas-service/pkg/api/log"
	"atlas-service/pkg/api/obs"
	"errors"
	"fmt"
	zk "github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
	"unsafe"
)

type Client struct {
	client       *zk.Conn
	waitIndex    uint64
	WatchHandler WatchHandler
}
type WatchHandler interface {
	NodeCreate(EventNodeCreated)
	NodeDeleted(EventNodeDeleted)
	NodeDataChanged(EventNodeDataChanged)
	NodeChildrenChanged(EventNodeChildrenChanged)
}

func New(machines []string) (*Client, error) {
	zkClient, _, err := zk.Connect(machines, time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{client: zkClient, waitIndex: 0}, nil
}
func NewWithHandler(machines []string, handler WatchHandler) (*Client, error) {
	zkClient, _, err := zk.Connect(machines, 20*time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{client: zkClient, waitIndex: 0, WatchHandler: handler}, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	resp, _, err := c.client.Get(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetString(path string) (s string, err error) {
	bytes, err := c.Get(path)
	s = *(*string)(unsafe.Pointer(&bytes))
	return
}

//深度搜索叶子
func nodeWalk(prefix string, c *Client, vars map[string]string) error {
	l, stat, err := c.client.Children(prefix)
	if err != nil {
		return err
	}

	if stat.NumChildren == 0 {
		b, _, err := c.client.Get(prefix)
		if err != nil {
			return err
		}
		vars[prefix] = string(b)

	} else {
		for _, key := range l {
			s := prefix + "/" + key
			_, stat, err := c.client.Exists(s)
			if err != nil {
				return err
			}
			if stat.NumChildren == 0 {
				b, _, err := c.client.Get(s)
				if err != nil {
					return err
				}
				vars[s] = string(b)
			} else {
				return nodeWalk(s, c, vars)
			}
		}
	}
	return nil
}

func (c *Client) GetValues(key string, keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, v := range keys {
		v = fmt.Sprintf("%s/%s", key, v)
		v = strings.Replace(v, "/*", "", -1)
		_, _, err := c.client.Exists(v)
		if err != nil {
			return vars, err
		}
		if v == "/" {
			v = ""
		}
		err = nodeWalk(v, c, vars)
		if err != nil {
			return vars, err
		}
	}
	return vars, nil
}

func (c *Client) createParents(key string) error {
	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)

	if key[0] != '/' {
		return errors.New("Invalid path")
	}

	payload := []byte("")
	pathString := ""
	pathNodes := strings.Split(key, "/")
	for i := 1; i < len(pathNodes); i++ {
		pathString += "/" + pathNodes[i]
		_, err := c.client.Create(pathString, payload, flags, acl)
		if err != nil && err != zk.ErrNodeExists && err != zk.ErrNoAuth {
			return err
		}
	}
	return nil
}
func (c *Client) CreateParentNode(node string) (err error) {
	exists, _, err := c.client.Exists(node)
	if !exists {
		_, err = c.client.Create(node, nil, 0, zk.WorldACL(zk.PermAll))
		//fmt.Print(err)
	}
	return
}
func (c *Client) CreateNode(node string, data string) (err error) {
	exists, _, err := c.client.Exists(node)
	if !exists {
		_, err = c.client.Create(node, []byte(data), 0, zk.WorldACL(zk.PermAll))
		fmt.Print(err)
	}
	return
}
func (c *Client) GetChildren(parentNode string) (list []string, err error) {
	list, _, err = c.client.Children(parentNode)
	return
}
func (c *Client) Set(path string, newData string) (err error) {
	_, sate, err := c.client.Get(path)
	if err != nil {
		return
	}
	_, err = c.client.Set(path, []byte(newData), sate.Version)
	if err != nil {
		fmt.Printf("数据修改失败: %v\n", err)
	}
	fmt.Println("数据修改成功")
	return
}

func (c *Client) CreateOrSetChildren(node string, data string) (err error) {
	exists, _, err := c.client.Exists(node)
	if !exists {
		err = c.CreateNode(node, data)
	} else {
		err = c.Set(node, data)
	}
	return
}

func (c *Client) DeleteNode(path string) (err error) {
	_, sate, _ := c.client.Get(path)
	err = c.client.Delete(path, sate.Version)
	if err != nil {
		fmt.Printf("删除失败: %v\n", err)
	}
	return
}
func (c *Client) EventRegistry() error {
	return obs.GetAsyncEventBus().Register(c.WatchHandler)
}
func (c *Client) Watch(path string) {
	go func() {
		for {
			_, _, e, err := c.client.GetW(path)
			if err != nil {
				log.Error(err.Error())
			}
			event := <-e
			var k interface{}
			base := BaseEvent{path, c}
			switch event.Type {
			case zk.EventNodeCreated:
				k = EventNodeCreated{base}
			case zk.EventNodeDeleted:
				k = EventNodeDeleted{base}
			case zk.EventNodeDataChanged:
				k = EventNodeDataChanged{base}
			case zk.EventNodeChildrenChanged:
				k = EventNodeChildrenChanged{base}
			}
			obs.GetAsyncEventBus().PostWithObs(k, c.WatchHandler)
		}
	}()

}

func (c *Client) Close() {
	c.client.Close()
}
