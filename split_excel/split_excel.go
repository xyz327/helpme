package split_excel

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/yankeguo/rg"
	"helpme/utils/arrutils"
	"helpme/utils/excel"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func getCurrentExecDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(file)
}
func getCurrentExecutableDir() (path string) {
	defer func() {
		os.Create(path)
	}()
	return filepath.Join(os.TempDir(), "/helpme_tmp/", fmt.Sprintf("%v", time.Now().UnixNano()))
}
func Split(fileName string, _bytes []byte) (zipFile *bytes.Buffer, err error) {
	defer rg.Guard(&err)
	//teamPath := "./team.csv"
	//dataPath := "./settle.xlsx"
	pwd := getCurrentExecutableDir()
	defer func() {
		os.RemoveAll(pwd)
	}()
	fmt.Println(fmt.Sprintf("current exec dir: %v", pwd))
	fmap := map[string][]*Info{}

	xlsx, err := excelize.OpenReader(bytes.NewBuffer(_bytes))
	if err != nil {
		err = fmt.Errorf("read xlsx failed: %v", err)
		return
	}

	for _, sheetName := range xlsx.GetSheetMap() {
		rows, err := xlsx.Rows(sheetName)
		if err != nil {
			return nil, err
		}
		rows.Next()
		headers := rows.Columns()
		records := make([][]string, 0)
		for rows.Next() {
			columns := rows.Columns()
			if columns[0] == "" {
				continue
			}
			records = append(records, columns)
		}
		groupByRecord := arrutils.GroupBy(records, func(in []string) string {
			return strings.Trim(in[0], "")
		})
		fmap[fileName] = append(fmap[fileName], &Info{
			Headers:     headers,
			SheetName:   sheetName,
			GroupedData: groupByRecord,
		})
	}
	fcache := map[string]*excelize.File{}
	getFile := func(teamId string) *excelize.File {
		f, ok := fcache[teamId]
		if !ok {
			f = excelize.NewFile()
			fcache[teamId] = f
		}
		return f
	}
	ctx := context.Background()
	for fileName, infos := range fmap {
		for _, info := range infos {
			for teamId, data := range info.GroupedData {
				xlsx := getFile(teamId)
				sheetName := fmt.Sprintf("%v-%v", info.SheetName, fileName)
				excel.WriteToXlsx(ctx, xlsx, sheetName, getHeader(info.Headers), data)
			}
		}
	}
	//zipFile, err = os.Create("split.zip")
	//if err != nil {
	//	return nil, err
	//}
	//defer zipFile.Close()
	zipFile = bytes.NewBuffer([]byte{})
	//第二步，创建一个新的 *Writer 对象
	zipWriter := zip.NewWriter(zipFile)
	// 关闭 zip writer，将所有数据写入指向基础 zip 文件的数据流
	defer zipWriter.Close()
	//
	for teamId, file := range fcache {
		//第三步，向 zip 文件中添加第一个文件
		w, err := zipWriter.Create(fmt.Sprintf("%v.xlsx", teamId))
		if err != nil {
			return nil, err
		}
		//向 zip 文件中添加一个文件，返回一个待压缩的文件内容应写入的 Writer
		_, err = file.WriteTo(w)
		if err != nil {
			return nil, err
		}
	}

	return
}

type Info struct {
	Headers     []string
	SheetName   string
	GroupedData map[string][][]string
}

func (i *Info) String() string {
	return fmt.Sprintf("sheet: %v, data: %v", i.SheetName, i.GroupedData)
}

func getHeader(header []string) []*excel.HeaderInfo[[]string] {
	idx := 0
	return arrutils.Map(header, func(in string) *excel.HeaderInfo[[]string] {
		defer func() {
			idx++
		}()
		return func(idx int) *excel.HeaderInfo[[]string] {
			return &excel.HeaderInfo[[]string]{
				Name: in,
				ValueMapper: func(d []string) interface{} {
					if len(d) <= idx {
						return ""
					}
					return d[idx]
				},
			}
		}(idx)
	})
}
