package com.example.ordermanagement.api.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import java.util.List;

public class CreateOrderRequest {
  @Email
  @NotNull
  private String customerEmail;

  @NotNull
  @Size(min = 1)
  private List<OrderItemDto> items;

  public String getCustomerEmail() { return customerEmail; }
  public void setCustomerEmail(String customerEmail) { this.customerEmail = customerEmail; }
  public List<OrderItemDto> getItems() { return items; }
  public void setItems(List<OrderItemDto> items) { this.items = items; }
}
