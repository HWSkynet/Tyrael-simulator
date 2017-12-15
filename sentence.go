package main

//import "fmt"
import "math/rand"

//A:385629844458831873
//熏:385677254841204736
//花:379168019613745154
//猪:385751185036017665
//JR:379514394649821186
//斑:385363719003045888
//云:385366895496265729
//芋:385641722727628800
//m：

// 老司机最近天天看电影啊。跟哪个妹子去的
// 所以快发小仙女的皂片

var idle = []string{
	"<:xyx:389356458539614208>",
	"续命成功",
	"快女装",
	"[discord红包]我发了一个“专享红包”，请使用新版手机discord查收红包。",
}

var pic = []string{
	"女的？",
	"男的？",
}

var double7 = []string{
	"噫。七七",
	"七七嫉妒了",
	"萌七！么么艹",
	"萌七",
	"七七好厉害",
	"又黑萌七",
	"七七买买买",
}

var m = []string{
	"m啊。不要撩骚了",
	"m不怂了啊",
	"跳m",
	"m各种跳",
	"完了，m又装死了",
	"m不要怂了，快去吧",
	"m开始跳了啊，你的七七呢",
	"毕竟m",
	"m还是留给暖司机了",
}
var laoshifu = []string{
	"老师傅长的漂亮，发啥都对",
}
var bingo = []string{
	"Bingo不要这么不自信",
	"Bingo小仙女",
}
var azi = []string{
	"A子不怂了啊",
	"A子还是留给暖司机了",
}
var xunxun = []string{
	"还是熏熏",
}
var rika = []string{
	"六花你的小裙子呢",
}
var pig = []string{
	"长老沉迷发财，忘记了七七",
	"没想到你是这样的猪长老",
}
var jr = []string{
	"JR很懂的样子",
}
var banban = []string{
	"斑斑小心被请喝茶",
	"斑斑好厉害",
	"斑斑不是在看书嘛",
	"斑斑的学习就是水群",
}
var qianyunzi = []string{
	"浅云还没女装么",
	"浅云你的小裙子呢",
}
var xiangyu = []string{
	"香芋你的小裙子呢",
}
var xianii = []string{
	//"还有下限不会的嘛？穿上小裙子让六花教你",
	"我觉得下限说得对",
}

var name2str = map[string][]string{
	"A子":    azi,
	"熏熏":    xunxun,
	"花花":    rika,
	"猪哥":    pig,
	"JR":    jr,
	"斑斑":    banban,
	"浅云":    qianyunzi,
	"香芋":    xiangyu,
	"下限":    xianii,
	"七七":    double7,
	"m":     m,
	"老师傅":   laoshifu,
	"bingo": bingo,
}

var name2id = map[string]string{
	"A子": "385629844458831873",
	"熏熏": "385677254841204736",
	"花花": "379168019613745154",
	"猪哥": "385751185036017665",
	"JR": "379514394649821186",
	"斑斑": "385363719003045888",
	"浅云": "385366895496265729",
	"香芋": "385641722727628800",
	"下限": "377366407089881088",
}

var id2name = map[string]string{
	"385629844458831873": "A子",
	"385677254841204736": "熏熏",
	"379168019613745154": "花花",
	"385751185036017665": "猪哥",
	"379514394649821186": "JR",
	"385363719003045888": "斑斑",
	"385366895496265729": "浅云",
	"385641722727628800": "香芋",
	"377366407089881088": "下限",
}

func IsVip(id string) bool {
	if _, ok := id2name[id]; ok {
		return true
	} else {
		return false
	}
}

func GetIdle() string {
	return idle[rand.Intn(len(idle))]
}

func GetRandom(id string) string {
	arr := name2str[id2name[id]]
	return arr[rand.Intn(len(arr))]
}
