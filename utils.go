package alipanSdk

func formScopesStr(scope int) []string {
	var scopes []string
	for k, v := range scoreMap {
		if scope|k > 0 {
			scopes = append(scopes, v)
		}
	}
	return scopes
}

func formDriversStr(drive int) []string {
	var drivers []string
	for k, v := range driveMap {
		if drive|k > 0 {
			drivers = append(drivers, v)
		}
	}
	return drivers
}

func strArray2StrComma(strA []string) string {
	if len(strA) == 0 {
		return ""
	}
	result := strA[0]
	for i := 1; i < len(strA); i++ {
		result += "," + strA[i]
	}
	return result
}

func bool2Str(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
