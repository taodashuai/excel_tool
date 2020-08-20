package controller

import (
	"awesome2/util"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/tealeg/xlsx"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type IndexController struct {
	Ctx iris.Context
}

func (c *IndexController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "Index")
	b.Handle("POST", "/upload", "Upload")
	b.Handle("GET", "/excel/read", "ExcelRead")
}

func (c *IndexController) ExcelRead() {
	fileName := c.Ctx.URLParam("name")
	xlFile, err := xlsx.OpenFile(util.LocalPath() + fileName)
	if err != nil {
		fmt.Println(err)
		c.Ctx.JSON("error")
		return
	}
	//  按名字分组
	var nameResult = make(map[string][][]string, 0)
	var names = make([]string, 0)

	for _, v := range xlFile.Sheet {
		if len(v.Rows) == 0 {
			continue
		}

		for index, row := range v.Rows {
			var temp = make([]string, len(row.Cells))
			for indexI, cell := range row.Cells {
				temp[indexI] = cell.String()
			}
			if index == 0 {
				continue
			}
			if nameResult[row.Cells[1].String()] == nil {
				var n = make([][]string, 0)
				nameResult[row.Cells[1].String()] = n
			}
			nameResult[row.Cells[1].String()] = append(nameResult[row.Cells[1].String()], temp)
			for indexI, cell := range row.Cells {
				if indexI == 1 {
					names = append(names, cell.String())
				}
			}
		}
	}
	result1 := make([]map[string]interface{}, 0)
	for name, v0 := range nameResult {
		temp := make([][]string, 0)
		ids:=make([]string,0)
		for i := 0; i < len(v0); i++ {
			for j := i + 1; j < len(v0); j++ {
				// 如果出现相同的型号，就将这个型号丢进去
				if v0[i][2] == v0[j][2] {
					if !isContain(v0[i][0],ids) {
						ids=append(ids, v0[i][0])
						temp = append(temp, v0[i])
					}
					if !isContain(v0[j][0],ids) {
						ids=append(ids, v0[j][0])
						temp = append(temp, v0[j])
					}
				}
			}
		}
		if len(temp) == 0 {
			continue
		}
		m := make(map[string]interface{}, 0)
		m["name"] = name
		m["data"] = temp
		result1 = append(result1, m)
	}
	c.Ctx.JSON(result1)
}
func isContain(a string, all []string) bool {
	for _, v := range all {
		if a == v {
			return true
		}
	}
	return false
}
func (c *IndexController) Index() {
	c.Ctx.View("index.html")
}

// 上传
func (c *IndexController) Upload() {
	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		fmt.Println(err)
		c.Ctx.JSON("error")
		return
	}
	defer func() { err = file.Close() }()
	if info.Size > 100*1024*1024 {
		c.Ctx.JSON("error")
		return
	}
	fileType := strings.Split(info.Filename, ".")
	if fileType[len(fileType)-1] != "xls" && fileType[len(fileType)-1] != "csv" && fileType[len(fileType)-1] != "xlsx" {
		c.Ctx.JSON("error")
		return
	}

	unixTime := time.Now().UnixNano()
	fileName := strconv.FormatInt(unixTime, 10) + "." + fileType[len(fileType)-1]
	path := "web/upload"
	_, err = os.Stat(util.LocalPath() + path)
	if os.IsNotExist(err) {
		err := os.Mkdir(util.LocalPath()+path, os.ModePerm)
		fmt.Println(err)
		if err != nil {
			c.Ctx.JSON("error")
			return
		}
	}
	filePath := util.LocalPath() + "web/upload/" + fileName
	out, err := os.OpenFile(filePath,
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		c.Ctx.JSON("error")
		return
	}
	defer func() { err = out.Close() }()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
	}
	c.Ctx.JSON("/web/upload/" + fileName)
}
