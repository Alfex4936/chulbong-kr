package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Alfex4936/chulbong-kr/dto/notification"
	"github.com/Alfex4936/chulbong-kr/model"
	"github.com/jmoiron/sqlx"

	sonic "github.com/bytedance/sonic"
	"github.com/redis/rueidis"
)

type NotificationService struct {
	DB    *sqlx.DB
	Redis *RedisService
}

func NewNotificationService(db *sqlx.DB, redis *RedisService) *NotificationService {
	return &NotificationService{
		DB:    db,
		Redis: redis,
	}
}

type (
	Notification      = model.Notification
	NotificationRedis = notification.NotificationRedis
)

// PostNotification posts a new notification into the database
func (s *NotificationService) PostNotification(userID, notificationType, title, message string, metadata json.RawMessage) error {
	result, err := s.DB.Exec(
		`INSERT INTO Notifications (UserId, NotificationType, Title, Message, Metadata, Viewed, CreatedAt, UpdatedAt) 
         VALUES (?, ?, ?, ?, ?, FALSE, NOW(), NOW())`,
		userID, notificationType, title, message, metadata,
	)
	if err != nil {
		return err
	}

	notificationId, _ := result.LastInsertId()

	var channelName string
	// Determine the appropriate channel based on notification type
	if notificationType == "Like" || notificationType == "Comment" {
		channelName = "notifications:user:" + userID
	} else {
		channelName = "notifications:broadcast"
	}

	// Publish notification to Redis
	notificationData := NotificationRedis{
		NotificationId:   notificationId,
		UserId:           userID,
		NotificationType: notificationType,
		Title:            title,
		Message:          message,
		Metadata:         metadata,
	}
	jsonData, err := sonic.Marshal(notificationData)
	if err != nil {
		return err
	}

	err = s.Redis.Core.Client.Do(context.Background(), s.Redis.Core.Client.B().Publish().Channel(channelName).Message(rueidis.BinaryString(jsonData)).Build()).Error()
	if err != nil {
		return err
	}

	return nil
}

// GetNotifications retrieves notifications for a specific user (unviewed)
func (s *NotificationService) GetNotifications(userID string) ([]NotificationRedis, error) {
	var notifications []Notification
	const query = `(SELECT * FROM Notifications 
		WHERE UserId = ? AND Viewed = FALSE 
		ORDER BY CreatedAt DESC)
		UNION ALL
		(SELECT * FROM Notifications 
		WHERE NotificationType IN ('NewMarker', 'System', 'Other') AND Viewed = FALSE 
		ORDER BY CreatedAt DESC)`
	err := s.DB.Select(&notifications, query, userID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make([]NotificationRedis, len(notifications))
	errors := make(chan error, len(notifications))

	for i, n := range notifications {
		wg.Add(1)
		go func(idx int, notif Notification) {
			defer wg.Done()
			if notif.NotificationType == "Like" || notif.NotificationType == "Comment" {
				if !notif.Viewed {
					results[idx] = mapToNotificationRedis(notif)
				}
			} else {
				viewed, err := s.IsNotificationViewed(notif.NotificationId, userID)
				if err != nil {
					errors <- err
					return
				}
				if !viewed {
					results[idx] = mapToNotificationRedis(notif)
				}
			}
		}(i, n)
	}
	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed checking notification viewed status")
	}

	filteredNotifications := []NotificationRedis{}
	for _, res := range results {
		if res.NotificationId != 0 {
			filteredNotifications = append(filteredNotifications, res)
		}
	}
	return filteredNotifications, nil
}

// markNotificationAsViewed(notification, userID)
func (s *NotificationService) MarkNotificationAsViewed(nid int64, ntype, userID string) {
	if ntype == "Like" || ntype == "Comment" {
		if err := s.MarkPersonalNotificationViewed(nid, userID); err != nil {
			log.Printf("Error marking personal notification as viewed: %v", err)
		}
	} else {
		s.MarkBroadcastNotificationViewed(nid, userID)
	}
}

func (s *NotificationService) MarkPersonalNotificationViewed(notificationId int64, userID string) error {
	_, err := s.DB.Exec(`UPDATE Notifications SET Viewed = TRUE WHERE NotificationId = ? AND UserId = ?`, notificationId, userID)
	return err
}

// REDIS
func (s *NotificationService) MarkBroadcastNotificationViewed(notificationId int64, userID string) {
	key := fmt.Sprintf("viewed:notification:%d:%s", notificationId, userID)
	ctx := context.Background()
	s.Redis.Core.Client.Do(ctx, s.Redis.Core.Client.B().Sadd().Key(key).Member("viewed").Build())

	expireCmd := s.Redis.Core.Client.B().Expire().Key(key).Seconds(864000).Build() // 10 days
	s.Redis.Core.Client.Do(ctx, expireCmd)

}

func (s *NotificationService) IsNotificationViewed(notificationId int64, userID string) (bool, error) {
	key := fmt.Sprintf("viewed:notification:%d:%s", notificationId, userID)
	exists, err := s.Redis.Core.Client.Do(context.Background(), s.Redis.Core.Client.B().Sismember().Key(key).Member("viewed").Build()).AsBool()
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SubscribeNotification subscribes to user-specific notification channels and handles messages as they arrive.
func (s *NotificationService) SubscribeNotification(userID string, handleMessage func(msg rueidis.PubSubMessage)) (cancelFunc func(), err error) {
	dedicatedClient, cancel := s.Redis.Core.Client.Dedicate()

	userChannel := "notifications:user:" + userID
	broadcastChannel := "notifications:broadcast"

	wait := dedicatedClient.SetPubSubHooks(rueidis.PubSubHooks{
		OnMessage: handleMessage, // pass the handling function directly
	})

	// Subscribe to the user-specific notifications channel
	if err := dedicatedClient.Do(context.Background(), dedicatedClient.B().Subscribe().Channel(userChannel, broadcastChannel).Build()).Error(); err != nil {
		cancel()
		log.Printf("🎯 Subscription failed: %v", err)
		return nil, err
	}

	// Return a cancel function and handle the subscription lifecycle in a non-blocking way
	return func() {
		cancel() // Clean up the dedicated connection when needed
		go func() {
			err := <-wait // Wait for the subscription to close and log any errors
			if err != nil {
				log.Printf("🎯 Subscription closed with error: %v", err)
			}
		}()
	}, nil
}

func HandleNotification(data []byte) {
	var notification NotificationRedis
	if err := sonic.Unmarshal(data, &notification); err != nil {
		fmt.Printf("Error decoding notification: %v\n", err)
		return
	}

	fmt.Printf("Notification: %+v\n", notification)
}

func mapToNotificationRedis(n Notification) NotificationRedis {
	return NotificationRedis{
		NotificationId:   n.NotificationId,
		UserId:           n.UserId,
		NotificationType: n.NotificationType,
		Title:            n.Title,
		Message:          n.Message,
		Metadata:         n.Metadata,
	}
}
