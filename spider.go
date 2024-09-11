type SpiderIP struct {
	IP          string
	Description string
	IsCIDR      bool
}

func (s *SpiderIP) Match(ip string) bool {
	if !s.IsCIDR {
		return s.IP == ip
	}
	_, ipNet, err := net.ParseCIDR(s.IP)
	if err != nil {
		panic(err)
	}
	return ipNet.Contains(net.ParseIP(ip))

}

type ResultRecord struct {
	IP          string
	Count       int
	Description string
}

var items = []SpiderIP{
	{
		IP:          "60.28.22.0",
		Description: "百度蜘蛛",
	},
	{
		IP:          "60.172.229.61",
		Description: "这个ip段百度蜘蛛IP造访，准备抓取你东西，抓取网页的百度蜘蛛。",
	},
	{
		IP:          "61.129.45.72",
		Description: "这个ip段百度蜘蛛IP造访，准备抓取你东西，抓取网页的百度蜘蛛。",
	},
	{
		IP:          "61.135.162.0/23",
		Description: "这个ip段百度蜘蛛IP造访，准备抓取你东西，抓取网页的百度蜘蛛。",
		IsCIDR:      true,
	},
	{
		IP:          "61.135.164.0/22",
		Description: "这个ip段百度蜘蛛IP造访，准备抓取你东西，抓取网页的百度蜘蛛。",
		IsCIDR:      true,
	},
	{
		IP:          "61.135.168.0/23",
		Description: "这个ip段百度蜘蛛IP造访，准备抓取你东西，抓取网页的百度蜘蛛。",
		IsCIDR:      true,
	},
	{
		IP:          "61.135.186.0/23",
		Description: "百度图片爬虫",
		IsCIDR:      true,
	},
	{
		IP:          "61.135.188.0/23",
		Description: "百度图片爬虫",
		IsCIDR:      true,
	},
	{
		IP:          "61.135.190.0/24",
		Description: "百度图片爬虫",
		IsCIDR:      true,
	},
	{
		IP:          "111.206.198.0/24",
		Description: "百度渲染蜘蛛",
		IsCIDR:      true,
	},
	{
		IP:          "111.206.221.0/24",
		Description: "百度渲染蜘蛛",
		IsCIDR:      true,
	},
	{
		IP:          "116.179.32.12",
		Description: "与220开头的类似、新版百度蜘蛛，高权重段，一般抓取文章页。",
	},
	{
		IP:          "116.179.32.95",
		Description: "重点抓取内页，爬过的文章或首页基本24小时放出来。",
	},
	{
		IP:          "116.179.37.0/24",
		Description: "百度渲染蜘蛛,巡查合规，类同惩罚蜘蛛。",
		IsCIDR:      true,
	},
	{
		IP:          "119.188.14.13",
		Description: "百度蜘蛛，IP段位于济南市",
	},
	{
		IP:          "119.188.14.35",
		Description: "百度蜘蛛，IP段位于济南市",
	},
	{
		IP:          "121.14.89.0/24",
		Description: "这个ip段作为度过新站考察期，基本上是网站无排名。",
		IsCIDR:      true,
	},
	{
		IP:          "123.125.66.0/24",
		Description: "代表百度蜘蛛IP造访，准备抓取你东西",
		IsCIDR:      true,
	},
	{
		IP:          "123.125.68.0/24",
		Description: "这个蜘蛛经常来,别的来的少,表示网站可能要进入沙盒了，或被者降权。",
		IsCIDR:      true,
	},
	{
		IP:          "123.125.71.95",
		Description: "抓取内页收录的，权重较低，爬过此段的内页文章不会很快放出来，因不是原创或采集文章。",
	},
	{
		IP:          "123.125.71.97",
		Description: "抓取内页收录的，权重较低，爬过此段的内页文章不会很快放出来，因不是原创或采集文章。",
	},
	{
		IP:          "123.125.71.106",
		Description: "抓取内页收录的，权重较低，爬过此段的内页文章不会很快放出来，因不是原创或采集文章。",
	},
	{
		IP:          "123.125.71.117",
		Description: "抓取内页收录的，权重较低，爬过此段的内页文章不会很快放出来，因不是原创或采集文章。",
	},
	{
		IP:          "123.181.108.77",
		Description: "抓取内页收录的， 权重较低,爬过此段的内页文章不会很快放出来,因不是原创",
	},
	{
		IP:          "124.166.232.0/24",
		Description: "可能为新版新站专属百度蜘蛛，或低质量蜘蛛。",
		IsCIDR:      true,
	},
	{
		IP:          "125.90.88.0/24",
		Description: "主要造成成分，是新上线站较多，还有使用过站长工具，或SEO综合检测造成的",
		IsCIDR:      true,
	},
	{
		IP:          "159.226.50.0/24",
		Description: "百度其他蜘蛛，北京联通，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "162.105.207.0/24",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "180.76.5.87",
		Description: "百度其他蜘蛛，北京电信，功能未详！",
	},
	{
		IP:          "180.76.15.0/24",
		Description: "降权蜘蛛，有这个ip说明网站不会在收录了，一直到这个ip段消失",
		IsCIDR:      true,
	},
	{
		IP:          "180.149.133.0/24",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "183.91.40.144",
		Description: "这个ip段出现在新站或站点有不正常现象后",
	},
	{
		IP:          "202.108.249.0/24",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "202.108.250.0/24",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "203.208.60.0/24",
		Description: "此ip段为异常蜘蛛，通常由于网站服务器问题或其他违规行为会引起它来爬取",
		IsCIDR:      true,
	},
	{
		IP:          "210.72.225.0/24",
		Description: "此ip段为日常巡逻蜘蛛，只要网站没有问题，没有违规操作就行。",
		IsCIDR:      true,
	},
	{
		IP:          "218.30.118.102",
		Description: "每天这个IP 段只增不减很有可能进沙盒或K站",
	},
	{
		IP:          "220.181.7.0/24",
		Description: "代表百度蜘蛛IP造访，准备抓取你东西。",
		IsCIDR:      true,
	},
	{
		IP:          "220.181.19.0/24",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "220.181.32.0/22",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "220.181.36.0/23",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},

	{
		IP:          "220.181.38.0/24",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},

	{
		IP:          "220.181.108.0/24",
		Description: "百度其他蜘蛛，功能未详！",
		IsCIDR:      true,
	},
	{
		IP:          "220.181.158.107",
		Description: "百度其他蜘蛛，功能未详！",
	},
	{
		IP:          "220.181.68.0/24",
		Description: "每天这个IP 段只增不减很有可能进沙盒或K站降权。",
		IsCIDR:      true,
	},
	{
		IP:          "220.181.108.75",
		Description: "重点抓取更新文章的内页达到90%,8%的抓取首页,2%其他权重ip段,抓过的文章或首页基本24小时放出来。",
	},
	{
		IP:          "220.181.108.77",
		Description: "专用抓取首页IP权重段，一般返回代码是30400代表未更新。",
	},
	{
		IP:          "220.181.108.80",
		Description: "专用抓取首页IP权重段，一般返回代码是30400代表未更新。",
	},
	{
		IP:          "220.181.108.82",
		Description: "抓取tag页面。",
	},
	{
		IP:          "220.181.108.83",
		Description: "专用抓取首页IP权重段，一般返回代码是30400代表未更新。",
	},
	{
		IP:          "220.181.108.86",
		Description: "专用抓取首页IP权重段，一般返回代码是30400代表未更新。",
	},
	{
		IP:          "220.181.108.89",
		Description: "专用抓取首页IP权重段，一般返回代码是30400代表未更新。",
	},
	{
		IP:          "220.181.108.91",
		Description: "属于综合的。主要抓取首页和内页或者其它页面。属于权重IP段, 抓过的文章或首页基本24小时放出来",
	}, {
		IP:          "220.181.108.92",
		Description: "属于综合的。主要抓取首页和内页或者其它页面。属于权重IP段, 抓过的文章或首页基本24小时放出来",
	},
	{
		IP:          "220.181.108.93",
		Description: "重点抓取内页，爬过的文章或首页基本24小时放出来。",
	},
	{
		IP:          "220.181.108.94",
		Description: "专用抓取首页IP权重段，一般返回代码是304 0 0代表未更新。",
	},
	{
		IP:          "220.181.108.95",
		Description: "这个是百度抓取首页的专用IP，基本来说你的网站会天天隔夜快照。",
	},
	{
		IP:          "220.181.108.97",
		Description: "专用抓取首页IP权重段，一般返回代码是304 0 0代表未更新。",
	},
	{
		IP:          "220.181.108.115",
		Description: "重点抓取内页，爬过的文章或首页基本24小时放出来。",
	},
	{
		IP:          "220.181.108.119",
		Description: "专用抓取首页IP权重段，一般返回代码是304",
	},
	{
		IP:          "220.181.108.156",
		Description: "重点抓取内页，爬过的文章或首页基本24小时放出来。",
	},
	{
		IP:          "220.181.108.158",
		Description: "重点抓取内页，爬过的文章或首页基本24小时放出来。",
	},
	{
		IP:          "220.181.108.180",
		Description: "重点抓取内页，爬过的文章或首页基本24小时放出来。",
	},
	{
		IP:          "220.181.108.184",
		Description: "重点抓取内页，爬过的文章或首页基本24小时放出来。",
	},
}
func main(){
  	result := make(map[string]*ResultRecord)
	data, err := os.ReadFile("ip.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	ips := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	for _, ip := range ips {
		for _, item := range items {
			if item.Match(ip) {
				result[ip].Count += 1
				result[ip].IP = ip
				result[ip].Description = item.Description
			}
		}
	}
	for _, v := range result {
		fmt.Printf("IP:%s    次数:%d    描述:%s\n", v.IP, v.Count, v.Description)
	}

  
}
