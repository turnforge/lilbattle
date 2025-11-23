package gormbe

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GenId struct {
	Class      string `gorm:"primaryKey"`
	Id         string `gorm:"primaryKey"`
	CreatedAt  time.Time
	VerifiedAt time.Time
	Released   bool
}

func randid() string {
	max_id := int64(math.Pow(36, 8))
	randval := rand.Int63() % max_id
	return strconv.FormatInt(randval, 36)
}

func GetID(storage *gorm.DB, cls string, id string) *GenId {
	return nil
}

// Generate 1 New ID
func NewID(storage *gorm.DB, cls string, existingId string) string {
	// if an existing Id was provided then see if it exists
	existingId = strings.ToLower(existingId)
	if existingId != "" {
		gid := GetID(storage, cls, existingId)
		if gid != nil {
			return ""
		}
		gid = &GenId{Id: existingId, Class: cls, CreatedAt: time.Now()}
		err := storage.Create(gid).Error
		if err == nil {
			return gid.Id
		} else {
			log.Println("Error blocking existing ID: ", err)
			return ""
		}
	}
	for i := range 5 {
		gid := GenId{Id: randid(), Class: cls, CreatedAt: time.Now()}
		err := storage.Create(gid).Error
		if err == nil {
			return gid.Id
		} else {
			log.Printf("Trial: %d, ID Create Error: %v", i, err)
		}
	}
	return ""
}

/**
 * Create N IDs in batch.
 */
func NewIDs(storage *gorm.DB, cls string, numids int) (out []string) {
	for i := range numids {
		for {
			gid := GenId{Id: randid(), Class: cls, CreatedAt: time.Now()}
			err := storage.Create(gid).Error
			if err != nil {
				log.Println("ID Create Error: ", i, err)
			} else {
				out = append(out, gid.Id)
				break
			}
		}
	}
	return
}

func VerifyID(storage *gorm.DB, cls string, id string) error {
	var gid GenId
	id = strings.ToLower(id)
	err := storage.First(&gid, "cls = ? and id = ?", cls, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	gid.VerifiedAt = time.Now()
	return storage.Updates(gid).Error
}

func ReleaseID(storage *gorm.DB, cls string, id string) error {
	var gid GenId
	id = strings.ToLower(id)
	err := storage.First(&gid, "cls = ? and id = ?", cls, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	gid.Released = true
	gid.VerifiedAt = time.Now()
	return storage.Updates(gid).Error
}
