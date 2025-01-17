package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config database configurations
type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

// Database database helper
type Database struct {
	*gorm.DB
}

// NewDBConnection new database connection
func NewDBConnection(c *Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&AccessKeys{},
		&Users{},
		&PlayURLCache{},
		&THSeasonCache{},
		&THSubtitleCache{},
	)

	return &Database{db}, err
}

// GetKey get access key data
func (db *Database) GetKey(key string) (*AccessKeys, error) {
	var data AccessKeys
	err := db.Where(&AccessKeys{Key: key}).First(&data).Error
	return &data, err
}

// InsertOrUpdateKey insert or update access key data
func (db *Database) InsertOrUpdateKey(key string, uid int) (int64, error) {
	data, err := db.GetKey(key)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := db.Create(&AccessKeys{Key: key, UID: uid})
		return result.RowsAffected, result.Error
	} else if err != nil {
		return 0, err
	}
	result := db.Model(data).Updates(AccessKeys{Key: key, UID: uid})
	return result.RowsAffected, result.Error
}

// GetUser get user from uid
func (db *Database) GetUser(uid int) (*Users, error) {
	var data Users
	err := db.Where(&Users{UID: uid}).First(&data).Error
	return &data, err
}

// InsertOrUpdateUser insert or update user data
func (db *Database) InsertOrUpdateUser(uid int, name string, vipDueDate time.Time) (int64, error) {
	data, err := db.GetUser(uid)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := db.Create(&Users{UID: uid, Name: name, VIPDueDate: vipDueDate})
		return result.RowsAffected, result.Error
	} else if err != nil {
		return 0, err
	}
	result := db.Model(data).Updates(Users{UID: uid, Name: name, VIPDueDate: vipDueDate})
	return result.RowsAffected, result.Error
}

// GetPlayURLCache get play url caching with device type, area, cid or episode ID
func (db *Database) GetPlayURLCache(deviceType DeviceType, area Area, isVIP bool, cid int, episodeID int) (*PlayURLCache, error) {
	var data PlayURLCache
	err := db.Where(&PlayURLCache{
		DeviceType: deviceType,
		Area:       area,
		IsVip:      isVIP,
		CID:        cid,
		EpisodeID:  episodeID,
	}).First(&data).Error
	return &data, err
}

// InsertOrUpdatePlayURLCache insert or update play url cache data
func (db *Database) InsertOrUpdatePlayURLCache(deviceType DeviceType, area Area, isVIP bool, cid int, episodeID int, jsonData string) (int64, error) {
	data, err := db.GetPlayURLCache(deviceType, area, isVIP, cid, episodeID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := db.Create(&PlayURLCache{
			DeviceType: deviceType,
			Area:       area,
			IsVip:      isVIP,
			CID:        cid,
			EpisodeID:  episodeID,
			JSONData:   jsonData,
		})
		return result.RowsAffected, result.Error
	} else if err != nil {
		return 0, err
	}
	result := db.Model(data).Updates(PlayURLCache{
		DeviceType: deviceType,
		Area:       area,
		CID:        cid,
		EpisodeID:  episodeID,
		JSONData:   jsonData,
	})
	return result.RowsAffected, result.Error
}

// GetTHSeasonCache get season api cache from season id
func (db *Database) GetTHSeasonCache(seasonID int) (*THSeasonCache, error) {
	var data THSeasonCache
	err := db.Where(&THSeasonCache{SeasonID: seasonID}).First(&data).Error
	return &data, err
}

// InsertOrUpdateTHSeasonCache insert or update season api cache
func (db *Database) InsertOrUpdateTHSeasonCache(seasonID int, jsonData string) (int64, error) {
	data, err := db.GetTHSeasonCache(seasonID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := db.Create(&THSeasonCache{SeasonID: seasonID, JSONData: jsonData})
		return result.RowsAffected, result.Error
	} else if err != nil {
		return 0, err
	}
	result := db.Model(data).Updates(THSeasonCache{SeasonID: seasonID, JSONData: jsonData})
	return result.RowsAffected, result.Error
}

// GetTHSubtitleCache get th subtitle api cache from season id
func (db *Database) GetTHSubtitleCache(episodeID int) (*THSubtitleCache, error) {
	var data THSubtitleCache
	err := db.Where(&THSubtitleCache{EpisodeID: episodeID}).First(&data).Error
	return &data, err
}

// InsertOrUpdateTHSubtitleCache insert or update th subtitle api cache
func (db *Database) InsertOrUpdateTHSubtitleCache(episodeID int, jsonData string) (int64, error) {
	data, err := db.GetTHSubtitleCache(episodeID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := db.Create(&THSubtitleCache{EpisodeID: episodeID, JSONData: jsonData})
		return result.RowsAffected, result.Error
	} else if err != nil {
		return 0, err
	}
	result := db.Model(data).Updates(THSubtitleCache{EpisodeID: episodeID, JSONData: jsonData})
	return result.RowsAffected, result.Error
}
