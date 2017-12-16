package main

//import "fmt"
import "math/rand"
import "strings"

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
	"[女装福袋]我发了一个“女装福袋”，请使用新版手机discord查收福袋。",
}

var pic = []string{
	"女的？",
	"男的？",
}

type keywords struct {
	defaults     []string
	keywords_map map[string][]string
}

var xianii keywords = keywords{
	defaults: []string{"我觉的下限说的很对"},
	keywords_map: map[string][]string{
		"下限不会": {"还有下限不会的嘛？穿上小裙子让六花教你", "还有下限不会的嘛？让六花穿上小裙子教你呀"},
		"不会":   {"还有下限不会的？让六花穿上小裙子教你", "还有下限不会的嘛？让六花穿上小裙子教你呀"},
	},
}

var double7 keywords = keywords{
	defaults: []string{
		"噫。七七",
		"七七嫉妒了",
		"萌七！么么艹",
		"萌七",
		"七七好厉害",
		"又黑萌七",
		"七七买买买",
	},
}

var m keywords = keywords{
	defaults: []string{
		"m啊。不要撩骚了",
		"m不怂了啊",
		"跳m",
		"m各种跳",
		"完了，m又装死了",
		"m不要怂了，快去吧",
		"m开始跳了啊，你的七七呢",
		"毕竟m",
		"m还是留给暖司机了",
	},
}
var laoshifu keywords = keywords{
	defaults: []string{
		"老师傅长的漂亮，发啥都对",
	},
}
var bingo keywords = keywords{
	defaults: []string{
		"Bingo不要这么不自信",
		"Bingo小仙女",
	},
}
var azi keywords = keywords{
	defaults: []string{
		"A子不怂了啊",
		"A子还是留给暖司机了",
	},
}
var xunxun keywords = keywords{
	defaults: []string{
		"还是熏熏",
	},
}
var rika keywords = keywords{
	defaults: []string{
		"六花你的小裙子呢",
	},
}
var pig keywords = keywords{
	defaults: []string{
		"长老沉迷发财，忘记了七七",
		"没想到你是这样的猪长老",
	},
}
var jr keywords = keywords{
	defaults: []string{
		"JR很懂的样子",
	},
}
var banban keywords = keywords{
	defaults: []string{
		"斑斑好厉害",
		"斑斑说的也没错啊",
	},
	keywords_map: map[string][]string{
		"女装": {"斑斑的学习就是水群", "斑斑爆照呀"},
		"裙":  {"斑斑不是在看书嘛"},
		"大佬": {"斑斑小心被请喝茶"},
		"飞":  {"斑斑小心被请喝茶"},
		"▽":  {"斑斑真可爱"},
		"☆":  {"斑斑太可爱了"},
	},
}
var qianyunzi keywords = keywords{
	defaults: []string{
		"噫。浅云",
	},
	keywords_map: map[string][]string{
		"女装": {"浅云还没女装么"},
		"裙":  {"浅云你的小裙子呢"},
	},
}
var xiangyu keywords = keywords{
	defaults: []string{
		"香芋快女装",
		"香芋怎么还没女装呢？(歪头",
	},
	keywords_map: map[string][]string{
		"裙": {"香芋你的小裙子呢"},
	},
}

var name2str = map[string]keywords{
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

func PicTalk() string {
	return pic[rand.Intn(len(pic))]
}

func IdleTalk() string {
	return idle[rand.Intn(len(idle))]
}

func Talk(id string, word string) string {
	if IsVip(id) {
		// 获取小可爱的反馈结构体
		cute := name2str[id2name[id]]
		// 找到关键词则使用对应关键词的回复
		for k, v := range cute.keywords_map {
			if strings.Contains(word, k) {
				return v[rand.Intn(len(v))]
			}
		}
		// 否则使用默认列表回复
		return cute.defaults[rand.Intn(len(cute.defaults))]
	} else {
		return ""
	}
}
