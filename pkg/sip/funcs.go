package sip

import (
	"regexp"
	"strings"
)

func parametersToString(pmap map[string]string) string {
	if pmap == nil {
		return ""
	}
	res := ""
	for k, v := range pmap {
		res = res + ";" + k
		if v == "" {
			continue
		} else {
			res = res + "=" + v
		}
	}
	return res
}

// ;a=b;c;d=f
var reg_parameter = regexp.MustCompile(`\s*(?:;?)\s*(.*?)\s*(?:;|$)`)

func stringToParameters(parameterStr string) map[string]string {
	res := make(map[string]string)
	if parameterStr == "" {
		return res
	}
	valueGroups := reg_parameter.FindAllStringSubmatch(parameterStr, -1)
	for _, valueGroup := range valueGroups {
		for _, pair := range valueGroup[1:] {
			pairs := strings.Split(pair, "=")
			if len(pairs) == 1 {
				res[strings.Trim(pairs[0], " ")] = ""
			} else {
				res[strings.Trim(pairs[0], " ")] = strings.Trim(pairs[1], " ")
			}
		}
	}
	return res
}

func valuesToString(values []string) string {
	res := ""
	for i, v := range values {
		if i > 0 {
			res = res + ","
		}
		res = res + v
	}
	return res
}
