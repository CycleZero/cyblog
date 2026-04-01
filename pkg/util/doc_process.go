package util

import (
	"io"
	"leyline-doc-backend/log"
	"path"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/ZeroHawkeye/wordZero/pkg/markdown"
)

func DocxToText(file io.ReadCloser) string {

	memory, err := document.OpenFromMemory(file)

	//fmt.Println(p)
	e := markdown.NewExporter(nil)
	opt := markdown.DefaultExportOptions()
	opt.PreserveLineBreaks = true
	s, err := e.ExportToString(memory, opt)
	if err != nil {
		//fmt.Println(err)
		return ""
	}
	//fmt.Println(s)
	return s
}

func InArray[T comparable](obj T, arr []T) bool {
	for _, v := range arr {
		if obj == v {
			return true
		}
	}
	return false
}

func IsExcel(fileName string) bool {
	if InArray(path.Ext(fileName), []string{".xlsx"}) {
		return true
	}
	return false
}

func IsDocx(fileName string) bool {
	log.SugaredLogger().Debugf("文件名：%s", fileName)
	log.SugaredLogger().Debugf("文件后缀：%s", path.Ext(fileName))
	if InArray(path.Ext(fileName), []string{".docx", ".doc"}) {
		return true
	}
	return false
}
