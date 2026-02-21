-- 1. 従業員データの投入
INSERT INTO employees (id, name, email, password_hash, role) VALUES
('11111111-1111-1111-1111-111111111111', '山田 店長', 'manager@example.com', 'hashed_password_here', 'manager'),
('22222222-2222-2222-2222-222222222222', '佐藤 エージェント', 'agent1@example.com', 'hashed_password_here', 'agent'),
('33333333-3333-3333-3333-333333333333', '鈴木 エージェント', 'agent2@example.com', 'hashed_password_here', 'agent'),
('00000000-0000-0000-0000-000000000000', '具志堅 管理者', 'admin@example.com', 'hashed_password_here', 'admin');

-- 2. 顧客データの投入
INSERT INTO customers (id, name, email, phone) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '田中 太郎', 'tanaka@example.com', '090-1234-5678'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '伊藤 花子', 'ito@example.com', '080-9876-5432');

-- 3. 物件データの投入
INSERT INTO properties (id, name, rent, address, layout, status) VALUES
('xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx', 'メゾン渋谷', 115000, '東京都渋谷区神南1-1-1', '1LDK', 'available'),
('yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy', 'グランド世田谷', 78000, '東京都世田谷区三軒茶屋2-2-2', '1K', 'available');

-- 4. 案件データの投入（カンバン表示用）
INSERT INTO deals (id, customer_id, property_id, assignee_id, status, move_in_date) VALUES
-- 田中さんの案件：山田店長が担当し、内見予定
('dddddddd-1111-dddd-1111-dddddddddddd', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx', '11111111-1111-1111-1111-111111111111', 'viewing_scheduled', '2026-03-01'),
-- 伊藤さんの案件：佐藤エージェントが担当し、追客中
('dddddddd-2222-dddd-2222-dddddddddddd', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy', '22222222-2222-2222-2222-222222222222', 'following_up', '2026-04-01');