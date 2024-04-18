package services

import (
	"chulbong-kr/database"
	"chulbong-kr/dto/notification"
	"chulbong-kr/models"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/goccy/go-json"
	"github.com/redis/rueidis"
)

type Notification = models.Notification
type NotificationRedis = notification.NotificationRedis

// PostNotification posts a new notification into the database
func PostNotification(userId, notificationType, title, message string, metadata json.RawMessage) error {
	result, err := database.DB.Exec(
		`INSERT INTO Notifications (UserId, NotificationType, Title, Message, Metadata, Viewed, CreatedAt, UpdatedAt) 
         VALUES (?, ?, ?, ?, ?, FALSE, NOW(), NOW())`,
		userId, notificationType, title, message, metadata,
	)
	if err != nil {
		return err
	}

	notificationId, _ := result.LastInsertId()

	var channelName string
	// Determine the appropriate channel based on notification type
	if notificationType == "Like" || notificationType == "Comment" {
		channelName = "notifications:user:" + userId
	} else {
		channelName = "notifications:broadcast"
	}

	// Publish notification to Redis
	notificationData := NotificationRedis{
		NotificationId:   notificationId,
		UserId:           userId,
		NotificationType: notificationType,
		Title:            title,
		Message:          message,
		Metadata:         metadata,
	}
	jsonData, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}

	err = RedisStore.Do(context.Background(), RedisStore.B().Publish().Channel(channelName).Message(rueidis.BinaryString(jsonData)).Build()).Error()
	if err != nil {
		return err
	}

	return nil
}

// GetNotifications retrieves notifications for a specific user (unviewed)
func GetNotifications(userId string) ([]NotificationRedis, error) {
	var notifications []Notification
	const query = `(SELECT * FROM Notifications 
		WHERE UserId = ? AND Viewed = FALSE 
		ORDER BY CreatedAt DESC)
		UNION ALL
		(SELECT * FROM Notifications 
		WHERE NotificationType IN ('NewMarker', 'System', 'Other') AND Viewed = FALSE 
		ORDER BY CreatedAt DESC)`
	err := database.DB.Select(&notifications, query, userId)
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
				viewed, err := IsNotificationViewed(notif.NotificationId, userId)
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

// markNotificationAsViewed(notification, userId)
func MarkNotificationAsViewed(nid int64, ntype, userId string) {
	if ntype == "Like" || ntype == "Comment" {
		if err := MarkPersonalNotificationViewed(nid, userId); err != nil {
			log.Printf("Error marking personal notification as viewed: %v", err)
		}
	} else {
		MarkBroadcastNotificationViewed(nid, userId)
	}
}

func MarkPersonalNotificationViewed(notificationId int64, userId string) error {
	_, err := database.DB.Exec(`UPDATE Notifications SET Viewed = TRUE WHERE NotificationId = ? AND UserId = ?`, notificationId, userId)
	return err
}

// REDIS
func MarkBroadcastNotificationViewed(notificationId int64, userId string) {
	key := fmt.Sprintf("viewed:notification:%d:%s", notificationId, userId)
	ctx := context.Background()
	RedisStore.Do(ctx, RedisStore.B().Sadd().Key(key).Member("viewed").Build())

	expireCmd := RedisStore.B().Expire().Key(key).Seconds(864000).Build() // 10 days
	RedisStore.Do(ctx, expireCmd)

}

func IsNotificationViewed(notificationId int64, userId string) (bool, error) {
	key := fmt.Sprintf("viewed:notification:%d:%s", notificationId, userId)
	exists, err := RedisStore.Do(context.Background(), RedisStore.B().Sismember().Key(key).Member("viewed").Build()).AsBool()
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SubscribeNotification subscribes to user-specific notification channels and handles messages as they arrive.
func SubscribeNotification(userId string, handleMessage func(msg rueidis.PubSubMessage)) (cancelFunc func(), err error) {
	dedicatedClient, cancel := RedisStore.Dedicate()

	userChannel := "notifications:user:" + userId
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
	if err := json.Unmarshal(data, &notification); err != nil {
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
