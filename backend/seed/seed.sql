-- CATEGORIAS
INSERT OR IGNORE INTO categorias (id, nome, cor, ordem) VALUES
(1, 'Proteínas',         '#C0392B', 1),
(2, 'Carboidratos',      '#C4850A', 2),
(3, 'Laticínios / Frios','#2980B9', 3),
(4, 'Hortifruti',        '#2D7A5F', 4),
(5, 'Bebidas',           '#6C5CE7', 5),
(6, 'Temperos / Básicos','#C4850A', 6),
(7, 'Limpeza / Higiene', '#636E72', 7),
(8, 'Snacks',            '#E17055', 8);

-- PRODUTOS
INSERT OR IGNORE INTO produtos (id, nome, categoria_id, unidade, unidade_preco, preco_ref) VALUES
(1, 'Filé de peito frango Sadia IQF 1kg', 1, 'un', 'un', 17.98),
(2, 'Empanado de frango', 1, 'un', 'un', 2.00),
(3, 'Salsicha 500g', 1, 'pct', 'pct', 6.00),
(4, 'Acém (panela)', 1, 'kg', 'kg', 36.00),
(5, 'Ovos (dúzia)', 1, 'dz', 'dz', 9.93),
(6, 'Mortadela Sadia def.fat. 200g', 1, 'un', 'un', 8.78),
(7, 'Arroz parboilizado 3kg', 2, 'pct', 'pct', 20.00),
(8, 'Arroz integral Urbano 1kg', 2, 'un', 'un', 4.38),
(9, 'Macarrão 500g', 2, 'pct', 'pct', 3.31),
(10, 'Farinha de trigo 5kg', 2, 'pct', 'pct', 13.67),
(11, 'Pão Bauducco 37 integral 390g', 2, 'un', 'un', 8.98),
(12, 'Batata inglesa', 2, 'kg', 'kg', 4.46),
(13, 'Leite Tirol desnatado 1L', 3, 'un', 'L', 4.98),
(14, 'Queijo fatiado 300g', 3, 'emb', 'emb', 17.00),
(15, 'Margarina Qualy 500g', 3, 'un', 'un', 11.00),
(16, 'Presunto Sadia 150g', 3, 'emb', 'emb', 11.00),
(17, 'Requeijão Laclelo 180g', 3, 'un', 'un', 7.18),
(18, 'Laranja pera', 4, 'kg', 'kg', 3.47),
(19, 'Banana branca', 4, 'kg', 'kg', 6.89),
(20, 'Tomate longa vida', 4, 'kg', 'kg', 5.24),
(21, 'Cebola', 4, 'kg', 'kg', 3.98),
(22, 'Alho 200g', 4, 'un', 'un', 7.19),
(23, 'Cenoura', 4, 'kg', 'kg', 3.86),
(24, 'Alface americana hidropônica 150g', 4, 'un', 'un', 6.48),
(25, 'Repolho verde', 4, 'un', 'un', 6.18),
(26, 'Beterraba', 4, 'kg', 'kg', 8.48),
(27, 'Limão tahiti', 4, 'kg', 'kg', 3.58),
(28, 'Mamão papaya', 4, 'un', 'un', 6.98),
(29, 'Coca-Cola s/ açúcar 2L', 5, 'un', 'un', 10.98),
(30, 'Óleo girassol Liza 900ml', 6, 'un', 'un', 18.68),
(31, 'Óleo soja 900ml', 6, 'un', 'un', 7.06),
(32, 'Açúcar 5kg', 6, 'pct', 'pct', 17.84),
(33, 'Café 500g', 6, 'pct', 'pct', 23.52),
(34, 'Sal 1kg', 6, 'un', 'un', 1.99),
(35, 'Extrato tomate Elefante 300g', 6, 'un', 'un', 5.48),
(36, 'Ketchup Heinz 397g', 6, 'un', 'un', 12.98),
(37, 'Maionese', 6, 'un', 'un', 6.31),
(38, 'Vinagre', 6, 'un', 'un', 2.01),
(39, 'Detergente 500ml', 7, 'un', 'un', 1.91),
(40, 'Sabão em pó', 7, 'un', 'un', 7.07),
(41, 'Água sanitária 1L', 7, 'un', 'un', 3.12),
(42, 'Desinfetante 500ml', 7, 'un', 'un', 3.77),
(43, 'Sabão barra 5un', 7, 'pct', 'pct', 11.67),
(44, 'Papel higiênico 4 rolos', 7, 'pct', 'pct', 5.95),
(45, 'Creme dental', 7, 'un', 'un', 3.13),
(46, 'Sabonete', 7, 'un', 'un', 1.68),
(47, 'Shampoo', 7, 'un', 'un', 8.31),
(48, 'Desodorante', 7, 'un', 'un', 9.50),
(49, 'Baconzitos 86g', 8, 'un', 'un', 8.98),
(50, 'Biscoito Isabela maizena 350g', 8, 'un', 'un', 4.98);

-- LISTA ABRIL 2026
INSERT OR IGNORE INTO listas (id, nome, mes, ano, renda) VALUES (1, 'Abril 2026', 4, 2026, 6000);

-- ITENS DA LISTA
INSERT OR IGNORE INTO lista_itens (lista_id, produto_id, qtd, preco_estimado, duracao_meses, observacao) VALUES
(1, 1, 10, 17.98, 2, 'Frango congelado — rende 2 meses'),
(1, 2, 30, 2.00, 1, NULL),
(1, 3, 3, 6.00, 1, NULL),
(1, 4, 1.5, 36.00, 1, NULL),
(1, 5, 4, 9.93, 1, NULL),
(1, 6, 1, 8.78, 1, NULL),
(1, 7, 2, 20.00, 1, NULL),
(1, 8, 1, 4.38, 1, NULL),
(1, 9, 6, 3.31, 1, NULL),
(1, 10, 0.5, 13.67, 1, NULL),
(1, 11, 10, 8.98, 1, 'Pão integral — estoque mês inteiro'),
(1, 12, 3, 4.46, 1, NULL),
(1, 13, 12, 4.98, 1, NULL),
(1, 14, 2, 17.00, 1, NULL),
(1, 15, 2, 11.00, 1, NULL),
(1, 16, 3, 11.00, 1, NULL),
(1, 17, 1, 7.18, 1, NULL),
(1, 18, 5, 3.47, 1, NULL),
(1, 19, 3, 6.89, 1, NULL),
(1, 20, 3, 5.24, 1, NULL),
(1, 21, 2, 3.98, 1, NULL),
(1, 22, 2, 7.19, 1, NULL),
(1, 23, 2, 3.86, 1, NULL),
(1, 24, 6, 6.48, 1, NULL),
(1, 25, 1, 6.18, 1, NULL),
(1, 26, 1, 8.48, 1, NULL),
(1, 27, 1, 3.58, 1, NULL),
(1, 28, 1, 6.98, 1, NULL),
(1, 29, 7, 10.98, 1, NULL),
(1, 30, 1, 18.68, 1, NULL),
(1, 31, 2, 7.06, 1, NULL),
(1, 32, 0.5, 17.84, 1, NULL),
(1, 33, 2, 23.52, 1, NULL),
(1, 34, 1, 1.99, 1, NULL),
(1, 35, 3, 5.48, 1, NULL),
(1, 36, 1, 12.98, 1, NULL),
(1, 37, 1, 6.31, 1, NULL),
(1, 38, 1, 2.01, 1, NULL),
(1, 39, 3, 1.91, 1, NULL),
(1, 40, 2, 7.07, 1, NULL),
(1, 41, 2, 3.12, 1, NULL),
(1, 42, 2, 3.77, 1, NULL),
(1, 43, 1, 11.67, 1, NULL),
(1, 44, 3, 5.95, 1, NULL),
(1, 45, 2, 3.13, 1, NULL),
(1, 46, 6, 1.68, 1, NULL),
(1, 47, 2, 8.31, 1, NULL),
(1, 48, 2, 9.50, 1, NULL),
(1, 49, 1, 8.98, 1, NULL),
(1, 50, 2, 4.98, 1, NULL);

-- COMPRA #1 — Supermercado Local 03/04/2026
INSERT OR IGNORE INTO compras (id, lista_id, local, data, total_nfe, observacao) VALUES
(1, 1, 'Supermercado Local', '2026-04-03', 444.96, 'Compra do mês — exemplo');

-- ITENS COMPRA #1
INSERT OR IGNORE INTO compra_itens (compra_id, produto_id, nome_nfe, qtd, unidade, preco_unit, preco_total) VALUES
(1, 21, 'CEBOLA SOLTA AV', 0.90, 'kg', 3.98, 3.58),
(1, 26, 'BETERRABA AV', 0.48, 'kg', 8.48, 4.07),
(1, 20, 'TOMATE LONGA VIDA AV', 0.414, 'kg', 5.24, 2.17),
(1, 19, 'BANANA BRANCA AV', 1.086, 'kg', 6.89, 7.48),
(1, 27, 'LIMAO TAHITI AV', 0.464, 'kg', 3.58, 1.66),
(1, 12, 'BATATA INGLESA BRANCA SOLTA AV', 0.744, 'kg', 4.46, 3.32),
(1, 18, 'LARANJA PERA AGRANEL AV', 2.924, 'kg', 3.47, 10.16),
(1, 25, 'REPOLHO VERDE AV 1UN', 1.0, 'un', 6.18, 6.18),
(1, 24, 'ALFACE AMERIC.COR.HIDR.BLUMENAU 150G', 1.0, 'un', 6.48, 6.48),
(1, 28, 'MAMAO PAPAYA/AMAZONAS AV 1UN', 1.0, 'un', 6.98, 6.98),
(1, 8, 'ARROZ URBANO INTEGRAL T-1 PE 1KG', 1.0, 'un', 4.38, 4.38),
(1, 11, 'PAO BAUDUCCO 37 INTEGRAL PE 390G', 10.0, 'un', 8.98, 89.80),
(1, 17, 'REQUEIJAO LACLELO TRADICIONAL 180G', 1.0, 'un', 7.18, 7.18),
(1, 49, 'SALG. BACONZITOS CLASSICO PE 86G', 1.0, 'un', 8.98, 8.98),
(1, 1, 'FILE DE PEITO FGO.SADIA IQF PE 1KG', 10.0, 'un', 17.98, 179.80),
(1, 13, 'LEITE TIROL DESN.C/TP TP 1LT', 5.0, 'un', 4.98, 24.90),
(1, 29, 'REFRIG. COCA COLA S/ACUCAR FC 2LT', 2.0, 'un', 10.98, 21.96),
(1, 36, 'KETCHUP HEINZ TRADICIONAL FC 397G', 1.0, 'un', 12.98, 12.98),
(1, 30, 'OLEO GIRASSOL LIZA FC 900ML', 1.0, 'un', 18.68, 18.68),
(1, 35, 'EXTRATO TOMATE ELEFANTE 300G', 1.0, 'un', 5.48, 5.48),
(1, 50, 'BISC.ISABELA MAIZENA PE 350G', 2.0, 'un', 4.98, 9.96),
(1, 6, 'MORT. SADIA DEF.FAT.SOLTISSIMO 200G', 1.0, 'un', 8.78, 8.78);

-- Marcar itens da lista que foram comprados na compra #1
-- ATENÇÃO: se o banco já foi populado sem este UPDATE, é necessário deletar o
-- volume Docker e recriar o container para que o seed seja reexecutado.
-- Ex: docker compose down -v && docker compose up --build
UPDATE lista_itens SET comprado = 1
WHERE lista_id = 1 AND produto_id IN (
    SELECT produto_id FROM compra_itens WHERE compra_id = 1 AND produto_id IS NOT NULL
);
