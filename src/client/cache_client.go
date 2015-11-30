// cache_client
package main

import (
	"fmt"
	"gomemcache/memcache"
)

func main() {
	mc := memcache.New("127.0.0.1:9891")
	mc.Set(&memcache.Item{Key: "siteName", Value: []byte("sudops.com")})
	//mc.Delete("siteName")
	if v, err := mc.Get("siteName"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", string(v.Value))
	}
}
