package main

import (
	"fmt"
	"kp_data/pkg/parser"
	"kp_data/pkg/setting"
)

func init() {
	setting.Setup()
	parser.Setup()
}

func main() {

	// open an existing file

	//parser.GetPanelData()

	fmt.Println(parser.PanelParser())

	//wb, err := xlsx.FileToSlice("KP_Data.xlsx")
	//if err != nil {
	//	panic(err)
	//}
	// wb now contains a reference to the workbook
	// show all the sheets in the workbook
	//fmt.Println("Sheets in this file:")
	//for i, sh := range wb.Sheets {
	//	fmt.Println(i, sh.Name)
	//}
	//
	//fmt.Println(wb.Sheets[0])
	//fmt.Println("----")

	//fmt.Println(wb[0])
}
