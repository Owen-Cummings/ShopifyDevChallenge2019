/* Product table defined to store numerical IDs, titles limited to 50 characters, a price
up to 2 decimals and inventory count constrained to be non-negative */

CREATE TABLE PRODUCTS (
  PRODUCT_ID INT PRIMARY KEY,
  TITLE VARCHAR(50) NOT NULL,
  PRICE DECIMAL(10,2) NOT NULL,
  INVENTORY_COUNT INT NOT NULL,
  CONSTRAINT POSITIVE_COUNT CHECK (INVENTORY_COUNT >= 0)
);