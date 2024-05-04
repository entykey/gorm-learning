use ntex::web::{self, HttpResponse};
use serde::Serialize;
use sqlx::Row;  // for row.get("field_name")

use crate::{models::models::{Order, Customer, OrderViewModel, OrderData, OrderItemViewModel, Product}, AppState};

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/orders")
            .route("simple", web::get().to(load_orders_simple))
            .route("simple-manual-mapping", web::get().to(load_orders_simple_manual_mapping))
            .route("simple-manual-mapping-1", web::get().to(load_orders_manual_1))
            .route("joined-props-non-optimized", web::get().to(load_orders_joined_props_3_queries))
            .route("joined_props_not_grouped", web::get().to(load_orders_joined_props_not_grouped))
            .route("joined_props_hashmap", web::get().to(load_orders_joined_props_hashmap))
            .route("joined-props-optimized", web::get().to(retrieve_orders_optimized)),
    );
}


/// Load Records with LRU cache
///
/// 2nd line of desc: Retrieve all Order
///
/// One could call the api endpoint with following curl.
/// ```text
/// curl http://192.168.1.11:4000/api/orders/simple
/// ```
/// autocannon test:
/// ```
/// autocannon http://192.168.1.11:4000/api/orders/simple -d 10 -c 300 -w 4
/// ```
// #[utoipa::path(
//     get,                       // important
//     path = "/api/orders/simple",        // important
//     responses(
//       (status = 200, description = "success response", body = [Vec<Order>]),
//     ),
// )]

async fn load_orders_simple(app_state: web::types::State<AppState>) -> HttpResponse {
    // Fetch Orders from the database with query_as! mapping feature
    // NOTE: `sqlx::query_as!()` requires running "$ cargo sqlx prepare --database-url <db-url>" to compile
    let orders_result: Result<Vec<Order>, sqlx::Error> =
        sqlx::query_as!(Order, "SELECT order_id, customer_id, CAST(order_date AS CHAR) AS order_date FROM orders")
            .fetch_all(&app_state.pool)
            .await;

    // Handle the result of the database query
    let response = match orders_result {

        Ok(orders_list) => {
            HttpResponse::Ok().json(&orders_list)
        }

        // Handle RowNotFound error
        Err(sqlx::Error::RowNotFound) => HttpResponse::NotFound().body("Orders not found"),

        // Catching other errors
        Err(e) => {
            eprintln!("Error fetching orders: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    };

    response
}

async fn load_orders_simple_manual_mapping(app_state: web::types::State<AppState>) -> HttpResponse {
    // Fetch orders from the database
    // NOTE: `sqlx::query()` DOES NOT requires running "$ cargo sqlx prepare --database-url <db-url>" to compile
    let query_result = sqlx::query(
        r#"
        SELECT order_id, customer_id, CAST(order_date AS CHAR) AS order_date_char FROM orders
        "#
    )
    .fetch_all(&app_state.pool)
    .await;

    // Handle the result of the database query
    let response = match query_result {
        Ok(rows) => {
            // Map the rows to Vec<Order>
            let orders_list: Vec<Order> = rows
                .into_iter()
                .map(|row| {
                    let order_id = row.get("order_id"); // must import "sqlx::Row"
                    let customer_id = row.get("customer_id"); // must import "sqlx::Row"
                    let order_date_char = row.get("order_date_char"); // must import "sqlx::Row"
                    
                    // Parse order_date_char into a DateTime type as needed
                    // For example, if it's a string in a specific format, you can use a parser like chrono
                    // let order_date = parse_order_date(&order_date_char);

                    Order {
                        order_id,
                        customer_id,
                        order_date: order_date_char, // Replace this with parsed order_date
                    }
                })
                .collect();

            HttpResponse::Ok().json(&orders_list)
        }
        Err(sqlx::Error::RowNotFound) => HttpResponse::NotFound().body("Orders not found"),
        Err(e) => {
            eprintln!("Error fetching orders: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    };

    response
}


// worked, but unwanted json format
async fn load_orders_joined_props_not_grouped(app_state: web::types::State<AppState>) -> HttpResponse {
    // the sqlx::query() return type `Result<Vec<MySqlRow>, Error>`
    // since `sqlx::mysql::MySqlRow` does not implment Serialize, so we create a wrapper struct for it:
    // Define a wrapper struct that implements Serialize for MySqlRow
    #[derive(Debug, Serialize)]
    struct SerializableMySqlRow {
        // Define fields to hold data from MySqlRow
        order_id: i64,
        order_date: String,
        customer_id: i64,
        customer__customer_id: Option<i64>,
        customer_name: Option<String>,
        customer_email: Option<String>,
        order_items__order_item_id: Option<i64>,
        order_items__product_id: Option<i64>,
        order_items__quantity: Option<i64>,
        order_items__product__product_id: Option<i64>,
        order_items__product__name: Option<String>,
        order_items__product__price: Option<f64>,
        // Add other fields as needed
    }

    impl From<sqlx::mysql::MySqlRow> for SerializableMySqlRow {
        fn from(row: sqlx::mysql::MySqlRow) -> Self {
            // Extract data from MySqlRow and initialize SerializableMySqlRow fields
            let order_id: i64 = row.get("order_id");
            let order_date: String = row.get("order_date");
            let customer_id = row.get("customer_id");
            let customer__customer_id = row.get("customer__customer_id");
            let customer_name = row.get("customer_name");
            let customer_email = row.get("customer_email");
            let order_items__order_item_id = row.get("order_items__order_item_id");
            let order_items__product_id = row.get("order_items__product_id");
            let order_items__quantity = row.get("order_items__quantity");
            let order_items__product__product_id = row.get("order_items__product__product_id");
            let order_items__product__name = row.get("order_items__product__name");
            let order_items__product__price = row.get("order_items__product__price");

            // Extract other fields as needed

            SerializableMySqlRow {
                order_id,
                order_date,
                customer_id,
                customer__customer_id,
                customer_name,
                customer_email,
                order_items__order_item_id,
                order_items__product_id,
                order_items__quantity,
                order_items__product__product_id,
                order_items__product__name,
                order_items__product__price,
                // Initialize other fields as needed
            }
        }
    }

    let orders_query = sqlx::query(
        r#"
        SELECT
            orders.order_id AS "order_id",
            -- orders.order_date AS "order_date",
            CAST(orders.order_date AS CHAR) AS "order_date", -- Convert DATETIME to CHAR
            orders.customer_id AS "customer_id",
            customers.customer_id AS "customer__customer_id",
            customers.name AS "customer_name",
            customers.email AS "customer_email",
            order_items.order_item_id AS "order_items__order_item_id",
            order_items.product_id AS "order_items__product_id",
            order_items.quantity AS "order_items__quantity",
            products.product_id AS "order_items__product__product_id",
            products.name AS "order_items__product__name",
            products.price AS "order_items__product__price"
        FROM orders
        LEFT JOIN customers ON orders.customer_id = customers.customer_id
        LEFT JOIN order_items ON orders.order_id = order_items.order_id
        LEFT JOIN products ON order_items.product_id = products.product_id
        "#
    );
    
    // Execute the query and fetch results
    let orders_result = orders_query.fetch_all(&app_state.pool).await;

    // Handle the result of the database query
    match orders_result {
        Ok(orders_list) => {
            let mut buf: Vec<u8> = Vec::new();
            let formatter: serde_json::ser::PrettyFormatter = serde_json::ser::PrettyFormatter::with_indent(b"  ");
            let mut ser: serde_json::Serializer<&mut Vec<u8>, serde_json::ser::PrettyFormatter> = serde_json::Serializer::with_formatter(&mut buf, formatter);
            let serialized_rows: Vec<SerializableMySqlRow> = orders_list.into_iter().map(|row| row.into()).collect();
            serialized_rows.serialize(&mut ser).unwrap();   // the method `serialize` exists for struct `Vec<MySqlRow>`
            // let buf_clone: Vec<u8> = buf.clone();
            println!("{}", String::from_utf8(buf).unwrap());
            // HttpResponse::Ok().body(buf_clone)

            // if use `sqlx::query()` => the trait bound `MySqlRow: Serialize` is not satisfied
            // let serialized_rows: Vec<SerializableMySqlRow> = orders_list.into_iter().map(|row| row.into()).collect();
            
            HttpResponse::Ok().json(&serialized_rows)
        },
        Err(sqlx::Error::RowNotFound) => HttpResponse::NotFound().body("Orders not found"),
        Err(e) => {
            eprintln!("Error fetching orders: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
    
}

async fn load_orders_joined_props_hashmap(app_state: web::types::State<AppState>) -> HttpResponse {
    // the sqlx::query() return type `Result<Vec<MySqlRow>, Error>`
    // since `sqlx::mysql::MySqlRow` does not implment Serialize, so we create a wrapper struct for it:
    // Define a wrapper struct that implements Serialize for MySqlRow
    #[derive(Debug, Serialize)]
    struct SerializableMySqlRow {
        // Define fields to hold data from MySqlRow
        order_id: i64,
        order_date: String,
        customer_id: i64,
        customer__customer_id: Option<i64>,
        customer_name: Option<String>,
        customer_email: Option<String>,
        order_items__order_item_id: Option<i64>,
        order_items__product_id: Option<i64>,
        order_items__quantity: Option<i64>,
        order_items__product__product_id: Option<i64>,
        order_items__product__name: Option<String>,
        order_items__product__price: Option<f64>,
        // Add other fields as needed
    }

    impl From<sqlx::mysql::MySqlRow> for SerializableMySqlRow {
        fn from(row: sqlx::mysql::MySqlRow) -> Self {
            // Extract data from MySqlRow and initialize SerializableMySqlRow fields
            let order_id: i64 = row.get("order_id");
            let order_date: String = row.get("order_date");
            let customer_id = row.get("customer_id");
            let customer__customer_id = row.get("customer__customer_id");
            let customer_name = row.get("customer_name");
            let customer_email = row.get("customer_email");
            let order_items__order_item_id = row.get("order_items__order_item_id");
            let order_items__product_id = row.get("order_items__product_id");
            let order_items__quantity = row.get("order_items__quantity");
            let order_items__product__product_id = row.get("order_items__product__product_id");
            let order_items__product__name = row.get("order_items__product__name");
            let order_items__product__price = row.get("order_items__product__price");

            // Extract other fields as needed

            SerializableMySqlRow {
                order_id,
                order_date,
                customer_id,
                customer__customer_id,
                customer_name,
                customer_email,
                order_items__order_item_id,
                order_items__product_id,
                order_items__quantity,
                order_items__product__product_id,
                order_items__product__name,
                order_items__product__price,
                // Initialize other fields as needed
            }
        }
    }

    let orders_query = sqlx::query(
        r#"
        SELECT
            orders.order_id AS "order_id",
            -- orders.order_date AS "order_date",
            CAST(orders.order_date AS CHAR) AS "order_date", -- Convert DATETIME to CHAR
            orders.customer_id AS "customer_id",
            customers.customer_id AS "customer__customer_id",
            customers.name AS "customer_name",
            customers.email AS "customer_email",
            order_items.order_item_id AS "order_items__order_item_id",
            order_items.product_id AS "order_items__product_id",
            order_items.quantity AS "order_items__quantity",
            products.product_id AS "order_items__product__product_id",
            products.name AS "order_items__product__name",
            products.price AS "order_items__product__price"
        FROM orders
        LEFT JOIN customers ON orders.customer_id = customers.customer_id
        LEFT JOIN order_items ON orders.order_id = order_items.order_id
        LEFT JOIN products ON order_items.product_id = products.product_id
        "#
    );
    
    // Execute the query and fetch results
    let orders_result = orders_query.fetch_all(&app_state.pool).await;

    // Handle the result of the database query
    match orders_result {
        Ok(orders_list) => {

            // Create a HashMap to group order items by order ID
            let mut orders_map = std::collections::HashMap::new();

            // Iterate over fetched rows and populate the HashMap
            for row in orders_list {
                // Convert the row into a SerializableMySqlRow
                let order_row: SerializableMySqlRow = row.into();

                // Get the order ID
                let order_id = order_row.order_id;

                // Check if the order ID already exists in the HashMap
                let order = orders_map.entry(order_id).or_insert_with(|| {
                    // If the order ID doesn't exist, create a new OrderViewModel
                    OrderViewModel {
                        order_id,
                        order_date: order_row.order_date,
                        customer_id: order_row.customer_id,
                        customer: Customer {
                            customer_id: order_row.customer__customer_id.unwrap_or_default(),
                            name: order_row.customer_name.clone().unwrap_or_default(),
                            email: order_row.customer_email.clone().unwrap_or_default(),
                        },
                        order_items: Vec::new(), // Initialize an empty Vec for order items
                        total: 0.0, // Initialize total for the order
                    }
                });

                // Create an OrderItemViewModel for the current row
                let order_item = OrderItemViewModel {
                    order_item_id: order_row.order_items__order_item_id.unwrap_or_default(),
                    product_id: order_row.order_items__product_id.unwrap_or_default(),
                    quantity: order_row.order_items__quantity.unwrap_or_default(),
                    product: Product {
                        product_id: order_row.order_items__product__product_id.unwrap_or_default(),
                        name: order_row.order_items__product__name.clone().unwrap_or_default(),
                        price: order_row.order_items__product__price.unwrap_or_default(),
                    },
                };

                // Push the order item into the order's order_items Vec
                order.order_items.push(order_item.clone());
            }

            // Calculate the total for each order after all order items are added
            for order in orders_map.values_mut() {
                for order_item in &order.order_items {
                    order.total += order_item.quantity as f64 * order_item.product.price;
                }
            }

            // Convert the HashMap values into a Vec<OrderViewModel>
            let orders: Vec<OrderViewModel> = orders_map.into_values().collect();

            let mut buf: Vec<u8> = Vec::new();
            let formatter: serde_json::ser::PrettyFormatter = serde_json::ser::PrettyFormatter::with_indent(b"  ");
            let mut ser: serde_json::Serializer<&mut Vec<u8>, serde_json::ser::PrettyFormatter> = serde_json::Serializer::with_formatter(&mut buf, formatter);
            orders.serialize(&mut ser).unwrap();   // the method `serialize` exists for struct `Vec<MySqlRow>`
            // let buf_clone: Vec<u8> = buf.clone();
            println!("{}", String::from_utf8(buf).unwrap());
            
            HttpResponse::Ok().json(&orders)
        },
        Err(sqlx::Error::RowNotFound) => HttpResponse::NotFound().body("Orders not found"),
        Err(e) => {
            eprintln!("Error fetching orders: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
    
}

async fn load_orders_manual_1(app_state: web::types::State<AppState>) -> HttpResponse {
    //*
    // the sqlx::query() return type `Result<Vec<MySqlRow>, Error>`
    // since `sqlx::mysql::MySqlRow` does not implment Serialize, so we create a wrapper struct for it:
    // Define a wrapper struct that implements Serialize for MySqlRow
    #[derive(Debug, Serialize)]
    struct SerializableMySqlRow {
        // Define fields to hold data from MySqlRow
        order_id: i64,
        order_date: String,
        customer_id: i64,
        // customer__customer_id: Option<i64>,
        customer_name: Option<String>,
        customer_email: Option<String>,
        // order_items__order_item_id: Option<i64>,
        // order_items__product_id: Option<i64>,
        // order_items__quantity: Option<i64>,
        // order_items__product__product_id: Option<i64>,
        // order_items__product__name: Option<String>,
        // order_items__product__price: Option<f64>,
        // Add other fields as needed
    }

    impl From<sqlx::mysql::MySqlRow> for SerializableMySqlRow {
        fn from(row: sqlx::mysql::MySqlRow) -> Self {
            // Extract data from MySqlRow and initialize SerializableMySqlRow fields
            let order_id: i64 = row.get("order_id");
            let order_date: String = row.get("order_date");
            let customer_id: i64 = row.get("customer_id");
            let customer_name: Option<String> = row.get("customer_name");
            let customer_email: Option<String> = row.get("customer_email");

            // Extract other fields as needed

            SerializableMySqlRow {
                order_id,
                order_date,
                customer_id,
                customer_name,
                customer_email
                // Initialize other fields as needed
            }
        }
    }

    let orders_result = sqlx::query(
        r#"
        SELECT
            o.order_id,
            CAST(o.order_date AS CHAR) AS "order_date", 
            o.customer_id,
            c.name AS customer_name,
            c.email AS customer_email,
            oi.order_item_id,
            oi.product_id,
            oi.quantity,
            p.name AS product_name,
            p.price AS product_price
        FROM
            orders o
            LEFT JOIN customers c ON o.customer_id = c.customer_id
            LEFT JOIN order_items oi ON o.order_id = oi.order_id
            LEFT JOIN products p ON oi.product_id = p.product_id
        ORDER BY
            o.order_id, oi.order_item_id
        "#,
    )
    .fetch_all(&app_state.pool)
    .await;


    // Handle the result of the database query
    match orders_result {
        
        Ok(rows) => {
            // Convert MySqlRow instances to SerializableMySqlRow
            let serialized_rows: Vec<SerializableMySqlRow> = rows.into_iter().map(|row| row.into()).collect();
            HttpResponse::Ok().json(&serialized_rows)
        },
        Err(sqlx::Error::RowNotFound) => HttpResponse::NotFound().body("Orders not found"),
        Err(e) => {
            eprintln!("Error fetching orders: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}

// attempt to fix from above working code:
async fn load_orders_joined_props_3_queries(app_state: web::types::State<AppState>) -> HttpResponse {
    #[derive(Debug, Serialize)]
    struct Order {
        order_id: i64,
        order_date: String,
        customer: Customer,
        order_items: Vec<OrderItem>,
        total: f64,
    }

    // Add a method to calculate the total for the order
    impl Order {
        fn calculate_total(&self) -> f64 {
            self.order_items.iter().map(|item| item.quantity as f64 * item.product.price).sum()
        }
    }

    #[derive(Debug, Serialize)]
    struct Customer {
        customer_id: i64,
        name: String,
        email: String,
    }

    #[derive(Debug, Serialize)]
    struct OrderItem {
        order_item_id: i64,
        product: Product,
        quantity: i64,
    }

    #[derive(Debug, Serialize)]
    struct Product {
        product_id: i64,
        name: String,
        price: f64,
    }

    let orders_query = sqlx::query!(
        r#"
        SELECT 
            orders.order_id AS "order_id!",
            CAST(orders.order_date AS CHAR) AS "order_date!",
            orders.customer_id AS "customer_id!"
        FROM orders
        "#
    );

    let order_query_time = std::time::Instant::now();
    let orders_result = orders_query.fetch_all(&app_state.pool).await;
    
    // stop timer & print to terminal
    let duration: std::time::Duration = order_query_time.elapsed();
    let elapsed_ms: f64 = duration.as_secs_f64() * 1000.0;
    let elapsed_seconds: f64 = elapsed_ms / 1000.0;
    print!(
        "\norder_query time: {:?} ({:?} ms) ({:.8} s)\n",
        duration, elapsed_ms, elapsed_seconds
    );

    match orders_result {
        Ok(orders_list) => {
            let mut response: Vec<Order> = Vec::new();

            for order_row in orders_list {
                let customer_id = order_row.customer_id;
                let order_id = order_row.order_id;

                let customer_query = sqlx::query!(
                    r#"
                    SELECT 
                        customers.customer_id AS "customer_id!",
                        customers.name AS "name!",
                        customers.email AS "email!"
                    FROM customers
                    WHERE customers.customer_id = ?
                    "#,
                    customer_id
                );


                let customer_query_time = std::time::Instant::now();
                let customer_result = customer_query.fetch_one(&app_state.pool).await;

                // stop timer & print to terminal
                let duration: std::time::Duration = customer_query_time.elapsed();
                let elapsed_ms: f64 = duration.as_secs_f64() * 1000.0;
                let elapsed_seconds: f64 = elapsed_ms / 1000.0;
                println!(
                    "customer_query time: {:?} ({:?} ms) ({:.8} s)",
                    duration, elapsed_ms, elapsed_seconds
                );

                match customer_result {
                    Ok(customer_row) => {
                        let customer = Customer {
                            customer_id: customer_row.customer_id,
                            name: customer_row.name,
                            email: customer_row.email,
                        };

                        let order_items_query = sqlx::query!(
                            r#"
                            SELECT 
                                order_items.order_item_id AS "order_item_id!",
                                order_items.product_id AS "product_id!",
                                order_items.quantity AS "quantity!",
                                products.product_id AS "product__product_id!",
                                products.name AS "product__name!",
                                products.price AS "product__price!"
                            FROM order_items
                            LEFT JOIN products ON order_items.product_id = products.product_id
                            WHERE order_items.order_id = ?
                            "#,
                            order_id
                        );
                        
                        let order_item_query_time = std::time::Instant::now();
                        let order_items_result = order_items_query.fetch_all(&app_state.pool).await;
                        
                        // stop timer & print to terminal
                        let duration: std::time::Duration = order_item_query_time.elapsed();
                        let elapsed_ms: f64 = duration.as_secs_f64() * 1000.0;
                        let elapsed_seconds: f64 = elapsed_ms / 1000.0;
                        println!(
                            "order_item_query time: {:?} ({:?} ms) ({:.8} s)",
                            duration, elapsed_ms, elapsed_seconds
                        );

                        match order_items_result {
                            Ok(order_items_list) => {
                                let mut order_items: Vec<OrderItem> = Vec::new();

                                for order_item_row in order_items_list {
                                    let product = Product {
                                        product_id: order_item_row.product_id,
                                        name: order_item_row.product__name,
                                        price: order_item_row.product__price,
                                    };

                                    let order_item: OrderItem = OrderItem {
                                        order_item_id: order_item_row.order_item_id,
                                        product,
                                        quantity: order_item_row.quantity,
                                    };

                                    order_items.push(order_item);
                                }

                                // Calculate total
                                let total = order_items.iter().map(|item| item.quantity as f64 * item.product.price).sum();

                                let order = Order {
                                    order_id,
                                    order_date: order_row.order_date,
                                    customer,
                                    order_items,
                                    total,
                                };

                                response.push(order);
                            }
                            Err(e) => {
                                eprintln!("Error fetching order items: {:?}", e);
                                return HttpResponse::InternalServerError().body("Internal Server Error");
                            }
                        }
                    }
                    Err(e) => {
                        eprintln!("Error fetching customer: {:?}", e);
                        return HttpResponse::InternalServerError().body("Internal Server Error");
                    }
                }
            }

            HttpResponse::Ok().json(&response)
        }
        Err(e) => {
            eprintln!("Error fetching orders: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}

async fn retrieve_orders_optimized(app_state: web::types::State<AppState>) -> HttpResponse {
    // removed o.order_date, due to NaiveDateTime failure
    let result = sqlx::query!(
        r#"
        SELECT
            o.order_id,
            CAST(o.order_date AS CHAR) AS "order_date", 
            o.customer_id,
            c.name AS customer_name,
            c.email AS customer_email,
            oi.order_item_id,
            oi.product_id,
            oi.quantity,
            p.name AS product_name,
            p.price AS product_price
        FROM
            orders o
            LEFT JOIN customers c ON o.customer_id = c.customer_id
            LEFT JOIN order_items oi ON o.order_id = oi.order_id
            LEFT JOIN products p ON oi.product_id = p.product_id
        ORDER BY
            o.order_id, oi.order_item_id
        "#,
    )
    .fetch_all(&app_state.pool)
    .await;

    match result {
        Ok(rows) => {
            let mut orders: Vec<OrderViewModel> = Vec::new();
            let mut current_order_id: Option<i64> = None;
            let mut current_order: Option<OrderViewModel> = None;

            for row in rows {
                let data = OrderData {
                    order_id: row.order_id,
                    // order_date: row.order_date, // err: mismatched types: expected struct `std::string::String` found enum `std::option::Option<std::string::String>`r
                    order_date: row.order_date.unwrap(), // without sql CAST => error: "no method named `unwrap` found for struct `PrimitiveDateTime` in the current scope. method not found in `PrimitiveDateTime`r"
                    // order_date: row.order_date.format("%Y-%m-%dT%H:%M:%S").to_string(),  // without sql CAST => err

                    // order_date: match row.order_date.format("%Y-%m-%dT%H:%M:%S") {
                    //     Ok(formatted_date) => formatted_date.to_string(),
                    //     Err(_) => String::new(), // Provide a default value in case of an error
                    // },

                    // order_date: match row.order_date {
                    //     Some(dt) => format!("%Y-%m-%dT%H:%M:%S").to_string(),
                    //     None => "".to_string(),
                    // },
                    
                    customer_id: row.customer_id,
                    customer_name: row.customer_name.unwrap(),
                    customer_email: row.customer_email.unwrap(),
                    order_item_id: row.order_item_id.unwrap(),
                    product_id: row.product_id.unwrap(),
                    quantity: row.quantity.unwrap(),
                    product_name: row.product_name.unwrap(),
                    product_price: row.product_price.unwrap(),
                };

                if current_order_id != Some(data.order_id) {
                    if let Some(order) = current_order.take() {
                        orders.push(order);
                    }
                    current_order_id = Some(data.order_id);
                    let customer = Customer {
                        customer_id: data.customer_id,
                        name: data.customer_name,
                        email: data.customer_email,
                    };
                    let order_items: Vec<OrderItemViewModel> = Vec::new();
                    let total: f64 = 0.0; // Initialize total for the order
                    let order = OrderViewModel {
                        order_id: data.order_id,
                        order_date: data.order_date,
                        customer_id: data.customer_id,
                        customer,
                        order_items,
                        total,
                    };
                    current_order = Some(order);
                }

                if let Some(ref mut order) = current_order {
                    let product = Product {
                        product_id: data.product_id,
                        name: data.product_name,
                        price: data.product_price,
                    };
                    let order_item = OrderItemViewModel {
                        order_item_id: data.order_item_id,
                        product_id: data.product_id,
                        quantity: data.quantity,
                        product,
                    };
                    order.order_items.push(order_item);

                    // Update total for the order
                    order.total += data.quantity as f64 * data.product_price;
                }
            }

            if let Some(order) = current_order.take() {
                orders.push(order);
            }

            HttpResponse::Ok().json(&orders)
        }
        Err(e) => {
            eprintln!("Error retrieving orders: {:?}", e);
            HttpResponse::InternalServerError().finish()
        }
    }
}

/*
// failed
async fn load_orders_joined_props_v1(app_state: web::types::State<AppState>) -> HttpResponse {
    // let orders_query = sqlx::query_as::<_, JoinedOrderResponse>(
    let orders_query = sqlx::query_as!(JoinedOrderResponse,
        r#"
        SELECT
            orders.order_id AS "order_id!",
            CAST(orders.order_date AS CHAR) AS "order_date!", -- Convert DATETIME to CHAR
            customers.customer_id AS "customer.customer_id!",
            customers.name AS "customer.name!",
            customers.email AS "customer.email!",
            order_items.order_item_id AS "order_items__order_item_id!",
            order_items.product_id AS "order_items__product_id!",
            order_items.quantity AS "order_items__quantity!",
            products.product_id AS "order_items__product.product_id!",
            products.name AS "order_items__product.name!",
            products.price AS "order_items__product.price!"
        FROM orders
        LEFT JOIN customers ON orders.customer_id = customers.customer_id
        LEFT JOIN order_items ON orders.order_id = order_items.order_id
        LEFT JOIN products ON order_items.product_id = products.product_id
        "#
    );

    // Execute the query and fetch results
    let orders_result = orders_query.fetch_all(&app_state.pool).await;

    // Handle the result of the database query
    match orders_result {
        Ok(orders_list) => HttpResponse::Ok().json(&orders_list),
        Err(sqlx::Error::RowNotFound) => HttpResponse::NotFound().body("Orders not found"),
        Err(e) => {
            eprintln!("Error fetching orders: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}
*/
