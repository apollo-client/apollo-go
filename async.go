package apollo

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	chstr = make(chan string)
	mns   sync.Map
)

func (c *Client) asyncApollo(namespace string, cb WatchCallback) error {
	mns.Store(namespace, &Notifcation{NamespaceName: namespace, NotifcationID: -1})
	status, apol, err := c.getConfigs(namespace, "")
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}
	if err = safeCallback(&apol, cb); err != nil {
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}

	go func() {
		for nsp := range chstr {
			fmt.Printf("namespace: %s\n", nsp)
			if !strings.EqualFold(nsp, namespace) {
				continue
			}
			ns, na, ne := c.getConfigs(namespace, apol.ReleaseKey)
			if ne != nil || ns != http.StatusOK {
				continue
			}
			fmt.Printf("namespace: %s, na: %+v, err: %+v\n", nsp, na, err)
			apol = na
			_ = safeCallback(&apol, cb)
		}
	}()
	return nil
}

func (c *Client) asyncNotifications() {
	go func() {
		ticker := time.NewTicker(c.opts.WatchInterval)
		for range ticker.C {
			ns := make([]*Notifcation, 0)
			mns.Range(func(key, value interface{}) bool {
				n, ok := value.(*Notifcation)
				if !ok {
					return false
				}
				ns = append(ns, n)
				return true
			})
			if len(ns) <= 0 {
				continue
			}

			nns, nnn, nne := c.getNotifications(ns)
			if nne != nil || nns != http.StatusOK {
				continue
			}
			for _, n := range nnn {
				if n == nil {
					continue
				}
				mns.Store(n.NamespaceName, n)
				chstr <- n.NamespaceName
			}
		}
	}()
}
