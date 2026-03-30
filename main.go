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

// ── Estilos: Fusión de Cuadros + Layout de Foto ───────────────────────────

var (
	accentColor = lipgloss.Color("#7C3AED") // Morado original
	cyanColor   = lipgloss.Color("#22D3EE") // Cian de la foto
	grayColor   = lipgloss.Color("#6B7280")
	whiteColor  = lipgloss.Color("#F9FAFB")

	// Estilo para la columna izquierda (Arte ASCII)
	asciiStyle = lipgloss.NewStyle().Foreground(grayColor).MarginRight(4)

	// Estilo de los cuadros (Tus cuadros originales)
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			Width(60)

	titleStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	labelStyle = lipgloss.NewStyle().Foreground(cyanColor).Bold(true)
	dimStyle   = lipgloss.NewStyle().Foreground(grayColor)

	// Pestañas
	activeTabStyle = lipgloss.NewStyle().Foreground(whiteColor).Background(accentColor).Padding(0, 1)
	tabStyle       = lipgloss.NewStyle().Foreground(grayColor).Padding(0, 1)

	// Estilo para el cursor cian (corrección para .Render)
	cyanTextStyle = lipgloss.NewStyle().Foreground(cyanColor)
)

// ── Arte ASCII ────────────────────────────────────────────────────────────

const avatarASCII = `                                                                                                                                                                          
                    h                                               
          ########                                                                
         ##      #                                                              
        #         ##                                                            
       #           ##                                                           
       #            #      ########      ###                                    
        #            #   ##        ##   #  ###                                  
         #            ####          ## ##  ###                                  
         #              #   #        ####    #                                  
          # ##          #         #          #                                  
      #             #          #        #                                   
     #            ####    ######  ###   #                                   
      #         #                    ####                                   
          #   ####                           ##                                
            ##   ##                             #                               
           #    # #          ####                ##                             
         ##     # ##        ######                 #                            
        ##      ####        ######        #        #                            
        #       ####                               ##                           
        #       ####                     #           #                          
        #        ###            #                     #                         
       #                       #    #       ##         #                        
       #       #            #       #       ## #         #                       
      #                    #               ####         #                       
       #      #       ###                    ###         #                       
      #            #####                         #      #                       
      #            # ###                         #      #                       
      #     #       ###                                #                        
        #                                               #                        
        #                     ####            #      #                          
         #                    ####           ##     #                           
         ##          #                     #       #                            
          #         ####              #   #                                     
           ##                    ###             #                              
             ###                              ##                                
                 ##                         ##                                  
                   ####            ########                                     
                          #####                                                 
                                                                                                                             
                                                                                                    
`
// ── Datos ─────────────────────────────────────────────────────────────────

const nombre = "Alba Costas Fernández"
const rol = "Developer"

type Proyecto struct {
	Nombre string
	Desc   string
	Stack  string
	URL    string
}

var proyectos = []Proyecto{
	{ 	
		Nombre: "SSH Portfolio", 
		Desc: "Interactive ssh server.", 
		Stack: "Go · Wish", 
		URL: "github.com/albacostas/Portfolio-ssh",
	},
	{
		Nombre: "Interactive-Notch", 
		Desc: "App for macOS notch functions.", 
		Stack: "Swift · SwiftUI", 
		URL: "github.com/albacostas/Notch",
	},
	{
		Nombre: "Expenses", 
		Desc: "Track your expenses on macOS.", 
		Stack: "Swift · CoreData", 
		URL: "github.com/albacostas/Expenses",
	},
}

var tabs = []string{"About me", "Projects", "Contact"}

// ── Modelo ────────────────────────────────────────────────────────────────

type model struct {
	tab     int
	proyIdx int
	width   int
	height  int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
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

// ── Renderizado ───────────────────────────────────────────────────────────

func (m model) View() string {
	if m.width == 0 {
		return "Cargando..."
	}

	// 1. Columna Izquierda: Arte ASCII
	colIzquierda := asciiStyle.Render(avatarASCII)

	// 2. Cabecera (Nombre y Ubicación)
	header := titleStyle.Render(fmt.Sprintf("%s — %s", nombre, rol)) + "\n" +
		dimStyle.Render("📍 Santiago de Compostela, Spain")

	// 3. Pestañas
	var tabsRow string
	for i, t := range tabs {
		if i == m.tab {
			tabsRow += activeTabStyle.Render(t) + " "
		} else {
			tabsRow += tabStyle.Render(t) + " "
		}
	}

	// 4. Contenido en CUADRO (boxStyle)
	var dynamicContent string
	switch m.tab {
	case 0: // About
		dynamicContent = boxStyle.Render(
			labelStyle.Render("$ whoami") + "\n\n" +
				"Apasionada por construir aplicaciones desde cero,\n" +
				"explorar nuevos lenguajes y llevar ideas al código.")
	case 1: // Projects
		dynamicContent = m.viewProyectos()
	case 2: // Contact
		dynamicContent = boxStyle.Render(
			labelStyle.Render("$ contact --list") + "\n\n" +
				"Email:    albacostasfernandez@gmail.com\n" +
				"GitHub:   github.com/albacostas\n" +
				"LinkedIn: linkedin.com/in/albacostasfernandez")
	}

	// 5. Unimos la columna derecha (Header + Tabs + Cuadro)
	colDerecha := lipgloss.JoinVertical(lipgloss.Left, header, "\n", tabsRow, dynamicContent)

	// 6. Resultado final: ASCII a la izquierda, todo lo demás a la derecha
	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, colIzquierda, colDerecha)

	footer := dimStyle.Render("\n  ←/→ tabs   ↑/↓ navegar   q salir")

	return lipgloss.JoinVertical(lipgloss.Left, mainLayout, footer)
}

func (m model) viewProyectos() string {
	p := proyectos[m.proyIdx]

	// Lista de navegación de proyectos
	lista := ""
	for i, proj := range proyectos {
		if i == m.proyIdx {
			lista += cyanTextStyle.Render("▶ ") + proj.Nombre + "\n"
		} else {
			lista += "  " + proj.Nombre + "\n"
		}
	}

	contenido := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("PROJECTS"),
		"\n"+lista,
		"\n"+labelStyle.Render("DETAILS"),
		p.Desc,
		dimStyle.Render("Stack: ")+p.Stack,
		dimStyle.Render("URL:   ")+p.URL,
	)

	return boxStyle.Render(contenido)
}

// ── Servidor SSH ───────────────────────────────────────────────────────────

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:2222"),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bbtea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return model{}, []tea.ProgramOption{tea.WithAltScreen()}
			}),
		),
	)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🚀 Portfolio SSH en :2222")
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			fmt.Println("Error:", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}