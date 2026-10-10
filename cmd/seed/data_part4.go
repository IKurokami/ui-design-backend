package main

import "time"

func getStoriesPart4() []StorySeed {
	return []StorySeed{
		{
			ID:          "book-thanh-vien",
			Title:       "Thánh Viên Vị Hoa",
			Description: "Thiếu nữ tóc hồng với đôi cánh thiên thần trong vắt giữa biển trời xanh biếc. Thánh viên vị hoa nở rộ mỗi chu kỳ ngàn năm một lần, báo hiệu sự trở lại của kỷ nguyên thánh tích cổ xưa.",
			CoverURL:    "/mock-images/ln_12.png",
			Category:    "Kỳ Ảo",
			Tags:        "Xuyên Không, Phép Thuật, Thiên Sứ",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -2),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Tiếng chuông nơi thánh điện hoa trắng",
					Illustrations: []string{"/mock-images/ln_12.png"},
					Content: `<p>Khi tia nắng ban mai đầu tiên chiếu rọi qua vòm kính màu của Thánh Điện Vị Hoa, hàng ngàn đóa hoa thủy tinh nở rộ đồng loạt, phát ra những thanh âm trong trẻo tựa chuông gió thiên đình.</p>
<p>Thiếu nữ đứng giữa tâm điện khẽ vươn đôi cánh lông vũ hồng phấn. Đôi mắt nàng đọng lại ánh sao của kỷ nguyên trước:</p>
<p>— "Một ngàn năm... cuối cùng phong ấn của Đấng Sáng Tạo cũng đã lung lay."</p>
<p>Dưới chân tháp ngọc, các kỵ sĩ dòng Đền đang quỳ rạp, thành kính dâng lên thanh gươm ánh sáng thiêng liêng.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Cánh chim xé toạc tầng mây",
					Illustrations: []string{"/mock-images/ln_13.jpg"},
					Content: `<p>Cơn bão ma lực từ Vực Thẳm Đen bất ngờ ập đến, xé toạc bầu trời thanh bình phía trên lãnh địa. Những sinh vật hắc ám gầm rú lao xuống từ khe nứt không gian.</p>
<p>Không hề nao núng, thánh nữ vỗ cánh bay vút lên bầu trời cao vạn trượng. Đôi tay nàng vẽ nên trận đồ quang minh khổng lồ, thanh tẩy toàn bộ tà khí chỉ trong một cái chớp mắt.</p>`,
				},
			},
		},
		{
			ID:          "book-tieu-dien",
			Title:       "Tiếu Diện Giáo Thụ",
			Description: "Bí ẩn đằng sau nụ cười của vị giáo sư trẻ tuổi tại học viện pháp thuật tối cao. Không một ai biết được danh tính thực sự và mục đích sâu xa phía sau những bài giảng kỳ dị của ngài.",
			CoverURL:    "/mock-images/ln_1.png",
			Category:    "Bí Ẩn",
			Tags:        "Đấu Trí, Học Viện, Pháp Thuật",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -5),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Nụ cười nửa miệng trong buổi khảo hạch",
					Illustrations: []string{"/mock-images/ln_1.png"},
					Content: `<p>Giáo sư William luôn xuất hiện với một nụ cười hòa nhã trên môi, cặp kính gọng vàng che giấu đi sự sắc sảo tột cùng trong ánh mắt. Trong buổi khảo hạch đầu năm của Viện Pháp Thuật Hoàng Gia, ngài đặt lên bàn một chiếc hộp sắt rỉ sét:</p>
<p>— "Ai có thể mở chiếc hộp này trong vòng năm phút mà không dùng ma lực, người đó sẽ trở thành đệ tử chân truyền của ta."</p>
<p>Tất cả các học viên thiên tài đều dùng đủ mọi ma thuật bạo lực nhưng chiếc hộp không hề suy chuyển. Cho đến khi một học viên vô danh bước lên và nhận ra chiếc chìa khóa thực chất đã được vị giáo sư giấu ngay trong chiếc nhẫn bạc ngài đang đeo.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Cuốn sổ tay bọc da cổ tự",
					Illustrations: []string{"/mock-images/ln_22.jpg"},
					Content: `<p>Đêm khuya tại phòng làm việc riêng, ánh nến lung linh hắt bóng giáo sư William lên kệ sách ngập tràn tài liệu cấm kỵ. Ngài mở cuốn sổ tay bọc da rồng cổ xưa, những dòng chữ ma pháp màu huyết sắc hiện lên:</p>
<p>— "Bàn cờ đã được sắp đặt. Hoàng tộc và Giáo hội... hãy xem ai sẽ là kẻ ngã xuống trước."</p>`,
				},
			},
		},
		{
			ID:          "book-lao-su",
			Title:       "Lão Sư Vạn Năng",
			Description: "Người thầy dẫn dắt thế hệ đệ tử vượt qua ranh giới sinh tử. Tưởng chừng là kẻ lười biếng ham ăn mê ngủ, nhưng mỗi khi ra tay lại làm chấn động cả Cửu Châu thiên hạ.",
			CoverURL:    "/mock-images/ln_2.png",
			Category:    "Hài Hước",
			Tags:        "Dị Giới, Tu Tiên, Hài Bựa",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -8),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Thu nhận tên đồ đệ phế vật đầu tiên",
					Illustrations: []string{"/mock-images/ln_2.png"},
					Content: `<p>Lâm Bắc nằm dài trên ghế đu tre, tay phe phẩy chiếc quạt nan rách, miệng ngậm cọng cỏ nhìn tên thiếu niên đang quỳ gối dập đầu trước cổng viện:</p>
<p>— "Sư phụ! Đan điền của con bị vỡ, kinh mạch tắc nghẽn, ba đại môn phái đều đuổi con ra ngoài! Xin sư phụ rủ lòng từ bi thu nhận con!"</p>
<p>Lâm Bắc ngáp dài một cái, ném cho hắn một chiếc xẻng sắt:</p>
<p>— "Được rồi, vào xới hết mười mẫu ruộng rau phía sau núi đi. Nhớ là xới bằng tay không dùng linh lực, khi nào xong thì được ăn cơm chiều."</p>
<p>Tên đồ đệ ngơ ngác, không hề hay biết rằng mỗi nhát xẻng xới xuống đều đang vô tình khai mở ba mươi sáu đại huyệt đạo chí tôn!</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Đạp nát sơn môn đối thủ chỉ bằng một chiếc dép",
					Illustrations: []string{"/mock-images/ln_11.png"},
					Content: `<p>Khi Chưởng môn Thanh Vân Tông kéo hàng trăm đệ tử tinh anh đến đòi san bằng ngọn núi của Lâm Bắc, vị lão sư lười biếng chỉ bực bội vì bị phá hỏng giấc ngủ trưa.</p>
<p>Cởi chiếc dép rơm dưới chân, hắn nhẹ nhàng quăng ra một cái. Chiếc dép xé gió bay vút, hóa thành một quả núi khổng lồ đập tan trận pháp hộ tông của đối thủ thành tro bụi trong tiếng kinh hoàng bạt vía của bách tính xung quanh.</p>`,
				},
			},
		},
		{
			ID:          "book-bach-hop",
			Title:       "Bách Hợp Viên Thánh Á",
			Description: "Khoảnh khắc yên bình giữa hoa viên rực rỡ và những câu chuyện chưa kể. Nơi hai tâm hồn đồng điệu tìm thấy nơi trú ẩn sau những giông bão của cuộc đời.",
			CoverURL:    "/mock-images/ln_3.jpg",
			Category:    "Đời Thường",
			Tags:        "Học Đường, Chữa Lành, Thanh Xuân",
			Status:      "draft",
			CreatedAt:   time.Now().AddDate(0, 0, -12),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Giọt sương mai trên cánh hoa bách hợp",
					Illustrations: []string{"/mock-images/ln_3.jpg"},
					Content: `<p>Khu vườn bí mật phía sau nhà kính thực vật của trường nữ sinh Thánh Á luôn ngập tràn mùi hương ngọt ngào của những đóa bách hợp trắng muốt.</p>
<p>An và Chi tình cờ gặp nhau vào một buổi chiều mưa rào bất chợt. Dưới mái hiên kính đọng đầy những hạt nước trong veo, giai điệu đàn violin du dương đã kết nối hai trái tim cô đơn tìm về bến đỗ bình yên.</p>`,
				},
			},
		},
		{
			ID:          "book-anh-dao",
			Title:       "Anh Đào Mỹ Đại",
			Description: "Giai điệu mùa xuân dịu dàng dưới tàng cây anh đào bung nở. Tình yêu tuổi học trò ngọt ngào như viên kẹo bông gòn giữa tiết trời se lạnh.",
			CoverURL:    "/mock-images/ln_10.png",
			Category:    "Tình Cảm",
			Tags:        "Rom-Com, Thanh Xuân, Học Đường",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -14),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Cơn mưa cánh hoa đầu mùa",
					Illustrations: []string{"/mock-images/ln_10.png"},
					Content: `<p>Con đường dốc dẫn lên trường trung học ngập tràn sắc hồng phớt của hoa anh đào vào ngày đầu tiên của năm học mới. Gió thổi tung những cánh hoa lướt qua tà áo đồng phục tinh khôi.</p>
<p>Một cú va chạm bất ngờ ở ngã rẽ khiến cả chồng sách vở rơi tung tóe. Khi ngẩng đầu lên, ánh mắt hai người chạm nhau, một cảm giác xao xuyến lạ kỳ bỗng chốc đọng lại mãi trong ký ức thanh xuân tươi đẹp.</p>`,
				},
			},
		},
		{
			ID:          "book-huyet-da",
			Title:       "Huyết Dạ Ma Nữ",
			Description: "Nữ hoàng bóng đêm thức giấc sau ngàn năm phong ấn dưới lòng đất. Với ánh mắt đỏ rực như máu và lưỡi hái bóng đêm, nàng bắt đầu cuộc thanh trừng những kẻ phản bội đế chế xưa.",
			CoverURL:    "/mock-images/ln_7.png",
			Category:    "Dark Fantasy",
			Tags:        "Bá Đạo, Ma Thuật, Huyết Tộc",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -17),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Xiềng xích phong ấn vỡ vụn",
					Illustrations: []string{"/mock-images/ln_7.png"},
					Content: `<p>Dưới đáy sâu của lăng mộ ngàn năm, nơi không một tia sáng mặt trời nào có thể chạm tới, cỗ quan tài bằng đá hắc diệp thạch bỗng rung chuyển dữ dội.</p>
<p>Những sợi xích ma pháp khắc ấn cổ tự lần lượt đứt gãy với tiếng nổ vang rền như sấm sét. Nắp quan tài bật tung, một đôi mắt đỏ thẫm rực sáng trong bóng đêm đặc quánh:</p>
<p>— "Ngàn năm đằng đẵng... Giờ là lúc vương triều bóng đêm trở lại làm chủ thế giới này!"</p>`,
				},
			},
		},
	}
}
