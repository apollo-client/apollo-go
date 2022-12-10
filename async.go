package apollo

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/apollo-client/apollo-go/log"
)

var (
	chstr = make(chan string) // namespace channel
	mns   sync.Map            // notifications sync map
)

// asyncApollo async get apollo config
func (c *Client) asyncApollo(namespace string, cb WatchCallback) error {
	// init sync map notifcation
	mns.Store(namespace, &Notifcation{NamespaceName: namespace, NotifcationID: -1})
	// get apollo config first
	status, apol, err := c.getConfigs(namespace, "")
	if err != nil || status != http.StatusOK {
		log.Errorf("watch namespace:%s, err:%v", namespace, err)
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}
	// if success, callback function
	if err = safeCallback(&apol, cb); err != nil {
		log.Errorf("watch namespace:%s, err:%v", namespace, err)
		return fmt.Errorf("watch namespace:%s, err:%v", namespace, err)
	}

	go func() {
		// listen namespace channel
		for nsp := range chstr {
			if !strings.EqualFold(nsp, namespace) {
				continue
			}
			ns, na, ne := c.getConfigs(namespace, apol.ReleaseKey)
			if ne != nil || ns != http.StatusOK {
				continue
			}
			apol = na
			_ = safeCallback(&apol, cb)
		}
	}()
	return nil
}

// asyncNotifications async get notifications
func (c *Client) asyncNotifications() {
	go func() {
		ticker := time.NewTicker(c.opts.WatchInterval)
		for range ticker.C {
			// get all notifications
			ns := make([]*Notifcation, 0)
			mns.Range(func(key, value interface{}) bool {
				n, ok := value.(*Notifcation)
				if !ok {
					log.Warnf("namespace notification err, namespace: %s", key)
					return false
				}
				ns = append(ns, n)
				return true
			})
			if len(ns) <= 0 {
				continue
			}

			// get remote notifications
			nns, nnn, nne := c.getNotifications(ns)
			if nne != nil || nns != http.StatusOK {
				continue
			}
			for _, n := range nnn {
				if n == nil {
					continue
				}
				// store notification and send namespace channel
				mns.Store(n.NamespaceName, n)
				chstr <- n.NamespaceName
			}
		}
	}()
}
