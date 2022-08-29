package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListCuboids(c *gin.Context) {
	var cuboids []models.Cuboid
	if r := db.CONN.Find(&cuboids); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, cuboids)
}

func GetCuboid(c *gin.Context) {
	cuboidId := c.Param("cuboidId")

	cuboid, err := findCuboidById(cuboidId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cuboid)
}

func findCuboidById(id string) (models.Cuboid, error) {
	var cuboid models.Cuboid
	if r := db.CONN.Preload("Bag").First(&cuboid, id); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return cuboid, errors.New("not found")
		} else {
			return cuboid, errors.New(r.Error.Error())
		}
	}
	return cuboid, nil
}

func findBagById(id uint) (models.Bag, error) {
	var bag models.Bag
	if r := db.CONN.First(&bag, id); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return bag, errors.New("not found")
		} else {
			return bag, errors.New(r.Error.Error())
		}
	}
	return bag, nil
}

func CreateCuboid(c *gin.Context) {
	cuboid := models.Cuboid{}

	if err := c.BindJSON(&cuboid); err != nil {
		return
	}

	bag, err := findBagById(cuboid.BagID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if bag.Disable {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		return
	}

	if cuboid.PayloadVolume() > bag.AvailableVolume() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		return
	}

	if r := db.CONN.Create(&cuboid); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	c.JSON(http.StatusCreated, &cuboid)
}

func UpdateCuboid(c *gin.Context) {
	cuboidId := c.Param("cuboidId")

	cuboid, err := findCuboidById(cuboidId)
	if err != nil {
		return
	}

	if err := c.BindJSON(&cuboid); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cuboid.Bag.Disable {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		return
	}

	if cuboid.PayloadVolume() > cuboid.Bag.AvailableVolume() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		return
	}

	if r := db.CONN.Updates(&cuboid); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, &cuboid)
}

func DeleteCuboid(c *gin.Context) {
	cuboidId := c.Param("cuboidId")

	cuboid, err := findCuboidById(cuboidId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if r := db.CONN.Delete(&cuboid); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Done"})
}
