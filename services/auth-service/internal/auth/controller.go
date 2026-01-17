package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc: svc}
}

// Detailed methods

// Register godoc
// @Summary      User Registration
// @Description  Register a new user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      auth.RegisterRequest  true  "User Registration Data"
// @Success      201   {object}  auth.RegisterResponse
// @Failure      400   {object}  map[string]string
// @Router       /auth/register [post]
func (ctrl *Controller) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ctrl.svc.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ctrl.svc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
