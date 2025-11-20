from locust import HttpUser, task, between, events
import json
import random

class UserManagementUser(HttpUser):
    wait_time = between(1, 3)
    
    def on_start(self):
        """Register and login user"""
        self.register()
        self.login()
    
    def register(self):
        """Register a new user"""
        user_id = random.randint(1000, 999999)
        self.user_data = {
            "email": f"testuser{user_id}@example.com",
            "username": f"testuser{user_id}",
            "password": "TestPassword123!",
            "first_name": "Test",
            "last_name": "User",
            "phone": f"+1555{random.randint(1000000, 9999999)}"
        }
        
        response = self.client.post(
            "/api/v1/auth/register",
            json=self.user_data,
            name="/auth/register"
        )
        
        if response.status_code == 201:
            print(f"Registered user: {self.user_data['email']}")
    
    def login(self):
        """Login user"""
        response = self.client.post(
            "/api/v1/auth/login",
            json={
                "email_or_username": self.user_data["email"],
                "password": self.user_data["password"]
            },
            name="/auth/login"
        )
        
        if response.status_code == 200:
            data = response.json()
            self.token = data.get("data", {}).get("access_token")
            self.headers = {"Authorization": f"Bearer {self.token}"}
            print(f"Logged in user: {self.user_data['email']}")
    
    @task(3)
    def get_profile(self):
        """Get user profile"""
        if hasattr(self, 'headers'):
            self.client.get(
                "/api/v1/users/me",
                headers=self.headers,
                name="/users/me"
            )
    
    @task(2)
    def update_profile(self):
        """Update user profile"""
        if hasattr(self, 'headers'):
            self.client.put(
                "/api/v1/users/me",
                headers=self.headers,
                json={
                    "first_name": "Updated",
                    "last_name": "Name"
                },
                name="/users/me [PUT]"
            )

class ProductCatalogUser(HttpUser):
    wait_time = between(1, 2)
    
    @task(5)
    def list_products(self):
        """List products"""
        params = {
            "limit": random.randint(10, 50),
            "skip": random.randint(0, 100)
        }
        self.client.get(
            "/api/v1/products",
            params=params,
            name="/products"
        )
    
    @task(3)
    def get_product(self):
        """Get product details"""
        product_id = random.randint(1, 1000)
        self.client.get(
            f"/api/v1/products/{product_id}",
            name="/products/{id}"
        )
    
    @task(2)
    def search_products(self):
        """Search products"""
        search_terms = ["laptop", "phone", "tablet", "headphones", "watch"]
        self.client.get(
            "/api/v1/products",
            params={"search": random.choice(search_terms)},
            name="/products [search]"
        )
    
    @task(2)
    def list_categories(self):
        """List categories"""
        self.client.get(
            "/api/v1/categories",
            name="/categories"
        )
    
    @task(1)
    def get_inventory(self):
        """Check inventory"""
        product_id = random.randint(1, 1000)
        self.client.get(
            f"/api/v1/inventory/{product_id}",
            name="/inventory/{product_id}"
        )

class OrderManagementUser(HttpUser):
    wait_time = between(2, 5)
    
    def on_start(self):
        """Login before making orders"""
        # Simulated login
        self.token = "test-token"
        self.headers = {"Authorization": f"Bearer {self.token}"}
    
    @task(5)
    def list_orders(self):
        """List user orders"""
        self.client.get(
            "/api/v1/orders",
            headers=self.headers,
            name="/orders"
        )
    
    @task(2)
    def create_order(self):
        """Create a new order"""
        order_data = {
            "items": [
                {
                    "product_id": str(random.randint(1, 100)),
                    "quantity": random.randint(1, 5),
                    "price": round(random.uniform(10, 1000), 2)
                }
                for _ in range(random.randint(1, 5))
            ],
            "shipping_address": {
                "street": "123 Main St",
                "city": "New York",
                "state": "NY",
                "zip": "10001",
                "country": "US"
            }
        }
        
        self.client.post(
            "/api/v1/orders",
            headers=self.headers,
            json=order_data,
            name="/orders [POST]"
        )
    
    @task(3)
    def get_order(self):
        """Get order details"""
        order_id = random.randint(1, 1000)
        self.client.get(
            f"/api/v1/orders/{order_id}",
            headers=self.headers,
            name="/orders/{id}"
        )

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print("Load test starting...")

@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    print("Load test completed!")
    print(f"Total requests: {environment.stats.total.num_requests}")
    print(f"Total failures: {environment.stats.total.num_failures}")
    print(f"Average response time: {environment.stats.total.avg_response_time}ms")
    print(f"RPS: {environment.stats.total.current_rps}")
