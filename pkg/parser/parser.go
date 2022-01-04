package parser

import (
	"errors"
	"fmt"
	"kp_data/pkg/check"
	"kp_data/pkg/setting"
	"kp_data/pkg/xlsx"
	"log"
	"sort"
	"strconv"
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

func GetData() ([]map[string]string, error) {
	lotMap := lotParser()

	fmt.Println(lotMap)
	panelMap, _ := panelParser()

	fmt.Println(panelMap)
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

	return nil, nil
}

func lotParser() map[string]string {
	const valueIndex = 0
	const keyIndex = 1

	result := make(map[string]string)
	wb := (*WorkBook)[lotWbIndex]

	for _, cellData := range wb {
		key := cellData[keyIndex]
		value := cellData[valueIndex]

		if len(key) != 0 && len(value) != 0 {
			result[key] = value
		}

	}

	return result
}

func panelParser() (map[string][]map[string]interface{}, error) {

	const keyIndex = 0
	const valueIndex = 1

	result := make(map[string][]map[string]interface{})
	var chipGroup []map[string]string

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

		chipGroup = append(chipGroup, tmpMap)

		// 4 chips in a group
		if (index % 4) == 0 {

			var corner []string
			var sortedChipGroup []map[string]string
			var compareChip string

			// sorting the chip by corner and check if the chip ID is the same .
			for _, c := range chipGroup {
				//check key existed
				if err := check.MapKeyExist("Corner", c); err != nil {
					return nil, err
				}

				if err := check.MapKeyExist("Chip", c); err != nil {
					return nil, err
				}

				if compareChip == "" {
					compareChip = c["Chip"]
				}

				if c["Chip"] != compareChip {
					return nil, errors.New("the chip ids in the chip group are not the same")
				}

				corner = append(corner, c["Corner"])
			}
			sort.Strings(corner)

			for _, sortedID := range corner {
				for _, c := range chipGroup {
					id := c["Corner"]
					if id == sortedID {
						sortedChipGroup = append(sortedChipGroup, c)
					}
				}
			}

			//update some calculate value
			updateChipGroup, err := calcGroupValue(sortedChipGroup)
			if err != nil {
				return nil, err
			}

			result[chipGroup[0]["Chip"]] = updateChipGroup

			// reset parameter
			chipGroup = nil
			compareChip = ""
		}
	}
	return result, nil
}

func calcGroupValue(chipGroup []map[string]string) ([]map[string]interface{}, error) {

	var x_r, y_r float32
	var err error

	lotData := lotParser()
	x_r, err = mapValueToFloat("X_Variable Raito", lotData)
	y_r, err = mapValueToFloat("X_Variable Raito", lotData)

	if err != nil {
		return nil, err
	}

	var updateGroup []map[string]interface{}

	for _, chipData := range chipGroup {

		var xMeasure, yMeasure, xDesign, yDesign float32
		var xVar, yVar, xDraw, yDraw float32

		updateChip := make(map[string]interface{})
		xMeasure, err = mapValueToFloat("X_Measure", chipData)
		yMeasure, err = mapValueToFloat("Y_Measure", chipData)
		xDesign, err = mapValueToFloat("X_Design", chipData)
		yDesign, err = mapValueToFloat("Y_Design", chipData)

		if err != nil {
			return nil, err
		}

		xVar = xMeasure - (xDesign / 1000)
		yVar = yMeasure - (yDesign / 1000)
		xDraw = xDesign - (xVar * x_r)
		yDraw = yDesign - (yVar * y_r)

		for key, value := range chipData {
			switch key {
			case "X_Var":
				updateChip[key] = xVar
			case "Y_Var":
				updateChip[key] = yVar
			case "X_Draw":
				updateChip[key] = xDraw
			case "Y_Draw":
				updateChip[key] = yDraw
			default:
				updateChip[key] = value
			}
		}

		updateGroup = append(updateGroup, updateChip)

	}

	return updateGroup, nil
}

func mapValueToFloat(key string, mapData map[string]string) (float32, error) {

	if err := check.MapKeyExist(key, mapData); err != nil {
		return 0.0, err
	}
	value, err := strconv.ParseFloat(mapData[key], 32)
	if err != nil {
		return 0.0, err
	}

	float := float32(value)

	return float, nil

}
