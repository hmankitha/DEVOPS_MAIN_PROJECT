from pydantic import BaseModel, Field, validator
from typing import Optional, List
from datetime import datetime

class ProductBase(BaseModel):
    sku: str
    name: str
    description: Optional[str] = None
    short_description: Optional[str] = None
    category_id: str
    brand: Optional[str] = None
    price: float = Field(..., gt=0)
    compare_at_price: Optional[float] = None
    cost_price: Optional[float] = None
    currency: str = "USD"
    is_active: bool = True
    is_featured: bool = False
    weight: Optional[float] = None
    weight_unit: str = "kg"
    dimensions: Optional[dict] = None
    images: Optional[List[str]] = []
    tags: Optional[List[str]] = []
    metadata: Optional[dict] = None
    seo_title: Optional[str] = None
    seo_description: Optional[str] = None
    seo_keywords: Optional[str] = None

class ProductCreate(ProductBase):
    pass

class ProductUpdate(BaseModel):
    name: Optional[str] = None
    description: Optional[str] = None
    short_description: Optional[str] = None
    category_id: Optional[str] = None
    brand: Optional[str] = None
    price: Optional[float] = None
    compare_at_price: Optional[float] = None
    is_active: Optional[bool] = None
    is_featured: Optional[bool] = None
    images: Optional[List[str]] = None
    tags: Optional[List[str]] = None

class ProductResponse(ProductBase):
    id: str
    rating_average: float
    rating_count: int
    view_count: int
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True

class CategoryBase(BaseModel):
    name: str
    slug: str
    description: Optional[str] = None
    parent_id: Optional[str] = None
    image_url: Optional[str] = None
    is_active: bool = True
    sort_order: int = 0
    metadata: Optional[dict] = None

class CategoryCreate(CategoryBase):
    pass

class CategoryUpdate(BaseModel):
    name: Optional[str] = None
    description: Optional[str] = None
    is_active: Optional[bool] = None
    image_url: Optional[str] = None
    sort_order: Optional[int] = None

class CategoryResponse(CategoryBase):
    id: str
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True

class InventoryBase(BaseModel):
    product_id: str
    quantity: int = Field(..., ge=0)
    reserved_quantity: int = Field(default=0, ge=0)
    reorder_point: int = Field(default=10, ge=0)
    reorder_quantity: int = Field(default=50, ge=0)
    warehouse_location: Optional[str] = None

class InventoryCreate(InventoryBase):
    pass

class InventoryUpdate(BaseModel):
    quantity: Optional[int] = Field(None, ge=0)
    reserved_quantity: Optional[int] = Field(None, ge=0)
    reorder_point: Optional[int] = Field(None, ge=0)
    reorder_quantity: Optional[int] = Field(None, ge=0)
    warehouse_location: Optional[str] = None

class InventoryResponse(InventoryBase):
    id: str
    available_quantity: int
    last_restocked_at: Optional[datetime] = None
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True

class ReviewBase(BaseModel):
    product_id: str
    user_id: str
    rating: int = Field(..., ge=1, le=5)
    title: Optional[str] = None
    comment: Optional[str] = None
    images: Optional[List[str]] = []

class ReviewCreate(ReviewBase):
    pass

class ReviewResponse(ReviewBase):
    id: str
    is_verified_purchase: bool
    is_approved: bool
    helpful_count: int
    created_at: datetime
    
    class Config:
        from_attributes = True
