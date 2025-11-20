package com.example.ordermanagement.domain;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "orders")
public class Order {
  @Id
  @GeneratedValue(strategy = GenerationType.IDENTITY)
  private Long id;

  @NotNull
  @Size(min = 3, max = 120)
  private String customerEmail;

  @NotNull
  private Instant createdAt = Instant.now();

  @Enumerated(EnumType.STRING)
  private OrderStatus status = OrderStatus.PENDING;

  @OneToMany(mappedBy = "order", cascade = CascadeType.ALL, orphanRemoval = true)
  private List<OrderItem> items = new ArrayList<>();

  @OneToOne(mappedBy = "order", cascade = CascadeType.ALL, orphanRemoval = true, fetch = FetchType.LAZY)
  private Payment payment;

  @OneToOne(mappedBy = "order", cascade = CascadeType.ALL, orphanRemoval = true, fetch = FetchType.LAZY)
  private Shipment shipment;

  public Long getId() { return id; }
  public String getCustomerEmail() { return customerEmail; }
  public void setCustomerEmail(String customerEmail) { this.customerEmail = customerEmail; }
  public Instant getCreatedAt() { return createdAt; }
  public OrderStatus getStatus() { return status; }
  public void setStatus(OrderStatus status) { this.status = status; }
  public List<OrderItem> getItems() { return items; }
  public void addItem(OrderItem item) { item.setOrder(this); items.add(item); }
  public Payment getPayment() { return payment; }
  public void setPayment(Payment payment) { this.payment = payment; if (payment != null) payment.setOrder(this); }
  public Shipment getShipment() { return shipment; }
  public void setShipment(Shipment shipment) { this.shipment = shipment; if (shipment != null) shipment.setOrder(this); }
}
