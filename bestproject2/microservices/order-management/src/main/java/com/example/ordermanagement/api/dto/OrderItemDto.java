package com.example.ordermanagement.api.dto;

import jakarta.validation.constraints.NotNull;
import java.math.BigDecimal;

public class OrderItemDto {
  @NotNull
  private Long productId;
  @NotNull
  private Integer quantity;
  @NotNull
  private BigDecimal unitPrice;
  public Long getProductId() { return productId; }
  public void setProductId(Long productId) { this.productId = productId; }
  public Integer getQuantity() { return quantity; }
  public void setQuantity(Integer quantity) { this.quantity = quantity; }
  public BigDecimal getUnitPrice() { return unitPrice; }
  public void setUnitPrice(BigDecimal unitPrice) { this.unitPrice = unitPrice; }
}
