package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"mercado-app/internal/handlers"
	"mercado-app/internal/middleware"
	"mercado-app/internal/repository"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	dbPath := envOr("DB_PATH", "./data/mercado.db")
	dbType := envOr("DB_TYPE", "sqlite")
	port := envOr("PORT", "3001")
	host := envOr("HOST", "0.0.0.0")
	staticDir := envOr("STATIC_DIR", "./static")

	// Set defaults based on dbType, allow override via env vars
	defaultMigrationsDir := "./migrations/sqlite"
	defaultSeedFile := "./seed/seed.sql"
	if dbType == "postgres" {
		defaultMigrationsDir = "./migrations/postgres"
		defaultSeedFile = "./seed/postgres_seed.sql"
	}
	migrationsDir := envOr("MIGRATIONS", defaultMigrationsDir)
	seedFile := envOr("SEED_FILE", defaultSeedFile)

	dsn := dbPath
	if dbType == "postgres" {
		dsn = envOr("DATABASE_URL", "postgres://gastou:gastou@localhost:5432/gastou?sslmode=disable")
	} else {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			log.Fatalf("failed to create data directory: %v", err)
		}
	}

	repo, err := repository.New(dbType, dsn)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}
	defer repo.Close()

	if err := repo.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	empty, err := repo.IsEmpty()
	if err != nil {
		log.Printf("could not check if DB is empty: %v", err)
	} else if empty {
		if err := repo.RunSeed(seedFile); err != nil {
			log.Printf("failed to run seed: %v", err)
		}
	}

	h := handlers.NewHandlers(repo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/categorias", h.ListCategorias)
		r.Get("/categorias/{id}", h.GetCategoria)
		r.Post("/categorias", h.CreateCategoria)
		r.Put("/categorias/{id}", h.UpdateCategoria)
		r.Delete("/categorias/{id}", h.DeleteCategoria)

		r.Get("/produtos", h.ListProdutos)
		r.Get("/produtos/{id}", h.GetProduto)
		r.Post("/produtos", h.CreateProduto)
		r.Put("/produtos/{id}", h.UpdateProduto)
		r.Get("/produtos/{id}/precos", h.GetHistoricoPrecos)

		r.Get("/listas", h.ListListas)
		r.Get("/listas/{id}", h.GetLista)
		r.Post("/listas", h.CreateLista)
		r.Put("/listas/{id}", h.UpdateLista)
		r.Delete("/listas/{id}", h.DeleteLista)

		r.Get("/listas/{id}/itens", h.ListItens)
		r.Post("/listas/{id}/itens", h.AddItem)
		r.Put("/listas/{id}/itens/{itemId}", h.UpdateItem)
		r.Patch("/listas/{id}/itens/{itemId}/check", h.ToggleCheck)
		r.Delete("/listas/{id}/itens/{itemId}", h.DeleteItem)

		r.Get("/compras", h.ListCompras)
		r.Get("/compras/{id}", h.GetCompra)
		r.Post("/compras", h.CreateCompra)
		r.Put("/compras/{id}", h.UpdateCompra)
		r.Delete("/compras/{id}", h.DeleteCompra)

		r.Post("/compras/{id}/itens", h.AddCompraItem)
		r.Put("/compras/{id}/itens/{itemId}", h.UpdateCompraItem)
		r.Delete("/compras/{id}/itens/{itemId}", h.DeleteCompraItem)

		r.Get("/listas/{id}/dashboard/resumo", h.GetResumo)
		r.Get("/listas/{id}/dashboard/comparativo", h.GetComparativo)
		r.Get("/dashboard/evolucao", h.GetEvolucao)
	})

	fileServer := http.FileServer(http.Dir(staticDir))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticDir, r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	log.Printf("Gastou.app running on %s:%s", host, port)
	if err := http.ListenAndServe(host+":"+port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
