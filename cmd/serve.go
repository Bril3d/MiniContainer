package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	pst "github.com/Bril3d/minicontainer/internal/preset"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var servePort int
var noGui bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start both MiniContainer API and GUI",
	Long:  "Start the backend REST API server and launch the Tauri GUI. Use --no-gui to run only the API.",
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		// Configure Gin
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(gin.Recovery())

		// Robust CORS middleware
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:1420", "http://localhost:8080"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "tauri://")
			},
			MaxAge: 12 * time.Hour,
		}))

		manager, err := pst.NewManager(pst.GetDefaultPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load presets: %v\n", err)
		}

		api := r.Group("/api")
		{
			api.GET("/presets", func(c *gin.Context) {
				if manager == nil {
					c.String(http.StatusInternalServerError, "Presets not loaded")
					return
				}
				c.JSON(http.StatusOK, manager.GetAll())
			})

			// Pause operations
			api.POST("/pause/:id", func(c *gin.Context) {
				id := c.Param("id")
				if err := podman.Pause(id); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "paused"})
			})

			api.POST("/unpause/:id", func(c *gin.Context) {
				id := c.Param("id")
				if err := podman.Unpause(id); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "unpaused"})
			})

			api.POST("/restart/:id", func(c *gin.Context) {
				id := c.Param("id")
				if err := podman.Restart(id); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "restarted"})
			})

			// Exec operations
			api.POST("/exec", func(c *gin.Context) {
				var opts rt.ExecOptions
				if err := c.ShouldBindJSON(&opts); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := podman.Exec(opts.Container, opts.Command, opts.Interactive); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "executed"})
			})

			api.POST("/build", func(c *gin.Context) {
				var opts rt.BuildOptions
				if err := c.ShouldBindJSON(&opts); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := podman.Build(opts); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "build started"})
			})

			api.GET("/ps", func(c *gin.Context) {
				containers, err := podman.List()
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.JSON(http.StatusOK, containers)
			})

			api.GET("/images", func(c *gin.Context) {
				images, err := podman.Images()
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.JSON(http.StatusOK, images)
			})

			api.POST("/rmi", func(c *gin.Context) {
				id := c.Query("id")
				if id == "" {
					c.String(http.StatusBadRequest, "Missing image ID")
					return
				}
				if err := podman.RemoveImage(id); err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.Status(http.StatusOK)
			})

			api.GET("/stats", func(c *gin.Context) {
				stats, err := podman.Stats()
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.JSON(http.StatusOK, stats)
			})

			api.POST("/deploy", func(c *gin.Context) {
				presetName := c.Query("id")
				if presetName == "" {
					c.String(http.StatusBadRequest, "Missing preset name")
					return
				}

				if manager == nil {
					c.String(http.StatusInternalServerError, "Presets not available")
					return
				}

				p, ok := manager.Find(presetName)
				if !ok {
					c.String(http.StatusNotFound, fmt.Sprintf("Preset '%s' not found", presetName))
					return
				}

				var resolvedVolumes []string
				for _, vol := range p.Volumes {
					// Handle volume string which might contain multiple colons (e.g., C:\path:/app)
					lastColon := strings.LastIndex(vol, ":")
					if lastColon > 0 {
						hostPath := vol[:lastColon]
						containerPath := vol[lastColon+1:]
						
						absPath, err := filepath.Abs(hostPath)
						if err == nil {
							// If on Windows and using WSL Podman, convert D:\ to /mnt/d/
							if runtime.GOOS == "windows" {
								absPath = filepath.ToSlash(absPath)
								if len(absPath) > 2 && absPath[1] == ':' {
									drive := strings.ToLower(string(absPath[0]))
									absPath = "/mnt/" + drive + absPath[2:]
								}
							}
							resolvedVolumes = append(resolvedVolumes, fmt.Sprintf("%s:%s", absPath, containerPath))
							continue
						}
					}
					resolvedVolumes = append(resolvedVolumes, vol)
				}

				opts := rt.RunOptions{
					Image:   p.Image,
					Name:    fmt.Sprintf("%s-%d", presetName, time.Now().Unix()%10000),
					Ports:   p.Ports,
					Volumes: resolvedVolumes,
					Env:     p.Env,
					Cmd:     strings.Fields(p.Cmd),
					Detach:  true,
				}

				id, err := podman.Run(opts)
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.JSON(http.StatusOK, gin.H{"id": id})
			})

			api.POST("/start", func(c *gin.Context) {
				name := c.Query("name")
				if name == "" {
					c.String(http.StatusBadRequest, "Missing name parameter")
					return
				}
				if err := podman.Start(name); err != nil {
					c.String(http.StatusInternalServerError, "Failed to start: "+err.Error())
					return
				}
				c.Status(http.StatusOK)
			})

			api.POST("/stop", func(c *gin.Context) {
				name := c.Query("name")
				if name == "" {
					c.String(http.StatusBadRequest, "Missing name parameter")
					return
				}
				if err := podman.Stop(name); err != nil {
					c.String(http.StatusInternalServerError, "Failed to stop: "+err.Error())
					return
				}
				c.Status(http.StatusOK)
			})

			api.POST("/remove", func(c *gin.Context) {
				name := c.Query("name")
				if name == "" {
					c.String(http.StatusBadRequest, "Missing name parameter")
					return
				}
				if err := podman.Remove(name, true); err != nil {
					c.String(http.StatusInternalServerError, "Failed to remove: "+err.Error())
					return
				}
				c.Status(http.StatusOK)
			})

			api.POST("/pause", func(c *gin.Context) {
				id := c.Query("id")
				if id == "" {
					c.String(http.StatusBadRequest, "Missing ID parameter")
					return
				}
				if err := podman.Pause(id); err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.Status(http.StatusOK)
			})

			api.POST("/unpause", func(c *gin.Context) {
				id := c.Query("id")
				if id == "" {
					c.String(http.StatusBadRequest, "Missing ID parameter")
					return
				}
				if err := podman.Unpause(id); err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.Status(http.StatusOK)
			})

			api.POST("/pull", func(c *gin.Context) {
				image := c.Query("image")
				if image == "" {
					c.String(http.StatusBadRequest, "Missing image name")
					return
				}
				if err := podman.Pull(image); err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.Status(http.StatusOK)
			})

			api.GET("/version", func(c *gin.Context) {
				v, err := podman.Version()
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				c.JSON(http.StatusOK, gin.H{"version": v})
			})
		}

		if noGui {
			fmt.Printf("MiniContainer API Daemon starting on http://localhost:%d\n", servePort)
			if err := r.Run(fmt.Sprintf(":%d", servePort)); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				os.Exit(1)
			}
		} else {
			// Start API in background
			go func() {
				if err := r.Run(fmt.Sprintf(":%d", servePort)); err != nil {
					fmt.Fprintf(os.Stderr, "Error starting API: %s\n", err)
					os.Exit(1)
				}
			}()

			// Wait for API to be ready
			waitForAPI(servePort)

			fmt.Println("Launching MiniContainer GUI...")
			if err := launchGUI(); err != nil {
				fmt.Fprintf(os.Stderr, "Error launching GUI: %s\n", err)
				os.Exit(1)
			}
		}
	},
}

func waitForAPI(port int) {
	url := fmt.Sprintf("http://localhost:%d/api/version", port)
	for i := 0; i < 20; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Println("Warning: API is taking longer than expected to start. Attempting to launch GUI anyway...")
}

func launchGUI() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "npm", "run", "tauri", "dev")
	} else {
		cmd = exec.Command("npm", "run", "tauri", "dev")
	}

	cmd.Dir = "gui"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Port to listen on")
	serveCmd.Flags().BoolVar(&noGui, "no-gui", false, "Skip launching the GUI")
	rootCmd.AddCommand(serveCmd)
}
