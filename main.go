package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bbtea "github.com/charmbracelet/wish/bubbletea"
)

// ── Estilos ────────────────────────────────────────────────────────────────

var (
	accent = lipgloss.Color("#7C3AED")
	muted  = lipgloss.Color("#6B7280")
	white  = lipgloss.Color("#F9FAFB")
	green  = lipgloss.Color("#10B981")

	titleStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			PaddingBottom(1)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(accent).
			PaddingLeft(2).PaddingRight(2).
			Bold(true)

	tabStyle = lipgloss.NewStyle().
			Foreground(muted).
			PaddingLeft(2).PaddingRight(2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(1, 2).
			MarginTop(1)

	labelStyle = lipgloss.NewStyle().Foreground(green).Bold(true)
	textStyle  = lipgloss.NewStyle().Foreground(white)
	dimStyle   = lipgloss.NewStyle().Foreground(muted)
)

// ── Datos del portafolio — EDITA ESTO ─────────────────────────────────────

const nombre = "Alba Costas"
const rol = "Developer"
const ubicacion = "España"

var sobre = []string{
	"Apasionada por construir aplicaciones desde cero, explorar nuevos lenguajes de programación y llevar ideas al código. ",
	"",
	"Actualmente aprendiendo Go y explorando el mundo de las TUI.",
}

type Proyecto struct {
	Nombre string
	Desc   string
	Stack  string
	URL    string
}

var proyectos = []Proyecto{
	{
		Nombre: "SSH Portfolio",
		Desc:   "Este mismo portafolio: un servidor SSH interactivo.",
		Stack:  "Go · Wish · Bubble Tea · Lip Gloss",
		URL:    "github.com/albacostas/Portfolio-ssh",
	},
	{
		Nombre: "Interactive-Notch",
		Desc:   "App para macOS desarrollada en Swift que añade utilidades al notch mostrando accesos rápidos a notas y a un spotify mini player.",
		Stack:  "Swift · SwiftUI",
		URL:    "github.com/albacostas/Notch",
	},
	{
		Nombre: "Expenses",
		Desc:   "App para macOS desarrollada en Swift que permite gestionar gastos e ingresos.",
		Stack:  "Swift · SwiftUI · CoreData",
		URL:    "github.com/albacostas/Expenses",
	},
}

var contacto = []string{
	"Email:    albacostasfernandez@gmail.com",
	"GitHub:   github.com/albacostas",
	"LinkedIn: www.linkedin.com/in/albacostasfernandez",
}

// ── Modelo ─────────────────────────────────────────────────────────────────

type model struct {
	tab     int
	proyIdx int
	width   int
	height  int
}

var tabs = []string{"Sobre mí", "Proyectos", "Contacto"}

func newModel() model { return model{} }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "right", "l", "tab":
			m.tab = (m.tab + 1) % len(tabs)
			m.proyIdx = 0
		case "left", "h", "shift+tab":
			m.tab = (m.tab - 1 + len(tabs)) % len(tabs)
			m.proyIdx = 0
		case "down", "j":
			if m.tab == 1 {
				m.proyIdx = (m.proyIdx + 1) % len(proyectos)
			}
		case "up", "k":
			if m.tab == 1 {
				m.proyIdx = (m.proyIdx - 1 + len(proyectos)) % len(proyectos)
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	w := m.width
	if w == 0 {
		w = 80
	}
	_ = w

	header := titleStyle.Render(fmt.Sprintf("  %s  —  %s", nombre, rol)) + "\n" +
		dimStyle.Render(fmt.Sprintf("  📍 %s", ubicacion))

	tabBar := ""
	for i, t := range tabs {
		if i == m.tab {
			tabBar += activeTabStyle.Render(t)
		} else {
			tabBar += tabStyle.Render(t)
		}
	}

	var content string
	switch m.tab {
	case 0:
		content = m.viewSobre()
	case 1:
		content = m.viewProyectos()
	case 2:
		content = m.viewContacto()
	}

	footer := dimStyle.Render("  ←/→ tabs   ↑/↓ navegar   q salir")
	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s", header, tabBar, content, footer)
}

func (m model) viewSobre() string {
	body := labelStyle.Render("$ whoami") + "\n\n"
	for _, l := range sobre {
		body += textStyle.Render(l) + "\n"
	}
	return boxStyle.Render(body)
}

func (m model) viewProyectos() string {
	lista := ""
	for i, p := range proyectos {
		cursor := "  "
		if i == m.proyIdx {
			cursor = "▶ "
		}
		lista += fmt.Sprintf("%s%s\n", cursor, p.Nombre)
	}

	p := proyectos[m.proyIdx]
	detalle := labelStyle.Render(p.Nombre) + "\n\n" +
		textStyle.Render(p.Desc) + "\n\n" +
		dimStyle.Render("Stack: ") + textStyle.Render(p.Stack) + "\n" +
		dimStyle.Render("URL:   ") + textStyle.Render(p.URL)

	left := boxStyle.Width(22).Render(lista)
	right := boxStyle.Width(48).Render(detalle)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
}

func (m model) viewContacto() string {
	body := labelStyle.Render("$ contact --list") + "\n\n"
	for _, l := range contacto {
		body += textStyle.Render(l) + "\n"
	}
	return boxStyle.Render(body)
}

// ── Servidor SSH ───────────────────────────────────────────────────────────
func main() {
	srv, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:2222"),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bbtea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return newModel(), []tea.ProgramOption{tea.WithAltScreen()}
			}),
		),
	)
	if err != nil {
		fmt.Println("Error creando servidor:", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("   SSH Portfolio escuchando en :2222")
	fmt.Println("   Conéctate con: ssh localhost -p 2222")

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			fmt.Println("Servidor detenido:", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("\n Servidor cerrado.")
}