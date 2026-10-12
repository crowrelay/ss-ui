<div dir="rtl">

# ss-ui — پنل S-UI همراه با هستهٔ sing-box مجهز به پروتکل Sudoku

[English](../en/README.md) | **فارسی**

> توجه: متن این صفحه راست‌به‌چپ است، اما بلوک‌های کد و دستورها به‌صورت چپ‌به‌راست (LTR) نمایش داده می‌شوند.

**ss-ui** یک انشعاب (fork) از پروژهٔ [s-ui](https://github.com/alireza0/s-ui) است؛ s-ui یک پنل وب و مدیریت کاربر برای هستهٔ پروکسی sing-box است. در ss-ui این هسته به‌جای نسخهٔ رسمی، با انشعابی از sing-box که پروتکل **Sudoku** را داراست کامپایل می‌شود.

- پنل: همین مخزن.
- هسته: [`SUDOKU-ASCII/sudoku-sing-box`](https://github.com/SUDOKU-ASCII/sudoku-sing-box) (نسخهٔ sing-box 1.14.2 به‌همراه Sudoku).
- پروتکل: [`SUDOKU-ASCII/sudoku`](https://github.com/SUDOKU-ASCII/sudoku).

هستهٔ sing-box در اینجا یک فایل اجرایی جداگانه نیست؛ به‌صورت کتابخانه (library) داخل فایل اجرایی `ss-ui` لینک می‌شود. جای‌گزینی هسته فقط یک دستور `replace` در `go.mod` است، بنابراین پنل و هسته همیشه به‌صورت یک فایل اجرایی واحد عرضه می‌شوند.

---

## فهرست مطالب

1. [پروتکل Sudoku چیست](#۱-پروتکل-sudoku-چیست)
2. [ss-ui چیست](#۲-ss-ui-چیست)
3. [نصب](#۳-نصب)
4. [اولین اجرا و ساختِ inbound سو‌دوکو](#۴-اولین-اجرا-و-ساختinbound-سودوکو)
5. [اتصال کلاینت](#۵-اتصال-کلاینت)
6. [به‌روزرسانی و حذف](#۶-بهروزرسانی-و-حذف)
7. [به‌روزرسانی خودکار](#۷-بهروزرسانی-خودکار)
8. [پرسش‌های متداول و عیب‌یابی](#۸-پرسشهای-متداول-و-عیبیابی)

---

## ۱. پروتکل Sudoku چیست

Sudoku یک پروتکل مبهم‌سازی ترافیک است که به‌جای تصادفی‌بودن، روی **حل معمای سودوکوی ۴×۴** بنا شده است. هر بایت دادهٔ ورودی (۲۵۶ حالت) به یکی از ۲۸۸ حالت غیرهم‌شکلِ سودوکوی ۴×۴ نگاشت می‌شود. چون هر معما را می‌توان با چند مجموعهٔ معتبر از سرنخ‌ها (clue) ساخت، دادهٔ یکسان به شکل‌های متفاوتی روی سیم رمزگذاری می‌شود و می‌توان ظاهر ترافیک را آزادانه تعیین کرد.

ویژگی‌های مهم:

- **پدینگ دلخواه** در هر موقعیت، بدون نیاز به توافق قبلی روی طول پدینگ.
- **حالت‌های متن‌آشکار/کم‌آنتروپی:** `prefer_ascii` (حدود ۹۸٪ کاراکتر چاپی ASCII) یا `prefer_entropy` (حدود ۶۵٪ ASCII و وزن همینگ پایین) — کمتر از بازه‌ای که در گزارش GFW نشانهٔ مسدودسازی است.
- **اثر انگشت قابل تعیین توسط کاربر:** با `custom_table` / `custom_tables` می‌توانید ظاهر بایت‌ها را انتخاب کنید. استفادهٔ ترکیبی میان کاربران توصیه می‌شود.
- **HTTPMask:** می‌توان جریان را شبیه ترافیک HTTP(S)/WebSocket درآورد.
- **بازگشت دفاعی (fallback):** وقتی دست‌دهی مشکوک باشد (اسکنر، کلید اشتباه، دادهٔ نامعتبر)، سرور اتصال را بی‌صدا به یک آدرس طعمه (`fallback_address`) هدایت می‌کند — مثلاً یک صفحهٔ Nginx، یا **یک سرور Sudoku داخلی** برای راه‌اندازی دو‌لایه.
- **رمزنگاری:** در لایهٔ زیرین مبهم‌سازی، به‌صورت اختیاری AEAD به‌کار می‌رود — `chacha20-poly1305` (پیش‌فرض) یا `aes-128-gcm`، با کلیدِ هر اتصال و محافظت ضدِ بازپخش (anti-replay).
- **UDP روی TCP (UoT):** ترافیک UDP داخل TCP حمل می‌شود؛ شنوندهٔ UDP خام وجود ندارد.

هشدارها:

- مبتنی بر TCP است؛ UDP از طریق UoT منتقل می‌شود.
- هستهٔ خالی تنها SOCKS5/HTTP ارائه می‌دهد مگر از کلاینت‌های رسمی گرافیکی استفاده کنید.
- پشتیبانی رسمی کلاینت محدود است (Mihomo، Sudodroid، sudoku4x4، Anywhere) — بخش [۵](#۵-اتصال-کلاینت) را ببینید.

### کلید: حالت مشترک در برابر حالت جداگانه

Sudoku با یک «کلید» احراز هویت می‌کند.

- **حالت کلید مشترک (shared):** کلاینت و سرور از **یک رشتهٔ یکسان** به‌عنوان `key` استفاده می‌کنند. ساده‌ترین حالت.
- **حالت کلید جداگانه (split):** سرور از **کلید عمومی (public)** و کلاینت از **کلید خصوصی (private)** متناظر استفاده می‌کند. یک جفت‌کلید با ss-ui بسازید و فقط کلید خصوصی را به کلاینت‌ها بدهید.

---

## ۲. ss-ui چیست

ss-ui همان پنل s-ui است که هستهٔ Sudoku در آن تنیده شده. همهٔ امکانات s-ui در اینجا هست — مدیریت چند کلاینت، اشتراک (subscription)، آمار ترافیک، مسیریابی، TLS — به‌علاوهٔ:

- پشتیبانی از پروتکل Sudoku در ویرایشگر inbound و outbound.
- **مولد جفت‌کلید** Sudoku.
- تولید اشتراک برای کلاینت‌های Sudoku.
- [گردش‌کار به‌روزرسانی خودکار](#۷-بهروزرسانی-خودکار).

نقاط دسترسی پیش‌فرض پس از نصب:

| مورد | مقدار |
|------|-------|
| پورت پنل | 2095 |
| مسیر پنل | `/app/` |
| پورت اشتراک | 2096 |
| مسیر اشتراک | `/sub/` |
| نام کاربری/رمز پیش‌فرض | `admin` / `admin` |

> بلافاصله پس از اولین ورود، اطلاعات ورود پیش‌فرض را تغییر دهید.

---

## ۳. نصب

### لینوکس / مک

```sh
bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

اگر مخزن شما با مالک/نام دیگری میزبانی می‌شود، آن را بازنویسی کنید:

```sh
SS_UI_REPO=owner/ss-ui bash <(curl -Ls https://raw.githubusercontent.com/owner/ss-ui/main/install.sh)
```

نصب‌کننده به‌صورت پیش‌فرض انگلیسی است؛ برای فارسی `SUI_LANG=fa` را تنظیم کنید:

```sh
SUI_LANG=fa bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

### آلپاین لینوکس (OpenRC)

آلپاین به‌صورت پیش‌فرض `bash` ندارد:

```sh
apk add bash
bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

مدیریت سرویس با OpenRC:

```sh
rc-service ss-ui start|stop|restart
rc-update add ss-ui default
```

### داکر

با Compose:

```sh
mkdir ss-ui && cd ss-ui
wget -q https://raw.githubusercontent.com/crowrelay/ss-ui/main/docker-compose.yml
docker compose up -d
```

یا با `docker run`:

```sh
mkdir ss-ui && cd ss-ui
docker run -itd \
    -p 2095:2095 -p 2096:2096 -p 443:443 -p 80:80 \
    -v $PWD/db/:/app/db/ \
    -v $PWD/cert/:/root/cert/ \
    --name ss-ui --restart=unless-stopped \
    ghcr.io/crowrelay/ss-ui:latest
```

### ویندوز

1. آخرین بستهٔ ویندوز (`ss-ui-windows-amd64.zip`) را از صفحهٔ انتشارها دانلود کنید.
2. آن را استخراج کنید.
3. فایل `install-windows.bat` را **با دسترسی Administrator** اجرا و ویزارد را دنبال کنید.
4. پنل را در `http://localhost:2095/app` باز کنید.

### ساخت از منبع

پیش‌نیازها: Go (نسخه در `go.mod`)، Node.js 26 و زنجیره‌ابزار C برای CGO.

```sh
git clone https://github.com/crowrelay/ss-ui
cd ss-ui
git submodule update --init --recursive

# یک‌مرحله‌ای: فرانت‌اند و بک‌اند را در ./ss-ui می‌سازد
./build.sh
```

یا گام‌به‌گام:

```sh
(cd frontend && npm ci && npm run build)
rm -fr web/html/* && cp -R frontend/dist/* web/html/
. ./build-tags.sh
go build -ldflags "$(ldflags_for dev)" -tags "$(tags_for dev)" -o ss-ui main.go
./ss-ui
```

> پس از تغییر نسخهٔ هسته در `go.mod`، یک‌بار `go mod tidy` را اجرا کنید تا `go.sum` پیش از ساخت، انشعاب را ثبت کند.

---

## ۴. اولین اجرا و ساختِ inbound سو‌دوکو

سرویس را راه‌اندازی کنید (`systemctl enable ss-ui --now`، `rc-service ss-ui start` یا اجرای `./ss-ui`) و در `http://<host>:2095/app` با `admin` / `admin` وارد شوید.

### ساخت جفت‌کلید

یا از خط فرمان:

```sh
./ss-ui generate sudoku-keypair
# PrivateKey: ...
# PublicKey:  ...
```

یا از پنل: **Inbounds → Add → نوع `Sudoku` → آیکن بازآوری (refresh) کنار «Key».**

- **حالت کلید مشترک:** همان رشته را در **Key** اینباند و **Key** کلاینت بگذارید.
- **حالت کلید جداگانه:** **PublicKey** را در اینباند و **PrivateKey** را در کلاینت (`out_json`) قرار دهید تا درز اطلاعاتِ سرور، کلید کلاینت را فاش نکند.

### فیلدهای inbound

| فیلد | توضیح |
|------|-------|
| `key` | رمز مشترک، یا کلید عمومی در حالت جداگانه. |
| `aead` | `chacha20-poly1305` (پیش‌فرض)، `aes-128-gcm` یا `none` (فقط آزمایشی). |
| `ascii` | `prefer_ascii` یا `prefer_entropy`. |
| `padding_min` / `padding_max` | کرانه‌های پدینگ، ۰ تا ۱۰۰. |
| `custom_table` | رشتهٔ چیدمان، مثلاً `xpxvvpvv` (۲ تا `x`، ۲ تا `p`، ۴ تا `v`). |
| `custom_tables` | چند چیدمان؛ ترکیب میان کاربران توصیه می‌شود. |
| `enable_pure_downlink` | `true` (پیش‌فرض) = دانلینک خالص سو‌دوکو؛ `false` = دانلینک بهینه از نظر پهنا (نیازمند `aead != none`). |
| `multiplex` | `off` یا `auto`. |
| `httpmask` | شبیه‌سازی جریان به‌صورت HTTP(S)/WebSocket: `mode` (`auto`/`legacy`/`ws`)، `path_root` و در کلاینت `host` و `tls`. |
| `handshake_timeout` | مهلت دست‌دهی بر حسب ثانیه. |
| `suspicious_action` | برای اتصال‌های کاوشگر: `fallback` (پیش‌فرض) یا `silent`. |
| `fallback_address` | مقصد هدایت اتصال مشکوک (صفحهٔ طعمه، یا یک سرور Sudoku داخلی). |

یک اینباند حداقلی:

```json
{
  "type": "sudoku",
  "tag": "sudoku-in",
  "listen": "0.0.0.0",
  "listen_port": 8080,
  "key": "YOUR_PUBLIC_OR_SHARED_KEY",
  "aead": "chacha20-poly1305",
  "ascii": "prefer_entropy",
  "padding_min": 5,
  "padding_max": 15,
  "enable_pure_downlink": true,
  "handshake_timeout": 5,
  "fallback_address": "127.0.0.1:80",
  "suspicious_action": "fallback"
}
```

### راه‌اندازی دو‌لایه (زنجیره‌ای)

یک سرور Sudoku بیرونی را به‌عنوان نقطهٔ ورود عمومی قرار دهید (استتار و مقاومت در برابر کاوش). به آن کلیدی **متفاوت** از سرور داخلی بدهید و `suspicious_action: "fallback"` را با `fallback_address` که به سرور Sudoku داخلی اشاره می‌کند تنظیم کنید. کلاینت‌های عادی در دست‌دهیِ سرور بیرونی شکست می‌خورند و به سرور داخلی هدایت می‌شوند که دست‌دهی واقعی را انجام می‌دهد. کلاینت به آدرس بیرونی وصل می‌شود اما از **کلید سرور داخلی** استفاده می‌کند.

---

## ۵. اتصال کلاینت

Sudoku توسط این کلاینت‌ها پشتیبانی می‌شود:

| پلتفرم | کلاینت |
|--------|--------|
| اندروید | [CMFA](https://github.com/MetaCubeX/ClashMetaForAndroid) (Mihomo)، [Sudodroid](https://github.com/SUDOKU-ASCII/sudoku-android) |
| iOS | [Anywhere Proxy](https://apps.apple.com/us/app/anywhere-proxy/id6758235178) |
| دسکتاپ | [sudoku4x4](https://github.com/SUDOKU-ASCII/sudoku-desktop) |
| روتر | [OpenClash](https://github.com/vernesong/OpenClash) |
| سرور/ابزار | [easy-install](https://github.com/SUDOKU-ASCII/easy-install) |

در ss-ui یک **کلاینت** بسازید، اینباند Sudoku را به آن تخصیص دهید و اشتراک را به اشتراک بگذارید. outboundِ Sudoku که کلاینت دریافت می‌کند از بلوک **Client `out_json`** اینباند ساخته می‌شود — اینباند را باز کنید، به تب **Client** بروید و برای حالت کلید جداگانه، `key` کلاینت را آنجا تنظیم کنید.

ساختار outbound کلاینت:

```json
{
  "type": "sudoku",
  "tag": "sudoku-out",
  "server": "example.com",
  "server_port": 443,
  "key": "YOUR_PRIVATE_OR_SHARED_KEY",
  "aead": "chacha20-poly1305",
  "ascii": "prefer_entropy",
  "padding_min": 5,
  "padding_max": 15,
  "enable_pure_downlink": true,
  "httpmask": { "disable": false, "mode": "auto", "tls": true, "host": "example.com", "path_root": "aabbcc" }
}
```

---

## ۶. به‌روزرسانی و حذف

به‌روزرسانی به آخرین نسخه:

```sh
bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh)
```

نصب یک نسخهٔ مشخص:

```sh
VERSION=v1.6.4-sudoku.1 bash <(curl -Ls https://raw.githubusercontent.com/crowrelay/ss-ui/main/install.sh) $VERSION
```

حذف (systemd):

```sh
systemctl disable ss-ui --now
rm -f /etc/systemd/system/ss-ui.service
systemctl daemon-reload
rm -fr /usr/local/ss-ui /usr/bin/ss-ui
```

آلپاین (OpenRC):

```sh
rc-service ss-ui stop
rc-update del ss-ui default
rm -f /etc/init.d/ss-ui
rm -fr /usr/local/ss-ui /usr/bin/ss-ui
```

---

## ۷. به‌روزرسانی خودکار

گردش‌کار `.github/workflows/sync-upstream.yml` روزانه (و در صورت درخواست) اجرا می‌شود:

1. آخرین برچسب (tag) `sudoku-sing-box` را پیدا و خط `replace` را در `go.mod` بازنویسی می‌کند.
2. مخزن `alireza0/s-ui` را دریافت و به‌صورت بهترین‌تلاش ادغام می‌کند (`-X ours`، بنابراین فایل‌های سفارشی‌شدهٔ این انشعاب و اشاره‌گر زیرماژول فرانت‌اند برنده می‌شوند).
3. `go mod tidy` و یک بررسی ساخت و تست‌های هسته/ابزار را اجرا می‌کند.
4. کامیت کرده، به `main` پوش می‌کند و یک برچسب نسخه می‌سازد — که `release.yml` را برای ساخت همهٔ پلتفرم‌ها و انتشار نسخه فعال می‌کند.

در صورت بروز تعارض در ادغام با بالادست، ادغام لغو می‌شود (فقط به‌روزرسانی هسته منتشر می‌شود) و یک issue با برچسب `sync-conflict` برای بررسی انسان باز می‌شود. فایل `.gitattributes` مسیر `frontend` را `merge=ours` علامت می‌زند؛ فرانت‌اند از انشعاب خودش بازآوری می‌شود.

---

## ۸. پرسش‌های متداول و عیب‌یابی

**هسته بالا نمی‌آید / خطای «missing key».**
اینباند Sudoku به `key` غیرخالی نیاز دارد. یکی بسازید و جای‌گذاری کنید.

**کلاینت وصل می‌شود اما چیزی بار نمی‌شود.**
بررسی کنید کلید کلاینت و سرور یکسان باشند (حالت مشترک)، یا کلید خصوصی کلاینت با کلید عمومی سرور جفت باشد (حالت جداگانه)، و مقادیر `aead`، `ascii` و `custom_table` سازگار باشند.

**`prefer_ascii` یا `prefer_entropy`؟**
`prefer_ascii` شبیه متن آشکار است (حدود ۹۸٪ ASCII)؛ `prefer_entropy` آنتروپی پایین‌تر و ظاهرِ کمتر قابل‌چاپ دارد. در صورت تردید، چیدمان‌ها را میان کاربران ترکیب کنید.

**Fallback هرگز فعال نمی‌شود.**
مطمئن شوید `suspicious_action: "fallback"` و یک `fallback_address` در دسترس تنظیم شده است.

**مستندات Sudoku کجاست؟**
مستندات رسمی sing-box پروتکل Sudoku را پوشش نمی‌دهد؛ لینک «Doc» در پنل به مخزن [sudoku-sing-box](https://github.com/SUDOKU-ASCII/sudoku-sing-box) اشاره می‌کند.

</div>
