-- 테스트 데이터
INSERT INTO accounts (id, account_number, balance) VALUES (1, 'AC001', 1000000) ON CONFLICT (id) DO NOTHING;
INSERT INTO holdings (account_id, stock_code, quantity) VALUES (1, 'STOCK01', 100) ON CONFLICT (account_id, stock_code) DO NOTHING;