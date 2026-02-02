package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func AddFakeEnv(r chi.Router) {

	r.Get("/.env", func(w http.ResponseWriter, req *http.Request) {
		
		w.Header().Set("ETag", `"version-1"`)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		
		const port = "PORT=3000\n"
		const web = "WEB_URL=https://lesi.dev\n"
		const dbUrl = "DATABASE_URL=postgresql://postgres.diwjafhwdeuagjjnaanlkl:f0d9adwdwam1m2c3cfff5b6f6kmkm8@aws-0-eu-west-2.pooler.supabase.com:5432/postgres\n"
		const supabase = "NEXT_PUBLIC_SUPABASE_URL=https://diwjafhwdeuagjjnaanlkl.supabase.co\n"
		const supabaseAnon = "NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlIiwiaWF0IjoxNzcwMDM4MTMwLCJleHAiOjE5Mjc3MTgxMzAsImlzX3JlYWwiOiBmYWxzZSwgIndhdGNoIjogIk9uZSBQaWVjZSJ9.-VgB0srFkndb2-2cnz-Moc72N9oyRtBuL3A9NKehZX0\n"
		const supabaseService = "NEXT_PUBLIC_SUPABASE_SERVICE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic2VydmljZV9yb2xlIiwiaXNzIjoic3VwYWJhc2UiLCJpYXQiOjE3NzAwMzgxMzAsImV4cCI6MTkyNzcxODEzMCwgImlzX3JlYWwiOiBmYWxzZSwgIndhdGNoIjogIk9uZSBQaWVjZSJ9.jnjpsm3JYLgvaZYvLHrGvPFKdmXg6TzKQoLFsjtaYEs\n"
		const jwtSecret = "JWT_SECRET=l+BAvNSZdek2yrj2KzXbse0HpZNQQa8xGd3lSDGJ\n"
		const pgMeta = "PG_META_CRYPTO_KEY=GBAw63g3QD7zEkRdFeH2z6cfAszc50j6\n"

		const enc = "ENCRYPTION_KEY=444ac301bd419132a6b1ce33ef2211fbe78ae0d3cdba3f80517b296cb74a3a0a\n"
		const openAi = "OPENAI_API_KEY=sk-hunter2\n"

		combined := port + web + dbUrl + supabase + supabaseAnon + supabaseService + jwtSecret + pgMeta + enc + openAi
		utils.TextResponse(w, http.StatusOK, combined)
	})
}