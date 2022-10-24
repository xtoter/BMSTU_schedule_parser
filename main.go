package main

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Group struct {
	group_name string
	group_addr string
}
type Schedule struct {
	group_name  string
	cabinet     string
	num         int
	teacher     string
	itemname    string
	itemtype    string
	day         int
	isnumerator bool
}

func main() {
	links := get_group()
	var out []Schedule
	for i := 0; i < len(links); i++ {
		out = append(out, get_schedule(links[i])...)
		//	println(links[i].group_name)
	}

	for _, a := range out {
		if a.teacher == "Царев А. С." {
			println("group_name ", a.group_name,
				"cabinet", a.cabinet,
				"num", a.num,
				"teacher", a.teacher,
				"itemname", a.itemname,
				"itemtype ", a.itemtype,
				"day", a.day,
				"isnumerator", a.isnumerator)
		}

	}
	println(len(links))

}
func get_group() []Group {
	resp, _ := http.Get("https://lks.bmstu.ru/schedule/list")
	doc, _ := html.Parse(resp.Body)

	links := visit(nil, doc)
	resp.Body.Close()
	return links
}
func visit(links []Group, n *html.Node) []Group {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				t := strings.TrimSpace(n.FirstChild.Data)
				if len(a.Val) > 10 && a.Val[:10] == "/schedule/" && (t[:6] != "ИУК") && (t[:4] != "МК") && (t[:2] != "К") && (t[:4] != "ЛТ") {

					//println(t[:2])
					links = append(links, Group{t, a.Val})
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
func get_schedule(in Group) []Schedule {
	resp, _ := http.Get("https://lks.bmstu.ru/" + in.group_addr)
	doc, _ := html.Parse(resp.Body)
	day := -1
	schedul := parse_schedule(nil, doc, in.group_name, &day)
	//links := visit(nil, doc)
	resp.Body.Close()
	return schedul

}
func parse_schedule(schedule []Schedule, n *html.Node, namegroup string, day *int) []Schedule {
	//println(n.Data)

	if n.Type == html.ElementNode && n.Data == "table" {
		//println(n.Data)
		for _, a := range n.Attr {
			if n.Parent.Attr[0].Val == "col-md-6 hidden-sm hidden-md hidden-lg" &&
				a.Val == "table table-bordered text-center table-responsive" {
				*day++
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Data == "tbody" {

						numdate := 0
						for h := c.FirstChild; h != nil; h = h.NextSibling {

							if h.Data == "tr" {
								for j := h.FirstChild; j != nil; j = j.NextSibling {
									if j.Data == "td" {

										//println(j.Attr[0].Val)
										for _, a1 := range j.Attr {

											if a1.Val == "2" {
												numdate++
												count := 0

												var curitem Schedule
												curitem.group_name = namegroup
												curitem.num = numdate
												curitem.day = *day
												for j1 := j.FirstChild; j1 != nil; j1 = j1.NextSibling.NextSibling {

													if j1.FirstChild != nil {

														//println(j1.FirstChild.Data)
														count++
														switch count {
														case 1:
															curitem.itemtype = j1.FirstChild.Data
														case 2:
															curitem.itemname = j1.FirstChild.Data
														case 3:
															curitem.cabinet = j1.FirstChild.Data

														case 4:
															curitem.teacher = j1.FirstChild.Data
														}

													}
												}
												curitem.isnumerator = false
												schedule = append(schedule, curitem)

												curitem.isnumerator = true

												schedule = append(schedule, curitem)
											}
											if a1.Val == "text-success" {
												numdate++
												count := 0

												var curitem Schedule
												curitem.group_name = namegroup
												curitem.num = numdate
												curitem.day = *day

												for j1 := j.FirstChild; j1 != nil; j1 = j1.NextSibling {
													if j1.Data == "i" || j1.Data == "span" {
														count++
														if j1.FirstChild != nil {
															switch count {
															case 1:
																curitem.itemtype = j1.FirstChild.Data
															case 2:
																curitem.itemname = j1.FirstChild.Data
															case 3:
																curitem.cabinet = j1.FirstChild.Data

															case 4:
																curitem.teacher = j1.FirstChild.Data
															}
														}

													}

												}

												curitem.isnumerator = true
												if curitem.itemname != "" {
													schedule = append(schedule, curitem)
												}
											}
											if a1.Val == "text-info" {
												count := 0

												var curitem Schedule
												curitem.group_name = namegroup
												curitem.num = numdate
												curitem.day = *day

												for j1 := j.FirstChild; j1 != nil; j1 = j1.NextSibling {
													if j1.Data == "i" || j1.Data == "span" {
														count++
														if j1.FirstChild != nil {
															switch count {
															case 1:
																curitem.itemtype = j1.FirstChild.Data
															case 2:
																curitem.itemname = j1.FirstChild.Data
															case 3:
																curitem.cabinet = j1.FirstChild.Data

															case 4:
																curitem.teacher = j1.FirstChild.Data
															}
														}

													}

												}

												curitem.isnumerator = false
												if curitem.itemname != "" {
													schedule = append(schedule, curitem)
												}
											}
										}
									}

								}
							}
						}
					}
				}

			}
		}
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			schedule = parse_schedule(schedule, c, namegroup, day)
		}
	}
	return schedule
}
