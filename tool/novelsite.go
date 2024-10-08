package tool

import (
	"regexp"
	"strings"
)

func TestHtmlOK(html string) bool {
	isOK := false
	if len(html) > 3 {
		lowCase := strings.ToLower(html)
		if strings.Contains(lowCase, "</html>") { // html 下载完毕
			isOK = true
		} else {
			if !strings.Contains(lowCase, "doctype") { // json
				isOK = true
			}
		}
	}
	return isOK
}

func GetTOCLast(html string) [][]string {
	if !TestHtmlOK(html) {
		return nil
	}
	lastCount := 50 // 取倒数50个链接

	// 链接列表 lks
	lks := regexp.MustCompile("(?smi)<a[^>]*?href=[\"|']([^\"']*?)[\"|'][^>]*?>([^<]*)<").FindAllStringSubmatch(html, -1)
	if nil == lks {
		return nil
	}

	lastIDX := len(lks) - 1
	firstIDX := lastIDX - lastCount
	firstURLLen := len(lks[firstIDX][1])
	firstURLLenB := firstURLLen + 1 // URL长度余量

	var nowLen, endIDX int
	bAdded := false
	for i, lk := range lks {
		if i < firstIDX {
			continue
		}
		nowLen = len(lk[1])
		if bAdded {
			if nowLen != firstURLLenB {
				break
			}
		} else {
			if nowLen != firstURLLen {
				if nowLen != firstURLLenB {
					break
				} else {
					bAdded = true
				}
			}
		}
		endIDX = i
	}

	// return lks[firstIDX : 1+endIDX]
	return filterIT(lks[firstIDX : 1+endIDX])
}

func GetTOC(html string) [][]string {
	if !TestHtmlOK(html) {
		return nil
	}
	var nowLen, nowCount int
	var isKeyExist bool

	// 链接列表 lks
	lks := regexp.MustCompile("(?smi)<a[^>]*?href=[\"|']([^\"']*?)[\"|'][^>]*?>([^<]*)<").FindAllStringSubmatch(html, -1)
	if nil == lks {
		return nil
	}

	// 获取 mapLenCount
	mapLenCount := make(map[int]int)
	for _, lk := range lks {
		nowLen = len(lk[1])
		nowCount, isKeyExist = mapLenCount[nowLen]
		if isKeyExist {
			mapLenCount[nowLen] = 1 + nowCount
		} else {
			nowCount = 1
			mapLenCount[nowLen] = 1
		}
		//		p(lk[1], lk[2])
	}

	// mapLenCount 获取 maxLen
	var maxLen, maxCount int = 0, 0
	for k, v := range mapLenCount {
		if v > maxCount {
			maxCount = v
			maxLen = k
		}
		//		p(k, "=", v)
	}
	//	p("maxLen =", maxLen, "  maxCount =", maxCount)
	var halfPos int = len(lks) / 2 // lks 中间点
	//	p(halfPos)
	// 中间往前找开始点
	startLen := maxLen - 1 // 边界点长度
	endLen := maxLen + 1   // 结束边界点长度
	startIDX := halfPos    // 开始索引，包含
	prevLen := len(lks[halfPos][1])
	for i := halfPos; i >= 0; i-- {
		nowLen = len(lks[i][1])
		if nowLen == prevLen {
			startIDX = i
			prevLen = nowLen
			continue
		} else if nowLen == startLen {
			startIDX = i
			prevLen = nowLen
			continue
		} else {
			break
		}
	}
	//	p("startIDX =", startIDX, "  startURL =", lks[startIDX][1], "  startName =", lks[startIDX][2])
	// 中间往后找结束点
	endIDX := halfPos // 结束索引，包含
	prevLen = len(lks[halfPos][1])
	allIDX := len(lks) - 1
	for i := halfPos; i <= allIDX; i++ {
		nowLen = len(lks[i][1])
		if nowLen == prevLen {
			endIDX = i
			prevLen = nowLen
			continue
		} else if nowLen == endLen {
			endIDX = i
			prevLen = nowLen
			continue
		} else {
			break
		}
	}
	//	p("endIDX =", endIDX, "  endURL =", lks[endIDX][1], "  endName =", lks[endIDX][2])

	// return lks[startIDX : 1+endIDX]
	return filterIT(lks[startIDX : 1+endIDX])
}

// 2020-11-11: 遍历links，比较url以/开头或http开头的数量与links的数量差异以确定url是什么类型的,方便删除多余链接，主要针对miaobige
func filterIT(links [][]string) [][]string {
	countLinks := len(links)
	countURLStartWithSlash := 0
	countURLStartWithHTTP := 0

	for _, l := range links {
		if strings.HasPrefix(l[1], "/") {
			countURLStartWithSlash += 1
			continue
		}
		if strings.HasPrefix(l[1], "http") {
			countURLStartWithHTTP += 1
			continue
		}
		// p(l[1], l[2])
	}

	if countURLStartWithSlash+countURLStartWithHTTP > 0 {
		countURLOther := countLinks - countURLStartWithSlash - countURLStartWithHTTP
		if 0 == countURLOther {
			if countURLStartWithHTTP > countURLStartWithSlash {
				for i, l := range links {
					if strings.HasPrefix(l[1], "/") {
						links = append(links[:i], links[i+1:]...) // 删除
					}
				}
			} else {
				for i, l := range links {
					if strings.HasPrefix(l[1], "http") {
						links = append(links[:i], links[i+1:]...) // 删除
					}
				}
			}
		}
	}
	return links
}

func GetContent(html string) string {
	bd := regexp.MustCompile("(?smi)<body[^>]*?>(.*)</body>").FindStringSubmatch(html)
	if nil != bd {
		html = bd[1]
	}
	// 替换无用标签，可根据需要自行添加
	html = regexp.MustCompile("(?smi)<script[^>]*?>.*?</script>").ReplaceAllString(html, "")
	html = regexp.MustCompile("(?smi)<style[^>]*?>.*?</style>").ReplaceAllString(html, "")
	html = regexp.MustCompile("(?smi)<a[^>]*?>").ReplaceAllString(html, "<a>")
	html = regexp.MustCompile("(?smi)<div[^>]*?>").ReplaceAllString(html, "<div>")
	// 下面这两个是必需的
	html = regexp.MustCompile("(?smi)[\r\n]*").ReplaceAllString(html, "")
	html = regexp.MustCompile("(?smi)</div>").ReplaceAllString(html, "</div>\n")

	// 获取最长的行maxLine
	lines := strings.Split(html, "\n")
	maxLine := ""
	nowLen, maxLen := 0, 0
	for _, line := range lines {
		nowLen = len(line)
		if nowLen > maxLen {
			maxLen = nowLen
			maxLine = line
		}
	}
	html = maxLine
	// 替换内容里面的html标签
	html = strings.Replace(html, "\t", "", -1)
	html = strings.Replace(html, "&nbsp;", " ", -1)
	html = strings.Replace(html, "</p>", "\n", -1)
	html = strings.Replace(html, "</P>", "\n", -1)
	html = strings.Replace(html, "<p>", "\n", -1)
	html = strings.Replace(html, "<P>", "\n", -1)
	html = strings.Replace(html, "<br>", "\n", -1)
	html = strings.Replace(html, "<br/>", "\n", -1)
	html = strings.Replace(html, "<br />", "\n", -1)
	html = strings.Replace(html, "<div>", "", -1)
	html = strings.Replace(html, "</div>", "", -1)
	html = regexp.MustCompile("(?i)<br[ /]*>").ReplaceAllString(html, "\n")
	html = regexp.MustCompile("(?smi)<a[^>]*?>.*?</a>").ReplaceAllString(html, "\n")
	html = regexp.MustCompile("(?mi)<[^>]*?>").ReplaceAllString(html, "") // 删除所有标签
	html = regexp.MustCompile("(?mi)^[ \t]*").ReplaceAllString(html, "") // 删除所有标签
	strings.TrimLeft(html, "\n ")
	html = strings.Replace(html, "\n\n", "\n", -1)
	html = strings.Replace(html, "　　", "", -1)
	return html
}

func IsQidanTOCURL_Desk8(iURL string) bool {
	return strings.Contains(iURL, "www.qidian.com/book/")
}

func IsQidanTOCURL_Touch8(iURL string) bool { // https://m.qidian.com/book/1038324115/catalog/
	return strings.Contains(iURL, "m.qidian.com/book/")
}
/*
func IsQidanTOCURL_Touch7_Ajax(iURL string) bool {
	// https://m.qidian.com/majax/book/category?bookId=1016319872
	return strings.Contains(iURL, "m.qidian.com/majax/book/category")
}
*/

func Qidian_GetTOC_Desk8(html string) [][]string {
	lks := regexp.MustCompile("(?smi)<li [^>]+>[^<]*<a [^>]+ href=\"(//www.qidian.com/chapter/[^\"]+)\"[^>]+>([^<]+)</a>(.*?)</li>").FindAllStringSubmatch(html, -1)
	if nil == lks {
		return nil
	}
	var olks [][]string // [] ["", pageurl, pagename]
	for _, lk := range lks {
		if strings.Contains(lk[3], "chapter-locked") {
			olks = append(olks, []string{"", "https:" + lk[1], "VIP: " + lk[2]})
			break
		} else {
			olks = append(olks, []string{"", "https:" + lk[1], lk[2]})
		}
	}
	return olks
}

func Qidian_GetTOC_Touch8(html string) [][]string {
	jsonStr := regexp.MustCompile("(?smi)\"application/json\">(.+)</script>").FindStringSubmatch(html)
	mID := regexp.MustCompile("(?i)\"bookId\":\"([0-9]+)\",").FindStringSubmatch(jsonStr[1])

	// {"uuid":1,"cN":"001巴克的早餐","uT":"2019-09-06  14:05","cnt":2089,"cU":"","id":491020997,"sS":1},
	// {"uuid":126,"cN":"第125章 死团子当活团子医","uT":"2019-10-18  12:12","cnt":3142,"cU":"","id":498495980,"sS":0},
	lks := regexp.MustCompile("(?mi)\"cN\":\"([^\"]+)\".*?\"id\":([0-9]+).*?\"sS\":([01])").FindAllStringSubmatch(jsonStr[1], -1)
	if nil == lks {
		return nil
	}
	var olks [][]string // [] ["", pageurl, pagename]
	for _, lk := range lks {
		if "1" == lk[3] {
			olks = append(olks, []string{"", Qidian_getContentURL_Touch8(lk[2], mID[1]), lk[1]})
		} else {
			olks = append(olks, []string{"", Qidian_getContentURL_Touch8(lk[2], mID[1]), "VIP: " + lk[1]})
			break
		}
	}
	return olks
}

/*
func Qidian_GetTOC_Touch7_Ajax(jsonStr string) [][]string {
	mID := regexp.MustCompile("(?i)\"bookId\":\"([0-9]+)\",").FindStringSubmatch(jsonStr)

	// {"uuid":1,"cN":"001巴克的早餐","uT":"2019-09-06  14:05","cnt":2089,"cU":"","id":491020997,"sS":1},
	// {"uuid":126,"cN":"第125章 死团子当活团子医","uT":"2019-10-18  12:12","cnt":3142,"cU":"","id":498495980,"sS":0},
	lks := regexp.MustCompile("(?mi)\"cN\":\"([^\"]+)\".*?\"id\":([0-9]+).*?\"sS\":([01])").FindAllStringSubmatch(jsonStr, -1)
	if nil == lks {
		return nil
	}
	var olks [][]string // [] ["", pageurl, pagename]
	for _, lk := range lks {
		if "1" == lk[3] {
			olks = append(olks, []string{"", Qidian_getContentURL_Touch7_Ajax(lk[2], mID[1]), lk[1]})
		} else {
			olks = append(olks, []string{"", Qidian_getContentURL_Touch7_Ajax(lk[2], mID[1]), "VIP: " + lk[1]})
			break
		}
	}
	return olks
}
*/
func IsQidanContentURL_Desk8(iURL string) bool {
	return strings.Contains(iURL, "www.qidian.com/chapter/")
}
func IsQidanContentURL_Touch8(iURL string) bool {
	return strings.Contains(iURL, "m.qidian.com/chapter/")
}
func Qidian_getContentURL_Touch8(pageID string, bookID string) string {
//	return "https://m.qidian.com/chapter/" + bookID + "/" + pageID + "/"
	return "/chapter/" + bookID + "/" + pageID + "/"
}
/*
func IsQidanContentURL_Touch7_Ajax(iURL string) bool {
	// https://m.qidian.com/majax/chapter/getChapterInfo?bookId=1015209014&chapterId=462636481
	return strings.Contains(iURL, "m.qidian.com/majax/chapter/getChapterInfo")
}

func Qidian_getContentURL_Touch7_Ajax(pageID string, bookID string) string {
	return "https://m.qidian.com/majax/chapter/getChapterInfo?bookId=" + bookID + "&chapterId=" + pageID
}
*/
func Qidian_GetContent_Desk8(html string) string {
	ctt := html
	fs := regexp.MustCompile("(?smi)<main[^>]+>(.+?)</main>").FindStringSubmatch(html)
	if nil != fs {
		ctt = fs[1]
	}
	if strings.HasPrefix(ctt, "<p>") {
		ctt = strings.TrimLeft(ctt, "<p>")
	}
	ctt = strings.Replace(ctt, "<p>", "\n", -1)
	ctt = strings.Replace(ctt, "　　", "", -1)
	ctt = strings.Replace(ctt, "</p>", "\n", -1)
	ctt = strings.Replace(ctt, "\n\n", "\n", -1)
	return ctt
}
/*
func Qidian_GetContent_Touch7_Ajax(jsonStr string) string {
	fs := regexp.MustCompile("(?smi)\"content\":\"([^\"]+)\",").FindStringSubmatch(jsonStr)
	if nil != fs {
		jsonStr = fs[1]
	}

	if strings.HasPrefix(jsonStr, "<p>　　") {
		jsonStr = strings.TrimLeft(jsonStr, "<p>　　")
	}
	jsonStr = strings.Replace(jsonStr, "<p>　　", "\n", -1)
	return jsonStr
}
*/

// var p = fmt.Println

// func main() {
// 	bb, _ := tool.ReadFile("T:/index.html")

// 	lk := GetTOCLast(string(bb))
// 	for _, l := range lk {
// 		p(l[1], l[2])
// 	}
// 	p("len =", len(lk))

// 	//	p( getContent(string(bb)) )

// 	//	var aaa int
// 	//	fmt.Scanf("%c",&aaa)
// }
