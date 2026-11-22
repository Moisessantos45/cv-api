package routes

import (
	"cv_api/config"
	"cv_api/config/db"
	"cv_api/internal/features/auth"
	"cv_api/internal/features/post"
	"cv_api/internal/features/profile"
	"cv_api/internal/features/project"
	"cv_api/internal/features/stream"
	"cv_api/internal/shared/middleware"
	"cv_api/internal/shared/service"
	"cv_api/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

var (
	maker     *utils.PasetoMaker
	authUc    *auth.AuthUseCase
	profileUC profile.ProfileService
	projectUC project.ProjectService
	postUC    post.PostService
	videoUC   stream.VideoService
)

func Init() {
	rd := config.Rdb
	maker = utils.NewPasetoMaker()
	cache := service.NewCacheService(rd)

	aRp := auth.NewPostgresRepository(db.DB)
	authUc = auth.NewAuthUseCase(aRp, rd, maker)

	uRp := profile.NewPostgresRepository(db.DB)
	profileUC = profile.NewProfileUseCase(uRp, authUc, cache)

	pRp := project.NewPostgresRepository(db.DB)
	projectUC = project.NewProjectService(pRp, cache)

	poRp := post.NewPostgresRepository(db.DB)
	postUC = post.NewPostService(poRp, cache, profileUC)

	vRp := stream.NewPostgresRepository(db.DB)
	videoUC = stream.NewVideoUseCase(vRp, cache)
}

func authMiddleware() gin.HandlerFunc {
	return middleware.AuthMiddleware(maker, config.Rdb)
}

func preAuthMiddleware() gin.HandlerFunc {
	return middleware.PreAuthMiddleware(config.Rdb)
}

func idempotencyMiddleware() gin.HandlerFunc {
	return middleware.IdempotencyMiddleware(config.Rdb)
}
