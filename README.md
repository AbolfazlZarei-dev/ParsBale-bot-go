# ParsBale-bot-go
 # **مستندات کتابخانه ``ParsBale``**



<p align="center">
  <img src="https://s8.uupload.ir/files/parsbale_c9n2.jpg" alt="ParsBale Banner" width="30%">
</p>

<h1 align="center">ParsBale Bot API</h1>
<p align="center">
  <strong>⚡️ یک Wrapper قدرتمند، سریع و مدرن برای Bot API مسنجر بله با زبان Go (Golang)</strong>
</p>

<p align="center">
    <a href="https://goreportcard.com/report/github.com/AbolfazlZarei-dev/ParsBale-bot-go/v1"><img src="https://goreportcard.com/badge/github.com/AbolfazlZarei-dev/ParsBale-bot-go/v1" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/AbolfazlZarei-dev/ParsBale-bot-go/v1"><img src="https://godoc.org/github.com/AbolfazlZarei-dev/ParsBale-bot-go/v1?status.svg" alt="GoDoc"></a>
    <a href="https://choosealicense.com/licenses/mit/"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License"></a>
    <a href="https://github.com/AbolfazlZarei-dev/ParsBale-bot-go/issues"><img src="https://img.shields.io/github/issues/AbolfazlZarei-dev/ParsBale-bot-go.svg" alt="Issues"></a>
</p>

<p align="center">
  <a href="#ویژگی‌ها">📚 ویژگی‌ها</a> •
  <a href="#نصب">🚀 نصب</a> •
  <a href="#شروع-سریع">⚡️ شروع سریع</a> •
  <a href="#مستندات">📖 مستندات</a> •
  <a href="#مجوز">🛡️ مجوز</a>
</p>

---

## 📋 فهرست مطالب
- [✨ ویژگی‌ها](#-ویژگی‌ها)
- [🚀 نصب](#-نصب)
- [⚡️ شروع سریع](#️-شروع-سریع)
- [📖 مستندات](#-مستندات)
  - [🔀 روتینگ و دیسپچر](#-روتینگ-و-دیسپچر)
  - [🧠 مدیریت وضعیت (State/FSM)](#-مدیریت-وضعیت-statefsm)
  - [🛡️ میدل‌ور‌ها (Middlewares)](#️-میدل‌ور‌ها-middlewares)
  - [💬 ارسال پیام و مدیا](#-ارسال-پیام-و-مدیا)
  - [⌨️ کیبوردها](#️-کیبوردها)
  - [💳 پرداخت و فاکتور](#-پرداخت-و-فاکتور)
  - [📱 سرویس سفیر (Safir)](#-سرویس-سفیر-safir)
  - [🌐 وب‌هوک (Webhook)](#-وب‌هوک-webhook)
- [🤝 مشارکت](#-مشارکت)
- [🛡️ مجوز](#️-مجوز)

---

## ✨ ویژگی‌ها

| قابلیت | توضیحات |
| :--- | :--- |
| 🤖 **API کامل** | پشتیبانی کامل از تمام متدهای Bot API بله (پیام، مدیا، مدیریت گروه، پین کردن و...). |
| ⚡️ **دیسپچر پیشرفته** | سیستم هندلر و روتینگ داخلی برای مدیریت آسان‌تر آپدیت‌ها و فیلترگذاری. |
| 🧠 **مدیریت State** | مدیریت وضعیت کاربران (FSM) به صورت پیش‌فرض در حافظه با قابلیت توسعه. |
| 🔐 **میدل‌ور** | قابلیت اضافه کردن لایه‌های میانی (احراز هویت، لاگینگ) قبل از هندلرها. |
| 📱 **سرویس Safir** | ارسال پیامک OTP و آپلود فایل از طریق API سفیر بله. |
| 🌍 **WebApp** | اعتبارسنجی داده‌های Mini App (Web App). |
| 🔁 **Polling & Webhook** | پشتیبانی همزمان از Long Polling و Webhook. |

---

## 🚀 نصب

برای نصب کتابخانه کافیست دستور زیر را در ترمینال اجرا کنید:

```bash
go get github.com/AbolfazlZarei-dev/ParsBale-bot-go/v1
```

---

## ⚡️ شروع سریع

یک ربات ساده که به دستور `/start` پاسخ می‌دهد و پیام‌های متنی را اکو می‌کند:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	ParsBale "github.com/AbolfazlZarei-dev/ParsBale-bot-go/v1"
)

func main() {
	// 🤖 توکن ربات خود را اینجا قرار دهید
	bot, err := ParsBale.NewBot("YOUR_BOT_TOKEN")
	if err != nil {
		log.Fatalf("🚫 Error creating bot: %v", err)
	}

	dp := ParsBale.NewDispatcher(bot)

	// 🟢 هندلر دستور /start
	dp.OnCommand("start", func(bot *ParsBale.Bot, u ParsBale.Update) {
		chatID := u.Message.Chat.ID
		bot.SendMessage(chatID, "👋 سلام! به ربات ParsBale خوش آمدید.", nil)
	})

	// 🔄 هندلر متن (اکو)
	dp.OnText(".*", func(bot *ParsBale.Bot, u ParsBale.Update) {
		chatID := u.Message.Chat.ID
		bot.SendMessage(chatID, "📢 شما گفتید: "+u.Message.Text, nil)
	})

	// 🛑 مدیریت خاتمه برنامه
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Println("🤖 Bot is running...")
	if err := dp.StartPolling(ctx); err != nil {
		log.Fatal(err)
	}
}
```

---

## 📖 مستندات

### 🔀 روتینگ و دیسپچر
دیسپچر (`Dispatcher`) مسئول دریافت آپدیت‌ها و هدایت آن‌ها به هندلرهای مناسب است.

#### 🎯 هندلرهای آماده
کتابخانه هندلرهای از پیش تعریف شده‌ای دارد:

```go
// 🔹 دستورات
dp.OnCommand("help", handlerHelp)

// 🔹 متن با رگولار اکسپرشن
dp.OnText("^[a-zA-Z]+$", handlerOnlyLetters)

// 🔹 دکمه‌های شیشه‌ای (Callback)
dp.OnCallback("action_", handlerCallback)

// 🔹 مدیریت وضعیت
dp.OnState("waiting_email", handlerEmail)
```

#### 🛠️ هندلر سفارشی
می‌توانید با استفاده از `Handle` فیلتر کاملاً سفارشی بنویسید:

```go
dp.Handle(func(u ParsBale.Update) bool {
    // فقط پیام‌های تکست که طولشان بیشتر از ۵ حرف است
    return u.Message != nil && len(u.Message.Text) > 5
}, func(bot *ParsBale.Bot, u ParsBale.Update) {
    // logic ...
})
```

---

### 🧠 مدیریت وضعیت (State/FSM)
برای ساخت ربات‌های چندمرحله‌ای (مثل فرم ثبت‌نام)، از سیستم State استفاده کنید.

```go
// 🟢 مرحله ۱: شروع فرآیند
dp.OnCommand("register", func(bot *ParsBale.Bot, u ParsBale.Update) {
    userID := u.Message.From.ID
    dp.State.Set(userID, "waiting_name") // ست کردن وضعیت
    bot.SendMessage(u.Message.Chat.ID, "📝 نام خود را وارد کنید:", nil)
})

// 🟡 مرحله ۲: دریافت نام
dp.OnState("waiting_name", func(bot *ParsBale.Bot, u ParsBale.Update) {
    // logic save name...
    dp.State.Set(u.Message.From.ID, "waiting_phone")
    bot.SendMessage(u.Message.Chat.ID, "📞 شماره تماس خود را وارد کنید:", nil)
})

// 🔴 مرحله ۳: پایان
dp.OnState("waiting_phone", func(bot *ParsBale.Bot, u ParsBale.Update) {
    // logic...
    dp.State.Delete(u.Message.From.ID) // حذف وضعیت
    bot.SendMessage(u.Message.Chat.ID, "✅ ثبت نام با موفقیت انجام شد!", nil)
})
```

---

### 🛡️ میدل‌ور‌ها (Middlewares)
لاژیک قبل از هندلر اصلی (مثل احراز هویت یا لاگینگ).

```go
// 📊 میدل‌ور لاگر
func LoggerMiddleware(next ParsBale.HandlerFunc) ParsBale.HandlerFunc {
    return func(bot *ParsBale.Bot, u ParsBale.Update) {
        log.Printf("📩 Update received: %d", u.UpdateID)
        next(bot, u)
    }
}

// 🔐 میدل‌ور ادمین
func AdminOnlyMiddleware(next ParsBale.HandlerFunc) ParsBale.HandlerFunc {
    return func(bot *ParsBale.Bot, u ParsBale.Update) {
        if u.Message.From.ID == 12345678 {
            next(bot, u)
        } else {
            bot.SendMessage(u.Message.Chat.ID, "⛔️ دسترسی غیرمجاز!", nil)
        }
    }
}

func main() {
    dp := ParsBale.NewDispatcher(bot)
    dp.Use(LoggerMiddleware) // اعمال سراسری
}
```

---

### 💬 ارسال پیام و مدیا

#### 📝 ارسال متن
```go
bot.SendMessage(chatID, "سلام", nil)

// با فرمت Markdown
bot.SendMessage(chatID, "*متن بولد*", &ParsBale.SendMessageOptions{
    ParseMode: ParsBale.ModeMarkdown,
})
```

#### 🖼️ ارسال عکس و فایل
```go
// با File ID
bot.SendPhoto(chatID, "FILE_ID", "توضیحات", nil)

// آپلود فایل
file, _ := os.Open("photo.jpg")
bot.SendPhoto(chatID, ParsBale.FileUpload{Name: "photo.jpg", Content: file}, "عکس جدید", nil)
```

#### ✨ فرمت‌های Markdown
```go
text := ParsBale.Bold("مهم") + " و " + ParsBale.Link("گوگل", "https://google.com")
bot.SendMessage(chatID, text, &ParsBale.SendMessageOptions{ParseMode: ParsBale.ModeMarkdown})
```

---

### ⌨️ کیبوردها

#### 💎 کیبورد شیشه‌ای (Inline)
```go
keyboard := ParsBale.NewInlineKeyboard(
    []ParsBale.InlineKeyboardButton{
        {Text: "دکمه ۱", CallbackData: "btn1"},
    },
    []ParsBale.InlineKeyboardButton{
        {Text: "🌐 سایت گوگل", URL: "https://google.com"},
    },
)
bot.SendMessage(chatID, "منو:", &ParsBale.SendMessageOptions{ReplyMarkup: keyboard})
```

#### 📱 کیبورد Reply
```go
keyboard := ParsBale.NewKeyboard(
    []ParsBale.KeyboardButton{
        {Text: "📞 ارسال شماره تماس", RequestContact: true},
    },
)
bot.SendMessage(chatID, "شمارت رو بفرست:", &ParsBale.SendMessageOptions{ReplyMarkup: keyboard})
```

---

### 💳 پرداخت و فاکتور

```go
// 🧾 ارسال فاکتور
msg, _ := bot.SendInvoice(chatID, "خرید اشتراک", "پرداخت ۱۰ هزار تومانی", 
    "payload_123", "PROVIDER_TOKEN", 
    []ParsBale.LabeledPrice{{Label: "قیمت", Amount: 100000}}, "")

// ✅ استعلام تراکنش
transaction, _ := bot.InquireTransaction("transaction_id")
fmt.Println("وضعیت:", transaction.Status)
```

---

### 📱 سرویس سفیر (Safir)

#### 📨 ارسال پیامک OTP
```go
otpClient := ParsBale.NewSafirOTPClient("CLIENT_ID", "CLIENT_SECRET")
err := otpClient.SendOTP("09123456789", "12345")
```

#### 📤 ارسال پیام متنی (V3)
```go
safir := ParsBale.NewSafirClient("ACCESS_KEY")
resp, err := safir.SendMessage(botID, "09123456789", ParsBale.SafirMessageData{
    Message: &ParsBale.SafirMessage{Text: "سلام از طرف ربات"},
})
```

---

### 🌐 وب‌هوک (Webhook)

```go
func main() {
    bot, _ := ParsBale.NewBot("TOKEN")
    dp := ParsBale.NewDispatcher(bot)
    ctx := context.Background()
    
    // تنظیم وب‌هوک
    bot.SetWebhook("https://your-domain.com/bot")

    // شروع سرور
    log.Fatal(dp.StartWebhook(ctx, ":8080", "/bot"))
}
```

---

## 🤝 مشارکت
پروژه متن‌باز است. اگر مایل به بهبود آن هستید، Pull Request های شما بازرسانی می‌شود! ❤️

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 🛡️ مجوز
این پروژه تحت لایسنس MIT منتشر شده است. برای اطلاعات بیشتر فایل `LICENSE` را مطالعه کنید.

<p align="center">
  Made with ❤️ by <a href="https://github.com/AbolfazlZarei-dev">Abolfazl Zarei</a>
</p>
```
