-- User Table
CREATE TABLE User (
    user_id INT PRIMARY KEY,
    username VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    address VARCHAR(255),
    phone_number VARCHAR(20)
);

-- Product Table
CREATE TABLE Product (
    product_id INT PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    price DECIMAL(10, 2),
    image_url VARCHAR(255),
    category_id INT,
    FOREIGN KEY (category_id) REFERENCES Category(category_id)
);

-- Category Table
CREATE TABLE Category (
    category_id INT PRIMARY KEY,
    name VARCHAR(255)
);


INSERT INTO Category (category_id, name)
VALUES
    (1, 'Clothing'),
    (2, 'Accessories'),
    (3, 'Shoes');
ALTER TABLE product MODIFY product_id INT AUTO_INCREMENT;


INSERT INTO Product (product_id, name, description, price, image_url, category_id)
VALUES
    (1, 'T-Shirt', 'Comfortable cotton t-shirt', 15.99, 'tshirt.jpg', 1),
    (2, 'Jeans', 'Slim-fit denim jeans', 39.99, 'jeans.jpg', 1),
    (3, 'Sneakers', 'Classic canvas sneakers', 29.99, 'sneakers.jpg', 3),
    (4, 'Handbag', 'Stylish leather handbag', 49.99, 'handbag.jpg', 2),
    (5, 'Dress', 'Elegant evening dress', 59.99, 'dress.jpg', 1),
    (6, 'Watch', 'Analog wristwatch', 79.99, 'watch.jpg', 2),
    (7, 'Tie', 'Silk necktie', 24.99, 'tie.jpg', 2),
    (8, 'Sunglasses', 'UV protection sunglasses', 34.99, 'sunglasses.jpg', 3),
    (9, 'Hat', 'Stylish fedora hat', 19.99, 'hat.jpg', 2),
    (10, 'Wallet', 'Genuine leather wallet', 39.99, 'wallet.jpg', 2),
    (11, 'Shirt', 'Formal dress shirt', 29.99, 'shirt.jpg', 1),
    (12, 'Skirt', 'Floral print skirt', 34.99, 'skirt.jpg', 1),
    (13, 'Blouse', 'Chiffon blouse', 24.99, 'blouse.jpg', 1),
    (14, 'Pants', 'Casual trousers', 39.99, 'pants.jpg', 1),
    (15, 'Jacket', 'Lightweight bomber jacket', 49.99, 'jacket.jpg', 1),
    (16, 'Sandals', 'Comfortable sandals', 24.99, 'sandals.jpg', 3),
    (17, 'Earrings', 'Fashionable earrings', 14.99, 'earrings.jpg', 2),
    (18, 'Necklace', 'Statement necklace', 19.99, 'necklace.jpg', 2),
    (19, 'Bracelet', 'Stylish bracelet', 12.99, 'bracelet.jpg', 2),
    (20, 'Hat', 'Sun hat for summer', 29.99, 'sunhat.jpg', 2);


ALTER TABLE products
ADD COLUMN rating DECIMAL(3, 2) DEFAULT 0.0;


CREATE TABLE comments (
  id INT AUTO_INCREMENT PRIMARY KEY,
  product_id INT,
  user_id INT,
  comment_text TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (product_id) REFERENCES products (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
);


-- Comment 1 for Product 1 by User 1
INSERT INTO comments (product_id, user_id, comment_text) VALUES (1, 1, 'Great product!');

-- Comment 2 for Product 1 by User 2
INSERT INTO comments (product_id, user_id, comment_text) VALUES (1, 2, 'I love this product!');

-- Comment 3 for Product 2 by User 3
INSERT INTO comments (product_id, user_id, comment_text) VALUES (2, 3, 'Highly recommended.');

-- Comment 4 for Product 2 by User 4
INSERT INTO comments (product_id, user_id, comment_text) VALUES (2, 4, 'Good quality.');

-- Comment 5 for Product 3 by User 5
INSERT INTO comments (product_id, user_id, comment_text) VALUES (3, 5, 'Impressive features.');

-- Comment 6 for Product 3 by User 6
INSERT INTO comments (product_id, user_id, comment_text) VALUES (3, 6, 'Worth the price.');

-- Comment 7 for Product 4 by User 7
INSERT INTO comments (product_id, user_id, comment_text) VALUES (4, 7, 'Fast shipping.');

-- ... Repeat the above INSERT statements for the remaining comments ...


CREATE TABLE ratings ( 
    id INT PRIMARY KEY AUTO_INCREMENT, 
    product_id INT, 
    user_id INT, 
    rate INT, 
FOREIGN KEY (product_id) REFERENCES product (product_id), 
FOREIGN KEY (user_id) REFERENCES user (user_id) );
    
CREATE TRIGGER update_average_rating 
AFTER INSERT ON ratings 
FOR EACH ROW 
BEGIN 
UPDATE product 
SET ratings = 
( SELECT AVG(rate) 
FROM ratings 
WHERE product_id = NEW.product_id ) 
WHERE product_id = NEW.product_id;
 END;


 CREATE TABLE Cart (
    cart_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES User(user_id)
);

CREATE TABLE CartItem (
    cart_item_id INT PRIMARY KEY AUTO_INCREMENT,
    cart_id INT,
    product_id INT,
    FOREIGN KEY (cart_id) REFERENCES Cart(cart_id),
    FOREIGN KEY (product_id) REFERENCES Product(product_id)
);


DELIMITER //
CREATE TRIGGER after_cartitem_insert
AFTER INSERT ON cartitem
FOR EACH ROW
BEGIN
    UPDATE cart
    SET totalprice = (
        SELECT SUM(p.price)
        FROM cartitem ci
        JOIN product p ON ci.product_id = p.product_id
        WHERE ci.cart_id = NEW.cart_id
    )
    WHERE cart_id = NEW.cart_id;
END //
DELIMITER ;
