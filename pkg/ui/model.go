package ui

import (
	"devterminal/pkg/config"
	"devterminal/pkg/domain"
	"devterminal/pkg/service"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SessionState defines the current view
type SessionState int

const (
	StateDashboard SessionState = iota
	StateProjectSelect
	StateProjectActions // Submenu for launch options
	StateScanning
	StateFirstRun
	StateContextGen
	StateDependencyDoctor
	StateNgrok
)

type NgrokStep int

const (
	NgrokMainMenu NgrokStep = iota // New Start Screen
	NgrokCheckInstall
	NgrokModeSelect // Ask user if installed or not
	NgrokManualPath // New: Input path
	NgrokCheckAuth
	NgrokAuth
	NgrokAskPort
	NgrokRunning
)

type MainModel struct {
	Config *domain.Config
	// Services
	Scanner  *service.Scanner
	TreeGen  *service.TreeGenerator
	Launcher *service.Launcher
	Doctor   *service.Doctor
	State    SessionState

	// Components
	List          list.Model
	Table         table.Model
	Spinner       spinner.Model
	FirstRunInput textinput.Model

	// Ngrok
	NgrokService    *service.NgrokService
	NgrokStep       NgrokStep
	NgrokPathInput  textinput.Model
	NgrokPortInput  textinput.Model
	NgrokTokenInput textinput.Model
	NgrokCmd        string // final command to run

	// Data
	Projects []domain.Project
	Selected *domain.Project

	// Error handling
	Err    error
	Width  int
	Height int

	// Feedback Flags
	CopiedSuccess       bool
	AllPackagesUpToDate bool
}

type copiedResetMsg struct{}

func NewMainModel() *MainModel {
	cfg, err := config.LoadConfig() // Should inject this properly in main
	if err != nil {
		// handle fatal error or use empty config
		cfg = &domain.Config{}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = TitleStyle // Reuse style

	tiPort := textinput.New()
	tiPort.Placeholder = "3000"
	tiPort.SetValue("3000")
	tiPort.CharLimit = 5
	tiPort.Width = 10

	tiToken := textinput.New()
	tiToken.Placeholder = "Authtoken"
	tiToken.CharLimit = 100
	tiToken.Width = 40
	tiToken.EchoMode = textinput.EchoPassword

	tiPath := textinput.New()
	tiPath.Placeholder = "C:\\path\\to\\ngrok.exe"
	tiPath.Width = 60

	// First Run Input
	tiFirstRun := textinput.New()
	tiFirstRun.Placeholder = "C:\\Projelerim (Tırnak işaretleri temizlenir)"
	tiFirstRun.Width = 60
	tiFirstRun.Focus()

	initialState := StateScanning
	if len(cfg.ProjectsPaths) == 0 {
		initialState = StateFirstRun
	}

	return &MainModel{
		Config:          cfg,
		Scanner:         service.NewScanner(cfg),
		TreeGen:         service.NewTreeGenerator(cfg),
		Launcher:        service.NewLauncher(cfg),
		Doctor:          service.NewDoctor(cfg),
		NgrokService:    service.NewNgrokService(cfg),
		NgrokPathInput:  tiPath,
		NgrokPortInput:  tiPort,
		NgrokTokenInput: tiToken,
		FirstRunInput:   tiFirstRun,
		State:           initialState,
		Spinner:         s,
		List:            list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		Table:           newTable(),
	}
}

func newTable() table.Model {
	columns := []table.Column{
		{Title: "Paket", Width: 20},
		{Title: "Mevcut", Width: 10},
		{Title: "İstenen", Width: 10},
		{Title: "Son", Width: 10},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return t
}

func (m *MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, spinner.Tick)

	if m.State == StateScanning {
		cmds = append(cmds, m.scanProjectsCmd())
	}
	return tea.Batch(cmds...)
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Context Global Back
		if msg.String() == "esc" {
			if m.State == StateProjectActions {
				m.State = StateProjectSelect // Back to list
				return m, nil
			}

		}

		// State Specific Handling
		switch m.State {
		case StateFirstRun:
			var cmd tea.Cmd
			m.FirstRunInput, cmd = m.FirstRunInput.Update(msg)
			if msg.String() == "enter" {
				path := strings.Trim(m.FirstRunInput.Value(), "\"")
				if path != "" {
					m.Config.ProjectsPaths = []string{path}
					// Proje yolunu config'e kaydet
					if err := config.SaveConfig(m.Config); err != nil {
						m.Err = err
					} else {
						m.State = StateScanning
						cmds = append(cmds, m.scanProjectsCmd())
					}
					return m, tea.Batch(cmds...)
				}
			}
			return m, cmd

		case StateDashboard:
			switch msg.String() {
			case "1":
				m.State = StateProjectSelect
				m.List.Title = "Proje Seçiniz"
				cmds = append(cmds, m.List.StartSpinner())
			case "2":
				m.State = StateProjectSelect
				m.List.Title = "Bağımlılık Kontrolü İçin Proje Seç"
				// Seçimden sonra StateDependencyDoctor'a geçiş yukarıdaki logic'te
			case "3":
				// Start Ngrok Flow
				m.State = StateNgrok
				m.NgrokStep = NgrokMainMenu
				return m, nil
			case "q":
				return m, tea.Quit
			}

		case StateDependencyDoctor:
			if msg.String() == "esc" {
				m.State = StateProjectActions
				return m, nil
			}
			if msg.String() == "q" {
				return m, tea.Quit
			}

		case StateNgrok:
			switch m.NgrokStep {
			case NgrokMainMenu:
				switch msg.String() {
				case "1":
					m.NgrokStep = NgrokManualPath
					m.NgrokPathInput.Focus()
				case "2":
					// Try everything until something works
					url := "https://ngrok.com/download/windows"
					go func() {
						// Strategy 1: Explorer
						if err := exec.Command("explorer", url).Start(); err == nil {
							return
						}
						// Strategy 2: Rundll32
						if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start(); err == nil {
							return
						}
						// Strategy 3: Cmd Start (Standard)
						if err := exec.Command("cmd", "/c", "start", "", url).Start(); err == nil {
							return
						}
						if err := exec.Command("powershell", "-c", "Start-Process", fmt.Sprintf("'%s'", url)).Start(); err == nil {
							return
						}
						// Strategy 4: PowerShell
						_ = exec.Command("pwsh", "-c", "Start-Process", fmt.Sprintf("'%s'", url)).Start()
					}()

					m.NgrokCmd = "downloading"
					m.NgrokStep = NgrokManualPath
					m.NgrokPathInput.Focus()
				case "esc":
					m.State = StateProjectActions // Return to Project Menu
				}

			case NgrokCheckInstall:
				// Handled by Cmd
			case NgrokModeSelect:
				// Removed/Unused but keeping safe
				if msg.String() == "esc" {
					m.State = StateProjectActions
				}
			case NgrokManualPath:
				var cmd tea.Cmd
				m.NgrokPathInput, cmd = m.NgrokPathInput.Update(msg)
				if msg.String() == "enter" {
					path := strings.Trim(m.NgrokPathInput.Value(), "\"")
					if m.NgrokService.ValidatePath(path) {
						m.NgrokService.SavePath(path)
						m.NgrokStep = NgrokAuth   // Skip check, go straight to Auth
						m.NgrokTokenInput.Focus() // FIX: Focus the input!
						return m, nil
					}
					// Invalid? Shake or show error (for now just stay)
				}
				if msg.String() == "esc" {
					m.NgrokStep = NgrokMainMenu
					m.NgrokPathInput.Blur()
				}
				return m, cmd

			case NgrokCheckAuth:
				// Removed/Skipped
				m.NgrokStep = NgrokAuth
				return m, nil

			case NgrokAuth:
				var cmd tea.Cmd
				m.NgrokTokenInput, cmd = m.NgrokTokenInput.Update(msg)

				switch msg.String() {
				case "enter":
					token := m.NgrokTokenInput.Value()
					if token != "" {
						_ = m.NgrokService.SetAuthToken(token)
						m.NgrokStep = NgrokAskPort
						m.NgrokPortInput.Focus()
						return m, nil
					}
				case "ctrl+o":
					// Open Dashboard
					url := "https://dashboard.ngrok.com/get-started/your-authtoken"
					go func() {
						_ = exec.Command("cmd", "/c", "start", "", url).Start()
					}()
					return m, nil
				case "esc":
					m.NgrokStep = NgrokMainMenu
					m.NgrokTokenInput.Blur()
				}
				return m, cmd
			case NgrokAskPort:
				var cmd tea.Cmd
				m.NgrokPortInput, cmd = m.NgrokPortInput.Update(msg)
				if msg.String() == "enter" {
					m.NgrokStep = NgrokRunning
					port := m.NgrokPortInput.Value()
					exe := m.NgrokService.GetExecutable() // Use resolved path
					return m, func() tea.Msg {
						// We need to quote the exe path in case of spaces
						// wt command: wt new-tab ... cmd /k "path/to/ngrok" http 3000
						c := exec.Command("wt.exe", "new-tab", "-d", ".", "--title", "Ngrok "+port, "cmd", "/k", "\""+exe+"\"", "http", port)
						c.Start()
						return nil
					}
				}
				if msg.String() == "ctrl+r" {
					// Force Reset: Clear path and go to main menu
					m.NgrokService.SavePath("") // Clear config
					m.NgrokStep = NgrokMainMenu
					return m, nil
				}
				if msg.String() == "esc" {
					// Smart Back: If we are fully configured, go back to Project Menu.
					// If we are setting up manually, go back to Main Menu.
					if m.NgrokService.GetExecutable() != "" {
						m.State = StateProjectActions
					} else {
						m.NgrokStep = NgrokMainMenu
					}
					m.NgrokPortInput.Blur()
				}
				return m, cmd
			case NgrokRunning:
				if msg.String() == "esc" {
					m.State = StateProjectActions
				}
			}

			// Global Ngrok Quit
			if msg.String() == "q" {
				return m, tea.Quit
			}

		case StateProjectActions:
			switch msg.String() {
			case "1", "f":
				return m, func() tea.Msg { m.Launcher.LaunchProject(m.Selected, "frontend"); return nil }
			case "2", "b":
				return m, func() tea.Msg { m.Launcher.LaunchProject(m.Selected, "backend"); return nil }
			case "3", "l":
				return m, func() tea.Msg { m.Launcher.LaunchProject(m.Selected, "full"); return nil }
			case "4":
				// Ngrok Flow - Smart Skip
				m.State = StateNgrok
				if m.NgrokService.GetExecutable() != "" {
					// Path is known, assume detailed setup is done -> Jump to Port
					m.NgrokStep = NgrokAskPort
					m.NgrokPortInput.Focus()
				} else {
					// First time? Show Setup Menu
					m.NgrokStep = NgrokMainMenu
				}
				return m, nil
			case "5", "c":
				// AI Context
				return m, func() tea.Msg {
					tree, err := m.TreeGen.GenerateTree(m.Selected.Path)
					if err != nil {
						return errMsg(err)
					}
					return contextMsg(tree)
				}
			case "6":
				// Doctor
				m.State = StateDependencyDoctor
				m.Err = nil
				m.Table.SetRows([]table.Row{}) // Clear old results
				m.AllPackagesUpToDate = false  // Reset flag
				return m, tea.Batch(m.Spinner.Tick, m.checkDependenciesCmd())
			case "\"", "backspace":
				m.State = StateProjectSelect
			case "q":
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height - 10)

	case projectMsg:
		m.Projects = msg
		m.State = StateProjectSelect // Direkt listeye git
		// Listeyi burada başlat
		items := make([]list.Item, len(m.Projects))
		for i, p := range m.Projects {
			// Technology icon helper
			getTechIcon := func(techType domain.ProjectType) string {
				icons := map[domain.ProjectType]string{
					domain.TypeNext:        "⚡",
					domain.TypeReact:       "⚛️",
					domain.TypeVue:         "💚",
					domain.TypeVite:        "⚡",
					domain.TypeReactNative: "📱",
					domain.TypeMobile:      "📱",
					domain.TypeHTML:        "🌐",
					domain.TypeTypeScript:  "🔷",
					domain.TypeNest:        "🐱",
					domain.TypeExpress:     "🚂",
					domain.TypeGo:          "🐹",
					domain.TypeDjango:      "🐍",
					domain.TypeFlask:       "🧪",
					domain.TypeLaravel:     "🐘",
					domain.TypeSpring:      "☕",
					domain.TypePHP:         "🐘",
					domain.TypeDocker:      "🐳",
				}
				if icon, ok := icons[techType]; ok {
					return icon
				}
				return ""
			}

			// Build combined icon (Frontend + Backend)
			var iconParts []string
			if p.FrontendType != "" && p.FrontendType != domain.TypeUnknown {
				if ic := getTechIcon(p.FrontendType); ic != "" {
					iconParts = append(iconParts, ic)
				}
			}
			if p.BackendType != "" && p.BackendType != domain.TypeUnknown {
				if ic := getTechIcon(p.BackendType); ic != "" {
					iconParts = append(iconParts, ic)
				}
			}
			// Docker indicator
			if p.HasDocker {
				iconParts = append(iconParts, "🐳")
			}

			icon := "📁 "
			if len(iconParts) > 0 {
				icon = strings.Join(iconParts, "") + " "
			}

			// Build technology description (Frontend + Backend names)
			var techParts []string
			if p.FrontendType != "" && p.FrontendType != domain.TypeUnknown {
				techParts = append(techParts, string(p.FrontendType))
			}
			if p.BackendType != "" && p.BackendType != domain.TypeUnknown {
				techParts = append(techParts, string(p.BackendType))
			}

			techDesc := "Bilinmeyen"
			if len(techParts) > 0 {
				techDesc = strings.Join(techParts, " + ")
			}

			// Title: Icon + Name
			items[i] = item{title: icon + p.Name, desc: techDesc + " | " + p.Path, project: &m.Projects[i]}
		}

		// List Configuration
		delegate := list.NewDefaultDelegate()
		delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.BorderLeftForeground(lipgloss.Color("#bd93f9")).Foreground(lipgloss.Color("#bd93f9"))
		delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.BorderLeftForeground(lipgloss.Color("#bd93f9")).Foreground(lipgloss.Color("#6272a4"))

		// Setup Help Styles to match standard
		// delegate.Styles.HelpStyle ... (Usually internal, but we can verify)

		m.List = list.New(items, delegate, m.Width, m.Height)
		m.List.Title = "🚀 PROJELER"
		m.List.SetShowTitle(true)
		m.List.SetStatusBarItemName("Proje", "Proje")
		m.List.FilterInput.Prompt = "🔍 Ara: "
		m.List.DisableQuitKeybindings()

		// Translate KeyMap (Help)
		m.List.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{
				key.NewBinding(
					key.WithKeys("q"),
					key.WithHelp(
						lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render("q"),
						lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render("Çıkış"),
					),
				),
			}
		}

		// Custom Key Bindings
		m.List.KeyMap.CursorUp.SetKeys("up")
		m.List.KeyMap.CursorDown.SetKeys("down")
		m.List.KeyMap.Filter.SetKeys("tab")
		m.List.KeyMap.ClearFilter.SetKeys("tab", "esc") // Tab toggle and Esc fallback
		m.List.KeyMap.ShowFullHelp.SetKeys(",")
		m.List.KeyMap.CloseFullHelp.SetKeys(",")

		m.List.KeyMap.CursorUp.SetHelp("↑", "Yukarı")
		m.List.KeyMap.CursorDown.SetHelp("↓", "Aşağı")
		m.List.KeyMap.Filter.SetHelp("tab", "Ara")
		m.List.KeyMap.ClearFilter.SetHelp("tab/esc", "Vazgeç")
		m.List.KeyMap.AcceptWhileFiltering.SetHelp("enter", "Seç")
		m.List.KeyMap.ShowFullHelp.SetHelp(",", "Daha Fazla")
		m.List.KeyMap.CloseFullHelp.SetHelp(",", "Kapat")
		m.List.KeyMap.Quit.SetHelp("q", "Çıkış") // Standart quit

		// Bizim custom q implementasyonumuzu menüde kırmızı göstermek için ekledik.
		// Ancak standart Help de aktif. Onu da yönetelim.

	case contextMsg:
		clipboard.WriteAll(string(msg))
		m.CopiedSuccess = true
		// Reset flag after 4 seconds
		cmd = tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return copiedResetMsg{}
		})
		cmds = append(cmds, cmd)
		m.State = StateProjectActions // Stay here

	case copiedResetMsg:
		m.CopiedSuccess = false

	case errMsg:
		m.Err = error(msg)
		// Clear any ongoing operations

	case doctorMsg:
		if len(msg) == 0 {
			// No outdated packages - keep table empty and set flag
			m.Table.SetRows([]table.Row{})
			m.AllPackagesUpToDate = true
		} else {
			// Outdated packages exist - populate table
			rows := []table.Row{}
			for pkg, info := range msg {
				rows = append(rows, table.Row{pkg, info.Current, info.Wanted, info.Latest})
			}
			m.Table.SetRows(rows)
			m.AllPackagesUpToDate = false
		}
		m.Err = nil // Clear any previous errors
		// Tablo güncellendi
	}

	// Alt bileşenleri güncelle
	switch m.State {
	case StateScanning, StateNgrok, StateDependencyDoctor:
		m.Spinner, cmd = m.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	case StateProjectSelect:
		// "Tab" ile filtreleme modu kapatma (Toggle)
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "tab" && m.List.FilterState() == list.Filtering {
			msg = tea.KeyMsg{Type: tea.KeyEsc}
		}

		// "q" ile çıkış (eğer filtreleme modunda değilse)
		if m.List.FilterState() == list.Filtering {
			// Filtreleme modundaysa listeye bırak
		} else if key, ok := msg.(tea.KeyMsg); ok && key.String() == "q" {
			return m, tea.Quit
		}

		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)

		// Liste seçimini yönet
		if m.State == StateProjectSelect {
			if val, ok := msg.(tea.KeyMsg); ok && val.String() == "enter" {
				// Proje seçildi
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Selected = i.project
					if m.List.Title == "Bağımlılık Kontrolü İçin Proje Seç" {
						// Doktoru çalıştır
						m.State = StateDependencyDoctor
						m.Err = nil
						// Spinner başlat ve komutu tetikle
						cmds = append(cmds, m.Spinner.Tick, m.checkDependenciesCmd())
					} else {
						m.State = StateProjectActions // Alt menüye git
					}
				}
			}
		}

	}

	return m, tea.Batch(cmds...)
}

func (m *MainModel) checkDependenciesCmd() tea.Cmd {
	return func() tea.Msg {
		res, err := m.Doctor.CheckDependencies(m.Selected.Path)
		if err != nil {
			return errMsg(err)
		}
		return doctorMsg(res)
	}
}

func (m *MainModel) View() string {
	if m.Err != nil {
		return "Hata: " + m.Err.Error()
	}

	switch m.State {
	case StateScanning:
		return fmt.Sprintf("\n\n   %s Taranıyor... %d yol bulundu.\n\n", m.Spinner.View(), len(m.Config.ProjectsPaths))
	case StateFirstRun:
		return fmt.Sprintf("\n\n  👋 Hoşgeldiniz! \n\n  Lütfen projelerinizin bulunduğu ana klasör yolunu giriniz:\n  (Tırnak işareti ile yapıştırabilirsiniz)\n\n  %s\n\n  [Enter] Kaydet\n", m.FirstRunInput.View())
	case StateDashboard:
		return m.dashboardView()
	case StateProjectSelect:
		v := m.List.View()
		v = strings.ReplaceAll(v, "No Proje found", "Proje bulunamadı")
		v = strings.ReplaceAll(v, "No Proje", "Proje yok") // Fallback
		return v
	case StateProjectActions:
		return m.actionsView()
	case StateDependencyDoctor:
		footer := m.renderFooter("Esc", "Geri Dön")

		// Show error if exists
		if m.Err != nil {
			errorMsg := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff5555")).
				Render(fmt.Sprintf("\n  ⚠️  Hata: %s", m.Err.Error()))
			return fmt.Sprintf("\n  %s İçin Doktor Raporu\n%s\n\n  %s", m.Selected.Name, errorMsg, footer)
		}

		// Show loading or results
		var tableView string
		if len(m.Table.Rows()) == 0 {
			// Empty table - either loading or all up to date
			if m.AllPackagesUpToDate {
				// All packages are up to date
				tableView = "\n  " + lipgloss.NewStyle().
					Foreground(lipgloss.Color("#50fa7b")).
					Render("✅ Tüm paketler güncel!")
			} else {
				// Still loading
				tableView = "\n  " + m.Spinner.View() + " Paketler kontrol ediliyor..."
			}
		} else {
			// Results arrived - show table
			tableView = "\n" + m.Table.View()
		}

		return fmt.Sprintf("\n  %s İçin Doktor Raporu%s\n\n  %s", m.Selected.Name, tableView, footer)
	case StateNgrok:
		return m.ngrokView()
	}

	return "Bilinmeyen Durum"
}

// Msg types
type errMsg error
type contextMsg string
type doctorMsg service.NpmOutdatedResult
type ngrokInstalledMsg string // changed to string (path)
type ngrokAuthMsg bool

func (m *MainModel) checkNgrokInstallCmd() tea.Cmd {
	return func() tea.Msg {
		// Temporary Bypass to verify UI flow
		c := make(chan string, 1)
		go func() {
			// m.NgrokService.CheckCommonPaths()
			time.Sleep(500 * time.Millisecond) // Simulate work
			c <- ""                            // Return empty (Not found)
		}()

		select {
		case res := <-c:
			return ngrokInstalledMsg(res)
		case <-time.After(1000 * time.Millisecond):
			return ngrokInstalledMsg("")
		}
	}
}

func (m *MainModel) checkNgrokAuthCmd() tea.Cmd {
	return func() tea.Msg {
		return ngrokAuthMsg(m.NgrokService.HasAuthToken())
	}
}

func (m *MainModel) ngrokView() string {
	var content, footer string

	s := "\n  🌐 Ngrok Bağlantı Sihirbazı\n\n"

	switch m.NgrokStep {
	case NgrokMainMenu:
		content = s + "  Ne yapmak istersiniz?\n\n" +
			"  [1] Manuel Yol Gir\n" +
			"  [2] Ngrok İndir\n"
		footer = m.renderFooter("Esc", "İptal")

	case NgrokCheckInstall:
		content = s + "  " + m.Spinner.View() + " Ngrok kontrol ediliyor..."
		// No footer here

	case NgrokModeSelect:
		// ...

	case NgrokManualPath:
		content = s + "  Lütfen 'ngrok.exe' dosyasının tam yolunu girin:\n"
		if m.NgrokCmd == "downloading" {
			content += lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("  (Tarayıcı açıldı. İndirdikten sonra kurulum yolunu girin)") + "\n"
		} else {
			content += "  (Örn: C:\\ProgramData\\chocolatey\\bin\\ngrok.exe)\n"
		}
		content += "\n" + fmt.Sprintf("  Yol: %s\n", m.NgrokPathInput.View())

		footer = m.renderFooter("Enter", "Kaydet", "Esc", "Geri")

	case NgrokCheckAuth:
		content = s + "  Authtoken ekranına geçiliyor..."

	case NgrokAuth:
		content = s + "  🔐 Ngrok Authtoken Gerekiyor\n\n" +
			"  1. Tarayıcıda şu adrese gidin:\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Render("     https://dashboard.ngrok.com/get-started/your-authtoken") + "\n\n" +
			"  2. 'Your Authtoken' yazısının altındaki " + lipgloss.NewStyle().Bold(true).Render("Copy") + " butonuna basın.\n" +
			"  3. Çıkan kodu kopyalayıp aşağıya yapıştırın:\n" +
			"\n" + fmt.Sprintf("  Token: %s\n", m.NgrokTokenInput.View())

		footer = m.renderFooter("Enter", "Kaydet", "Ctrl+O", "Siteye Git", "Esc", "Geri")

	case NgrokAskPort:
		content = s + "  Bağlantı Portu:\n" +
			fmt.Sprintf("  %s\n", m.NgrokPortInput.View())

		footer = m.renderFooter("Enter", "Başlat", "Ctrl+R", "Yeniden Kur", "Esc", "Geri")

	case NgrokRunning:
		content = s + "  🚀 Ngrok Çalışıyor!\n" +
			"  Yeni sekmede tünel açıldı."

		footer = m.renderFooter("Esc", "Projeye Dön")
	}

	// Calculate vertical space needed to push footer to bottom
	if footer != "" {
		hContent := lipgloss.Height(content)
		hFooter := lipgloss.Height(footer)
		gap := m.Height - hContent - hFooter - 1 // -1 for safety margin
		if gap > 0 {
			content += strings.Repeat("\n", gap)
		} else {
			content += "\n\n" // Fallback minimum spacing
		}
		content += "  " + footer // Add some left padding
	}

	return content
}

// -- Helper Types & Cmds --

type projectMsg []domain.Project

func (m *MainModel) scanProjectsCmd() tea.Cmd {
	return func() tea.Msg {
		projs := m.Scanner.ScanProjects()
		return projectMsg(projs)
	}
}

// List Item Adapter
type item struct {
	title, desc string
	project     *domain.Project
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.project.Name }

// ... (List initialization loop in Update) ...
// Instead of editing the loop here (which is in Update), I will edit the FilterValue method first.
// Wait, I can do both in one REPLACE if they are close, or separate tools.
// They are far apart (Update around 200, Item at bottom).
// I will do FilterValue first.
