package trtree

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"tr369-wss-client/pkg/api"
)

func FindKeyInMap(data map[string]interface{}, paths []string, fpath string) (interface{}, string, bool) {

	if len(paths) == 1 {
		if value, ok := data[paths[0]]; ok {
			if _, ok := value.(map[string]interface{}); ok {
				return nil, "", false
			}
			return value, fpath + paths[0], true
		} else if paths[0] == "" {
			return data, fpath, true
		} else if strings.Contains(paths[0], "[") {
			return nil, "", false
		}
	}

	path := paths[0]
	if value, ok := data[path]; ok {
		if dataMap, ok := value.(map[string]interface{}); ok {
			if val, fpath1, found := FindKeyInMap(dataMap, paths[1:], fpath+path+"."); found {
				return val, fpath1, found
			} else {
				return nil, "", false
			}
		} else {
			return nil, "", false
		}
	} else if strings.Contains(path, "[") {
		subPath := path[1 : len(path)-1]
		subPaths := strings.Split(subPath, "&&")
		cIndex := "0"

		for key := range data {
			total := 0
			if value, ok := data[key].(map[string]interface{}); ok {

				for _, tpath := range subPaths {
					tpaths := strings.Split(tpath, "==")
					if len(tpaths) != 2 {
						return nil, "", false
					}

					if value[tpaths[0]] == strings.ReplaceAll(tpaths[1], "\"", "") {
						total = total + 1
					}
				}
			}

			if total == len(subPaths) {
				cIndex = key
			}
		}

		if cIndex == "0" {
			return nil, "", false
		}

		dataMap, _ := data[cIndex].(map[string]interface{})
		if val, fpath1, found := FindKeyInMap(dataMap, paths[1:], fpath+cIndex+"."); found {
			return val, fpath1, found
		} else {
			return nil, "", false
		}

	}

	return nil, "", false
}

func changeToString(data interface{}) string {
	if value, ok := data.(string); ok {
		return value
	}

	if value, ok := data.(bool); ok {
		strValue := fmt.Sprintf("%t", value)
		return strValue
	}

	if value, ok := data.(int); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(int8); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(int16); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(int32); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(int64); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(uint); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(uint8); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(uint16); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(uint32); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(uint64); ok {
		strValue := fmt.Sprintf("%d", value)
		return strValue
	}

	if value, ok := data.(float32); ok {
		strValue := fmt.Sprintf("%d", int(value))
		return strValue
	}

	if value, ok := data.(float64); ok {
		strValue := fmt.Sprintf("%d", int(value))
		return strValue
	}

	return ""
}

func constructResultParams(data map[string]interface{}, tmap map[string]string, path string) map[string]string {

	for key := range data {
		if value, ok := data[key].(map[string]interface{}); ok {
			tmap = constructResultParams(value, tmap, path+key+".")
		} else {
			// 到达末尾
			tmap[path+key] = changeToString(data[key])
		}
	}

	return tmap

}

func handleGetRequest(data map[string]interface{}, path string) *api.GetResp_ResolvedPathResult {

	tmap := make(map[string]string)
	// 判断是否以.结尾
	paths := strings.Split(path, ".")

	if len(paths) < 2 {
		return nil
	}

	value, fpath, found := FindKeyInMap(data, paths, "")
	if !found {
		return nil
	}
	if paths[len(paths)-1] == "" {
		tmap = constructResultParams(value.(map[string]interface{}), tmap, "")

		return &api.GetResp_ResolvedPathResult{
			ResolvedPath: fpath,
			ResultParams: tmap,
		}
	}

	tmap[paths[len(paths)-1]] = changeToString(value)

	return &api.GetResp_ResolvedPathResult{
		ResolvedPath: fpath[0 : len(fpath)-len(paths[len(paths)-1])],
		ResultParams: tmap,
	}
}

func ConstructGetResp(data map[string]interface{}, paths []string) api.Response_GetResp {
	response := api.Response_GetResp{
		GetResp: &api.GetResp{
			ReqPathResults: []*api.GetResp_RequestedPathResult{},
		},
	}

	for _, path := range paths {
		requestedPathResult := api.GetResp_RequestedPathResult{
			RequestedPath:       path,
			ResolvedPathResults: []*api.GetResp_ResolvedPathResult{},
		}
		result := handleGetRequest(data, path)

		if result != nil {
			requestedPathResult.ResolvedPathResults = append(requestedPathResult.ResolvedPathResults, result)
		}
		response.GetResp.ReqPathResults = append(response.GetResp.ReqPathResults, &requestedPathResult)
	}

	return response
}

func SetValueInMap(data map[string]interface{}, paths []string, value string) {

	if len(paths) == 1 {
		data[paths[0]] = value
		return
	}

	path := paths[0]
	if oldValue, ok := data[path]; ok {
		if dataMap, ok := oldValue.(map[string]interface{}); ok {
			SetValueInMap(dataMap, paths[1:], value)
		} else {
			return
		}
	} else if strings.Contains(path, "[") {
		subPath := path[1 : len(path)-1]
		subPaths := strings.Split(subPath, "&&")
		cIndex := "0"

		for key := range data {
			total := 0
			if value, ok := data[key].(map[string]interface{}); ok {

				for _, tpath := range subPaths {
					tpaths := strings.Split(tpath, "==")
					if len(tpaths) != 2 {
						return
					}

					if value[tpaths[0]] == strings.ReplaceAll(tpaths[1], "\"", "") {
						total = total + 1
					}
				}
			}

			if total == len(subPaths) {
				cIndex = key
			}
		}

		if cIndex == "0" {
			return
		}

		dataMap, _ := data[cIndex].(map[string]interface{})
		SetValueInMap(dataMap, paths[1:], value)

	} else {
		tempMap := make(map[string]interface{})
		data[path] = tempMap
		SetValueInMap(tempMap, paths[1:], value)
	}

	return

}

func HandleSetRequest(data map[string]interface{}, path string, key string, value string) {

	paths := strings.Split(path+key, ".")
	SetValueInMap(data, paths, value)

}

func IsExistPath(data map[string]interface{}, path string) (isSuccess bool, tpath string) {
	if strings.Contains(path, "[") {
		paths := strings.Split(path, ".")
		_, fpath, found := FindKeyInMap(data, paths, "")
		if !found {
			return false, ""
		}
		return true, fpath
	}
	return true, path

}

func GetNewInstance(data map[string]interface{}, path string) (tpath string) {
	paths := strings.Split(path, ".")
	value, _, found := FindKeyInMap(data, paths, "")
	if !found {
		return path + "1."
	}

	var nums []int
	for key := range value.(map[string]interface{}) {
		if index, err := strconv.Atoi(key); err == nil {
			nums = append(nums, index)
		}
	}
	if len(nums) == 0 {
		return path + "1."
	}
	sort.Ints(nums)
	lastNumber := strconv.Itoa(nums[len(nums)-1] + 1)
	return path + lastNumber + "."
}

func deleteInstance(data map[string]interface{}, paths []string) {

	if len(paths) == 2 && paths[1] == "" {
		delete(data, paths[0])
		return
	}

	if len(paths) == 1 {
		delete(data, paths[0])
		return
	}

	path := paths[0]
	if value, ok := data[path]; ok {
		if dataMap, ok := value.(map[string]interface{}); ok {
			deleteInstance(dataMap, paths[1:])
		}
	}
	return

}

func HandleDeleteRequest(data map[string]interface{}, path string) (tpath string, isFound bool) {
	paths := strings.Split(path, ".")
	_, fpath, found := FindKeyInMap(data, paths, "")
	if !found {
		return "", false
	}

	paths = strings.Split(fpath, ".")
	deleteInstance(data, paths)

	return fpath, true

}

func CloneTrtree(trtree map[string]interface{}) map[string]interface{} {

	cloneTrtree := make(map[string]interface{})
	for k, v := range trtree {
		if value, ok := v.(map[string]interface{}); ok {
			cloneTrtree[k] = CloneTrtree(value)
		} else {
			cloneTrtree[k] = v
		}
	}
	return cloneTrtree
}
