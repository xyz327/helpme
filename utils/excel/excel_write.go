package excel

import (
	"context"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

type HeaderInfo[T any] struct {
	Name        string
	ValueMapper func(T) interface{}
}
type Config struct {
	CellWidth float64
}
type option func(config *Config)

var WithCellWidth = func(width float64) option {
	return func(config *Config) {
		config.CellWidth = width
	}
}

func WriteToXlsx[T any](ctx context.Context, xlsx *excelize.File, sheetName string, headers []*HeaderInfo[T], rowsData []T, opts ...option) {
	if idx := xlsx.GetSheetIndex(sheetName); idx == 0 {
		// sheet 不存在
		_ = xlsx.NewSheet(sheetName)
	}
	if sheetName != "Sheet1" {
		// 设置默认的 sheet
		if idx := xlsx.GetSheetIndex("Sheet1"); idx != 0 {
			// 设置默认的 sheet
			xlsx.DeleteSheet("Sheet1")
		}
	}
	c := &Config{CellWidth: float64(15)}
	rowHeight := float64(20)
	// write header
	xlsx.SetColWidth(sheetName, excelize.ToAlphaString(0), excelize.ToAlphaString(len(headers)), c.CellWidth)
	xlsx.SetRowHeight(sheetName, 1, rowHeight)
	for x, header := range headers {
		axis := excelize.ToAlphaString(x) + strconv.Itoa(1)
		xlsx.SetCellValue(sheetName, axis, header.Name)

	}

	// write data
	for y, d := range rowsData {
		row := y + 2 // 跳过第一行
		for x, r := range headers {
			axis := excelize.ToAlphaString(x) + strconv.Itoa(row)
			xlsx.SetCellValue(sheetName, axis, r.ValueMapper(d))
		}
		xlsx.SetRowHeight(sheetName, row, rowHeight)
	}
	return
}
