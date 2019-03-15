package main

import (
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func main() {
	var arg string
	// -fオプション flag.Arg(0)だとファイル名が展開されてしまうようなので
	flag.StringVar(&arg, "f", "", "SearchPattern")
	// コマンドライン引数を解析
	flag.Parse()

	// カレントディレクトリの取得
	var curDir, _ = os.Getwd()
	curDir += "/"

	// 引数が取得できなければ、カレントディレクトリを使用
	if arg == "" {
		arg = curDir
	}

	// ディレクトリとファイルパターンに分割して格納
	var dirName, filePattern = path.Split(arg)

	// ディレクトリが無いならばカレントディレクトリを使用
	if dirName == "" {
		dirName = curDir
	}

	// 取得しようとしているパスがディレクトリかチェック
	var isDir, _ = IsDirectory(dirName + filePattern)

	// ディレクトリならば、そのディレクトリ配下のファイルを調べる。
	if isDir == true {
		dirName = dirName + filePattern
		filePattern = ""
	}

	fileInfos, _ := GetFileInfos(dirName, filePattern)

	//Excelのやーつ
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Println(err.Error())
	}

	iLoop := 1
	sheet.Cell(0, 0).Value = "ファイル名"
	sheet.Cell(0, 1).Value = "ファイル日付"
	sheet.Cell(0, 2).Value = "ファイルサイズ"
	for _, fileInfo := range fileInfos {
		sheet.Cell(iLoop, 0).Value = fileInfo.Name()
		sheet.Cell(iLoop, 1).Value = fmt.Sprint(fileInfo.ModTime())
		sheet.Cell(iLoop, 2).Value = fmt.Sprint(fileInfo.Size())
		//sheet.Cell(iLoop, 3).Value = "A1"
		//sheet.Cell(iLoop, 4).Value = "B1"
		//sheet.Cell(iLoop, 5).Value = "A2"
		iLoop++
	}

	//TODO:ファイル名のフォーマットを修正する
	err = file.Save(time.Now().Format("2006-01-01-0000000000") + ".xlsx")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetFileInfos(dirName string, filePattern string) (fileInfos []os.FileInfo, err error) {
	// ディレクトリ内のファイル情報の読み込み[]*os.FileInfoが返る。
	fileInfos, err = ioutil.ReadDir(dirName)
	// ディレクトリの読み込みに失敗したらエラーで終了
	if err != nil {
		fmt.Errorf("Directory cannot read %s\n", err)
		os.Exit(1)
	}
	//蛇足
	// ファイル情報を一つずつ表示する
	for _, fileInfo := range fileInfos {
		// *FileInfo型
		findName := fileInfo.Name()
		matched := true
		// lsのようなワイルドカード検索を行うため、path.Matchを呼び出す
		if filePattern != "" {
			matched, _ = path.Match(filePattern, findName)
		}
		// path.Matchでマッチした場合、ファイル名を表示
		if matched == true {
			fmt.Printf("一致%s\n", findName)
		}
		fmt.Println(fileInfo.Name())
		fmt.Println(fileInfo.IsDir())
		fmt.Println(fileInfo.Mode())
	}
	return
}

// 指定されたファイル名がディレクトリかどうか調べる
func IsDirectory(name string) (isDir bool, err error) {
	fInfo, err := os.Stat(name) // FileInfo型が返る。
	if err != nil {
		return false, err // もしエラーならエラー情報を返す
	}
	// ディレクトリかどうかチェック
	return fInfo.IsDir(), nil
}
