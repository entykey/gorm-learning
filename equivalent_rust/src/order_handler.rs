use ntex::web::{self, HttpResponse};
use serde::{Deserialize, Serialize};

use crate::AppState;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/orders")
            // .route("simple", web::get().to(load_orders))
            .route("joined-props", web::get().to(load_orders_joined_props_3_queries))
            // .route("joined-props-optimize", web::get().to(load_orders_joined_props_1_query)),
    );
}

// async fn load_roles(pool: web::Data<MySqlPool>) -> Result<HttpResponse> {
//     let roles_result = sqlx::query("SELECT * FROM AspNetRoles")
//         .fetch_all(pool.get_ref())
//         .await;

//     match roles_result {
//         Ok(roles_list) => {
//             Ok(HttpResponse::Ok().json(roles_list))
//         }

//         // Handle RowNotFound error
//         Err(sqlx::Error::RowNotFound) => {
//             Ok(HttpResponse::NotFound().body("Roles not found"))
//         }

//         // Catching other errors
//         Err(e) => {
//             eprintln!("Error fetching roles: {:?}", e);
//             Ok(HttpResponse::InternalServerError().body("Internal Server Error"))
//         }
//     }
// }

/// Load Roles with LRU cache
///
/// 2nd line of desc: Retrieve all AspNetRoles
///
/// One could call the api endpoint with following curl.
/// ```text
/// curl http://192.168.1.11:4000/api/roles
/// ```
/// autocannon test:
/// ```
/// autocannon http://192.168.1.11:4000/api/roles -d 10 -c 300 -w 4
/// ```
// #[utoipa::path(
//     get,                       // important
//     path = "/api/roles",        // important
//     responses(
//       (status = 200, description = "success response", body = [Vec<AspNetRole>]),
//     ),
// )]
/*
async fn load_orders(app_state: web::types::State<AppState>) -> HttpResponse {
    // If roles are not in the cache, fetch them from the database
    let orders_result: Result<Vec<Order>, sqlx::Error> =
        sqlx::query_as!(Order, "SELECT order_id, customer_id FROM orders")
            .fetch_all(&app_state.pool)
            .await;

    // Handle the result of the database query
    let response = match orders_result {
        Ok(orders_list) => {
            // Update the cache with the fetched roles
            // let mut cache = app_state.role_cache.lock().unwrap();
            // cache.put(cache_key.to_string(), orders_list.clone());

            // println!("Roles put to cache");

            HttpResponse::Ok().json(&orders_list)
        }
        Err(sqlx::Error::RowNotFound) => HttpResponse::NotFound().body("Roles not found"),
        Err(e) => {
            eprintln!("Error fetching roles: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    };

    response
}
*/

// worked, but unwanted json format
/*
async fn load_orders_joined_props(app_state: web::types::State<AppState>) -> HttpResponse {
    // Define a new struct to represent the joined data from the query
    #[derive(Debug, Serialize)]
    struct JoinedOrder {
        order_id: i64,
        order_date: String,
        customer_id: i64,
        customer__customer_id: Option<i64>,
        customer__name: Option<String>,
        customer__email: Option<String>,
        order_items__order_item_id: Option<i64>,
        order_items__product_id: Option<i64>,
        order_items__quantity: Option<i64>,
        order_items__product__product_id: Option<i64>,
        order_items__product__name: Option<String>,
        order_items__product__price: Option<f64>,
    }

    // Fetch orders with joined properties from the database
    let orders_query = sqlx::query_as!(
        JoinedOrder, // Specify the type of the result (should be a struct that implements Serialize)
        r#"
        SELECT
            orders.order_id AS "order_id!",
            -- orders.order_date AS "order_date!",
            CAST(orders.order_date AS CHAR) AS "order_date!", -- Convert DATETIME to CHAR
            orders.customer_id AS "customer_id!",
            customers.customer_id AS "customer__customer_id!",
            customers.name AS "customer__name!",
            customers.email AS "customer__email!",
            order_items.order_item_id AS "order_items__order_item_id!",
            order_items.product_id AS "order_items__product_id!",
            order_items.quantity AS "order_items__quantity!",
            products.product_id AS "order_items__product__product_id!",
            products.name AS "order_items__product__name!",
            products.price AS "order_items__product__price!"
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

    let orders_result = orders_query.fetch_all(&app_state.pool).await;

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

                let customer_result = customer_query.fetch_one(&app_state.pool).await;

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

                        let order_items_result = order_items_query.fetch_all(&app_state.pool).await;

                        match order_items_result {
                            Ok(order_items_list) => {
                                let mut order_items: Vec<OrderItem> = Vec::new();

                                for order_item_row in order_items_list {
                                    let product = Product {
                                        product_id: order_item_row.product_id,
                                        name: order_item_row.product__name,
                                        price: order_item_row.product__price,
                                    };

                                    let order_item = OrderItem {
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

/*

// async fn load_roles(
//     app_state: web::types::State<AppState>, // ntex
// ) -> HttpResponse {
//     let roles_result: Result<Vec<AspNetRole>, sqlx::Error> = sqlx::query_as::<_, AspNetRole>("SELECT * FROM AspNetRoles")
//         .fetch_all(&app_state.pool)
//         .await;

//     // println!("{:?}", roles_result); // Ok([AspNetRole { Id: "094cd9d9-7388-4804-89d9-f748148bce17", Name: "admin", NormalizedName: "ADMIN" }, AspNetRole { Id: "1489bac2-c238-49d0-92c0-6f0049301157", Name: "staff", NormalizedName: "STAFF" }])

//     match roles_result {
//         Ok(roles_list) => {
//             HttpResponse::Ok().json(&roles_list)
//         }
//         // Handle RowNotFound error
//         Err(sqlx::Error::RowNotFound) => {
//             HttpResponse::NotFound().body("Roles not found")
//         }
//         // Catching other errors
//         Err(e) => {
//             eprintln!("Error fetching roles: {:?}", e);
//             HttpResponse::InternalServerError().body("Internal Server Error")
//         }
//     }
// }

/// Load Roles By UserName
///
/// 2nd line of desc: Retrieve all AspNetRoles of an AspNetUser
///
/// One could call the api endpoint with following curl.
/// ```text
/// curl localhost:4000/api/roles/{username}
/// ```
#[utoipa::path(
    get,                       // important
    path = "/api/roles/{username}",        // important
    params(
        ("username", description = "Type in UserName", example = "testuser")
    ),
    responses(
      (status = 200, description = "success response",body = UserInRoleInfo, content_type = "application/json", example = json!(
        [
            {
            "RoleId": "xxx",
            "RoleName": "xxx",
            "IsInRole": false
            },
            {
            "RoleId": "xxx",
            "RoleName": "xxx",
            "IsInRole": false
            },
      ])),
    ),
)]
async fn load_roles_by_username(
    app_state: web::types::State<AppState>, // ntex
    path: web::types::Path<String>,
) -> HttpResponse {

    let username: &String = &*path;
    let result: Result<Vec<UserInRoleInfo>, sqlx::Error> = sqlx::query_as::<_, UserInRoleInfo>(
        r#"
        SELECT
            r.Id as RoleId,
            r.Name as RoleName,
            COALESCE(ur.UserId, 0) as IsInRole
        FROM AspNetRoles r
        LEFT JOIN (
            SELECT RoleId, 1 as UserId
            FROM AspNetUserRoles
            WHERE UserId = (SELECT Id FROM AspNetUsers WHERE UserName = ?)
        ) ur ON r.Id = ur.RoleId
        "#,
    )
    .bind(username)
    .fetch_all(&app_state.pool)
    .await;

    match result {
        Ok(user_roles) => HttpResponse::Ok().json(&user_roles),
        // Err(_) => web::HttpResponse::InternalServerError().finish()

        // Catching other errors
        Err(e) => {
            eprintln!("Error fetching roles of designated user: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}

// async fn update_user_roles(
//     app_state: web::types::State<AppState>,
//     data: web::types::Json<Vec<UserInRoleInfo>>,
//     path: web::types::Path<String>,
// ) -> HttpResponse {
//     use sqlx::Acquire;

//     let username: &String = &*path;
//     let mut conn = match app_state.pool.acquire().await {
//         Ok(conn) => conn,
//         Err(e) => {
//             eprintln!("Failed to acquire database connection: {:?}", e);
//             return HttpResponse::InternalServerError().body("Internal Server Error, Failed to acquire database connection");
//         }
//     };

//     let mut transaction = match conn.begin().await {
//         Ok(transaction) => transaction,
//         Err(e) => {
//             eprintln!("Failed to begin database transaction: {:?}", e);
//             return HttpResponse::InternalServerError().body("Internal Server Error, Failed to begin database transaction");
//         }
//     };

//     // Clear all existing roles of targeted user
//     if let Err(e) = sqlx::query("DELETE FROM AspNetUserRoles WHERE UserId = (SELECT Id FROM AspNetUsers WHERE UserName = ?)")
//         .bind(username)
//         .execute(&mut *transaction)
//         .await {
//         eprintln!("Failed to execute DELETE query: {:?}", e);
//         return HttpResponse::InternalServerError().body("Internal Server Error, Failed to execute DELETE query");
//     }

//     // Add selected roles for the user (all roles in the request body will be assigned to user)
//     for user_role in &data.into_inner() {
//         if let Err(e) = sqlx::query(
//             "INSERT INTO AspNetUserRoles (UserId, RoleId) VALUES ((SELECT Id FROM AspNetUsers WHERE UserName = ?), ?)",
//         )
//         .bind(username.clone())
//         .bind(&user_role.RoleId)
//         .execute(&mut transaction)
//         .await {
//             eprintln!("Failed to execute INSERT query: {:?}", e);
//             return HttpResponse::InternalServerError().body("Internal Server Error, Failed to execute INSERT query");
//         }
//     }

//     if let Err(e) = transaction.commit().await {
//         eprintln!("Failed to commit transaction: {:?}", e);
//         return HttpResponse::InternalServerError().body("Internal Server Error, Failed to commit transaction");
//     }

//     // Respond with JSON body
//     HttpResponse::Ok().json(&serde_json::json!({
//         "message": "Updated the user roles.",
//     }))
// }
// // Example request body:
// // 1. Add this 1 role to the targeted user (username in path param)
// // [
// //     {
// //         "RoleId": "1489bac2-c238-49d0-92c0-6f0049301157",
// //         "RoleName": "staff",
// //         "IsInRole": true
// //     }
// // ]

// // 2. Add 2 roles to the targeted user (username in path param) (doesn't care about the IsInRole is true or false)
// // [
// //     {
// //         "RoleId": "094cd9d9-7388-4804-89d9-f748148bce17",
// //         "RoleName": "admin",
// //         "IsInRole": false
// //     },
// //     {
// //         "RoleId": "1489bac2-c238-49d0-92c0-6f0049301157",
// //         "RoleName": "staff",
// //         "IsInRole": true
// //     }
// // ]


/// Assign User Roles
///
/// 2nd line of desc: Update / Assign User Roles
///
/// One could call the api endpoint with following curl.
/// ```text
/// curl localhost:4000/api/roles/{username}
/// ```
#[utoipa::path(
    put,                                // important
    request_body = UserInRoleInfo,
    request_body(content = UserInRoleInfo, description = "Json Vec<UserInRoleInfo> request body", content_type = "application/json", example = json!(
        [
            {
                "RoleId": "094cd9d9-7388-4804-89d9-f748148bce17",
                "RoleName": "admin",
                "IsInRole": false
              },
              {
                "RoleId": "1489bac2-c238-49d0-92c0-6f0049301157",
                "RoleName": "staff",
                "IsInRole": true
              }
      ]
    )),
    path = "/api/roles/{username}",     // important
    params(
        ("username", description = "Type in UserName", example = "testuser")
    ),
    responses(
      (status = 200, description = "success response"),
    ),
)]
async fn update_user_roles(
    app_state: web::types::State<AppState>,
    data: web::types::Json<Vec<UserInRoleInfo>>,
    path: web::types::Path<String>,
) -> HttpResponse {
    use sqlx::Acquire;

    // Clear the cache for the key "users_with_roles"
    {
        let mut cache = app_state.user_with_roles_cache.lock().unwrap();
        cache.pop("users_with_roles");
    }

    let username: &String = &*path;
    let mut conn = match app_state.pool.acquire().await {
        Ok(conn) => conn,
        Err(e) => {
            eprintln!("Failed to acquire database connection: {:?}", e);
            return HttpResponse::InternalServerError().body("Internal Server Error, Failed to acquire database connection");
        }
    };

    let mut transaction = match conn.begin().await {
        Ok(transaction) => transaction,
        Err(e) => {
            eprintln!("Failed to begin database transaction: {:?}", e);
            return HttpResponse::InternalServerError().body("Internal Server Error, Failed to begin database transaction");
        }
    };

    // Clear all existing roles of the targeted user
    let delete_result = sqlx::query("DELETE FROM AspNetUserRoles WHERE UserId = (SELECT Id FROM AspNetUsers WHERE UserName = ?)")
        .bind(username)
        .execute(&mut *transaction)
        .await;

    if let Err(e) = delete_result {
        eprintln!("Failed to execute DELETE query: {:?}", e);
        return HttpResponse::InternalServerError().body("Internal Server Error, Failed to execute DELETE query");
    }

    // // Prepare the SQL query outside the loop
    // let base_query = sqlx::query(
    //     "INSERT INTO AspNetUserRoles (UserId, RoleId) VALUES ((SELECT Id FROM AspNetUsers WHERE UserName = ?), ?)",
    // )
    // .bind(username.clone());

    // Add or remove roles based on IsInRole field in the request body
    for user_role in &data.into_inner() {
        let insert_query = sqlx::query("INSERT INTO AspNetUserRoles (UserId, RoleId) VALUES ((SELECT Id FROM AspNetUsers WHERE UserName = ?), ?)",
        )
        .bind(username.clone());

        let insert_result = if user_role.IsInRole {
            // If IsInRole is true, execute the query
            insert_query
                .bind(&user_role.RoleId)
                .execute(&mut *transaction)
                .await
                .map(|_| ()) // Map the Ok(_) variant to Ok(())
        } else {
            // If IsInRole is false, skip the query
            Ok(())
        };

        if let Err(e) = insert_result {
            eprintln!("Failed to execute INSERT query: {:?}", e);
            return HttpResponse::InternalServerError().body("Internal Server Error, Failed to execute INSERT query");
        }
    }

    if let Err(e) = transaction.commit().await {
        eprintln!("Failed to commit transaction: {:?}", e);
        return HttpResponse::InternalServerError().body("Internal Server Error, Failed to commit transaction");
    }

    // Respond with JSON body
    HttpResponse::Ok().json(&serde_json::json!({
        "message": "Updated the user roles.",
    }))
}
// Example request body:
/// 1. Assign all roles to targeted user:
// [
//     {
//         "RoleId": "094cd9d9-7388-4804-89d9-f748148bce17",
//         "RoleName": "admin",
//         "IsInRole": true
//     },
//     {
//         "RoleId": "1489bac2-c238-49d0-92c0-6f0049301157",
//         "RoleName": "staff",
//         "IsInRole": true
//     }
// ]

/// 2. Remove all roles from targeted user:
// [
//     {
//         "RoleId": "094cd9d9-7388-4804-89d9-f748148bce17",
//         "RoleName": "admin",
//         "IsInRole": false
//     },
//     {
//         "RoleId": "1489bac2-c238-49d0-92c0-6f0049301157",
//         "RoleName": "staff",
//         "IsInRole": false
//     }
// ]



async fn add_role(
    app_state: web::types::State<AppState>,
    new_role: web::types::Json<NewAspNetRole>,
) -> HttpResponse {
    let name = &new_role.name;
    if name.is_empty() {
        return HttpResponse::BadRequest().body("Role name cannot be empty");
    }

    // Clear the cache for the key "roles"
    {
        let mut cache = app_state.role_cache.lock().unwrap();
        cache.pop("roles");
    }

    let insert_result: Result<_, _> = sqlx::query!(
        "INSERT INTO AspNetRoles (Id, Name, NormalizedName) VALUES (?, ?, ?)",
        uuid::Uuid::new_v4().to_string(),
        name.clone(),
        name.to_uppercase()
    )
    .execute(&app_state.pool)
    .await;

    match insert_result {
        Ok(_) => HttpResponse::Created().body("Role added successfully"),
        Err(e) => {
            eprintln!("Error adding role: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}


async fn update_role(
    app_state: web::types::State<AppState>,
    path: web::types::Path<String>,
    new_role: web::types::Json<NewAspNetRole>,
) -> HttpResponse {
    let role_id: &String = &*path;
    let name: &String = &new_role.name;

    // Clear the cache for the key "roles"
    {
        let mut cache = app_state.role_cache.lock().unwrap();
        cache.pop("roles");
    }

    let update_result: Result<_, sqlx::Error> = sqlx::query!(
        "UPDATE AspNetRoles SET Name = ?, NormalizedName = ? WHERE Id = ?",
        name,
        name.to_uppercase(),
        role_id
    )
    .execute(&app_state.pool)
    .await;

    match update_result {
        Ok(_) => {
            HttpResponse::Ok().body("Role updated successfully")
        }
        Err(sqlx::Error::RowNotFound) => {
            HttpResponse::NotFound().body("Role not found")
        }
        Err(e) => {
            eprintln!("Error updating role: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}

async fn delete_role(
    app_state: web::types::State<AppState>,
    path: web::types::Path<String>,
) -> HttpResponse {
    let role_id = &*path;

    // Clear the cache for the key "roles"
    {
        let mut cache = app_state.role_cache.lock().unwrap();
        cache.pop("roles");
    }

    let delete_result: Result<_, sqlx::Error> = sqlx::query!("DELETE FROM AspNetRoles WHERE Id = ?", role_id)
        .execute(&app_state.pool)
        .await;

    match delete_result {
        Ok(_) => {
            HttpResponse::Ok().body("Role deleted successfully")
        }
        Err(sqlx::Error::RowNotFound) => {
            HttpResponse::NotFound().body("Role not found")
        }
        Err(e) => {
            eprintln!("Error deleting role: {:?}", e);
            HttpResponse::InternalServerError().body("Internal Server Error")
        }
    }
}
*/
