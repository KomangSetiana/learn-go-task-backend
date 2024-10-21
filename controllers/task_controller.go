package controllers

import (
	"net/http"
	"os"
	"strconv"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

func (t *TaskController) Create(c *gin.Context) {
	task := models.Task{}

	errBindJson := c.ShouldBindJSON(&task)

	if errBindJson != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBindJson.Error()})
		return
	}

	errDB := t.DB.Create(&task).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, task)

}

func (t *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")

	task := models.Task{}

	if err := t.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	errDB := t.DB.Delete(&models.Task{}, id).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	if task.Attachment != "" {
		os.Remove("attachments" + task.Attachment)
	}

	c.JSON(http.StatusOK, "Deleted")

}

func (t *TaskController) Submit(c *gin.Context) {
	task := models.Task{}
	id := c.Param("id")
	submitDate := c.PostForm("submitDate")
	file, errFile := c.FormFile("attachment")

	if errFile != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFile.Error()})
		return
	}
	if err := t.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//hapus data lama attachment
	attachment := task.Attachment

	fileInfo, _ := os.Stat("attachments/" + attachment)

	if fileInfo != nil {
		os.Remove("attachments/" + attachment)
	}

	//simpan attachments
	attachment = file.Filename
	errSave := c.SaveUploadedFile(file, "attachments/"+attachment)

	if errSave != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errSave.Error()})
		return
	}

	errDB := t.DB.Where("id=?", id).Updates(models.Task{
		Status:     "Review",
		SubmitDate: submitDate,
		Attachment: attachment,
	}).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, "Submit to review")

}

func (t *TaskController) Reject(c *gin.Context) {
	task := models.Task{}
	id := c.Param("id")
	rejectedDate := c.PostForm("rejectedDate")
	reason := c.PostForm("reason")

	if err := t.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	errDB := t.DB.Where("id=?", id).Updates(models.Task{
		Status:       "Rejected",
		Reason:       reason,
		RejectedDate: rejectedDate,
	}).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, "Rejected")

}

func (t *TaskController) Fix(c *gin.Context) {
	id := c.Param("id")
	revision, err := strconv.Atoi(c.PostForm("revision"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err := t.DB.First(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	errDB := t.DB.Where("id=?", id).Updates(models.Task{
		Status:   "Queue",
		Revision: int8(revision),
	}).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, "Fix to Queue")

}

func (t *TaskController) Approve(c *gin.Context) {
	id := c.Param("id")
	approveDate := c.PostForm("approveDate")

	if err := t.DB.First(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	errDB := t.DB.Where("id=?", id).Updates(models.Task{
		Status:       "Aproved",
		ApprovedDate: approveDate,
	}).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, "Approved")

}

func (t *TaskController) FindById(c *gin.Context) {
	tusk := models.Task{}
	id := c.Param("id")

	if err := t.DB.First(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	errDB := t.DB.Preload("User").Find(&tusk, id).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, tusk)

}

func (t *TaskController) NeedToBeReview(c *gin.Context) {
	tusks := []models.Task{}

	errDB := t.DB.Preload("User").Where("status=?", "review").Order("submit_date ASC").Limit(2).Find(&tusks).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, tusks)

}
func (t *TaskController) ProgressTasks(c *gin.Context) {
	tusks := []models.Task{}
	userId := c.Param("userId")
	errDB := t.DB.Preload("User").Where("(status!=? AND user_id=?) OR (revision!=? AND user_id=?)", "Queue", userId, 0, userId).Order("updated_at DESC").Limit(5).Find(&tusks).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, tusks)

}

func (t *TaskController) Statistic(c *gin.Context) {
	userId := c.Param("userId")

	stat := []map[string]interface{}{}

	errDB := t.DB.Model(models.Task{}).Select("status, count(status) as total").Where("user_id=?", userId).Group("status").Limit(5).Find(&stat).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, stat)

}

func (t *TaskController) FindByUserAndStatus(c *gin.Context) {
	tusks := []models.Task{}
	userId := c.Param("userId")
	status := c.Param("status")
	errDB := t.DB.Where("user_id=? AND status=?", userId, status).Find(&tusks).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, tusks)

}
