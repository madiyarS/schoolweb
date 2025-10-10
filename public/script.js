document.addEventListener('DOMContentLoaded', function() {

    // Translations
    const translations = {
        ru: {
            "title": "Начальная школа Академия - Астана",
            "lang.ru": "Русский",
            "lang.kz": "Қазақша",
            "lang.en": "English",
            "nav.about": "О школе",
            "nav.programs": "Программы",
            "nav.achievements": "Достижения",
            "nav.teachers": "Педагоги",
            "nav.news": "Новости",
            "nav.contact": "Контакты",
            "nav.apply": "Поступить",
            "hero.title": "Начальная школа Академия",
            "hero.subtitle": "Благоприятная среда для академических и личных успехов каждого ребёнка",
            "hero.subtitle.en": "A nurturing environment for academic and personal success of every child",
            "hero.btnPrograms": "Наши программы",
            "hero.btnAbout": "Узнать больше",
            "about.title": "Наша миссия",
            "about.mission": "Мы создадим благоприятную и стимулирующую учебную среду для достижения академических и личных успехов каждого ребёнка.",
            "about.values.title": "Наши ценности:",
            "about.values.v1": "Обучение, ориентированное на каждого ученика",
            "about.values.v2": "Трудолюбие и дух сотрудничества",
            "about.values.v3": "Самоуважение и уважение к другим",
            "about.values.v4": "Мы учим быть ответственным — за свой выбор, за свой успех",
            "about.values.v5": "Поддержание семейных ценностей",
            "about.stat1": "Года работы",
            "about.stat2": "Учеников",
            "about.stat3": "Педагогов",
            "about.stat4": "Учеников в классе",
            "programs.title": "Образовательные программы",
            "programs.subtitle": "Начальная школа для детей 1-4 классов",
            "programs.card1.title": "Математика на английском",
            "programs.card1.text": "Уникальная программа изучения математики на английском языке для развития билингвального мышления.",
            "programs.card2.title": "Развитие навыков обучения",
            "programs.card2.text": "Learning how to learn — обучение навыкам самообразования и развитие эмоционального интеллекта.",
            "programs.card3.title": "Основная программа",
            "programs.card3.text": "Все основные предметы начальной школы с индивидуальным подходом к каждому ученику.",
            "programs.activities.title": "Дополнительные активности",
            "programs.activities.chess": "Шахматы",
            "programs.activities.robotics": "Робототехника",
            "programs.activities.vocal": "Вокал",
            "programs.activities.dance": "Танцы",
            "programs.activities.mental": "Ментальная арифметика",
            "programs.activities.acting": "Актерское мастерство",
            "achievements.title": "Достижения наших учеников",
            "achievements.subtitle": "15 призовых мест на олимпиадах за последние 3 года",
            "achievements.math.title": "Математические олимпиады",
            "achievements.math.text": "Скорняков Всеволод — 2 место на городской олимпиаде по математике среди школьников",
            "achievements.abay.title": "«Абаевские чтения»",
            "achievements.abay.text": "Республиканский конкурс: Марат Альтаир — 1 место, Кононов Иван — 2 место",
            "achievements.akbota.title": "«Акбота»",
            "achievements.akbota.text": "Республиканский уровень: Кононов Иван — 1 место, Тургумбаева Севиль — 3 место, Кафеджис Владислав — 1 место",
            "achievements.kangaroo.title": "«Кенгуру-математика»",
            "achievements.kangaroo.text": "Республиканский уровень: Проскрякова Лена — 1 место и многие другие призёры",
            "teachers.title": "Наши педагоги",
            "teachers.subtitle": "Опытные преподаватели с индивидуальным подходом к каждому ребёнку",
            "teachers.t1.name": "Сидоренко Светлана Игоревна",
            "teachers.t1.position": "Учитель начальных классов",
            "teachers.t1.experience": "Опыт работы: 15 лет",
            "teachers.t1.achievements": "Высшая категория. Сертификаты за подготовку к олимпиадам, участие в курсах и городских конференциях",
            "teachers.t2.name": "Михненко Марина Филипповна",
            "teachers.t2.position": "Учитель начальных классов, учитель английского языка",
            "teachers.t2.experience": "Опыт работы: 25 лет",
            "teachers.t2.achievements": "Сертификаты за подготовку к олимпиадам, участие в курсах и городских конференциях",
            "teachers.t3.name": "Прядко Стелла Валерьевна",
            "teachers.t3.position": "Директор школы",
            "teachers.t3.experience": "Опыт работы в образовании",
            "teachers.t3.achievements": "Руководитель образовательного учреждения",
            "teachers.t4.name": "Залесская Дарья Олеговна",
            "teachers.t4.position": "Учитель начальных классов",
            "teachers.t4.experience": "Опыт работы: 12 лет",
            "teachers.t4.achievements": "Модератор. Сертификаты за подготовку к олимпиадам, участие в курсах и городских конференциях",
            "news.title": "Последние новости",
            "news.noNews": "Новостей пока нет.",
            "news.error": "Не удалось загрузить новости.",
            "news.published": "Опубликовано:",
            "contact.title": "Контакты и поступление",
            "contact.info.title": "Свяжитесь с нами",
            "contact.info.text": "Наша приемная комиссия готова ответить на все ваши вопросы и помочь с процессом поступления.",
            "contact.info.address": "г. Астана, район Сарыарка, ул. Шыганак, 7",
            "contact.info.hours": "Пн-Пт: 08:00 - 17:00",
            "contact.social.title": "Мы в социальных сетях",
            "contact.license": "Лицензия № KZ88LAA00031985 от 07.09.2021г",
            "contact.form.title": "Задать вопрос",
            "contact.form.name": "Ваше имя",
            "contact.form.email": "Ваш Email",
            "contact.form.phone": "Номер телефона",
            "contact.form.message": "Ваше сообщение",
            "contact.form.submit": "Отправить",
            "contact.form.errorName": "Пожалуйста, введите ваше имя.",
            "contact.form.errorEmail": "Пожалуйста, введите корректный email.",
            "contact.form.errorPhone": "Пожалуйста, введите номер телефона.",
            "contact.form.errorMessage": "Пожалуйста, введите ваше сообщение.",
            "contact.form.sending": "Отправка...",
            "contact.form.success": "Сообщение успешно получено и сохранено!",
            "contact.form.error": "Не удалось отправить сообщение. Попробуйте позже.",
            "footer.description": "Благоприятная и стимулирующая учебная среда для достижения успехов каждого ребёнка.",
            "footer.founded": "Основана в 2021 году",
            "footer.links.title": "Быстрые ссылки",
            "footer.social.title": "Мы в соцсетях",
            "footer.copy": "© 2024 ТОО «Начальная школа Академия». Все права защищены."
        },
        kz: {
            "title": "Академия бастауыш мектебі - Астана",
            "lang.ru": "Орысша",
            "lang.kz": "Қазақша",
            "lang.en": "English",
            "nav.about": "Мектеп туралы",
            "nav.programs": "Бағдарламалар",
            "nav.achievements": "Жетістіктер",
            "nav.teachers": "Мұғалімдер",
            "nav.news": "Жаңалықтар",
            "nav.contact": "Байланыс",
            "nav.apply": "Түсу",
            "hero.title": "Академия бастауыш мектебі",
            "hero.subtitle": "Әр баланың академиялық және жеке табыстарына жету үшін қолайлы орта",
            "hero.subtitle.en": "A nurturing environment for academic and personal success of every child",
            "hero.btnPrograms": "Біздің бағдарламалар",
            "hero.btnAbout": "Толығырақ",
            "about.title": "Біздің миссиямыз",
            "about.mission": "Біз әр баланың академиялық және жеке табыстарына жету үшін қолайлы және ынталандырушы оқу ортасын құрамыз.",
            "about.values.title": "Біздің құндылықтарымыз:",
            "about.values.v1": "Әр оқушыға бағытталған оқыту",
            "about.values.v2": "Еңбекқорлық және ынтымақтастық рухы",
            "about.values.v3": "Өзін-өзі құрметтеу және басқаларды құрметтеу",
            "about.values.v4": "Біз жауапкершілікті үйретеміз — өз таңдауы үшін, өз табысы үшін",
            "about.values.v5": "Отбасылық құндылықтарды сақтау",
            "about.stat1": "Жұмыс жылдары",
            "about.stat2": "Оқушылар",
            "about.stat3": "Мұғалімдер",
            "about.stat4": "Сыныптағы оқушылар",
            "programs.title": "Білім беру бағдарламалары",
            "programs.subtitle": "1-4 сынып балаларына арналған бастауыш мектеп",
            "programs.card1.title": "Ағылшын тіліндегі математика",
            "programs.card1.text": "Қостілді ойлауды дамыту үшін математиканы ағылшын тілінде оқытудың бірегей бағдарламасы.",
            "programs.card2.title": "Оқу дағдыларын дамыту",
            "programs.card2.text": "Learning how to learn — өздігінен білім алу дағдыларын үйрету және эмоционалды интеллектті дамыту.",
            "programs.card3.title": "Негізгі бағдарлама",
            "programs.card3.text": "Әр оқушыға жеке көзқарас қолданылатын бастауыш мектептің барлық негізгі пәндері.",
            "programs.activities.title": "Қосымша белсенділіктер",
            "programs.activities.chess": "Шахмат",
            "programs.activities.robotics": "Робототехника",
            "programs.activities.vocal": "Вокал",
            "programs.activities.dance": "Би",
            "programs.activities.mental": "Ментальды арифметика",
            "programs.activities.acting": "Актерлік шеберлік",
            "achievements.title": "Оқушыларымыздың жетістіктері",
            "achievements.subtitle": "Соңғы 3 жылда олимпиадаларда 15 жүлделі орын",
            "achievements.math.title": "Математикалық олимпиадалар",
            "achievements.math.text": "Скорняков Всеволод — қалалық математика олимпиадасында 2-орын",
            "achievements.abay.title": "«Абай оқулары»",
            "achievements.abay.text": "Республикалық конкурс: Марат Альтаир — 1-орын, Кононов Иван — 2-орын",
            "achievements.akbota.title": "«Ақбота»",
            "achievements.akbota.text": "Республикалық деңгей: Кононов Иван — 1-орын, Тургумбаева Севиль — 3-орын, Кафеджис Владислав — 1-орын",
            "achievements.kangaroo.title": "«Кенгуру-математика»",
            "achievements.kangaroo.text": "Республикалық деңгей: Проскрякова Лена — 1-орын және басқа да көптеген жүлдегерлер",
            "teachers.title": "Біздің мұғалімдер",
            "teachers.subtitle": "Әр балаға жеке көзқарас қолданатын тәжірибелі мұғалімдер",
            "teachers.t1.name": "Сидоренко Светлана Игоревна",
            "teachers.t1.position": "Бастауыш сынып мұғалімі",
            "teachers.t1.experience": "Жұмыс тәжірибесі: 15 жыл",
            "teachers.t1.achievements": "Жоғары санат. Олимпиадаларға дайындау сертификаттары, курстар мен қалалық конференцияларға қатысу",
            "teachers.t2.name": "Михненко Марина Филипповна",
            "teachers.t2.position": "Бастауыш сынып мұғалімі, ағылшын тілі мұғалімі",
            "teachers.t2.experience": "Жұмыс тәжірибесі: 25 жыл",
            "teachers.t2.achievements": "Олимпиадаларға дайындау сертификаттары, курстар мен қалалық конференцияларға қатысу",
            "teachers.t3.name": "Прядко Стелла Валерьевна",
            "teachers.t3.position": "Мектеп директоры",
            "teachers.t3.experience": "Білім беру саласындағы жұмыс тәжірибесі",
            "teachers.t3.achievements": "Білім беру мекемесінің басшысы",
            "teachers.t4.name": "Залесская Дарья Олеговна",
            "teachers.t4.position": "Бастауыш сынып мұғалімі",
            "teachers.t4.experience": "Жұмыс тәжірибесі: 12 жыл",
            "teachers.t4.achievements": "Модератор. Олимпиадаларға дайындау сертификаттары, курстар мен қалалық конференцияларға қатысу",
            "news.title": "Соңғы жаңалықтар",
            "news.noNews": "Жаңалықтар әзірше жоқ.",
            "news.error": "Жаңалықтарды жүктеу мүмкін болмады.",
            "news.published": "Жарияланды:",
            "contact.title": "Байланыстар және қабылдау",
            "contact.info.title": "Бізбен байланысыңыз",
            "contact.info.text": "Біздің қабылдау комиссиясы барлық сұрақтарыңызға жауап беруге және түсу процесіне көмектесуге дайын.",
            "contact.info.address": "Астана қ., Сарыарқа ауданы, Шығанақ к-сі, 7",
            "contact.info.hours": "Дс-Жм: 08:00 - 17:00",
            "contact.social.title": "Біз әлеуметтік желілерде",
            "contact.license": "Лицензия № KZ88LAA00031985 07.09.2021ж",
            "contact.form.title": "Сұрақ қою",
            "contact.form.name": "Сіздің атыңыз",
            "contact.form.email": "Сіздің Email",
            "contact.form.phone": "Телефон нөмірі",
            "contact.form.message": "Сіздің хабарламаңыз",
            "contact.form.submit": "Жіберу",
            "contact.form.errorName": "Атыңызды енгізіңіз.",
            "contact.form.errorEmail": "Дұрыс email енгізіңіз.",
            "contact.form.errorPhone": "Телефон нөмірін енгізіңіз.",
            "contact.form.errorMessage": "Хабарламаңызды енгізіңіз.",
            "contact.form.sending": "Жіберілуде...",
            "contact.form.success": "Хабарлама сәтті қабылданды және сақталды!",
            "contact.form.error": "Хабарламаны жіберу мүмкін болмады. Кейінірек қайталап көріңіз.",
            "footer.description": "Әр баланың табыстарына жету үшін қолайлы және ынталандырушы оқу ортасы.",
            "footer.founded": "2021 жылы құрылған",
            "footer.links.title": "Жылдам сілтемелер",
            "footer.social.title": "Біз әлеуметтік желілерде",
            "footer.copy": "© 2024 «Академия бастауыш мектебі» ЖШС. Барлық құқықтар қорғалған."
        },
        en: {
            "title": "Akademia Primary School - Astana",
            "lang.ru": "Russian",
            "lang.kz": "Kazakh",
            "lang.en": "English",
            "nav.about": "About School",
            "nav.programs": "Programs",
            "nav.achievements": "Achievements",
            "nav.teachers": "Teachers",
            "nav.news": "News",
            "nav.contact": "Contacts",
            "nav.apply": "Apply",
            "news.title": "Latest News",
            "news.noNews": "No news yet.",
            "news.error": "Failed to load news.",
            "news.published": "Published:",
            "hero.title": "Akademia Primary School",
            "hero.subtitle": "A nurturing environment for academic and personal success of every child",
            "hero.subtitle.en": "A nurturing environment for academic and personal success of every child",
            "hero.btnPrograms": "Our Programs",
            "hero.btnAbout": "Learn More",
            "about.title": "Our Mission",
            "about.mission": "We create a favorable and stimulating learning environment for the academic and personal success of every child.",
            "about.values.title": "Our Values:",
            "about.values.v1": "Student-centered learning",
            "about.values.v2": "Hard work and spirit of cooperation",
            "about.values.v3": "Self-respect and respect for others",
            "about.values.v4": "We teach responsibility — for one's choices, for one's success",
            "about.values.v5": "Maintaining family values",
            "about.stat1": "Years of Experience",
            "about.stat2": "Students",
            "about.stat3": "Teachers",
            "about.stat4": "Students per Class",
            "programs.title": "Educational Programs",
            "programs.subtitle": "Primary school for children grades 1-4",
            "programs.card1.title": "Math in English",
            "programs.card1.text": "Unique program for learning mathematics in English to develop bilingual thinking.",
            "programs.card2.title": "Learning Skills Development",
            "programs.card2.text": "Learning how to learn — teaching self-education skills and developing emotional intelligence.",
            "programs.card3.title": "Core Curriculum",
            "programs.card3.text": "All primary school core subjects with an individual approach to each student.",
            "programs.activities.title": "Additional Activities",
            "programs.activities.chess": "Chess",
            "programs.activities.robotics": "Robotics",
            "programs.activities.vocal": "Vocal",
            "programs.activities.dance": "Dance",
            "programs.activities.mental": "Mental Arithmetic",
            "programs.activities.acting": "Acting",
            "achievements.title": "Student Achievements",
            "achievements.subtitle": "15 prize-winning places in olympiads over the past 3 years",
            "achievements.math.title": "Math Olympiads",
            "achievements.math.text": "Skornyakov Vsevolod — 2nd place at city math olympiad among schoolchildren",
            "achievements.abay.title": "\"Abay Readings\"",
            "achievements.abay.text": "Republican competition: Marat Altair — 1st place, Kononov Ivan — 2nd place",
            "achievements.akbota.title": "\"Akbota\"",
            "achievements.akbota.text": "Republican level: Kononov Ivan — 1st place, Turgumbayeva Sevil — 3rd place, Kafejis Vladislav — 1st place",
            "achievements.kangaroo.title": "\"Kangaroo Math\"",
            "achievements.kangaroo.text": "Republican level: Proskryakova Lena — 1st place and many other prize winners",
            "teachers.title": "Our Teachers",
            "teachers.subtitle": "Experienced teachers with individual approach to each child",
            "teachers.t1.name": "Sidorenko Svetlana Igorevna",
            "teachers.t1.position": "Primary School Teacher",
            "teachers.t1.experience": "Work experience: 15 years",
            "teachers.t1.achievements": "Highest category. Certificates for olympiad preparation, participation in courses and city conferences",
            "teachers.t2.name": "Mikhnenko Marina Filippovna",
            "teachers.t2.position": "Primary School Teacher, English Teacher",
            "teachers.t2.experience": "Work experience: 25 years",
            "teachers.t2.achievements": "Certificates for olympiad preparation, participation in courses and city conferences",
            "teachers.t3.name": "Pryadko Stella Valeryevna",
            "teachers.t3.position": "School Director",
            "teachers.t3.experience": "Educational work experience",
            "teachers.t3.achievements": "Head of educational institution",
            "teachers.t4.name": "Zalesskaya Darya Olegovna",
            "teachers.t4.position": "Primary School Teacher",
            "teachers.t4.experience": "Work experience: 12 years",
            "teachers.t4.achievements": "Moderator. Certificates for olympiad preparation, participation in courses and city conferences",
            "contact.title": "Contacts and Admission",
            "contact.info.title": "Contact Us",
            "contact.info.text": "Our admissions office is ready to answer all your questions and help with the enrollment process.",
            "contact.info.address": "Astana, Saryarka district, Shyganak st., 7",
            "contact.info.hours": "Mon-Fri: 08:00 - 17:00",
            "contact.social.title": "Follow Us",
            "contact.license": "License № KZ88LAA00031985 dated 07.09.2021",
            "contact.form.title": "Ask a Question",
            "contact.form.name": "Your Name",
            "contact.form.email": "Your Email",
            "contact.form.phone": "Phone Number",
            "contact.form.message": "Your Message",
            "contact.form.submit": "Send",
            "contact.form.errorName": "Please enter your name.",
            "contact.form.errorEmail": "Please enter a valid email.",
            "contact.form.errorPhone": "Please enter phone number.",
            "contact.form.errorMessage": "Please enter your message.",
            "contact.form.sending": "Sending...",
            "contact.form.success": "Message successfully received and saved!",
            "contact.form.error": "Failed to send message. Please try again later.",
            "footer.description": "A favorable and stimulating learning environment for the success of every child.",
            "footer.founded": "Founded in 2021",
            "footer.links.title": "Quick Links",
            "footer.social.title": "Follow Us",
            "footer.copy": "© 2024 Akademia Primary School LLP. All rights reserved."
        }
    };

    // News Carousel Implementation
    const newsCarousel = {
        currentSlide: 0,
        slides: [],
        newsPerSlide: 3,
        autoplayInterval: null,
        
        init() {
            this.loadNews();
        },
        
        async loadNews() {
            const container = document.getElementById('news-grid-container');
            
            try {
                // Try to load from API
                const response = await fetch('/api/news');
                if (response.ok) {
                    const articles = await response.json();
                    if (Array.isArray(articles) && articles.length > 0) {
                        this.createCarousel(articles, container);
                        return;
                    }
                }
                
                // If API is unavailable, use static news
                this.loadStaticNews(container);
                
            } catch (error) {
                console.error('Error loading news:', error);
                // Load static news as fallback
                this.loadStaticNews(container);
            }
        },

        loadStaticNews(container) {
            // Static news for demonstration
            const staticArticles = [
                {
                    id: 1,
                    title: "Открытие нового учебного года",
                    content: "Начальная школа Академия с радостью открывает двери для новых и вернувшихся учеников. Мы готовы к новому учебному году с обновленными программами и свежими идеями.",
                    created_at: new Date().toISOString(),
                    image_url: "https://images.unsplash.com/photo-1497633762265-9d179a990aa6?q=80&w=2073&auto=format&fit=crop"
                },
                {
                    id: 2,
                    title: "Победа в городской олимпиаде",
                    content: "Наш ученик Скорняков Всеволод занял 2 место на городской олимпиаде по математике. Поздравляем с выдающимся достижением!",
                    created_at: new Date(Date.now() - 86400000).toISOString(),
                    image_url: "https://images.unsplash.com/photo-1567427017947-545c5f8d16ad?q=80&w=2053&auto=format&fit=crop"
                },
                {
                    id: 3,
                    title: "День открытых дверей",
                    content: "Приглашаем родителей и будущих учеников познакомиться с нашей школой, встретиться с педагогами и узнать о наших образовательных программах.",
                    created_at: new Date(Date.now() - 172800000).toISOString(),
                    image_url: "https://images.unsplash.com/photo-1503676260728-1c00da094a0b?q=80&w=2022&auto=format&fit=crop"
                },
                {
                    id: 4,
                    title: "Занятия по робототехнике",
                    content: "В школе стартовали занятия по робототехнике для учеников всех классов. Дети с энтузиазмом осваивают новые технологии.",
                    created_at: new Date(Date.now() - 259200000).toISOString(),
                    image_url: "https://images.unsplash.com/photo-1485827404703-89b55fcc595e?q=80&w=2070&auto=format&fit=crop"
                },
                {
                    id: 5,
                    title: "Шахматный турнир",
                    content: "Состоялся школьный шахматный турнир, в котором приняли участие ученики всех классов. Победители получили дипломы и призы.",
                    created_at: new Date(Date.now() - 345600000).toISOString(),
                    image_url: "https://images.unsplash.com/photo-1529699211952-734e80c4d42b?q=80&w=2071&auto=format&fit=crop"
                },
                {
                    id: 6,
                    title: "Концерт ко Дню учителя",
                    content: "Ученики подготовили праздничный концерт для наших любимых учителей. Спасибо за ваш труд и преданность!",
                    created_at: new Date(Date.now() - 432000000).toISOString(),
                    image_url: "https://images.unsplash.com/photo-1511578314322-379afb476865?q=80&w=2069&auto=format&fit=crop"
                }
            ];

            this.createCarousel(staticArticles, container);
        },
        
        createCarousel(articles, container) {
            // Group articles by 3 per slide
            this.slides = [];
            for (let i = 0; i < articles.length; i += this.newsPerSlide) {
                this.slides.push(articles.slice(i, i + this.newsPerSlide));
            }
            
            if (this.slides.length === 0) {
                container.innerHTML = '<div class="news-no-content">Новостей пока нет.</div>';
                return;
            }
            
            const carouselHTML = `
                <div class="news-carousel-container">
                    <div class="news-carousel">
                        <div class="news-slider" id="news-slider">
                            ${this.slides.map((slide, slideIndex) => this.createSlide(slide, slideIndex)).join('')}
                        </div>
                    </div>
                    ${this.slides.length > 1 ? `
                    <div class="news-carousel-controls">
                        <button class="carousel-btn" id="prev-btn">
                            <i class="fas fa-chevron-left"></i>
                        </button>
                        <div class="carousel-dots">
                            ${this.slides.map((_, index) => 
                                `<div class="carousel-dot ${index === 0 ? 'active' : ''}" 
                                     data-slide="${index}"></div>`
                            ).join('')}
                        </div>
                        <button class="carousel-btn" id="next-btn">
                            <i class="fas fa-chevron-right"></i>
                        </button>
                    </div>
                    ` : ''}
                </div>
            `;
            
            container.innerHTML = carouselHTML;
            
            if (this.slides.length > 1) {
                this.bindEvents();
                this.updateControls();
                this.startAutoPlay();
            }
        },
        
        createSlide(articles, slideIndex) {
            const slideArticles = [...articles];
            // Fill empty spots if there are fewer articles than needed
            while (slideArticles.length < this.newsPerSlide) {
                slideArticles.push(null);
            }
            
            return `
                <div class="news-slide" data-slide="${slideIndex}">
                    ${slideArticles.map(article => article ? this.createNewsCard(article) : '<div class="news-card" style="visibility: hidden;"></div>').join('')}
                </div>
            `;
        },
        
        createNewsCard(article) {
            const date = new Date(article.created_at).toLocaleDateString('ru-RU', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
            
            const excerpt = article.content && article.content.length > 100 
                ? article.content.substring(0, 100) + '...' 
                : article.content || 'Читать полностью...';
                
            const imageHTML = article.image_url 
                ? `<img src="${article.image_url}" alt="${article.title}" class="news-card-image" loading="lazy">`
                : `<div class="news-card-image" style="background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%); display: flex; align-items: center; justify-content: center; color: #6c757d; font-size: 3rem;"><i class="fas fa-newspaper"></i></div>`;
            
            return `
                <div class="news-card" onclick="window.location.href='news_article.html?id=${article.id}'">
                    ${imageHTML}
                    <div class="news-card-content">
                        <h3 class="news-card-title">${article.title || 'Без заголовка'}</h3>
                        <p class="news-card-excerpt">${excerpt}</p>
                        <div class="news-card-date">${date}</div>
                    </div>
                </div>
            `;
        },
        
        bindEvents() {
            const prevBtn = document.getElementById('prev-btn');
            const nextBtn = document.getElementById('next-btn');
            const dots = document.querySelectorAll('.carousel-dot');
            
            if (prevBtn) {
                prevBtn.addEventListener('click', () => this.previousSlide());
            }
            
            if (nextBtn) {
                nextBtn.addEventListener('click', () => this.nextSlide());
            }
            
            dots.forEach((dot, index) => {
                dot.addEventListener('click', () => this.goToSlide(index));
            });
        },
        
        nextSlide() {
            if (this.currentSlide < this.slides.length - 1) {
                this.currentSlide++;
            } else {
                this.currentSlide = 0; // Loop
            }
            this.updateSlider();
            this.resetAutoPlay();
        },
        
        previousSlide() {
            if (this.currentSlide > 0) {
                this.currentSlide--;
            } else {
                this.currentSlide = this.slides.length - 1; // Loop
            }
            this.updateSlider();
            this.resetAutoPlay();
        },
        
        goToSlide(slideIndex) {
            this.currentSlide = slideIndex;
            this.updateSlider();
            this.resetAutoPlay();
        },
        
        updateSlider() {
            const slider = document.getElementById('news-slider');
            if (slider) {
                slider.style.transform = `translateX(-${this.currentSlide * 100}%)`;
                this.updateControls();
            }
        },
        
        updateControls() {
            // Update dots
            document.querySelectorAll('.carousel-dot').forEach((dot, index) => {
                dot.classList.toggle('active', index === this.currentSlide);
            });
        },
        
        startAutoPlay() {
            if (this.slides.length > 1) {
                this.autoplayInterval = setInterval(() => {
                    this.nextSlide();
                }, 5000);
            }
        },
        
        resetAutoPlay() {
            if (this.autoplayInterval) {
                clearInterval(this.autoplayInterval);
                this.startAutoPlay();
            }
        }
    };

    // Initialize news carousel
    newsCarousel.init();

    // Language Switcher
    const languageSwitcher = document.querySelector('.language-switcher');
    const langButton = document.querySelector('.lang-button');
    const currentLangSpan = document.querySelector('.current-lang');
    const langOptions = document.querySelectorAll('.lang-option');

    const setLanguage = (lang) => {
        document.documentElement.lang = lang;
        document.querySelectorAll('[data-i18n-key]').forEach(element => {
            const key = element.getAttribute('data-i18n-key');
            if (translations[lang] && translations[lang][key]) {
                if (element.hasAttribute('data-i18n-key-placeholder')) {
                    element.placeholder = translations[lang][key];
                } else {
                    element.innerHTML = translations[lang][key];
                }
            }
        });
        // Update button text
        const langLabels = {
            'ru': 'РУС',
            'kz': 'ҚАЗ',
            'en': 'ENG'
        };
        currentLangSpan.textContent = langLabels[lang] || 'РУС';
        localStorage.setItem('language', lang);
        languageSwitcher.classList.remove('active');
        
    };

    if (langButton) {
        langButton.addEventListener('click', (e) => {
            e.stopPropagation();
            languageSwitcher.classList.toggle('active');
        });
    }

    document.addEventListener('click', () => {
        if (languageSwitcher && languageSwitcher.classList.contains('active')) {
            languageSwitcher.classList.remove('active');
        }
    });

    langOptions.forEach(option => {
        option.addEventListener('click', (e) => {
            e.preventDefault();
            const selectedLang = option.getAttribute('data-lang');
            setLanguage(selectedLang);
        });
    });

    // Load saved language or default to 'ru'
    const savedLang = localStorage.getItem('language') || 'ru';
    setLanguage(savedLang);

    // Header Scroll Effect
    const header = document.getElementById('main-header');
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            header.classList.add('scrolled');
        } else {
            header.classList.remove('scrolled');
        }
    });

    // Mobile Menu
    const menuToggle = document.querySelector('.mobile-menu-toggle');
    const navMenu = document.querySelector('.nav-menu');
    
    if (menuToggle && navMenu) {
        menuToggle.addEventListener('click', () => {
            navMenu.classList.toggle('active');
            const icon = menuToggle.querySelector('i');
            icon.classList.toggle('fa-bars');
            icon.classList.toggle('fa-times');
        });
        
        document.querySelectorAll('.nav-link').forEach(link => {
            link.addEventListener('click', () => {
                if (navMenu.classList.contains('active')) {
                    navMenu.classList.remove('active');
                    const icon = menuToggle.querySelector('i');
                    icon.classList.add('fa-bars');
                    icon.classList.remove('fa-times');
                }
            });
        });
    }

    // Active Nav Link Highlighting on Scroll
    const sections = document.querySelectorAll('section[id]');
    const navLinks = document.querySelectorAll('.nav-menu a');
    
    window.addEventListener('scroll', () => {
        let current = '';
        sections.forEach(section => {
            const sectionTop = section.offsetTop;
            if (pageYOffset >= sectionTop - 150) {
                current = section.getAttribute('id');
            }
        });
        
        navLinks.forEach(link => {
            link.classList.remove('active');
            if (link.getAttribute('href').includes(current)) {
                link.classList.add('active');
            }
        });
    });

    // Scroll Animations
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
            }
        });
    }, { threshold: 0.1 });
    
    document.querySelectorAll('.fade-in').forEach(el => {
        observer.observe(el);
    });

    // Animated Counters
    const counters = document.querySelectorAll('.stat-number');
    const speed = 200;
    
    const counterObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const counter = entry.target;
                const updateCount = () => {
                    const target = +counter.getAttribute('data-target');
                    const count = +counter.innerText;
                    const inc = target / speed;
                    
                    if (count < target) {
                        counter.innerText = Math.ceil(count + inc);
                        setTimeout(updateCount, 10);
                    } else {
                        counter.innerText = target;
                    }
                };
                updateCount();
                observer.unobserve(counter);
            }
        });
    }, { threshold: 0.5 });
    
    counters.forEach(counter => counterObserver.observe(counter));

    // Contact Form Submission
    const form = document.getElementById('main-contact-form');
    
    if (form) {
        form.addEventListener('submit', async function(event) {
            event.preventDefault();

            let isValid = true;
            let errorMessages = [];
            const name = document.getElementById('name');
            const email = document.getElementById('email');
            const phone = document.getElementById('phone');
            const message = document.getElementById('message');

            // Clear previous errors
            [name, email, phone, message].forEach(el => {
                el.classList.remove('invalid');
                el.style.borderColor = '';
            });

            // Validate name
            if (name.value.trim() === '') {
                isValid = false;
                name.classList.add('invalid');
                errorMessages.push('• ' + translations[savedLang]['contact.form.errorName']);
            }
            
            // Validate email
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(email.value)) {
                isValid = false;
                email.classList.add('invalid');
                errorMessages.push('• ' + translations[savedLang]['contact.form.errorEmail']);
            }
            
            // Validate phone
            if (phone.value.trim() === '') {
                isValid = false;
                phone.classList.add('invalid');
                errorMessages.push('• ' + translations[savedLang]['contact.form.errorPhone']);
            }
            
            // Validate message
            if (message.value.trim() === '') {
                isValid = false;
                message.classList.add('invalid');
                errorMessages.push('• ' + translations[savedLang]['contact.form.errorMessage']);
            }

            if (!isValid) {
                showValidationError(errorMessages.join('\n'));
                return;
            }

            const formData = {
                name: name.value,
                email: email.value,
                phone: phone.value,
                message: message.value
            };

            const submitButton = form.querySelector('button[type="submit"]');
            submitButton.disabled = true;
            submitButton.textContent = translations[savedLang]['contact.form.sending'];

            try {
                const response = await fetch('/api/contact', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData),
                });

                if (response.ok) {
                    const data = await response.json();
                    form.reset();
                    showSuccessMessage(data.message || translations[savedLang]['contact.form.success']);
                } else {
                    let errorMessage = translations[savedLang]['contact.form.error'];
                    try {
                        const errorText = await response.text();
                        if (errorText) errorMessage = errorText;
                    } catch (e) {}
                    showErrorMessage(errorMessage);
                }
            } catch (error) {
                showErrorMessage(translations[savedLang]['contact.form.error']);
            } finally {
                submitButton.disabled = false;
                submitButton.textContent = translations[savedLang]['contact.form.submit'];
            }
        });
    }

    function showValidationError(messages) {
        const errorDiv = document.createElement('div');
        errorDiv.style.cssText = 'position: fixed; top: 100px; left: 50%; transform: translateX(-50%); background: #fee2e2; color: #991b1b; padding: 1rem 1.5rem; border-radius: 0.5rem; border: 2px solid #ef4444; z-index: 10000; max-width: 500px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); white-space: pre-line; text-align: left; font-size: 0.9rem; line-height: 1.6;';
        errorDiv.innerHTML = `<strong style="display: block; margin-bottom: 0.5rem;">Пожалуйста, исправьте ошибки:</strong>${messages}`;
        document.body.appendChild(errorDiv);
        
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.style.opacity = '0';
                errorDiv.style.transition = 'opacity 0.3s ease';
                setTimeout(() => errorDiv.parentNode.removeChild(errorDiv), 300);
            }
        }, 5000);
    }

    function showSuccessMessage(message) {
        const successDiv = document.createElement('div');
        successDiv.textContent = message;
        successDiv.style.cssText = 'position: fixed; top: 100px; left: 50%; transform: translateX(-50%); background: #d1fae5; color: #065f46; padding: 1rem 1.5rem; border-radius: 0.5rem; border: 2px solid #10b981; z-index: 10000; max-width: 500px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); font-weight: 500;';
        document.body.appendChild(successDiv);
        
        setTimeout(() => {
            if (successDiv.parentNode) {
                successDiv.style.opacity = '0';
                successDiv.style.transition = 'opacity 0.3s ease';
                setTimeout(() => successDiv.parentNode.removeChild(successDiv), 300);
            }
        }, 5000);
    }
    
    function showErrorMessage(message) {
        const errorDiv = document.createElement('div');
        errorDiv.textContent = message;
        errorDiv.style.cssText = 'position: fixed; top: 100px; left: 50%; transform: translateX(-50%); background: #fee2e2; color: #991b1b; padding: 1rem 1.5rem; border-radius: 0.5rem; border: 2px solid #ef4444; z-index: 10000; max-width: 500px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); font-weight: 500;';
        document.body.appendChild(errorDiv);
        
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.style.opacity = '0';
                errorDiv.style.transition = 'opacity 0.3s ease';
                setTimeout(() => errorDiv.parentNode.removeChild(errorDiv), 300);
            }
        }, 5000);
    }

    // Smooth Scrolling for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
});