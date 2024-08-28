//黑名单
const blackes = [
    "baidu.com",
    "1688.com",
    "qcc.com",
    "gov.cn",
    "zhipin.com",
    "tianyancha.com",
    "sogou.com",
    "douyin.com",
    "sina.cn",
    "zhaopinku.cn",
    "google.com",
    "coovee.com",
    "job5156.com",
    "1637.com",
    "amap.com",
    "163.com",
    "cfi.cn",
    "cfi.net.cn",
    "amz123.com",
    "facebook.com",
    "sohu.com",
    "yingjiesheng.com",
];
async function* generateKeryword() {
    const cities = [
        "北京",
        "上海",
        "重庆",
        "广州",
        "深圳",
        "广东",
        "杭州",
        "浙江",
        "武汉",
        "湖北",
        "成都",
        "四川",
        "南京",
        "江苏",
        "沈阳",
        "辽宁",
        "长沙",
        "湖南",
        "石家庄",
        "河北",
        "郑州",
        "河南",
        "济南",
        "山东",
        "哈尔滨",
        "黑龙江",
        "长春",
        "吉林",
        "西安",
        "陕西",
        "福州",
        "福建",
        "合肥",
        "安徽",
        "南昌",
        "江西",
        "昆明",
        "云南",
        "呼和浩特",
        "内蒙古",
        "南宁",
        "广西",
        "太原",
        "山西",
        "乌鲁木齐",
        "新疆",
        "贵阳",
        "贵州",
        "兰州",
        "甘肃",
        "西宁",
        "青海",
        "海口",
        "海南",
        "银川",
        "宁夏",
        "拉萨",
        "西藏",
    ];
    for (const city of cities) {
        yield city + "光电有限公司";
    }

}
async function* generatePage() {
    for (let i = 0; i < 20; i++) {
        yield i;
    }
}
function sleep(time) {
    return new Promise(function (resolve) {
        setTimeout(resolve, time);
    });
}
for await (let keywrod of generateKeryword()) {
    for await (let page of generatePage()) {
        let start = page * 10;
        const url = `https://www.google.com/search?q=${keywrod}&start=${start}`;
        const resp = await fetch(url);
        if (resp.status == 429) {
            await sleep(60000)
        }
        const text = await resp.text();
        const parser = new DOMParser();
        const htmlDoc = parser.parseFromString(text, 'text/html');
        if(htmlDoc.querySelectorAll("a[data-ved]").length===0){
            break
        }
        htmlDoc.querySelectorAll("a[data-ved]").forEach(item => {
            const href = item.getAttribute("href");
            if (!href) {
                return
            }
            if (href.indexOf("http") !== 0) {
                return
            }
            const u = new URL(href);
            for (const black of blackes) {
                if (u.hostname.indexOf(black) !== -1) {
                    return
                }
            }
            console.log("%s,%s",u.protocol+"//"+u.hostname,  item.querySelector("h3")?.innerText);
        });
    }

}
