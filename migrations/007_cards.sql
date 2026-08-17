CREATE TABLE carts(
    id SERIAL PRIMARY KEY, 
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Sepetin kime ait olduğu 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_items(
    id SERIAL PRIMARY KEY,
    cart_id INTEGER NOT NULL REFERENCES carts(id) ON DELETE CASCADE, -- Ürün hangi sepette duruyor ? Sepet silinirse içindekilerde silinir
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE, -- Sepete eklene ürün hangisi ?
    quantity INTEGER NOT NULL -- Kaç tane ürün eklendi ?
);