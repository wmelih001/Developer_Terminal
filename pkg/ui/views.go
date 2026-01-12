package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Text Styles
	LabelStyle = lipgloss.NewStyle().Foreground(ColorWhite)
	ValueStyle = lipgloss.NewStyle().Foreground(ColorCyan)
	IconStyle  = lipgloss.NewStyle().Foreground(ColorYellow)

	HeaderStyle = lipgloss.NewStyle().Foreground(ColorPurple).Bold(true)
)

func (m *MainModel) actionsView() string {
	p := m.Selected

	// Helper to render checkmark
	renderCheck := func(exists bool) string {
		if exists {
			return lipgloss.NewStyle().Foreground(ColorGreen).Render("✅ VAR")
		}
		return lipgloss.NewStyle().Foreground(ColorRed).Render("❌ YOK")
	}

	// Helper to get technology icon
	getTechIcon := func(techType string) string {
		icons := map[string]string{
			"Next.js":      "⚡",
			"React":        "⚛️",
			"Vue":          "💚",
			"Vite":         "⚡",
			"React Native": "📱",
			"Mobile":       "📱",
			"HTML":         "🌐",
			"TypeScript":   "🔷",
			"Angular":      "🅰️",
			"Svelte":       "🔥",
			"SolidJS":      "💎",
			"Astro":        "🚀",
			"Remix":        "💿",
			"Nuxt":         "💚",
			"NestJS":       "🐱",
			"Express":      "🚂",
			"Go":           "🐹",
			"Django":       "🐍",
			"Flask":        "🧪",
			"Laravel":      "🐘",
			"Spring":       "☕",
			"PHP":          "🐘",
			"FastAPI":      "⚡",
			"Fiber":        "🔷",
			"Hono":         "🔥",
			"Koa":          "🥝",
			"Flutter":      "🦋",
			"Expo":         "📱",
			"Docker":       "🐳",
			"Bilinmeyen":   "📦",
		}
		if icon, ok := icons[techType]; ok {
			return icon
		}
		return "📦"
	}

	// Frontend label with technology name
	frontendLabel := "🖥️ Frontend"
	if p.FrontendType != "" && p.FrontendType != "Bilinmeyen" {
		frontendLabel = fmt.Sprintf("%s Frontend (%s)", getTechIcon(string(p.FrontendType)), p.FrontendType)
	}

	// Backend label with technology name
	backendLabel := "⚙️ Backend"
	if p.BackendType != "" && p.BackendType != "Bilinmeyen" {
		backendLabel = fmt.Sprintf("%s Backend (%s)", getTechIcon(string(p.BackendType)), p.BackendType)
	}

	// Sol Kolon (Frontend)
	frontVer := "Yok"
	if p.FrontendVer != "" {
		frontVer = p.FrontendVer
	}

	// Sağ Kolon (Backend)
	backVer := "Yok"
	if p.BackendVer != "" {
		backVer = p.BackendVer
	}

	// Table Dimensions
	const (
		totalWidth = 72
		innerW     = 70
		col1W      = 35
		col2W      = 34 // 35+1+34 = 70
	)

	// Styles for alignment
	cellStyle := lipgloss.NewStyle().Padding(0, 1).Width(col1W)  // 35 width
	cellStyleR := lipgloss.NewStyle().Padding(0, 1).Width(col2W) // 34 width
	fullRowStyle := lipgloss.NewStyle().Padding(0, 1).Width(innerW)

	borderColor := ColorGrey

	// 1. Top Border
	topBorder := lipgloss.NewStyle().Foreground(borderColor).Render("┌" + strings.Repeat("─", innerW) + "┐")

	// 2. Name Row with Health Score
	healthIcon := "🔴"
	healthColor := ColorRed
	if p.HealthScore >= 80 {
		healthIcon = "🟢"
		healthColor = ColorGreen
	} else if p.HealthScore >= 50 {
		healthIcon = "🟡"
		healthColor = ColorYellow
	}
	healthScoreStr := lipgloss.NewStyle().Foreground(healthColor).Render(fmt.Sprintf("%s %d/100", healthIcon, p.HealthScore))
	nameContent := "📂 PROJE: " + IconStyle.Render(p.Name) + "  " + healthScoreStr
	nameRowStr := lipgloss.NewStyle().Width(innerW).Padding(0, 1).Render(nameContent)
	nameRow := lipgloss.NewStyle().Foreground(borderColor).Render("│") + nameRowStr + lipgloss.NewStyle().Foreground(borderColor).Render("│")

	// 3. Separator 1 (Split)
	sep1 := lipgloss.NewStyle().Foreground(borderColor).Render("├" + strings.Repeat("─", col1W) + "┬" + strings.Repeat("─", col2W) + "┤")

	// 4. Version Row - Dynamic labels based on detected tech
	frontVerLabel := "📦 Versiyon"
	if p.FrontendType != "" && p.FrontendType != "Bilinmeyen" {
		frontVerLabel = fmt.Sprintf("%s %s", getTechIcon(string(p.FrontendType)), p.FrontendType)
	}
	backVerLabel := "📦 Versiyon"
	if p.BackendType != "" && p.BackendType != "Bilinmeyen" {
		backVerLabel = fmt.Sprintf("%s %s", getTechIcon(string(p.BackendType)), p.BackendType)
	}

	vLeftStr := cellStyle.Render(fmt.Sprintf("%s: %s", frontVerLabel, ValueStyle.Render(frontVer)))
	vRightStr := cellStyleR.Render(fmt.Sprintf("%s: %s", backVerLabel, ValueStyle.Render(backVer)))
	verRow := lipgloss.NewStyle().Foreground(borderColor).Render("│") + vLeftStr + lipgloss.NewStyle().Foreground(borderColor).Render("│") + vRightStr + lipgloss.NewStyle().Foreground(borderColor).Render("│")

	// Helper: Sürüm numarası mı yoksa sadece "Var" mı kontrol et
	hasRealVersion := func(ver string) bool {
		// "Var", "iOS", "Android", "iOS & Android" gibi değerler sürüm değil
		nonVersionValues := []string{"Var", "iOS", "Android", "iOS & Android"}
		for _, nv := range nonVersionValues {
			if ver == nv {
				return false
			}
		}
		return true
	}

	// Teknolojileri sürümü olanlar ve olmayanlar olarak ayır
	var frontendWithVersion, frontendWithoutVersion []struct {
		Type    string
		Version string
	}
	var backendWithVersion, backendWithoutVersion []struct {
		Type    string
		Version string
	}

	for _, ft := range p.DetectedFrontendTechs {
		tech := struct {
			Type    string
			Version string
		}{string(ft.Type), ft.Version}
		if hasRealVersion(ft.Version) {
			frontendWithVersion = append(frontendWithVersion, tech)
		} else {
			frontendWithoutVersion = append(frontendWithoutVersion, tech)
		}
	}

	for _, bt := range p.DetectedBackendTechs {
		tech := struct {
			Type    string
			Version string
		}{string(bt.Type), bt.Version}
		if hasRealVersion(bt.Version) {
			backendWithVersion = append(backendWithVersion, tech)
		} else {
			backendWithoutVersion = append(backendWithoutVersion, tech)
		}
	}

	// 5. Versioned tech rows (sürümü olanlar - versiyon satırının altına)
	var versionedTechRows []string
	maxVersionedRows := len(frontendWithVersion)
	if len(backendWithVersion) > maxVersionedRows {
		maxVersionedRows = len(backendWithVersion)
	}

	for i := 0; i < maxVersionedRows; i++ {
		var frontTechStr, backTechStr string
		if i < len(frontendWithVersion) {
			ft := frontendWithVersion[i]
			frontTechStr = fmt.Sprintf("  %s %s: %s", getTechIcon(ft.Type), ft.Type, ValueStyle.Render(ft.Version))
		}
		if i < len(backendWithVersion) {
			bt := backendWithVersion[i]
			backTechStr = fmt.Sprintf("  %s %s: %s", getTechIcon(bt.Type), bt.Type, ValueStyle.Render(bt.Version))
		}
		leftCell := cellStyle.Render(frontTechStr)
		rightCell := cellStyleR.Render(backTechStr)
		row := lipgloss.NewStyle().Foreground(borderColor).Render("│") + leftCell + lipgloss.NewStyle().Foreground(borderColor).Render("│") + rightCell + lipgloss.NewStyle().Foreground(borderColor).Render("│")
		versionedTechRows = append(versionedTechRows, row)
	}

	// 6. Separator 2 (Cross)
	sep2 := lipgloss.NewStyle().Foreground(borderColor).Render("├" + strings.Repeat("─", col1W) + "┼" + strings.Repeat("─", col2W) + "┤")

	// 7. Status Row - Dynamic labels
	sLeftStr := cellStyle.Render(fmt.Sprintf("%s: %s", frontendLabel, renderCheck(p.HasFrontend)))
	sRightStr := cellStyleR.Render(fmt.Sprintf("%s: %s", backendLabel, renderCheck(p.HasBackend)))
	statRow := lipgloss.NewStyle().Foreground(borderColor).Render("│") + sLeftStr + lipgloss.NewStyle().Foreground(borderColor).Render("│") + sRightStr + lipgloss.NewStyle().Foreground(borderColor).Render("│")

	// 8. Non-versioned tech rows (sadece "Var" olanlar - status satırının altına)
	var nonVersionedTechRows []string
	maxNonVersionedRows := len(frontendWithoutVersion)
	if len(backendWithoutVersion) > maxNonVersionedRows {
		maxNonVersionedRows = len(backendWithoutVersion)
	}

	for i := 0; i < maxNonVersionedRows; i++ {
		var frontTechStr, backTechStr string
		if i < len(frontendWithoutVersion) {
			ft := frontendWithoutVersion[i]
			frontTechStr = fmt.Sprintf("  %s %s: %s", getTechIcon(ft.Type), ft.Type, lipgloss.NewStyle().Foreground(ColorGreen).Render("✅ VAR"))
		}
		if i < len(backendWithoutVersion) {
			bt := backendWithoutVersion[i]
			backTechStr = fmt.Sprintf("  %s %s: %s", getTechIcon(bt.Type), bt.Type, lipgloss.NewStyle().Foreground(ColorGreen).Render("✅ VAR"))
		}
		leftCell := cellStyle.Render(frontTechStr)
		rightCell := cellStyleR.Render(backTechStr)
		row := lipgloss.NewStyle().Foreground(borderColor).Render("│") + leftCell + lipgloss.NewStyle().Foreground(borderColor).Render("│") + rightCell + lipgloss.NewStyle().Foreground(borderColor).Render("│")
		nonVersionedTechRows = append(nonVersionedTechRows, row)
	}

	// 9. Docker Row (if exists)
	var dockerRow string
	var sep3 string
	if p.HasDocker {
		sep3 = lipgloss.NewStyle().Foreground(borderColor).Render("├" + strings.Repeat("─", innerW) + "┤")
		dockerContent := fullRowStyle.Render(fmt.Sprintf("🐳 Docker: %s", lipgloss.NewStyle().Foreground(ColorGreen).Render("✅ VAR")))
		dockerRow = lipgloss.NewStyle().Foreground(borderColor).Render("│") + dockerContent + lipgloss.NewStyle().Foreground(borderColor).Render("│")
	}

	// 10. Monorepo alt projeleri (varsa)
	var monorepoRows []string
	var sep4 string
	if p.IsMonorepo && (len(p.AllFrontends) > 1 || len(p.AllBackends) > 1) {
		sep4 = lipgloss.NewStyle().Foreground(borderColor).Render("├" + strings.Repeat("─", innerW) + "┤")

		// Başlık
		monorepoHeader := fullRowStyle.Render(lipgloss.NewStyle().Foreground(ColorPurple).Bold(true).Render("📦 MONOREPO ALT PROJELERİ"))
		monorepoRows = append(monorepoRows, lipgloss.NewStyle().Foreground(borderColor).Render("│")+monorepoHeader+lipgloss.NewStyle().Foreground(borderColor).Render("│"))

		// Frontend alt projeleri
		for i, sub := range p.AllFrontends {
			prefix := "  "
			if i == 0 {
				prefix = "→ " // Ana proje
			}
			subStr := fullRowStyle.Render(fmt.Sprintf("%s%s %s: %s", prefix, getTechIcon(string(sub.Type)), sub.Name, ValueStyle.Render(sub.Version)))
			monorepoRows = append(monorepoRows, lipgloss.NewStyle().Foreground(borderColor).Render("│")+subStr+lipgloss.NewStyle().Foreground(borderColor).Render("│"))
		}

		// Backend alt projeleri
		for i, sub := range p.AllBackends {
			prefix := "  "
			if i == 0 {
				prefix = "→ " // Ana proje
			}
			subStr := fullRowStyle.Render(fmt.Sprintf("%s%s %s: %s", prefix, getTechIcon(string(sub.Type)), sub.Name, ValueStyle.Render(sub.Version)))
			monorepoRows = append(monorepoRows, lipgloss.NewStyle().Foreground(borderColor).Render("│")+subStr+lipgloss.NewStyle().Foreground(borderColor).Render("│"))
		}
	}

	// 11. Bottom Border
	var botBorder string
	if p.HasDocker || len(monorepoRows) > 0 {
		botBorder = lipgloss.NewStyle().Foreground(borderColor).Render("└" + strings.Repeat("─", innerW) + "┘")
	} else {
		botBorder = lipgloss.NewStyle().Foreground(borderColor).Render("└" + strings.Repeat("─", col1W) + "┴" + strings.Repeat("─", col2W) + "┘")
	}

	// Assemble
	var boxParts []string
	boxParts = append(boxParts, topBorder, nameRow, sep1, verRow)
	// Sürümü olan teknolojiler (versiyon satırı altına)
	for _, row := range versionedTechRows {
		boxParts = append(boxParts, row)
	}
	boxParts = append(boxParts, sep2, statRow)
	// Sürümü olmayan teknolojiler (status satırı altına)
	for _, row := range nonVersionedTechRows {
		boxParts = append(boxParts, row)
	}
	if p.HasDocker {
		boxParts = append(boxParts, sep3, dockerRow)
	}
	// Port uyarıları
	if len(p.PortWarnings) > 0 {
		sep5 := lipgloss.NewStyle().Foreground(borderColor).Render("├" + strings.Repeat("─", innerW) + "┤")
		boxParts = append(boxParts, sep5)
		for _, warning := range p.PortWarnings {
			warningContent := fullRowStyle.Render(lipgloss.NewStyle().Foreground(ColorYellow).Render(warning))
			warningRow := lipgloss.NewStyle().Foreground(borderColor).Render("│") + warningContent + lipgloss.NewStyle().Foreground(borderColor).Render("│")
			boxParts = append(boxParts, warningRow)
		}
	}
	// Monorepo alt projeleri
	if len(monorepoRows) > 0 {
		boxParts = append(boxParts, sep4)
		boxParts = append(boxParts, monorepoRows...)
	}
	boxParts = append(boxParts, botBorder)

	finalBox := lipgloss.JoinVertical(lipgloss.Left, boxParts...)

	// --- SEÇENEKLER ---
	var b strings.Builder

	b.WriteString("\n" + finalBox + "\n\n")

	// 1. Başlatma
	b.WriteString(HeaderStyle.Render("🚀 BAŞLATMA SEÇENEKLERİ") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(ColorGrey).Render("───────────────────────") + "\n")
	b.WriteString("[1] 🖥️  Sadece Frontend\n")
	b.WriteString("[2] ⚙️  Sadece Backend\n")
	b.WriteString("[3] 🔥  Full Stack (İkisi Bir Arada)\n\n")

	// 2. Uzak Erişim
	b.WriteString(HeaderStyle.Render("🌍 UZAK ERİŞİM") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(ColorGrey).Render("───────────────────────") + "\n")
	b.WriteString("[4] 📡  Canlı Bağlantı (Ngrok Public)\n\n")

	// 3. Genel Araçlar
	b.WriteString(HeaderStyle.Render("🛠️ GENEL") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(ColorGrey).Render("───────────────────────") + "\n")

	if m.CopiedSuccess {
		b.WriteString(lipgloss.NewStyle().Foreground(ColorGreen).Render("[5] ✅  Kopyalandı! (Panoya Hazır)") + "\n")
	} else {
		b.WriteString("[5] 🧬  AI Context (Ağacı Kopyala)\n")
	}

	b.WriteString("[6] 🩺  Dependency Doctor (Paket Güncelle)\n")
	b.WriteString("[H] 🏥  Sağlık Skoru Hesapla\n")
	b.WriteString("[E] 📂  Explorer'da Aç\n")

	// 3.5. Task Runner (Scriptler varsa)
	if len(m.Selected.Scripts) > 0 {
		b.WriteString("[7] 📜  Script Çalıştır (Task Runner)\n")
	}
	b.WriteString("\n")

	// 4. Veritabanı Araçları (sadece varsa göster)
	hasDbTools := m.Selected.HasPrisma || m.Selected.HasDrizzle || m.Selected.HasHasura || m.Selected.HasSupabase || m.Selected.HasStorybook
	if hasDbTools {
		b.WriteString(HeaderStyle.Render("🧠 VERİTABANI & UI ARAÇLARI") + "\n")
		b.WriteString(lipgloss.NewStyle().Foreground(ColorGrey).Render("───────────────────────") + "\n")
		if m.Selected.HasPrisma {
			b.WriteString("[F1] ◮  Prisma Studio\n")
		}
		if m.Selected.HasDrizzle {
			b.WriteString("[F2] 🌧️  Drizzle Studio\n")
		}
		if m.Selected.HasHasura {
			b.WriteString("[F3] 🦅  Hasura Console\n")
		}
		if m.Selected.HasSupabase {
			b.WriteString("[F4] ⚡  Supabase Status\n")
		}
		if m.Selected.HasStorybook {
			b.WriteString("[F5] 📕  Storybook (UI Dev)\n")
		}
		b.WriteString("\n")
	}

	// Seçenekleri bitir ve input satırını ekle

	// Apply global left padding to main content
	content := lipgloss.NewStyle().PaddingLeft(2).Render(b.String())

	// Footer Oluştur
	footer := m.renderFooter("Esc", "Geri Dön")

	// Sticky Footer Logic (En alta it)
	if footer != "" {
		hContent := lipgloss.Height(content)
		hFooter := lipgloss.Height(footer)
		gap := m.Height - hContent - hFooter - 1 // -1 safety
		if gap > 0 {
			content += strings.Repeat("\n", gap)
		} else {
			content += "\n\n"
		}
		// Footer padding must match content padding (2 spaces)
		content += "  " + footer
	}

	return content
}

func (m *MainModel) dashboardView() string {
	return ""
}

// renderFooter generates a standardized footer with custom keys + global keys (q, ,)
func (m *MainModel) renderFooter(pairs ...string) string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#909090"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	dot := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(" • ")

	qKey := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render("q")
	qDesc := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render("Çıkış")

	var parts []string

	// 1. Custom Keys
	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			k, d := pairs[i], pairs[i+1]
			parts = append(parts, keyStyle.Render(k)+" "+descStyle.Render(d))
		}
	}

	// 2. Global Keys
	parts = append(parts, qKey+" "+qDesc)
	parts = append(parts, keyStyle.Render(",")+" "+descStyle.Render("Daha Fazla"))

	return strings.Join(parts, dot)
}

func (m *MainModel) taskRunnerView() string {
	doc := strings.Builder{}

	doc.WriteString("\n")

	// Helper logic to style the list
	// The list component handles its own rendering
	listView := m.TaskRunnerList.View()
	listView = strings.Replace(listView, "filtered", "sonuç", 1) // Hacky localization
	listView = strings.Replace(listView, "Nothing matched", "Sonuç bulunamadı", 1)
	doc.WriteString(listView)

	return doc.String()
}

func (m *MainModel) splashView() string {
	art := `
  ____                 _                         _____                    _             _
 |  _ \  _____   _____| | ___  _ __   ___ _ __  |_   _|__ _ __ _ __ ___ (_)_ __   __ _| |
 | | | |/ _ \ \ / / _ \ |/ _ \| '_ \ / _ \ '__|   | |/ _ \ '__| '_ \ _ \| | '_ \ / _\ | |
 | |_| |  __/\ V /  __/ | (_) | |_) |  __/ |      | |  __/ |  | | | | | | | | | | (_| | |
 |____/ \___| \_/ \___|_|\___/| .__/ \___|_|      |_|\___|_|  |_| |_| |_|_|_| |_|\__,_|_|
                              |_|
`
	// 1. Solid Color Logo (Cool Dark Purple/Blue)
	// Havalı koyu stil: #bd93f9 (Dracula Purple) veya #6272a4 (Comment Blue/Gray)
	// Kullanıcı "Havalı koyu bir renk" dedi.
	styledLogo := lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Bold(true).Render(art)

	version := lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4")).Italic(true).Render("Developer Terminal v1.0.5")

	// 2. Dynamic Progress Bar
	width := 40
	completed := int(float64(width) * m.SplashProgress)
	if completed > width {
		completed = width
	}
	remaining := width - completed
	if remaining < 0 {
		remaining = 0
	}

	// Bar Gradient Color
	var barColor lipgloss.Color
	if m.SplashProgress < 0.3 {
		barColor = lipgloss.Color("#ff5555") // Red
	} else if m.SplashProgress < 0.7 {
		barColor = lipgloss.Color("#f1fa8c") // Yellow
	} else {
		barColor = lipgloss.Color("#50fa7b") // Green
	}

	barStyle := lipgloss.NewStyle().Foreground(barColor)
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#44475a"))

	barStr := barStyle.Render(strings.Repeat("█", completed)) + emptyStyle.Render(strings.Repeat("░", remaining))
	percent := int(m.SplashProgress * 100)

	// 3. Random Loading Messages
	messages := []string{
		"Kuantum evreni taranıyor...",
		"Kahve hazırlanıyor...",
		"Matrix'e bağlanılıyor...",
		"Node_modules ağırlığı hesaplanıyor...",
		"Yapay zeka motoru ısıtılıyor...",
		"Sistem kaynakları optimize ediliyor...",
		"Geliştirici modu etkinleştiriliyor...",
	}
	// Pick message based on progress to cycle through them
	// Show all 7 messages evenly distributed over the 6 seconds
	msgIndex := int(m.SplashProgress * float64(len(messages)))
	if msgIndex >= len(messages) {
		msgIndex = len(messages) - 1
	}
	loadingMsg := messages[msgIndex]

	// 4. Layout Assembly
	content := lipgloss.JoinVertical(lipgloss.Center,
		styledLogo,
		version,
		"",
		"",
		barStr,
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Render(fmt.Sprintf("%s (%d%%)", loadingMsg, percent)),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4")).Faint(true).Render("Atlamak için 'Space' veya 'Enter'"),
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}
