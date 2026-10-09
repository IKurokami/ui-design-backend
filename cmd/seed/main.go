package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/repository"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type DBUser struct {
	ID          string
	DisplayName string
	Email       string
}

func main() {
	dbPath := "data/auth.db"
	db, err := repository.NewDB("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Lỗi mở database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Không thể kết nối SQLite: %v", err)
	}

	ctx := context.Background()

	// 1. Tìm kiếm user trong database
	rows, err := db.QueryContext(ctx, "SELECT id, display_name, email FROM users")
	if err != nil {
		log.Fatalf("Lỗi truy vấn users: %v", err)
	}
	defer rows.Close()

	var users []DBUser
	for rows.Next() {
		var u DBUser
		if err := rows.Scan(&u.ID, &u.DisplayName, &u.Email); err == nil {
			users = append(users, u)
		}
	}

	var targetUser DBUser
	if len(users) == 0 {
		// Tạo user mẫu nếu database trống
		targetUser = DBUser{
			ID:          uuid.New().String(),
			DisplayName: "Kuro",
			Email:       "kuro@gmail.com",
		}
		now := time.Now().UTC()
		_, err := db.ExecContext(ctx, `INSERT INTO users (id, display_name, email, password_hash, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)`,
			targetUser.ID, targetUser.DisplayName, targetUser.Email,
			"$2a$10$abcdefghijklmnopqrstuvwxyz123456", now, now,
		)
		if err != nil {
			log.Fatalf("Lỗi tạo user mặc định: %v", err)
		}
		fmt.Printf(" [DB] Đã tạo user mới: %s (%s)\n", targetUser.DisplayName, targetUser.Email)
	} else {
		// Chọn ngẫu nhiên một user có sẵn trong database
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		targetUser = users[r.Intn(len(users))]
		fmt.Printf(" [DB] Đã chọn ngẫu nhiên user có sẵn: %s (ID: %s, Email: %s)\n",
			targetUser.DisplayName, targetUser.ID, targetUser.Email)
	}

	// 2. Cập nhật / tạo hồ sơ cá nhân hoàn chỉnh cho user này
	var profileCount int
	_ = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_profiles WHERE user_id = ?", targetUser.ID).Scan(&profileCount)
	now := time.Now().UTC()
	if profileCount == 0 {
		_, _ = db.ExecContext(ctx, `INSERT INTO user_profiles 
			(id, user_id, avatar_url, banner_url, bio, headline, theme_color, website_url, location, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			uuid.New().String(), targetUser.ID,
			"/mock-images/ln_1.png",
			"/mock-images/ln_18.jpg",
			"Tác giả văn học mạng & Dịch giả Light Novel. Đam mê xây dựng thế giới kỳ ảo, hồi hộp và những câu chuyện tình cảm sâu sắc.",
			"Tác giả & Dịch giả Light Novel",
			"#3B82F6",
			"https://github.com/funny",
			"Việt Nam",
			now, now,
		)
		fmt.Println(" [DB] Đã khởi tạo hồ sơ tác giả đẹp mắt cho user.")
	}

	// 3. Xóa dữ liệu truyện cũ của user này để chèn lại dữ liệu mới sạch sẽ
	_, _ = db.ExecContext(ctx, "DELETE FROM chapters WHERE story_id IN (SELECT id FROM stories WHERE user_id = ?)", targetUser.ID)
	_, _ = db.ExecContext(ctx, "DELETE FROM stories WHERE user_id = ?", targetUser.ID)

	// 4. Tổng hợp tất cả mock data từ các phần
	allStories := append([]StorySeed{}, getStoriesPart1()...)
	allStories = append(allStories, getStoriesPart2()...)
	allStories = append(allStories, getStoriesPart3()...)
	allStories = append(allStories, getStoriesPart4()...)

	fmt.Printf(" [SEED] Bắt đầu thêm %d bộ truyện vào database cho user %s...\n", len(allStories), targetUser.DisplayName)

	totalChapters := 0
	for _, s := range allStories {
		storyID := s.ID
		if storyID == "" {
			storyID = uuid.New().String()
		}

		createdAt := s.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := createdAt.Add(time.Hour * 12)

		_, err := db.ExecContext(ctx, `INSERT INTO stories 
			(id, user_id, title, description, cover_url, category, tags, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			storyID, targetUser.ID, s.Title, s.Description, s.CoverURL, s.Category, s.Tags, s.Status, createdAt, updatedAt,
		)
		if err != nil {
			log.Fatalf("Lỗi chèn truyện %s: %v", s.Title, err)
		}

		for _, ch := range s.Chapters {
			chapterID := uuid.New().String()
			illusJSON, _ := json.Marshal(ch.Illustrations)

			chCreatedAt := createdAt.Add(time.Duration(ch.Number) * time.Hour * 24)
			chUpdatedAt := chCreatedAt

			_, err := db.ExecContext(ctx, `INSERT INTO chapters
				(id, story_id, chapter_number, title, content, illustrations, created_at, updated_at, has_content)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				chapterID, storyID, ch.Number, ch.Title, ch.Content, string(illusJSON), chCreatedAt, chUpdatedAt, readableFlag(ch.Content),
			)
			if err != nil {
				log.Fatalf("Lỗi chèn chương %d của truyện %s: %v", ch.Number, s.Title, err)
			}
			totalChapters++
		}
	}

	fmt.Printf("\n HOÀN TẤT THÀNH CÔNG!\n")
	fmt.Printf("===================================================\n")
	fmt.Printf(" User sở hữu: %s (%s)\n", targetUser.DisplayName, targetUser.Email)
	fmt.Printf(" Tổng số truyện đã thêm:  %d bộ truyện\n", len(allStories))
	fmt.Printf(" Tổng số chương đã thêm:  %d chương đầy đủ nội dung\n", totalChapters)
	fmt.Printf(" Database: %s\n", dbPath)
	fmt.Printf("===================================================\n")
}

func readableFlag(content string) int {
	if domain.HasChapterContent(content) {
		return 1
	}
	return 0
}
