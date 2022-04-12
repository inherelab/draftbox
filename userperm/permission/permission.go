package permission

import "path"

func match(urlPath string) {
	path.Match("*", urlPath)
}
