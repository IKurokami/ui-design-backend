package main

import "time"

func getStoriesPart3() []StorySeed {
	return []StorySeed{
		{
			ID:          "ln-konosuba",
			Title:       "KonoSuba: Mở Ra Thế Giới Tuyệt Vời!",
			Description: "Satou Kazuma chết vì trụy tim sau khi tưởng lầm máy kéo là xe tải cán người. Gặp nữ thần Aqua kênh kiệu ở thế giới bên kia, cậu quyết định lôi cổ ả đi cùng sang thế giới fantasy. Nhưng thay vì một hành trình anh hùng vĩ đại, Kazuma phải đối mặt với một tổ đội toàn những kẻ lập dị vô dụng.",
			CoverURL:    "/mock-images/ln_63.png",
			Category:    "isekai",
			Tags:        "Isekai, Hài Bựa, Tổ Đội Thảm Họa, Phép Thuật",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -45),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Cái chết hài hước nhất lịch sử thế giới",
					Illustrations: []string{"/mock-images/ln_64.jpg", "/mock-images/ln_63.png"},
					Content: `<p>Mọi chuyện bắt đầu bằng một hành động tưởng như nghĩa hiệp: Satou Kazuma lao ra đẩy một cô gái thoát khỏi chiếc "xe tải" đang lao tới. Nhưng thực tế cay đắng là chiếc xe đó chỉ là một chiếc máy cày mini di chuyển với tốc độ rùa bò, cô gái chẳng hề gặp nguy hiểm, còn Kazuma thì chết vì... sốc và trụy tim do quá sợ hãi!</p>
<p>Tại phòng chờ luân hồi, nữ thần tóc xanh Aqua ôm bụng cười ngặt nghẽo, chế giễu cái chết nhục nhã của cậu trước mặt toàn thể linh hồn.</p>
<p>— "Hahaha! Ngươi thậm chí còn tè dầm ra quần trước khi chết nữa chứ! Các bác sĩ ở bệnh viện đã cười vào mặt gia đình ngươi đấy!"</p>
<p>Cơn thịnh nộ bùng nổ. Khi Aqua kiêu ngạo cho phép Kazuma chọn mang theo một món đồ tùy ý sang thế giới mới (một thánh kiếm, ma pháp tối thượng hay kỹ năng thần thánh), cậu chỉ thẳng tay vào mặt ả:</p>
<p>— "Món đồ ta chọn mang theo... chính là CÔ! Lên đường thôi mụ nữ thần chết tiệt!"</p>
<p>Một vòng tròn triệu hồi rực sáng kéo cả hai rơi thẳng xuống thị trấn khởi đầu Axel trong tiếng la hét thất thanh của Aqua.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Tổ đội thảm họa toàn diện",
					Illustrations: []string{"/mock-images/ln_9.png"},
					Content: `<p>Để sinh tồn và trả nợ cho những thiệt hại do Aqua gây ra, Kazuma đăng tuyển đồng đội tại Hội Mạo Hiểm Giả. Và đây là những thành viên "ưu tú" gia nhập:</p>
<p>Thứ nhất: Megumin – một Ma Pháp Sư Hồng Ma Tộc mang phong cách chunibyo cuồng bạo. Nàng sở hữu Ma Pháp Bộc Phá (Explosion) có sức tàn phá khủng khiếp san phẳng cả ngọn núi, nhưng sau khi bắn đúng một phát thì mana cạn kiệt và nằm liệt đất như một khúc gỗ mục.</p>
<p>Thứ hai: Darkness – một Thập Tự Kỵ Sĩ (Crusader) mang dòng máu quý tộc lộng lẫy, sở hữu lớp giáp dày cộm và chỉ số phòng ngự vô địch. Nhưng vấn đề là kiếm thuật của cô nàng chém một trăm nhát thì trượt cả một trăm, và sở thích bệnh hoạn của cô là được quái vật... đánh đập và lăng nhục!</p>
<p>— "Trời ơi! Tôi muốn một dàn harem phiêu lưu mộng mơ cơ mà, sao lại thành ra trại tâm thần thế này?!" – Kazuma ôm đầu gào khóc trong tuyệt vọng.</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Thảm họa ếch khổng lồ mùa sinh sản",
					Illustrations: []string{"/mock-images/ln_65.png", "/mock-images/ln_66.png"},
					Content: `<p>Nhiệm vụ đầu tiên tiêu diệt loài Cự Oa (Ếch khổng lồ) ở đồng cỏ ngoại ô nhanh chóng biến thành thảm họa bi hài. Aqua hăng hái lao lên dùng "Thần Quyền" nhưng bị con ếch nuốt chửng nửa người, chất nhầy nhơn nhớt bao phủ khắp người khiến ả khóc lóc cầu cứu thảm thiết.</p>
<p>Megumin tung ra một phát Bộc Phá nổ tung con ếch thứ hai rồi nằm gục xuống chờ chết, chỉ còn lại Kazuma phải dùng kỹ năng móc túi (Steal) và cung tên lén lút kéo từng người tháo chạy thục mạng.</p>
<p>Trở về quán trọ với thân thể dính đầy nhớt ếch tanh tưởi, cả bọn ngồi quanh chiếc bàn gỗ chia nhau vài đồng tiền thưởng còm cõi. Thế nhưng, giữa những tiếng cãi vã ầm ĩ và trò hề ngớ ngẩn, Kazuma bất giác nhận ra: cuộc sống mới này, dù khốn khổ trăm bề, nhưng lại vui vẻ hơn nhiều so với chuỗi ngày cô độc giam mình trong phòng ngủ kiếp trước.</p>`,
				},
			},
		},
		{
			ID:          "ln-overlord",
			Title:       "Overlord: Kẻ Thống Trị Lăng Mộ",
			Description: "Vào thời khắc tựa game DMMO-RPG Yggdrasil đóng máy chủ vĩnh viễn, Momonga - hội trưởng bang Ainz Ooal Gown - quyết định ở lại đến giây phút cuối cùng. Khi đồng hồ điểm 00:00, lăng mộ ngầm Nazarick cùng toàn bộ các NPC hộ vệ trung thành bất ngờ bị dịch chuyển đến một thế giới hoàn toàn xa lạ.",
			CoverURL:    "/mock-images/ln_7.png",
			Category:    "isekai",
			Tags:        "Isekai, Bá Đạo, Chiến Thuật, Nhân Vật Phản Diện",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -48),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Giây phút chia tay và tiếng chuông điểm nửa đêm",
					Illustrations: []string{"/mock-images/ln_68.png", "/mock-images/ln_7.png"},
					Content: `<p>23:59:50... 23:59:55...</p>
<p>Momonga ngồi trên Chiếc Ngai Tráng Lệ của Đại Lăng Tộ Nazarick, tay cầm Quyền Trượng Thần Thánh của Ainz Ooal Gown – minh chứng cho vinh quang của 41 vị Đấng Tối Cao. Cậu nhắm mắt lại, chờ đợi khoảnh khắc trò chơi đóng cửa và bị đẩy về với thực tại khắc nghiệt của một nhân viên công sở quèn.</p>
<p>00:00:00... 00:00:03...</p>
<p>Không có màn hình đăng xuất. Đồng hồ hệ thống đã biến mất. Nhưng điều kinh khủng hơn là các giác quan của cậu hoàn toàn chân thực: mùi hương thanh khiết của hoa cỏ, làn gió nhẹ lay động chiếc áo choàng ma thuật, và đặc biệt là giọng nói ngọt ngào, tràn ngập tình cảm của Albedo – vị Tổng quản Thủ vệ tầng:</p>
<p>— "Ngài có điều gì phiền lòng sao, Momonga-sama?"</p>
<p>Cậu giật mình nhìn xuống: bàn tay cậu không còn là da thịt con người, mà là những đốt xương trắng bệch phát ra luồng ma lực tăm tối của một Chúa Tể Hư Không (Overlord)!</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Lòng trung thành tuyệt đối của các Thủ vệ tầng",
					Illustrations: []string{"/mock-images/ln_69.png"},
					Content: `<p>Tại Đấu trường Tầng 6, toàn bộ các Thủ vệ tầng hùng mạnh nhất tập hợp quỳ gối dưới chân Momonga: Ma cà rồng Shalltear Bloodfallen, quái vật côn trùng Cocytus, song sinh hắc ám Aura và Mare, ác quỷ xảo quyệt Demiurge và nữ chúa Albedo.</p>
<p>— "Thưa Đấng Tối Cao duy nhất còn ở lại với chúng thần! Nguyện dâng hiến toàn bộ sinh mạng và linh hồn để phục vụ ý chí của ngài!"</p>
<p>Nhìn vào ánh mắt cuồng tín của những tạo vật do chính tay các người bạn cũ tạo nên, Momonga phải dốc hết sức kìm nén sự lo sợ của một người bình thường. Kỹ năng thụ động "Khống Chế Cảm Xúc" của tộc Undead lập tức kích hoạt, tỏa ra một làn sáng xanh xoa dịu tâm trí cậu, biến cậu trở thành một bậc đế vương lạnh lùng, uy nghiêm và tàn nhẫn đúng nghĩa.</p>
<p>— "Từ hôm nay, hãy gọi ta là Ainz Ooal Gown! Mục tiêu của chúng ta là biến cái tên này thành huyền thoại bất diệt của thế giới này!"</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Tàn sát giải cứu làng Carne",
					Illustrations: []string{"/mock-images/ln_70.png", "/mock-images/ln_71.png"},
					Content: `<p>Qua quả cầu viễn thị ma thuật, Ainz quan sát một ngôi làng nhỏ bé đang bị quân đội đế quốc thảm sát tàn bạo. Nhớ lại ân nghĩa của người bạn Touch Me – người luôn chủ trương cứu giúp kẻ yếu – Ainz quyết định can thiệp.</p>
<p>Bước qua cánh cổng không gian "Gate", vị Chúa Tể Xương hiện ra giữa ngôi làng ngập tràn khói lửa. Đứng trước đội kỵ sĩ thiết giáp hung hãn, Ainz chỉ giơ nhẹ một ngón tay xương xẩu lên:</p>
<p>— "Phép thuật bậc 9: Nắm Đấm Tử Thần (Grasp Heart)!"</p>
<p>Trái tim của tên chỉ huy kỵ sĩ lập tức bị bóp nát trong hư vô, gã ngã gục xuống đất chết không kịp ngáp. Những con Tử Kỵ Sĩ (Death Knight) khổng lồ được triệu hồi, mở ra một màn săn lùng kinh hoàng quét sạch toàn bộ quân xâm lược trong chớp mắt.</p>
<p>Đó là bước chân đầu tiên của vị Vua Phù Thủy trên con đường chinh phục và định hình lại bản đồ của Tân Thế Giới.</p>`,
				},
			},
		},
		{
			ID:          "ln-zeroth-maria",
			Title:       "Hộp Rỗng Và Maria Thứ 0",
			Description: "Hoshino Kazuki coi trọng cuộc sống bình thường của mình hơn bất cứ điều gì. Nhưng sự xuất hiện của Otonashi Maria vào ngày 2 tháng 3 đã xáo trộn tất cả. Cô nữ sinh chuyển trường bước lên bục giảng và tuyên bố: 'Tôi ở đây để tiêu diệt cậu, Hoshino Kazuki'. Đó là sự khởi đầu của vòng lặp thứ 27.756.",
			CoverURL:    "/mock-images/ln_73.jpg",
			Category:    "mystery",
			Tags:        "Tâm Lý Tối Tăm, Vòng Lặp, Hồi Hộp, Bí Ẩn",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -51),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Ngày 2 tháng 3 lần thứ 27.756",
					Illustrations: []string{"/mock-images/ln_74.png", "/mock-images/ln_73.jpg"},
					Content: `<p>Mọi thứ diễn ra y hệt như mọi ngày: tiếng chuông báo thức reo lúc 6h30, bánh mì nướng bơ tỏi, cuộc trò chuyện quen thuộc với cô bạn thời thơ ấu Kokone trên đường đến trường.</p>
<p>Cho đến khi giáo viên chủ nhiệm dẫn một học sinh mới bước vào lớp học vào ngày 2 tháng 3.</p>
<p>Cô gái ấy sở hữu vẻ đẹp phi thực tế: làn da trắng xanh xao, đôi mắt đen sâu thẳm không chút gợn sóng và mái tóc cắt ngắn ngang vai. Nàng cầm viên phấn viết lên bảng ba chữ: "Otonashi Maria".</p>
<p>Sau đó, nàng nhìn thẳng vào mắt Hoshino Kazuki, chỉ tay vào ngực cậu trước sự sững sờ của toàn bộ bạn bè:</p>
<p>— "Tôi ở đây để tiêu diệt cậu, Hoshino Kazuki."</p>
<p>— "Hả?! Cậu đang nói cái quái gì thế?!"</p>
<p>Maria bước từng bước chậm rãi đến trước bàn cậu, cúi xuống thì thầm bằng một tông giọng lạnh toát chứa đựng hàng vạn năm mệt mỏi:</p>
<p>— "Đây là lần thứ 27.756 tôi chuyển đến ngôi trường này. Và cậu chính là kẻ sở hữu chiếc 'Hộp Khước Từ' giam cầm chúng ta trong ngày 2 tháng 3 vô tận này."</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Bản chất của chiếc Hộp Ma Thuật",
					Illustrations: []string{"/mock-images/ln_75.jpg"},
					Content: `<p>Thế giới bên ngoài lớp học bắt đầu rạn nứt như kính vỡ. Những người xung quanh cử động như những con rối cơ học lặp đi lặp lại một hành vi duy nhất.</p>
<p>Trên sân thượng lộng gió của trường học, Maria giải thích cho Kazuki về sự tồn tại của "Hộp" (Box) – những tạo vật siêu nhiên sinh ra từ ước vọng điên cuồng của con người do một thực thể tên là "O" ban phát.</p>
<p>— "Một ai đó trong số bạn bè của cậu đã ước rằng 'tháng ngày hạnh phúc này sẽ không bao giờ kết thúc'. Chiếc hộp đã đáp ứng ước nguyện đó bằng cách xóa sổ tương lai và nhốt tất cả vào một ngày duy nhất. Cứ mỗi lần kim đồng hồ điểm 24:00, toàn bộ ký ức sẽ bị xóa sạch, ngoại trừ tôi – người mang chiếc Hộp Số 0 có khả năng chống lại sự can thiệp."</p>
<p>Kazuki kinh hoàng khi nhận ra những người bạn thân thiết nhất của mình đều đang che giấu những góc tối bệnh hoạn trong tâm hồn mà cậu chưa từng nhận ra.</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Đập tan ảo ảnh",
					Illustrations: []string{"/mock-images/ln_76.jpg", "/mock-images/ln_77.png"},
					Content: `<p>Cuộc đối đầu tâm lý nghẹt thở diễn ra giữa Kazuki, Maria và chủ nhân thực sự của chiếc hộp giam cầm thời gian. Từng lớp mặt nạ bị bóc trần trong đau đớn và nước mắt.</p>
<p>Để phá hủy chiếc hộp, Kazuki phải chấp nhận một sự thật nghiệt ngã: cuộc sống bình yên mà cậu hằng nâng niu thực chất chỉ là một ảo tưởng được dựng nên bởi sự trốn tránh thực tại.</p>
<p>— "Maria, hãy đâm vào tim tôi đi!"</p>
<p>Lưỡi dao xé toạc không gian, âm thanh vỡ vụn giòn tan vang lên như hàng triệu tinh thể tuyết tan rã dưới ánh bình minh. Vòng lặp thứ 27.756 chính thức khép lại. Ngày 3 tháng 3 đã đến, nhưng một cuộc chiến tinh thần mới tàn khốc hơn vừa mới bắt đầu.</p>`,
				},
			},
		},
		{
			ID:          "ln-three-days",
			Title:       "Ba Ngày Hạnh Phúc",
			Description: "Kusunoki, một thanh niên 20 tuổi tuyệt vọng và cạn kiệt tiền bạc, tìm đến một cửa hàng mua bán tuổi thọ, sức khỏe và thời gian. Khi biết phần đời còn lại của mình chỉ đáng giá vỏn vẹn 300.000 yên (chưa đầy 30 yên mỗi ngày), cậu quyết định bán đi 30 năm tuổi thọ và giữ lại 3 tháng cuối cùng cùng cô gái giám sát Miyagi.",
			CoverURL:    "/mock-images/ln_79.jpg",
			Category:    "romcom",
			Tags:        "Đoản Thiên, Ý Nghĩa Sống, Tình Cảm, Cảm Xúc",
			Status:      "completed",
			CreatedAt:   time.Now().AddDate(0, 0, -55),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Ba mươi năm đổi lấy ba mươi vạn yên",
					Illustrations: []string{"/mock-images/ln_80.jpg", "/mock-images/ln_79.jpg"},
					Content: `<p>Trong căn phòng u ám trên tầng bốn của một tòa nhà cũ kỹ không biển hiệu, người nhân viên thẩm định đẩy một tờ giấy về phía Kusunoki:</p>
<p>— "Sau khi xem xét kỹ lưỡng về thành tựu, mức độ đóng góp cho xã hội và sự hạnh phúc mà cậu mang lại cho người khác trong tương lai... Giá trị tuổi thọ của cậu được ấn định là 10.000 yên mỗi năm. Cậu còn ba mươi năm mười tháng cuộc đời. Nếu bán ba mươi năm, cậu sẽ nhận được đúng ba mươi vạn yên tiền mặt."</p>
<p>Ba mươi vạn yên (tương đương khoảng 60 triệu đồng Việt Nam). Một con số rẻ mạt đến nực cười. Một cuộc đời ba mươi năm không bằng giá tiền của một chiếc máy tính xách tay cao cấp.</p>
<p>Nhưng Kusunoki không hề do dự. Cậu ký tên vào bản giao kèo, nhận chiếc phong bì dày cộm và giữ lại vỏn vẹn ba tháng cuối cùng để sống.</p>
<p>Và người đồng hành cùng cậu trong những ngày ngắn ngủi này là Miyagi – một nữ giám sát viên có nhiệm vụ theo dõi để đảm bảo cậu không làm điều gì điên rồ ảnh hưởng đến trật tự xã hội trước khi chết.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Người giám sát vô hình Miyagi",
					Illustrations: []string{"/mock-images/ln_81.jpg"},
					Content: `<p>Miyagi là một người vô hình đối với cả thế giới, chỉ có người sắp chết như Kusunoki mới có thể nhìn thấy và trò chuyện với cô. Nàng luôn mặc một bộ vest đen trang nhã, mái tóc buộc đuôi ngựa gọn gàng và đôi mắt chứa đựng nỗi buồn man mác.</p>
<p>Những ngày đầu tiên, Kusunoki cố gắng bù đắp cho những ước mơ dở dang: ăn những món ngon đắt tiền, đi du lịch đến những nơi chưa từng đến. Nhưng cảm giác trống rỗng và cô độc vẫn không ngừng bủa vây lấy cậu.</p>
<p>Chỉ khi cậu bắt đầu trò chuyện với Miyagi về tuổi thơ, về những sai lầm và nỗi ân hận trong quá khứ, nhịp tim của cả hai mới tìm thấy sự đồng điệu kỳ lạ:</p>
<p>— "Miyagi này, tại sao cô lại làm công việc giám sát viên này?"</p>
<p>— "Bởi vì mẹ tôi đã bán hết tuổi thọ của bà ấy để trả nợ và bỏ rơi tôi. Tôi phải làm công việc này suốt ba mươi năm để trả món nợ tuổi thọ đó của bà..."</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Ba ngày hạnh phúc nhất trần gian",
					Illustrations: []string{"/mock-images/ln_82.jpg"},
					Content: `<p>Khi thời hạn ba tháng gần cạn, Kusunoki đã phát hiện ra một bí mật: nếu cậu dùng số tiền còn lại và ba ngày tuổi thọ cuối cùng của mình để chuộc lại tự do cho Miyagi, cô sẽ được sống một cuộc đời trọn vẹn của một con người bình thường.</p>
<p>Cậu đã bán nốt toàn bộ thời gian còn lại, chỉ giữ lại đúng ba ngày cuối cùng.</p>
<p>Và điều kỳ diệu đã xảy ra: giá trị của ba ngày đó, nhờ vào tình yêu thương và sự hy sinh cao cả dành cho Miyagi, đã tăng vọt lên hàng tỷ yên – trở thành những ngày tháng quý giá nhất mà cửa hàng từng thẩm định.</p>
<p>Dưới ánh hoàng hôn mùa hè rực rỡ bên bờ biển, hai người nắm chặt tay nhau bước đi trên cát mịn:</p>
<p>— "Kusunoki... ba ngày này, chúng ta hãy sống thật hạnh phúc nhé?"</p>
<p>— "Ừ, ba ngày này... sẽ đáng giá hơn cả một đời người."</p>`,
				},
			},
		},
	}
}
