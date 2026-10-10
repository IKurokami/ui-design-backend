package main

import "time"

func getStoriesPart1() []StorySeed {
	return []StorySeed{
		{
			ID:          "ln-mushoku",
			Title:       "Mushoku Tensei: Thất Nghiệp Chuyển Sinh",
			Description: "Một gã NEET 34 tuổi bị đuổi khỏi nhà sau khi cha mẹ qua đời. Khi cố gắng cứu những thiếu niên khỏi chiếc xe tải đang lao tới, hắn đã chết và tái sinh vào một thế giới kiếm và ma thuật với thân phận Rudeus Greyrat. Mang theo ký ức kiếp trước, cậu thề sẽ sống một cuộc đời nghiêm túc không hối tiếc.",
			CoverURL:    "/mock-images/ln_12.png",
			Category:    "isekai",
			Tags:        "Isekai, Kỳ Ảo, Phiêu Lưu, Phép Thuật",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -15),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Tái sinh tại thế giới xa lạ",
					Illustrations: []string{"/mock-images/ln_13.jpg", "/mock-images/ln_12.png"},
					Content: `<p>Ánh nắng sớm xuyên qua rèm cửa sổ bằng vải thô, rọi những vệt sáng vàng óng lên sàn gỗ mộc mạc. Gió phương Bắc mang theo hơi thở se lạnh của mùa đông vừa tan, len qua khe cửa khẽ làm lay động trang giấy còn vương mùi mực tươi.</p>
<p>— <em>"Rốt cuộc thì... đây không phải là một giấc mơ."</em></p>
<p>Rudeus khẽ đưa bàn tay nhỏ bé của một đứa trẻ sơ sinh lên trước mắt mình. Cảm giác da thịt chân thực, nhịp đập nơi lồng ngực và luồng sinh khí ấm áp đang cuộn trào khắp các mạch máu không thể nào làm giả được. Kiếp trước, hắn từng là một kẻ thất bại thảm hại, ba mươi tư năm chôn vùi thanh xuân trong bốn bức tường u ám, không dám bước chân ra thế giới bên ngoài, chịu đựng sự khinh miệt của gia đình và người đời cho đến tận giây phút bị chiếc xe tải nghiền nát...</p>
<p>Nhưng giờ đây, số phận đã trao cho hắn một cơ hội thứ hai. Một thế giới hoàn toàn mới, nơi ma thuật không phải là điều huyễn tưởng nằm trong những trang tiểu thuyết.</p>
<h3>1. Cảm nhận ma lực sơ khai</h3>
<p>Ở thế giới này, dòng chảy mana len lỏi trong từng giọt nước, từng ngọn cỏ và cả trong hơi thở. Nhắm mắt lại, cậu tập trung ý thức vào khoảng không giữa tâm trí. Ban đầu chỉ là một đốm sáng li ti lấp lánh như đom đóm, nhưng khi cậu thử điều tức và dẫn dắt nó theo những kiến thức kiếp trước về năng lượng học, đốm sáng bỗng bùng lên thành một vòng xoáy dịu dàng.</p>
<p>— <em>"Nước ơi, hãy tụ lại trong lòng bàn tay ta..."</em></p>
<p>Một quả cầu nước trong vắt từ từ ngưng tụ, lơ lửng ngay phía trên đầu ngón tay cậu. Nó dao động theo từng nhịp thở, phản chiếu ánh bình minh rạng rỡ. Không cần xướng ngôn ma pháp, hoàn toàn bằng khả năng kiểm soát ma lực thuần túy!</p>
<blockquote>"Nếu ta bắt đầu luyện tập ma lực ngay từ khi còn là một đứa trẻ, lượng mana tối đa của cơ thể sẽ không ngừng mở rộng theo cấp số nhân."</blockquote>
<p>Một nụ cười khẽ nở trên môi đứa trẻ. Cuộc đời mới này, cậu nhất định sẽ dốc hết toàn lực sống một cách xứng đáng!</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Vị gia sư nhỏ tuổi Roxy Migurdia",
					Illustrations: []string{"/mock-images/ln_14.png"},
					Content: `<p>Vào một buổi sáng đầu hè ngát hương cỏ non, cánh cổng gỗ của ngôi nhà nhà Greyrat mở ra. Một thiếu nữ vóc dáng nhỏ nhắn bước vào sân. Mái tóc màu xanh lam buộc bím hai bên buông xõa qua vai, chiếc áo choàng pháp sư rộng thùng thình và cây quyền trượng dài gần bằng chiều cao cơ thể.</p>
<p>— "Tôi là Roxy Migurdia, theo lời mời của ngài Paul đến làm gia sư ma pháp cho thiếu gia Rudeus."</p>
<p>Giọng nói nàng bình thản, mang chút lạnh lùng đặc trưng của tộc Migurd – một chủng tộc quỷ nhân nổi tiếng với khả năng giao tiếp bằng thần giao cách cảm và tuổi thọ hàng trăm năm. Khi ánh mắt nàng chạm vào Rudeus, một thoáng hoài nghi hiện lên trong đáy mắt:</p>
<p>— "Cậu bé này chỉ mới năm tuổi thôi sao? Ngài Paul, ngài không đùa đấy chứ?"</p>
<p>— "Hahaha! Con trai ta là thiên tài đấy cô giáo trẻ ạ! Cứ kiểm tra nó xem!" – Paul khoanh tay cười lớn đầy tự đắc.</p>
<p>Roxy thở dài, khẽ nâng cây quyền trượng lên. Một vòng tròn ma pháp xanh lam hiện ra dưới chân nàng:</p>
<p>— "Thủy Cầu Thuật (Water Ball)!"</p>
<p>Quả cầu nước bắn vút ra, đập tan bia gỗ cách đó ba mươi mét. Roxy quay sang nhìn Rudeus với ánh mắt ra hiệu: "Đến lượt cậu đấy, làm theo xướng ngôn của tôi."</p>
<p>Nhưng thay vì đọc thần chú, Rudeus chỉ khẽ chớp mắt. Trong tích tắc, không cần xướng ngôn, một mũi tên nước xoáy rít lên với vận tốc kinh người, xuyên thủng tấm bia đá kiên cố phía sau bia gỗ, để lại một lỗ hổng trơn nhẵn bốc hơi nghi ngút.</p>
<p>Sân vườn hoàn toàn rơi vào im lặng tuyệt đối. Chiếc mũ phù thủy của Roxy rơi tụt xuống đất.</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Lời từ biệt và bầu trời giông bão",
					Illustrations: []string{"/mock-images/ln_15.png"},
					Content: `<p>Hai năm thấm thoát trôi qua dưới sự chỉ dạy tận tình của Roxy. Rudeus đã học được cách thi triển ma pháp cấp Thánh và vượt qua bài kiểm tra tốt nghiệp trong cơn mưa rào dữ dội ngoài đồng hoang.</p>
<p>Đêm trước ngày khởi hành, Roxy ngồi trên bậu cửa sổ, ngắm nhìn ánh trăng vằng vặc in bóng trên dòng sông ngoại ô làng Buena.</p>
<p>— "Rudeus, thế giới ngoài kia rộng lớn hơn những gì cậu tưởng tượng rất nhiều. Đừng bao giờ tự mãn với tài năng của mình."</p>
<p>— "Sư phụ... người thật sự phải đi sao?"</p>
<p>Roxy mỉm cười, một nụ cười hiền từ mà trước đây nàng chưa từng bộc lộ. Nàng tháo chiếc mặt dây chuyền khắc huy hiệu tộc Migurd trên cổ mình, cẩn thận đeo vào cổ Rudeus:</p>
<p>— "Ta chẳng còn gì để dạy cho cậu nữa. Cậu đã là một pháp sư vĩ đại hơn ta rồi. Hãy giữ lấy vật này, nó sẽ bảo hộ cho cậu trên những chặng đường tương lai."</p>
<p>Sáng hôm sau, khi tia nắng đầu tiên chưa kịp xua tan màn sương, bóng dáng nhỏ bé mang chiếc áo choàng xanh đã khuất dần phía chân trời. Rudeus siết chặt mặt dây chuyền trong tay, lòng dâng lên một quyết tâm sắt đá. Một ngày nào đó, cậu sẽ trở thành một người đàn ông đáng tin cậy để gặp lại người thầy đầu tiên của mình.</p>`,
				},
				{
					Number:        4,
					Title:         "Chương 4: Thành phố Roa và cô tiểu thư kiêu ngạo Eris",
					Illustrations: []string{"/mock-images/ln_16.png", "/mock-images/ln_17.png"},
					Content: `<p>Cỗ xe ngựa lọc cọc dừng lại trước cánh cổng dinh thự bề thế của Lãnh chúa Boreas Greyrat tại thành phố phồn hoa Roa. Nhiệm vụ mới của Rudeus là làm gia sư cho cô con gái duy nhất của gia tộc: Eris Boreas Greyrat.</p>
<p>— "Rầm!"</p>
<p>Cửa phòng bật mở dữ dội. Một thiếu nữ tóc đỏ rực lửa, đôi mắt màu ruby ánh lên tia nhìn hung hăng như một chú thú hoang lao thẳng về phía Rudeus:</p>
<p>— "Kẻ nào dám cả gan bảo dạy bảo ta?! Nhận lấy một đòn này!"</p>
<p>Cú đấm kèm theo luồng đấu khí cuồn cuộn xé gió lao tới. Rudeus giật mình lùi lại, nhanh như cắt tạo ra một bức tường đất chặn đứng cú đấm. Cả gian phòng rung chuyển, bụi bay mù mịt.</p>
<p>— "Hừm, không tệ đấy đồ lùn! Nhưng muốn làm thầy của Eris này, ngươi còn phải trải qua thử thách nữa!"</p>
<p>Cuộc gặp gỡ định mệnh giữa vị pháp sư trẻ tuổi và nàng kiếm sĩ hung bạo bắt đầu từ đây, mở ra chuỗi ngày đầy rẫy sóng gió và những hiểm nguy rình rập ngoài tầm kiểm soát của cả hai...</p>`,
				},
			},
		},
		{
			ID:          "ln-rezero",
			Title:       "Re:Zero − Bắt Đầu Lại Ở Thế Giới Khác",
			Description: "Natsuki Subaru, một học sinh bình thường đột ngột bị triệu hồi đến thế giới khác trên đường về từ cửa hàng tiện lợi. Năng lực duy nhất cậu nhận được là 'Trở Về Từ Cõi Chết' – quay ngược thời gian mỗi khi sinh mạng kết thúc. Liệu cậu có thể cứu vãn những người thân yêu khỏi lưỡi hái tử thần?",
			CoverURL:    "/mock-images/ln_18.jpg",
			Category:    "isekai",
			Tags:        "Isekai, Vòng Lặp, Tâm Lý, Bi Kịch",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -18),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Mở mắt giữa vương đô Lugnica",
					Illustrations: []string{"/mock-images/ln_19.jpg", "/mock-images/ln_18.jpg"},
					Content: `<p>Chớp mắt một cái, khung cảnh cửa hàng tiện lợi quen thuộc bỗng chốc biến mất. Thay vào đó là một khu chợ sầm uất với những cỗ xe rồng kéo đồ sộ, tiếng rao hàng huyên náo và những người dân mang tai thú kỳ lạ.</p>
<p>— "Hả?! Chuyện quái gì thế này? Triệu hồi thế giới khác sao?!"</p>
<p>Natsuki Subaru lục tung chiếc túi ni lông trên tay: một gói mì ăn liền, hai túi khoai tây chiên, một lon trà xanh và chiếc điện thoại nắp gập hoàn toàn mất sóng. Không có thánh kiếm, không có bảng kỹ năng bá đạo, thậm chí không có lấy một đồng xu tiền địa phương trong túi.</p>
<p>Ngay khi vừa rẽ vào một con hẻm vắng để định thần, ba tên du côn hung hãn đã chặn đầu khóa đuôi:</p>
<p>— "Này thằng nhóc lạ mặt, có tiền thì ngoan ngoãn nôn ra, không thì nếm thử lưỡi dao này!"</p>
<p>Dù cố gắng áp dụng các đòn cận chiến tự học, Subaru nhanh chóng bị hạ đo ván. Ngay trong khoảnh khắc cận kề lưỡi dao găm sắc lạnh, một giọng nói trong trẻo như chuông bạc vang lên xua tan bóng tối:</p>
<p>— "Dừng tay lại! Mau trả đồ cho người khác!"</p>
<p>Đứng ở đầu hẻm là một thiếu nữ tuyệt mỹ với mái tóc màu bạch kim lấp lánh dưới ánh chiều tà, đôi mắt màu thạch anh tím ngập tràn vẻ kiên định, và bên vai là một chú tinh linh mèo nhỏ nhắn đang lơ lửng.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Máu tươi nơi nhà kho tăm tối",
					Illustrations: []string{"/mock-images/ln_8.png"},
					Content: `<p>Đêm buông xuống trên khu ổ chuột vương đô. Subaru cùng thiếu nữ tự xưng là "Satella" lần theo dấu vết của kẻ trộm huy hiệu đến một nhà kho chứa đồ phế thải cũ kỹ.</p>
<p>Mùi rỉ sét nồng nặc và không khí đặc quánh sự bất an. Cánh cửa gỗ khẽ kêu kẽo kẹt khi Subaru đẩy nhẹ vào.</p>
<p>— "Này... có ai ở nhà không? Chúng tôi đến để thương lượng mua lại chiếc huy..."</p>
<p>Lời nói nghẹn lại nơi cổ họng. Dưới sàn nhà ẩm mốc là một vũng máu đen ngòm đang lan rộng. Chủ nhà kho - ông lão Rom khổng lồ - đã nằm bất động với vết rạch chí mạng nơi cổ họng.</p>
<p>— "Chạy mau, Satella!" – Subaru hét lên trong hoảng loạn.</p>
<p>Nhưng đã quá muộn. Một cái bóng đen uốn lượn như vũ nữ lướt qua trong màn đêm. Lưỡi dao cong lưỡi liềm lóe lên ánh thép lạnh ngắt. Cơn đau thấu xương ập đến trước khi Subaru kịp nhận thức chuyện gì vừa xảy ra. Phần bụng cậu bị xé rách, máu trào ra xối xả.</p>
<p>— "A... đau quá... lạnh quá..."</p>
<p>Bên cạnh cậu, thiếu nữ tóc bạc cũng ngã xuống, bàn tay cô run rẩy với về phía cậu trong vô vọng. Subaru dốc hết chút sinh lực tàn tạ, nắm lấy ngón tay lạnh ngắt của nàng:</p>
<p>— "Đừng sợ... tôi nhất định... sẽ cứu cô..."</p>
<p>Bóng tối nuốt chửng toàn bộ ý thức.</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Trở về từ cõi chết",
					Illustrations: []string{"/mock-images/ln_20.jpg"},
					Content: `<p>— "Này cậu em, có mua táo không thì bảo? Đứng ngây người ra chắn hết lối vào quầy hàng của tôi rồi đấy!"</p>
<p>Subaru bừng tỉnh, ngã nhào ra đất, thở dốc dữ dội. Mồ hôi lạnh đầm đìa ướt đẫm lưng áo. Cậu vội vàng vạch áo kiểm tra bụng: không có vết thương, không có máu, thậm chí tấm vải áo còn nguyên vẹn chưa từng bị xé rách.</p>
<p>Mặt trời lại một lần nữa đứng bóng trên đỉnh đầu. Những người dân tai thú vẫn tấp nập qua lại. Gói snack trên tay vẫn còn nguyên vẹn.</p>
<p>— "Mình... chưa chết? Hay là vừa rồi chỉ là một cơn ác mộng chân thực?"</p>
<p>Nhưng cơn đau xé ruột xé gan lúc nãy vẫn còn ám ảnh từng tế bào thần kinh. Cậu nhìn sang góc phố: ba tên du côn quen thuộc đang ngóng tìm con mồi. Thiếu nữ tóc bạc sắp sửa đi qua đây trong ít phút nữa.</p>
<p>— "Không phải mơ... Mình đã quay ngược thời gian! Cái chết không phải là kết thúc, mà là điểm khởi đầu của một vòng lặp!"</p>
<p>Bàn tay Subaru run lên bần bật, nhưng đôi mắt cậu rực lên một ngọn lửa kiên định chưa từng thấy. Lần này, cậu sẽ không để bi kịch lặp lại!</p>`,
				},
				{
					Number:        4,
					Title:         "Chương 4: Tia hy vọng từ dinh thự Roswaal",
					Illustrations: []string{"/mock-images/ln_1.png", "/mock-images/ln_21.png"},
					Content: `<p>Sau trận chiến sinh tử đánh bại kẻ săn ruột Elsa Granhirte nhờ sự trợ giúp của Kiếm Thánh Reinhard, Subaru mở mắt ra trên chiếc giường phủ nệm êm ái tại dinh thự quý tộc Roswaal.</p>
<p>Bên cạnh giường là hai cô hầu gái sinh đôi với mái tóc ngắn màu xanh lam và hồng đào, mặc trang phục maid hoàn hảo không tì vết.</p>
<p>— "Chị gái, vị khách khả nghi cuối cùng cũng tỉnh lại rồi kìa."</p>
<p>— "Rem, có vẻ như tên ngốc này đang nhìn chằm chằm vào chúng ta bằng ánh mắt đồi bại đấy."</p>
<p>Subaru bật cười bất lực trước màn đối đáp chua cay nhưng vô cùng đáng yêu của Rem và Ram. Cậu đã vượt qua cánh cửa tử thần đầu tiên, nhưng cậu không hề hay biết rằng, những chiếc bẫy chết người tinh vi hơn, tàn khốc hơn đang chờ đợi cậu ngay bên trong những hành lang lộng lẫy của dinh thự này...</p>`,
				},
			},
		},
		{
			ID:          "ln-cote",
			Title:       "Classroom of the Elite: Lớp Học Thực Lực",
			Description: "Trường Cao trung Nâng cao Tokyo là nơi nuôi dưỡng những nhân tài hàng đầu Nhật Bản với tỷ lệ 100% đỗ đại học và có việc làm. Nhưng tại lớp 1-D, nơi tập hợp những học sinh 'khiếm khuyết' của hệ thống giáo dục, Ayanokouji Kiyotaka bắt đầu che giấu năng lực thực sự của mình.",
			CoverURL:    "/mock-images/ln_10.png",
			Category:    "romcom",
			Tags:        "Học Đường, Tâm Lý, Đấu Trí, Harem",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -21),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Buổi khai giảng và 100.000 điểm tích lũy",
					Illustrations: []string{"/mock-images/ln_22.jpg", "/mock-images/ln_10.png"},
					Content: `<p>Mỗi tháng một học sinh được cấp 100.000 điểm cá nhân, tương đương 100.000 yên tiền mặt để tiêu xài tùy thích trong khuôn viên trường rộng lớn như một thành phố thu nhỏ. Không kỷ luật gắt gao, không bài tập về nhà ép buộc, tự do chơi game và ngủ gật trong giờ học.</p>
<p>Đa phần các học sinh lớp 1-D coi đây là thiên đường trần gian. Họ điên cuồng vung tiền vào những bữa ăn sang trọng, quần áo hàng hiệu và máy chơi game đắt tiền.</p>
<p>Tuy nhiên, ở góc cuối lớp bên cạnh cửa sổ, Ayanokouji Kiyotaka chỉ lặng lẽ quan sát chiếc điện thoại hiển thị số dư tài khoản của mình. Bên cạnh cậu, Horikita Suzune – cô gái lạnh lùng có mái tóc đen dài mượt mà – đang lật từng trang sách dày cộp mà không thèm để mắt đến xung quanh.</p>
<p>— "Cậu không đi mua sắm cùng mọi người à, Horikita?"</p>
<p>— "Tôi không có sở thích kết giao với những kẻ ngây thơ lãng phí thời gian. Và cậu cũng thế, Ayanokouji-kun. Đừng bắt chuyện với tôi như thể chúng ta thân thiết lắm."</p>
<p>Một câu trả lời sắc lẹm chuẩn phong cách Horikita. Nhưng cả hai đều lờ mờ nhận ra: bữa tiệc nào rồi cũng phải tàn, và cái giá phải trả cho sự tự do thái quá này sẽ không hề rẻ.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Cú sốc ngày đầu tháng Năm",
					Illustrations: []string{"/mock-images/ln_23.png"},
					Content: `<p>Ngày 1 tháng Năm. Tiếng chuông báo tin nhắn chuyển điểm vang lên đồng loạt trong lớp học. Mọi ánh mắt háo hức nhìn vào màn hình điện thoại.</p>
<p>— "Ơ... sao số dư tháng này lại là 0 điểm?!"</p>
<p>— "Lỗi hệ thống à cô Chabashira? Sao chúng em không nhận được đồng nào thế này?!"</p>
<p>Cô giáo chủ nhiệm Chabashira Sae bước lên bục giảng, điếu thuốc trên môi khẽ nhả làn khói xám uể oải. Ánh mắt cô nhìn cả lớp như nhìn một bầy sâu bọ đáng thương:</p>
<p>— "Một lũ ngốc nghếch. Ngôi trường này hoạt động dựa trên cơ chế 'Thực Lực Tối Thượng'. Điểm số hàng tháng của các người phản ánh giá trị và thành tích của cả tập thể lớp."</p>
<p>Cô bật máy chiếu lên bảng. Một loạt con số đỏ rực hiện ra:</p>
<p>— "Trong tháng qua: 98 lần đi muộn, 342 lần nói chuyện riêng trong giờ, 0 giờ học nhóm ngoại khóa. Điểm đánh giá lớp 1-D tụt xuống đúng bằng 0. Từ hôm nay, các người sẽ phải sống sót với 0 yên cho đến kỳ thi giữa kỳ!"</p>
<p>Tiếng la hét và tuyệt vọng bao trùm căn phòng. Trò chơi sinh tồn trí tuệ chính thức bắt đầu.</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Kỳ thi trắc nghiệm và ván cờ sau cánh gà",
					Illustrations: []string{"/mock-images/ln_24.png", "/mock-images/ln_25.jpg"},
					Content: `<p>Nếu có bất kỳ ai bị điểm liệt trong kỳ thi giữa kỳ sắp tới, người đó sẽ lập tức bị đuổi học không khoan nhượng. Ba tên học sinh cá biệt đội sổ của lớp gồm Sudo, Yamauchi và Ike đang đứng trước bờ vực diệt vong.</p>
<p>Horikita cố gắng mở lớp phụ đạo nhưng bị từ chối thẳng thừng bởi thái độ kiêu ngạo của cô.</p>
<p>Đêm muộn tại cầu thang thoát hiểm tầng ba, Ayanokouji gặp một đàn anh lớp khóa trên:</p>
<p>— "Anh có đề thi cũ của kỳ thi giữa kỳ ba năm gần nhất đúng không? Tôi muốn mua lại nó bằng 15.000 điểm cá nhân."</p>
<p>Đàn anh cười khẩy:</p>
<p>— "Nhóc con lớp D biết điều đấy. Ở trường này, đề thi cũ không hề bị cấm lưu hành, nhưng không ai nói cho lũ tân binh các cậu biết cả."</p>
<p>Ayanokouji nhận lấy xấp đề thi, nhưng cậu không mang nó ra với tư cách của mình. Cậu trao nó cho Horikita và dàn dựng một kịch bản hoàn hảo để cô trở thành vị cứu tinh của cả lớp trong mắt mọi người.</p>
<blockquote>"Đối với tôi, tất cả mọi người chỉ là những công cụ. Miễn là cuối cùng tôi giành chiến thắng, phương thức nào cũng không thành vấn đề."</blockquote>`,
				},
			},
		},
		{
			ID:          "ln-86",
			Title:       "86 - Eighty Six: Chiến Trường Khắc Nghiệt",
			Description: "Cộng hòa San Magnolia tuyên bố đã phát triển thành công drone chiến đấu không người lái 'Juggernaut' để chống lại quân xâm lược Legion. Thế nhưng, đằng sau lớp vỏ bọc không thương vong ấy là Khu Vực 86 - nơi những con người bị tước đoạt nhân quyền phải chiến đấu và ngã xuống từng ngày.",
			CoverURL:    "/mock-images/ln_27.png",
			Category:    "scifi",
			Tags:        "Mecha, Viễn Tưởng, Chiến Tranh, Tình Cảm",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -24),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Bản báo cáo không thương vong",
					Illustrations: []string{"/mock-images/ln_28.jpg", "/mock-images/ln_27.png"},
					Content: `<p>— "Báo cáo chiến sự ngày hôm nay: Mặt trận phía Đông, tiểu đội Spearhead đã đẩy lùi hoàn toàn đợt tiến công của quân đoàn Legion. Thiệt hại quân ta: Không có thương vong về người."</p>
<p>Tiếng phát thanh viên dịu dàng vang lên trên khắp các đại lộ của thủ đô Charité. Người dân da trắng thuần chủng Alba nâng ly sâm panh chúc mừng một nền hòa bình giả tạo được che đậy bằng máu và nước mắt.</p>
<p>Tại phòng chỉ huy quân sự trung tâm, Thiếu úy Vladilena Milizé – một cô gái trẻ với mái tóc bạch kim ánh bạc – siết chặt bản báo cáo chiến trường trong tay. Cô biết rất rõ sự thật tàn nhẫn đằng sau cụm từ "không có thương vong".</p>
<p>Những cỗ máy Juggernaut bốn chân bằng hợp kim nhôm giòn tan không hề tự động lái. Bên trong khoang lái chật hẹp, ngột ngạt nhiệt độ lên tới bốn mươi độ C là những thiếu niên, thiếu nữ thuộc chủng tộc màu da bị tước đoạt quốc tịch và quyền con người – những kẻ mang con số '86'.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Tiếng gọi qua thiết bị Para-RAID",
					Illustrations: []string{"/mock-images/ln_5.png"},
					Content: `<p>Chiếc khuyên tai Para-RAID khẽ rung lên, kích hoạt xung thần kinh kết nối trực tiếp vào thùy cảm giác của Lena. Một tiếng rè rè của sóng âm vang lên, rồi một giọng nói trầm ổn, lạnh lùng như băng tuyết phương Bắc truyền đến:</p>
<p>— "Handler One, đây là Undertaker. Tiểu đội Spearhead nghe lệnh."</p>
<p>— "Chào buổi tối, Undertaker. Tôi là Thiếu úy Vladilena Milizé, người vừa được bổ nhiệm làm chỉ huy điều hành mới của các bạn."</p>
<p>Đầu dây bên kia im lặng một thoáng. Một tiếng cười khẽ mỉa mai của một thành viên khác vọng lại:</p>
<p>— "Lại một con heo trắng nữa thích chơi trò cứu rỗi đạo đức giả rồi..."</p>
<p>— "Im đi, Kurena." – Undertaker lạnh lùng ngắt lời cấp dưới, rồi quay lại với Lena – "Thưa Thiếu úy, nếu ngài muốn chỉ huy chúng tôi qua màn hình vô tuyền an toàn cách xa hàng trăm cây số, chúng tôi không có ý kiến. Nhưng xin ngài hãy chuẩn bị tinh thần: chiến trường này không có chỗ cho sự thương hại viển vông."</p>
<p>Đó là lần đầu tiên Lena lắng nghe giọng nói của Shin – chàng trai mang biệt danh "Tử Thần", người mang trên mình sứ mệnh tiễn đưa những đồng đội đã ngã xuống về nơi an nghỉ cuối cùng.</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Hoa bỉ ngạn nở rộ bên chiến hào",
					Illustrations: []string{"/mock-images/ln_29.png", "/mock-images/ln_30.jpg"},
					Content: `<p>Mặt trời lặn nhuộm đỏ rực cánh đồng hoa bỉ ngạn trải dài vô tận ngoài căn cứ Khu Vực 86. Shin ngồi trên nắp khoang lái của cỗ Juggernaut mang biểu tượng bộ xương không đầu cầm lưỡi hái, tay cầm con dao găm tỉ mỉ khắc tên người đồng đội vừa hy sinh sáng nay lên một mảnh nhôm nhỏ.</p>
<p>Năm trăm bảy mươi sáu mảnh nhôm. Đó là số lượng linh hồn mà cậu đang gánh vác trên lưng.</p>
<p>— "Undertaker... cậu có nghe thấy không?" – Giọng nói của Lena qua Para-RAID khẽ run rẩy.</p>
<p>— "Tôi đang nghe."</p>
<p>— "Tôi... tôi xin lỗi vì không thể làm gì để bảo vệ các bạn. Những chiếc máy bay không người lái của Cộng hòa thực chất chỉ là những chiếc quan tài di động..."</p>
<p>Shin ngẩng đầu nhìn về phía vầng trăng khuyết treo lơ lửng trên nền trời tím sẫm:</p>
<p>— "Đừng khóc vì chúng tôi, Thiếu úy. Chúng tôi chiến đấu không phải vì Cộng hòa hay vì danh dự của những kẻ da trắng. Chúng tôi chiến đấu vì lòng tự tôn của chính mình, để khi ngã xuống vẫn ngẩng cao đầu như một con người."</p>`,
				},
			},
		},
		{
			ID:          "ln-oregairu",
			Title:       "Oregairu: Chuyện Tình Thanh Xuân Bi Hài",
			Description: "Hikigaya Hachiman là một nam sinh trung học có cái nhìn méo mó và cay đắng về thanh xuân, tình bạn và tình yêu. Bị giáo viên chủ nhiệm ép tham gia CLB Tình Nguyện, cậu phải đối mặt với Yukinoshita Yukino - một đóa hoa trên đỉnh núi cao với sự sắc bén đến đau lòng.",
			CoverURL:    "/mock-images/ln_32.png",
			Category:    "romcom",
			Tags:        "Rom-Com, Học Đường, Đời Thường, Triết Lý",
			Status:      "published",
			CreatedAt:   time.Now().AddDate(0, 0, -27),
			Chapters: []ChapterSeed{
				{
					Number:        1,
					Title:         "Chương 1: Bài luận thanh xuân là một trò lừa dối",
					Illustrations: []string{"/mock-images/ln_33.jpg", "/mock-images/ln_32.png"},
					Content: `<p><em>"Thanh xuân là một điều dối trá, là một sự tha hóa độc ác. Những kẻ luôn miệng ca ngợi thanh xuân chỉ đang cố lừa dối bản thân và những người xung quanh..."</em></p>
<p>Đó là những dòng mở đầu trong bài luận tốt nghiệp môn Ngữ văn của Hikigaya Hachiman. Và phần thưởng cho sự chân thật quá mức đó là một cú đấm vào bụng từ cô giáo chủ nhiệm Hiratsuka Shizuka.</p>
<p>— "Hikigaya, đôi mắt cá chết và cái tâm hồn mục ruỗng này của em cần phải được cải tạo ngay lập tức!"</p>
<p>Bị lôi xềnh xệch dọc hành lang trường Sobu, Hachiman bị tống vào căn phòng câu lạc bộ nằm ở góc khuất nhất của tầng bốn. Ánh hoàng hôn rọi qua khung cửa sổ, chiếu lên bóng hình một thiếu nữ đang ngồi đọc sách đơn độc.</p>
<p>Mái tóc đen dài óng ả buông xõa qua làn da trắng sứ, đôi mắt toát lên vẻ thông tuệ sắc sảo khiến bất cứ ai nhìn vào cũng phải chùn bước.</p>
<p>— "Cô Hiratsuka, em đã nói nhiều lần rồi, đừng tự tiện mang những thứ rác rưởi không rõ nguồn gốc vào phòng câu lạc bộ của em."</p>
<p>Đó là lời chào đầu tiên của Tuyệt sắc giai nhân trường Sobu – Yukinoshita Yukino – dành cho Hachiman.</p>`,
				},
				{
					Number:        2,
					Title:         "Chương 2: Yuigahama Yui và mẻ bánh quy cháy khét",
					Illustrations: []string{"/mock-images/ln_2.png"},
					Content: `<p>Thành viên thứ ba xuất hiện như một tia nắng xua tan bầu không khí ngột ngạt giữa hai con người lập dị. Yuigahama Yui – cô nàng tóc búi củ tỏi ngây thơ, luôn nở nụ cười rạng rỡ và lo sợ bị người khác cô lập.</p>
<p>— "Tớ... tớ muốn học cách làm bánh quy để tặng cho một người!"</p>
<p>Sau ba tiếng đồng hồ vật lộn trong phòng bếp thực hành, mẻ bánh quy ra lò trông không khác gì những cục than tổ ong đen sì bốc khói cay xè.</p>
<p>Yukinoshita nếm thử một góc nhỏ rồi bình luận không thương tiếc:</p>
<p>— "Độ xốp bằng không, vị đắng chiếm 90%, trông như tàn tích sau một vụ nổ hóa học."</p>
<p>Yui cúi gầm mặt, chực khóc. Nhưng Hachiman chỉ nhón lấy một miếng cắn vụn rồi nhai chậm rãi:</p>
<p>— "Thực ra... bánh quy làm thủ công mà hoàn hảo như đồ mua ngoài tiệm thì chẳng có ý nghĩa gì. Chính cái vẻ vụng về, cháy xém này mới chứng minh người làm đã dồn tâm huyết vào nó. Con trai tụi này... đôi khi chỉ cần nhận được một thứ như vậy là đủ vui cả ngày rồi."</p>
<p>Cả căn phòng bỗng chốc lặng ngắt. Yukino ngước lên nhìn Hachiman bằng một ánh mắt hoàn toàn khác, còn mặt Yui thì đỏ bừng như quả cà chua chín.</p>`,
				},
				{
					Number:        3,
					Title:         "Chương 3: Cuộc đối thoại tìm kiếm điều chân thực",
					Illustrations: []string{"/mock-images/ln_34.png", "/mock-images/ln_35.jpg"},
					Content: `<p>Mùa đông lạnh buốt bao trùm khuôn viên trường. Sau những biến cố đau lòng tại kỳ lễ hội văn hóa và chuyến dã ngoại Kyoto, bức tường ngăn cách giữa ba người họ ngày càng trở nên mong manh và đau đớn.</p>
<p>Trong căn phòng câu lạc bộ không bật lò sưởi, Hachiman siết chặt nắm tay, đối diện với Yukino và Yui:</p>
<p>— "Tôi không muốn sự thấu hiểu giả tạo bằng những lời xã giao sáo rỗng. Tôi không cần một mối quan hệ hòa thuận nếu cái giá phải trả là sự im lặng trước những vết thương..."</p>
<p>Giọng cậu nghẹn lại, đôi mắt cá chết ngấn nước – điều chưa từng xảy ra trong suốt cuộc đời cậu:</p>
<p>— "Tôi... tôi muốn một thứ gì đó chân thực (genuine)!"</p>
<p>Yukino sững sờ. Đôi bàn tay nhỏ bé của nàng khẽ run lên. Lần đầu tiên, nàng không tìm ra được bất kỳ câu từ sắc bén nào để phản bác lại chàng trai trước mặt.</p>`,
				},
			},
		},
	}
}
