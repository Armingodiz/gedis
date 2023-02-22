package kvcontroller

import (
	"gedis/models"
	"gedis/services/ttlManager"
	"gedis/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KvController struct {
	Store      store.Store
	TtlManager ttlManager.Manager
}

func (u *KvController) CreateKv() gin.HandlerFunc {
	return func(c *gin.Context) {
		var kv models.KeyValue
		err := c.ShouldBindJSON(&kv)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		kv.UserName = c.MustGet("user_name").(string)
		err = u.Store.CreateKv(kv)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = u.TtlManager.SetTtl(kv.UserName, kv.Key, -1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "key value set successfully!"})
	}
}

func (u *KvController) GetKvs() gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.MustGet("user_name").(string)
		kvs, err := u.Store.GetKvs(userName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		available := []models.KeyValue{}
		for _, kv := range kvs {
			ttlCheck, err := u.TtlManager.CheckKey(userName, kv.Key)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if ttlCheck {
				available = append(available, kv)
			}
		}
		c.JSON(http.StatusOK, available)
	}
}

func (u *KvController) GetKv() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		userName := c.MustGet("user_name").(string)
		ttlCheck, err := u.TtlManager.CheckKey(userName, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ttlCheck {
			c.JSON(http.StatusBadRequest, gin.H{"error": "key not found or expired"})
			return
		}
		Kv, err := u.Store.GetKv(userName, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, Kv)
	}
}

func (u *KvController) SetTtl() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ttl models.TtlRequest
		err := c.ShouldBindJSON(&ttl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userName := c.MustGet("user_name").(string)
		err = u.TtlManager.SetTtl(userName, ttl.Key, ttl.Ttl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Ttl set successfully!"})
	}
}
