--
-- PostgreSQL database dump
--

\restrict Aa25jvB42U36swgA8bHT3UgtSjvH2lLwn7cbdU3urMX9IGH05e0TvsqOpC21X7d

-- Dumped from database version 17.10
-- Dumped by pg_dump version 17.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.comments DROP CONSTRAINT IF EXISTS comments_article_id_fkey;
DROP INDEX IF EXISTS public.comments_article_id_created_at_idx;
ALTER TABLE IF EXISTS ONLY public.comments DROP CONSTRAINT IF EXISTS comments_pkey;
ALTER TABLE IF EXISTS ONLY public.articles DROP CONSTRAINT IF EXISTS articles_slug_key;
ALTER TABLE IF EXISTS ONLY public.articles DROP CONSTRAINT IF EXISTS articles_pkey;
ALTER TABLE IF EXISTS ONLY public.admins DROP CONSTRAINT IF EXISTS admins_username_key;
ALTER TABLE IF EXISTS ONLY public.admins DROP CONSTRAINT IF EXISTS admins_pkey;
DROP TABLE IF EXISTS public.comments;
DROP TABLE IF EXISTS public.articles;
DROP TABLE IF EXISTS public.admins;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: admins; Type: TABLE; Schema: public; Owner: blog
--

CREATE TABLE public.admins (
    id bigint NOT NULL,
    username character varying(100) NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.admins OWNER TO blog;

--
-- Name: admins_id_seq; Type: SEQUENCE; Schema: public; Owner: blog
--

ALTER TABLE public.admins ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.admins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: articles; Type: TABLE; Schema: public; Owner: blog
--

CREATE TABLE public.articles (
    id bigint NOT NULL,
    title character varying(200) NOT NULL,
    slug character varying(200) NOT NULL,
    summary text DEFAULT ''::text NOT NULL,
    content text NOT NULL,
    status character varying(20) DEFAULT 'draft'::character varying NOT NULL,
    published_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT articles_status_check CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'published'::character varying])::text[])))
);


ALTER TABLE public.articles OWNER TO blog;

--
-- Name: articles_id_seq; Type: SEQUENCE; Schema: public; Owner: blog
--

ALTER TABLE public.articles ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: comments; Type: TABLE; Schema: public; Owner: blog
--

CREATE TABLE public.comments (
    id bigint NOT NULL,
    article_id bigint NOT NULL,
    author_name character varying(50) NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT comments_author_name_not_blank CHECK ((length(TRIM(BOTH FROM author_name)) > 0)),
    CONSTRAINT comments_content_not_blank CHECK ((length(TRIM(BOTH FROM content)) > 0))
);


ALTER TABLE public.comments OWNER TO blog;

--
-- Name: comments_id_seq; Type: SEQUENCE; Schema: public; Owner: blog
--

ALTER TABLE public.comments ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.comments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Data for Name: admins; Type: TABLE DATA; Schema: public; Owner: blog
--

COPY public.admins (id, username, password_hash, created_at) FROM stdin;
1	admin	$2a$10$Gu5vK.ZENG.CKMEKKcYY.OWkpwBgz9it51ij8.YKo415jnmxUs5US	2026-08-18 13:46:48.504204+00
\.


--
-- Data for Name: articles; Type: TABLE DATA; Schema: public; Owner: blog
--

COPY public.articles (id, title, slug, summary, content, status, published_at, created_at, updated_at) FROM stdin;
1	我的第一篇文章	my-first-post	这是测试摘要	# Hello PostgreSQL	draft	\N	2026-08-15 06:17:28.080471+00	2026-08-15 06:17:28.080471+00
5	我的第一篇文章	my-first-post-1	这是我的第一篇博客文章	# Hello Blog	published	\N	2026-08-16 03:52:31.469899+00	2026-08-16 03:52:31.469899+00
6	探索未知的边际：为什么我们依然热爱徒步？	1	在屏幕与算法定义日常的时代，双手触碰真实的泥土、双脚丈量起伏的山脊，正成为一种罕见的精神疗愈。徒步（Hiking）早已超越了单纯的体育运动，它演变成现代人主动选择的“数字断连”与自我对话。	# 探索未知的边际：为什么我们依然热爱徒步？\n\n在屏幕与算法定义日常的时代，双手触碰真实的泥土、双脚丈量起伏的山脊，正成为一种罕见的精神疗愈。徒步（Hiking）早已超越了单纯的体育运动，它演变成现代人主动选择的“数字断连”与自我对话。\n\n---\n\n## 一、 徒步的物理与心理图谱\n\n为何越来越多的人选择背起行囊走入荒野？科学与心理学给出了清晰的解释：\n\n* **皮质醇下降**：研究表明，在森林环境（Shinrin-yoku）中行走超过 20 分钟，人体内的压力荷尔蒙——皮质醇水平会显著降低。\n* **内源性内啡肽释放**：有氧运动配合自然的视觉刺激，能持续激发大脑产生愉悦感与专注力。\n* **时空感的重构**：离线状态下，时间的度量从“秒与通知推送”回归到“日照倾斜的角度与脚步的节奏”。\n\n> “山就在那里，但走过山的人，不再是原来的自己。”\n\n---\n\n## 二、 基础装备选择对比\n\n对于初学者而言，装备的选择决定了户外体验的舒适度与安全性。以下是核心三件套的挑选维度：\n\n| 装备类型 | 关键指标 | 推荐选择建议 | 避坑提示 |\n| :--- | :--- | :--- | :--- |\n| **徒步鞋** | 抓地力、支撑性、防水性 | 中高帮 Gore-Tex 材质，Vibram 大底 | 切忌穿未经磨合的新鞋走长线 |\n| **背包** | 背负系统、容量、自重 | 单日徒步选择 **18L - 25L**，含腰带支撑 | 避免重量集中在双肩 |\n| **穿搭系统** | 排汗、保暖、防风 | 遵循**三层穿衣法**（基础层+中间层+外壳层） | 严禁穿纯棉内衣（湿后易失温） |\n\n---\n\n## 三、 行前准备步骤\n\n一次安全且愉悦的徒步，离不开严谨的行前规划。\n\n<Sequence>\n\n  <Step subtitle="安全第一" title="规划路线与评估难度">\n    检查路线的升降落差、里程及地形（泥地、碎石或岩石），确保难度与自身体力匹配。\n  </Step>\n\n  <Step title="查验天气与日落时间">\n    查看山区的微气候预测，根据日落时间倒排出发时刻，预留至少 2 小时安全余量。\n  </Step>\n\n  <Step title="准备物资与轨迹下载">\n    准备充足的水（建议 1.5L-2L）、高热量路餐（如坚果、能量棒），并将离线轨迹导入导航软件（如两步路、户外助手）。\n  </Step>\n\n  <Step title="报备行程（留守人制度）">\n    将预定路线、预计返回时间及紧急联系人告知未同行的家人或朋友。\n  </Step>\n\n</Sequence>\n\n---\n\n## 四、 践行无痕山林（Leave No Trace）\n\n走到野外，我们既是探索者，也是守护者。谨记 **LNT 七大原则** 的核心：\n\n1. **带走所有垃圾**：包括果皮与纸巾（果皮在野外分解可能需要数月至数年）。\n2. **保持环境原貌**：不折采花木，不带走石头与古迹。\n3. **尊重野生动物**：保持距离观察，切勿喂食。\n\n再次出发时，愿你在山野间找回属于自己的节奏与平静。	published	\N	2026-08-16 04:58:37.148523+00	2026-08-16 04:58:37.148523+00
7	快速排序（Quick Sort）算法详解	2	快速排序是一种基于**分治法（Divide and Conquer）**的高效排序算法。它的平均时间复杂度为 $O(n \\log n)$，在大多数实际场景中表现优异。	# 快速排序（Quick Sort）算法详解\n\n快速排序是一种基于**分治法（Divide and Conquer）**的高效排序算法。它的平均时间复杂度为 $O(n \\log n)$，在大多数实际场景中表现优异。\n\n---\n\n## 一、算法核心思想\n\n1. **选择基准值（Pivot）**：从数组中挑选一个元素作为基准。\n2. **分区（Partition）**：重新排序数组，所有比基准值小的元素摆放在基准前面，所有比基准值大的元素摆在基准后面。\n3. **递归排序（Recursively Sort）**：递归地把小于基准值元素的子序列和大于基准值元素的子序列排序。\n\n---\n\n## 二、代码实现（Go 语言）\n\n```go\npackage main\n\nimport "fmt"\n\n// QuickSort 快速排序主入口\nfunc QuickSort(arr []int) []int {\nif len(arr) <= 1 {\nreturn arr\n}\nquickSortHelper(arr, 0, len(arr)-1)\nreturn arr\n}\n\n// quickSortHelper 递归切分与排序\nfunc quickSortHelper(arr []int, low, high int) {\nif low < high {\n// p 为分区后基准值的最终正确位置\np := partition(arr, low, high)\n\n// 递归排序基准左侧与右侧的子数组\nquickSortHelper(arr, low, p-1)\nquickSortHelper(arr, p+1, high)\n}\n}\n\n// partition 原地分区函数（以最后一个元素作为 Pivot）\nfunc partition(arr []int, low, high int) int {\npivot := arr[high]\ni := low - 1\n\nfor j := low; j < high; j++ {\nif arr[j] < pivot {\ni++\narr[i], arr[j] = arr[j], arr[i]\n}\n}\n\n// 将基准值放到正确的位置上\narr[i+1], arr[high] = arr[high], arr[i+1]\nreturn i + 1\n}\n\nfunc main() {\nnumbers := []int{34, 7, 23, 32, 5, 62}\nfmt.Println("排序前:", numbers)\n\nQuickSort(numbers)\nfmt.Println("排序后:", numbers)\n}	published	\N	2026-08-16 09:05:32.289873+00	2026-08-16 09:05:32.289873+00
8	《沉思录》摘要	3	古罗马皇帝奥勒留是希腊五贤帝之一	## 《沉思录》\n> 你的心灵也许会破碎，但你依然会像从前一样生活\n\n> 找不到比灵魂更安静、更无忧无虑的避风港了\n\n> 最好的报复，不是让自己变得跟上海你的人一样，而是超越他，成为一个跟他完全不同的人\n\n> 专注生命的美好，仰望星空，想象自己和它们一起奔跑	published	2026-08-18 14:00:54.211996+00	2026-08-18 14:00:46.113593+00	2026-08-18 14:07:54.781776+00
9	KDA算子	kda-suan-zi	KDA算子是KimiK3模型用的线性注意力算子，是GDN算子的一个变形，主要对阿尔法进行了改造，把原来标量模式改为了主对角线的模式	略	published	2026-08-18 14:31:24.829352+00	2026-08-18 14:31:10.499299+00	2026-08-18 14:31:24.829352+00
\.


--
-- Data for Name: comments; Type: TABLE DATA; Schema: public; Owner: blog
--

COPY public.comments (id, article_id, author_name, content, created_at) FROM stdin;
1	8	Steven	hello	2026-08-20 07:45:42.800164+00
2	8	指头女巫	即使引导早已破碎，也请你当上埃尔登之王	2026-08-20 07:48:19.572481+00
\.


--
-- Name: admins_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blog
--

SELECT pg_catalog.setval('public.admins_id_seq', 1, true);


--
-- Name: articles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blog
--

SELECT pg_catalog.setval('public.articles_id_seq', 9, true);


--
-- Name: comments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blog
--

SELECT pg_catalog.setval('public.comments_id_seq', 2, true);


--
-- Name: admins admins_pkey; Type: CONSTRAINT; Schema: public; Owner: blog
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_pkey PRIMARY KEY (id);


--
-- Name: admins admins_username_key; Type: CONSTRAINT; Schema: public; Owner: blog
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_username_key UNIQUE (username);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: blog
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);


--
-- Name: articles articles_slug_key; Type: CONSTRAINT; Schema: public; Owner: blog
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_slug_key UNIQUE (slug);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: blog
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: comments_article_id_created_at_idx; Type: INDEX; Schema: public; Owner: blog
--

CREATE INDEX comments_article_id_created_at_idx ON public.comments USING btree (article_id, created_at);


--
-- Name: comments comments_article_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blog
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_article_id_fkey FOREIGN KEY (article_id) REFERENCES public.articles(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict Aa25jvB42U36swgA8bHT3UgtSjvH2lLwn7cbdU3urMX9IGH05e0TvsqOpC21X7d

