package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TestimonyController struct {
	ts *service.TestimonyService
}

func NewTestimonyController(ts *service.TestimonyService) *TestimonyController {
	return &TestimonyController{
		ts: ts,
	}
}

// SetTestimony
//
// @Summary			Set testimony
// @Description		Set testimony
// @Tags			testimonies
// @Produce			json
// @Security		BearerToken
// @Router			/testimonies		[post]
// @Param			data	body		dto.Testimony	true	"body"
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (t *TestimonyController) SetTestimony(ctx *gin.Context) {
	var body dto.Testimony

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)
	if err := t.ts.SetTestimony(ctx.Request.Context(), claims.Id, body); err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.TestimonyEmptyField) {
			ctx.JSON(http.StatusNoContent, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "testimony successfully created",
	})
}

// GetTestimony
//
// @Summary			Get testimony
// @Description		Get testimony
// @Tags			testimonies
// @Produce			json
// @Router			/testimonies		[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (t *TestimonyController) GetTestimony(ctx *gin.Context) {
	result, err := t.ts.GetTestimony(ctx.Request.Context())

	if err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get testimony",
	})
}
