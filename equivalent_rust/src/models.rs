pub mod models {
    use serde::{Deserialize, Serialize};
    // use utoipa::ToSchema;
    // use chrono::prelude::*;
    // use chrono::{DateTime, Utc};
    // use sqlx::prelude::Type;
    // use sqlx::MySql;

    // Define models here (all must be public to allow main to access)

    // #[derive(Debug, sqlx::FromRow, Serialize, Deserialize, ToSchema, Clone)] // Clone for cache putting
    // #[allow(non_snake_case)]
    // pub struct AspNetRole {
    //     pub Id: String,
    //     pub Name: String,
    //     pub NormalizedName: String,
    // }

    // Define the Order model
    #[derive(Debug, sqlx::FromRow, sqlx::types::Type, Serialize, Deserialize)]
    #[allow(non_snake_case)]
    pub struct Order {
        pub order_id: i64,
        // pub order_date: Option<DateTime<Utc>>,
        // pub order_date: String  // => err: the trait `From<std::option::Option<std::string::String>>` is not implemented for `std::string::String`
        pub order_date: Option<String>,
        pub customer_id: i64,
        // pub customer: Customer,
        // pub order_items: Vec<OrderItem>,
        // pub total: f64,
    }

    // Define the Customer model
    #[derive(Debug, Serialize, Deserialize)]
    pub struct Customer {
        pub customer_id: i64,
        pub name: String,
        pub email: String,
    }

    /*
    // Define the OrderItem model
    #[derive(Debug, Serialize, Deserialize)]
    pub struct OrderItem {
        order_item_id: i32,
        order_id: i32,
        product_id: i32,
        quantity: i32,
        product: Product,
    }

    // Define the Product model
    #[derive(Debug, Serialize, Deserialize)]
    pub struct Product {
        product_id: i32,
        name: String,
        price: f64,
    }

    // Define the ViewModel for simplified order representation
    #[derive(Debug, Serialize, Deserialize)]
    pub struct OrderViewModel {
        order_id: i32,
        // order_date: DateTime<Utc>,
        // order_date: String,
        customer_id: i32,
        total: f64,
        customer: Customer,
        order_items: Vec<OrderItem>,
    }
    */

    #[derive(Debug, Serialize, sqlx::FromRow, sqlx::types::Type, sqlx::Decode)]
    pub struct JoinedOrderResponse {
        pub order_id: i64,
        pub order_date: String,
        pub customer: Customer,
        pub order_items: Vec<OrderItem>,
    }

    #[derive(Debug, Serialize, Deserialize)]
    pub struct OrderViewModel {
        pub order_id: i64,
        pub order_date: String,
        pub customer_id: i64,
        pub customer: Customer,
        pub order_items: Vec<OrderItemViewModel>,
        pub total: f64,
    }
    #[derive(Debug, Serialize, Deserialize, Clone)]
    pub struct OrderItemViewModel {
        pub order_item_id: i64,
        pub product_id: i64,
        pub quantity: i64,
        pub product: Product,
    }

    #[derive(Debug, Serialize, Deserialize)]
    pub struct OrderData {
        pub order_id: i64,
        pub order_date: String,
        pub customer_id: i64,
        pub customer_name: String,
        pub customer_email: String,
        pub order_item_id: i64,
        pub product_id: i64,
        pub quantity: i64,
        pub product_name: String,
        pub product_price: f64,
    }

    #[derive(Debug, Serialize, sqlx::FromRow, sqlx::types::Type, sqlx::Decode)]
    pub struct OrderItem {
        pub order_item_id: i64,
        pub product_id: i64,
        pub quantity: i64,
        pub product: Product,
    }

    #[derive(Debug, Serialize, Deserialize, sqlx::FromRow, sqlx::Type, sqlx::Decode, Clone)]
    pub struct Product {
        pub product_id: i64,
        pub name: String,
        pub price: f64,
    }

    /* failed
    // Note:    DataType is not directly accessible in the sqlx module. Instead, you should use sqlx::types::Type to represent data types.
    //          The TypeInfo trait cannot be used as a trait object directly. You need to implement it for a concrete type.
    // Implement Type and Decode traits for Vec<OrderItem>
    impl Type<MySql> for Vec<OrderItem> {
        fn type_info() -> dyn sqlx::TypeInfo {
            sqlx::TypeInfo::with_name("BLOB")
        }

        fn compatible(ty: &dyn sqlx::TypeInfo) -> bool {
            ty.name() == "BLOB"
        }
    }

    impl<'r> sqlx::Decode<'r, MySql> for Vec<OrderItem> {
        fn decode(
            value: sqlx::mysql::MySqlValueRef<'r>,
        ) -> Result<Self, Box<dyn std::error::Error + 'static + Send + Sync>> {
            let bytes = <Vec<u8> as sqlx::Decode<MySql>>::decode(value)?;
            bincode::deserialize(&bytes).map_err(Into::into)
        }
    }
    */

}
