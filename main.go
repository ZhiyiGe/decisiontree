package main

import (
	"fmt"
	"math"
	"sort"
)

func createDataSet() (titles []string, dataset [][]string) {
	titles = []string{"weather", "Temperature", "Humidity", "Wind", "PlayTennis"}
	dataset = [][]string{
		{"Sunny", "Hot", "High", "Weak", "no"},
		{"Sunny", "Hot", "High", "Strong", "no"},
		{"Overcast", "Hot", "High", "Weak", "yes"},
		{"Rainy", "Mild", "High", "Weak", "yes"},
		{"Rainy", "Cool", "Normal", "Weak", "yes"},
		{"Rainy", "Cool", "Normal", "Strong", "no"},
		{"Overcast", "Cool", "Normal", "Strong", "yes"},
		{"Sunny", "Mild", "High", "Weak", "no"},
		{"Sunny", "Cool", "Normal", "Weak", "yes"},
		{"Rainy", "Mild", "Normal", "Weak", "yes"},
		{"Sunny", "Mild", "Normal", "Strong", "yes"},
		{"Overcast", "Mild", "High", "Strong", "yes"},
		{"Overcast", "Hot", "Normal", "Weak", "yes"},
		{"Rainy", "Mild", "High", "Strong", "no"},
	}
	return titles, dataset
}

func calcShannonEnt(dataset [][]string) float64 {
	num := len(dataset)
	// 保存分类的字典
	labelCount := make(map[string]int, 20)
	for _, feavec := range dataset {
		// 取出最后一列数据
		label := feavec[len(feavec)-1]
		// 统计最后一列中元素的个数
		_, ok := labelCount[label]
		if !ok {
			labelCount[label] = 0
		}
		labelCount[label] += 1
	}

	// 计算信息熵
	entropy := 0.0
	for key := range labelCount {
		// 事件发生的概率
		p := float64(labelCount[key]) / float64(num)
		entropy += p * math.Log2(p)
	}
	return -entropy
}

func splitDataset(dataset [][]string, feature int, value string) [][]string {
	var retDataset [][]string
	var reduceFeavec []string
	for _, feavec := range dataset {
		if feavec[feature] == value {
			reduceFeavec=[]string{}
			for i := 0; i < len(feavec); i++ {
				if i != feature{
					reduceFeavec=append(reduceFeavec,feavec[i])
				}
			}
			retDataset = append(retDataset, reduceFeavec)
		}
	}
	return retDataset
}

func uniqueSlice(slice []string) []string {
	sort.Strings(slice)
	var j int
	for i := 0;i<len(slice)-1;i++ {
		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice = append(slice[:i+1], slice[j:]...)
	}
	return slice
}

func chooseBestFeature(dataset [][]string) int {
	// 获得特征的总数
	numFeature := len(dataset[0]) - 1
	// 获得初始信息熵
	initentropy := calcShannonEnt(dataset)
	// 最佳特征
	bestFeature := 0
	// 最大信息增益
	bestInfoGain := 0.0
	for i := 0; i <numFeature; i++ {
		// 获得第i个特征的一列
		var feaList []string
		for _, data := range dataset {
			feaList = append(feaList, data[i])
		}
		// 获得一个无序不重复元素集
		uniqueValues := uniqueSlice(feaList)
		newentropy:=0.0
		for _, value := range uniqueValues {
			// 根据特征i划分子集
			subDataset := splitDataset(dataset, i, value)
			// 分支集合在原集合中出现的概率
			p := float64(len(subDataset)) / float64(len(dataset))
			// 获得新的信息熵：分支集合出现的概率*分支集合的信息熵
			newentropy =newentropy+p * calcShannonEnt(subDataset)
			}
		// 获得信息增益：初始信息熵-新的信息熵
		infoGain := initentropy - newentropy
		// 获得最大信息增益
		if infoGain > bestInfoGain {
			bestInfoGain = infoGain
			bestFeature = i
		}
	}
	return bestFeature
}

func max(inputMap map[string]int) string {
	valuemax := 0
	var keymax string
	for key := range inputMap {
		if inputMap[key] > valuemax {
			keymax = key
			valuemax = inputMap[key]
		}
	}
	return keymax
}

func majorityCnt(classList []string) string {
	classCount := make(map[string]int, 20)
	// 定义一个字典
	// 统计数组中的元素及其出现的次数，并存在字典中
	for _, vote := range classList {
		_, ok := classCount[vote]
		if !ok {
			classCount[vote] = 0
		}
		classCount[vote] += 1
	}
	return max(classCount)
}

func createTree(dataSet [][]string, titles []string) interface{} {
	var classList []string
	for _, data := range dataSet {
		classList = append(classList, data[len(data)-1])
	}
	// 如果数组中所有元素相同，停止划分，并返回一个数组中元素的类别
	if len(uniqueSlice(classList)) == 1 {
		return classList[0]
	}
	// 长度为1，返回出现次数最多的类别
	if len(dataSet[0]) == 1 {
		return majorityCnt(classList)
	}
	// 按照信息增益最高选取分类特征属性
	bestFeat := chooseBestFeature(dataSet) // 返回分类的特征序号
	bestFeatTitle := titles[bestFeat]                          // 该特征的title
	mytree := make(map[string]interface{}, 5)                  // 构建树的字典
	// 从titles的list中删除该title
	var reducestitles []string
	for i := 0; i < len(titles); i++ {
		if i != bestFeat{
			reducestitles=append(reducestitles,titles[i])
		}
	}
	titles=reducestitles
	// 从dataset中获得bestFeat列得数据，得到一个新的list
	var featValues []string
	for _, example := range dataSet {
		featValues = append(featValues, example[bestFeat])
	}
	// 获得bestfeat列的不重复元素集
	uniqueVals := uniqueSlice(featValues)
	subtree := make(map[string]interface{}, 5)
	for _, value := range uniqueVals {
		subLables := titles[:] // 子集合
		// 构建数据的子集合，并进行递归
		subtree[value] = createTree(splitDataset(dataSet, bestFeat, value), subLables)
	}
	mytree[bestFeatTitle] = subtree
	return mytree

}

func getmapkey(inputmap map[string]interface{}) []string {
	var retclip []string
	for key := range inputmap {
		retclip = append(retclip, key)
	}
	return retclip
}
func getindex(inputslice []string, str string) int {
	ret := -1
	for i, s := range inputslice {
		if s == str {
			ret = i
		}
	}
	return ret
}
func classify(inputTree interface{}, titles,testVec []string) string {
	mytree:=inputTree.(map[string]interface{})
	firstStr := getmapkey(mytree)[0]     // 获取树的第一个特征属性
	featIndex := getindex(titles, firstStr) // 获取决策树第一层在featLables中的位置
	classLabel := "no"
	nextDict:=mytree[firstStr].(map[string]interface{})
	for key := range nextDict {
		if testVec[featIndex] == key{
			subtree,ok:=nextDict[key].(map[string]interface{})
			if ok{
				classLabel = classify(subtree, titles, testVec)
			}else{
				classLabel = nextDict[key].(string)
			}
		}
	}
	return classLabel
}
func main() {
	myTitle, myData := createDataSet()
	myTree := createTree(myData, myTitle)
	fmt.Println(myTree)
	classLabel1 := classify(myTree, myTitle, []string{"Sunny", "Hot", "High", "Weak"})
	classLabel2 := classify(myTree, myTitle, []string{"Overcast", "Hot", "High", "Weak"})
	fmt.Println(classLabel1)
	fmt.Println(classLabel2)
}
