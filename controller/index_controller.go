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
	}
	var result = make([][]string, 0)
	for _, v := range xlFile.Sheet {
		if len(v.Rows) == 0 {
			continue
		}

		for _, row := range v.Rows {
			var temp = make([]string, len(row.Cells))
			for index, cell := range row.Cells {
				temp[index] = cell.String()
			}
			result = append(result, temp)
		}
	}
	c.Ctx.JSON(result)
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
