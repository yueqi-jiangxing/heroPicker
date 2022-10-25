package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//双色球号码：6红+1蓝
var number string

//双色球开奖期号
var issue string

//开始序号
var begin int64

func init() {
	flag.StringVar(&number, "n", "", "")
	flag.StringVar(&issue, "i", "", "")
	flag.Int64Var(&begin, "s", 1, "")
}

//双色球数据返回值
type SSQResult struct {
	Result []struct {
		Blue string `json:"blue"`
		Red  string `json:"red"`
	}
}

func main() {
	flag.Parse()
	if number == "" {
		if issue == "" {
			fmt.Println("请输入期号或开奖号码")
			return
		} else {
			//未指定开奖号码，但是指定了开奖期号，则通过双色球福彩官网获取该期的开奖号码
			client := http.Client{Timeout: 3 * time.Second}
			url := "http://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=ssq&issueCount=&dayStart=&dayEnd="
			url = url + "&issueStart=" + issue + "&issueEnd=" + issue
			request, err := http.NewRequest(http.MethodGet, url, nil)
			content := []byte("")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			request.Header.Add("Expect", "")
			request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
			request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,ko;q=0.8")
			request.Header.Add("Cache-Control", "no-cache")
			request.Header.Add("Connection", "Connection")
			request.Header.Add("Cookie", "HMF_CI=a9be2e67fe0ea6387b29b855aa0409b23eb4549edc1b64adc702c6fff187d432d82804cba609dcb8b75f68055eed6bfce81374f843f481aac5200bb04fcfda9a98; 21_vq=5")
			request.Header.Add("Host", "http://www.cwl.gov.cn")
			request.Header.Add("Pragma", "no-cache")
			request.Header.Add("Upgrade-Insecure-Requests", "1")
			request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
			response, err := client.Do(request)
			if err == nil {
				content, err = ioutil.ReadAll(response.Body)
				defer response.Body.Close()
			}
			SSQResult := SSQResult{}
			err = json.Unmarshal(content, &SSQResult)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			//组装开奖号码
			number = SSQResult.Result[0].Red + "," + SSQResult.Result[0].Blue
		}
	}
	numbers := strings.Split(number, ",")
	//组装英雄基础数据，该数据实际应根据具体业务指标中获取，此处仅作示例
	heroes := make(map[string]HeroesData)
	for i := 1; i <= 50; i++ {
		D := HeroesData{}
		D.Surplus = 10000
		D.IsEmpty = 0
		heroes[strconv.Itoa(i)] = D
	}
	//调用选将核心算法
	instances := Hero{}.New(begin)
	data := instances.Choose(heroes, numbers)
	HeroesList := GetHeroesData()
	//打印出选将后的结果
	for i := 1; i < 11; i++ {
		fmt.Printf("N%d:【%d - %s】%s", i, data[strconv.Itoa(i)], HeroesList[strconv.FormatInt(data[strconv.Itoa(i)], 10)], "\n")
	}
	return
}
