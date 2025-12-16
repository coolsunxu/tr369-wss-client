// Package trtree 提供 TR181 树结构操作功能
package trtree

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"tr369-wss-client/pkg/api"
)

// FindKeyInMap 在 map 中查找指定路径的键值
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
			}
			return nil, "", false
		}
		return nil, "", false
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
						total++
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
		}
		return nil, "", false
	}

	return nil, "", false
}

// changeToString 将任意类型转换为字符串
func changeToString(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%d", int(v))
	case float64:
		return fmt.Sprintf("%d", int(v))
	default:
		return ""
	}
}

// constructResultParams 构建结果参数
func constructResultParams(data map[string]interface{}, tmap map[string]string, path string) map[string]string {
	for key := range data {
		if value, ok := data[key].(map[string]interface{}); ok {
			tmap = constructResultParams(value, tmap, path+key+".")
		} else {
			tmap[path+key] = changeToString(data[key])
		}
	}
	return tmap
}

// handleGetRequest 处理 GET 请求
func handleGetRequest(data map[string]interface{}, path string) *api.GetResp_ResolvedPathResult {
	tmap := make(map[string]string)
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

// ConstructGetResp 构建 GET 响应
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

// SetValueInMap 在 map 中设置值
func SetValueInMap(data map[string]interface{}, paths []string, value string) {
	if len(paths) == 1 {
		data[paths[0]] = value
		return
	}

	path := paths[0]
	if oldValue, ok := data[path]; ok {
		if dataMap, ok := oldValue.(map[string]interface{}); ok {
			SetValueInMap(dataMap, paths[1:], value)
		}
		return
	} else if strings.Contains(path, "[") {
		subPath := path[1 : len(path)-1]
		subPaths := strings.Split(subPath, "&&")
		cIndex := "0"

		for key := range data {
			total := 0
			if v, ok := data[key].(map[string]interface{}); ok {
				for _, tpath := range subPaths {
					tpaths := strings.Split(tpath, "==")
					if len(tpaths) != 2 {
						return
					}
					if v[tpaths[0]] == strings.ReplaceAll(tpaths[1], "\"", "") {
						total++
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
}

// HandleSetRequest 处理 SET 请求
func HandleSetRequest(data map[string]interface{}, path, key, value string) {
	paths := strings.Split(path+key, ".")
	SetValueInMap(data, paths, value)
}

// IsExistPath 检查路径是否存在
func IsExistPath(data map[string]interface{}, path string) (bool, string) {
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

// GetNewInstance 获取新实例路径
func GetNewInstance(data map[string]interface{}, path string) string {
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

// deleteInstance 删除实例
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
}

// HandleDeleteRequest 处理 DELETE 请求
func HandleDeleteRequest(data map[string]interface{}, path string) (string, bool) {
	paths := strings.Split(path, ".")
	_, fpath, found := FindKeyInMap(data, paths, "")
	if !found {
		return "", false
	}

	paths = strings.Split(fpath, ".")
	deleteInstance(data, paths)

	return fpath, true
}

// CloneTrtree 克隆树结构
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
