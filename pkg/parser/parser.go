package parser

import (
	"fmt"
	"kp_data/pkg/setting"
	"kp_data/pkg/xlsx"
	"log"
)

var WorkBook *[][][]string

const lotWbIndex = 0
const panelWbIndex = 1

type KpLoc struct {
	SEQ string
	X1  string
	Y1  string
	X2  string
	Y2  string
	X3  string
	Y3  string
	X4  string
	Y4  string
}

func Setup() {
	wb, err := xlsx.FileToSlice(setting.SystemSetting.Source)
	if err != nil {
		log.Fatalf("parser.Setup, fail to read '%s': %v", setting.SystemSetting.Source, err)
	}
	WorkBook = &wb
}

//func GetPanelData() ([]map[string]string, error) {
//	lotMap := lotParser()
//	panelMap := panelParser()
//
//	fmt.Println(lotMap)
//
//	//var result []map[string]string
//	var kpData KpLoc
//	var locIndex = 0
//
//	groupChipID := 1
//	for index, data := range panelMap {
//
//		//
//		if err:= check.MapKeyExist("Chip",data);err!=nil{
//			return nil,err
//		}
//
//		chipID, err := strconv.Atoi(data["Chip"])
//		if err != nil {
//			return nil, errors.New("the xlxs chip id convert to an integer has an error")
//		}
//
//		if index == 0 && chipID!=1{
//			return nil, errors.New("The first chip id is not '1' ")
//		}
//
//
//		if chipID == groupChipID {
//			switch locIndex{
//
//			case 0:
//				kpData.X1 = ""
//				kpData.Y1 = ""
//				kpData.SEQ = "1"
//			case 1:
//				kpData.X1 = ""
//				kpData.Y1 = ""
//				kpData.SEQ = "2"
//			case 2:
//				kpData.X1 = ""
//				kpData.Y1 = ""
//				kpData.SEQ = "3"
//			case 3:
//				kpData.X1 = ""
//				kpData.Y1 = ""
//				kpData.SEQ = "4"
//			}
//
//			locIndex+=1
//
//		//reset all tmp parameter
//		}else{
//
//			fmt.Println("===============")
//			kpData.X1 = ""
//			kpData.Y1 = ""
//			kpData.SEQ = "1"
//
//			groupChipID = chipID
//			locIndex = 0
//			kpData =KpLoc{}
//		}
//
//		fmt.Println(kpData)
//	}
//
//	return nil, nil
//}

func lotParser() map[string]string {
	const valueIndex = 0
	const keyIndex = 1

	result := make(map[string]string)
	wb := (*WorkBook)[lotWbIndex]

	for _, cellData := range wb {
		key := cellData[keyIndex]
		value := cellData[valueIndex]
		result[key] = value
	}

	return result
}

func PanelParser() (map[string]interface{}, error) {

	const keyIndex = 0
	const valueIndex = 1

	result := make(map[string]interface{})
	var chipGroup []interface{}

	wb := (*WorkBook)[panelWbIndex]
	//get columns name
	var keyArray = wb[keyIndex]
	var valueArray = wb[valueIndex:]

	for index, cellData := range valueArray {
		//setting  index start was '1'
		index += 1
		tmpMap := make(map[string]string)

		for i, data := range cellData {
			tmpMap[keyArray[i]] = data
		}

		if (index % 4) == 0 {
			// handler last data
			if (index) == len(valueArray) {
				chipGroup = append(chipGroup, tmpMap)
			}

			for _, c := range chipGroup {
				fmt.Println(c.(map[string]string)["Corner"])
			}

			result[chipGroup[0].(map[string]string)["Chip"]] = chipGroup
			chipGroup = nil

		}
		chipGroup = append(chipGroup, tmpMap)

	}

	return nil, nil
}
