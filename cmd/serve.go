package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	pst "github.com/Bril3d/minicontainer/internal/preset"
	rt "github.com/Bril3d/minicontainer/internal/runtime"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API daemon",
	Long:  "Start a REST API server to handle container operations for the browser frontend.",
	Run: func(cmd *cobra.Command, args []string) {
		podman := rt.NewPodmanRuntime()

		r := gin.Default()

		// Custom Gin CORS middleware
		r.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusOK)
				return
			}
			c.Next()
		})

		api := r.Group("/api")
		{
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

				manager, err := pst.NewManager(pst.GetDefaultPath())
				if err != nil {
					c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to load presets: %s", err))
					return
				}

				p, ok := manager.Find(presetName)
				if !ok {
					c.String(http.StatusNotFound, fmt.Sprintf("Preset '%s' not found", presetName))
					return
				}

				opts := rt.RunOptions{
					Image:   p.Image,
					Name:    fmt.Sprintf("%s-%d", presetName, time.Now().Unix()%10000),
					Ports:   p.Ports,
					Volumes: p.Volumes,
					Env:     p.Env,
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

		fmt.Printf("MiniContainer API Daemon starting on http://localhost:%d\n", servePort)
		fmt.Printf("CORS enabled for all origins (Development mode)\n")

		if err := r.Run(fmt.Sprintf(":%d", servePort)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Port to listen on")
	rootCmd.AddCommand(serveCmd)
}
