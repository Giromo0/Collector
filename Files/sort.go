package main

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

var countrySymbols = map[string][]string{
    "USA": {
        "United States", "USA", "US", "America", "آمریکا", "США", "Америка", "美国", "امریکا",
        "Alabama", "آلاباما", "阿拉巴马州", "Alaska", "آلاسکا", "阿拉斯加州", "Arizona", "آریزونا", "亚利桑那州",
        "Arkansas", "آرکانزاس", "阿肯色州", "California", "کالیفرنیا", "加利福尼亚州", "加州", "Colorado", "کلرادو", "科罗拉多州",
        "Connecticut", "کنتیکت", "康涅狄格州", "Delaware", "دلاور", "特拉华州", "Florida", "فلوریدا", "佛罗里达州",
        "Georgia", "جورجیا", "佐治亚州", "Hawaii", "هاوایی", "夏威夷州", "Idaho", "آیداهو", "爱达荷州",
        "Illinois", "ایلینوی", "伊利诺伊州", "Indiana", "ایندیانا", "印第安纳州", "Iowa", "آیووا", "艾奥瓦州",
        "Kansas", "کانزاس", "堪萨斯州", "Kentucky", "کنتاکی", "肯塔基州", "Louisiana", "لوئیزیانا", "路易斯安那州",
        "Maine", "مین", "缅因州", "Maryland", "مریلند", "马里兰州", "Massachusetts", "ماساچوست", "马萨诸塞州",
        "Michigan", "میشیگان", "密歇根州", "Minnesota", "مینه‌سوتا", "明尼苏达州", "Mississippi", "میسیسیپی", "密西西比州",
        "Missouri", "میزوری", "密苏里州", "Montana", "مونتانا", "蒙大拿州", "Nebraska", "نبراسکا", "内布拉斯加州",
        "Nevada", "نوادا", "内华达州", "New Hampshire", "نیوهمپشر", "新罕布什尔州", "New Jersey", "نیوجرسی", "新泽西州",
        "New Mexico", "نیومکزیکو", "新墨西哥州", "New York", "نیویورک", "纽约州", "North Carolina", "کارولینای شمالی", "北卡罗来纳州",
        "North Dakota", "داکوتای شمالی", "北达科他州", "Ohio", "اوهایو", "俄亥俄州", "Oklahoma", "اکلاهما", "俄克拉荷马州",
        "Oregon", "اورگن", "俄勒冈州", "Pennsylvania", "پنسیلوانیا", "宾夕法尼亚州", "Rhode Island", "رود آیلند", "罗得岛州",
        "South Carolina", "کارولینای جنوبی", "南卡罗来纳州", "South Dakota", "داکوتای جنوبی", "南达科他州", "Tennessee", "تنسی", "田纳西州",
        "Texas", "تگزاس", "得克萨斯州", "德州", "Utah", "یوتا", "犹他州", "Vermont", "ورمانت", "佛蒙特州",
        "Virginia", "ویرجینیا", "弗吉尼亚州", "Washington", "واشنگتن", "华盛顿州", "West Virginia", "ویرجینیای غربی", "西弗吉尼亚州",
        "Wisconsin", "ویسکانسین", "威斯康星州", "Wyoming", "وایومینگ", "怀俄明州", "🇺🇸",
    },
    "Afghanistan": {"Afghanistan", "AF", "افغانستان", "Афганистан", "阿富汗", "🇦🇫"},
    "Albania": {"Albania", "AL", "آلبانی", "Албания", "阿尔巴尼亚", "🇦🇱"},
    "Algeria": {"Algeria", "DZ", "الجزایر", "Алжир", "阿尔及利亚", "🇩🇿"},
    "Andorra": {"Andorra", "AD", "آندورا", "Андорра", "安道尔", "🇦🇩"},
    "Angola": {"Angola", "AO", "آنگولا", "Ангола", "安哥拉", "🇦🇴"},
    "AntiguaAndBarbuda": {"Antigua and Barbuda", "AG", "آنتیگوا و باربودا", "Антигуа и Барбуда", "安提瓜和巴布达", "🇦🇬"},
    "Argentina": {"Argentina", "AR", "آرژانتین", "Аргентина", "阿根廷", "🇦🇷"},
    "Armenia": {"Armenia", "AM", "ارمنستان", "Армения", "亚美尼亚", "🇦🇲"},
    "Australia": {"Australia", "AU", "استرالیا", "Австралия", "澳大利亚", "🇦🇺"},
    "Austria": {"Austria", "AT", "اتریش", "Австрия", "奥地利", "🇦🇹"},
    "Azerbaijan": {"Azerbaijan", "AZ", "آذربایجان", "Азербайджан", "阿塞拜疆", "🇦🇿"},
    "Bahamas": {"Bahamas", "BS", "باهاما", "Багамы", "巴哈马", "🇧🇸"},
    "Bahrain": {"Bahrain", "BH", "بحرین", "Бахрейн", "巴林", "🇧🇭"},
    "Bangladesh": {"Bangladesh", "BD", "بنگلادش", "Бангладеш", "孟加拉国", "🇧🇩"},
    "Barbados": {"Barbados", "BB", "باربادوس", "Барбадос", "巴巴多斯", "🇧🇧"},
    "Belarus": {"Belarus", "BY", "بلاروس", "Беларусь", "白俄罗斯", "🇧🇾"},
    "Belgium": {"Belgium", "BE", "بلژیک", "Бельгия", "比利时", "🇧🇪"},
    "Belize": {"Belize", "BZ", "بلیز", "Белиз", "伯利兹", "🇧🇿"},
    "Benin": {"Benin", "BJ", "بنین", "Бенин", "贝宁", "🇧🇯"},
    "Bhutan": {"Bhutan", "BT", "بوتان", "Бутан", "不丹", "🇧🇹"},
    "Bolivia": {"Bolivia", "BO", "بولیوی", "Боливия", "玻利维亚", "🇧🇴"},
    "BosniaAndHerzegovina": {"Bosnia and Herzegovina", "BA", "بوسنی و هرزگوین", "Босния и Герцеговина", "波斯尼亚和黑塞哥维那", "🇧🇦"},
    "Botswana": {"Botswana", "BW", "بوتسوانا", "Ботсвана", "博茨瓦纳", "🇧🇼"},
    "Brazil": {"Brazil", "BR", "برزیل", "Бразилия", "巴西", "🇧🇷"},
    "Brunei": {"Brunei", "BN", "برونئی", "Бруней", "文莱", "🇧🇳"},
    "Bulgaria": {"Bulgaria", "BG", "بلغارستان", "Болгария", "保加利亚", "🇧🇬"},
    "BurkinaFaso": {"Burkina Faso", "BF", "بورکینافاسو", "Буркина-Фасо", "布基纳法索", "🇧🇫"},
    "Burundi": {"Burundi", "BI", "بوروندی", "Бурунди", "布隆迪", "🇧🇮"},
    "CaboVerde": {"Cabo Verde", "CV", "کیپ ورد", "Кабо-Верде", "佛得角", "🇨🇻"},
    "Cambodia": {"Cambodia", "KH", "کامبوج", "Камбоджа", "柬埔寨", "🇰🇭"},
    "Cameroon": {"Cameroon", "CM", "کامرون", "Камерун", "喀麦隆", "🇨🇲"},
    "Canada": {"Canada", "CA", "کانادا", "Канада", "加拿大", "🇨🇦"},
    "CentralAfricanRepublic": {"Central African Republic", "CF", "جمهوری آفریقای مرکزی", "Центральноафриканская Республика", "中非共和国", "🇨🇫"},
    "Chad": {"Chad", "TD", "چاد", "Чад", "乍得", "🇹🇩"},
    "Chile": {"Chile", "CL", "شیلی", "Чили", "智利", "🇨🇱"},
    "China": {"China", "CN", "چین", "Китай", "中国", "🇨🇳"},
    "Colombia": {"Colombia", "CO", "کلمبیا", "Колумбия", "哥伦比亚", "🇨🇴"},
    "Comoros": {"Comoros", "KM", "کومور", "Коморы", "科摩罗", "🇰🇲"},
    "CongoBrazzaville": {"Congo (Brazzaville)", "CG", "کنگو برازاویل", "Конго (Браззавиль)", "刚果（布）", "🇨🇬"},
    "CongoKinshasa": {"Congo (Kinshasa)", "CD", "کنگو کینشاسا", "Конго (Киншаса)", "刚果（金）", "🇨🇩"},
    "CostaRica": {"Costa Rica", "CR", "کاستاریکا", "Коста-Рика", "哥斯达黎加", "🇨🇷"},
    "Croatia": {"Croatia", "HR", "کرواسی", "Хорватия", "克罗地亚", "🇭🇷"},
    "Cuba": {"Cuba", "CU", "کوبا", "Куба", "古巴", "🇨🇺"},
    "Cyprus": {"Cyprus", "CY", "قبرس", "Кипр", "塞浦路斯", "🇨🇾"},
    "Czechia": {"Czechia", "Czech Republic", "CZ", "جمهوری چک", "Чехия", "捷克", "🇨🇿"},
    "Denmark": {"Denmark", "DK", "دانمارک", "Дания", "丹麦", "🇩🇰"},
    "Djibouti": {"Djibouti", "DJ", "جیبوتی", "Джибути", "吉布提", "🇩🇯"},
    "Dominica": {"Dominica", "DM", "دومینیکا", "Доминика", "多米尼克", "🇩🇲"},
    "DominicanRepublic": {"Dominican Republic", "DO", "جمهوری دومینیکن", "Доминиканская Республика", "多米尼加共和国", "🇩🇴"},
    "Ecuador": {"Ecuador", "EC", "اکوادور", "Эквадор", "厄瓜多尔", "🇪🇨"},
    "Egypt": {"Egypt", "EG", "مصر", "Египет", "埃及", "🇪🇬"},
    "ElSalvador": {"El Salvador", "SV", "السالوادور", "Сальвадор", "萨尔瓦多", "🇸🇻"},
    "EquatorialGuinea": {"Equatorial Guinea", "GQ", "گینه استوایی", "Экваториальная Гвинея", "赤道几内亚", "🇬🇶"},
    "Eritrea": {"Eritrea", "ER", "اریتره", "Эритрея", "厄立特里亚", "🇪🇷"},
    "Estonia": {"Estonia", "EE", "استونی", "Эстония", "爱沙尼亚", "🇪🇪"},
    "Eswatini": {"Eswatini", "SZ", "اسواتینی", "Эсватини", "斯威士兰", "🇸🇿"},
    "Ethiopia": {"Ethiopia", "ET", "اتیوپی", "Эфиопия", "埃塞俄比亚", "🇪🇹"},
    "Fiji": {"Fiji", "FJ", "فیجی", "Фиджи", "斐济", "🇫🇯"},
    "Finland": {"Finland", "FI", "فنلاند", "Финляндия", "芬兰", "🇫🇮"},
    "France": {"France", "FR", "فرانسه", "Франция", "法国", "🇫🇷"},
    "Gabon": {"Gabon", "GA", "گابن", "Габон", "加蓬", "🇬🇦"},
    "Gambia": {"Gambia", "GM", "گامبیا", "Гамбия", "冈比亚", "🇬🇲"},
    "Georgia": {"Georgia", "GE", "گرجستان", "Грузия", "格鲁吉亚", "🇬🇪"},
    "Germany": {"Germany", "DE", "Deutschland", "آلمان", "Германия", "德国", "🇩🇪"},
    "Ghana": {"Ghana", "GH", "غنا", "Гана", "加纳", "🇬🇭"},
    "Greece": {"Greece", "GR", "یونان", "Греция", "希腊", "🇬🇷"},
    "Grenada": {"Grenada", "GD", "گرنادا", "Гренада", "格林纳达", "🇬🇩"},
    "Guatemala": {"Guatemala", "GT", "گواتمالا", "Гватемала", "危地马拉", "🇬🇹"},
    "Guinea": {"Guinea", "GN", "گینه", "Гвинея", "几内亚", "🇬🇳"},
    "GuineaBissau": {"Guinea-Bissau", "GW", "گینه بیسائو", "Гвинея-Бисау", "几内亚比绍", "🇬🇼"},
    "Guyana": {"Guyana", "GY", "گویان", "Гайана", "圭亚那", "🇬🇾"},
    "Haiti": {"Haiti", "HT", "هائیتی", "Гаити", "海地", "🇭🇹"},
    "Honduras": {"Honduras", "HN", "هندوراس", "Гондурас", "洪都拉斯", "🇭🇳"},
    "Hungary": {"Hungary", "HU", "مجارستان", "Венгрия", "匈牙利", "🇭🇺"},
    "Iceland": {"Iceland", "IS", "ایسلند", "Исландия", "冰岛", "🇮🇸"},
    "India": {"India", "IN", "هند", "Индия", "印度", "🇮🇳"},
    "Indonesia": {"Indonesia", "ID", "اندونزی", "Индонезия", "印度尼西亚", "🇮🇩"},
    "Iran": {"Iran", "IR", "ایران", "Иран", "伊朗", "🇮🇷"},
    "Iraq": {"Iraq", "IQ", "عراق", "Ирак", "伊拉克", "🇮🇶"},
    "Ireland": {"Ireland", "IE", "ایرلند", "Ирландия", "爱尔兰", "🇮🇪"},
    "Israel": {"Israel", "IL", "اسرائیل", "Израиль", "以色列", "🇮🇱"},
    "Italy": {"Italy", "IT", "ایتالیا", "Италия", "意大利", "🇮🇹"},
    "Jamaica": {"Jamaica", "JM", "جامائیکا", "Ямайка", "牙买加", "🇯🇲"},
    "Japan": {"Japan", "JP", "ژاپن", "Япония", "日本", "🇯🇵"},
    "Jordan": {"Jordan", "JO", "اردن", "Иордания", "约旦", "🇯🇴"},
    "Kazakhstan": {"Kazakhstan", "KZ", "قزاقستان", "Казахстан", "哈萨克斯坦", "🇰🇿"},
    "Kenya": {"Kenya", "KE", "کنیا", "Кения", "肯尼亚", "🇰🇪"},
    "Kiribati": {"Kiribati", "KI", "کیریباتی", "Кирибати", "基里巴斯", "🇰🇮"},
    "Kuwait": {"Kuwait", "KW", "کویت", "Кувейт", "科威特", "🇰🇼"},
    "Kyrgyzstan": {"Kyrgyzstan", "KG", "قرقیزستان", "Кыргызстан", "吉尔吉斯斯坦", "🇰🇬"},
    "Laos": {"Laos", "LA", "لائوس", "Лаос", "老挝", "🇱🇦"},
    "Latvia": {"Latvia", "LV", "لتونی", "Латвия", "拉脱维亚", "🇱🇻"},
    "Lebanon": {"Lebanon", "LB", "لبنان", "Ливан", "黎巴嫩", "🇱🇧"},
    "Lesotho": {"Lesotho", "LS", "لسوتو", "Лесото", "莱索托", "🇱🇸"},
    "Liberia": {"Liberia", "LR", "لیبریا", "Либерия", "利比里亚", "🇱🇷"},
    "Libya": {"Libya", "LY", "لیبی", "Ливия", "利比亚", "🇱🇾"},
    "Liechtenstein": {"Liechtenstein", "LI", "لیختن‌اشتاین", "Лихтенштейн", "列支敦士登", "🇱🇮"},
    "Lithuania": {"Lithuania", "LT", "لیتوانی", "Литва", "立陶宛", "🇱🇹"},
    "Luxembourg": {"Luxembourg", "LU", "لوکزامبورگ", "Люксембург", "卢森堡", "🇱🇺"},
    "Madagascar": {"Madagascar", "MG", "ماداگاسکار", "Мадагаскар", "马达加斯加", "🇲🇬"},
    "Malawi": {"Malawi", "MW", "مالاوی", "Малави", "马拉维", "🇲🇼"},
    "Malaysia": {"Malaysia", "MY", "مالزی", "Малайзия", "马来西亚", "🇲🇾"},
    "Maldives": {"Maldives", "MV", "مالدیو", "Мальдивы", "马尔代夫", "🇲🇻"},
    "Mali": {"Mali", "ML", "مالی", "Мали", "马里", "🇲🇱"},
    "Malta": {"Malta", "MT", "مالت", "Мальта", "马耳他", "🇲🇹"},
    "MarshallIslands": {"Marshall Islands", "MH", "جزایر مارشال", "Маршалловы Острова", "马绍尔群岛", "🇲🇭"},
    "Mauritania": {"Mauritania", "MR", "موریتانی", "Мавритания", "毛里塔尼亚", "🇲🇷"},
    "Mauritius": {"Mauritius", "MU", "موریس", "Маврикий", "毛里求斯", "🇲🇺"},
    "Mexico": {"Mexico", "MX", "مکزیک", "Мексика", "墨西哥", "🇲🇽"},
    "Micronesia": {"Micronesia", "FM", "میکرونزی", "Микронезия", "密克罗尼西亚", "🇫🇲"},
    "Moldova": {"Moldova", "MD", "مولداوی", "Молдова", "摩尔多瓦", "🇲🇩"},
    "Monaco": {"Monaco", "MC", "موناکو", "Монако", "摩纳哥", "🇲🇨"},
    "Mongolia": {"Mongolia", "MN", "مغولستان", "Монголия", "蒙古", "🇲🇳"},
    "Montenegro": {"Montenegro", "ME", "مونته‌نگرو", "Черногория", "黑山", "🇲🇪"},
    "Morocco": {"Morocco", "MA", "مراکش", "Марокко", "摩洛哥", "🇲🇦"},
    "Mozambique": {"Mozambique", "MZ", "موزامبیک", "Мозамбик", "莫桑比克", "🇲🇿"},
    "Myanmar": {"Myanmar", "Burma", "MM", "میانمار", "Мьянма", "缅甸", "🇲🇲"},
    "Namibia": {"Namibia", "NA", "نامیبیا", "Намибия", "纳米比亚", "🇳🇦"},
    "Nauru": {"Nauru", "NR", "نائورو", "Науру", "瑙鲁", "🇳🇷"},
    "Nepal": {"Nepal", "NP", "نپال", "Непал", "尼泊尔", "🇳🇵"},
    "Netherlands": {"Netherlands", "NL", "Holland", "هلند", "Нидерланды", "Голландия", "荷兰", "🇳🇱"},
    "NewZealand": {"New Zealand", "NZ", "نیوزلند", "Новая Зеландия", "新西兰", "🇳🇿"},
    "Nicaragua": {"Nicaragua", "NI", "نیکاراگوئه", "Никарагуа", "尼加拉瓜", "🇳🇮"},
    "Niger": {"Niger", "NE", "نیجر", "Нигер", "尼日尔", "🇳🇪"},
    "Nigeria": {"Nigeria", "NG", "نیجریه", "Нигерия", "尼日利亚", "🇳🇬"},
    "NorthKorea": {"North Korea", "KP", "کره شمالی", "Северная Корея", "朝鲜", "🇰🇵"},
    "NorthMacedonia": {"North Macedonia", "MK", "مقدونیه شمالی", "Северная Македония", "北马其顿", "🇲🇰"},
    "Norway": {"Norway", "NO", "نروژ", "Норвегия", "挪威", "🇳🇴"},
    "Oman": {"Oman", "OM", "عمان", "Оман", "阿曼", "🇴🇲"},
    "Pakistan": {"Pakistan", "PK", "پاکستان", "Пакистан", "巴基斯坦", "🇵🇰"},
    "Palau": {"Palau", "PW", "پالائو", "Палау", "帕劳", "🇵🇼"},
    "Palestine": {"Palestine", "PS", "فلسطین", "Палестина", "巴勒斯坦", "🇵🇸"},
    "Panama": {"Panama", "PA", "پاناما", "Панама", "巴拿马", "🇵🇦"},
    "PapuaNewGuinea": {"Papua New Guinea", "PG", "پاپوآ گینه نو", "Папуа - Новая Гвинея", "巴布亚新几内亚", "🇵🇬"},
    "Paraguay": {"Paraguay", "PY", "پاراگوئه", "Парагвай", "巴拉圭", "🇵🇾"},
    "Peru": {"Peru", "PE", "پرو", "Перу", "秘鲁", "🇵🇪"},
    "Philippines": {"Philippines", "PH", "فیلیپین", "Филиппины", "菲律宾", "🇵🇭"},
    "Poland": {"Poland", "PL", "لهستان", "Польша", "波兰", "🇵🇱"},
    "Portugal": {"Portugal", "PT", "پرتغال", "Португалия", "葡萄牙", "🇵🇹"},
    "Qatar": {"Qatar", "QA", "قطر", "Катар", "卡塔尔", "🇶🇦"},
    "Romania": {"Romania", "RO", "رومانی", "Румыния", "罗马尼亚", "🇷🇴"},
    "Russia": {"Russia", "RU", "روسیه", "Россия", "俄罗斯", "🇷🇺"},
    "Rwanda": {"Rwanda", "RW", "رواندا", "Руанда", "卢旺达", "🇷🇼"},
    "SaintKittsAndNevis": {"Saint Kitts and Nevis", "KN", "سنت کیتس و نویس", "Сент-Китс и Невис", "圣基茨和尼维斯", "🇰🇳"},
    "SaintLucia": {"Saint Lucia", "LC", "سنت لوسیا", "Сент-Люсия", "圣卢西亚", "🇱🇨"},
    "SaintVincentAndTheGrenadines": {"Saint Vincent and the Grenadines", "VC", "سنت وینسنت و گرنادین‌ها", "Сент-Винсент и Гренадины", "圣文森特和格林纳丁斯", "🇻🇨"},
    "Samoa": {"Samoa", "WS", "ساموآ", "Самоа", "萨摩亚", "🇼🇸"},
    "SanMarino": {"San Marino", "SM", "سان مارینو", "Сан-Марино", "圣马力诺", "🇸🇲"},
    "SaoTomeAndPrincipe": {"Sao Tome and Principe", "ST", "سائوتومه و پرنسیپ", "Сан-Томе и Принсипи", "圣多美和普林西比", "🇸🇹"},
    "SaudiArabia": {"Saudi Arabia", "SA", "عربستان سعودی", "Саудовская Аравия", "沙特阿拉伯", "🇸🇦"},
    "Senegal": {"Senegal", "SN", "سنگال", "Сенегал", "塞内加尔", "🇸🇳"},
    "Serbia": {"Serbia", "RS", "صربستان", "Сербия", "塞尔维亚", "🇷🇸"},
    "Seychelles": {"Seychelles", "SC", "سیشل", "Сейшелы", "塞舌尔", "🇸🇨"},
    "SierraLeone": {"Sierra Leone", "SL", "سیرالئون", "Сьерра-Леоне", "塞拉利昂", "🇸🇱"},
    "Singapore": {"Singapore", "SG", "سنگاپور", "Сингапур", "新加坡", "🇸🇬"},
    "Slovakia": {"Slovakia", "SK", "اسلواکی", "Словакия", "斯洛伐克", "🇸🇰"},
    "Slovenia": {"Slovenia", "SI", "اسلوونی", "Словения", "斯洛文尼亚", "🇸🇮"},
    "SolomonIslands": {"Solomon Islands", "SB", "جزایر سلیمان", "Соломоновы Острова", "所罗门群岛", "🇸🇧"},
    "Somalia": {"Somalia", "SO", "سومالی", "Сомали", "索马里", "🇸🇴"},
    "SouthAfrica": {"South Africa", "ZA", "آفریقای جنوبی", "Южная Африка", "南非", "🇿🇦"},
    "SouthKorea": {"South Korea", "KR", "کره جنوبی", "Южная Корея", "韩国", "🇰🇷"},
    "SouthSudan": {"South Sudan", "SS", "سودان جنوبی", "Южный Судан", "南苏丹", "🇸🇸"},
    "Spain": {"Spain", "ES", "اسپانیا", "Испания", "西班牙", "🇪🇸"},
    "SriLanka": {"Sri Lanka", "LK", "سریلانکا", "Шри-Ланка", "斯里兰卡", "🇱🇰"},
    "Sudan": {"Sudan", "SD", "سودان", "Судан", "苏丹", "🇸🇩"},
    "Suriname": {"Suriname", "SR", "سورینام", "Суринам", "苏里南", "🇸🇷"},
    "Sweden": {"Sweden", "SE", "سوئد", "Швеция", "瑞典", "🇸🇪"},
    "Switzerland": {"Switzerland", "CH", "سوئیس", "Швейцария", "瑞士", "🇨🇭"},
    "Syria": {"Syria", "SY", "سوریه", "Сирия", "叙利亚", "🇸🇾"},
    "Taiwan": {"Taiwan", "TW", "تایوان", "Тайвань", "台湾", "🇹🇼"},
    "Tajikistan": {"Tajikistan", "TJ", "تاجیکستان", "Таджикистан", "塔吉克斯坦", "🇹🇯"},
    "Tanzania": {"Tanzania", "TZ", "تانزانیا", "Танзания", "坦桑尼亚", "🇹🇿"},
    "Thailand": {"Thailand", "TH", "تایلند", "Таиланд", "泰国", "🇹🇭"},
    "TimorLeste": {"Timor-Leste", "TL", "تیمور شرقی", "Восточный Тимор", "东帝汶", "🇹🇱"},
    "Togo": {"Togo", "TG", "توگو", "Того", "多哥", "🇹🇬"},
    "Tonga": {"Tonga", "TO", "تونگا", "Тонга", "汤加", "🇹🇴"},
    "TrinidadAndTobago": {"Trinidad and Tobago", "TT", "ترینیداد و توباگو", "Тринидад и Тобаго", "特立尼达和多巴哥", "🇹🇹"},
    "Tunisia": {"Tunisia", "TN", "تونس", "Тунис", "突尼斯", "🇹🇳"},
    "Turkey": {"Türkiye", "Turkey", "TR", "Turkiye", "ترکیه", "Турция", "土耳其", "🇹🇷"},
    "Turkmenistan": {"Turkmenistan", "TM", "ترکمنستان", "Туркменистан", "土库曼斯坦", "🇹🇲"},
    "Tuvalu": {"Tuvalu", "TV", "تووالو", "Тувалу", "图瓦卢", "🇹🇻"},
    "Uganda": {"Uganda", "UG", "اوگاندا", "Уганда", "乌干达", "🇺🇬"},
    "Ukraine": {"Ukraine", "UA", "اوکراین", "Украина", "乌克兰", "🇺🇦"},
    "UAE": {"United Arab Emirates", "AE", "UAE", "امارات", "ОАЭ", "阿联酋", "🇦🇪"},
    "UK": {"United Kingdom", "GB", "UK", "England", "Britain", "انگلیس", "بریتانیا", "Великобритания", "Англия", "英国", "🇬🇧"},
    "Uruguay": {"Uruguay", "UY", "اروگوئه", "Уругвай", "乌拉圭", "🇺🇾"},
    "Uzbekistan": {"Uzbekistan", "UZ", "ازبکستان", "Узбекистан", "乌兹别克斯坦", "🇺🇿"},
    "Vanuatu": {"Vanuatu", "VU", "وانواتو", "Вануату", "瓦努阿图", "🇻🇺"},
    "VaticanCity": {"Vatican City", "VA", "واتیکان", "Ватикан", "梵蒂冈", "🇻🇦"},
    "Venezuela": {"Venezuela", "VE", "ونزوئلا", "Венесуэла", "委内瑞拉", "🇻🇪"},
    "Vietnam": {"Vietnam", "VN", "ویتنام", "Вьетнам", "越南", "🇻🇳"},
    "Yemen": {"Yemen", "YE", "یمن", "Йемен", "也门", "🇾🇪"},
    "Zambia": {"Zambia", "ZM", "زامبیا", "Замбия", "赞比亚", "🇿🇲"},
    "Zimbabwe": {"Zimbabwe", "ZW", "زیمبابوه", "Зимбабве", "津巴布韦", "🇿🇼"},
}

// کلمات گمراه‌کننده
var ignoreWords = []string{
    "test", "free", "premium", "vip", "cloud", "server", "proxy", "vpn",
    "archive", "android", "ios", "config", "fast", "slow", "channel", "group",
}

// لیست پروتکل‌ها
var protocols = []string{
    "vmess://", "vless://", "trojan://", "ss://", "ssr://", "hy2://", "hysteria2://", "tuic://", "warp://", "wireguard://",
}

// شناسایی کشور با اولویت‌بندی
func identifyCountry(config string) string {
    configLower := strings.ToLower(config)
    for _, ignore := range ignoreWords {
        configLower = strings.ReplaceAll(configLower, ignore, "")
    }

    flagToCountry := make(map[string]string)
    for country, symbols := range countrySymbols {
        for _, symbol := range symbols {
            if strings.HasPrefix(symbol, "🇦") || strings.HasPrefix(symbol, "🇧") ||
                strings.HasPrefix(symbol, "🇨") || strings.HasPrefix(symbol, "🇩") ||
                strings.HasPrefix(symbol, "🇪") || strings.HasPrefix(symbol, "🇫") ||
                strings.HasPrefix(symbol, "🇬") || strings.HasPrefix(symbol, "🇭") ||
                strings.HasPrefix(symbol, "🇮") || strings.HasPrefix(symbol, "🇯") ||
                strings.HasPrefix(symbol, "🇰") || strings.HasPrefix(symbol, "🇱") ||
                strings.HasPrefix(symbol, "🇲") || strings.HasPrefix(symbol, "🇳") ||
                strings.HasPrefix(symbol, "🇴") || strings.HasPrefix(symbol, "🇵") ||
                strings.HasPrefix(symbol, "🇶") || strings.HasPrefix(symbol, "🇷") ||
                strings.HasPrefix(symbol, "🇸") || strings.HasPrefix(symbol, "🇹") ||
                strings.HasPrefix(symbol, "🇺") || strings.HasPrefix(symbol, "🇻") ||
                strings.HasPrefix(symbol, "🇼") || strings.HasPrefix(symbol, "🇽") ||
                strings.HasPrefix(symbol, "🇾") || strings.HasPrefix(symbol, "🇿") {
                flagToCountry[symbol] = country
            }
        }
    }

    if idx := strings.Index(config, "#"); idx != -1 {
        remark := config[idx+1:]
        for flag, country := range flagToCountry {
            if strings.Contains(remark, flag) {
                return country
            }
        }
    }

    if idx := strings.Index(config, "?"); idx != -1 {
        query := config[idx+1:]
        for flag, country := range flagToCountry {
            if strings.Contains(query, flag) {
                return country
            }
        }
    }

    if idx := strings.Index(config, "#"); idx != -1 {
        remark := strings.ToLower(config[idx+1:])
        for country, symbols := range countrySymbols {
            for _, symbol := range symbols {
                if strings.Contains(remark, strings.ToLower(symbol)) {
                    return country
                }
            }
        }
    }

    if strings.HasPrefix(config, "vmess://") {
        encoded := strings.TrimPrefix(config, "vmess://")
        if len(encoded)%4 != 0 {
            encoded += strings.Repeat("=", 4-len(encoded)%4)
        }
        decoded, err := base64.StdEncoding.DecodeString(encoded)
        if err == nil {
            var vmess struct {
                Ps string `json:"ps"`
            }
            if err := json.Unmarshal(decoded, &vmess); err == nil && vmess.Ps != "" {
                psLower := strings.ToLower(vmess.Ps)
                for _, ignore := range ignoreWords {
                    psLower = strings.ReplaceAll(psLower, ignore, "")
                }
                for country, symbols := range countrySymbols {
                    for _, symbol := range symbols {
                        if strings.Contains(psLower, strings.ToLower(symbol)) {
                            return country
                        }
                    }
                }
            }
        }
    }

    for country, symbols := range countrySymbols {
        for _, symbol := range symbols {
            if strings.Contains(configLower, strings.ToLower(symbol)) {
                return country
            }
        }
    }

    return "unknown"
}

// تابع جدید برای جداسازی بر اساس پروتکل (بهبودیافته)
func sortByProtocol() {
    inputFile := "All_Configs_Sorted.txt"
    outputDir := "Splitted-By-Protocol"

    if err := os.MkdirAll(outputDir, 0755); err != nil {
        fmt.Printf("Error creating protocol output directory: %v\n", err)
        return
    }

    file, err := os.Open(inputFile)
    if err != nil {
        fmt.Printf("Error opening input file: %v\n", err)
        return
    }
    defer file.Close()

    protocolFiles := make(map[string]*os.File)
    protocolWriters := make(map[string]*bufio.Writer)
    protocolConfigCount := make(map[string]int)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        protocol := "unknown"
        for _, proto := range protocols {
            if strings.HasPrefix(strings.ToLower(line), strings.ToLower(proto)) {
                // اعتبارسنجی اضافی برای Shadowsocks
                if proto == "ss://" {
                    if !isValidShadowsocksConfig(line) {
                        continue // نادیده گرفتن کانفیگ‌های نامعتبر
                    }
                }
                // اعتبارسنجی اضافی برای vmess
                if proto == "vmess://" {
                    if !isValidVmessConfig(line) {
                        continue
                    }
                }
                // اعتبارسنجی اضافی برای vless
                if proto == "vless://" {
                    if !isValidVlessConfig(line) {
                        continue
                    }
                }
                protocol = strings.TrimSuffix(proto, "://")
                break
            }
        }

        if _, ok := protocolFiles[protocol]; !ok {
            filename := filepath.Join(outputDir, protocol+".txt")
            f, err := os.Create(filename)
            if err != nil {
                fmt.Printf("Error creating file for %s: %v\n", protocol, err)
                continue
            }
            protocolFiles[protocol] = f
            protocolWriters[protocol] = bufio.NewWriter(f)
            protocolConfigCount[protocol] = 0
        }

        if _, err := protocolWriters[protocol].WriteString(line + "\n"); err != nil {
            fmt.Printf("Error writing to %s: %v\n", protocol, err)
            continue
        }
        protocolConfigCount[protocol]++
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
    }

    for protocol, writer := range protocolWriters {
        writer.Flush()
        protocolFiles[protocol].Close()
        fmt.Printf("Wrote %d configs to %s.txt\n", protocolConfigCount[protocol], protocol)
    }
}

// اعتبارسنجی Shadowsocks
func isValidShadowsocksConfig(config string) bool {
    if !strings.HasPrefix(config, "ss://") {
        return false
    }
    encoded := strings.TrimPrefix(config, "ss://")
    if idx := strings.Index(encoded, "#"); idx != -1 {
        encoded = encoded[:idx]
    }
    if idx := strings.Index(encoded, "?"); idx != -1 {
        encoded = encoded[:idx]
    }
    // بررسی base64 معتبر
    if len(encoded)%4 != 0 {
        encoded += strings.Repeat("=", 4-len(encoded)%4)
    }
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return false
    }
    // بررسی فرمت Shadowsocks (method:password@server:port)
    parts := strings.Split(string(decoded), "@")
    if len(parts) != 2 {
        return false
    }
    auth := strings.Split(parts[0], ":")
    if len(auth) != 2 {
        return false
    }
    serverPort := strings.Split(parts[1], ":")
    if len(serverPort) != 2 {
        return false
    }
    return true
}

// اعتبارسنجی Vmess
func isValidVmessConfig(config string) bool {
    if !strings.HasPrefix(config, "vmess://") {
        return false
    }
    encoded := strings.TrimPrefix(config, "vmess://")
    if len(encoded)%4 != 0 {
        encoded += strings.Repeat("=", 4-len(encoded)%4)
    }
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return false
    }
    var vmess struct {
        V   string `json:"v"`
        Ps  string `json:"ps"`
        Add string `json:"add"`
        Port string `json:"port"`
    }
    if err := json.Unmarshal(decoded, &vmess); err != nil {
        return false
    }
    return vmess.Add != "" && vmess.Port != ""
}

// اعتبارسنجی Vless
func isValidVlessConfig(config string) bool {
    if !strings.HasPrefix(config, "vless://") {
        return false
    }
    parts := strings.Split(config, "@")
    if len(parts) != 2 {
        return false
    }
    serverPort := strings.Split(parts[1], "?")
    if len(serverPort) < 1 {
        return false
    }
    addrPort := strings.Split(serverPort[0], ":")
    if len(addrPort) != 2 {
        return false
    }
    return true
}

func sortConfigs() {
    inputFile := "All_Configs_Sub.txt"
    outputFile := "All_Configs_Sorted.txt"

    file, err := os.Open(inputFile)
    if err != nil {
        fmt.Printf("Error opening input file: %v\n", err)
        return
    }
    defer file.Close()

    var configs []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" && !strings.HasPrefix(line, "#") {
            configs = append(configs, line)
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
        return
    }

    seen := make(map[string]bool)
    var uniqueConfigs []string
    for _, config := range configs {
        if !seen[config] {
            seen[config] = true
            uniqueConfigs = append(uniqueConfigs, config)
        }
    }

    out, err := os.Create(outputFile)
    if err != nil {
        fmt.Printf("Error creating output file: %v\n", err)
        return
    }
    defer out.Close()

    writer := bufio.NewWriter(out)
    defer writer.Flush()

    fixedText := `#profile-title: base64:8J+GkyBHaXRodWIgfCBEYW5pYWwgU2FtYWRpIPCfkI0=
#profile-update-interval: 1
#support-url: https://github.com/Giromo0/Collector
#profile-web-page-url: https://github.com/Giromo0/Collector
`
    if _, err := writer.WriteString(fixedText); err != nil {
        fmt.Printf("Error writing header: %v\n", err)
        return
    }

    for _, config := range uniqueConfigs {
        if _, err := writer.WriteString(config + "\n"); err != nil {
            fmt.Printf("Error writing config: %v\n", err)
            return
        }
    }

    fmt.Printf("Sorted %d unique configs into %s\n", len(uniqueConfigs), outputFile)
}

func sortByCountry() {
    inputFile := "All_Configs_Sorted.txt"
    outputDir := "Splitted-By-Country"

    if err := os.MkdirAll(outputDir, 0755); err != nil {
        fmt.Printf("Error creating output directory: %v\n", err)
        return
    }

    file, err := os.Open(inputFile)
    if err != nil {
        fmt.Printf("Error opening input file: %v\n", err)
        return
    }
    defer file.Close()

    countryFiles := make(map[string]*os.File)
    countryWriters := make(map[string]*bufio.Writer)
    countryConfigCount := make(map[string]int)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        country := identifyCountry(line)
        if country == "" {
            country = "unknown"
        }

        if _, ok := countryFiles[country]; !ok {
            filename := filepath.Join(outputDir, country+".txt")
            f, err := os.Create(filename)
            if err != nil {
                fmt.Printf("Error creating file for %s: %v\n", country, err)
                continue
            }
            countryFiles[country] = f
            countryWriters[country] = bufio.NewWriter(f)
            countryConfigCount[country] = 0
        }

        if _, err := countryWriters[country].WriteString(line + "\n"); err != nil {
            fmt.Printf("Error writing to %s: %v\n", country, err)
            continue
        }
        countryConfigCount[country]++
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
    }

    for country, writer := range countryWriters {
        writer.Flush()
        countryFiles[country].Close()
        fmt.Printf("Wrote %d configs to %s.txt\n", countryConfigCount[country], country)
    }
}
