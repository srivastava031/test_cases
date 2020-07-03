package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

func slicesearch(p []interface{}, ks ...string) (retval interface{}, err error) {
	var ok bool
	i1, err := strconv.Atoi(ks[0])
	if err != nil {
		fmt.Println("cannot convert string to int")
	}
	retval = p[i1]
	if len(ks) <= 1 {
		return retval, nil
	} else if p, ok = retval.([]interface{}); !ok {

		if q, ok := retval.(map[string]interface{}); !ok {
			return retval, nil
		} else {
			return mapsearch(q, ks[1:]...)
		}
	} else {
		return slicesearch(p, ks[1:]...)
	}
	return nil, nil
}
func mapsearch(m map[string]interface{}, ks ...string) (retval interface{}, err error) {
	var ok bool
	if len(ks) == 0 {
		return nil, fmt.Errorf("needs at least one key")
	}
	if retval, ok = m[ks[0]]; !ok {
		return nil, fmt.Errorf("wrong key entered, keys: %v", ks)
	} else if len(ks) == 1 {
		return retval, nil
	} else if m, ok = retval.(map[string]interface{}); !ok {
		if p, ok := retval.([]interface{}); ok {
			return slicesearch(p, ks[1:]...)
		}
	} else {
		return mapsearch(m, ks[1:]...)
	}
	return
}
func main() {
	var m = make(map[string]interface{})
	data, _ := ioutil.ReadFile("/home/asus/go/src/nfcase/testcases/config/test_config/backup/tc_config.json")
	json.Unmarshal(data, &m)
	fmt.Println(m)
	nf, _ := mapsearch(m, "client", "A", "0", "port")
	fmt.Println(nf)
}
