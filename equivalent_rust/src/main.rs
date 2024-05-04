use ntex::web::{self, App};
// use serde::{Deserialize, Serialize};
use sqlx::mysql::{MySqlPool, MySqlPoolOptions};

mod models;
mod order_handler;


const HOST: &str = "127.0.0.1";
const PORT: &str = "4000";


#[derive(Clone)]
struct AppState {
    pool: MySqlPool,

    // Add an LRU cache
    // role_cache: std::sync::Arc<std::sync::Mutex<LruCache<String, Vec<AspNetRole>>>>,
    // user_cache: std::sync::Arc<std::sync::Mutex<LruCache<String, Vec<AspNetUser>>>>,
    // user_with_roles_cache: std::sync::Arc<std::sync::Mutex<LruCache<String, Vec<AspNetUserWithRoles>>>>,
}


// cargo sqlx prepare --database-url mysql://user:password@127.0.0.1:3306/gorm_learning
#[ntex::main]
async fn main() -> std::io::Result<()> {

    // let _database_url: String = env::var("DATABASE_URL").unwrap();
    const DATABASE_URL: &str = "mysql://user:password@127.0.0.1:3306/gorm_learning"; // "mysql://user:password@127.0.0.1:3306/actix_sqlx"
    
    // (FAILED) Attempt to read the DATABASE_URL environment variable
    #[allow(non_snake_case)]
    // let DATABASE_URL = std::env::var("DATABASE_URL").expect("DATABASE_URL is not set in the .env file");

    const MAX_DB_RETRIES: u32 = 5; // Maximum number of connection retries
    const RETRY_INTERVAL_SECS: u64 = 5; // Interval between retries in seconds

   
    // Log that the API is starting to establish a database connection
    println!("⌛️ Starting Server, establishing database connection...");

    let mut retries: u32 = 0;

    while retries < MAX_DB_RETRIES {
        // create connection pool
        match MySqlPoolOptions::new()
            .max_connections(10)
            .connect(&DATABASE_URL)
            .await
        {
            Ok(pool) => {

                // Create an LRU cache with a maximum capacity of 100 items
                // let role_cache: std::sync::Arc<std::sync::Mutex<LruCache<String, Vec<AspNetRole>>>> 
                //     = std::sync::Arc::new(std::sync::Mutex::new(LruCache::new(std::num::NonZeroUsize::new(100).unwrap())));

                // let user_cache: std::sync::Arc<std::sync::Mutex<LruCache<String, Vec<AspNetUser>>>> 
                //     = std::sync::Arc::new(std::sync::Mutex::new(LruCache::new(std::num::NonZeroUsize::new(100).unwrap())));

                // let user_with_roles_cache 
                //     = std::sync::Arc::new(std::sync::Mutex::new(LruCache::new(std::num::NonZeroUsize::new(100).unwrap())));

                // let app_state: AppState = AppState { pool, role_cache, user_cache, user_with_roles_cache };
                let app_state: AppState = AppState { pool };


                // Start the Actix server with the established database connection
                println!("✅ Database connection established successful! Starting Server...");
                
                // ntex tokio server

                let server = web::server(move || {
                    App::new()

                        .state(app_state.clone())

                        // .route("/", web::get().to(root))
                        // .service(web::resource("/").to(get_pool_info))

                        // use scope to wrap modules:
                        .service(web::scope("/api")
                            // .configure(openapi::configure)
                            .configure(order_handler::configure)
                        )
                })
                .workers(4)
                .bind(format!("{}:{}", HOST, PORT));

                match server {
                    Ok(server) => {
                        // Print the success message after the server starts
                        println!("🚀 Server is up and listening at http://{}:{}/api/explorer/", HOST, PORT);

                        // Start the server
                        if let Err(e) = server.run().await {
                            println!("❌ Server error: {}", e);
                        }
                    }
                    Err(e) => {
                        println!("❌ Failed to bind server: {}", e);
                        return Err(e);
                    }
                }

                return Ok(());
                
            }
            Err(e) => {
                // Log the error and wait before retrying
                eprintln!("❌ Failed to connect to the database: {}", e);
                retries += 1;

                if retries < MAX_DB_RETRIES {
                    println!("⌛️ Retrying in {} seconds...", RETRY_INTERVAL_SECS);
                    std::thread::sleep(std::time::Duration::from_secs(RETRY_INTERVAL_SECS));
                } else {
                    eprintln!("❌ Max connection retries reached. Exiting...");
                    return Err(std::io::Error::new(
                        std::io::ErrorKind::Other,
                        "❌ Failed to connect to the database",
                    ));
                }
            }
        }
    }

    Ok(())
}