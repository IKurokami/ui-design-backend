package main

import "time"

func getStoriesPart2() []StorySeed {
	return []StorySeed{
		{
			ID:          "ln-apothecary",
			Title:       "Dược Sư Tự Sự: Bí Ẩn Cung Đình",
			Description: "Maomao, một thiếu nữ hành nghề dược tại phố đèn đỏ, bị bắt cóc và bán vào Hậu cung làm cung nữ. Vốn muốn sống ẩn dật cho qua ngày, nhưng khi phát hiện con cái các phi tần lần lượt chết yểu vì phấn bột chì độc hại, tính tò mò đã lôi kéo cô vào vòng xoáy âm mưu chốn thâm cung.",
			CoverURL:    "/mock-images/ln_6.png",
			Category:    "mystery",
			Tags:        "Cổ Trang, Phá Án, Y Học, Cung Đình",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -30),
			Chapters: []ChapterSeed{
				{
					Number: 1,
					Title:  "Chương 1: Lời nguyền hoàng tử sơ sinh",
					Illustrations: []string{"/mock-images/ln_37.jpg", "/mock-images/ln_6.png"},
					Content: `<p>Hậu cung của Hoàng đế là một chiếc lồng son tráng lệ nhưng nồng nặc mùi tranh sủng và mưu mô tàn độc. Gần đây, một tin đồn đáng sợ lan truyền như dịch bệnh: những đứa trẻ sơ sinh của các quý phi nương nương lần lượt phát bệnh lạ, da dẻ xanh xao rồi ra đi trong đau đớn.</p>
<p>Người ta đồn rằng đó là lời nguyền của những linh hồn chết oan trong giếng ngọc.</p>
<p>Nhưng Maomao, khi vô tình nhìn thấy vết phát ban trên má của tiểu thái tử con Lệ phi và tiểu công chúa con Quý phi Á Đa, chỉ khẽ nhíu mày. Cô quan sát thứ phấn trắng phủ dày đặc trên gương mặt các bà mẹ trẻ đang khóc lóc thảm thiết:</p>
<p>— <em>"Bạch phấn chứa chì độc... Loại mỹ phẩm thượng hạng được các cung phi ưa chuộng để làm trắng da, nhưng đối với trẻ sơ sinh bú sữa mẹ, đó là chất kịch độc hủy hoại nội tạng."</em></p>
<p>Biết rằng nếu can thiệp vào chuyện hoàng tộc, một cung nữ thấp kém như mình sẽ dễ bị chặt đầu bịt miệng, Maomao quyết định viết một mảnh giấy cảnh báo giấu tên, kẹp vào cành hoa mẫu đơn đặt trước cửa cung của Lệ phi.</p>`,
				},
				{
					Number: 2,
					Title:  "Chương 2: Ánh mắt hồ ly của Tổng quản Nhâm Thị",
					Illustrations: []string{"/mock-images/ln_38.jpg"},
					Content: `<p>Mảnh giấy cảnh báo đã cứu sống tiểu thái tử của Lệ phi, nhưng nó cũng kích hoạt sự tò mò của vị quan giám sát quyền lực nhất Hậu cung – Nhâm Thị (Jinshi).</p>
<p>Với vẻ ngoài tuấn tú phi phàm, làn da mịn như ngọc và nụ cười quyến rũ có thể làm khuynh đảo cả nam lẫn nữ, Nhâm Thị triệu tập toàn bộ cung nữ nghi vấn để kiểm tra chữ viết.</p>
<p>Khi các cung nữ khác đều đỏ mặt ngượng ngùng trước nhan sắc của chàng, Maomao lại chỉ cúi đầu nhìn xuống đất với ánh mắt kinh tởm như nhìn một con ấu trùng nhớp nhúa. Biểu cảm hiếm có đó lập tức lọt vào tầm mắt của vị hoạn quan tinh quái.</p>
<p>— "Tên của ngươi là Maomao đúng không? Từ hôm nay, ngươi được thăng chức làm cung nữ nếm độc kiêm dược sư riêng của cung Ngọc Diệp."</p>
<p>— "Dạ... nông nữ xin tuân mệnh..." – Maomao ngoài mặt cung kính cúi đầu, nhưng trong lòng đang gào thét vì kế hoạch sống ẩn dật yên bình đã tan tành mây khói.</p>`,
				},
				{
					Number: 3,
					Title:  "Chương 3: Bữa tiệc hoa viên và chén súp nấm độc",
					Illustrations: []string{"/mock-images/ln_39.jpg", "/mock-images/ln_40.jpg"},
					Content: `<p>Đại tiệc ngắm hoa mùa thu diễn ra rực rỡ với sự tham gia của Hoàng đế và các đại thần. Là người chịu trách nhiệm nếm độc cho Ngọc Diệp nương nương, Maomao đứng nghiêm trang bên cạnh bàn tiệc thịnh soạn.</p>
<p>Một bát súp cá nấu măng tây bốc khói nghi ngút được đưa lên. Maomao dùng chiếc thìa bạc múc một ngụm nhỏ đưa vào miệng.</p>
<p>Hương vị ngọt ngào của nước dùng vừa lan tỏa, đầu lưỡi cô bỗng cảm nhận được một cảm giác tê rần kỳ lạ, kèm theo vị ngai ngái của nấm vân chi mọc trên gỗ mục có độc.</p>
<p>Thay vì hoảng sợ, đôi mắt Maomao bỗng sáng rực lên với sự hưng phấn tột độ của một kẻ cuồng độc học. Đôi má cô ửng hồng, khóe môi khẽ nhếch lên:</p>
<p>— "Món này... có độc rồi ạ!"</p>
<p>Cả bàn tiệc hoàng gia lập tức hỗn loạn. Các thị vệ tuốt kiếm bao vây khu vực bếp ăn, còn Nhâm Thị đứng từ xa mỉm cười đầy thích thú khi nhìn thấy bản chất thực sự của cô cung nữ quái chiêu.</p>`,
				},
			},
		},
		{
			ID:          "ln-slime",
			Title:       "Tensura: Chuyển Sinh Thành Slime",
			Description: "Mikami Satoru, 37 tuổi, bị đâm chết trên phố và tái sinh ở thế giới khác dưới hình dạng một con Slime dẻo quẹo. Nhờ kỹ năng tối thượng 'Đại Hiền Giả' và 'Kẻ Thôn Phệ', Rimuru Tempest từng bước kết bạn với Bạo Phong Long Veldora và gây dựng liên minh quái vật hùng mạnh.",
			CoverURL:    "/mock-images/ln_42.png",
			Category:    "isekai",
			Tags:        "Isekai, Xây Dựng Vương Quốc, Hài Hước, Quái Vật",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -33),
			Chapters: []ChapterSeed{
				{
					Number: 1,
					Title:  "Chương 1: Thân xác đàn hồi trong hang động ngầm",
					Illustrations: []string{"/mock-images/ln_43.jpg", "/mock-images/ln_42.png"},
					Content: `<p>Cảm giác đầu tiên sau khi chết là không hề có cảm giác đau đớn. Không thấy ánh sáng, không nghe thấy âm thanh, cơ thể mềm nhũn, trơn tuột và đàn hồi như một miếng thạch rau câu.</p>
<p>— <em>[Báo cáo: Kỹ năng độc nhất "Đại Hiền Giả" (Great Sage) đã được kích hoạt thành công.]</em></p>
<p>— <em>[Báo cáo: Kỹ năng độc nhất "Kẻ Thôn Phệ" (Predator) đã sẵn sàng hoạt động.]</em></p>
<p>Mikami Satoru nhận ra mình đã biến thành một con Slime – sinh vật được coi là yếu ớt nhất trong mọi game nhập vai RPG. Nhưng thay vì than vãn, cậu bắt đầu nhảy tưng tưng khắp hang động, nuốt chửng hàng tấn Thảo Dược Ma Lực và Quặng Ma Tinh để tích lũy nguyên liệu và kỹ năng phân giải.</p>
<p>Đột nhiên, một luồng áp lực kinh thiên động địa bao trùm không gian. Một tiếng gầm trầm vang như sấm sét xé toạc màng nhĩ hư ảo:</p>
<p>— "Này nhóc con tí hon, ngươi có nghe thấy tiếng của ta không?!"</p>
<p>Đó là Bạo Phong Long Veldora – một trong bốn Long Chủng tối cao của thế giới, kẻ đã bị Dũng Giả phong ấn suốt ba trăm năm trong bức tường kết giới bất khả xâm phạm.</p>`,
				},
				{
					Number: 2,
					Title:  "Chương 2: Lời hứa giữa Slime và Rồng",
					Illustrations: []string{"/mock-images/ln_44.png"},
					Content: `<p>Hai kẻ cô đơn nhanh chóng trở thành bạn thân qua những câu chuyện về thế giới loài người và văn hóa manga/anime mà Rimuru kể lại. Veldora, dù mang danh xưng rồng hủy diệt tàn bạo, thực chất lại là một "chú rồng tsundere" vô cùng dễ dỗ dành.</p>
<p>— "Ta muốn giải thoát cho cậu, Veldora. Ta sẽ dùng 'Kẻ Thôn Phệ' nuốt trọn cả phong ấn lẫn cậu vào dạ dày ảo vô tận của mình. Trong lúc đó, Đại Hiền Giả sẽ phân tích và phá hủy phong ấn từ bên trong."</p>
<p>— "Gwahaha! Được lắm! Vậy trước khi đi, chúng ta hãy đặt cho nhau một cái tên chung để minh chứng cho tình bạn này!"</p>
<p>Và thế là, cái tên "Rimuru Tempest" ra đời. Một con Slime nhỏ bé đã nuốt trọn một con rồng thần thoại vào bụng, rời khỏi hang động ngầm và bước chân ra thế giới rộng lớn bên ngoài.</p>`,
				},
				{
					Number: 3,
					Title:  "Chương 3: Đặt tên cho tộc Goblin khốn cùng",
					Illustrations: []string{"/mock-images/ln_45.png", "/mock-images/ln_46.png"},
					Content: `<p>Bước ra khỏi hang động, Rimuru chạm trán với một bộ tộc Goblin rách rưới đang quỳ rạp cầu xin sự che chở trước sự tấn công hung hãn của tộc Nha Lang (Sói Băng Răng Nanh).</p>
<p>Sau khi dùng kỹ năng tơ thép và lôi kích hạ gục thủ lĩnh bầy sói, Rimuru thu phục cả hai tộc và quyết định làm một việc mà sau này cậu vô cùng hối hận: đặt tên cho từng con quái vật!</p>
<p>— "Ngươi sẽ là Rigurd... ngươi là Gobta... ngươi là Ranga..."</p>
<p>Trong thế giới này, việc đặt tên cho ma vật sẽ tiêu tốn một lượng ma lực khổng lồ từ người đặt. Sau khi đặt tên cho hơn một trăm con Goblin và Sói, cơ thể Rimuru cạn kiệt mana đến mức bất tỉnh nhân sự suốt ba ngày.</p>
<p>Đến khi tỉnh lại, trước mắt cậu không còn là những con Goblin già nua còm cõi nữa, mà là một dàn trai xinh gái đẹp cơ bắp cuồn cuộn và tràn đầy sinh lực!</p>`,
				},
			},
		},
		{
			ID:          "ln-sao",
			Title:       "Sword Art Online: Thế Giới Kiếm Ảo",
			Description: "Năm 2022, tựa game thực tế ảo VRMMORPG đầu tiên Sword Art Online ra mắt. 10.000 người chơi đăng nhập và bàng hoàng nhận ra họ bị giam cầm trong thế giới ảo: cái chết trong game đồng nghĩa với cái chết ngoài đời thực. Cách duy nhất để thoát ra là chinh phục toàn bộ 100 tầng thành Aincrad.",
			CoverURL:    "/mock-images/ln_4.png",
			Category:    "scifi",
			Tags:        "Thế Giới Ảo, Kiếm Thuật, Hành Động, Tình Cảm",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -36),
			Chapters: []ChapterSeed{
				{
					Number: 1,
					Title:  "Chương 1: Hồi chuông tử thần tại Quảng trường Khởi nguyên",
					Illustrations: []string{"/mock-images/ln_48.jpg", "/mock-images/ln_4.png"},
					Content: `<p>Ngày 6 tháng 11 năm 2022. Mười ngàn người chơi đắm chìm trong vẻ đẹp ngoạn mục của thành phố nổi Aincrad. Nhưng nút "Đăng xuất" (Log Out) trong menu hệ thống đã biến mất không dấu vết.</p>
<p>Chuông giáo đường vang lên gióng giả. Toàn bộ người chơi bị dịch chuyển cưỡng chế về Quảng trường Khởi nguyên ở Tầng 1. Bầu trời chuyển sang màu đỏ rực như máu, và một bóng ma khổng lồ mặc áo choàng đỏ thẫm hiện ra trên không trung – Kayaba Akihiko, người sáng tạo ra SAO.</p>
<p>— "Chào mừng các bạn đến với thế giới của tôi. Đây không còn là một trò chơi đơn thuần. Nếu thanh HP của các bạn về 0, thiết bị NerveGear sẽ phóng ra sóng viba công suất cao phá hủy não bộ của các bạn ngoài đời thực. Và nếu bất kỳ ai bên ngoài cố tháo NerveGear, kết quả cũng tương tự."</p>
<p>Tiếng la hét, khóc lóc và hỗn loạn bùng nổ như ong vỡ tổ. Giữa biển người hoảng loạn, Kirito siết chặt thanh đoản kiếm khởi đầu, lập tức lao nhanh ra khỏi cổng thành trước khi quái vật bị săn hết. Cậu hiểu rằng: muốn sống sót, phải trở nên mạnh mẽ hơn bất kỳ ai!</p>`,
				},
				{
					Number: 2,
					Title:  "Chương 2: Lưỡi kiếm bạc và chiếc bánh mì quết kem",
					Illustrations: []string{"/mock-images/ln_11.png"},
					Content: `<p>Một tháng sau ngày thảm họa bắt đầu, hai ngàn người chơi đã thiệt mạng mà Tầng 1 vẫn chưa thể khai phá. Trong buổi họp mặt chiến lược tại quảng trường phế tích, Kirito chú ý đến một người chơi trùm kín áo choàng xám ngồi cô độc ở hàng ghế sau.</p>
<p>Sau buổi họp, Kirito tiến lại gần, đưa cho người đó một chiếc bánh mì đen cứng ngắc:</p>
<p>— "Này, ăn thử cái này đi."</p>
<p>Cậu mở hũ mứt kem hảo hạng mua từ nhiệm vụ đặc biệt, quết một lớp dày lên chiếc bánh. Người mang áo choàng ngập ngừng cắn một miếng nhỏ. Chiếc mũ trùm đầu khẽ rơi xuống, để lộ ra gương mặt thanh tú với làn da trắng mịn và mái tóc màu hạt dẻ bồng bềnh.</p>
<p>Đó là Asuna – một cô gái ngây thơ chưa từng chơi game, người sau này sẽ trở thành Phó bang chủ của Huyết Minh Kỵ Sĩ Đoàn khét tiếng.</p>`,
				},
				{
					Number: 3,
					Title:  "Chương 3: Quyết đấu Illfang The Kobold Lord",
					Illustrations: []string{"/mock-images/ln_49.png", "/mock-images/ln_50.png"},
					Content: `<p>Căn phòng Boss Tầng 1 tối tăm và ngột ngạt. Tiếng gầm của Illfang làm rung chuyển những cột đá cẩm thạch khổng lồ. Thủ lĩnh đội tiên phong Diabel trúng đòn chí mạng vì sự thay đổi vũ khí bất ngờ của Boss ở giai đoạn cuối.</p>
<p>— "Tất cả lùi lại! Để tôi và Asuna xử lý!" – Kirito hét lớn.</p>
<p>Hai bóng người lao vút lên như tia chớp. Lưỡi kiếm bạc của Asuna đâm liên hoàn với tốc độ không tưởng, trong khi thanh kiếm đen của Kirito tung ra kỹ năng trảm kích xé toạc lớp giáp của tên chúa tể quái thú.</p>
<p>Một đòn quyết định! Thanh HP cuối cùng của Illfang vỡ tan thành muôn vàn mảnh thủy tinh lấp lánh. Tầng 1 đã được giải phóng, nhưng danh xưng "Beater" (Kẻ gian lận thử nghiệm) bắt đầu bám theo bước chân cô độc của Kiếm Sĩ Áo Đen.</p>`,
				},
			},
		},
		{
			ID:          "ln-spice-wolf",
			Title:       "Sói Và Gia Vị: Hành Trình Thương Gia",
			Description: "Kraft Lawrence, một thương nhân lang bạt buôn bán lúa mì và lông thú trên cỗ xe ngựa, tình cờ phát hiện một thiếu nữ có tai và đuôi sói say ngủ trong đống lúa mạch. Đó là Holo - Hiền Lang xứ Yoitsu, vị thần mùa màng bị người dân lãng quên. Cả hai cùng bước vào chuyến hành trình về phương Bắc xa xôi.",
			CoverURL:    "/mock-images/ln_3.jpg",
			Category:    "fantasy",
			Tags:        "Fantasy, Kinh Tế, Du Ngoạn, Người Thú",
			Status:      "completed",
			CreatedAt:   time.Now().AddDate(0, 0, -39),
			Chapters: []ChapterSeed{
				{
					Number: 1,
					Title:  "Chương 1: Đêm hội thu hoạch ở làng Pasloe",
					Illustrations: []string{"/mock-images/ln_52.jpg", "/mock-images/ln_3.jpg"},
					Content: `<p>Hương lúa mạch chín vàng ngập tràn trong gió đêm của làng Pasloe. Sau khi chất đầy các kiện lúa mì lên cỗ xe ngựa của mình, Kraft Lawrence chuẩn bị rời làng để tránh né những lễ nghi ồn ào của ngày hội bội thu.</p>
<p>Khi kéo tấm bạt phủ che mưa trên thùng xe, Lawrence giật mình đánh thót một cái. Nằm cuộn tròn ấm áp giữa đống rơm khô là một thiếu nữ không một mảnh vải che thân, đôi tai sói màu hạt dẻ khẽ giật giật theo nhịp thở và chiếc đuôi dài phủ lông mềm mại đang đung đưa nhịp nhàng.</p>
<p>— "Ưm... ồn ào quá đấy chàng trai trẻ..."</p>
<p>Thiếu nữ khẽ cựa mình, đôi mắt màu hổ phách sắc sảo nhìn thẳng vào người thương nhân đang hóa đá vì kinh ngạc.</p>
<p>— "Ta là Holo, Hiền Lang xứ Yoitsu phương Bắc. Đã hàng trăm năm rồi người dân ở đây không còn cần đến ta nữa... Cậu có muốn đưa ta trở về quê hương không, thương nhân tốt bụng?"</p>`,
				},
				{
					Number: 2,
					Title:  "Chương 2: Cơn sốt giá bạc tại thị trấn Pazzio",
					Illustrations: []string{"/mock-images/ln_53.png"},
					Content: `<p>Bước vào thị trấn giao thương sầm uất Pazzio, Holo đội một chiếc mũ trùm đầu kín mít để che đi đôi tai sói. Với chiếc mũi cực kỳ thính nhạy và trí tuệ hàng trăm năm của một vị thần, nàng ngửi thấy mùi mưu mô trong từng đồng tiền xu bạc đang lưu hành trên thị trường.</p>
<p>— "Này Lawrence, đồng Trenni bạc này... hàm lượng bạc nguyên chất của nó đã bị giảm đi so với các đợt phát hành trước đây."</p>
<p>Lawrence sững sờ. Một âm mưu thao túng tiền tệ quy mô lớn giữa các nghiệp đoàn thương hội và lãnh chúa địa phương đang dần lộ diện. Nếu nắm bắt được thời cơ này, một thương nhân nhỏ bé như anh có thể đổi đời thành triệu phú; nhưng chỉ cần một bước đi sai lầm, tính mạng của cả hai sẽ vùi chôn dưới đáy sông sâu.</p>`,
				},
				{
					Number: 3,
					Title:  "Chương 3: Tiếng hú dưới ánh trăng phương Bắc",
					Illustrations: []string{"/mock-images/ln_54.png", "/mock-images/ln_55.png"},
					Content: `<p>Bị truy sát bởi toán sát thủ của thương hội Medio trong những con ngõ lầy lội, cỗ xe của Lawrence bị dồn vào ngõ cụt. Những mũi giáo sáng loáng chĩa thẳng vào cổ họng anh.</p>
<p>— "Đứng lùi lại phía sau ta, Lawrence."</p>
<p>Holo nuốt trọn một hạt lúa mì thiêng và ngụm rượu vang cay nồng. Cơ thể nàng phát ra ánh sáng chói lòa, xương cốt chuyển dịch với tiếng gầm rung chuyển đất trời. Một con cự lang lông trắng khổng lồ cao bằng cả tòa nhà hai tầng xuất hiện, đôi mắt đỏ rực như lửa thiêu đốt quân thù.</p>
<p>Chỉ với một cú quẫy đuôi, toàn bộ toán sát thủ bị thổi bay như lá rụng mùa thu. Nhưng khi quay lại nhìn Lawrence, nàng khẽ rên rỉ sợ hãi:</p>
<p>— "Chàng... chàng có sợ hình dạng thật của ta không...?"</p>
<p>Lawrence mỉm cười, bình thản bước tới và ôm lấy chiếc mũi to lớn của con sói thần thoại:</p>
<p>— "Dù nàng là ai, nàng vẫn là bạn đồng hành không thể thay thế của ta trên cỗ xe này."</p>`,
				},
			},
		},
		{
			ID:          "ln-monogatari",
			Title:       "Bakemonogatari: Quái Vật Truyện",
			Description: "Araragi Koyomi, một học sinh năm ba cao trung sống sót sau vụ tấn công của ma cà rồng. Một ngày nọ, cô bạn cùng lớp Senjougahara Hitagi rơi từ trên cầu thang xuống và rơi trúng tay cậu - chỉ để lộ ra một bí mật kinh hoàng: cơ thể cô không hề có trọng lượng, nhẹ bẫng tựa như một chiếc lông vũ.",
			CoverURL:    "/mock-images/ln_57.png",
			Category:    "mystery",
			Tags:        "Siêu Nhiên, Huyền Bí, Đối Thoại, Dị Nhân",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -42),
			Chapters: []ChapterSeed{
				{
					Number: 1,
					Title:  "Chương 1: Con cua nặng vỏn vẹn năm kilôgam",
					Illustrations: []string{"/mock-images/ln_58.jpg", "/mock-images/ln_57.png"},
					Content: `<p>Mọi chuyện bắt đầu vào ngày 8 tháng Năm. Khi Araragi Koyomi đang thong thả bước lên cầu thang xoắn ốc của trường trung học Naoetsu, một bóng hình trượt chân và rơi tự do từ tầng ba xuống.</p>
<p>Koyomi theo phản xạ lao tới đỡ lấy. Cậu chuẩn bị tinh thần chịu đựng một chấn động nặng nề của một cơ thể thiếu nữ rơi từ trên cao.</p>
<p>Nhưng khi hai cánh tay ôm trọn lấy cô bạn Senjougahara Hitagi, cảm giác duy nhất cậu nhận được là sự hẫng hụt tột cùng.</p>
<p>Cô gái ấy hoàn toàn không có trọng lượng. Quần áo, da thịt, mái tóc và cả thân xác nàng... tất cả gộp lại không nặng quá năm kilôgam, tương tự như một con búp bê nhồi bông rỗng tuếch!</p>
<p>Và ngay trong khoảnh khắc cậu chưa kịp hiểu chuyện gì đang xảy ra, một chiếc ghim bấm và một chiếc compa kim loại đã gí sát vào khóe mắt và vòm họng của Koyomi:</p>
<p>— "Araragi-kun... nếu cậu dám tiết lộ bí mật này cho bất kỳ ai, chiếc bấm kim này sẽ bấm dính lưỡi của cậu vào hàm ếch đấy."</p>`,
				},
				{
					Number: 2,
					Title:  "Chương 2: Tòa nhà luyện thi bỏ hoang và Oshino Meme",
					Illustrations: []string{"/mock-images/ln_59.jpg"},
					Content: `<p>Để cứu giúp Senjougahara thoát khỏi sự dày vò của quái dị, Koyomi dẫn cô đến một tòa nhà trung tâm luyện thi hoang phế cỏ mọc um tùm. Tại đây, trên sàn gỗ bám đầy bụi, một người đàn ông trung niên mặc áo sơ mi Hawaii sặc sỡ đang ngồi xếp bằng hút thuốc lá.</p>
<p>— "Ồ, có vẻ hôm nay cậu em mang đến một vị khách đặc biệt thú vị nhỉ? Tự hào bản thân tràn trề sức sống ghê ha, Araragi-kun?"</p>
<p>Oshino Meme – một chuyên gia về quái dị và các hiện tượng siêu nhiên – chỉ nhìn lướt qua Hitagi đã thấu rõ căn nguyên:</p>
<p>— "Thần Trọng Lượng (Omoi-gani) đã lấy đi gánh nặng của cô gái này. Nó không chỉ lấy đi cân nặng thể xác, mà lấy đi cả nỗi đau tinh thần, ký ức bi thương về một gia đình tan vỡ mà cô không muốn nhớ tới. Muốn lấy lại trọng lượng, cô bé phải sẵn sàng gánh vác lại toàn bộ nỗi đau đó."</p>`,
				},
				{
					Number: 3,
					Title:  "Chương 3: Bầu trời sao và lời hẹn ước đêm mùa hạ",
					Illustrations: []string{"/mock-images/ln_60.jpg", "/mock-images/ln_61.png"},
					Content: `<p>Sau khi trục xuất thành công con quái vật cua và lấy lại được thân xác trọn vẹn, Senjougahara mời Koyomi đi dã ngoại trên một ngọn đồi ngoại ô bằng chiếc xe cổ của cha mình.</p>
<p>Trải tấm bạt mỏng trên bãi cỏ dưới bầu trời mùa hạ không một gợn mây, Hitagi nằm ngửa mặt lên nhìn dải ngân hà lấp lánh như dát bạc:</p>
<p>— "Araragi-kun, hãy nhìn ngôi sao Deneb, Altair và Vega kìa. Đó là Tam giác Mùa Hè."</p>
<p>Nàng khẽ nghiêng đầu sang, ánh mắt dịu dàng phản chiếu ánh sao trời:</p>
<p>— "Tôi không có tiền tài, không có gia đình hoàn hảo, tính cách thì độc đoán và tàn nhẫn... Thứ duy nhất tôi có thể trao trọn cho cậu, chỉ là bầu trời đầy sao tuyệt đẹp này, và tình cảm chân thành nhất của tôi mà thôi."</p>`,
				},
			},
		},
	}
}
