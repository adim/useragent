package useragent

import (
	"bytes"
	"regexp"
	"strings"
)

// UserAgent struct containing all data extracted from parsed user-agent string
type UserAgent struct {
	VersionNo   VersionNo
	OSVersionNo VersionNo
	URL         string
	String      string
	Name        string
	Version     string
	OS          string
	OSVersion   string
	Device      string
	Mobile      bool
	Tablet      bool
	Desktop     bool
	Bot         bool
}

// Constants for browsers and operating systems for easier comparison
const (
	Windows        = "Windows"
	WindowsPhone   = "Windows Phone"
	WindowsNT      = "Windows NT"
	WindowsPhoneOS = "Windows Phone OS"
	Android        = "Android"
	MacOS          = "macOS"
	IOS            = "iOS"
	Linux          = "Linux"
	FreeBSD        = "FreeBSD"
	ChromeOS       = "ChromeOS"
	BlackBerry     = "BlackBerry"
	CrOS           = "CrOS"
	Harmony        = "Harmony"
	WebOS          = "webOS"
	Tizen          = "Tizen"

	Opera            = "Opera"
	OperaMini        = "Opera Mini"
	OperaTouch       = "Opera Touch"
	Chrome           = "Chrome"
	HeadlessChrome   = "Headless Chrome"
	Firefox          = "Firefox"
	InternetExplorer = "Internet Explorer"
	Safari           = "Safari"
	Edge             = "Edge"
	Vivaldi          = "Vivaldi"
	MobileSafari     = "Mobile Safari"
	NetFront         = "NetFront"
	Mozilla          = "Mozilla"
	Msie             = "MSIE"
	SamsungBrowser   = "Samsung Browser"
	LGBrowser         = "LG Browser"
	MiTVBrowser       = "Mi TV Browser"

	GoogleAdsBot        = "Google Ads Bot"
	Googlebot           = "Googlebot"
	Twitterbot          = "Twitterbot"
	FacebookExternalHit = "facebookexternalhit"
	Applebot            = "Applebot"
	Bingbot             = "Bingbot"
	YandexBot           = "YandexBot"
	YandexAdNet         = "YandexAdNet"

	FacebookApp  = "Facebook App"
	InstagramApp = "Instagram App"
	TiktokApp    = "TikTok App"

	Version = "Version"
	Mobile  = "Mobile"
	Tablet  = "Tablet"

	tablet = "tablet"
)

// Parse user agent string returning UserAgent struct
func Parse(userAgent string) UserAgent {
	ua := UserAgent{
		String: userAgent,
	}

	tokens := parse([]byte(userAgent))
	ua.URL = tokens.url

	// OS lookup
	switch {
	// Smart TV OS detection - check these first
	case tokens.exists("Web0S"):
		ua.OS = WebOS
		ua.OSVersion = tokens.get("Web0S")
		ua.Device = "Smart TV"

	case tokens.exists("SmartTV"):
		ua.OS = "SmartTV"
		ua.Device = "Smart TV"

	case tokens.exists(Android):
		ua.OS = Android
		var osIndex int
		osIndex, ua.OSVersion = tokens.getIndexValue(Android)
		ua.Tablet = strings.Contains(strings.ToLower(ua.String), tablet)
		ua.Device = tokens.findAndroidDevice(osIndex)

	case tokens.exists("iPhone"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPhone"
		ua.Mobile = true

	case tokens.exists("iPad"):
		ua.OS = IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPad"
		ua.Tablet = true

	case tokens.exists(WindowsNT):
		ua.OS = Windows
		ua.OSVersion = tokens.get(WindowsNT)
		ua.Desktop = true

	case tokens.exists(WindowsPhoneOS):
		ua.OS = WindowsPhone
		ua.OSVersion = tokens.get(WindowsPhoneOS)
		ua.Mobile = true

	case tokens.exists("Macintosh"):
		ua.OS = MacOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Desktop = true

	// Smart TV OS detection - check these before Linux
	case tokens.exists(Linux):
		ua.OS = Linux
		ua.OSVersion = tokens.get(Linux)
		ua.Desktop = true

	case tokens.exists(FreeBSD):
		ua.OS = FreeBSD
		ua.OSVersion = tokens.get(FreeBSD)
		ua.Desktop = true

	case tokens.exists(CrOS):
		ua.OS = ChromeOS
		ua.OSVersion = tokens.get(CrOS)
		ua.Desktop = true

	case tokens.exists(BlackBerry):
		ua.OS = BlackBerry
		ua.OSVersion = tokens.get(BlackBerry)
		ua.Mobile = true

	case tokens.exists("OpenHarmony"):
		ua.OS = Harmony
		ua.OSVersion = tokens.get("OpenHarmony")
		ua.Mobile = true
	}

	switch {
	case strings.Contains(ua.String, "AFTSSS") || strings.Contains(ua.String, "AFTBOXE1") || strings.Contains(ua.String, "AFTGAZL") || strings.Contains(ua.String, "OLED TV") || strings.Contains(ua.String, "Chromecast") || strings.Contains(ua.String, "MiTV-AFKR0") || strings.Contains(ua.String, "H96 Max") || strings.Contains(ua.String, "RK3318") || (strings.Contains(ua.String, "SMART-TV") && strings.Contains(ua.String, "Tizen")):
		uaString := ua.String
		ua.Device = "Smart TV"
		if strings.Contains(uaString, "AFTSSS") || strings.Contains(uaString, "AFTBOXE1") || strings.Contains(uaString, "AFTGAZL") {
			ua.OS = Android
			ua.OSVersion = tokens.get("Android")
			ua.Name = "Amazon Fire TV Browser"
			ua.Version = tokens.get("Chrome")
		} else if strings.Contains(uaString, "OLED TV") || strings.Contains(uaString, "H96 Max") || strings.Contains(uaString, "RK3318") {
			ua.OS = Android
			ua.OSVersion = tokens.get("Android")
			ua.Name = "Android TV Browser"
			// Handle both Chrome-based and Firefox-based Android TV devices
			if chromeVersion := tokens.get("Chrome"); chromeVersion != "" {
				ua.Version = chromeVersion
			} else if firefoxVersion := tokens.get("Firefox"); firefoxVersion != "" {
				ua.Version = firefoxVersion
			}
		} else if strings.Contains(uaString, "Chromecast") {
			ua.OS = Android
			ua.OSVersion = tokens.get("Android")
			ua.Name = "Chromecast"
			ua.Version = tokens.get("Chrome")
		} else if strings.Contains(uaString, "MiTV-AFKR0") {
			ua.OS = Android
			ua.OSVersion = tokens.get("Android")
			ua.Name = MiTVBrowser
			ua.Version = tokens.get("Chrome")
		} else if strings.Contains(uaString, "SMART-TV") && strings.Contains(uaString, "Tizen") {
			ua.OS = Tizen
			ua.OSVersion = tokens.get("Tizen")
			ua.Name = "Samsung TV Browser"
			if tokens.exists("SamsungBrowser") {
				ua.Version = tokens.get("SamsungBrowser")
			} else {
				// Handle newer Tizen format like "108.0.5359.1/8.0"
				if version := tokens.get("108.0.5359.1"); version != "" {
					ua.Version = "108.0.5359.1"
				}
			}
		}
		// Skip the rest of browser detection
		return ua

	case tokens.exists(Googlebot):
		ua.Name = Googlebot
		ua.Version = tokens.get(Googlebot)
		ua.Bot = true
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.existsAny("GoogleProber", "GoogleProducer"):
		if name := tokens.findBestMatch(false); name != "" {
			ua.Name = name
		}
		ua.Bot = true

	case tokens.exists("Bytespider"):
		ua.Name = "Bytespider"
		ua.Mobile = tokens.exists("Mobile Safari")
		ua.Bot = true

	case tokens.exists(Applebot):
		ua.Name = Applebot
		ua.Version = tokens.get(Applebot)
		ua.Bot = true
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)
		ua.OS = ""

	case tokens.get(OperaMini) != "":
		ua.Name = OperaMini
		ua.Version = tokens.get(OperaMini)
		ua.Mobile = true

	case tokens.get("OPR") != "":
		ua.Name = Opera
		ua.Version = tokens.get("OPR")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get("OPT") != "":
		ua.Name = OperaTouch
		ua.Version = tokens.get("OPT")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	// Opera on iOS
	case tokens.get("OPiOS") != "":
		ua.Name = Opera
		ua.Version = tokens.get("OPiOS")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	// Chrome on iOS
	case tokens.get("CriOS") != "":
		ua.Name = Chrome
		ua.Version = tokens.get("CriOS")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	// Firefox on iOS
	case tokens.get("FxiOS") != "":
		ua.Name = Firefox
		ua.Version = tokens.get("FxiOS")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get(Firefox) != "":
		ua.Name = Firefox
		ua.Version = tokens.get(Firefox)
		ua.Mobile = tokens.exists(Mobile)
		ua.Tablet = tokens.exists(Tablet)

	case tokens.get(Vivaldi) != "":
		ua.Name = Vivaldi
		ua.Version = tokens.get(Vivaldi)

	case tokens.exists(Msie):
		ua.Name = InternetExplorer
		ua.Version = tokens.get(Msie)

	case tokens.get("EdgiOS") != "":
		ua.Name = Edge
		ua.Version = tokens.get("EdgiOS")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get(Edge) != "":
		ua.Name = Edge
		ua.Version = tokens.get(Edge)
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get("Edg") != "":
		ua.Name = Edge
		ua.Version = tokens.get("Edg")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get("EdgA") != "":
		ua.Name = Edge
		ua.Version = tokens.get("EdgA")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get("bingbot") != "":
		ua.Name = Bingbot
		ua.Version = tokens.get("bingbot")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get(YandexBot) != "":
		ua.Name = YandexBot
		ua.Version = tokens.get(YandexBot)
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)
		ua.Bot = true

	case tokens.get(YandexAdNet) != "":
		ua.Name = YandexAdNet
		ua.Version = tokens.get(YandexAdNet)
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)
		ua.Bot = true

	
	case tokens.exists("SamsungBrowser") && tokens.exists("TV"):
		ua.Name = "Samsung TV Browser"
		ua.Version = tokens.get("SamsungBrowser")
		ua.Device = "Smart TV"
		// Skip the rest of browser detection
		return ua

	case tokens.exists("SMART-TV") && tokens.exists("Tizen"):
		ua.Name = "Samsung TV Browser"
		ua.Version = tokens.get("Safari")
		ua.Device = "Smart TV"
		// Skip the rest of browser detection
		return ua

	case tokens.exists("SMART-TV") && !tokens.exists("Tizen"):
		// Handle newer Tizen user agents that don't explicitly mention Tizen
		ua.Name = "Samsung TV Browser"
		ua.Version = tokens.findBestMatchVersion()
		ua.Device = "Smart TV"
		// Skip the rest of browser detection
		return ua

	case tokens.get("SamsungBrowser") != "":
		ua.Name = SamsungBrowser
		ua.Version = tokens.get("SamsungBrowser")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)
		ua.OS = Android

	case tokens.get("HeadlessChrome") != "":
		ua.Name = HeadlessChrome
		ua.Version = tokens.get("HeadlessChrome")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)
		ua.Bot = true

	case tokens.existsAny("AdsBot-Google-Mobile", "Mediapartners-Google", "AdsBot-Google"):
		ua.Name = GoogleAdsBot
		ua.Bot = true
		ua.Mobile = ua.IsAndroid() || ua.IsIOS()

	case tokens.exists("Yahoo Ad monitoring"):
		ua.Name = "Yahoo Ad monitoring"
		ua.Bot = true
		ua.Mobile = ua.IsAndroid() || ua.IsIOS()

	case tokens.exists("XiaoMi"):
		miui := tokens.get("XiaoMi")
		if strings.HasPrefix(miui, "MiuiBrowser") {
			ua.Name = "Miui Browser"
			ua.Version = strings.TrimPrefix(miui, "MiuiBrowser/")
			ua.Mobile = true
		}

	case tokens.exists("FBAN"):
		ua.Name = FacebookApp
		ua.Version = tokens.get("FBAN")

	case tokens.exists("FB_IAB"):
		ua.Name = FacebookApp
		ua.Version = tokens.get("FBAV")

	case tokens.startsWith("Instagram"):
		ua.Name = InstagramApp
		ua.Version = tokens.findInstagramVersion()

	case tokens.exists("BytedanceWebview"):
		ua.Name = TiktokApp
		ua.Version = tokens.get("app_version")

	case tokens.get("HuaweiBrowser") != "":
		ua.Name = "Huawei Browser"
		ua.Version = tokens.get("HuaweiBrowser")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.get("LG Browser") != "":
		ua.Name = LGBrowser
		ua.Version = tokens.get("LG Browser")
		ua.Device = "Smart TV"

	case tokens.exists(BlackBerry):
		ua.Name = BlackBerry
		ua.Version = tokens.get(Version)

	case tokens.exists(NetFront):
		ua.Name = NetFront
		ua.Version = tokens.get(NetFront)
		ua.Mobile = true

	
	// if Chrome and Safari defined, find any other token sent descr
	case tokens.exists(Chrome) && tokens.exists(Safari):
		name := tokens.findBestMatch(true)
		if name != "" {
			ua.Name = name
			ua.Version = tokens.get(name)
			break
		}
		fallthrough

	case tokens.exists(Chrome):
		ua.Name = Chrome
		ua.Version = tokens.get(Chrome)
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.exists("Brave Chrome"):
		ua.Name = Chrome
		ua.Version = tokens.get("Brave Chrome")
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	case tokens.exists(Safari):
		ua.Name = Safari
		v := tokens.get(Version)
		if v != "" {
			ua.Version = v
		} else {
			ua.Version = tokens.get(Safari)
		}
		ua.Mobile = tokens.existsAny(Mobile, MobileSafari)

	default:
		if ua.IsAndroid() && tokens.get(Version) != "" {
			ua.Name = "Android browser"
			ua.Version = tokens.get(Version)
			ua.Mobile = true
		} else {
			if name := tokens.findBestMatch(false); name != "" {
				ua.Name = name
				ua.Version = tokens.get(name)
			} else {
				ua.Name = ua.String
			}
			ua.Bot = strings.Contains(strings.ToLower(ua.Name), "bot")
			// If mobile flag has already been set, don't override it.
			if !ua.Mobile {
				ua.Mobile = tokens.existsAny(Mobile, MobileSafari)
			}
		}
	}

	if ua.IsAndroid() {
		ua.Mobile = true
	}

	// if tablet, switch mobile to off
	if ua.Tablet {
		ua.Mobile = false
	}

	// if not already bot, check some popular bots and whether URL is set
	if !ua.Bot {
		switch ua.Name {
		case Twitterbot, FacebookExternalHit, "facebookcatalog":
			ua.Bot = true
		default:
			ua.Bot = ua.URL != ""
		}
	}

	ua.VersionNo = parseVersion(ua.Version)
	ua.OSVersionNo = parseVersion(ua.OSVersion)

	return ua
}

// var buffPool = sync.Pool{New: func() interface{} {
// 	return bytes.NewBuffer(make([]byte, 0, 30))
// }}

func parse(userAgent []byte) properties {
	clients := properties{
		list: make([]property, 0, 8),
	}
	slash := false
	isURL := false
	// buff := buffPool.Get().(*bytes.Buffer)
	// val := buffPool.Get().(*bytes.Buffer)
	// buff.Reset()
	// val.Reset()

	buff := bytes.NewBuffer(make([]byte, 0, 30))
	val := bytes.NewBuffer(make([]byte, 0, 30))

	addToken := func() {
		if buff.Len() != 0 {
			s := string(bytes.TrimSpace(buff.Bytes()))
			if !ignore(s) {
				if isURL {
					clients.url = strings.TrimPrefix(s, "+")
					return
				}
				if val.Len() == 0 {
					// if value don't exists, try to get version from the token
					clients.list = append(clients.list, checkVer(s))
				} else {
					clients.list = append(clients.list, property{Key: s, Value: string(bytes.TrimSpace(val.Bytes()))})
				}
			}
		}
		buff.Reset()
		val.Reset()
		slash = false
		isURL = false
	}

	parOpen := false
	braOpen := false

	for i, c := range userAgent {
		switch {
		case c == 41: // )
			addToken()
			parOpen = false

		case (parOpen || braOpen) && c == 59: // ;
			addToken()

		case c == 59: // ;
			addToken()

		case c == 40: // (
			addToken()
			parOpen = true

		case c == 91: // [
			addToken()
			braOpen = true
		case c == 93: // ]
			addToken()
			braOpen = false

		case c == 58: // :
			if bytes.HasSuffix(buff.Bytes(), []byte("http")) || bytes.HasSuffix(buff.Bytes(), []byte("https")) {
				// If we are part of a URL just write the character.
				buff.WriteByte(c)
			} else if i != len(userAgent)-1 && userAgent[i+1] != ' ' {
				// If the following character is not a space, change to a space.
				buff.WriteByte(' ')
			}
			// Otherwise don't write as it's probably a badly formatted key value separator.

		case slash && c == 32:
			addToken()

		case slash:
			val.WriteByte(c)

		case c == 47 && !isURL: //   /
			if i != len(userAgent)-1 && userAgent[i+1] == 47 && (bytes.HasSuffix(buff.Bytes(), []byte("http:")) || bytes.HasSuffix(buff.Bytes(), []byte("https:"))) {
				buff.WriteByte(c)
				isURL = true
			} else {
				if ignore(buff.String()) {
					buff.Reset()
				} else {
					slash = true
				}
			}

		default:
			buff.WriteByte(c)
		}
	}
	addToken()

	// buffPool.Put(buff)
	// buffPool.Put(val)
	return clients
}

func checkVer(s string) property {
	i := strings.LastIndex(s, " ")
	if i == -1 {
		return property{Key: s, Value: ""}
	}

	switch s[:i] {
	case Linux, WindowsNT, WindowsPhoneOS, Msie, Android, "OpenHarmony":
		return property{Key: s[:i], Value: s[i+1:]}
	case "CrOS x86_64", "CrOS aarch64", "CrOS armv7l":
		j := strings.LastIndex(s[:i], " ")
		return property{Key: s[:j], Value: s[j+1 : i]}
	default:
		return property{Key: s, Value: ""}
	}
}

// ignore returns true if token should be ignored
func ignore(s string) bool {
	switch s {
	case "KHTML, like Gecko", "U", "compatible", Mozilla, "WOW64", "en", "en-us", "en-gb", "ru-ru", "Browser":
		return true
	default:
		return false
	}
}

type property struct {
	Key   string
	Value string
}
type properties struct {
	list []property
	url  string
}

func (p properties) get(key string) string {
	for _, prop := range p.list {
		if prop.Key == key {
			return prop.Value
		}
	}
	return ""
}

func (p properties) getIndexValue(key string) (int, string) {
	for i, prop := range p.list {
		if prop.Key == key {
			return i, prop.Value
		}
	}
	return -1, ""
}

func (p properties) exists(key string) bool {
	for _, prop := range p.list {
		if prop.Key == key {
			return true
		}
	}
	return false
}

// func (p properties) existsIgnoreCase(key string) bool {
// 	for _, prop := range p.list {
// 		if strings.EqualFold(prop.Key, key) {
// 			return true
// 		}
// 	}
// 	return false
// }

func (p properties) existsAny(keys ...string) bool {
	for _, k := range keys {
		for _, prop := range p.list {
			if prop.Key == k {
				return true
			}
		}
	}
	return false
}

func (p properties) getAny(keys ...string) (key, value string) {
	for _, k := range keys {
		for _, prop := range p.list {
			if prop.Key == k {
				return prop.Key, prop.Value
			}
		}
	}
	return "", ""
}

func (p properties) findMacOSVersion() string {
	for _, token := range p.list {
		if strings.Contains(token.Key, "OS") {
			if ver := findVersion(token.Value); ver != "" {
				return ver
			} else if ver = findVersion(token.Key); ver != "" {
				return ver
			}
		}

	}
	return ""
}

func (p properties) startsWith(value string) bool {
	for _, prop := range p.list {
		if strings.HasPrefix(prop.Key, value) {
			return true
		}
	}
	return false
}

func (p properties) findInstagramVersion() string {
	for _, token := range p.list {
		if strings.HasPrefix(token.Key, "Instagram") {
			if ver := findVersion(token.Value); ver != "" {
				return ver
			} else if ver = findVersion(token.Key); ver != "" {
				return ver
			}
		}

	}
	return ""
}

// findBestMatchVersion from the rest of the bunch
func (p properties) findBestMatchVersion() string {
	for _, prop := range p.list {
		switch prop.Key {
		case Chrome, Firefox, Safari, Version, Mobile, MobileSafari, Mozilla, "AppleWebKit", WindowsNT, WindowsPhoneOS, Android, "Macintosh", Linux, "GSA", CrOS, Tablet, "OpenHarmony":
		default:
			// Check if it looks like a version number
			if strings.Contains(prop.Value, ".") && len(prop.Value) > 3 {
				return prop.Value
			}
			// Check if the key itself looks like a version
			if strings.Contains(prop.Key, ".") && len(prop.Key) > 3 {
				return prop.Key
			}
		}
	}
	return ""
}

// findBestMatch from the rest of the bunch
// in first cycle only return key with version value
// if withVerValue is false, do another cycle and return any token
func (p properties) findBestMatch(withVerOnly bool) string {
	n := 2
	if withVerOnly {
		n = 1
	}
	for i := 0; i < n; i++ {
		for _, prop := range p.list {
			switch prop.Key {
			case Chrome, Firefox, Safari, Version, Mobile, MobileSafari, Mozilla, "AppleWebKit", WindowsNT, WindowsPhoneOS, Android, "Macintosh", Linux, "GSA", CrOS, Tablet, "OpenHarmony":
			default:
				// don't pick if starts with number
				if len(prop.Key) != 0 && prop.Key[0] >= 48 && prop.Key[0] <= 57 {
					break
				}
				if i == 0 {
					if prop.Value != "" { // in first check, only return keys with value
						return prop.Key
					}
				} else {
					return prop.Key
				}
			}
		}
	}
	return ""
}

var rxMacOSVer = regexp.MustCompile(`[_\d\.]+`)

func findVersion(s string) string {
	if ver := rxMacOSVer.FindString(s); ver != "" {
		return strings.Replace(ver, "_", ".", -1)
	}
	return ""
}

// findAndroidDevice in tokens
func (p *properties) findAndroidDevice(startIndex int) string {
	for i := startIndex; i < startIndex+1; i++ {
		if len(p.list) > i+1 {
			dev := p.list[i+1].Key
			if len(dev) == 2 || (len(dev) == 5 && dev[2] == '-') {
				// probably language tag (en-us etc..), ignore and continue loop
				continue
			}
			switch dev {
			case Chrome, Firefox, Safari, OperaMini, "Presto", Version, Mobile, MobileSafari, Mozilla, "AppleWebKit", WindowsNT, WindowsPhoneOS, Android, "Macintosh", Linux, CrOS:
				// ignore these tokens, not device names
			default:
				if strings.Contains(strings.ToLower(dev), tablet) {
					p.list[i+1].Key = Tablet // leave Tablet tag for later table detection
				} else {
					p.list = append(p.list[:i+1], p.list[i+2:]...)
				}
				return strings.TrimSpace(strings.TrimSuffix(dev, "Build"))
			}
		}
	}
	return ""
}
