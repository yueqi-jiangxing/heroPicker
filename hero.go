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
	Empty   int64 //轮空：0否,1是
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
	Address := make(map[string]int64)
	for index, val := range Numbers {
		n := 65 + index
		var r rune = rune(n)
		value, _ := strconv.ParseInt(val, 10, 64)
		Address[string(r)] = value
	}
	//2.构建抽奖参数T1~T10
	Parameters := make(map[string]int64)
	for i := int64(0); i < this.Selected; i++ {
		n := 65 + i
		field := "T" + strconv.FormatInt(i+1, 10)
		if i < 6 {
			Parameters[field] = Address[string(rune(n))] + Address[string(rune(71))]
		} else if i >= 6 && i < 9 {
			k := 59 + i
			var kv rune = rune(k)
			j := 76 - i
			var jv rune = rune(j)
			Parameters[field] = Address[string(kv)] + Address[string(jv)]
		} else {
			Parameters[field] = Address[string(rune(65))] + Address[string(rune(66))] + Address[string(rune(67))]
		}
	}
	//3.构建将领序号N1~N10
	SerialNumber := make(map[string]int64)
	N := this.Number
	//每轮的计数
	S := N
	for i := int64(0); i < this.Selected; i++ {
		index := strconv.FormatInt(i+1, 10)
		for j := int64(0); j < Parameters["T"+index]; j++ {
			S++
			if S > this.Total {
				S = 1
				continue
			}
			jj := strconv.FormatInt(S, 10)
			if Heroes[jj].Empty == 1 {
				S++
			}
		}
		for {
			if Heroes[strconv.FormatInt(S, 10)].Empty == 1 {
				S++
			} else {
				break
			}
		}
		HeroesData := HeroesData{}
		HeroesData.Surplus = Heroes[strconv.FormatInt(S, 10)].Surplus
		HeroesData.Empty = 1
		Heroes[strconv.FormatInt(S, 10)] = HeroesData
		SerialNumber[index] = S
	}
	return SerialNumber
}
