package main

import (
	"strconv"
)

type Hero struct {
	Number   int64 //本期开始的序号
	Selected int64 //选取将领数量
	Total    int64 //将领总数量
}

type HeroesData struct {
	Surplus int64 //将领剩余可用数量
	IsEmpty int64 //轮空：0否,1是
}

/**
 * Description: 实例化选将SDK
 * @param ...int param1=>number, param2=>selected, param3=>total
 * @return struct Hero
 */
func (Hero) New(opt ...int64) *Hero {
	var Number int64 = 1
	if len(opt) > 0 {
		Number = opt[0]
	}
	var Selected int64 = 10
	if len(opt) > 1 {
		Selected = opt[1]
	}
	var Total int64 = 50
	if len(opt) > 2 {
		Total = opt[2]
	}
	hero := &Hero{
		Number:   Number,
		Selected: Selected,
		Total:    Total,
	}
	return hero
}

/**
 * Description:选取将领
 * @param Heroes 将领编号及对应的剩余数量+是否轮空
 * @param Numbers 双色球各编号
 * @return 返回选取的将领序号
 */
func (this *Hero) Choose(Heroes map[string]HeroesData, Numbers []string) map[string]int64 {
	//具体规则
	//1.构建双色球的编号及其位置A~G
	address := make(map[string]int64)
	for index, val := range Numbers {
		n := 65 + index
		var r rune = rune(n)
		value, _ := strconv.ParseInt(val, 10, 64)
		address[string(r)] = value
	}
	//2.构建抽奖参数T1~T10
	parameters := make(map[string]int64)
	for i := int64(0); i < this.Selected; i++ {
		n := 65 + i
		field := "T" + strconv.FormatInt(i+1, 10)
		if i < 6 {
			parameters[field] = address[string(rune(n))] + address[string(rune(71))]
		} else if i >= 6 && i < 9 {
			k := 59 + i
			var kv rune = rune(k)
			j := 76 - i
			var jv rune = rune(j)
			parameters[field] = address[string(kv)] + address[string(jv)]
		} else {
			parameters[field] = address[string(rune(65))] + address[string(rune(66))] + address[string(rune(67))]
		}
	}
	//3.构建将领序号N1~N10
	serialNumber := make(map[string]int64)
	start := this.Number
	for i := int64(0); i < this.Selected; i++ {
		index := strconv.FormatInt(i+1, 10)
		for j := int64(0); j < parameters["T"+index]; j++ {
			start++
			if start > this.Total {
				start = 1
				continue
			}
			jj := strconv.FormatInt(start, 10)
			if Heroes[jj].IsEmpty == 1 {
				start++
			}
		}
		for {
			if Heroes[strconv.FormatInt(start, 10)].IsEmpty == 1 {
				start++
			} else {
				break
			}
		}
		HeroesData := HeroesData{}
		HeroesData.Surplus = Heroes[strconv.FormatInt(start, 10)].Surplus
		HeroesData.IsEmpty = 1
		Heroes[strconv.FormatInt(start, 10)] = HeroesData
		serialNumber[index] = start
	}
	return serialNumber
}
