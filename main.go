package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	//"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bbtea "github.com/charmbracelet/wish/bubbletea"
)

var (
	accentColor = lipgloss.Color("#7C3AED")
	cyanColor   = lipgloss.Color("#22D3EE")
	grayColor   = lipgloss.Color("#6B7280")
	whiteColor  = lipgloss.Color("#F9FAFB")
	greenColor  = lipgloss.Color("#10B981")

	asciiStyle = lipgloss.NewStyle().
		Foreground(grayColor).
		Padding(0).
		Margin(0)

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		Width(80)

	titleStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	labelStyle = lipgloss.NewStyle().Foreground(cyanColor).Bold(true)
	dimStyle   = lipgloss.NewStyle().Foreground(grayColor)

	activeTabStyle = lipgloss.NewStyle().Foreground(whiteColor).Background(accentColor).Padding(0, 1)
	tabStyle       = lipgloss.NewStyle().Foreground(grayColor).Padding(0, 1)
	cyanTextStyle  = lipgloss.NewStyle().Foreground(cyanColor)
	greenTextStyle = lipgloss.NewStyle().Foreground(greenColor)
)
/*
const avatarASCII = `
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
       #       #            #       #       ## #        # 
      #                    #               ####          #  
       #      #       ###                    ###         #    
      #            #####                         #       #  
      #            # ###                         #      #   
      #     #       ###                                #   
        #                                             #  
        #                     ####            #      #  
         #                    ####           ##     # 
         ##          #                     #       #                            
          #         ####              #   #       #  
           ##                    ###             #   
             ###                              ##   
                 ##                         ##     
                   ####            ########   
                          #####`
*/

                                                                                                                                                            
const avatarASCII = `                                                                                                                                                          
                           @@@@@@@                @@@@@                      
                         @@@      @@           @@@     @@                    
                       @@@         @@        @@         }@@                  
                      @@            @@      @@           @@                  
                     @@              @@    @@             @@                 
                     @@              @@    @@             @@                 
                    |@                @   @@              @@                 
                    @@                @@  @@              @@                 
                    @@                @@  @               @@                 
                    @@                @@ @@              @@                  
                    @@                @@ @@              @@                  
                    @@                @@ @@              @@                  
                    @@                @@ @@             @@                   
                     @                @@ @@             @@                   
                     @@               @@@@@            @@                    
                     @@                                @                     
                      @@                              @@                     
                   @@@@                               @@@                    
                 @@@                                     @@                  
               @@@                                         @@                
              @@@                                           @@               
              @@                                             @@              
             @@                                               @@             
             @                                                @@             
            @@                                                @@             
            @@                                                 @             
            @@                                                 @             
             @          @@@                       @@          @@             
             @@                                  @@@          @@             
             @@                                              @@              
              @@                    @                       @@               
               @@@                   @@@@@                 @@                
                 @@@               @@@@ @@@              @@@                 
                   @@@                                 @@@                   
                      @@@@@                         @@@                      
                           @@@@@@@         @@@@@@@@@                         
                                    @@@@@`                                                                         
                                                                                                                                                            
                                                                                                                                                            
                                                                                                                                                            
                                                                                                                                                            
                                                                                                                                                            
                                                                                                                                                            
// ── Data Structures ───────────────────────────────────────────────────────

type Project struct {
	Name  string
	Desc  string
	Stack string
	URL   string
}

type Skill struct {
	Category string
	Items    []string
}

type Experience struct {
	Title    string
	Company  string
	Period   string
	Desc     string
}

// ── Portfolio Data ────────────────────────────────────────────────────────

var projects = []Project{
	{
		Name: "SSH Portfolio",
		Desc: "Interactive SSH server with dynamic terminal UI\n" +
			"  Features: Real-time navigation, ASCII art rendering\n" +
			"  Deployed and accessible via SSH",
		Stack: "Go · Charmbracelet · Bubbletea · Wish",
		URL:   "github.com/albacostas/Portfolio-ssh",
	},
	{
		Name: "Interactive Notch",
		Desc: "macOS app for managing and controlling notch functions\n" +
			"  Features: Real-time status display, system integration\n",
		Stack: "Swift · SwiftUI · CoreData",
		URL:   "github.com/albacostas/Notch",
	},
	{
		Name: "Expenses Tracker",
		Desc: "Native macOS expense tracking application\n" +
			"  Features: Data persistence, charts, export to CSV\n" +
			"  Optimized for performance and user experience",
		Stack: "Swift · CoreData · AppKit",
		URL:   "github.com/albacostas/Expenses",
	},
	{
		Name: "3D Audio Visualizer",
		Desc: "Interactive 3D engine that transforms real-tiem audio signals into dynamic mesh deformactions. \n" +
			"  Features: 2,000+ reactive particles, frequency-based sphere pulsing, and custom lighting.\n" +
			"  Built with low-latency signal processing (512-frame buffer).",
		Stack: "C++ · OpenGL · PortAudio · GLUT",
		URL:   "https://github.com/albacostas/audio-visualizer",
	},
	{
		Name: "MyCalendar",
		Desc: "Calendar app for ios with a focus on simplicity and usability.\n"+
			"  Clean and intuitive design",
		Stack: "Swift · SwiftUI · EventKit",
		URL:   "github.com/albacostas/MyCalendar",
	},
}

var skills = []Skill{
	{
		Category: "Languages",
		Items:    []string{"Swift", "Java", "C", "C++", "Python"},
	},
	{
		Category: "Frontend",
		Items:    []string{"SwiftUI", "HTML/CSS"},
		//Items:    []string{"SwiftUI", "HTML/CSS", "Tailwind CSS"},
	},
	{
		Category: "Backend",
		Items:    []string{"Node.js", "REST APIs", "PostgreSQL"},
	},
	{
		Category: "Tools & DevOps",
		Items:    []string{"Git", "Docker", "Xcode", "VS Code", "GitHub Actions", "Linux"},
	},
	{
		Category: "Platforms",
		Items:    []string{"iOS", "Linux", "Web"},
	},
}

var experience = []Experience{
	{
		Title:   "Computer Science Student",
		Company: "University of Santiago de Compostela",
		Period:  "2023 - Present",
		Desc:    "Focusing on algorithms, databases and system design",
	},
	{
		Title:   "High School Diploma - Science Concentration",
		Company: "Colegio Apóstol Santiago - Jesuitas",
		Period:  "2021 - 2023",
		Desc:    "Advanced coursework in Mathematics, Physics and Chemistry.",
	},
}

var tabs = []string{"About", "Skills", "Projects", "Experience", "Contact"}

// ── Model ─────────────────────────────────────────────────────────────────

type model struct {
	tab     int
	projIdx int
	skillIdx int
	expIdx  int
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
			m.projIdx = 0
			m.skillIdx = 0
			m.expIdx = 0
		case "left", "h", "shift+tab":
			m.tab = (m.tab - 1 + len(tabs)) % len(tabs)
			m.projIdx = 0
			m.skillIdx = 0
			m.expIdx = 0
		case "down", "j":
			if m.tab == 2 {
				m.projIdx = (m.projIdx + 1) % len(projects)
			} else if m.tab == 1 {
				m.skillIdx = (m.skillIdx + 1) % len(skills)
			} else if m.tab == 3 {
				m.expIdx = (m.expIdx + 1) % len(experience)
			}
		case "up", "k":
			if m.tab == 2 {
				m.projIdx = (m.projIdx - 1 + len(projects)) % len(projects)
			} else if m.tab == 1 {
				m.skillIdx = (m.skillIdx - 1 + len(skills)) % len(skills)
			} else if m.tab == 3 {
				m.expIdx = (m.expIdx - 1 + len(experience)) % len(experience)
			}
		}
	}
	return m, nil
}

// ── Rendering ─────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.width == 0 {return "Loading..."}

	asciiContent := avatarASCII
	colLeft := asciiStyle.Width(80).Render(asciiContent)
	
	// Header
	header := titleStyle.Render("Alba Costas Fernández — Developer") + "\n" +
		dimStyle.Render("📍 Santiago de Compostela, Spain") + "\n" +
		dimStyle.Render("Building apps. Learning continuously. Open to opportunities.")

	// Tabs
	var tabsRow string
	for i, t := range tabs {
		if i == m.tab {
			tabsRow += activeTabStyle.Render(t) + " "
		} else {
			tabsRow += tabStyle.Render(t) + " "
		}
	}

	// Dynamic content
	var dynamicContent string
	switch m.tab {
		case 0:
			dynamicContent = m.viewAbout()
		case 1:
			dynamicContent = m.viewSkills()
		case 2:
			dynamicContent = m.viewProjects()
		case 3:
			dynamicContent = m.viewExperience()
		case 4:
			dynamicContent = m.viewContact()
	}

	// Right column
	colRight := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"\n",
		tabsRow,
		dynamicContent,
	)

	// Main layout
	mainLayout := lipgloss.JoinHorizontal(
			lipgloss.Top, 
			colLeft, 
			lipgloss.NewStyle().MarginLeft(2).Render(colRight), // Añade un margen manual aquí
		)
	footer := "\n" + dimStyle.Render("  ←/→ tabs   ↑/↓ navigate   q quit")

	return lipgloss.JoinVertical(lipgloss.Left, mainLayout, footer)
}

func (m model) viewAbout() string {
	content := labelStyle.Render("$ whoami") + "\n\n" +
		"Developer with a passion for building things and learning continuosly.\n" +
		"I'm passionate about exploring new technologies and languages, and creating products that solve real problems.\n\n" +
		"I actively participate in CTF competitions, honing my skills in problem-solving, and cybersecurity talks.\n\n" +
		"Beyond my university studies, I'm self-driven learner who pursues independent projects and explores emerging techonologies whenever my curiosity is sparked. \n\n" +
		"Currently pursuing my Computer Science degree while looking for the next challenge."
		greenTextStyle.Render("✓ Open to opportunities and collaborations")

	return boxStyle.Render(content)
}

func (m model) viewSkills() string {
	skill := skills[m.skillIdx]
	
	// Skill list
	skillList := ""
	for i, s := range skills {
		if i == m.skillIdx {
			skillList += cyanTextStyle.Render("▶ ") + s.Category + "\n"
		} else {
			skillList += "  " + s.Category + "\n"
		}
	}

	// Skill items
	skillItems := ""
	for _, item := range skill.Items {
		skillItems += "  • " + item + "\n"
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("CATEGORIES"),
		"\n"+skillList,
		"\n"+labelStyle.Render(skill.Category),
		skillItems,
	)

	return boxStyle.Render(content)
}

func (m model) viewProjects() string {
	p := projects[m.projIdx]
	
	// Project list
	projList := ""
	for i, proj := range projects {
		if i == m.projIdx {
			projList += cyanTextStyle.Render("▶ ") + proj.Name + "\n"
		} else {
			projList += "  " + proj.Name + "\n"
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("PROJECTS"),
		"\n"+projList,
		"\n"+labelStyle.Render("DETAILS"),
		p.Desc,
		"\n"+dimStyle.Render("Stack: ")+p.Stack,
		dimStyle.Render("GitHub: ")+p.URL,
	)

	return boxStyle.Render(content)
}

func (m model) viewExperience() string {
	exp := experience[m.expIdx]
	
	// Experience list
	expList := ""
	for i, e := range experience {
		if i == m.expIdx {
			expList += cyanTextStyle.Render("▶ ") + e.Title + "\n"
		} else {
			expList += "  " + e.Title + "\n"
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("EXPERIENCE"),
		"\n"+expList,
		"\n"+labelStyle.Render(exp.Title),
		dimStyle.Render(exp.Company+" • "+exp.Period),
		"\n"+exp.Desc,
	)

	return boxStyle.Render(content)
}

func (m model) viewContact() string {
	content := labelStyle.Render("$ contact --list") + "\n\n" +
		"Email:    albacostasfernandez@gmail.com\n" +
		"GitHub:   github.com/albacostas\n" +
		"LinkedIn: linkedin.com/in/albacostasfernandez\n" +
		"Twitter:  @albacostas\n\n" +
		greenTextStyle.Render("Let's build something amazing together! 🚀")

	return boxStyle.Render(content)
}

// ── Server ────────────────────────────────────────────────────────────────

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

	fmt.Println("🚀 SSH Portfolio running on :2222")
	fmt.Println("📍 Connect with: ssh -p 2222 localhost")
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
