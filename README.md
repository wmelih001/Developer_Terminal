# DEVELOPER TERMINAL

![Go Sürümü](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)
![Lisans](https://img.shields.io/badge/license-MIT-blue?style=flat-square)

Developer Terminal, eski PowerShell profil betiklerinin yerini almak üzere tasarlanmış, yüksek performanslı bir Geliştirici Kontrol Paneli ve CLI aracıdır. Go ve Bubble Tea framework'ü ile geliştirilen bu araç; proje yönetimi, yapay zeka bağlamı (context) oluşturma ve geliştirme ortamını izleme işlemleri için modern, klavye odaklı bir Terminal Kullanıcı Arayüzü (TUI) sunar.

## Özellikler

### 🚀 Proje Başlatıcı
Çalışma alanınızdaki dizinleri anında tarayın ve projelerinizi tek bir tuşla Windows Terminal'de başlatın. Proje türlerini otomatik olarak algılar ve özel başlatma modları sunar:
- **Frontend**: Proje dizinindeki `package.json` dosyasını analiz eder ve uygun başlatma komutunu (örn. `npm run dev`) otomatik olarak belirleyip yeni sekmede çalıştırır.
- **Backend**: Backend projesinin türünü (Go, NestJS vb.) algılar ve ilgili çalıştırma komutunu (örn. `go run .` veya `npm run start:dev`) yeni sekmede başlatır.
- **Full Stack**: Terminal penceresini ikiye bölerek her ikisini aynı anda çalıştırır.

### 📜 Gelişmiş Task Runner (Script Yöneticisi)
Proje kök dizinindeki veya `frontend/backend` alt klasörlerindeki `package.json` dosyalarını otomatik olarak tarar ve `scripts` komutlarını listeler.
- **Akıllı Çalıştırma:** Script ismine göre (`client:` veya `server:`) doğru çalışma dizinini (working directory) otomatik belirler ve komutu orada çalıştırır.
- **Hızlı Arama:** Binlerce script arasında kaybolmayın. `Tab` tuşu ile arama modunu açın ve istediğiniz komutu anında bulun.
- **Entegre Deneyim:** TUI'den ayrılmadan test, build, lint veya deploy işlemlerinizi tek tuşla başlatın.

### 🧠 Yapay Zeka Bağlam Oluşturucu
Büyük Dil Modelleri (LLM) için derinlemesine ve yapısal bağlamlar oluşturun. Bir proje seçin ve `.gitignore` kurallarına sadık kalarak kod tabanınızın temiz bir ASCII ağaç yapısını üretin. Çıktı anında panoya kopyalanır, yapay zekaya prompt girmek için hazırdır.

### 🩺 Bağımlılık Doktoru
Projelerinizi sağlıklı tutun. Developer Terminal, `package.json` dosyalarını analiz ederek temel framework'lerin (React, Next.js, NestJS) güncel sürümlerini görüntüler ve terminalden çıkmadan güncelliğini yitirmiş bağımlılıkları kontrol eder.

### 🏥 Proje Sağlık Skoru
Projelerinizin kalitesini ve standartlara uygunluğunu anlık olarak ölçün. "Sağlık Skoru Hesapla" özelliği, projenizi derinlemesine tarayarak (recursive) 100 üzerinden puanlar:
- **Kriterler:** Git durumu, README varlığı, Lisans dosyası, CI/CD yapılandırması, Docker kullanımı, Linter ayarları ve Env dosyaları.
- **Detaylı Rapor:** Eksik olan öğeleri ve puan kayıplarını listeleyerek iyileştirme önerileri sunar.
- **Akıllı Tarama:** Alt klasörlerdeki (`backend/schema.prisma` gibi) yapılandırmaları bile tespit eder.

### 🛠️ Geliştirici Araçları (Dev Tools)
Proje klasörlerinizde gömülü olan veritabanı ve UI araçlarını otomatik algılar ve tek tuşla başlatır. Komutlar, aracın bulunduğu alt klasörde (örn: `backend/`) otomatik olarak çalıştırılır:
- **[F1] Prisma Studio**: Prisma veritabanı yönetim panelini açar.
- **[F2] Drizzle Studio**: Drizzle ORM stüdyosunu başlatır.
- **[F3] Hasura Console**: Hasura GraphQL konsolunu açar.
- **[F4] Supabase Status**: Yerel Supabase durumunu görüntüler.
- **[F5] Storybook**: UI bileşen geliştirme ortamını başlatır.

### 🛡️ Port Çakışma Tespiti
Projeleri başlatmadan önce, gerekli portların (örn: 3000, 8080) dolu olup olmadığını kontrol eder. Çakışma varsa sizi uyararak "bind address already in use" hatalarının önüne geçer.

### 🚇 Ngrok Entegrasyonu
Yerel tünellerinizi doğrudan kontrol panelinden yönetin. Ngrok yolunuzu yapılandırın ve aktif tünel durumunu zahmetsizce görüntüleyin.

### 🎨 Modern TUI
Siberpunk esintili estetiğe sahip birinci sınıf bir geliştirici aracı deneyimi yaşayın.
- **Klavye Öncelikli**: Vim tarzı gezinme desteği.
- **Duyarlı (Responsive)**: Terminal yeniden boyutlandırma olaylarına dinamik olarak uyum sağlar.
- **Hızlı**: Anında açılış için tek bir yerel (native) binary olarak derlenmiştir.

### ✨ Akışkan Animasyonlar
Kullanıcı deneyimini en üst düzeye çıkaran görsel detaylar:
- **Sinematik Açılış:** Özel tasarım ASCII logo ve "cool dark" renk paleti ile profesyonel karşılama ekranı.
- **Dinamik Yükleme:** İşlem durumuna göre renk değiştiren progress bar ve sürekli güncellenen esprili yükleme mesajları ("Kuantum evreni taranıyor..." vb.).
- **Yumuşak Geçişler:** Liste ve menü geçişlerinde göz yormayan akıcı animasyonlar.

## Teknoloji Yığını

- **Çekirdek**: Go (Golang) 1.21+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Stil**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Yapılandırma**: [Viper](https://github.com/spf13/viper)

## Kurulum

### Gereksinimler
- Go 1.21 veya üzeri
- Windows Powershell 7+ (Tam başlatıcı entegrasyonu için önerilir)
- Nerd Fonts (Simgelerin görünmesi için önerilir)

### Kaynaktan Kurulum

```bash
git clone https://github.com/kullaniciadi/developer_terminal.git
cd devterminal
go install
```

## Yapılandırma

Developer Terminal, `~/.devterminal/config.yaml` konumunda bulunan bir YAML yapılandırma dosyası kullanır.

> **Not:** Windows'ta tam yol genellikle şöyledir: `C:\Users\KullaniciAdi\.devterminal\config.yaml`

Uygulamayı ilk kez çalıştırdığınızda, yapılandırma dosyası **otomatik olarak oluşturulur** ve size proje klasörlerinizin yolunu sorar. Manuel olarak oluşturmanıza gerek yoktur.

Örnek yapılandırma:

```yaml
# Proje klasörlerinin yolu
projects_paths:
  - M:\Projeler

# Tarama sırasında yok sayılacak klasörler
ignored_files:
  - .git
  - node_modules
  - dist
  - .next
  - .idea
  - .vscode

# Ngrok yolu (opsiyonel)
ngrok_path: C:\Users\KullaniciAdi\AppData\Local\Microsoft\WinGet\Links\ngrok.exe

# Başlatma komutları (Windows Terminal)
commands:
  launch_frontend: wt.exe -w 0 new-tab -d "{{.FrontendPath}}" cmd /k "{{.FrontendCmd}}"
  launch_backend: wt.exe -w 0 new-tab -d "{{.BackendPath}}" cmd /k "{{.BackendCmd}}"
  launch_full: wt.exe -w 0 new-tab -d "{{.FrontendPath}}" cmd /k "{{.FrontendCmd}}" ; split-pane -d "{{.BackendPath}}" cmd /k "{{.BackendCmd}}"

# Proje bazlı komut özelleştirmeleri (otomatik oluşturulur)
project_overrides:
  m:\projeler\my-nextjs-app:
    frontend: npm run dev
    backend: npm run start:dev
  m:\projeler\go-api:
    frontend: ""
    backend: go run .

# Son açılan projeler (otomatik oluşturulur)
last_opened:
  m:\projeler\my-project: 2026-01-12T19:00:00+03:00
```

## Lisans

MIT
