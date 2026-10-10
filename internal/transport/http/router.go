package http

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ui-design-backend/internal/config"
	"ui-design-backend/internal/middleware"
)

func NewRouter(cfg *config.Config, authHandler *AuthHandler, profileHandler *ProfileHandler, storyHandler *StoryHandler, authMW *middleware.AuthMiddleware, ebookHandlers ...*EbookHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.ClientOrigin, "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Static file serving for user uploads
	uploadPath := "./uploads"
	if cfg != nil && cfg.UploadDir != "" {
		uploadPath = cfg.UploadDir
	} else if profileHandler != nil && profileHandler.UploadDir != "" {
		uploadPath = profileHandler.UploadDir
	}
	uploadsDir := http.Dir(uploadPath)
	r.Get("/uploads/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(uploadsDir))
		fs.ServeHTTP(w, r)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	if authHandler != nil {
		r.Route("/api/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/admin/login", authHandler.AdminLogin)
			r.Post("/logout", authHandler.Logout)
			r.Post("/password-reset", authHandler.RequestPasswordReset)
			r.Post("/password-reset/validate", authHandler.ValidatePasswordReset)
			r.Post("/password-reset/confirm", authHandler.ConfirmPasswordReset)

			// Protected endpoints with transparent auto-refresh!
			if authMW != nil {
				r.Group(func(r chi.Router) {
					r.Use(authMW.RequireAuth)
					r.Get("/me", authHandler.GetMe)
				})
			}
		})
	}

	if authHandler != nil && authHandler.adminHandler != nil && authMW != nil {
		r.Route("/api/admin", func(r chi.Router) {
			r.Use(authMW.RequireAuth)
			r.Use(middleware.RequireStaff)
			r.Get("/stories", authHandler.adminHandler.ListStories)
			if storyHandler != nil { r.Get("/comments", storyHandler.PendingComments);r.Put("/comments/{id}",storyHandler.ModerateComment) }
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdmin)
				r.Get("/users", authHandler.adminHandler.ListUsers)
				r.Put("/users/{id}/role", authHandler.adminHandler.UpdateRole)
				r.Delete("/users/{id}", authHandler.adminHandler.DeleteUser)
			})
		})
	}

	if profileHandler != nil {
		// Profile endpoints
		r.Route("/api/profile", func(r chi.Router) {
			r.Get("/{userId}", profileHandler.GetPublicProfile)

			r.Group(func(r chi.Router) {
				if authMW != nil {
					r.Use(authMW.RequireAuth)
				}
				r.Get("/me", profileHandler.GetMyProfile)
				r.Put("/me", profileHandler.UpdateMyProfile)
			})
		})

		// Upload endpoint (protected)
		r.Group(func(r chi.Router) {
			if authMW != nil {
				r.Use(authMW.RequireAuth)
			}
			r.Post("/api/upload", profileHandler.UploadFile)
		})
	}

	if storyHandler != nil {
		r.Route("/api/stories", func(r chi.Router) {
			if authMW != nil {
				r.Use(authMW.RequireAuth)
			}
			r.Get("/", storyHandler.ListStories)

			r.Group(func(r chi.Router) {
				if authMW != nil {
					r.Use(authMW.RequireAuth)
				}
				r.Get("/mine", storyHandler.GetMyStories)
				r.Get("/{id}/reading-progress", storyHandler.GetReadingProgress)
				r.Put("/{id}/reading-progress", storyHandler.SaveReadingProgress)
				r.Post("/", storyHandler.CreateStory)
				r.Put("/{id}", storyHandler.UpdateStory)
				r.Delete("/{id}", storyHandler.DeleteStory)
				r.Post("/{id}/chapters", storyHandler.CreateChapter)
			})

			r.Get("/{id}/comments",storyHandler.Comments)
			r.Post("/{id}/comments",storyHandler.CreateComment)
			r.Get("/{id}", storyHandler.GetStory)
		})

		r.Route("/api/chapters", func(r chi.Router) {
			if authMW != nil {
				r.Use(authMW.RequireAuth)
			}
			r.Get("/{id}", storyHandler.GetChapter)

			r.Group(func(r chi.Router) {
				if authMW != nil {
					r.Use(authMW.RequireAuth)
				}
				r.Put("/{id}", storyHandler.UpdateChapter)
				r.Delete("/{id}", storyHandler.DeleteChapter)
			})
		})
	}

	if storyHandler != nil {
		r.Route("/api/favorites", func(r chi.Router) {
			if authMW != nil {
				r.Use(authMW.RequireAuth)
			}
			r.Get("/", storyHandler.ListFavorites)
			r.Get("/ids", storyHandler.GetFavoriteIDs)
			r.Post("/{id}", storyHandler.AddFavorite)
			r.Delete("/{id}", storyHandler.DeleteFavorite)
		})
	}

	if storyHandler != nil {
		r.Route("/api/follows",func(r chi.Router){
			if authMW!=nil{r.Use(authMW.RequireAuth)}
			r.Get("/ids",storyHandler.FollowIDs)
			r.Post("/{id}",storyHandler.FollowStory)
			r.Delete("/{id}",storyHandler.UnfollowStory)
		})
		r.Route("/api/notifications",func(r chi.Router){
			if authMW!=nil{r.Use(authMW.RequireAuth)}
			r.Get("/",storyHandler.Notifications)
			r.Put("/{id}/read",storyHandler.ReadNotification)
		})
		r.Route("/api/comments",func(r chi.Router){
			if authMW!=nil{r.Use(authMW.RequireAuth)}
			r.Delete("/{id}",storyHandler.DeleteComment)
		})
	}
	if len(ebookHandlers) > 0 && ebookHandlers[0] != nil {
		h := ebookHandlers[0]
		r.Route("/api/ebooks", func(r chi.Router) {
			if authMW != nil {
				r.Use(authMW.RequireAuth)
			}
			r.Get("/", h.List)
			r.Post("/", h.Create)
			r.Get("/{id}", h.Get)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	}
	return r
}
