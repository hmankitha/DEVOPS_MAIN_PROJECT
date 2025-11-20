from sqlalchemy import Column, String, Integer, Float, Boolean, DateTime, Text, ForeignKey, JSON, Index
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from app.database import Base
import uuid

def generate_uuid():
    return str(uuid.uuid4())

class Product(Base):
    __tablename__ = "products"
    
    id = Column(String, primary_key=True, default=generate_uuid, index=True)
    sku = Column(String(100), unique=True, nullable=False, index=True)
    name = Column(String(255), nullable=False, index=True)
    description = Column(Text)
    short_description = Column(String(500))
    category_id = Column(String, ForeignKey("categories.id"), nullable=False, index=True)
    brand = Column(String(100))
    price = Column(Float, nullable=False)
    compare_at_price = Column(Float)
    cost_price = Column(Float)
    currency = Column(String(3), default="USD")
    is_active = Column(Boolean, default=True, index=True)
    is_featured = Column(Boolean, default=False, index=True)
    weight = Column(Float)
    weight_unit = Column(String(10), default="kg")
    dimensions = Column(JSON)
    images = Column(JSON)
    tags = Column(JSON)
    meta_data = Column(JSON)
    seo_title = Column(String(255))
    seo_description = Column(String(500))
    seo_keywords = Column(String(500))
    rating_average = Column(Float, default=0.0)
    rating_count = Column(Integer, default=0)
    view_count = Column(Integer, default=0)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    deleted_at = Column(DateTime(timezone=True))
    
    category = relationship("Category", back_populates="products")
    inventory = relationship("Inventory", back_populates="product", uselist=False)
    reviews = relationship("Review", back_populates="product")
    
    __table_args__ = (
        Index('idx_product_name_search', 'name'),
        Index('idx_product_price', 'price'),
        Index('idx_product_created', 'created_at'),
    )

class Category(Base):
    __tablename__ = "categories"
    
    id = Column(String, primary_key=True, default=generate_uuid, index=True)
    name = Column(String(100), nullable=False, unique=True, index=True)
    slug = Column(String(100), nullable=False, unique=True, index=True)
    description = Column(Text)
    parent_id = Column(String, ForeignKey("categories.id"), index=True)
    image_url = Column(String(500))
    is_active = Column(Boolean, default=True, index=True)
    sort_order = Column(Integer, default=0)
    meta_data = Column(JSON)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    
    products = relationship("Product", back_populates="category")
    parent = relationship("Category", remote_side=[id], backref="children")

class Inventory(Base):
    __tablename__ = "inventory"
    
    id = Column(String, primary_key=True, default=generate_uuid, index=True)
    product_id = Column(String, ForeignKey("products.id"), nullable=False, unique=True, index=True)
    quantity = Column(Integer, nullable=False, default=0)
    reserved_quantity = Column(Integer, default=0)
    available_quantity = Column(Integer, default=0)
    reorder_point = Column(Integer, default=10)
    reorder_quantity = Column(Integer, default=50)
    warehouse_location = Column(String(100))
    last_restocked_at = Column(DateTime(timezone=True))
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    
    product = relationship("Product", back_populates="inventory")

class Review(Base):
    __tablename__ = "reviews"
    
    id = Column(String, primary_key=True, default=generate_uuid, index=True)
    product_id = Column(String, ForeignKey("products.id"), nullable=False, index=True)
    user_id = Column(String, nullable=False, index=True)
    rating = Column(Integer, nullable=False)
    title = Column(String(200))
    comment = Column(Text)
    is_verified_purchase = Column(Boolean, default=False)
    is_approved = Column(Boolean, default=False, index=True)
    helpful_count = Column(Integer, default=0)
    images = Column(JSON)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    
    product = relationship("Product", back_populates="reviews")
    
    __table_args__ = (
        Index('idx_review_product_rating', 'product_id', 'rating'),
        Index('idx_review_created', 'created_at'),
    )
