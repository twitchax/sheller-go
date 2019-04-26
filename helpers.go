package sheller

import "strings"

func (m EnvironmentVariableMap) clone() EnvironmentVariableMap {
	newMap := make(EnvironmentVariableMap)

	for k, v := range m {
		newMap[k] = v
	}

	return newMap
}

func (m EnvironmentVariableMap) toCommandEnv() []string {
	env := make([]string, 0, len(m))

	for k, v := range m {
		env = append(env, k+"="+v)
	}

	return env
}

func (l ArgumentList) clone() ArgumentList {
	newList := make(ArgumentList, 0, len(l))

	for _, v := range l {
		newList = append(newList, v)
	}

	return newList
}

func (l ArgumentList) appendMany(args ...string) ArgumentList {
	newList := make(ArgumentList, 0, len(l)+len(args))

	for _, v := range l {
		newList = append(newList, v)
	}

	for _, v := range args {
		newList = append(newList, v)
	}

	return newList
}

func environmentToMap(values []string) EnvironmentVariableMap {
	newMap := make(EnvironmentVariableMap)

	for _, v := range values {
		pair := strings.Split(v, "=")
		newMap[pair[0]] = pair[1]
	}

	return newMap
}

func trim(s string) string {
	return strings.Trim(s, " \t\n")
}
