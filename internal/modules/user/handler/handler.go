package handler

import (
	"net/http"
	"strconv"

	"erp/internal/modules/user/model"
	"erp/internal/modules/user/service"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler 用户处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建用户处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "用户注册信息"
// @Success 200 {object} response.Response{data=model.Response} "注册成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 409 {object} response.Response{error=string} "用户已存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /user/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	user, err := h.service.Register(c, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("用户注册成功", user))
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录并获取JWT令牌
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.LoginRequest true "用户登录信息"
// @Success 200 {object} response.Response{data=model.LoginResponse} "登录成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "用户名或密码错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	loginResp, err := h.service.Login(c, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("登录成功", loginResp))
}

// GetProfile godoc
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.Response} "获取成功"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "用户不存在"
// @Router /user/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.service.GetProfile(c, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", user))
}

// UpdateProfile godoc
// @Summary 更新用户信息
// @Description 更新当前登录用户的个人信息
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body model.UpdateProfileRequest true "用户更新信息"
// @Success 200 {object} response.Response{data=model.Response} "更新成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /user/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	user, err := h.service.UpdateProfile(c, userID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("用户信息更新成功", user))
}

// ChangePassword godoc
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body model.ChangePasswordRequest true "密码修改信息"
// @Success 200 {object} response.Response{data=string} "密码修改成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "原密码错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /user/change_password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	err := h.service.ChangePassword(c, userID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("密码修改成功", "密码已更新"))
}

// GetUsers godoc
// @Summary 获取用户列表
// @Description 获取所有用户列表（需要管理员权限）
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=model.UserListResponse} "获取成功"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 403 {object} response.Response{error=string} "权限不足"
// @Router /user/list [get]
func (h *Handler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, err := h.service.GetUsers(c, page, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", users))
}
