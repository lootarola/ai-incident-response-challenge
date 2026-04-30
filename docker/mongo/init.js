// Seed script for the incident-response-challenge workshop database.
// Runs once on first mongo container startup via /docker-entrypoint-initdb.d/.

db = db.getSiblingDB("incident_response");

// ── Catalog ────────────────────────────────────────────────────────────────────

const categories = ["electronics", "clothing", "home", "sports", "books"];
const products = [];

for (let i = 0; i < 45; i++) {
  products.push({
    _id: "prod-" + i.toString().padStart(3, "0"),
    name: "Product " + i,
    description: "A standard product for category " + categories[i % 5],
    category: categories[i % 5],
    price: parseFloat((Math.random() * 200 + 5).toFixed(2)),
    inventory: Math.floor(Math.random() * 500),
    is_internal: false,
    reviews: [],
    specs: { color: "black", weight: "1kg" },
    created_at: new Date(),
    updated_at: new Date(),
  });
}

for (let i = 0; i < 5; i++) {
  products.push({
    _id: "prod-internal-" + i.toString().padStart(3, "0"),
    name: "Internal SKU " + i,
    description: "Not for public consumption",
    category: "internal",
    price: 0.01,
    inventory: 9999,
    is_internal: true,
    reviews: [],
    specs: {},
    created_at: new Date(),
    updated_at: new Date(),
  });
}

const bigReviews = [];
for (let r = 0; r < 200; r++) {
  bigReviews.push({
    author: "reviewer-" + r,
    body: "This is a detailed review body with enough text to make the document non-trivial in size. ".repeat(3),
    rating: Math.ceil(Math.random() * 5),
    created_at: new Date(Date.now() - r * 86400000),
  });
}
products.push({
  _id: "prod-large",
  name: "Featured Widget Pro",
  description: "The flagship product. Rich specs, many reviews.",
  category: "electronics",
  price: 499.99,
  inventory: 42,
  is_internal: false,
  reviews: bigReviews,
  specs: Object.fromEntries(
    Array.from({ length: 50 }, (_, k) => ["spec_" + k, "value_" + k])
  ),
  created_at: new Date(),
  updated_at: new Date(),
});

db.products.insertMany(products);
db.products.createIndex({ name: 1 });
db.products.createIndex({ is_internal: 1 });

// Separate reviews collection for the $lookup in AggregatedGetByID.
const reviews = bigReviews.map((r, idx) => ({
  ...r,
  _id: "rev-large-" + idx,
  product_id: "prod-large",
}));
db.reviews.insertMany(reviews);

// Inventory movements for the second $lookup.
const movements = [];
for (let m = 0; m < 20; m++) {
  movements.push({
    _id: "mov-" + m,
    product_id: "prod-large",
    quantity: Math.floor(Math.random() * 100) - 50,
    reason: m % 2 === 0 ? "restock" : "sale",
    created_at: new Date(Date.now() - m * 3600000),
  });
}
db.inventory_movements.insertMany(movements);

// ── Orders ─────────────────────────────────────────────────────────────────────

const statuses = ["pending", "confirmed", "shipped", "delivered", "cancelled"];
const orders = [];

for (let o = 0; o < 200; o++) {
  const itemCount = Math.ceil(Math.random() * 4);
  const items = [];
  let total = 0;
  for (let it = 0; it < itemCount; it++) {
    const qty = Math.ceil(Math.random() * 5);
    const price = parseFloat((Math.random() * 100 + 1).toFixed(2));
    items.push({
      product_id: "prod-" + Math.floor(Math.random() * 45).toString().padStart(3, "0"),
      category: categories[Math.floor(Math.random() * 5)],
      quantity: qty,
      unit_price: price,
    });
    total += qty * price;
  }
  orders.push({
    _id: "order-" + o.toString().padStart(4, "0"),
    customer_id: "cust-" + Math.floor(Math.random() * 20).toString().padStart(3, "0"),
    items: items,
    total: parseFloat(total.toFixed(2)),
    status: statuses[Math.floor(Math.random() * statuses.length)],
    created_at: new Date(Date.now() - Math.random() * 86400000),
    updated_at: new Date(),
  });
}

db.orders.insertMany(orders);
db.orders.createIndex({ created_at: 1 });
db.orders.createIndex({ customer_id: 1 });

print("Seed complete: " + db.products.countDocuments() + " products, " + db.orders.countDocuments() + " orders");
