CREATE TABLE IF NOT EXISTS categorias (
    id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    cor TEXT NOT NULL DEFAULT '#636E72',
    icone TEXT,
    ordem INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS produtos (
    id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    categoria_id INTEGER NOT NULL REFERENCES categorias(id),
    unidade TEXT NOT NULL DEFAULT 'un',
    unidade_preco TEXT NOT NULL DEFAULT 'un',
    preco_ref REAL NOT NULL DEFAULT 0,
    ativo INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS listas (
    id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    mes INTEGER NOT NULL,
    ano INTEGER NOT NULL,
    renda REAL NOT NULL DEFAULT 6000,
    status TEXT NOT NULL DEFAULT 'ativa',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS lista_itens (
    id SERIAL PRIMARY KEY,
    lista_id INTEGER NOT NULL REFERENCES listas(id) ON DELETE CASCADE,
    produto_id INTEGER NOT NULL REFERENCES produtos(id),
    qtd REAL NOT NULL DEFAULT 1,
    preco_estimado REAL NOT NULL DEFAULT 0,
    duracao_meses INTEGER NOT NULL DEFAULT 1,
    comprado INTEGER NOT NULL DEFAULT 0,
    observacao TEXT
);

CREATE TABLE IF NOT EXISTS compras (
    id SERIAL PRIMARY KEY,
    lista_id INTEGER REFERENCES listas(id),
    local TEXT NOT NULL DEFAULT '',
    data DATE NOT NULL DEFAULT CURRENT_DATE,
    total_nfe REAL,
    observacao TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS compra_itens (
    id SERIAL PRIMARY KEY,
    compra_id INTEGER NOT NULL REFERENCES compras(id) ON DELETE CASCADE,
    produto_id INTEGER REFERENCES produtos(id),
    nome_nfe TEXT NOT NULL,
    qtd REAL NOT NULL DEFAULT 1,
    unidade TEXT NOT NULL DEFAULT 'un',
    preco_unit REAL NOT NULL DEFAULT 0,
    preco_total REAL NOT NULL DEFAULT 0,
    lista_item_id INTEGER REFERENCES lista_itens(id)
);

CREATE INDEX IF NOT EXISTS idx_lista_itens_lista ON lista_itens(lista_id);
CREATE INDEX IF NOT EXISTS idx_compra_itens_compra ON compra_itens(compra_id);
CREATE INDEX IF NOT EXISTS idx_compras_lista ON compras(lista_id);
CREATE INDEX IF NOT EXISTS idx_produtos_categoria ON produtos(categoria_id);
