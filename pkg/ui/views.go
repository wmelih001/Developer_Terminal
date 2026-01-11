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

	borderColor := ColorGrey

	// 1. Top Border
	// ┌──────────────────────────────────────────────────────────────────────┐
	topBorder := lipgloss.NewStyle().Foreground(borderColor).Render("┌" + strings.Repeat("─", innerW) + "┐")

	// 2. Name Row
	// │ PROJE: ...                                                           │
	nameContent := "📂 PROJE: " + IconStyle.Render(p.Name)
	nameRowStr := lipgloss.NewStyle().Width(innerW).Padding(0, 1).Render(nameContent)
	nameRow := lipgloss.NewStyle().Foreground(borderColor).Render("│") + nameRowStr + lipgloss.NewStyle().Foreground(borderColor).Render("│")

	// 3. Separator 1 (Split)
	// ├───────────────────────────────────┬──────────────────────────────────┤
	// 35 dashes + 1 (┬) + 34 dashes
	sep1 := lipgloss.NewStyle().Foreground(borderColor).Render("├" + strings.Repeat("─", col1W) + "┬" + strings.Repeat("─", col2W) + "┤")

	// 4. Version Row
	// │ Next.js: ...                      │ Nest.js: ...                     │
	vLeftStr := cellStyle.Render(fmt.Sprintf("📦 Next.js: %s", ValueStyle.Render(frontVer)))
	vRightStr := cellStyleR.Render(fmt.Sprintf("📦 Nest.js: %s", ValueStyle.Render(backVer)))
	verRow := lipgloss.NewStyle().Foreground(borderColor).Render("│") + vLeftStr + lipgloss.NewStyle().Foreground(borderColor).Render("│") + vRightStr + lipgloss.NewStyle().Foreground(borderColor).Render("│")

	// 5. Separator 2 (Cross)
	// ├───────────────────────────────────┼──────────────────────────────────┤
	sep2 := lipgloss.NewStyle().Foreground(borderColor).Render("├" + strings.Repeat("─", col1W) + "┼" + strings.Repeat("─", col2W) + "┤")

	// 6. Status Row
	// │ Frontend: VAR                     │ Backend: VAR                     │
	sLeftStr := cellStyle.Render(fmt.Sprintf("️🖥️ Frontend: %s", renderCheck(p.HasFrontend)))
	sRightStr := cellStyleR.Render(fmt.Sprintf("⚙️ Backend: %s", renderCheck(p.HasBackend)))
	statRow := lipgloss.NewStyle().Foreground(borderColor).Render("│") + sLeftStr + lipgloss.NewStyle().Foreground(borderColor).Render("│") + sRightStr + lipgloss.NewStyle().Foreground(borderColor).Render("│")

	// 7. Bottom Border (Join)
	// └───────────────────────────────────┴──────────────────────────────────┘
	botBorder := lipgloss.NewStyle().Foreground(borderColor).Render("└" + strings.Repeat("─", col1W) + "┴" + strings.Repeat("─", col2W) + "┘")

	// Assemble
	finalBox := lipgloss.JoinVertical(lipgloss.Left,
		topBorder,
		nameRow,
		sep1,
		verRow,
		sep2,
		statRow,
		botBorder,
	)

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

	// 3. Yapay Zeka
	b.WriteString(HeaderStyle.Render("🧠 ARAÇLAR") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(ColorGrey).Render("───────────────────────") + "\n")

	if m.CopiedSuccess {
		b.WriteString(lipgloss.NewStyle().Foreground(ColorGreen).Render("[5] ✅  Kopyalandı! (Panoya Hazır)") + "\n")
	} else {
		b.WriteString("[5] 🧬  AI Context (Ağacı Kopyala)\n")
	}

	b.WriteString("[6] 🩺  Dependency Doctor (Paket Güncelle)\n\n")

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
